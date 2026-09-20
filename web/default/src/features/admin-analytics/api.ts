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
import { t } from 'i18next'

import { api } from '@/lib/api'

import type { AdminAnalyticsReport } from './types'

type AdminAnalyticsResponse = {
  success: boolean
  message?: string
  data?: AdminAnalyticsReport
}

export async function getAdminAnalyticsReport(params: {
  start_timestamp: number
  end_timestamp: number
  tz_offset_minutes: number
}): Promise<AdminAnalyticsReport> {
  const response = await api.get<AdminAnalyticsResponse>(
    '/api/analytics/admin/report',
    {
      params,
      skipBusinessError: true,
    }
  )
  if (!response.data.success || !response.data.data) {
    throw new Error(response.data.message || t('Unable to load analytics'))
  }
  return response.data.data
}
