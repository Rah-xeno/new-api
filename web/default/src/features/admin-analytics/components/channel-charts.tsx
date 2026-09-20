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
  formatAnalyticsPercent,
  formatAnalyticsQuota,
  formatAnalyticsSeconds,
} from '@/features/admin-analytics/lib'
import type {
  AdminAnalyticsReport,
  AnalyticsChannelLatency,
  AnalyticsChannelQuota,
  AnalyticsChannelRate,
} from '@/features/admin-analytics/types'

import { AnalyticsChartCard } from './analytics-chart-card'

type ChannelChartDatum = {
  channel_name: string
  metric: string
  value: number
  success_count?: number
  failure_count?: number
  request_count?: number
}

const EMPTY_CHANNEL_RATES: AnalyticsChannelRate[] = []
const EMPTY_CHANNEL_LATENCIES: AnalyticsChannelLatency[] = []
const EMPTY_CHANNEL_CONSUMPTIONS: AnalyticsChannelQuota[] = []

export function ChannelCharts(props: {
  report?: AdminAnalyticsReport
  loading: boolean
}) {
  const { t, i18n } = useTranslation()
  const channelRates = props.report?.channel_rates ?? EMPTY_CHANNEL_RATES
  const channelLatencies =
    props.report?.channel_latencies ?? EMPTY_CHANNEL_LATENCIES
  const channelConsumptions =
    props.report?.channel_consumptions ?? EMPTY_CHANNEL_CONSUMPTIONS

  const rateData = useMemo(
    () =>
      channelRates.flatMap<ChannelChartDatum>((item) => [
        {
          channel_name: item.channel_name,
          metric: t('Success rate'),
          value: Number(item.success_rate || 0),
          success_count: item.success_count,
          failure_count: item.failure_count,
        },
        {
          channel_name: item.channel_name,
          metric: t('Failure rate'),
          value: Number(item.failure_rate || 0),
          success_count: item.success_count,
          failure_count: item.failure_count,
        },
      ]),
    [channelRates, t]
  )
  const rateSpec = useMemo(
    () => ({
      type: 'bar' as const,
      direction: 'horizontal' as const,
      stack: true,
      data: [{ id: 'channel-rates', values: rateData }],
      xField: 'value',
      yField: 'channel_name',
      seriesField: 'metric',
      legends: { visible: true, orient: 'top' as const },
      axes: [
        {
          orient: 'bottom' as const,
          label: {
            formatMethod: (value: number) => formatAnalyticsPercent(value),
          },
        },
        {
          orient: 'left' as const,
          label: { autoHide: true, autoLimit: true },
        },
      ],
      color: ['#16a34a', '#dc2626'],
      tooltip: {
        mark: {
          content: [
            {
              key: (datum: ChannelChartDatum) => datum.channel_name,
              value: (datum: ChannelChartDatum) => datum.metric,
            },
            {
              key: t('Rate'),
              value: (datum: ChannelChartDatum) =>
                formatAnalyticsPercent(datum.value),
            },
            {
              key: t('Success / failure'),
              value: (datum: ChannelChartDatum) =>
                `${formatAnalyticsNumber(datum.success_count || 0)} / ${formatAnalyticsNumber(datum.failure_count || 0)}`,
            },
          ],
        },
      },
    }),
    [rateData, t]
  )

  const latencyData = useMemo(
    () =>
      channelLatencies.flatMap<ChannelChartDatum>((item) => [
        {
          channel_name: item.channel_name,
          metric: t('Average'),
          value: Number(item.avg_use_time || 0),
          request_count: item.request_count,
        },
        {
          channel_name: item.channel_name,
          metric: 'P95',
          value: Number(item.p95_use_time || 0),
          request_count: item.request_count,
        },
        {
          channel_name: item.channel_name,
          metric: 'P99',
          value: Number(item.p99_use_time || 0),
          request_count: item.request_count,
        },
      ]),
    [channelLatencies, t]
  )
  const latencySpec = useMemo(
    () => ({
      type: 'bar' as const,
      direction: 'horizontal' as const,
      data: [{ id: 'channel-latencies', values: latencyData }],
      xField: 'value',
      yField: 'channel_name',
      seriesField: 'metric',
      legends: { visible: true, orient: 'top' as const },
      axes: [
        {
          orient: 'bottom' as const,
          label: {
            formatMethod: (value: number) =>
              formatAnalyticsSeconds(value, i18n.language),
          },
        },
        {
          orient: 'left' as const,
          label: { autoHide: true, autoLimit: true },
        },
      ],
      color: ['#2563eb', '#d97706', '#dc2626'],
      tooltip: {
        mark: {
          content: [
            {
              key: (datum: ChannelChartDatum) => datum.channel_name,
              value: (datum: ChannelChartDatum) => datum.metric,
            },
            {
              key: t('Response time'),
              value: (datum: ChannelChartDatum) =>
                formatAnalyticsSeconds(datum.value, i18n.language),
            },
            {
              key: t('Samples'),
              value: (datum: ChannelChartDatum) =>
                formatAnalyticsNumber(datum.request_count || 0),
            },
          ],
        },
      },
    }),
    [i18n.language, latencyData, t]
  )

  const consumptionSpec = useMemo(
    () => ({
      type: 'bar' as const,
      direction: 'horizontal' as const,
      data: [{ id: 'channel-consumption', values: channelConsumptions }],
      xField: 'quota',
      yField: 'channel_name',
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
              key: (datum: AnalyticsChannelQuota) => datum.channel_name,
              value: (datum: AnalyticsChannelQuota) =>
                formatAnalyticsQuota(datum.quota),
            },
          ],
        },
      },
    }),
    [channelConsumptions]
  )

  return (
    <>
      <div className='grid grid-cols-1 gap-4 xl:grid-cols-2'>
        <AnalyticsChartCard
          title={t('Channel success and failure rates')}
          description={t(
            'Compares request reliability across upstream channels.'
          )}
          loading={props.loading}
          hasData={channelRates.length > 0}
          spec={rateSpec}
        />
        <AnalyticsChartCard
          title={t('Channel average, P95 and P99 latency')}
          description={t(
            'Highlights typical response time and long-tail latency.'
          )}
          loading={props.loading}
          hasData={channelLatencies.length > 0}
          spec={latencySpec}
        />
      </div>
      <AnalyticsChartCard
        title={t('Consumption by channel')}
        description={t(
          'Aggregates billed quota from successful requests by channel.'
        )}
        loading={props.loading}
        hasData={channelConsumptions.length > 0}
        spec={consumptionSpec}
      />
    </>
  )
}
