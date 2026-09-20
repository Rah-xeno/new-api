package model

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/dev"
)

type AnalyticsNamedCount struct {
	Name  string `json:"name" gorm:"column:name"`
	Count int64  `json:"count" gorm:"column:total_count"`
}

type AnalyticsChannelRate struct {
	ChannelId    int     `json:"channel_id" gorm:"column:channel_id"`
	ChannelName  string  `json:"channel_name"`
	SuccessCount int64   `json:"success_count"`
	FailureCount int64   `json:"failure_count"`
	SuccessRate  float64 `json:"success_rate"`
	FailureRate  float64 `json:"failure_rate"`
}

type AnalyticsChannelQuota struct {
	ChannelId   int    `json:"channel_id" gorm:"column:channel_id"`
	ChannelName string `json:"channel_name"`
	Quota       int64  `json:"quota" gorm:"column:total_quota"`
}

type AnalyticsChannelLatency struct {
	ChannelId    int     `json:"channel_id"`
	ChannelName  string  `json:"channel_name"`
	RequestCount int64   `json:"request_count"`
	AvgUseTime   float64 `json:"avg_use_time"`
	P95UseTime   int     `json:"p95_use_time"`
	P99UseTime   int     `json:"p99_use_time"`
}

type AnalyticsUserQuota struct {
	UserId   int    `json:"user_id" gorm:"column:user_id"`
	Username string `json:"username"`
	Quota    int64  `json:"quota" gorm:"column:total_quota"`
}

type AnalyticsGroupStat struct {
	GroupName       string  `json:"group_name" gorm:"column:group_name"`
	ConsumeQuota    int64   `json:"consume_quota"`
	TopupMoney      float64 `json:"topup_money"`
	TopupOrderCount int64   `json:"topup_order_count"`
}

type AnalyticsOverallRate struct {
	SuccessCount int64   `json:"success_count"`
	FailureCount int64   `json:"failure_count"`
	SuccessRate  float64 `json:"success_rate"`
	FailureRate  float64 `json:"failure_rate"`
}

type AnalyticsDayQuota struct {
	Day   string `json:"day"`
	Quota int64  `json:"quota"`
}

type AnalyticsDayMoney struct {
	Day   string  `json:"day"`
	Money float64 `json:"money"`
}

type AnalyticsDayUserStat struct {
	Day             string `json:"day"`
	NewUserCount    int64  `json:"new_user_count"`
	PaidUserCount   int64  `json:"paid_user_count"`
	ActiveUserCount int64  `json:"active_user_count"`
}

type AnalyticsReferralOverview struct {
	InviteTopupMoney      float64 `json:"invite_topup_money"`
	InviteOrderCount      int64   `json:"invite_order_count"`
	InvitePaidUserCount   int64   `json:"invite_paid_user_count"`
	InviteRewardQuota     int64   `json:"invite_reward_quota"`
	InviteRewardCount     int64   `json:"invite_reward_count"`
	AgentTopupMoney       float64 `json:"agent_topup_money"`
	AgentOrderCount       int64   `json:"agent_order_count"`
	AgentPaidUserCount    int64   `json:"agent_paid_user_count"`
	AgentCommissionAmount int64   `json:"agent_commission_amount"`
	AgentCommissionCount  int64   `json:"agent_commission_count"`
	AgentNetMoney         float64 `json:"agent_net_money"`
}

type AnalyticsCacheOverview struct {
	PromptTokens          int64   `json:"prompt_tokens"`
	InputTokensTotal      int64   `json:"input_tokens_total"`
	CacheHitTokens        int64   `json:"cache_hit_tokens"`
	CacheCreationTokens   int64   `json:"cache_creation_tokens"`
	CacheCreationTokens5m int64   `json:"cache_creation_tokens_5m"`
	CacheCreationTokens1h int64   `json:"cache_creation_tokens_1h"`
	CacheWriteTokens      int64   `json:"cache_write_tokens"`
	CacheHitRate          float64 `json:"cache_hit_rate"`
	CacheSavedQuota       int64   `json:"cache_saved_quota"`
}

type AdminAnalyticsSummary struct {
	TotalRequestCount  int64   `json:"total_request_count"`
	TotalSuccessCount  int64   `json:"total_success_count"`
	TotalFailureCount  int64   `json:"total_failure_count"`
	TotalConsumeQuota  int64   `json:"total_consume_quota"`
	TotalTopupMoney    float64 `json:"total_topup_money"`
	TotalActiveUsers   int64   `json:"total_active_users"`
	NewUserTotal       int64   `json:"new_user_total"`
	PaidUserTotal      int64   `json:"paid_user_total"`
	CacheHitRate       float64 `json:"cache_hit_rate"`
	CacheSavedQuota    int64   `json:"cache_saved_quota"`
	TotalRedeemedQuota int64   `json:"total_redeemed_quota"`
}

