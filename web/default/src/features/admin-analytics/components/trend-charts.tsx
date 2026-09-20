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
import { useMemo } from 'react'
import { useTranslation } from 'react-i18next'

import {
  formatAnalyticsMoney,
  formatAnalyticsNumber,
  formatAnalyticsQuota,
} from '@/features/admin-analytics/lib'
import type {
  AdminAnalyticsReport,
  AnalyticsDayMoney,
  AnalyticsDayQuota,
  AnalyticsDayUserStat,
} from '@/features/admin-analytics/types'

import { AnalyticsChartCard } from './analytics-chart-card'

type UserTrendDatum = {
  day: string
  metric: string
  value: number
}

const EMPTY_DAILY_CONSUMES: AnalyticsDayQuota[] = []
const EMPTY_DAILY_TOPUPS: AnalyticsDayMoney[] = []
const EMPTY_DAILY_USERS: AnalyticsDayUserStat[] = []

export function TrendCharts(props: {
  report?: AdminAnalyticsReport
  loading: boolean
}) {
  const { t } = useTranslation()
  const dailyConsumes = props.report?.daily_consumes ?? EMPTY_DAILY_CONSUMES
  const dailyTopups = props.report?.daily_topups ?? EMPTY_DAILY_TOPUPS
  const dailyUsers = props.report?.daily_user_stats ?? EMPTY_DAILY_USERS

  const consumeSpec = useMemo(
    () => ({
      type: 'line' as const,
      data: [{ id: 'daily-consumption', values: dailyConsumes }],
      xField: 'day',
      yField: 'quota',
      point: { visible: true },
      area: { visible: true, style: { fillOpacity: 0.16 } },
      axes: [
        {
          orient: 'bottom' as const,
          label: { autoRotate: true, autoHide: true },
        },
        {
          orient: 'left' as const,
          label: {
            formatMethod: (value: number) => formatAnalyticsQuota(value),
          },
        },
      ],
      tooltip: {
        mark: {
          content: [
            {
              key: (datum: AnalyticsDayQuota) => datum.day,
              value: (datum: AnalyticsDayQuota) =>
                formatAnalyticsQuota(datum.quota),
            },
          ],
        },
      },
    }),
    [dailyConsumes]
  )

  const topupSpec = useMemo(
    () => ({
      type: 'line' as const,
      data: [{ id: 'daily-topups', values: dailyTopups }],
      xField: 'day',
      yField: 'money',
      point: { visible: true },
      area: { visible: true, style: { fillOpacity: 0.16 } },
      axes: [
        {
          orient: 'bottom' as const,
          label: { autoRotate: true, autoHide: true },
        },
        {
          orient: 'left' as const,
          label: {
            formatMethod: (value: number) => formatAnalyticsMoney(value),
          },
        },
      ],
      tooltip: {
        mark: {
          content: [
            {
              key: (datum: AnalyticsDayMoney) => datum.day,
              value: (datum: AnalyticsDayMoney) =>
                formatAnalyticsMoney(datum.money),
            },
          ],
        },
      },
    }),
    [dailyTopups]
  )

  const userData = useMemo(
    () =>
      dailyUsers.flatMap<UserTrendDatum>((item) => [
        {
          day: item.day,
          metric: t('New users'),
          value: Number(item.new_user_count || 0),
        },
        {
          day: item.day,
          metric: t('First-time paying users'),
          value: Number(item.paid_user_count || 0),
        },
        {
          day: item.day,
          metric: t('Active users'),
          value: Number(item.active_user_count || 0),
        },
      ]),
    [dailyUsers, t]
  )
  const userSpec = useMemo(
    () => ({
      type: 'line' as const,
      data: [{ id: 'daily-users', values: userData }],
      xField: 'day',
      yField: 'value',
      seriesField: 'metric',
      legends: { visible: true, orient: 'top' as const },
      point: { visible: true },
      axes: [
        {
          orient: 'bottom' as const,
          label: { autoRotate: true, autoHide: true },
        },
        {
          orient: 'left' as const,
          label: {
            formatMethod: (value: number) => formatAnalyticsNumber(value),
          },
        },
      ],
      tooltip: {
        mark: {
          content: [
            {
              key: (datum: UserTrendDatum) => datum.day,
              value: (datum: UserTrendDatum) =>
                `${datum.metric} ${formatAnalyticsNumber(datum.value)}`,
            },
          ],
        },
      },
    }),
    [userData]
  )

  return (
    <>
      <div className='grid grid-cols-1 gap-4 xl:grid-cols-2'>
        <AnalyticsChartCard
          title={t('Daily consumption')}
          description={t('Shows billed quota for each day in the period.')}
          loading={props.loading}
          hasData={dailyConsumes.some((item) => Number(item.quota || 0) > 0)}
          spec={consumeSpec}
        />
        <AnalyticsChartCard
          title={t('Daily top-up income')}
          description={t('Includes successfully completed top-up orders only.')}
          loading={props.loading}
          hasData={dailyTopups.some((item) => Number(item.money || 0) > 0)}
          spec={topupSpec}
        />
      </div>
      <AnalyticsChartCard
        title={t('New, paying and active user trends')}
        description={t(
          'Tracks account creation, first payment and distinct active users by day.'
        )}
        loading={props.loading}
        hasData={dailyUsers.some(
          (item) =>
            Number(item.new_user_count || 0) > 0 ||
            Number(item.paid_user_count || 0) > 0 ||
            Number(item.active_user_count || 0) > 0
        )}
        spec={userSpec}
      />
    </>
  )
}
