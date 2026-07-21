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
import { useQuery } from '@tanstack/react-query'
import { BadgePercent, Copy, Users, WalletCards } from 'lucide-react'
import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'

import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'

import { getAgentDashboard, getAgentRecords, getAgentTopUps } from './api'

function formatMoney(cents: number | undefined) {
  return `¥${((cents ?? 0) / 100).toFixed(2)}`
}

function formatDate(timestamp: number | undefined) {
  if (!timestamp) return '-'
  return new Date(timestamp * 1000).toLocaleString()
}

function Pager(props: {
  page: number
  total: number
  onChange: (page: number) => void
}) {
  const { t } = useTranslation()
  const pages = Math.max(1, Math.ceil(props.total / 10))
  return (
    <div className='flex items-center justify-end gap-2 pt-4'>
      <Button
        variant='outline'
        size='sm'
        disabled={props.page <= 1}
        onClick={() => props.onChange(props.page - 1)}
      >
        {t('Previous')}
      </Button>
      <span className='text-muted-foreground text-sm'>
        {props.page} / {pages}
      </span>
      <Button
        variant='outline'
        size='sm'
        disabled={props.page >= pages}
        onClick={() => props.onChange(props.page + 1)}
      >
        {t('Next')}
      </Button>
    </div>
  )
}

