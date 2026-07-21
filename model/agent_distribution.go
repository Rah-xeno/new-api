package model

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/setting/operation_setting"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

const (
	ReferralModeInvite            = "normal_invite"
	ReferralModeAgentDistribution = "agent_distribution"

	AgentCommissionSourceTopup      = "topup"
	AgentCommissionSourceRedemption = "redemption"
	AgentCommissionSourceWithdraw   = "withdraw"
)

var ErrAgentCommissionBalanceInsufficient = errors.New("agent commission balance insufficient")

// AgentRate preserves the old project's decimal(10,4) schema on MySQL and
// PostgreSQL. SQLite does not enforce decimal precision, and the pure-Go
// SQLite migrator cannot re-parse type declarations containing a comma, so a
// comma-free decimal declaration is used for newly created SQLite databases.
type AgentRate float64

func (AgentRate) GormDataType() string {
	return "decimal"
}

func (AgentRate) GormDBDataType(db *gorm.DB, _ *schema.Field) string {
	if db.Dialector.Name() == "sqlite" {
		return "decimal"
	}
	return "decimal(10,4)"
}

type AgentCommissionRecord struct {
	Id int `json:"id"`

	AgentId   int `json:"agent_id" gorm:"index;uniqueIndex:idx_agent_commission_source,priority:1"`
	InviteeId int `json:"invitee_id" gorm:"index"`
	TopUpId   int `json:"top_up_id" gorm:"index"`

	SourceType    string `json:"source_type" gorm:"type:varchar(16);not null;uniqueIndex:idx_agent_commission_source,priority:2;index;index:idx_agent_commission_source_type_created_at,priority:1"`
	SourceTradeNo string `json:"source_trade_no" gorm:"type:varchar(255);not null;uniqueIndex:idx_agent_commission_source,priority:3"`

	SourceAmount     int64     `json:"source_amount" gorm:"type:bigint;not null"`
	CommissionRate   AgentRate `json:"commission_rate" gorm:"not null"`
	CommissionAmount int64     `json:"commission_amount" gorm:"type:bigint;not null"`
	IsFirstTopup     bool      `json:"is_first_topup" gorm:"index"`
	CreatedAt        int64     `json:"created_at" gorm:"type:bigint;index;index:idx_agent_commission_source_type_created_at,priority:2"`
}

type AgentCommissionOutcome struct {
	AgentId           int
	InviteeId         int
	InviteeDisplay    string
	SourceTradeNo     string
	SourceAmount      int64
	CommissionRate    float64
	CommissionAmount  int64
	IsFirstTopupOrder bool
}

type AgentDashboardStats struct {
	InviteeTotal         int64 `json:"invitee_total"`
	PaidInviteeTotal     int64 `json:"paid_invitee_total"`
	TotalTopupAmount     int64 `json:"total_topup_amount"`
	CommissionOrderTotal int64 `json:"commission_order_total"`
}

type AgentInviteeStat struct {
	InviteeId             int    `json:"invitee_id"`
	Username              string `json:"username"`
	DisplayName           string `json:"display_name"`
	Email                 string `json:"email"`
	Status                int    `json:"status"`
	SuccessfulTopupCount  int64  `json:"successful_topup_count"`
	SuccessfulTopupAmount int64  `json:"successful_topup_amount"`
	LastTopupAt           int64  `json:"last_topup_at"`
	LastTopupTradeNo      string `json:"last_topup_trade_no"`
	LastTopupAmount       int64  `json:"last_topup_amount"`
	TotalCommissionAmount int64  `json:"total_commission_amount"`
}

type AgentInviteeTopUp struct {
	TopUpId           int     `json:"top_up_id"`
	InviteeId         int     `json:"invitee_id"`
	Username          string  `json:"username"`
	DisplayName       string  `json:"display_name"`
	Email             string  `json:"email"`
	TradeNo           string  `json:"trade_no"`
	PaymentMethod     string  `json:"payment_method"`
	Money             float64 `json:"money"`
	PaymentAmount     int64   `json:"payment_amount"`
	CreateTime        int64   `json:"create_time"`
	CompleteTime      int64   `json:"complete_time"`
	CommissionAmount  int64   `json:"commission_amount"`
	CommissionRate    float64 `json:"commission_rate"`
	HasCommission     bool    `json:"has_commission"`
	IsFirstTopupOrder bool    `json:"is_first_topup_order"`
}