type AdminAnalyticsReport struct {
	StartTimestamp      int64                     `json:"start_timestamp"`
	EndTimestamp        int64                     `json:"end_timestamp"`
	Summary             AdminAnalyticsSummary     `json:"summary"`
	Overall             AnalyticsOverallRate      `json:"overall"`
	ModelRanking        []AnalyticsNamedCount     `json:"model_ranking"`
	ChannelRates        []AnalyticsChannelRate    `json:"channel_rates"`
	ChannelConsumptions []AnalyticsChannelQuota   `json:"channel_consumptions"`
	ChannelLatencies    []AnalyticsChannelLatency `json:"channel_latencies"`
	UserConsumptions    []AnalyticsUserQuota      `json:"user_consumptions"`
	GroupStats          []AnalyticsGroupStat      `json:"group_stats"`
	DailyConsumes       []AnalyticsDayQuota       `json:"daily_consumes"`
	DailyTopups         []AnalyticsDayMoney       `json:"daily_topups"`
	DailyUserStats      []AnalyticsDayUserStat    `json:"daily_user_stats"`
	ReferralOverview    AnalyticsReferralOverview `json:"referral_overview"`
	CacheOverview       AnalyticsCacheOverview    `json:"cache_overview"`
}

type analyticsChannelTypeCount struct {
	ChannelId  int   `gorm:"column:channel_id"`
	LogType    int   `gorm:"column:log_type"`
	TotalCount int64 `gorm:"column:total_count"`
}

type analyticsChannelLatencyBucket struct {
	ChannelId  int   `gorm:"column:channel_id"`
	UseTime    int   `gorm:"column:use_time"`
	TotalCount int64 `gorm:"column:total_count"`
}

type analyticsDayCount struct {
	DayIndex   int64 `gorm:"column:day_index"`
	TotalCount int64 `gorm:"column:total_count"`
}

type analyticsDayQuotaRow struct {
	DayIndex int64 `gorm:"column:day_index"`
	Quota    int64 `gorm:"column:total_quota"`
}

type analyticsGroupQuotaRow struct {
	GroupName    string `gorm:"column:group_name"`
	ConsumeQuota int64  `gorm:"column:consume_quota"`
}

type analyticsGroupTopupRow struct {
	GroupName       string  `gorm:"column:group_name"`
	TopupMoney      float64 `gorm:"column:topup_money"`
	TopupOrderCount int64   `gorm:"column:topup_order_count"`
}

type analyticsReferralTopupRow struct {
	ReferralMode  string  `gorm:"column:referral_mode"`
	TopupMoney    float64 `gorm:"column:topup_money"`
	OrderCount    int64   `gorm:"column:order_count"`
	PaidUserCount int64   `gorm:"column:paid_user_count"`
}

type analyticsRewardRow struct {
	RewardQuota int64 `gorm:"column:reward_quota"`
	RecordCount int64 `gorm:"column:record_count"`
}

type analyticsCommissionRow struct {
	CommissionAmount int64 `gorm:"column:commission_amount"`
	RecordCount      int64 `gorm:"column:record_count"`
}

type analyticsCacheOverviewRow struct {
	PromptTokens          float64 `gorm:"column:prompt_tokens"`
	InputTokensTotal      float64 `gorm:"column:input_tokens_total"`
	CacheHitTokens        float64 `gorm:"column:cache_hit_tokens"`
	CacheRateDenominator  float64 `gorm:"column:cache_rate_denominator"`
	CacheCreationTokens   float64 `gorm:"column:cache_creation_tokens"`
	CacheCreationTokens5m float64 `gorm:"column:cache_creation_tokens_5m"`
	CacheCreationTokens1h float64 `gorm:"column:cache_creation_tokens_1h"`
	CacheWriteTokens      float64 `gorm:"column:cache_write_tokens"`
	CacheSavedQuota       float64 `gorm:"column:cache_saved_quota"`
}

func analyticsDayIndexExpr(column string, databaseType common.DatabaseType, timezoneOffsetMinutes int) string {
	offsetSeconds := timezoneOffsetMinutes * 60
	switch databaseType {
	case common.DatabaseTypePostgreSQL:
		return fmt.Sprintf("CAST(FLOOR((%s + %d) / 86400.0) AS BIGINT)", column, offsetSeconds)
	case common.DatabaseTypeMySQL:
		return fmt.Sprintf("CAST(FLOOR((%s + %d) / 86400.0) AS SIGNED)", column, offsetSeconds)
	case common.DatabaseTypeClickHouse:
		return fmt.Sprintf("toInt64(intDiv(%s + %d, 86400))", column, offsetSeconds)
	default:
		return fmt.Sprintf("CAST((%s + %d) / 86400 AS INTEGER)", column, offsetSeconds)
	}
}

