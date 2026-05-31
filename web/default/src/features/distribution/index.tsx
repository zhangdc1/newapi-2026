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
import { useEffect, useMemo, useState } from 'react'
import { Copy, RefreshCw, WalletCards } from 'lucide-react'
import { toast } from 'sonner'
import { useTranslation } from 'react-i18next'
import { SectionPageLayout } from '@/components/layout'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { formatQuota, formatTimestampToDate } from '@/lib/format'
import {
  getDistributionCommissionRecords,
  getDistributionInvites,
  getDistributionOverview,
  settleDistributionCommission,
} from './api'
import type {
  CommissionRecord,
  DistributionOverview,
  UserReferral,
} from './types'

function absoluteInviteLink(path: string) {
  if (!path) return ''
  if (/^https?:\/\//.test(path)) return path
  return `${window.location.origin}${path}`
}

function statusLabel(status: number) {
  if (status === 3) return 'Settled'
  if (status === 4) return 'Invalid'
  return 'Pending'
}

export function Distribution() {
  const { t } = useTranslation()
  const [overview, setOverview] = useState<DistributionOverview | null>(null)
  const [invites, setInvites] = useState<UserReferral[]>([])
  const [records, setRecords] = useState<CommissionRecord[]>([])
  const [loading, setLoading] = useState(true)
  const [settling, setSettling] = useState(false)

  const inviteLink = useMemo(
    () => absoluteInviteLink(overview?.invite_link ?? ''),
    [overview?.invite_link]
  )

  const loadData = async () => {
    setLoading(true)
    try {
      const [overviewRes, invitesRes, recordsRes] = await Promise.all([
        getDistributionOverview(),
        getDistributionInvites(1, 8),
        getDistributionCommissionRecords(1, 8),
      ])
      if (overviewRes.success) setOverview(overviewRes.data)
      if (invitesRes.success) setInvites(invitesRes.data.items ?? [])
      if (recordsRes.success) setRecords(recordsRes.data.items ?? [])
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadData()
  }, [])

  const copyInviteLink = async () => {
    await navigator.clipboard.writeText(inviteLink)
    toast.success(t('Copied'))
  }

  const settle = async () => {
    setSettling(true)
    try {
      const res = await settleDistributionCommission()
      if (res.success) {
        toast.success(t('Commission settled successfully'))
        await loadData()
      } else {
        toast.error(res.message || t('Settlement failed'))
      }
    } finally {
      setSettling(false)
    }
  }

  return (
    <SectionPageLayout>
      <SectionPageLayout.Title>{t('Distribution Center')}</SectionPageLayout.Title>
      <SectionPageLayout.Description>
        {t('Manage invite links, commission records, and settlement.')}
      </SectionPageLayout.Description>
      <SectionPageLayout.Actions>
        <Button variant='outline' size='sm' onClick={loadData} disabled={loading}>
          <RefreshCw className='size-4' />
          {t('Refresh')}
        </Button>
      </SectionPageLayout.Actions>
      <SectionPageLayout.Content>
        <div className='mx-auto flex w-full max-w-7xl flex-col gap-4'>
          <div className='grid gap-4 md:grid-cols-2 xl:grid-cols-4'>
            {[
              ['Total earned', overview?.total_earned],
              ['Pending commission', overview?.pending_commission],
              ['Available commission', overview?.available_commission],
              ['Settled commission', overview?.settled_commission],
            ].map(([label, value]) => (
              <Card key={label as string}>
                <CardHeader className='pb-2'>
                  <CardTitle className='text-sm font-medium text-muted-foreground'>
                    {t(label as string)}
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className='text-2xl font-semibold'>
                    {formatQuota(Number(value ?? 0))}
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>

          <Card>
            <CardHeader>
              <CardTitle>{t('Invite link')}</CardTitle>
            </CardHeader>
            <CardContent className='flex flex-col gap-3 md:flex-row'>
              <Input value={inviteLink} readOnly />
              <Button onClick={copyInviteLink} disabled={!inviteLink}>
                <Copy className='size-4' />
                {t('Copy')}
              </Button>
              <Button onClick={settle} disabled={settling}>
                <WalletCards className='size-4' />
                {t('Settle commission')}
              </Button>
            </CardContent>
          </Card>

          <div className='grid gap-4 xl:grid-cols-2'>
            <Card>
              <CardHeader>
                <CardTitle>{t('Invited users')}</CardTitle>
              </CardHeader>
              <CardContent>
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>{t('User ID')}</TableHead>
                      <TableHead>{t('Source')}</TableHead>
                      <TableHead>{t('Bound at')}</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {invites.map((item) => (
                      <TableRow key={item.id}>
                        <TableCell>{item.referred_user_id}</TableCell>
                        <TableCell>{item.source}</TableCell>
                        <TableCell>{formatTimestampToDate(item.bound_at)}</TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>{t('Commission records')}</CardTitle>
              </CardHeader>
              <CardContent>
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>{t('Amount')}</TableHead>
                      <TableHead>{t('Source')}</TableHead>
                      <TableHead>{t('Status')}</TableHead>
                      <TableHead>{t('Created at')}</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {records.map((item) => (
                      <TableRow key={item.id}>
                        <TableCell>{formatQuota(item.commission_amount)}</TableCell>
                        <TableCell>{item.source}</TableCell>
                        <TableCell>
                          <Badge variant='outline'>{t(statusLabel(item.status))}</Badge>
                        </TableCell>
                        <TableCell>{formatTimestampToDate(item.created_at)}</TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </CardContent>
            </Card>
          </div>
        </div>
      </SectionPageLayout.Content>
    </SectionPageLayout>
  )
}