type AgentCommissionRecordItem struct {
	Id               int     `json:"id"`
	SourceType       string  `json:"source_type"`
	InviteeId        int     `json:"invitee_id"`
	InviteeName      string  `json:"invitee_name"`
	InviteeEmail     string  `json:"invitee_email"`
	SourceTradeNo    string  `json:"source_trade_no"`
	SourceAmount     int64   `json:"source_amount"`
	CommissionRate   float64 `json:"commission_rate"`
	CommissionAmount int64   `json:"commission_amount"`
	IsFirstTopup     bool    `json:"is_first_topup"`
	CreatedAt        int64   `json:"created_at"`
}

type AgentInsightOverview struct {
	OrderTotal           int64   `json:"order_total"`
	SuccessOrderTotal    int64   `json:"success_order_total"`
	PendingOrderTotal    int64   `json:"pending_order_total"`
	FailedOrderTotal     int64   `json:"failed_order_total"`
	SuccessRate          float64 `json:"success_rate"`
	SuccessAmount        int64   `json:"success_amount"`
	CommissionAmount     int64   `json:"commission_amount"`
	CommissionOrderTotal int64   `json:"commission_order_total"`
}

type AgentInsightDay struct {
	Date                 string `json:"date"`
	Label                string `json:"label"`
	SuccessOrderTotal    int64  `json:"success_order_total"`
	SuccessAmount        int64  `json:"success_amount"`
	CommissionAmount     int64  `json:"commission_amount"`
	CommissionOrderTotal int64  `json:"commission_order_total"`
}

type AgentInsights struct {
	RangeDays   int                  `json:"range_days"`
	GeneratedAt int64                `json:"generated_at"`
	Overview    AgentInsightOverview `json:"overview"`
	Daily       []AgentInsightDay    `json:"daily"`
}

func (record *AgentCommissionRecord) BeforeCreate(_ *gorm.DB) error {
	if record.CreatedAt == 0 {
		record.CreatedAt = common.GetTimestamp()
	}
	record.SourceType = strings.TrimSpace(strings.ToLower(record.SourceType))
	return nil
}

func NormalizeReferralMode(value string) string {
	switch strings.TrimSpace(strings.ToLower(value)) {
	case ReferralModeInvite:
		return ReferralModeInvite
	case ReferralModeAgentDistribution:
		return ReferralModeAgentDistribution
	default:
		return ""
	}
}

// normalizeAgentDistributionUsers fills only NULL values introduced when an
// existing new-api users table is upgraded. Databases migrated from the old
// project already contain these columns and values, so their data is left
// untouched. The updates use portable GORM expressions for SQLite, MySQL and
// PostgreSQL.
func normalizeAgentDistributionUsers() error {
	defaults := []struct {
		column string
		value  interface{}
	}{
		{"referral_mode", ""},
		{"agent_enabled", false},
		{"agent_use_default_rates", true},
		{"agent_first_topup_rate", 0},
		{"agent_repeat_topup_rate", 0},
		{"agent_commission_balance", 0},
		{"agent_commission_total", 0},
		{"agent_commission_withdrawn", 0},
		{"first_payment_at", 0},
		{"first_payment_type", ""},
		{"first_payment_trade_no", ""},
	}
	for _, item := range defaults {
		if err := DB.Model(&User{}).
			Where(item.column+" IS NULL").
			UpdateColumn(item.column, item.value).Error; err != nil {
			return fmt.Errorf("normalize agent distribution column %s: %w", item.column, err)
		}
	}
	return nil
}

func resolveReferralModeByInviterTx(tx *gorm.DB, inviterId int) (string, error) {
	if inviterId <= 0 {
		return "", nil
	}
	var inviter User
	if err := tx.Select("id", "agent_enabled").Where("id = ?", inviterId).First(&inviter).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil
		}
		return "", err
	}
	if inviter.AgentEnabled {
		return ReferralModeAgentDistribution, nil
	}
	return ReferralModeInvite, nil
}

func GetUserEffectiveAgentRates(user *User) (float64, float64) {
	setting := operation_setting.GetAgentDistributionSetting()
	if user == nil || user.AgentUseDefaultRates {
		return setting.DefaultFirstTopupRate, setting.DefaultRepeatTopupRate
	}
	return float64(user.AgentFirstTopupRate), float64(user.AgentRepeatTopupRate)
}