func analyticsJSONNumberExpr(column string, key string, databaseType common.DatabaseType) string {
	switch databaseType {
	case common.DatabaseTypePostgreSQL:
		return "CASE WHEN " + column + " IS NULL OR " + column + " = '' THEN 0 ELSE COALESCE((" + column + "::jsonb ->> '" + key + "')::double precision, 0) END"
	case common.DatabaseTypeMySQL:
		return "COALESCE(CAST(JSON_UNQUOTE(JSON_EXTRACT(" + column + ", '$." + key + "')) AS DECIMAL(30,8)), 0)"
	case common.DatabaseTypeClickHouse:
		return "toFloat64OrZero(if(empty(" + column + "), '0', JSONExtractRaw(" + column + ", '" + key + "')))"
	default:
		return "COALESCE(CAST(json_extract(" + column + ", '$." + key + "') AS REAL), 0)"
	}
}

func analyticsJSONStringExpr(column string, key string, databaseType common.DatabaseType) string {
	switch databaseType {
	case common.DatabaseTypePostgreSQL:
		return "CASE WHEN " + column + " IS NULL OR " + column + " = '' THEN '' ELSE COALESCE(" + column + "::jsonb ->> '" + key + "', '') END"
	case common.DatabaseTypeMySQL:
		return "COALESCE(JSON_UNQUOTE(JSON_EXTRACT(" + column + ", '$." + key + "')), '')"
	case common.DatabaseTypeClickHouse:
		return "if(empty(" + column + "), '', JSONExtractString(" + column + ", '" + key + "'))"
	default:
		return "COALESCE(json_extract(" + column + ", '$." + key + "'), '')"
	}
}

func analyticsDayFromIndex(dayIndex int64) string {
	return time.Unix(dayIndex*86400, 0).UTC().Format("2006-01-02")
}

func analyticsRoundPercent(value float64) float64 {
	if value <= 0 || math.IsNaN(value) {
		return 0
	}
	return math.Round(value*100) / 100
}

func analyticsRoundMoney(value float64) float64 {
	return math.Round(value*100) / 100
}

func analyticsRoundInt64(value float64) int64 {
	if math.IsNaN(value) {
		return 0
	}
	rounded := math.Round(value)
	if rounded >= math.MaxInt64 {
		return math.MaxInt64
	}
	if rounded <= math.MinInt64 {
		return math.MinInt64
	}
	return int64(rounded)
}

func analyticsNormalizeGroup(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "default"
	}
	return value
}

func analyticsChannelNames(channelIDs []int) (map[int]string, error) {
	uniqueIDs := make([]int, 0, len(channelIDs))
	seen := make(map[int]struct{}, len(channelIDs))
	for _, channelID := range channelIDs {
		if channelID <= 0 {
			continue
		}
		if _, ok := seen[channelID]; ok {
			continue
		}
		seen[channelID] = struct{}{}
		uniqueIDs = append(uniqueIDs, channelID)
	}

	type channelNameRow struct {
		Id   int    `gorm:"column:id"`
		Name string `gorm:"column:name"`
	}
	rows := make([]channelNameRow, 0, len(uniqueIDs))
	if len(uniqueIDs) > 0 {
		if err := DB.Table("channels").Select("id, name").Where("id IN ?", uniqueIDs).Scan(&rows).Error; err != nil {
			return nil, err
		}
	}

	names := make(map[int]string, len(rows))
	for _, row := range rows {
		names[row.Id] = row.Name
	}
	return names, nil
}

func analyticsChannelName(channelID int, names map[int]string) string {
	if channelID <= 0 {
		return "Unassigned channel"
	}
	if name := strings.TrimSpace(names[channelID]); name != "" {
		return name
	}
	return fmt.Sprintf("Channel #%d", channelID)
}

func analyticsPercentile(buckets []analyticsChannelLatencyBucket, total int64, percentile float64) int {
	if len(buckets) == 0 || total <= 0 {
		return 0
	}
	threshold := int64(math.Ceil(float64(total) * percentile))
	var cumulative int64
	for _, bucket := range buckets {
		cumulative += bucket.TotalCount
		if cumulative >= threshold {
			return bucket.UseTime
		}
	}
	return buckets[len(buckets)-1].UseTime
}

