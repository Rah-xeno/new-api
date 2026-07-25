/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

For commercial licensing, please contact support@quantumnous.com
*/
import { useQuery } from '@tanstack/react-query'
import { AlertCircle, BarChart3, RefreshCw } from 'lucide-react'
import { useMemo, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'

import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { IconBadge } from '@/components/ui/icon-badge'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Spinner } from '@/components/ui/spinner'

import { getAdminAnalyticsReport } from './api'
import { ChannelCharts } from './components/channel-charts'
import { GroupCharts } from './components/group-charts'
import {
  MetricSection,
  SummaryMetricCards,
  type AnalyticsMetric,
} from './components/metric-cards'
import { RankingCharts } from './components/ranking-charts'
import { TrendCharts } from './components/trend-charts'
import {
  analyticsRangeFromDays,
  analyticsRangeIsValid,
  analyticsRangeTimestamps,
  formatAnalyticsCents,
  formatAnalyticsMoney,
  formatAnalyticsNumber,
  formatAnalyticsPercent,
  formatAnalyticsQuota,
  todayAnalyticsRange,
} from './lib'
import type { AnalyticsDateRange } from './types'

const QUICK_RANGES = [
  { key: 'today', label: 'Today', getRange: todayAnalyticsRange },
  {
    key: 'yesterday',
    label: 'Yesterday',
    getRange: () => analyticsRangeFromDays(1, 1),
  },
  {
    key: 'seven-days',
    label: 'Last 7 days',
    getRange: () => analyticsRangeFromDays(6),
  },
  {
    key: 'thirty-days',
    label: 'Last 30 days',
    getRange: () => analyticsRangeFromDays(29),
  },
] as const