func IsAgentPortalVisible(user *User) bool {
	return user != nil && (user.AgentEnabled || user.AgentCommissionBalance > 0 || user.AgentCommissionTotal > 0 || user.AgentCommissionWithdrawn > 0)
}

func yuanAmountToCents(amount float64) int64 {
	if amount <= 0 {
		return 0
	}
	return decimal.NewFromFloat(amount).Mul(decimal.NewFromInt(100)).Round(0).IntPart()
}

func buildAgentInviteeDisplayName(user *User) string {
	if user == nil {
		return "-"
	}
	if value := strings.TrimSpace(user.DisplayName); value != "" {
		return value
	}
	if value := strings.TrimSpace(user.Username); value != "" {
		return value
	}
	return fmt.Sprintf("用户#%d", user.Id)
}

func markFirstAgentPaymentTx(tx *gorm.DB, user *User, sourceType string, tradeNo string) (bool, error) {
	if user.FirstPaymentAt > 0 {
		return false, nil
	}
	now := common.GetTimestamp()
	result := tx.Model(&User{}).
		Where("id = ? AND (first_payment_at = 0 OR first_payment_at IS NULL)", user.Id).
		Updates(map[string]interface{}{
			"first_payment_at":       now,
			"first_payment_type":     sourceType,
			"first_payment_trade_no": tradeNo,
		})
	if result.Error != nil {
		return false, result.Error
	}
	if result.RowsAffected > 0 {
		user.FirstPaymentAt = now
		return true, nil
	}
	return false, nil
}

func handleAgentCommissionTx(tx *gorm.DB, userId int, topUpId int, sourceType string, tradeNo string, sourceAmount int64) (*AgentCommissionOutcome, error) {
	if tx == nil || userId <= 0 || strings.TrimSpace(tradeNo) == "" || sourceAmount <= 0 {
		return nil, nil
	}
	var user User
	if err := lockForUpdate(tx).Where("id = ?", userId).First(&user).Error; err != nil {
		return nil, err
	}
	if user.InviterId <= 0 || user.InviterId == user.Id || NormalizeReferralMode(user.ReferralMode) != ReferralModeAgentDistribution {
		return nil, nil
	}
	var agent User
	if err := lockForUpdate(tx).
		Select("id", "username", "display_name", "agent_enabled", "agent_use_default_rates", "agent_first_topup_rate", "agent_repeat_topup_rate").
		Where("id = ?", user.InviterId).First(&agent).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	if !agent.AgentEnabled {
		return nil, nil
	}
	isFirst, err := markFirstAgentPaymentTx(tx, &user, sourceType, tradeNo)
	if err != nil {
		return nil, err
	}
	firstRate, repeatRate := GetUserEffectiveAgentRates(&agent)
	rate := repeatRate
	if isFirst {
		rate = firstRate
	}
	if rate <= 0 || rate > 100 {
		return nil, nil
	}
	commissionAmount := decimal.NewFromInt(sourceAmount).Mul(decimal.NewFromFloat(rate)).Div(decimal.NewFromInt(100)).Round(0).IntPart()
	if commissionAmount <= 0 {
		return nil, nil
	}
	record := &AgentCommissionRecord{
		AgentId: agent.Id, InviteeId: user.Id, TopUpId: topUpId,
		SourceType: sourceType, SourceTradeNo: tradeNo, SourceAmount: sourceAmount,
		CommissionRate: AgentRate(rate), CommissionAmount: commissionAmount, IsFirstTopup: isFirst,
	}
	result := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(record)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	if err := tx.Model(&User{}).Where("id = ?", agent.Id).Updates(map[string]interface{}{
		"agent_commission_balance": gorm.Expr("agent_commission_balance + ?", commissionAmount),
		"agent_commission_total":   gorm.Expr("agent_commission_total + ?", commissionAmount),
	}).Error; err != nil {
		return nil, err
	}
	return &AgentCommissionOutcome{
		AgentId: agent.Id, InviteeId: user.Id, InviteeDisplay: buildAgentInviteeDisplayName(&user),
		SourceTradeNo: tradeNo, SourceAmount: sourceAmount, CommissionRate: rate,
		CommissionAmount: commissionAmount, IsFirstTopupOrder: isFirst,
	}, nil
}

func HandleAgentCommissionForTopUpTx(tx *gorm.DB, topUp *TopUp) (*AgentCommissionOutcome, error) {
	if topUp == nil {
		return nil, nil
	}
	return handleAgentCommissionTx(tx, topUp.UserId, topUp.Id, AgentCommissionSourceTopup, topUp.TradeNo, yuanAmountToCents(topUp.Money))
}