func GetAdminAnalyticsModelRanking(startTimestamp int64, endTimestamp int64, limit int) ([]AnalyticsNamedCount, error) {
	items := make([]AnalyticsNamedCount, 0)
	query := LOG_DB.Table("logs").
		Select("model_name AS name, COUNT(*) AS total_count").
		Where("type = ? AND created_at >= ? AND created_at <= ?", LogTypeConsume, startTimestamp, endTimestamp).
		Where("model_name <> ''").
		Group("model_name").
		Order("total_count DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if err := query.Scan(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func GetAdminAnalyticsChannelRates(startTimestamp int64, endTimestamp int64) ([]AnalyticsChannelRate, AnalyticsOverallRate, error) {
	rows := make([]analyticsChannelTypeCount, 0)
	if err := LOG_DB.Table("logs").
		Select("channel_id, type AS log_type, COUNT(*) AS total_count").
		Where("created_at >= ? AND created_at <= ?", startTimestamp, endTimestamp).
		Where("type IN ?", []int{LogTypeConsume, LogTypeError}).
		Group("channel_id, type").
		Scan(&rows).Error; err != nil {
		return nil, AnalyticsOverallRate{}, err
	}

	channelIDs := make([]int, 0, len(rows))
	statsByChannel := make(map[int]*AnalyticsChannelRate, len(rows))
	overall := AnalyticsOverallRate{}
	for _, row := range rows {
		stats := statsByChannel[row.ChannelId]
		if stats == nil {
			stats = &AnalyticsChannelRate{ChannelId: row.ChannelId}
			statsByChannel[row.ChannelId] = stats
			channelIDs = append(channelIDs, row.ChannelId)
		}
		if row.LogType == LogTypeConsume {
			stats.SuccessCount += row.TotalCount
			overall.SuccessCount += row.TotalCount
		} else {
			stats.FailureCount += row.TotalCount
			overall.FailureCount += row.TotalCount
		}
	}

	names, err := analyticsChannelNames(channelIDs)
	if err != nil {
		return nil, AnalyticsOverallRate{}, err
	}
	result := make([]AnalyticsChannelRate, 0, len(channelIDs))
	for _, channelID := range channelIDs {
		stats := statsByChannel[channelID]
		total := stats.SuccessCount + stats.FailureCount
		if total > 0 {
			stats.SuccessRate = analyticsRoundPercent(float64(stats.SuccessCount) * 100 / float64(total))
			stats.FailureRate = analyticsRoundPercent(float64(stats.FailureCount) * 100 / float64(total))
		}
		stats.ChannelName = analyticsChannelName(channelID, names)
		result = append(result, *stats)
	}
	sort.Slice(result, func(i, j int) bool {
		leftTotal := result[i].SuccessCount + result[i].FailureCount
		rightTotal := result[j].SuccessCount + result[j].FailureCount
		if leftTotal == rightTotal {
			return result[i].ChannelId < result[j].ChannelId
		}
		return leftTotal > rightTotal
	})

	total := overall.SuccessCount + overall.FailureCount
	if total > 0 {
		overall.SuccessRate = analyticsRoundPercent(float64(overall.SuccessCount) * 100 / float64(total))
		overall.FailureRate = analyticsRoundPercent(float64(overall.FailureCount) * 100 / float64(total))
	}
	return result, overall, nil
}

func GetAdminAnalyticsChannelConsumptions(startTimestamp int64, endTimestamp int64) ([]AnalyticsChannelQuota, error) {
	items := make([]AnalyticsChannelQuota, 0)
	if err := LOG_DB.Table("logs").
		Select("channel_id, SUM(quota) AS total_quota").
		Where("type = ? AND created_at >= ? AND created_at <= ?", LogTypeConsume, startTimestamp, endTimestamp).
		Group("channel_id").
		Order("total_quota DESC").
		Scan(&items).Error; err != nil {
		return nil, err
	}

	channelIDs := make([]int, 0, len(items))
	for _, item := range items {
		channelIDs = append(channelIDs, item.ChannelId)
	}
	names, err := analyticsChannelNames(channelIDs)
	if err != nil {
		return nil, err
	}
	for index := range items {
		items[index].ChannelName = analyticsChannelName(items[index].ChannelId, names)
	}
	return items, nil
}

func GetAdminAnalyticsChannelLatencies(startTimestamp int64, endTimestamp int64) ([]AnalyticsChannelLatency, error) {
	rows := make([]analyticsChannelLatencyBucket, 0)
	if err := LOG_DB.Table("logs").
		Select("channel_id, use_time, COUNT(*) AS total_count").
		Where("created_at >= ? AND created_at <= ?", startTimestamp, endTimestamp).
		Where("type IN ?", []int{LogTypeConsume, LogTypeError}).
		Where("use_time > 0").
		Group("channel_id, use_time").
		Order("channel_id ASC, use_time ASC").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	channelIDs := make([]int, 0)
	bucketsByChannel := make(map[int][]analyticsChannelLatencyBucket)
	for _, row := range rows {
		if _, exists := bucketsByChannel[row.ChannelId]; !exists {
			channelIDs = append(channelIDs, row.ChannelId)
		}
		bucketsByChannel[row.ChannelId] = append(bucketsByChannel[row.ChannelId], row)
	}
	names, err := analyticsChannelNames(channelIDs)
	if err != nil {
		return nil, err
	}

	items := make([]AnalyticsChannelLatency, 0, len(channelIDs))
	for _, channelID := range channelIDs {
		buckets := bucketsByChannel[channelID]
		var totalCount int64
		var totalUseTime int64
		for _, bucket := range buckets {
			totalCount += bucket.TotalCount
			totalUseTime += int64(bucket.UseTime) * bucket.TotalCount
		}
		if totalCount == 0 {
			continue
		}
		items = append(items, AnalyticsChannelLatency{
			ChannelId:    channelID,
			ChannelName:  analyticsChannelName(channelID, names),
			RequestCount: totalCount,
			AvgUseTime:   analyticsRoundPercent(float64(totalUseTime) / float64(totalCount)),
			P95UseTime:   analyticsPercentile(buckets, totalCount, 0.95),
			P99UseTime:   analyticsPercentile(buckets, totalCount, 0.99),
		})
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].RequestCount == items[j].RequestCount {
			return items[i].ChannelId < items[j].ChannelId
		}
		return items[i].RequestCount > items[j].RequestCount
	})
	return items, nil
}

