package dev

import (
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type legacyInviteRewardRecordIndex struct {
	PlanId    int `gorm:"uniqueIndex:idx_invite_reward_plan_user,priority:1"`
	InviterId int `gorm:"uniqueIndex:idx_invite_reward_plan_user,priority:2"`
	InviteeId int `gorm:"uniqueIndex:idx_invite_reward_plan_user,priority:3"`
}

func (legacyInviteRewardRecordIndex) TableName() string {
	return "invite_reward_records"
}

func TestAutoMigrateInvitePlansReplacesLegacyUniqueIndex(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, AutoMigrateInvitePlans(db))
	require.NoError(t, db.Migrator().CreateIndex(&legacyInviteRewardRecordIndex{}, "idx_invite_reward_plan_user"))
	require.True(t, db.Migrator().HasIndex(&InviteRewardRecord{}, "idx_invite_reward_plan_user"))

	require.NoError(t, migrateInviteRewardRepeatIndex(db))
	assert.False(t, db.Migrator().HasIndex(&InviteRewardRecord{}, "idx_invite_reward_plan_user"))
	assert.True(t, db.Migrator().HasIndex(&InviteRewardRecord{}, "idx_invite_reward_plan_inviter_invitee"))

	records := []InviteRewardRecord{
		{
			PlanId: 1, InviterId: 2, InviteeId: 3,
			SourceType: InvitePlanTriggerTopup, SourceTradeNo: "topup-1",
			TriggerType: InvitePlanTriggerTopup, RewardType: InvitePlanRewardQuota,
		},
		{
			PlanId: 1, InviterId: 2, InviteeId: 3,
			SourceType: InvitePlanTriggerTopup, SourceTradeNo: "topup-2",
			TriggerType: InvitePlanTriggerTopup, RewardType: InvitePlanRewardQuota,
		},
	}
	require.NoError(t, db.Create(&records).Error)
	assert.NotEqual(t, records[0].Id, records[1].Id)
}