func HandleAgentCommissionForRedemptionTx(tx *gorm.DB, userId int, tradeNo string, sourceAmountCents int64) (*AgentCommissionOutcome, error) {
	return handleAgentCommissionTx(tx, userId, 0, AgentCommissionSourceRedemption, tradeNo, sourceAmountCents)
}

func ApplyAgentCommissionSideEffects(outcome *AgentCommissionOutcome) {
	if outcome == nil || outcome.AgentId <= 0 || outcome.CommissionAmount <= 0 {
		return
	}
	stage := "复充"
	if outcome.IsFirstTopupOrder {
		stage = "首充"
	}
	RecordLog(outcome.AgentId, LogTypeSystem, fmt.Sprintf("代理分销佣金到账 ¥%.2f，用户 %s %s ¥%.2f，订单号 %s",
		float64(outcome.CommissionAmount)/100, outcome.InviteeDisplay, stage, float64(outcome.SourceAmount)/100, outcome.SourceTradeNo))
}

func WithdrawAgentCommission(agentId int, amountCents int64) (*User, error) {
	if agentId <= 0 || amountCents <= 0 {
		return nil, ErrAgentCommissionBalanceInsufficient
	}
	var updated User
	err := DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&User{}).Where("id = ? AND agent_commission_balance >= ?", agentId, amountCents).Updates(map[string]interface{}{
			"agent_commission_balance":   gorm.Expr("agent_commission_balance - ?", amountCents),
			"agent_commission_withdrawn": gorm.Expr("agent_commission_withdrawn + ?", amountCents),
		})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return ErrAgentCommissionBalanceInsufficient
		}
		record := &AgentCommissionRecord{
			AgentId: agentId, SourceType: AgentCommissionSourceWithdraw,
			SourceTradeNo:    fmt.Sprintf("withdraw-%d-%d-%s", agentId, common.GetTimestamp(), common.GetRandomString(8)),
			CommissionAmount: -amountCents,
		}
		if err := tx.Create(record).Error; err != nil {
			return err
		}
		return tx.Where("id = ?", agentId).First(&updated).Error
	})
	if err != nil {
		return nil, err
	}
	if err := updateUserCache(updated); err != nil {
		common.SysError("update agent cache failed: " + err.Error())
	}
	return &updated, nil
}

func agentInviteesQuery(agentId int, keyword string) *gorm.DB {
	query := DB.Model(&User{}).Where("inviter_id = ? AND referral_mode = ?", agentId, ReferralModeAgentDistribution)
	keyword = strings.TrimSpace(strings.ToLower(keyword))
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("LOWER(username) LIKE ? OR LOWER(display_name) LIKE ? OR LOWER(email) LIKE ?", like, like, like)
	}
	return query
}

func GetAgentDashboardStats(agentId int) (*AgentDashboardStats, error) {
	stats := &AgentDashboardStats{}
	if err := agentInviteesQuery(agentId, "").Count(&stats.InviteeTotal).Error; err != nil {
		return nil, err
	}
	type aggregate struct {
		Amount float64 `gorm:"column:amount"`
		Users  int64   `gorm:"column:users"`
	}
	var row aggregate
	if err := DB.Model(&TopUp{}).Joins("JOIN users ON users.id = top_ups.user_id").
		Where("users.inviter_id = ? AND users.referral_mode = ? AND top_ups.status = ?", agentId, ReferralModeAgentDistribution, common.TopUpStatusSuccess).
		Select("COALESCE(SUM(top_ups.money), 0) AS amount, COUNT(DISTINCT top_ups.user_id) AS users").Scan(&row).Error; err != nil {
		return nil, err
	}
	stats.TotalTopupAmount = yuanAmountToCents(row.Amount)
	stats.PaidInviteeTotal = row.Users
	if err := DB.Model(&AgentCommissionRecord{}).Where("agent_id = ? AND source_type IN ?", agentId, []string{AgentCommissionSourceTopup, AgentCommissionSourceRedemption}).Count(&stats.CommissionOrderTotal).Error; err != nil {
		return nil, err
	}
	return stats, nil
}