func GetAdminAnalyticsUserConsumptions(startTimestamp int64, endTimestamp int64, limit int) ([]AnalyticsUserQuota, error) {
	items := make([]AnalyticsUserQuota, 0)
	query := LOG_DB.Table("logs").
		Select("user_id, username, SUM(quota) AS total_quota").
		Where("type = ? AND created_at >= ? AND created_at <= ?", LogTypeConsume, startTimestamp, endTimestamp).
		Group("user_id, username").
		Order("total_quota DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if err := query.Scan(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func GetAdminAnalyticsGroupStats(startTimestamp int64, endTimestamp int64) ([]AnalyticsGroupStat, error) {
	logGroupExpr := logGroupCol
	quotaRows := make([]analyticsGroupQuotaRow, 0)
	if err := LOG_DB.Table("logs").
		Select(logGroupExpr+" AS group_name, COALESCE(SUM(quota), 0) AS consume_quota").
		Where("type = ? AND created_at >= ? AND created_at <= ?", LogTypeConsume, startTimestamp, endTimestamp).
		Group(logGroupExpr).
		Scan(&quotaRows).Error; err != nil {
		return nil, err
	}

	userGroupExpr := "users." + commonGroupCol
	topupRows := make([]analyticsGroupTopupRow, 0)
	if err := DB.Table("top_ups").
		Joins("JOIN users ON users.id = top_ups.user_id").
		Select(userGroupExpr+" AS group_name, COALESCE(SUM(top_ups.money), 0) AS topup_money, COUNT(*) AS topup_order_count").
		Where("top_ups.status = ? AND top_ups.complete_time >= ? AND top_ups.complete_time <= ?", common.TopUpStatusSuccess, startTimestamp, endTimestamp).
		Group(userGroupExpr).
		Scan(&topupRows).Error; err != nil {
		return nil, err
	}

	statsByGroup := make(map[string]*AnalyticsGroupStat)
	for _, row := range quotaRows {
		groupName := analyticsNormalizeGroup(row.GroupName)
		statsByGroup[groupName] = &AnalyticsGroupStat{
			GroupName:    groupName,
			ConsumeQuota: row.ConsumeQuota,
		}
	}
	for _, row := range topupRows {
		groupName := analyticsNormalizeGroup(row.GroupName)
		stats := statsByGroup[groupName]
		if stats == nil {
			stats = &AnalyticsGroupStat{GroupName: groupName}
			statsByGroup[groupName] = stats
		}
		stats.TopupMoney = analyticsRoundMoney(row.TopupMoney)
		stats.TopupOrderCount = row.TopupOrderCount
	}

	result := make([]AnalyticsGroupStat, 0, len(statsByGroup))
	for _, stats := range statsByGroup {
		result = append(result, *stats)
	}
	sort.Slice(result, func(i, j int) bool {
		leftScore := result[i].TopupMoney*100 + float64(result[i].ConsumeQuota)
		rightScore := result[j].TopupMoney*100 + float64(result[j].ConsumeQuota)
		if leftScore == rightScore {
			return result[i].GroupName < result[j].GroupName
		}
		return leftScore > rightScore
	})
	return result, nil
}

func GetAdminAnalyticsDailyConsumes(startTimestamp int64, endTimestamp int64, timezoneOffsetMinutes int) ([]AnalyticsDayQuota, error) {
	dayExpr := analyticsDayIndexExpr("created_at", common.LogDatabaseType(), timezoneOffsetMinutes)
	rows := make([]analyticsDayQuotaRow, 0)
	if err := LOG_DB.Table("logs").
		Select(dayExpr+" AS day_index, COALESCE(SUM(quota), 0) AS total_quota").
		Where("type = ? AND created_at >= ? AND created_at <= ?", LogTypeConsume, startTimestamp, endTimestamp).
		Group(dayExpr).
		Order("day_index ASC").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	items := make([]AnalyticsDayQuota, 0, len(rows))
	for _, row := range rows {
		items = append(items, AnalyticsDayQuota{Day: analyticsDayFromIndex(row.DayIndex), Quota: row.Quota})
	}
	return items, nil
}

func GetAdminAnalyticsDailyTopups(startTimestamp int64, endTimestamp int64, timezoneOffsetMinutes int) ([]AnalyticsDayMoney, error) {
	type topupRow struct {
		CompleteTime int64   `gorm:"column:complete_time"`
		Money        float64 `gorm:"column:money"`
	}
	rows := make([]topupRow, 0)
	if err := DB.Table("top_ups").
		Select("complete_time, money").
		Where("status = ? AND complete_time >= ? AND complete_time <= ?", common.TopUpStatusSuccess, startTimestamp, endTimestamp).
		Order("complete_time ASC").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	location := time.FixedZone("analytics", timezoneOffsetMinutes*60)
	moneyByDay := make(map[string]float64)
	for _, row := range rows {
		day := time.Unix(row.CompleteTime, 0).In(location).Format("2006-01-02")
		moneyByDay[day] += row.Money
	}
	items := make([]AnalyticsDayMoney, 0, len(moneyByDay))
	for day, money := range moneyByDay {
		items = append(items, AnalyticsDayMoney{Day: day, Money: analyticsRoundMoney(money)})
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].Day < items[j].Day
	})
	return items, nil
}

