package dev

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/logger"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// --- Trigger / Reward type constants ---

const (
	InvitePlanTriggerRedemption   = "redemption"
	InvitePlanTriggerTopup        = "topup"
	InvitePlanTriggerSubscription = "subscription"
)

const (
	InvitePlanRewardQuota        = "quota"
	InvitePlanRewardSubscription = "subscription"
)

// --- Models (compatible with old project SQL schema) ---

// InvitePlan represents an invite reward plan configured by admin.
type InvitePlan struct {
	Id int `json:"id"`

	Name    string `json:"name" gorm:"type:varchar(128);not null"`
	Remark  string `json:"remark" gorm:"type:varchar(255);default:''"`
	Enabled bool   `json:"enabled" gorm:"default:true;index"`

	Priority int `json:"priority" gorm:"type:int;default:0;index"`

	TriggerType               string `json:"trigger_type" gorm:"type:varchar(16);not null;index"`
	TriggerTopupQuota         int64  `json:"trigger_topup_quota" gorm:"type:bigint;not null;default:0"`
	TriggerSubscriptionPlanId int    `json:"trigger_subscription_plan_id" gorm:"type:int;default:0;index"`

	RewardType               string  `json:"reward_type" gorm:"type:varchar(16);not null"`
	RewardPercent            float64 `json:"reward_percent" gorm:"type:decimal(10,4);not null;default:0"`
	RewardMaxQuota           int64   `json:"reward_max_quota" gorm:"type:bigint;not null;default:0"`
	RewardSubscriptionPlanId int     `json:"reward_subscription_plan_id" gorm:"type:int;default:0;index"`

	MaxInviteesPerInviter int   `json:"max_invitees_per_inviter" gorm:"type:int;default:0"`
	EffectiveDays         int   `json:"effective_days" gorm:"type:int;default:0"`
	ExpiresAt             int64 `json:"expires_at" gorm:"bigint;index"`
	CreatedAt             int64 `json:"created_at" gorm:"bigint"`
	UpdatedAt             int64 `json:"updated_at" gorm:"bigint"`
}

// InviteRewardRecord logs each invite reward payout.
type InviteRewardRecord struct {
	Id int `json:"id"`

	PlanId    int `json:"plan_id" gorm:"index;index:idx_invite_reward_plan_inviter_invitee,priority:1"`
	InviterId int `json:"inviter_id" gorm:"index;index:idx_invite_reward_plan_inviter,priority:1;index:idx_invite_reward_plan_inviter_invitee,priority:2"`
	InviteeId int `json:"invitee_id" gorm:"index;index:idx_invite_reward_plan_inviter_invitee,priority:3"`

	SourceType    string `json:"source_type" gorm:"type:varchar(16);not null;uniqueIndex:idx_invite_reward_source,priority:1;index"`
	SourceTradeNo string `json:"source_trade_no" gorm:"type:varchar(255);not null;uniqueIndex:idx_invite_reward_source,priority:2"`

	TriggerType               string `json:"trigger_type" gorm:"type:varchar(16);not null;index"`
	TriggerQuota              int64  `json:"trigger_quota" gorm:"type:bigint;not null;default:0"`
	TriggerSubscriptionPlanId int    `json:"trigger_subscription_plan_id" gorm:"type:int;default:0;index"`

	RewardType               string `json:"reward_type" gorm:"type:varchar(16);not null"`
	RewardQuota              int64  `json:"reward_quota" gorm:"type:bigint;not null;default:0"`
	RewardSubscriptionPlanId int    `json:"reward_subscription_plan_id" gorm:"type:int;default:0;index"`

	CreatedAt int64 `json:"created_at" gorm:"bigint;index"`
}

// --- Sentinel errors ---

var (
	ErrInvitePlanNotFound   = errors.New("invite plan not found")
	ErrInvitePlanHasRecords = errors.New("cannot delete plan with reward records")
)

// --- Normalization helpers ---

func NormalizeInvitePlanTriggerType(value string) string {
	switch strings.TrimSpace(strings.ToLower(value)) {
	case InvitePlanTriggerTopup:
		return InvitePlanTriggerTopup
	case InvitePlanTriggerSubscription:
		return InvitePlanTriggerSubscription
	case InvitePlanTriggerRedemption:
		return InvitePlanTriggerRedemption
	default:
		return ""
	}
}

