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
export type PagedData<T> = {
  page: number
  page_size: number
  total: number
  items: T[]
}

export type AgentInvitee = {
  invitee_id: number
  username: string
  display_name: string
  email: string
  status: number
  successful_topup_count: number
  successful_topup_amount: number
  last_topup_at: number
  total_commission_amount: number
}

export type AgentDashboard = {
  agent_enabled: boolean
  agent_portal_visible: boolean
  agent_use_default_rates: boolean
  effective_first_topup_rate: number
  effective_repeat_topup_rate: number
  agent_commission_balance: number
  agent_commission_total: number
  agent_commission_withdrawn: number
  aff_code: string
  stats: {
    invitee_total: number
    paid_invitee_total: number
    total_topup_amount: number
    commission_order_total: number
  }
  invitees: PagedData<AgentInvitee>
}

export type AgentCommissionRecord = {
  id: number
  source_type: 'topup' | 'redemption' | 'withdraw'
  invitee_name: string
  source_trade_no: string
  source_amount: number
  commission_rate: number
  commission_amount: number
  is_first_topup: boolean
  created_at: number
}

export type AgentTopUp = {
  top_up_id: number
  invitee_id: number
  display_name: string
  trade_no: string
  payment_method: string
  payment_amount: number
  complete_time: number
  commission_amount: number
  commission_rate: number
  has_commission: boolean
  is_first_topup_order: boolean
}