func GetAdminAnalyticsDailyUserStats(startTimestamp int64, endTimestamp int64, timezoneOffsetMinutes int) ([]AnalyticsDayUserStat, int64, error) {
	mainDayExpr := analyticsDayIndexExpr("created_at", common.MainDatabaseType(), timezoneOffsetMinutes)
	paymentDayExpr := analyticsDayIndexExpr("first_payment_at", common.MainDatabaseType(), timezoneOffsetMinutes)
	logDayExpr := analyticsDayIndexExpr("created_at", common.LogDatabaseType(), timezoneOffsetMinutes)

	newUserRows := make([]analyticsDayCount, 0)
	if err := DB.Model(&User{}).
		Select(mainDayExpr+" AS day_index, COUNT(*) AS total_count").
		Where("created_at >= ? AND created_at <= ?", startTimestamp, endTimestamp).
		Group(mainDayExpr).
		Order("day_index ASC").
		Scan(&newUserRows).Error; err != nil {
		return nil, 0, err
	}

	paidUserRows := make([]analyticsDayCount, 0)
	if err := DB.Model(&User{}).
		Select(paymentDayExpr+" AS day_index, COUNT(*) AS total_count").
		Where("first_payment_at > 0 AND first_payment_at >= ? AND first_payment_at <= ?", startTimestamp, endTimestamp).
		Group(paymentDayExpr).
		Order("day_index ASC").
		Scan(&paidUserRows).Error; err != nil {
		return nil, 0, err
	}

	activeUserRows := make([]analyticsDayCount, 0)
	if err := LOG_DB.Table("logs").
		Select(logDayExpr+" AS day_index, COUNT(DISTINCT user_id) AS total_count").
		Where("created_at >= ? AND created_at <= ?", startTimestamp, endTimestamp).
		Where("type IN ?", []int{LogTypeConsume, LogTypeError}).
		Where("user_id > 0").
		Group(logDayExpr).
		Order("day_index ASC").
		Scan(&activeUserRows).Error; err != nil {
		return nil, 0, err
	}

	var totalActiveUsers int64
	if err := LOG_DB.Table("logs").
		Where("created_at >= ? AND created_at <= ?", startTimestamp, endTimestamp).
		Where("type IN ?", []int{LogTypeConsume, LogTypeError}).
		Where("user_id > 0").
		Distinct("user_id").
		Count(&totalActiveUsers).Error; err != nil {
		return nil, 0, err
	}

	statsByDay := make(map[string]*AnalyticsDayUserStat)
	mergeRows := func(rows []analyticsDayCount, apply func(*AnalyticsDayUserStat, int64)) {
		for _, row := range rows {
			day := analyticsDayFromIndex(row.DayIndex)
			stats := statsByDay[day]
			if stats == nil {
				stats = &AnalyticsDayUserStat{Day: day}
				statsByDay[day] = stats
			}
			apply(stats, row.TotalCount)
		}
	}
	mergeRows(newUserRows, func(stats *AnalyticsDayUserStat, count int64) {
		stats.NewUserCount = count
	})
	mergeRows(paidUserRows, func(stats *AnalyticsDayUserStat, count int64) {
		stats.PaidUserCount = count
	})
	mergeRows(activeUserRows, func(stats *AnalyticsDayUserStat, count int64) {
		stats.ActiveUserCount = count
	})

	items := make([]AnalyticsDayUserStat, 0, len(statsByDay))
	for _, stats := range statsByDay {
		items = append(items, *stats)
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].Day < items[j].Day
	})
	return items, totalActiveUsers, nil
}

