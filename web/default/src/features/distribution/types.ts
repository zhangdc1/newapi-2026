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
export type PageResponse<T> = {
  page: number
  page_size: number
  total: number
  items: T[]
}

export type DistributionOverview = {
  enabled: boolean
  dist_id: string
  invite_link: string
  total_earned: number
  pending_commission: number
  available_commission: number
  settled_commission: number
  invite_count: number
  effective_invite_count: number
  min_settlement_threshold: number
}

export type DistributionSettings = {
  id?: number
  enabled: boolean
  reward_rate: number
  referral_limit: number
  reward_mode: number
  fixed_reward_amount: number
  freeze_days: number
  referred_user_reward: number
  min_settlement_threshold: number
  admin_recharge_trigger: boolean
  updated_by?: number
  updated_at?: number
}

export type UserReferral = {
  id: number
  referrer_user_id: number
  referred_user_id: number
  status: number
  bound_at: number
  source: string
  referrer_dist_id: string
}

export type CommissionRecord = {
  id: number
  referrer_user_id: number
  referred_user_id: number
  source: string
  source_id: string
  increased_amount: number
  commission_type: number
  commission_amount: number
  status: number
  freeze_until: number
  created_at: number
  settled_at: number
  admin_id: number
  reason: string
  before_balance: number
  after_balance: number
}