export function AdminAnalytics() {
  const { t } = useTranslation()
  const [draftRange, setDraftRange] =
    useState<AnalyticsDateRange>(todayAnalyticsRange)
  const [appliedRange, setAppliedRange] =
    useState<AnalyticsDateRange>(todayAnalyticsRange)
  const timestamps = useMemo(
    () => analyticsRangeTimestamps(appliedRange),
    [appliedRange]
  )
  const timezoneOffsetMinutes = -new Date().getTimezoneOffset()
  const reportQuery = useQuery({
    queryKey: [
      'admin-analytics',
      timestamps.start_timestamp,
      timestamps.end_timestamp,
      timezoneOffsetMinutes,
    ],
    queryFn: () =>
      getAdminAnalyticsReport({
        ...timestamps,
        tz_offset_minutes: timezoneOffsetMinutes,
      }),
    staleTime: 60_000,
  })
  const report = reportQuery.data

  const applyRange = (range: AnalyticsDateRange) => {
    if (!analyticsRangeIsValid(range)) {
      toast.error(t('Select a valid start and end date.'))
      return
    }
    setDraftRange(range)
    if (range.start === appliedRange.start && range.end === appliedRange.end) {
      void reportQuery.refetch()
      return
    }
    setAppliedRange(range)
  }

  const summaryMetrics = useMemo<AnalyticsMetric[]>(() => {
    const summary = report?.summary
    return [
      {
        key: 'requests',
        label: t('Total requests'),
        value: formatAnalyticsNumber(summary?.total_request_count || 0),
        hint: `${t('Successful')} ${formatAnalyticsNumber(summary?.total_success_count || 0)} / ${t('Failed')} ${formatAnalyticsNumber(summary?.total_failure_count || 0)}`,
      },
      {
        key: 'success-rate',
        label: t('Overall success rate'),
        value: formatAnalyticsPercent(report?.overall.success_rate || 0),
        hint: `${t('Failure rate')} ${formatAnalyticsPercent(report?.overall.failure_rate || 0)}`,
      },
      {
        key: 'consumption',
        label: t('Total consumption'),
        value: formatAnalyticsQuota(summary?.total_consume_quota || 0),
        hint: t('Calculated from successful usage logs.'),
      },
      {
        key: 'topups',
        label: t('Total top-up income'),
        value: formatAnalyticsMoney(summary?.total_topup_money || 0),
        hint: t('Calculated from completed top-up orders.'),
      },
      {
        key: 'active-users',
        label: t('Active users'),
        value: formatAnalyticsNumber(summary?.total_active_users || 0),
        hint: t('Distinct users with successful or failed requests.'),
      },
      {
        key: 'new-paid-users',
        label: t('New / first-time paying users'),
        value: `${formatAnalyticsNumber(summary?.new_user_total || 0)} / ${formatAnalyticsNumber(summary?.paid_user_total || 0)}`,
        hint: t('Based on account creation and first payment time.'),
      },
      {
        key: 'cache-rate',
        label: t('Cache hit rate'),
        value: formatAnalyticsPercent(summary?.cache_hit_rate || 0),
        hint: t('Cache-read tokens divided by eligible input tokens.'),
      },
      {
        key: 'cache-savings',
        label: t('Quota saved by cache'),
        value: formatAnalyticsQuota(summary?.cache_saved_quota || 0),
        hint: t('Estimated from model, group and cache ratios.'),
      },
      {
        key: 'redemptions',
        label: t('Redeemed quota'),
        value: formatAnalyticsQuota(summary?.total_redeemed_quota || 0),
        hint: t('Counts redemption codes used during the selected period.'),
      },
    ]
  }, [report, t])

  const referral = report?.referral_overview
  const inviteMetrics: AnalyticsMetric[] = [
    {
      key: 'invite-topups',
      label: t('Top-up income'),
      value: formatAnalyticsMoney(referral?.invite_topup_money || 0),
      hint: `${t('Orders')} ${formatAnalyticsNumber(referral?.invite_order_count || 0)} · ${t('Paying users')} ${formatAnalyticsNumber(referral?.invite_paid_user_count || 0)}`,
    },
    {
      key: 'invite-rewards',
      label: t('Rewards issued'),
      value: formatAnalyticsQuota(referral?.invite_reward_quota || 0),
      hint: `${t('Reward records')} ${formatAnalyticsNumber(referral?.invite_reward_count || 0)}`,
    },
  ]
  const agentMetrics: AnalyticsMetric[] = [
    {
      key: 'agent-topups',
      label: t('Top-up income'),
      value: formatAnalyticsMoney(referral?.agent_topup_money || 0),
      hint: `${t('Orders')} ${formatAnalyticsNumber(referral?.agent_order_count || 0)} · ${t('Paying users')} ${formatAnalyticsNumber(referral?.agent_paid_user_count || 0)}`,
    },
    {
      key: 'agent-commission',
      label: t('Commission / net income'),
      value: `${formatAnalyticsCents(referral?.agent_commission_amount || 0)} / ${formatAnalyticsMoney(referral?.agent_net_money || 0)}`,
      hint: `${t('Commission records')} ${formatAnalyticsNumber(referral?.agent_commission_count || 0)}`,
    },
  ]

  const cache = report?.cache_overview
  const cacheCreationTokens =
    Number(cache?.cache_write_tokens || 0) ||
    Number(cache?.cache_creation_tokens || 0) +
      Number(cache?.cache_creation_tokens_5m || 0) +
      Number(cache?.cache_creation_tokens_1h || 0)
  const cacheMetrics: AnalyticsMetric[] = [
    {
      key: 'input-tokens',
      label: t('Total input tokens'),
      value: formatAnalyticsNumber(
        cache?.input_tokens_total || cache?.prompt_tokens || 0
      ),
    },
    {
      key: 'cache-hit-tokens',
      label: t('Cache-read tokens'),
      value: formatAnalyticsNumber(cache?.cache_hit_tokens || 0),
      hint: formatAnalyticsPercent(cache?.cache_hit_rate || 0),
    },
    {
      key: 'cache-saved',
      label: t('Quota saved'),
      value: formatAnalyticsQuota(cache?.cache_saved_quota || 0),
    },
    {
      key: 'cache-created',
      label: t('Cache-write tokens'),
      value: formatAnalyticsNumber(cacheCreationTokens),
      hint: `base ${formatAnalyticsNumber(cache?.cache_creation_tokens || 0)} · 5m ${formatAnalyticsNumber(cache?.cache_creation_tokens_5m || 0)} · 1h ${formatAnalyticsNumber(cache?.cache_creation_tokens_1h || 0)}`,
    },
  ]

  return (
    <div className='flex flex-1 flex-col gap-4'>
      <Card>
        <CardHeader>
          <div className='flex flex-col gap-4 xl:flex-row xl:items-start xl:justify-between'>
            <div className='flex items-start gap-3'>
              <IconBadge tone='primary' size='lg'>
                <BarChart3 />
              </IconBadge>
              <div>
                <CardTitle className='text-xl'>{t('Data analytics')}</CardTitle>
                <CardDescription className='mt-1 max-w-2xl'>
                  {t(
                    'Explore model usage, channel reliability, income, user growth and cache efficiency for one period.'
                  )}
                </CardDescription>
              </div>
            </div>
            <Badge variant='secondary'>
              {t('Current period')}: {appliedRange.start} – {appliedRange.end}
            </Badge>
          </div>
        </CardHeader>
        <CardContent className='flex flex-col gap-3'>
          <div className='flex flex-col gap-3 lg:flex-row lg:items-end'>
            <div className='grid flex-1 grid-cols-1 gap-3 sm:grid-cols-2'>
              <div className='flex flex-col gap-1.5'>
                <Label htmlFor='analytics-start-date'>{t('Start date')}</Label>
                <Input
                  id='analytics-start-date'
                  type='date'
                  value={draftRange.start}
                  max={draftRange.end}
                  onChange={(event) =>
                    setDraftRange((current) => ({
                      ...current,
                      start: event.target.value,
                    }))
                  }
                />
              </div>
              <div className='flex flex-col gap-1.5'>
                <Label htmlFor='analytics-end-date'>{t('End date')}</Label>
                <Input
                  id='analytics-end-date'
                  type='date'
                  value={draftRange.end}
                  min={draftRange.start}
                  onChange={(event) =>
                    setDraftRange((current) => ({
                      ...current,
                      end: event.target.value,
                    }))
                  }
                />
              </div>
            </div>
            <Button
              onClick={() => applyRange(draftRange)}
              disabled={reportQuery.isFetching}
            >
              {reportQuery.isFetching ? (
                <Spinner data-icon='inline-start' />
              ) : (
                <RefreshCw data-icon='inline-start' />
              )}
              {t('Refresh report')}
            </Button>
          </div>
          <div className='flex flex-wrap gap-2'>
            {QUICK_RANGES.map((range) => (
              <Button
                key={range.key}
                type='button'
                variant='outline'
                size='sm'
                onClick={() => applyRange(range.getRange())}
              >
                {t(range.label)}
              </Button>
            ))}
          </div>
        </CardContent>
      </Card>

      {reportQuery.error ? (
        <Alert variant='destructive'>
          <AlertCircle />
          <AlertTitle>{t('Unable to load analytics')}</AlertTitle>
          <AlertDescription>{reportQuery.error.message}</AlertDescription>
        </Alert>
      ) : null}

      <SummaryMetricCards
        metrics={summaryMetrics}
        loading={reportQuery.isLoading}
      />

      <div className='grid grid-cols-1 gap-4 xl:grid-cols-3'>
        <MetricSection
          title={t('Invite referrals')}
          description={t('Top-ups and rewards from standard invite referrals.')}
          metrics={inviteMetrics}
          loading={reportQuery.isLoading}
        />
        <MetricSection
          title={t('Agent distribution')}
          description={t('Top-ups, commissions and net agent contribution.')}
          metrics={agentMetrics}
          loading={reportQuery.isLoading}
        />
        <MetricSection
          title={t('Cache efficiency')}
          description={t('Cache reads, writes and estimated quota savings.')}
          metrics={cacheMetrics}
          loading={reportQuery.isLoading}
        />
      </div>

      <RankingCharts report={report} loading={reportQuery.isLoading} />
      <ChannelCharts report={report} loading={reportQuery.isLoading} />
      <GroupCharts report={report} loading={reportQuery.isLoading} />
      <TrendCharts report={report} loading={reportQuery.isLoading} />
    </div>
  )
}
