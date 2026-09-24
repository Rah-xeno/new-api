package model

import (
	"errors"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"gorm.io/gorm"
)

var (
	ErrUserInviterAlreadyBound = errors.New("user already has an inviter")
	ErrInvalidInviteCode       = errors.New("invalid invite code or inviter unavailable")
	ErrUserInviterSelf         = errors.New("user cannot be their own inviter")
)

// BindUserInviter fills a missing referral relationship without replaying
// registration rewards or historical payments.
func BindUserInviter(userId int, affCode string) (*User, error) {
	affCode = strings.TrimSpace(affCode)
	if affCode == "" || len(affCode) > 32 {
		return nil, ErrInvalidInviteCode
	}
	var user User
	err := DB.Transaction(func(tx *gorm.DB) error {
		if err := lockForUpdate(tx).Where("id = ?", userId).First(&user).Error; err != nil {
			return err
		}
		if user.InviterId != 0 {
			return ErrUserInviterAlreadyBound
		}

		var inviter User
		if err := tx.Where("aff_code = ? AND status = ?", affCode, common.UserStatusEnabled).First(&inviter).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrInvalidInviteCode
			}
			return err
		}
		if inviter.Id == user.Id {
			return ErrUserInviterSelf
		}
		user.InviterId = inviter.Id
		user.ReferralMode = ReferralModeInvite
		if inviter.AgentEnabled {
			user.ReferralMode = ReferralModeAgentDistribution
		}

		// The conditional write also protects SQLite, which has no FOR UPDATE.
		result := tx.Model(&User{}).
			Where("id = ? AND (inviter_id = ? OR inviter_id IS NULL)", user.Id, 0).
			Updates(map[string]interface{}{
				"inviter_id":    user.InviterId,
				"referral_mode": user.ReferralMode,
			})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return ErrUserInviterAlreadyBound
		}
		result = tx.Model(&User{}).Where("id = ? AND status = ?", inviter.Id, common.UserStatusEnabled).
			UpdateColumn("aff_count", gorm.Expr("COALESCE(aff_count, 0) + ?", 1))
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return ErrInvalidInviteCode
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &user, nil
}