func NormalizeInvitePlanRewardType(value string) string {
	switch strings.TrimSpace(strings.ToLower(value)) {
	case InvitePlanRewardQuota:
		return InvitePlanRewardQuota
	case InvitePlanRewardSubscription:
		return InvitePlanRewardSubscription
	default:
		return ""
	}
}

// --- GORM hooks ---

func (p *InvitePlan) BeforeCreate(tx *gorm.DB) error {
	now := common.GetTimestamp()
	if p.CreatedAt == 0 {
		p.CreatedAt = now
	}
	p.UpdatedAt = now
	p.TriggerType = NormalizeInvitePlanTriggerType(p.TriggerType)
	p.RewardType = NormalizeInvitePlanRewardType(p.RewardType)
	return nil
}

func (p *InvitePlan) BeforeUpdate(tx *gorm.DB) error {
	p.UpdatedAt = common.GetTimestamp()
	p.TriggerType = NormalizeInvitePlanTriggerType(p.TriggerType)
	p.RewardType = NormalizeInvitePlanRewardType(p.RewardType)
	return nil
}

func (r *InviteRewardRecord) BeforeCreate(tx *gorm.DB) error {
	if r.CreatedAt == 0 {
		r.CreatedAt = common.GetTimestamp()
	}
	r.TriggerType = NormalizeInvitePlanTriggerType(r.TriggerType)
	r.RewardType = NormalizeInvitePlanRewardType(r.RewardType)
	r.SourceType = NormalizeInvitePlanTriggerType(r.SourceType)
	return nil
}

// --- AutoMigrate ---

func AutoMigrateInvitePlans(db *gorm.DB) error {
	if err := migrateInviteRewardRepeatIndex(db); err != nil {
		return err
	}
	return db.AutoMigrate(&InvitePlan{}, &InviteRewardRecord{})
}

func migrateInviteRewardRepeatIndex(db *gorm.DB) error {
	const legacyUniqueIndex = "idx_invite_reward_plan_user"
	if db.Migrator().HasTable(&InviteRewardRecord{}) && db.Migrator().HasIndex(&InviteRewardRecord{}, legacyUniqueIndex) {
		return db.Migrator().DropIndex(&InviteRewardRecord{}, legacyUniqueIndex)
	}
	return nil
}

// --- CRUD ---

func CreateInvitePlan(db *gorm.DB, plan *InvitePlan) error {
	return db.Create(plan).Error
}