func GetAdminAnalyticsReferralOverview(startTimestamp int64, endTimestamp int64) (AnalyticsReferralOverview, error) {
	report := AnalyticsReferralOverview{}
	topupRows := make([]analyticsReferralTopupRow, 0, 2)
	if err := DB.Table("top_ups").
		Joins("JOIN users ON users.id = top_ups.user_id").
		Select("users.referral_mode AS referral_mode, COALESCE(SUM(top_ups.money), 0) AS topup_money, COUNT(*) AS order_count, COUNT(DISTINCT top_ups.user_id) AS paid_user_count").
		Where("top_ups.status = ? AND top_ups.complete_time >= ? AND top_ups.complete_time <= ?", common.TopUpStatusSuccess, startTimestamp, endTimestamp).
		Where("users.inviter_id > 0 AND users.referral_mode IN ?", []string{ReferralModeInvite, ReferralModeAgentDistribution}).
		Group("users.referral_mode").
		Scan(&topupRows).Error; err != nil {
		return report, err
	}
	for _, row := range topupRows {
		if row.ReferralMode == ReferralModeInvite {
			report.InviteTopupMoney = analyticsRoundMoney(row.TopupMoney)
			report.InviteOrderCount = row.OrderCount
			report.InvitePaidUserCount = row.PaidUserCount
			continue
		}
		report.AgentTopupMoney = analyticsRoundMoney(row.TopupMoney)
		report.AgentOrderCount = row.OrderCount
		report.AgentPaidUserCount = row.PaidUserCount
	}

	inviteReward := analyticsRewardRow{}
	if err := DB.Model(&dev.InviteRewardRecord{}).
		Select("COALESCE(SUM(reward_quota), 0) AS reward_quota, COUNT(*) AS record_count").
		Where("created_at >= ? AND created_at <= ?", startTimestamp, endTimestamp).
		Scan(&inviteReward).Error; err != nil {
		return report, err
	}
	agentCommission := analyticsCommissionRow{}
	if err := DB.Model(&AgentCommissionRecord{}).
		Select("COALESCE(SUM(commission_amount), 0) AS commission_amount, COUNT(*) AS record_count").
		Where("source_type = ? AND created_at >= ? AND created_at <= ?", AgentCommissionSourceTopup, startTimestamp, endTimestamp).
		Scan(&agentCommission).Error; err != nil {
		return report, err
	}

	report.InviteRewardQuota = inviteReward.RewardQuota
	report.InviteRewardCount = inviteReward.RecordCount
	report.AgentCommissionAmount = agentCommission.CommissionAmount
	report.AgentCommissionCount = agentCommission.RecordCount
	report.AgentNetMoney = analyticsRoundMoney(report.AgentTopupMoney - float64(report.AgentCommissionAmount)/100)
	return report, nil
}

