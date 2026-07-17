package dev
package dev

import (
	"errors"
	"fmt"
	"time"

	"github.com/QuantumNous/new-api/common"

	"gorm.io/gorm"
)

// InvitePlan represents an invite reward plan configured by admin.
type InvitePlan struct {
	Id                     int            `json:"id" gorm:"primaryKey;autoIncrement"`
	Name                   string         `json:"name" gorm:"type:varchar(255);not null"`
	Remark                 string         `json:"remark" gorm:"type:varchar(255);default:''"`
	Enabled                bool           `json:"enabled" gorm:"default:true"`
	Priority               int            `json:"priority" gorm:"default:0"`
	TriggerType            string         `json:"trigger_type" gorm:"type:varchar(32);not null;default:'topup'"`
	TriggerTopupQuota      int            `json:"trigger_topup_quota" gorm:"default:0"`
	TriggerSubscriptionPlanId int         `json:"trigger_subscription_plan_id" gorm:"default:0"`
	RewardType             string         `json:"reward_type" gorm:"type:varchar(32);not null;default:'quota'"`
	RewardPercent          float64        `json:"reward_percent" gorm:"default:0"`
	RewardMaxQuota         int            `json:"reward_max_quota" gorm:"default:0"`
	RewardSubscriptionPlanId int          `json:"reward_subscription_plan_id" gorm:"default:0"`
	MaxInviteesPerInviter  int            `json:"max_invitees_per_inviter" gorm:"default:0"`
	EffectiveDays          int            `json:"effective_days" gorm:"default:0"`
	ExpiresAt              int64          `json:"expires_at" gorm:"bigint;default:0"`
	CreatedAt              int64          `json:"created_at" gorm:"bigint;autoCreateTime:milli"`
	UpdatedAt              int64          `json:"updated_at" gorm:"bigint;autoUpdateTime:milli"`
	DeletedAt              gorm.DeletedAt `json:"-" gorm:"index"`
}

// InviteRewardRecord logs each invite reward payout.
type InviteRewardRecord struct {
	Id              int    `json:"id" gorm:"primaryKey;autoIncrement"`
	PlanId          int    `json:"plan_id" gorm:"index"`
	PlanName        string `json:"plan_name" gorm:"type:varchar(255)"`
	InviterId       int    `json:"inviter_id" gorm:"index"`
	InviterName     string `json:"inviter_name" gorm:"type:varchar(255)"`
	InviteeId       int    `json:"invitee_id" gorm:"index"`
	InviteeName     string `json:"invitee_name" gorm:"type:varchar(255)"`
	SourceType      string `json:"source_type" gorm:"type:varchar(32)"`   // topup / redemption
	SourceTradeNo   string `json:"source_trade_no" gorm:"type:varchar(255)"` // redemption key or trade no
	TriggerQuota    int    `json:"trigger_quota" gorm:"default:0"`         // the quota that triggered this
	RewardType      string `json:"reward_type" gorm:"type:varchar(32)"`    // quota
	RewardQuota     int    `json:"reward_quota" gorm:"default:0"`
	RewardPercent   float64 `json:"reward_percent" gorm:"default:0"`
	CreatedAt       int64  `json:"created_at" gorm:"bigint;autoCreateTime:milli"`
}

var (
	ErrInvitePlanNotFound     = errors.New("invite plan not found")
	ErrInvitePlanHasRecords   = errors.New("cannot delete plan with reward records")
	ErrInvitePlanExpired      = errors.New("invite plan has expired")
)

// DB is the database connection, set by model package during init.
var DB *gorm.DB

// AutoMigrateInvitePlans creates the invite_plans and invite_reward_records tables.
func AutoMigrateInvitePlans(db *gorm.DB) error {
	return db.AutoMigrate(&InvitePlan{}, &InviteRewardRecord{})
}

// --- InvitePlan CRUD ---

func CreateInvitePlan(plan *InvitePlan) error {
	return DB.Create(plan).Error
}

