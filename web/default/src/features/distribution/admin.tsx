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
import { useEffect, useState } from 'react'
import { Save } from 'lucide-react'
import { toast } from 'sonner'
import { useTranslation } from 'react-i18next'
import { SectionPageLayout } from '@/components/layout'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Switch } from '@/components/ui/switch'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Textarea } from '@/components/ui/textarea'
import { formatQuota, formatTimestampToDate } from '@/lib/format'
import {
  adjustCommission,
  getAdminCommissionRecords,
  getAdminDistributionSettings,
  getAdminReferrals,
  updateAdminDistributionSettings,
} from './api'
import type {
  CommissionRecord,
  DistributionSettings,
  UserReferral,
} from './types'

const defaultSettings: DistributionSettings = {
  enabled: false,
  reward_rate: 10,
  referral_limit: 0,
  reward_mode: 1,
  fixed_reward_amount: 0,
  freeze_days: 7,
  referred_user_reward: 0,
  min_settlement_threshold: 0,
  admin_recharge_trigger: true,
}

export function AdminDistribution() {
  const { t } = useTranslation()
  const [settings, setSettings] = useState<DistributionSettings>(defaultSettings)
  const [records, setRecords] = useState<CommissionRecord[]>([])
  const [referrals, setReferrals] = useState<UserReferral[]>([])
  const [referrerUserId, setReferrerUserId] = useState('')
  const [adjustUserId, setAdjustUserId] = useState('')
  const [adjustAction, setAdjustAction] = useState<'increase' | 'decrease' | 'clear'>('increase')
  const [adjustAmount, setAdjustAmount] = useState('')
  const [adjustReason, setAdjustReason] = useState('')

  const loadSettings = async () => {
    const res = await getAdminDistributionSettings()
    if (res.success) setSettings(res.data)
  }

  const loadRecords = async () => {
    const res = await getAdminCommissionRecords(1, 20)
    if (res.success) setRecords(res.data.items ?? [])
  }

  useEffect(() => {
    loadSettings()
    loadRecords()
  }, [])

  const saveSettings = async () => {
    const res = await updateAdminDistributionSettings(settings)
    if (res.success) toast.success(t('Saved'))
    else toast.error(res.message || t('Save failed'))
  }

  const searchReferrals = async () => {
    const id = Number(referrerUserId)
    if (!id) return
    const res = await getAdminReferrals(id, 1, 20)
    if (res.success) setReferrals(res.data.items ?? [])
  }

  const submitAdjustment = async () => {
    const userId = Number(adjustUserId)
    const amount = adjustAction === 'clear' ? 0 : Number(adjustAmount)
    if (!userId) return
    const res = await adjustCommission({
      user_id: userId,
      action: adjustAction,
      amount,
      reason: adjustReason,
    })
    if (res.success) {
      toast.success(t('Saved'))
      await loadRecords()
    } else {
      toast.error(res.message || t('Save failed'))
    }
  }

  return (
    <SectionPageLayout>
      <SectionPageLayout.Title>{t('Distribution Management')}</SectionPageLayout.Title>
      <SectionPageLayout.Description>
        {t('Configure distribution rules and manage commission accounts.')}
      </SectionPageLayout.Description>
      <SectionPageLayout.Content>
        <Tabs defaultValue='settings' className='mx-auto w-full max-w-7xl'>
          <TabsList>
            <TabsTrigger value='settings'>{t('Settings')}</TabsTrigger>
            <TabsTrigger value='records'>{t('Commission records')}</TabsTrigger>
            <TabsTrigger value='referrals'>{t('Referrals')}</TabsTrigger>
            <TabsTrigger value='adjust'>{t('Adjust commission')}</TabsTrigger>
          </TabsList>

          <TabsContent value='settings' className='mt-4'>
            <Card>
              <CardHeader>
                <CardTitle>{t('Distribution settings')}</CardTitle>
              </CardHeader>
              <CardContent className='grid gap-4 md:grid-cols-2'>
                <div className='flex items-center justify-between rounded-md border p-3'>
                  <Label>{t('Enable distribution')}</Label>
                  <Switch
                    checked={settings.enabled}
                    onCheckedChange={(enabled) =>
                      setSettings((prev) => ({ ...prev, enabled }))
                    }
                  />
                </div>
                <div className='flex items-center justify-between rounded-md border p-3'>
                  <Label>{t('Admin recharge triggers commission')}</Label>
                  <Switch
                    checked={settings.admin_recharge_trigger}
                    onCheckedChange={(admin_recharge_trigger) =>
                      setSettings((prev) => ({ ...prev, admin_recharge_trigger }))
                    }
                  />
                </div>
                {[
                  ['Reward rate', 'reward_rate'],
                  ['Referral limit', 'referral_limit'],
                  ['Reward mode', 'reward_mode'],
                  ['Fixed reward amount', 'fixed_reward_amount'],
                  ['Freeze days', 'freeze_days'],
                  ['Referred user reward', 'referred_user_reward'],
                  ['Minimum settlement threshold', 'min_settlement_threshold'],
                ].map(([label, key]) => (
                  <div key={key} className='space-y-2'>
                    <Label>{t(label)}</Label>
                    <Input
                      type='number'
                      value={String(settings[key as keyof DistributionSettings] ?? 0)}
                      onChange={(event) =>
                        setSettings((prev) => ({
                          ...prev,
                          [key]: Number(event.target.value),
                        }))
                      }
                    />
                  </div>
                ))}
                <div className='md:col-span-2'>
                  <Button onClick={saveSettings}>
                    <Save className='size-4' />
                    {t('Save')}
                  </Button>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value='records' className='mt-4'>
            <RecordsTable records={records} />
          </TabsContent>

          <TabsContent value='referrals' className='mt-4'>
            <Card>
              <CardHeader>
                <CardTitle>{t('Referral lookup')}</CardTitle>
              </CardHeader>
              <CardContent className='space-y-4'>
                <div className='flex gap-2'>
                  <Input
                    placeholder={t('Referrer user ID')}
                    value={referrerUserId}
                    onChange={(event) => setReferrerUserId(event.target.value)}
                  />
                  <Button onClick={searchReferrals}>{t('Search')}</Button>
                </div>
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>{t('Referrer')}</TableHead>
                      <TableHead>{t('Referred user')}</TableHead>
                      <TableHead>{t('Source')}</TableHead>
                      <TableHead>{t('Bound at')}</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {referrals.map((item) => (
                      <TableRow key={item.id}>
                        <TableCell>{item.referrer_user_id}</TableCell>
                        <TableCell>{item.referred_user_id}</TableCell>
                        <TableCell>{item.source}</TableCell>
                        <TableCell>{formatTimestampToDate(item.bound_at)}</TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value='adjust' className='mt-4'>
            <Card>
              <CardHeader>
                <CardTitle>{t('Adjust commission')}</CardTitle>
              </CardHeader>
              <CardContent className='grid gap-4 md:grid-cols-2'>
                <div className='space-y-2'>
                  <Label>{t('User ID')}</Label>
                  <Input value={adjustUserId} onChange={(e) => setAdjustUserId(e.target.value)} />
                </div>
                <div className='space-y-2'>
                  <Label>{t('Action')}</Label>
                  <select
                    className='h-9 rounded-md border bg-background px-3 text-sm'
                    value={adjustAction}
                    onChange={(e) =>
                      setAdjustAction(e.target.value as 'increase' | 'decrease' | 'clear')
                    }
                  >
                    <option value='increase'>{t('Increase')}</option>
                    <option value='decrease'>{t('Decrease')}</option>
                    <option value='clear'>{t('Clear unsettled')}</option>
                  </select>
                </div>
                <div className='space-y-2'>
                  <Label>{t('Amount')}</Label>
                  <Input
                    type='number'
                    disabled={adjustAction === 'clear'}
                    value={adjustAmount}
                    onChange={(e) => setAdjustAmount(e.target.value)}
                  />
                </div>
                <div className='space-y-2 md:col-span-2'>
                  <Label>{t('Reason')}</Label>
                  <Textarea value={adjustReason} onChange={(e) => setAdjustReason(e.target.value)} />
                </div>
                <div className='md:col-span-2'>
                  <Button onClick={submitAdjustment}>{t('Save')}</Button>
                </div>
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>
      </SectionPageLayout.Content>
    </SectionPageLayout>
  )
}

function RecordsTable({ records }: { records: CommissionRecord[] }) {
  const { t } = useTranslation()
  return (
    <Card>
      <CardHeader>
        <CardTitle>{t('Commission records')}</CardTitle>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>{t('Referrer')}</TableHead>
              <TableHead>{t('Referred user')}</TableHead>
              <TableHead>{t('Commission')}</TableHead>
              <TableHead>{t('Source')}</TableHead>
              <TableHead>{t('Created at')}</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {records.map((item) => (
              <TableRow key={item.id}>
                <TableCell>{item.referrer_user_id}</TableCell>
                <TableCell>{item.referred_user_id || '-'}</TableCell>
                <TableCell>{formatQuota(item.commission_amount)}</TableCell>
                <TableCell>{item.source}</TableCell>
                <TableCell>{formatTimestampToDate(item.created_at)}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  )
}
