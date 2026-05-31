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
  CommissionRecord,
  DistributionOverview,
  DistributionSettings,
  PageResponse,
  UserReferral,
} from './types'

export async function getDistributionOverview() {
  const res = await api.get<{ success: boolean; data: DistributionOverview }>(
    '/api/user/distribution/overview'
  )
  return res.data
}

export async function createDistributionInviteLink() {
  const res = await api.post('/api/user/distribution/invite-link')
  return res.data
}

export async function getDistributionInvites(page = 1, pageSize = 10) {
  const res = await api.get<{ success: boolean; data: PageResponse<UserReferral> }>(
    `/api/user/distribution/invites?p=${page}&page_size=${pageSize}`
  )
  return res.data
}

export async function getDistributionCommissionRecords(page = 1, pageSize = 10) {
  const res = await api.get<{
    success: boolean
    data: PageResponse<CommissionRecord>
  }>(`/api/user/distribution/commission-records?p=${page}&page_size=${pageSize}`)
  return res.data
}

export async function settleDistributionCommission() {
  const res = await api.post('/api/user/distribution/settle')
  return res.data
}

export async function getAdminDistributionSettings() {
  const res = await api.get<{ success: boolean; data: DistributionSettings }>(
    '/api/admin/distribution/settings'
  )
  return res.data
}

export async function updateAdminDistributionSettings(
  settings: DistributionSettings
) {
  const res = await api.put('/api/admin/distribution/settings', settings)
  return res.data
}

export async function getAdminCommissionRecords(
  page = 1,
  pageSize = 10,
  referrerUserId?: number
) {
  const params = new URLSearchParams({
    p: String(page),
    page_size: String(pageSize),
  })
  if (referrerUserId) params.set('referrer_user_id', String(referrerUserId))
  const res = await api.get<{
    success: boolean
    data: PageResponse<CommissionRecord>
  }>(`/api/admin/distribution/commission-records?${params.toString()}`)
  return res.data
}

export async function getAdminReferrals(
  referrerUserId: number,
  page = 1,
  pageSize = 10
) {
  const params = new URLSearchParams({
    p: String(page),
    page_size: String(pageSize),
    referrer_user_id: String(referrerUserId),
  })
  const res = await api.get<{ success: boolean; data: PageResponse<UserReferral> }>(
    `/api/admin/distribution/referrals?${params.toString()}`
  )
  return res.data
}

export async function manualGrantCommission(payload: {
  referrer_user_id: number
  referred_user_id?: number
  amount: number
  reason: string
}) {
  const res = await api.post('/api/admin/distribution/manual-grant', payload)
  return res.data
}

export async function adjustCommission(payload: {
  user_id: number
  action: 'increase' | 'decrease' | 'clear'
  amount: number
  reason: string
}) {
  const res = await api.post('/api/admin/distribution/commission-adjust', payload)
  return res.data
}
