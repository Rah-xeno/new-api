package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdminAnalyticsDailyUserStatsUsesRequestedTimezone(t *testing.T) {
	truncateTables(t)

	eventTime := time.Date(2025, time.January, 3, 23, 30, 0, 0, time.UTC).Unix()
	require.NoError(t, DB.Create(&User{
		Id:             101,
		Username:       "analytics-user",
		Password:       "analytics-password",
		CreatedAt:      eventTime,
		FirstPaymentAt: eventTime,
	}).Error)
	require.NoError(t, LOG_DB.Create(&Log{
		UserId:    101,
		CreatedAt: eventTime,
		Type:      LogTypeConsume,
	}).Error)

	location := time.FixedZone("UTC+2", 2*60*60)
	startTimestamp := time.Date(2025, time.January, 4, 0, 0, 0, 0, location).Unix()
	endTimestamp := time.Date(2025, time.January, 4, 23, 59, 59, 0, location).Unix()
	items, totalActiveUsers, err := GetAdminAnalyticsDailyUserStats(
		startTimestamp,
		endTimestamp,
		2*60,
	)

	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "2025-01-04", items[0].Day)
	assert.EqualValues(t, 1, items[0].NewUserCount)
	assert.EqualValues(t, 1, items[0].PaidUserCount)
	assert.EqualValues(t, 1, items[0].ActiveUserCount)
	assert.EqualValues(t, 1, totalActiveUsers)
}

func TestAdminAnalyticsCacheHitRateUsesAnthropicInputSemantics(t *testing.T) {
	truncateTables(t)

	now := time.Now().Unix()
	require.NoError(t, LOG_DB.Create(&Log{
		CreatedAt:    now,
		Type:         LogTypeConsume,
		PromptTokens: 100,
		Other:        `{"usage_semantic":"anthropic","cache_tokens":50}`,
	}).Error)

	overview, err := GetAdminAnalyticsCacheOverview(now-1, now+1)

	require.NoError(t, err)
	assert.EqualValues(t, 150, overview.InputTokensTotal)
	assert.EqualValues(t, 50, overview.CacheHitTokens)
	assert.InDelta(t, 33.33, overview.CacheHitRate, 0.001)
}
