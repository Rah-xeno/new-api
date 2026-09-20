package service

import (
	"testing"
	"time"

	"github.com/QuantumNous/new-api/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeAdminAnalyticsRangeRejectsOversizedRange(t *testing.T) {
	startTimestamp := time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC).Unix()
	endTimestamp := startTimestamp + (adminAnalyticsMaxRangeDays+1)*24*60*60

	_, _, err := normalizeAdminAnalyticsRange(startTimestamp, endTimestamp, time.UTC)

	require.Error(t, err)
	assert.ErrorContains(t, err, "cannot exceed 366 days")
}

func TestFillAdminAnalyticsDailySeriesFillsMissingDays(t *testing.T) {
	startTimestamp := time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC).Unix()
	endTimestamp := time.Date(2025, time.January, 3, 23, 59, 59, 0, time.UTC).Unix()

	consumes, topups, users := fillAdminAnalyticsDailySeries(
		startTimestamp,
		endTimestamp,
		time.UTC,
		[]model.AnalyticsDayQuota{{Day: "2025-01-02", Quota: 120}},
		[]model.AnalyticsDayMoney{{Day: "2025-01-03", Money: 5.5}},
		[]model.AnalyticsDayUserStat{{Day: "2025-01-01", NewUserCount: 2}},
	)

	require.Len(t, consumes, 3)
	require.Len(t, topups, 3)
	require.Len(t, users, 3)
	assert.Equal(t, []model.AnalyticsDayQuota{
		{Day: "2025-01-01"},
		{Day: "2025-01-02", Quota: 120},
		{Day: "2025-01-03"},
	}, consumes)
	assert.Equal(t, []model.AnalyticsDayMoney{
		{Day: "2025-01-01"},
		{Day: "2025-01-02"},
		{Day: "2025-01-03", Money: 5.5},
	}, topups)
	assert.Equal(t, []model.AnalyticsDayUserStat{
		{Day: "2025-01-01", NewUserCount: 2},
		{Day: "2025-01-02"},
		{Day: "2025-01-03"},
	}, users)
}