func UpdateInvitePlan(id int, updates map[string]interface{}) error {
	result := DB.Model(&InvitePlan{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrInvitePlanNotFound
	}
	return nil
}

func DeleteInvitePlan(id int) error {
	var count int64
	if err := DB.Model(&InviteRewardRecord{}).Where("plan_id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return ErrInvitePlanHasRecords
	}
	result := DB.Delete(&InvitePlan{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrInvitePlanNotFound
	}
	return nil
}

func GetInvitePlanByID(id int) (*InvitePlan, error) {
	var plan InvitePlan
	err := DB.First(&plan, id).Error
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

func GetAllInvitePlans() ([]InvitePlan, error) {
	var plans []InvitePlan
	err := DB.Order("priority DESC, id ASC").Find(&plans).Error
	return plans, err
}

func GetEnabledInvitePlans() ([]InvitePlan, error) {
	var plans []InvitePlan
	err := DB.Where("enabled = ?", true).Order("priority DESC, id ASC").Find(&plans).Error
	return plans, err
}

// --- Reward Records ---

func GetRewardRecords(planId int, page int, pageSize int) ([]InviteRewardRecord, int64, error) {
	var records []InviteRewardRecord
	var total int64
	query := model.DB.Model(&InviteRewardRecord{}).Where("plan_id = ?", planId)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&records).Error
	return records, total, err
}

// --- Core Reward Logic ---

// HandleInviteRewardForRedemption checks invite plans and grants rewards
// when an invitee redeems a code. Should be called inside the redemption DB transaction.
// inviteeId: the user who redeemed the code
// inviterId: the user who invited the invitee (0 means no inviter)
// inviteeName: username of the invitee
// redemptionKey: the redemption code key
// redemptionQuota: the quota the invitee received
// logFunc: optional callback for logging (inviterId, message)
// Returns the reward quota granted (0 if none).
func HandleInviteRewardForRedemption(tx *gorm.DB, inviteeId int, inviterId int, inviteeName string, redemptionKey string, redemptionQuota int, logFunc func(int, string)) (rewardQuota int, err error) {
	if inviterId == 0 {
		return 0, nil // no inviter
	}

	// Get enabled plans sorted by priority
	plans, err := GetEnabledInvitePlans()
	if err != nil {
		return 0, fmt.Errorf("get plans: %w", err)
	}

	now := common.GetTimestamp()

	for _, plan := range plans {
		// Skip plans that don't match redemption trigger
		if plan.TriggerType != "redemption" && plan.TriggerType != "topup" {
			continue
		}
		// Check minimum quota threshold
		if plan.TriggerTopupQuota > 0 && redemptionQuota < plan.TriggerTopupQuota {
			continue
		}
		// Check expiry
		if plan.ExpiresAt > 0 && plan.ExpiresAt < now {
			continue
		}
		// Check max invitees per inviter
		if plan.MaxInviteesPerInviter > 0 {
			var count int64
			if err := tx.Model(&InviteRewardRecord{}).
				Where("plan_id = ? AND inviter_id = ?", plan.Id, inviterId).
				Count(&count).Error; err != nil {
				return 0, fmt.Errorf("count inviter records: %w", err)
			}
			if int(count) >= plan.MaxInviteesPerInviter {
				continue
			}
		}

		// Calculate reward
		rewardPercent := plan.RewardPercent
		if rewardPercent <= 0 {
			continue
		}

		baseQuota := float64(redemptionQuota)
		reward := common.QuotaRound(baseQuota * rewardPercent / 100.0)

		// Apply max cap
		if plan.RewardMaxQuota > 0 && reward > plan.RewardMaxQuota {
			reward = plan.RewardMaxQuota
		}
		if reward <= 0 {
			continue
		}

		// Grant quota to inviter via raw SQL to avoid model package dependency
		if err := tx.Exec("UPDATE users SET quota = quota + ? WHERE id = ?", reward, inviterId).Error; err != nil {
			return 0, fmt.Errorf("grant quota: %w", err)
		}

		// Update inviter's AffQuota and AffHistoryQuota
		if err := tx.Exec("UPDATE users SET aff_quota = aff_quota + ?, aff_history = aff_history + ? WHERE id = ?",
			reward, reward, inviterId).Error; err != nil {
			common.SysError(fmt.Sprintf("update inviter aff_quota failed: %v", err))
		}

		// Record the reward
		record := &InviteRewardRecord{
			PlanId:        plan.Id,
			PlanName:      plan.Name,
			InviterId:     inviterId,
			InviteeId:     inviteeId,
			InviteeName:   inviteeName,
			SourceType:    "redemption",
			SourceTradeNo: redemptionKey,
			TriggerQuota:  redemptionQuota,
			RewardType:    "quota",
			RewardQuota:   reward,
			RewardPercent: rewardPercent,
			CreatedAt:     now,
		}
		if err := tx.Create(record).Error; err != nil {
			return 0, fmt.Errorf("create reward record: %w", err)
		}

		// Log
		if logFunc != nil {
			logFunc(inviterId, fmt.Sprintf("邀请返利: 邀请人#%d 获得 %d quota (计划:%s, 被邀请人#%d 兑换 %d)",
				inviterId, reward, plan.Name, inviteeId, redemptionQuota))
		}

		return reward, nil // first matching plan wins
	}

	return 0, nil
}
