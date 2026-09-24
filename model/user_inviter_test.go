package model

import (
	"errors"
	"strings"
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestBindUserInviter(t *testing.T) {
	for _, agent := range []bool{false, true} {
		mode := ReferralModeInvite
		if agent {
			mode = ReferralModeAgentDistribution
		}
		t.Run(mode, func(t *testing.T) {
			setupUserUpdateTestState(t)
			inviter := User{Username: "inviter", AffCode: "code", Status: common.UserStatusEnabled,
				AgentEnabled: agent, AffCount: 2, AffQuota: 100, AffHistoryQuota: 200, Quota: 300}
			invitee := User{Username: "invitee", AffCode: "own", Quota: 400, CreatedAt: 123,
				FirstPaymentAt: 456, FirstPaymentType: "topup", FirstPaymentTradeNo: "old-payment"}
			require.NoError(t, DB.Create(&inviter).Error)
			require.NoError(t, DB.Create(&invitee).Error)

			bound, err := BindUserInviter(invitee.Id, " code ")
			require.NoError(t, err)
			assert.Equal(t, inviter.Id, bound.InviterId)
			assert.Equal(t, mode, bound.ReferralMode)
			var stored User
			require.NoError(t, DB.First(&stored, invitee.Id).Error)
			assert.Equal(t, inviter.Id, stored.InviterId)
			assert.Equal(t, mode, stored.ReferralMode)
			assert.Equal(t, "own", stored.AffCode)
			assert.Equal(t, invitee.Quota, stored.Quota)
			assert.Equal(t, invitee.CreatedAt, stored.CreatedAt)
			assert.Equal(t, invitee.FirstPaymentAt, stored.FirstPaymentAt)
			assert.Equal(t, invitee.FirstPaymentType, stored.FirstPaymentType)
			assert.Equal(t, invitee.FirstPaymentTradeNo, stored.FirstPaymentTradeNo)

			other := User{Username: "other", AffCode: "other", Status: common.UserStatusEnabled}
			require.NoError(t, DB.Create(&other).Error)
			for _, code := range []string{"code", "other"} {
				_, err = BindUserInviter(invitee.Id, code)
				require.ErrorIs(t, err, ErrUserInviterAlreadyBound)
			}
			require.NoError(t, DB.First(&stored, invitee.Id).Error)
			assert.Equal(t, inviter.Id, stored.InviterId)
			require.NoError(t, DB.First(&inviter, inviter.Id).Error)
			assert.Equal(t, 3, inviter.AffCount)
			assert.Equal(t, int64(100), inviter.AffQuota)
			assert.Equal(t, int64(200), inviter.AffHistoryQuota)
			assert.Equal(t, int64(300), inviter.Quota)
			require.NoError(t, DB.First(&other, other.Id).Error)
			assert.Zero(t, other.AffCount)
		})
	}
}

func TestBindUserInviterRejectsInvalidRelationships(t *testing.T) {
	for _, tc := range []struct {
		name    string
		code    string
		status  int
		deleted bool
		wantErr error
	}{
		{name: "empty", code: "  ", wantErr: ErrInvalidInviteCode},
		{name: "too long", code: strings.Repeat("a", 33), wantErr: ErrInvalidInviteCode},
		{name: "unknown", code: "missing", wantErr: ErrInvalidInviteCode},
		{name: "self", code: "own", wantErr: ErrUserInviterSelf},
		{name: "disabled", code: "code", status: common.UserStatusDisabled, wantErr: ErrInvalidInviteCode},
		{name: "deleted", code: "code", status: common.UserStatusEnabled, deleted: true, wantErr: ErrInvalidInviteCode},
	} {
		t.Run(tc.name, func(t *testing.T) {
			setupUserUpdateTestState(t)
			invitee := User{Username: "invitee", AffCode: "own", Status: common.UserStatusEnabled}
			inviter := User{Username: "inviter", AffCode: "code", Status: tc.status}
			require.NoError(t, DB.Create(&invitee).Error)
			require.NoError(t, DB.Create(&inviter).Error)
			if tc.deleted {
				require.NoError(t, DB.Delete(&inviter).Error)
			}
			_, err := BindUserInviter(invitee.Id, tc.code)
			require.ErrorIs(t, err, tc.wantErr)
			require.NoError(t, DB.First(&invitee, invitee.Id).Error)
			assert.Zero(t, invitee.InviterId)
			assert.Empty(t, invitee.ReferralMode)
			require.NoError(t, DB.Unscoped().First(&inviter, inviter.Id).Error)
			assert.Zero(t, inviter.AffCount)
		})
	}
}

func TestBindUserInviterRollsBackWhenInviteCountFails(t *testing.T) {
	setupUserUpdateTestState(t)
	inviter := User{Username: "inviter", AffCode: "code", Status: common.UserStatusEnabled}
	invitee := User{Username: "invitee", AffCode: "own"}
	require.NoError(t, DB.Create(&inviter).Error)
	require.NoError(t, DB.Create(&invitee).Error)
	errCount := errors.New("invite count update failed")
	require.NoError(t, DB.Callback().Update().Before("gorm:update").Register("test:fail_invite_count", func(tx *gorm.DB) {
		if values, ok := tx.Statement.Dest.(map[string]interface{}); ok {
			if _, updatingCount := values["aff_count"]; updatingCount {
				tx.AddError(errCount)
			}
		}
	}))
	t.Cleanup(func() { require.NoError(t, DB.Callback().Update().Remove("test:fail_invite_count")) })

	_, err := BindUserInviter(invitee.Id, "code")
	require.ErrorIs(t, err, errCount)
	require.NoError(t, DB.First(&invitee, invitee.Id).Error)
	assert.Zero(t, invitee.InviterId)
	assert.Empty(t, invitee.ReferralMode)
	require.NoError(t, DB.First(&inviter, inviter.Id).Error)
	assert.Zero(t, inviter.AffCount)
}

func TestBindUserInviterConcurrentRequestsBindOnlyOnce(t *testing.T) {
	setupUserUpdateTestState(t)
	inviter := User{Username: "inviter", AffCode: "code", Status: common.UserStatusEnabled}
	invitee := User{Username: "invitee", AffCode: "own"}
	require.NoError(t, DB.Create(&inviter).Error)
	require.NoError(t, DB.Create(&invitee).Error)
	start := make(chan struct{})
	results := make(chan error, 2)
	for range 2 {
		go func() {
			<-start
			_, err := BindUserInviter(invitee.Id, "code")
			results <- err
		}()
	}
	close(start)
	first, second := <-results, <-results
	if first != nil {
		first, second = second, first
	}
	require.NoError(t, first)
	require.ErrorIs(t, second, ErrUserInviterAlreadyBound)
	require.NoError(t, DB.First(&inviter, inviter.Id).Error)
	assert.Equal(t, 1, inviter.AffCount)
}