func GetAgentInvitees(agentId int, pageInfo *common.PageInfo, keyword string) ([]AgentInviteeStat, int64, error) {
	query := agentInviteesQuery(agentId, keyword)
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	users := make([]User, 0)
	if err := query.Order("id DESC").Limit(pageInfo.GetPageSize()).Offset(pageInfo.GetStartIdx()).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	items := make([]AgentInviteeStat, 0, len(users))
	for _, user := range users {
		item := AgentInviteeStat{InviteeId: user.Id, Username: user.Username, DisplayName: buildAgentInviteeDisplayName(&user), Email: user.Email, Status: user.Status}
		var latest TopUp
		DB.Model(&TopUp{}).Where("user_id = ? AND status = ?", user.Id, common.TopUpStatusSuccess).Count(&item.SuccessfulTopupCount)
		var amount float64
		DB.Model(&TopUp{}).Where("user_id = ? AND status = ?", user.Id, common.TopUpStatusSuccess).Select("COALESCE(SUM(money), 0)").Scan(&amount)
		item.SuccessfulTopupAmount = yuanAmountToCents(amount)
		if DB.Where("user_id = ? AND status = ?", user.Id, common.TopUpStatusSuccess).Order("id DESC").First(&latest).Error == nil {
			item.LastTopupAt, item.LastTopupTradeNo, item.LastTopupAmount = latest.CompleteTime, latest.TradeNo, yuanAmountToCents(latest.Money)
		}
		DB.Model(&AgentCommissionRecord{}).Where("agent_id = ? AND invitee_id = ?", agentId, user.Id).Select("COALESCE(SUM(commission_amount), 0)").Scan(&item.TotalCommissionAmount)
		items = append(items, item)
	}
	return items, total, nil
}

