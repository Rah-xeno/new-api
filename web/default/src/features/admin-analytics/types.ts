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
export type AnalyticsNamedCount = {
  name: string
  count: number
}

export type AnalyticsChannelRate = {
  channel_id: number
  channel_name: string
  success_count: number
  failure_count: number
  success_rate: number
  failure_rate: number
}

export type AnalyticsChannelQuota = {
  channel_id: number
  channel_name: string
  quota: number
}

export type AnalyticsChannelLatency = {
  channel_id: number
  channel_name: string
  request_count: number
  avg_use_time: number
  p95_use_time: number
  p99_use_time: number
}

export type AnalyticsUserQuota = {
  user_id: number
  username: string
  quota: number
}

export type AnalyticsGroupStat = {
  group_name: string
  consume_quota: number
  topup_money: number
  topup_order_count: number
}

export type AnalyticsDayQuota = {
  day: string
  quota: number
}

export type AnalyticsDayMoney = {
  day: string
  money: number
}

export type AnalyticsDayUserStat = {
  day: string
  new_user_count: number
  paid_user_count: number
  active_user_count: number
}

export type AdminAnalyticsReport = {
  start_timestamp: number
  end_timestamp: number
  summary: {
    total_request_count: number
    total_success_count: number
    total_failure_count: number
    total_consume_quota: number
    total_topup_money: number
    total_active_users: number
    new_user_total: number
    paid_user_total: number
    cache_hit_rate: number
    cache_saved_quota: number
    total_redeemed_quota: number
  }
  overall: {
    success_count: number
    failure_count: number
    success_rate: number
    failure_rate: number
  }
  model_ranking: AnalyticsNamedCount[]
  channel_rates: AnalyticsChannelRate[]
  channel_consumptions: AnalyticsChannelQuota[]
  channel_latencies: AnalyticsChannelLatency[]
  user_consumptions: AnalyticsUserQuota[]
  group_stats: AnalyticsGroupStat[]
  daily_consumes: AnalyticsDayQuota[]
  daily_topups: AnalyticsDayMoney[]
  daily_user_stats: AnalyticsDayUserStat[]
  referral_overview: {
    invite_topup_money: number
    invite_order_count: number
    invite_paid_user_count: number
    invite_reward_quota: number
    invite_reward_count: number
    agent_topup_money: number
    agent_order_count: number
    agent_paid_user_count: number
    agent_commission_amount: number
    agent_commission_count: number
    agent_net_money: number
  }
  cache_overview: {
    prompt_tokens: number
    input_tokens_total: number
    cache_hit_tokens: number
    cache_creation_tokens: number
    cache_creation_tokens_5m: number
    cache_creation_tokens_1h: number
    cache_write_tokens: number
    cache_hit_rate: number
    cache_saved_quota: number
  }
}

export type AnalyticsDateRange = {
  start: string
  end: string
}
