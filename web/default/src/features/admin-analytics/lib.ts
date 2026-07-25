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
import { formatLocalCurrencyAmount } from '@/lib/currency'
import dayjs from '@/lib/dayjs'
import { formatNumber, formatPercent, formatQuota } from '@/lib/format'

import type { AnalyticsDateRange } from './types'

export function todayAnalyticsRange(): AnalyticsDateRange {
  const today = dayjs().format('YYYY-MM-DD')
  return { start: today, end: today }
}

export function analyticsRangeFromDays(
  startOffset: number,
  endOffset = 0
): AnalyticsDateRange {
  return {
    start: dayjs().subtract(startOffset, 'day').format('YYYY-MM-DD'),
    end: dayjs().subtract(endOffset, 'day').format('YYYY-MM-DD'),
  }
}

export function analyticsRangeTimestamps(range: AnalyticsDateRange) {
  return {
    start_timestamp: dayjs(range.start).startOf('day').unix(),
    end_timestamp: dayjs(range.end).endOf('day').unix(),
  }
}

export function analyticsRangeIsValid(range: AnalyticsDateRange): boolean {
  const start = dayjs(range.start)
  const end = dayjs(range.end)
  return start.isValid() && end.isValid() && !end.isBefore(start, 'day')
}

export function formatAnalyticsQuota(value: number): string {
  return formatQuota(Number(value || 0))
}

export function formatAnalyticsMoney(value: number): string {
  return formatLocalCurrencyAmount(Number(value || 0), {
    digitsLarge: 2,
    digitsSmall: 2,
    abbreviate: false,
  })
}

export function formatAnalyticsCents(value: number): string {
  return formatAnalyticsMoney(Number(value || 0) / 100)
}

export function formatAnalyticsNumber(value: number): string {
  return formatNumber(Number(value || 0))
}

export function formatAnalyticsPercent(value: number): string {
  return formatPercent(Number(value || 0))
}

export function formatAnalyticsSeconds(
  value: number,
  locale?: Intl.LocalesArgument
): string {
  const numericValue = Number(value || 0)
  if (!Number.isFinite(numericValue) || numericValue <= 0) return '0 s'
  return Intl.NumberFormat(locale, {
    style: 'unit',
    unit: 'second',
    unitDisplay: 'short',
    maximumFractionDigits: numericValue >= 10 ? 1 : 2,
  }).format(numericValue)
}
