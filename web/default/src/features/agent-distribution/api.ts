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
import { api } from '@/lib/api'

import type {
  AgentCommissionRecord,
  AgentDashboard,
  AgentTopUp,
  PagedData,
} from './types'

type ApiResponse<T> = {
  success: boolean
  message?: string
  data: T
}

export async function getAgentDashboard(page: number) {
  const response = await api.get<ApiResponse<AgentDashboard>>(
    '/api/user/agent/dashboard',
    { params: { p: page, page_size: 10 } }
  )
  return response.data.data
}

export async function getAgentRecords(page: number) {
  const response = await api.get<ApiResponse<PagedData<AgentCommissionRecord>>>(
    '/api/user/agent/records',
    { params: { p: page, page_size: 10 } }
  )
  return response.data.data
}

export async function getAgentTopUps(page: number) {
  const response = await api.get<ApiResponse<PagedData<AgentTopUp>>>(
    '/api/user/agent/topups',
    { params: { p: page, page_size: 10 } }
  )
  return response.data.data
}
