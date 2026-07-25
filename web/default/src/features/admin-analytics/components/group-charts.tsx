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
  AnalyticsGroupStat,
} from '@/features/admin-analytics/types'

import { AnalyticsChartCard } from './analytics-chart-card'

const EMPTY_GROUPS: AnalyticsGroupStat[] = []

export function GroupCharts(props: {
  report?: AdminAnalyticsReport
  loading: boolean
}) {
  const { t } = useTranslation()
  const groups = props.report?.group_stats ?? EMPTY_GROUPS

  const consumptionSpec = useMemo(
    () => ({
      type: 'bar' as const,
      direction: 'horizontal' as const,
      data: [{ id: 'group-consumption', values: groups }],
      xField: 'consume_quota',
      yField: 'group_name',
      axes: [
        {
          orient: 'bottom' as const,
          label: {
            formatMethod: (value: number) => formatAnalyticsQuota(value),
          },
        },
        {
          orient: 'left' as const,
          label: { autoHide: true, autoLimit: true },
        },
      ],
      bar: { style: { cornerRadius: 6 } },
      tooltip: {
        mark: {
          content: [
            {
              key: (datum: AnalyticsGroupStat) => datum.group_name,
              value: (datum: AnalyticsGroupStat) =>
                formatAnalyticsQuota(datum.consume_quota),
            },
          ],
        },
      },
    }),
    [groups]
  )

  const incomeSpec = useMemo(
    () => ({
      type: 'bar' as const,
      direction: 'horizontal' as const,
      data: [{ id: 'group-income', values: groups }],
      xField: 'topup_money',
      yField: 'group_name',
      axes: [
        {
          orient: 'bottom' as const,
          label: {
            formatMethod: (value: number) => formatAnalyticsMoney(value),
          },
        },
        {
          orient: 'left' as const,
          label: { autoHide: true, autoLimit: true },
        },
      ],
      bar: { style: { cornerRadius: 6 } },
      tooltip: {
        mark: {
          content: [
            {
              key: (datum: AnalyticsGroupStat) => datum.group_name,
              value: (datum: AnalyticsGroupStat) =>
                formatAnalyticsMoney(datum.topup_money),
            },
            {
              key: t('Orders'),
              value: (datum: AnalyticsGroupStat) =>
                formatAnalyticsNumber(datum.topup_order_count),
            },
          ],
        },
      },
    }),
    [groups, t]
  )

  return (
    <div className='grid grid-cols-1 gap-4 xl:grid-cols-2'>
      <AnalyticsChartCard
        title={t('Consumption by group')}
        description={t('Aggregates billed quota by the request group.')}
        loading={props.loading}
        hasData={groups.some((item) => Number(item.consume_quota || 0) > 0)}
        spec={consumptionSpec}
      />
      <AnalyticsChartCard
        title={t('Top-up income by group')}
        description={t(
          'Attributes successful top-ups to each user’s current group.'
        )}
        loading={props.loading}
        hasData={groups.some((item) => Number(item.topup_money || 0) > 0)}
        spec={incomeSpec}
      />
    </div>
  )
}
