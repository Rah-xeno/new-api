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
  formatAnalyticsNumber,
  formatAnalyticsQuota,
} from '@/features/admin-analytics/lib'
import type {
  AdminAnalyticsReport,
  AnalyticsNamedCount,
  AnalyticsUserQuota,
} from '@/features/admin-analytics/types'

import { AnalyticsChartCard } from './analytics-chart-card'

type RankingChartsProps = {
  report?: AdminAnalyticsReport
  loading: boolean
}

const EMPTY_MODEL_RANKING: AnalyticsNamedCount[] = []
const EMPTY_USER_RANKING: AnalyticsUserQuota[] = []

export function RankingCharts(props: RankingChartsProps) {
  const { t } = useTranslation()
  const modelRanking = props.report?.model_ranking ?? EMPTY_MODEL_RANKING
  const userRanking = props.report?.user_consumptions ?? EMPTY_USER_RANKING

  const modelSpec = useMemo(
    () => ({
      type: 'bar' as const,
      data: [{ id: 'model-ranking', values: modelRanking }],
      xField: 'name',
      yField: 'count',
      axes: [
        {
          orient: 'bottom' as const,
          label: { autoRotate: true, autoHide: true, autoLimit: true },
        },
        { orient: 'left' as const },
      ],
      label: {
        visible: true,
        position: 'top' as const,
        formatMethod: (value: number) => formatAnalyticsNumber(value),
      },
      bar: { style: { cornerRadius: [6, 6, 0, 0] } },
      tooltip: {
        mark: {
          content: [
            {
              key: (datum: AnalyticsNamedCount) => datum.name,
              value: (datum: AnalyticsNamedCount) =>
                formatAnalyticsNumber(datum.count),
            },
          ],
        },
      },
    }),
    [modelRanking]
  )

  const userSpec = useMemo(
    () => ({
      type: 'bar' as const,
      direction: 'horizontal' as const,
      data: [{ id: 'user-ranking', values: userRanking }],
      xField: 'quota',
      yField: 'username',
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
              key: (datum: AnalyticsUserQuota) => datum.username,
              value: (datum: AnalyticsUserQuota) =>
                formatAnalyticsQuota(datum.quota),
            },
          ],
        },
      },
    }),
    [userRanking]
  )

  const overallSpec = useMemo(() => {
    const values = [
      {
        category: t('Successful requests'),
        value: Number(props.report?.overall.success_count || 0),
      },
      {
        category: t('Failed requests'),
        value: Number(props.report?.overall.failure_count || 0),
      },
    ]
    return {
      type: 'pie' as const,
      data: [{ id: 'overall-rate', values }],
      valueField: 'value',
      categoryField: 'category',
      outerRadius: 0.76,
      innerRadius: 0.5,
      legends: { visible: true, orient: 'bottom' as const },
      label: {
        visible: true,
        formatMethod: (
          _value: number,
          datum: { category: string; value: number }
        ) => `${datum.category} ${formatAnalyticsNumber(datum.value)}`,
      },
      color: ['#16a34a', '#dc2626'],
    }
  }, [props.report, t])

  return (
    <>
      <div className='grid grid-cols-1 gap-4 xl:grid-cols-3'>
        <AnalyticsChartCard
          className='xl:col-span-2'
          title={t('Top models by calls')}
          description={t(
            'Ranks successful requests by model for the selected period.'
          )}
          loading={props.loading}
          hasData={modelRanking.length > 0}
          spec={modelSpec}
        />
        <AnalyticsChartCard
          title={t('Overall request success')}
          description={t(
            'Compares successful usage logs with failed request logs.'
          )}
          loading={props.loading}
          hasData={Number(props.report?.summary.total_request_count || 0) > 0}
          spec={overallSpec}
        />
      </div>
      <AnalyticsChartCard
        title={t('Top users by consumption')}
        description={t(
          'Ranks the ten highest-consuming users in the selected period.'
        )}
        loading={props.loading}
        hasData={userRanking.length > 0}
        spec={userSpec}
      />
    </>
  )
}