export function AgentDistributionPage() {
  const { t } = useTranslation()
  const [inviteePage, setInviteePage] = useState(1)
  const [recordPage, setRecordPage] = useState(1)
  const [topUpPage, setTopUpPage] = useState(1)

  const dashboardQuery = useQuery({
    queryKey: ['agent-distribution', 'dashboard', inviteePage],
    queryFn: () => getAgentDashboard(inviteePage),
  })
  const recordsQuery = useQuery({
    queryKey: ['agent-distribution', 'records', recordPage],
    queryFn: () => getAgentRecords(recordPage),
  })
  const topUpsQuery = useQuery({
    queryKey: ['agent-distribution', 'topups', topUpPage],
    queryFn: () => getAgentTopUps(topUpPage),
  })

  const dashboard = dashboardQuery.data
  const inviteUrl = dashboard?.aff_code
    ? `${window.location.origin}/sign-up?aff=${dashboard.aff_code}`
    : ''

  async function copyInviteLink() {
    if (!inviteUrl) return
    await navigator.clipboard.writeText(inviteUrl)
    toast.success(t('Copied'))
  }

  if (dashboardQuery.isLoading) {
    return (
      <div className='text-muted-foreground p-8 text-sm'>{t('Loading...')}</div>
    )
  }

  if (!dashboard) {
    return (
      <div className='text-muted-foreground p-8 text-sm'>
        {t('Agent distribution is not enabled for this account.')}
      </div>
    )
  }

  const stats = [
    {
      label: t('Available commission'),
      value: formatMoney(dashboard.agent_commission_balance),
      icon: WalletCards,
    },
    {
      label: t('Total commission'),
      value: formatMoney(dashboard.agent_commission_total),
      icon: BadgePercent,
    },
    {
      label: t('Agent customers'),
      value: String(dashboard.stats.invitee_total),
      icon: Users,
    },
  ]

  return (
    <div className='space-y-6 p-4 md:p-6'>
      <div>
        <h1 className='text-2xl font-semibold'>{t('Agent Distribution')}</h1>
        <p className='text-muted-foreground text-sm'>
          {t(
            'Invite customers and earn commission from their successful top-ups.'
          )}
        </p>
      </div>

      {!dashboard.agent_enabled && (
        <div className='rounded-md border border-amber-300 bg-amber-50 p-3 text-sm text-amber-900 dark:bg-amber-950/30 dark:text-amber-200'>
          {t(
            'Agent access is disabled. Historical commission data remains available.'
          )}
        </div>
      )}

      <div className='grid gap-4 md:grid-cols-3'>
        {stats.map((stat) => (
          <Card key={stat.label}>
            <CardContent className='flex items-center justify-between pt-6'>
              <div>
                <div className='text-muted-foreground text-sm'>
                  {stat.label}
                </div>
                <div className='mt-1 text-2xl font-semibold'>{stat.value}</div>
              </div>
              <stat.icon
                className='text-muted-foreground h-6 w-6'
                aria-hidden='true'
              />
            </CardContent>
          </Card>
        ))}
      </div>

      <Card>
        <CardHeader>
          <CardTitle>{t('Invitation link')}</CardTitle>
        </CardHeader>
        <CardContent className='space-y-3'>
          <div className='flex flex-col gap-2 sm:flex-row'>
            <code className='bg-muted min-w-0 flex-1 truncate rounded-md px-3 py-2 text-sm'>
              {inviteUrl}
            </code>
            <Button onClick={copyInviteLink} disabled={!inviteUrl}>
              <Copy className='mr-2 h-4 w-4' aria-hidden='true' />
              {t('Copy')}
            </Button>
          </div>
          <div className='text-muted-foreground text-sm'>
            {t('First top-up rate')}: {dashboard.effective_first_topup_rate}% ·{' '}
            {t('Repeat top-up rate')}: {dashboard.effective_repeat_topup_rate}%
          </div>
        </CardContent>
      </Card>

      <Tabs defaultValue='customers'>
        <TabsList>
          <TabsTrigger value='customers'>{t('Agent customers')}</TabsTrigger>
          <TabsTrigger value='topups'>{t('Top-up records')}</TabsTrigger>
          <TabsTrigger value='commissions'>
            {t('Commission records')}
          </TabsTrigger>
        </TabsList>

        <TabsContent value='customers'>
          <Card>
            <CardContent className='pt-6'>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>{t('User')}</TableHead>
                    <TableHead>{t('Successful top-ups')}</TableHead>
                    <TableHead>{t('Top-up amount')}</TableHead>
                    <TableHead>{t('Commission')}</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {dashboard.invitees.items.map((item) => (
                    <TableRow key={item.invitee_id}>
                      <TableCell>
                        <div className='font-medium'>{item.display_name}</div>
                        <div className='text-muted-foreground text-xs'>
                          {item.email || item.username}
                        </div>
                      </TableCell>
                      <TableCell>{item.successful_topup_count}</TableCell>
                      <TableCell>
                        {formatMoney(item.successful_topup_amount)}
                      </TableCell>
                      <TableCell>
                        {formatMoney(item.total_commission_amount)}
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
              <Pager
                page={inviteePage}
                total={dashboard.invitees.total}
                onChange={setInviteePage}
              />
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value='topups'>
          <Card>
            <CardContent className='pt-6'>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>{t('User')}</TableHead>
                    <TableHead>{t('Payment time')}</TableHead>
                    <TableHead>{t('Amount')}</TableHead>
                    <TableHead>{t('Commission')}</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {(topUpsQuery.data?.items ?? []).map((item) => (
                    <TableRow key={item.top_up_id}>
                      <TableCell>{item.display_name}</TableCell>
                      <TableCell>{formatDate(item.complete_time)}</TableCell>
                      <TableCell>{formatMoney(item.payment_amount)}</TableCell>
                      <TableCell>
                        {item.has_commission ? (
                          <Badge variant='secondary'>
                            {formatMoney(item.commission_amount)} ·{' '}
                            {item.commission_rate}%
                          </Badge>
                        ) : (
                          '-'
                        )}
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
              <Pager
                page={topUpPage}
                total={topUpsQuery.data?.total ?? 0}
                onChange={setTopUpPage}
              />
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value='commissions'>
          <Card>
            <CardContent className='pt-6'>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>{t('Source')}</TableHead>
                    <TableHead>{t('User')}</TableHead>
                    <TableHead>{t('Time')}</TableHead>
                    <TableHead>{t('Commission')}</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {(recordsQuery.data?.items ?? []).map((item) => (
                    <TableRow key={item.id}>
                      <TableCell>
                        <Badge variant='outline'>{t(item.source_type)}</Badge>
                      </TableCell>
                      <TableCell>{item.invitee_name || '-'}</TableCell>
                      <TableCell>{formatDate(item.created_at)}</TableCell>
                      <TableCell>
                        {formatMoney(item.commission_amount)}
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
              <Pager
                page={recordPage}
                total={recordsQuery.data?.total ?? 0}
                onChange={setRecordPage}
              />
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}