func UpdateInvitePlan(db *gorm.DB, id int, updates map[string]interface{}) error {
	result := db.Model(&InvitePlan{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrInvitePlanNotFound
	}
	return nil
}

func DeleteInvitePlan(db *gorm.DB, id int) error {
	var count int64
	if err := db.Model(&InviteRewardRecord{}).Where("plan_id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return ErrInvitePlanHasRecords
	}
	result := db.Delete(&InvitePlan{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrInvitePlanNotFound
	}
	return nil
}

func GetAllInvitePlans(db *gorm.DB) ([]InvitePlan, error) {
	var plans []InvitePlan
	err := db.Order("priority DESC, id ASC").Find(&plans).Error
	return plans, err
}

func GetEnabledInvitePlans(db *gorm.DB) ([]InvitePlan, error) {
	var plans []InvitePlan
	err := db.Where("enabled = ?", true).Order("priority DESC, id ASC").Find(&plans).Error
	return plans, err
}

func GetRewardRecords(db *gorm.DB, planId int, page int, pageSize int) ([]InviteRewardRecord, int64, error) {
	var records []InviteRewardRecord
	var total int64
	query := db.Model(&InviteRewardRecord{}).Where("plan_id = ?", planId)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&records).Error
	return records, total, err
}

// --- Core Reward Logic ---

// HandleInviteRewardForPayment checks invite plans and grants rewards to the
// inviter when the invitee makes a payment (topup or redemption).
//
// Logic:
//   - Plans are checked in priority DESC order, first match wins.
//   - A plan matches when: enabled, trigger_type matches sourceType,
//     trigger_topup_quota <= paymentQuota, plan not expired.
//   - EffectiveDays: if set (> 0), the invitee must have registered within
//     the last EffectiveDays days. After that window, no reward is given.
//   - MaxInviteesPerInviter: if set, caps total reward records for this
//     (plan, inviter) pair.
//   - Reward is rounded from paymentQuota * rewardPercent / 100, capped at
//     RewardMaxQuota.
//
// Called inside the payment DB transaction (tx). Returns rewardQuota > 0
// if a reward was granted.
func HandleInviteRewardForPayment(
	tx *gorm.DB,
	inviteeId int,
	inviterId int,
	inviteeName string,
	inviteeCreatedAt int64,
	sourceType string,
	sourceTradeNo string,
	paymentQuota int64,
	logFunc func(int, string),
) (rewardQuota int64, err error) {
	if inviterId == 0 {
		return 0, nil
	}

	sourceType = NormalizeInvitePlanTriggerType(sourceType)
	if sourceType == "" {
		return 0, nil
	}

	plans, err := GetEnabledInvitePlans(tx)
	if err != nil {
		return 0, fmt.Errorf("get plans: %w", err)
	}

	now := common.GetTimestamp()

	for _, plan := range plans {
		// Trigger type must match
		if plan.TriggerType != sourceType {
			continue
		}

		// Min trigger quota check
		if plan.TriggerTopupQuota > 0 && paymentQuota < plan.TriggerTopupQuota {
			continue
		}

		// Plan expiry check
		if plan.ExpiresAt > 0 && plan.ExpiresAt < now {
			continue
		}

		// EffectiveDays: reward only within N days of invitee registration.
		// 0 means no limit (always reward).
		if plan.EffectiveDays > 0 && inviteeCreatedAt > 0 {
			registrationTime := time.Unix(inviteeCreatedAt, 0)
			if time.Since(registrationTime) > time.Duration(plan.EffectiveDays)*24*time.Hour {
				continue
			}
		}

		// Max invitees per inviter
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

		rewardPercent := plan.RewardPercent
		if rewardPercent <= 0 {
			continue
		}

		reward, quotaErr := common.QuotaFromDecimal64Strict(
			decimal.NewFromInt(paymentQuota).
				Mul(decimal.NewFromFloat(rewardPercent)).
				Div(decimal.NewFromInt(100)),
		)
		if quotaErr != nil {
			return 0, quotaErr
		}
		if plan.RewardMaxQuota > 0 && reward > plan.RewardMaxQuota {
			reward = plan.RewardMaxQuota
		}
		if reward <= 0 {
			continue
		}

		// Grant quota to inviter
		if err := tx.Exec("UPDATE users SET quota = quota + ? WHERE id = ?", reward, inviterId).Error; err != nil {
			return 0, fmt.Errorf("grant quota: %w", err)
		}
		if err := tx.Exec("UPDATE users SET aff_quota = aff_quota + ?, aff_history = aff_history + ? WHERE id = ?",
			reward, reward, inviterId).Error; err != nil {
			common.SysError(fmt.Sprintf("update inviter aff_quota failed: %v", err))
		}

		record := &InviteRewardRecord{
			PlanId:        plan.Id,
			InviterId:     inviterId,
			InviteeId:     inviteeId,
			SourceType:    sourceType,
			SourceTradeNo: sourceTradeNo,
			TriggerType:   sourceType,
			TriggerQuota:  paymentQuota,
			RewardType:    InvitePlanRewardQuota,
			RewardQuota:   reward,
			CreatedAt:     now,
		}
		if err := tx.Create(record).Error; err != nil {
			return 0, fmt.Errorf("create reward record: %w", err)
		}

		if logFunc != nil {
			logFunc(inviterId, fmt.Sprintf("邀请返利到账 %s，活动：%s (被邀请人#%d %s %s)",
				logger.LogQuota(reward), plan.Name, inviteeId, sourceType, logger.LogQuota(paymentQuota)))
		}

		return reward, nil
	}

	return 0, nil
}