func GetAdminAnalyticsCacheOverview(startTimestamp int64, endTimestamp int64) (AnalyticsCacheOverview, error) {
	databaseType := common.LogDatabaseType()
	cacheTokensExpr := analyticsJSONNumberExpr("other", "cache_tokens", databaseType)
	cacheRatioExpr := analyticsJSONNumberExpr("other", "cache_ratio", databaseType)
	cacheCreationExpr := analyticsJSONNumberExpr("other", "cache_creation_tokens", databaseType)
	cacheCreation5mExpr := analyticsJSONNumberExpr("other", "cache_creation_tokens_5m", databaseType)
	cacheCreation1hExpr := analyticsJSONNumberExpr("other", "cache_creation_tokens_1h", databaseType)
	modelRatioExpr := analyticsJSONNumberExpr("other", "model_ratio", databaseType)
	groupRatioExpr := analyticsJSONNumberExpr("other", "group_ratio", databaseType)
	userGroupRatioExpr := analyticsJSONNumberExpr("other", "user_group_ratio", databaseType)
	modelPriceExpr := analyticsJSONNumberExpr("other", "model_price", databaseType)
	usageSemanticExpr := analyticsJSONStringExpr("other", "usage_semantic", databaseType)
	effectiveGroupRatioExpr := "(CASE WHEN " + userGroupRatioExpr + " > 0 THEN " + userGroupRatioExpr + " ELSE " + groupRatioExpr + " END)"
	rateDenominatorExpr := "(CASE WHEN " + usageSemanticExpr + " = 'anthropic' THEN prompt_tokens + " + cacheTokensExpr + " ELSE prompt_tokens END)"
	splitCreationExpr := "(" + cacheCreation5mExpr + " + " + cacheCreation1hExpr + ")"
	cacheWriteExpr := "(CASE WHEN " + splitCreationExpr + " > 0 THEN CASE WHEN " + cacheCreationExpr + " > " + splitCreationExpr + " THEN " + cacheCreationExpr + " ELSE " + splitCreationExpr + " END ELSE " + cacheCreationExpr + " END)"
	inputTokensExpr := "(CASE WHEN " + usageSemanticExpr + " = 'anthropic' THEN prompt_tokens + " + cacheTokensExpr + " + " + cacheWriteExpr + " ELSE prompt_tokens END)"
	savedQuotaExpr := "(CASE WHEN " + modelPriceExpr + " >= 0 THEN 0 ELSE " + cacheTokensExpr + " * CASE WHEN " + cacheRatioExpr + " < 1 THEN 1 - " + cacheRatioExpr + " ELSE 0 END * " + modelRatioExpr + " * " + effectiveGroupRatioExpr + " END)"

	row := analyticsCacheOverviewRow{}
	if err := LOG_DB.Table("logs").
		Select(
			"COALESCE(SUM(prompt_tokens), 0) AS prompt_tokens, "+
				"COALESCE(SUM("+inputTokensExpr+"), 0) AS input_tokens_total, "+
				"COALESCE(SUM("+cacheTokensExpr+"), 0) AS cache_hit_tokens, "+
				"COALESCE(SUM("+rateDenominatorExpr+"), 0) AS cache_rate_denominator, "+
				"COALESCE(SUM("+cacheCreationExpr+"), 0) AS cache_creation_tokens, "+
				"COALESCE(SUM("+cacheCreation5mExpr+"), 0) AS cache_creation_tokens_5m, "+
				"COALESCE(SUM("+cacheCreation1hExpr+"), 0) AS cache_creation_tokens_1h, "+
				"COALESCE(SUM("+cacheWriteExpr+"), 0) AS cache_write_tokens, "+
				"COALESCE(SUM("+savedQuotaExpr+"), 0) AS cache_saved_quota",
		).
		Where("type = ? AND created_at >= ? AND created_at <= ?", LogTypeConsume, startTimestamp, endTimestamp).
		Scan(&row).Error; err != nil {
		return AnalyticsCacheOverview{}, err
	}

	result := AnalyticsCacheOverview{
		PromptTokens:          analyticsRoundInt64(row.PromptTokens),
		InputTokensTotal:      analyticsRoundInt64(row.InputTokensTotal),
		CacheHitTokens:        analyticsRoundInt64(row.CacheHitTokens),
		CacheCreationTokens:   analyticsRoundInt64(row.CacheCreationTokens),
		CacheCreationTokens5m: analyticsRoundInt64(row.CacheCreationTokens5m),
		CacheCreationTokens1h: analyticsRoundInt64(row.CacheCreationTokens1h),
		CacheWriteTokens:      analyticsRoundInt64(row.CacheWriteTokens),
		CacheSavedQuota:       analyticsRoundInt64(row.CacheSavedQuota),
	}
	if row.CacheRateDenominator > 0 {
		result.CacheHitRate = analyticsRoundPercent(row.CacheHitTokens * 100 / row.CacheRateDenominator)
	}
	return result, nil
}

func GetAdminAnalyticsRedeemedQuota(startTimestamp int64, endTimestamp int64) (int64, error) {
	var total int64
	err := DB.Model(&Redemption{}).
		Where("status = ? AND redeemed_time >= ? AND redeemed_time <= ?", common.RedemptionCodeStatusUsed, startTimestamp, endTimestamp).
		Select("COALESCE(SUM(quota), 0)").
		Scan(&total).Error
	return total, err
}
