package service

import (
	"errors"
	"time"

	"github.com/QuantumNous/new-api/model"
	"golang.org/x/sync/errgroup"
)

const adminAnalyticsMaxRangeDays = 366

func normalizeAdminAnalyticsRange(startTimestamp int64, endTimestamp int64, location *time.Location) (int64, int64, error) {
	if location == nil {
		location = time.Local
	}
	now := time.Now().In(location)
	if startTimestamp <= 0 {
		startTimestamp = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, location).Unix()
	}
	if endTimestamp <= 0 {
		endTimestamp = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, location).Unix()
	}
	if endTimestamp < startTimestamp {
		return 0, 0, errors.New("end timestamp cannot be earlier than start timestamp")
	}
	if endTimestamp-startTimestamp > adminAnalyticsMaxRangeDays*24*60*60 {
		return 0, 0, errors.New("analytics range cannot exceed 366 days")
	}
	return startTimestamp, endTimestamp, nil
}

func fillAdminAnalyticsDailySeries(
	startTimestamp int64,
	endTimestamp int64,
	location *time.Location,
	consumeItems []model.AnalyticsDayQuota,
	topupItems []model.AnalyticsDayMoney,
	userItems []model.AnalyticsDayUserStat,
) ([]model.AnalyticsDayQuota, []model.AnalyticsDayMoney, []model.AnalyticsDayUserStat) {
	consumeByDay := make(map[string]int64, len(consumeItems))
	for _, item := range consumeItems {
		consumeByDay[item.Day] = item.Quota
	}
	topupByDay := make(map[string]float64, len(topupItems))
	for _, item := range topupItems {
		topupByDay[item.Day] = item.Money
	}
	usersByDay := make(map[string]model.AnalyticsDayUserStat, len(userItems))
	for _, item := range userItems {
		usersByDay[item.Day] = item
	}

	start := time.Unix(startTimestamp, 0).In(location)
	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, location)
	end := time.Unix(endTimestamp, 0).In(location)
	end = time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, location)

	consumes := make([]model.AnalyticsDayQuota, 0)
	topups := make([]model.AnalyticsDayMoney, 0)
	users := make([]model.AnalyticsDayUserStat, 0)
	for current := start; !current.After(end); current = current.AddDate(0, 0, 1) {
		day := current.Format("2006-01-02")
		consumes = append(consumes, model.AnalyticsDayQuota{Day: day, Quota: consumeByDay[day]})
		topups = append(topups, model.AnalyticsDayMoney{Day: day, Money: topupByDay[day]})
		userStat, ok := usersByDay[day]
		if !ok {
			userStat = model.AnalyticsDayUserStat{Day: day}
		}
		users = append(users, userStat)
	}
	return consumes, topups, users
}

func GetAdminAnalyticsReport(startTimestamp int64, endTimestamp int64, timezoneOffsetMinutes int) (*model.AdminAnalyticsReport, error) {
	location := time.FixedZone("analytics", timezoneOffsetMinutes*60)
	startTimestamp, endTimestamp, err := normalizeAdminAnalyticsRange(startTimestamp, endTimestamp, location)
	if err != nil {
		return nil, err
	}

	report := &model.AdminAnalyticsReport{
		StartTimestamp: startTimestamp,
		EndTimestamp:   endTimestamp,
	}
	var totalActiveUsers int64
	var totalRedeemedQuota int64
	group := new(errgroup.Group)
	group.Go(func() error {
		var queryErr error
		report.ModelRanking, queryErr = model.GetAdminAnalyticsModelRanking(startTimestamp, endTimestamp, 10)
		return queryErr
	})
	group.Go(func() error {
		var queryErr error
		report.ChannelRates, report.Overall, queryErr = model.GetAdminAnalyticsChannelRates(startTimestamp, endTimestamp)
		return queryErr
	})
	group.Go(func() error {
		var queryErr error
		report.ChannelConsumptions, queryErr = model.GetAdminAnalyticsChannelConsumptions(startTimestamp, endTimestamp)
		return queryErr
	})
	group.Go(func() error {
		var queryErr error
		report.ChannelLatencies, queryErr = model.GetAdminAnalyticsChannelLatencies(startTimestamp, endTimestamp)
		return queryErr
	})
	group.Go(func() error {
		var queryErr error
		report.UserConsumptions, queryErr = model.GetAdminAnalyticsUserConsumptions(startTimestamp, endTimestamp, 10)
		return queryErr
	})
	group.Go(func() error {
		var queryErr error
		report.GroupStats, queryErr = model.GetAdminAnalyticsGroupStats(startTimestamp, endTimestamp)
		return queryErr
	})
	group.Go(func() error {
		var queryErr error
		report.DailyConsumes, queryErr = model.GetAdminAnalyticsDailyConsumes(startTimestamp, endTimestamp, timezoneOffsetMinutes)
		return queryErr
	})
	group.Go(func() error {
		var queryErr error
		report.DailyTopups, queryErr = model.GetAdminAnalyticsDailyTopups(startTimestamp, endTimestamp, timezoneOffsetMinutes)
		return queryErr
	})
	group.Go(func() error {
		var queryErr error
		report.DailyUserStats, totalActiveUsers, queryErr = model.GetAdminAnalyticsDailyUserStats(startTimestamp, endTimestamp, timezoneOffsetMinutes)
		return queryErr
	})
	group.Go(func() error {
		var queryErr error
		report.ReferralOverview, queryErr = model.GetAdminAnalyticsReferralOverview(startTimestamp, endTimestamp)
		return queryErr
	})
	group.Go(func() error {
		var queryErr error
		report.CacheOverview, queryErr = model.GetAdminAnalyticsCacheOverview(startTimestamp, endTimestamp)
		return queryErr
	})
	group.Go(func() error {
		var queryErr error
		totalRedeemedQuota, queryErr = model.GetAdminAnalyticsRedeemedQuota(startTimestamp, endTimestamp)
		return queryErr
	})
	if err := group.Wait(); err != nil {
		return nil, err
	}

	report.DailyConsumes, report.DailyTopups, report.DailyUserStats = fillAdminAnalyticsDailySeries(
		startTimestamp,
		endTimestamp,
		location,
		report.DailyConsumes,
		report.DailyTopups,
		report.DailyUserStats,
	)

	var totalConsumeQuota int64
	for _, item := range report.DailyConsumes {
		totalConsumeQuota += item.Quota
	}
	var totalTopupMoney float64
	for _, item := range report.DailyTopups {
		totalTopupMoney += item.Money
	}
	var newUserTotal int64
	var paidUserTotal int64
	for _, item := range report.DailyUserStats {
		newUserTotal += item.NewUserCount
		paidUserTotal += item.PaidUserCount
	}

	report.Summary = model.AdminAnalyticsSummary{
		TotalRequestCount:  report.Overall.SuccessCount + report.Overall.FailureCount,
		TotalSuccessCount:  report.Overall.SuccessCount,
		TotalFailureCount:  report.Overall.FailureCount,
		TotalConsumeQuota:  totalConsumeQuota,
		TotalTopupMoney:    totalTopupMoney,
		TotalActiveUsers:   totalActiveUsers,
		NewUserTotal:       newUserTotal,
		PaidUserTotal:      paidUserTotal,
		CacheHitRate:       report.CacheOverview.CacheHitRate,
		CacheSavedQuota:    report.CacheOverview.CacheSavedQuota,
		TotalRedeemedQuota: totalRedeemedQuota,
	}
	return report, nil
}