func GetAgentTopUps(agentId int, pageInfo *common.PageInfo, keyword string) ([]AgentInviteeTopUp, int64, error) {
	type row struct {
		TopUpId       int `gorm:"column:top_up_id"`
		InviteeId     int `gorm:"column:invitee_id"`
		Username      string
		DisplayName   string
		Email         string
		TradeNo       string
		PaymentMethod string
		Money         float64
		CreateTime    int64
		CompleteTime  int64
	}
	query := DB.Table("top_ups").Joins("JOIN users ON users.id = top_ups.user_id").
		Where("users.inviter_id = ? AND users.referral_mode = ? AND top_ups.status = ?", agentId, ReferralModeAgentDistribution, common.TopUpStatusSuccess)
	if value := strings.TrimSpace(strings.ToLower(keyword)); value != "" {
		like := "%" + value + "%"
		query = query.Where("LOWER(users.username) LIKE ? OR LOWER(users.display_name) LIKE ? OR LOWER(users.email) LIKE ? OR LOWER(top_ups.trade_no) LIKE ?", like, like, like, like)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	rows := make([]row, 0)
	if err := query.Select("top_ups.id AS top_up_id, users.id AS invitee_id, users.username, users.display_name, users.email, top_ups.trade_no, top_ups.payment_method, top_ups.money, top_ups.create_time, top_ups.complete_time").
		Order("top_ups.id DESC").Limit(pageInfo.GetPageSize()).Offset(pageInfo.GetStartIdx()).Scan(&rows).Error; err != nil {
		return nil, 0, err
	}
	items := make([]AgentInviteeTopUp, 0, len(rows))
	for _, row := range rows {
		item := AgentInviteeTopUp{TopUpId: row.TopUpId, InviteeId: row.InviteeId, Username: row.Username, DisplayName: row.DisplayName, Email: row.Email, TradeNo: row.TradeNo, PaymentMethod: row.PaymentMethod, Money: row.Money, PaymentAmount: yuanAmountToCents(row.Money), CreateTime: row.CreateTime, CompleteTime: row.CompleteTime}
		var record AgentCommissionRecord
		if DB.Where("agent_id = ? AND source_type = ? AND source_trade_no = ?", agentId, AgentCommissionSourceTopup, row.TradeNo).First(&record).Error == nil {
			item.HasCommission, item.CommissionAmount, item.CommissionRate, item.IsFirstTopupOrder = true, record.CommissionAmount, float64(record.CommissionRate), record.IsFirstTopup
		}
		items = append(items, item)
	}
	return items, total, nil
}

func GetAgentCommissionRecords(agentId int, pageInfo *common.PageInfo, keyword string) ([]AgentCommissionRecordItem, int64, error) {
	query := DB.Model(&AgentCommissionRecord{}).Where("agent_id = ?", agentId)
	if value := strings.TrimSpace(strings.ToLower(keyword)); value != "" {
		like := "%" + value + "%"
		query = query.Joins("LEFT JOIN users ON users.id = agent_commission_records.invitee_id").Where("LOWER(users.username) LIKE ? OR LOWER(users.display_name) LIKE ? OR LOWER(users.email) LIKE ? OR LOWER(agent_commission_records.source_trade_no) LIKE ?", like, like, like, like)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	records := make([]AgentCommissionRecord, 0)
	if err := query.Select("agent_commission_records.*").Order("agent_commission_records.id DESC").Limit(pageInfo.GetPageSize()).Offset(pageInfo.GetStartIdx()).Find(&records).Error; err != nil {
		return nil, 0, err
	}
	items := make([]AgentCommissionRecordItem, 0, len(records))
	for _, record := range records {
		item := AgentCommissionRecordItem{Id: record.Id, SourceType: record.SourceType, InviteeId: record.InviteeId, SourceTradeNo: record.SourceTradeNo, SourceAmount: record.SourceAmount, CommissionRate: float64(record.CommissionRate), CommissionAmount: record.CommissionAmount, IsFirstTopup: record.IsFirstTopup, CreatedAt: record.CreatedAt}
		if record.SourceType == AgentCommissionSourceWithdraw {
			item.InviteeName = "提现结算"
		} else {
			var user User
			if DB.Select("id", "username", "display_name", "email").Where("id = ?", record.InviteeId).First(&user).Error == nil {
				item.InviteeName, item.InviteeEmail = buildAgentInviteeDisplayName(&user), user.Email
			}
		}
		items = append(items, item)
	}
	return items, total, nil
}

func GetAgentInsights(agentId int, days int) (*AgentInsights, error) {
	if days <= 0 {
		days = 7
	}
	if days > 90 {
		days = 90
	}
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, -(days - 1))
	result := &AgentInsights{RangeDays: days, GeneratedAt: common.GetTimestamp(), Daily: make([]AgentInsightDay, 0, days)}
	daily := make(map[string]*AgentInsightDay, days)
	for i := 0; i < days; i++ {
		date := start.AddDate(0, 0, i)
		result.Daily = append(result.Daily, AgentInsightDay{Date: date.Format("2006-01-02"), Label: date.Format("01-02")})
		daily[result.Daily[i].Date] = &result.Daily[i]
	}
	topups := make([]TopUp, 0)
	if err := DB.Model(&TopUp{}).Joins("JOIN users ON users.id = top_ups.user_id").Where("users.inviter_id = ? AND users.referral_mode = ? AND top_ups.complete_time >= ?", agentId, ReferralModeAgentDistribution, start.Unix()).Find(&topups).Error; err != nil {
		return nil, err
	}
	for _, topup := range topups {
		result.Overview.OrderTotal++
		switch topup.Status {
		case common.TopUpStatusSuccess:
			result.Overview.SuccessOrderTotal++
			amount := yuanAmountToCents(topup.Money)
			result.Overview.SuccessAmount += amount
			key := time.Unix(topup.CompleteTime, 0).In(now.Location()).Format("2006-01-02")
			if item := daily[key]; item != nil {
				item.SuccessOrderTotal++
				item.SuccessAmount += amount
			}
		case common.TopUpStatusPending:
			result.Overview.PendingOrderTotal++
		default:
			result.Overview.FailedOrderTotal++
		}
	}
	commissions := make([]AgentCommissionRecord, 0)
	if err := DB.Where("agent_id = ? AND source_type IN ? AND created_at >= ?", agentId, []string{AgentCommissionSourceTopup, AgentCommissionSourceRedemption}, start.Unix()).Find(&commissions).Error; err != nil {
		return nil, err
	}
	for _, record := range commissions {
		result.Overview.CommissionOrderTotal++
		result.Overview.CommissionAmount += record.CommissionAmount
		key := time.Unix(record.CreatedAt, 0).In(now.Location()).Format("2006-01-02")
		if item := daily[key]; item != nil {
			item.CommissionOrderTotal++
			item.CommissionAmount += record.CommissionAmount
		}
	}
	if result.Overview.OrderTotal > 0 {
		result.Overview.SuccessRate = float64(result.Overview.SuccessOrderTotal) * 100 / float64(result.Overview.OrderTotal)
	}
	sort.Slice(result.Daily, func(i, j int) bool { return result.Daily[i].Date < result.Daily[j].Date })
	return result, nil
}
