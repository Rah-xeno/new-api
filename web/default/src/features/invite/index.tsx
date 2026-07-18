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
import { Copy, Gift, Users, UserCheck, Coins, TrendingUp } from 'lucide-react'
import { useCallback, useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'

import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from '@/components/ui/table'
import { toast } from 'sonner'
import { api } from '@/lib/api'
import { formatQuota } from '@/lib/format'

interface InviteBenefit {
    id: number
    name: string
    remark: string
    trigger_type: string
    trigger_topup_quota: number
    reward_type: string
    reward_percent: number
    reward_max_quota: number
    max_invitees_per_inviter: number
    effective_days: number
    expires_at: number
}

interface InviteeItem {
    id: number
    username: string
    display_name: string
    status: number
    created_at: number
    rewarded: boolean
    reward_quota: number
    reward_plan_name: string
}

interface InviteDashboard {
    aff_code: string
    invite_total: number
    rewarded_total: number
    aff_quota: number
    aff_history_quota: number
    benefits: InviteBenefit[]
    invitees: {
        page: number
        page_size: number
        total: number
        items: InviteeItem[]
    }
}

export function InvitePage() {
    const { t } = useTranslation()
    const [dashboard, setDashboard] = useState<InviteDashboard | null>(null)
    const [loading, setLoading] = useState(true)

    const loadDashboard = useCallback(async () => {
        setLoading(true)
        try {
            const res = await api.get('/api/user/invite/dashboard')
            setDashboard(res.data?.data)
        } catch {
            toast.error(t('Failed to load'))
        } finally {
            setLoading(false)
        }
    }, [t])

    useEffect(() => {
        loadDashboard()
    }, [loadDashboard])

    const copyAffLink = () => {
        if (!dashboard?.aff_code) return
        const link = `${window.location.origin}/register?aff=${dashboard.aff_code}`
        navigator.clipboard.writeText(link).then(() => {
            toast.success(t('Copied'))
        })
    }

    if (loading) {
        return (
            <div className='flex items-center justify-center py-20'>
                <div className='text-muted-foreground text-sm'>{t('Loading...')}</div>
            </div>
        )
    }

    if (!dashboard) {
        return (
            <div className='flex items-center justify-center py-20'>
                <div className='text-muted-foreground text-sm'>{t('Failed to load')}</div>
            </div>
        )
    }

    const d = dashboard

    return (
        <div className='flex flex-col gap-6 p-6'>
            <div>
                <h1 className='text-2xl font-bold'>{t('Invite Program')}</h1>
                <p className='text-muted-foreground text-sm'>
                    {t('Invite friends and earn rewards when they top up or redeem codes.')}
                </p>
            </div>

            {/* Invite Code & Link */}
            <Card>
                <CardHeader>
                    <CardTitle className='flex items-center gap-2 text-lg'>
                        <Gift className='h-5 w-5' />
                        {t('Your Invite Code')}
                    </CardTitle>
                </CardHeader>
                <CardContent className='flex flex-wrap items-center gap-4'>
                    <code className='bg-muted rounded-md px-4 py-2 text-lg font-mono'>
                        {d.aff_code}
                    </code>
                    <Button variant='outline' size='sm' onClick={copyAffLink}>
                        <Copy className='mr-2 h-4 w-4' />
                        {t('Copy Invite Link')}
                    </Button>
                </CardContent>
            </Card>

            {/* Stats Cards */}
            <div className='grid grid-cols-2 gap-4 md:grid-cols-4'>
                <Card>
                    <CardContent className='flex items-center gap-3 pt-6'>
                        <Users className='text-primary h-8 w-8' />
                        <div>
                            <div className='text-muted-foreground text-xs'>{t('Invited')}</div>
                            <div className='text-2xl font-bold'>{d.invite_total}</div>
                        </div>
                    </CardContent>
                </Card>
                <Card>
                    <CardContent className='flex items-center gap-3 pt-6'>
                        <UserCheck className='text-green-600 h-8 w-8' />
                        <div>
                            <div className='text-muted-foreground text-xs'>{t('Rewarded')}</div>
                            <div className='text-2xl font-bold'>{d.rewarded_total}</div>
                        </div>
                    </CardContent>
                </Card>
                <Card>
                    <CardContent className='flex items-center gap-3 pt-6'>
                        <Coins className='text-amber-500 h-8 w-8' />
                        <div>
                            <div className='text-muted-foreground text-xs'>{t('Available')}</div>
                            <div className='text-2xl font-bold'>{formatQuota(d.aff_quota)}</div>
                        </div>
                    </CardContent>
                </Card>
                <Card>
                    <CardContent className='flex items-center gap-3 pt-6'>
                        <TrendingUp className='text-blue-500 h-8 w-8' />
                        <div>
                            <div className='text-muted-foreground text-xs'>{t('Total Earned')}</div>
                            <div className='text-2xl font-bold'>{formatQuota(d.aff_history_quota)}</div>
                        </div>
                    </CardContent>
                </Card>
            </div>

            {/* Active Benefits */}
            {d.benefits.length > 0 && (
                <Card>
                    <CardHeader>
                        <CardTitle className='text-lg'>{t('Active Benefits')}</CardTitle>
                    </CardHeader>
                    <CardContent>
                        <div className='flex flex-wrap gap-3'>
                            {d.benefits.map((b) => (
                                <Card key={b.id} className='min-w-[200px] flex-1'>
                                    <CardContent className='pt-4'>
                                        <div className='font-semibold'>{b.name}</div>
                                        {b.remark && (
                                            <div className='text-muted-foreground text-xs'>
                                                {b.remark}
                                            </div>
                                        )}
                                        <div className='mt-2 flex flex-wrap gap-1'>
                                            <Badge variant='outline'>
                                                {b.trigger_type === 'redemption'
                                                    ? t('Redemption')
                                                    : b.trigger_type === 'topup'
                                                        ? t('Top-up')
                                                        : b.trigger_type}
                                            </Badge>
                                            <Badge variant='secondary' className='text-green-600'>
                                                {b.reward_percent}%
                                            </Badge>
                                            {b.effective_days > 0 && (
                                                <Badge variant='outline'>
                                                    {b.effective_days}d
                                                </Badge>
                                            )}
                                        </div>
                                    </CardContent>
                                </Card>
                            ))}
                        </div>
                    </CardContent>
                </Card>
            )}

            {/* Invitees List */}
            <Card>
                <CardHeader>
                    <CardTitle className='text-lg'>
                        {t('Invitees')} ({d.invite_total})
                    </CardTitle>
                </CardHeader>
                <CardContent>
                    {d.invitees.items.length === 0 ? (
                        <div className='text-muted-foreground py-8 text-center text-sm'>
                            {t('No invitees yet. Share your invite link to start earning!')}
                        </div>
                    ) : (
                        <Table>
                            <TableHeader>
                                <TableRow>
                                    <TableHead>{t('User')}</TableHead>
                                    <TableHead>{t('Status')}</TableHead>
                                    <TableHead>{t('Reward')}</TableHead>
                                </TableRow>
                            </TableHeader>
                            <TableBody>
                                {d.invitees.items.map((inv) => (
                                    <TableRow key={inv.id}>
                                        <TableCell>
                                            <div className='font-medium'>
                                                {inv.display_name || inv.username}
                                            </div>
                                        </TableCell>
                                        <TableCell>
                                            {inv.status === 1 ? (
                                                <Badge variant='default' className='bg-green-600'>
                                                    {t('Active')}
                                                </Badge>
                                            ) : (
                                                <Badge variant='secondary'>
                                                    {t('Disabled')}
                                                </Badge>
                                            )}
                                        </TableCell>
                                        <TableCell>
                                            {inv.rewarded ? (
                                                <div>
                                                    <Badge className='bg-green-600'>
                                                        {formatQuota(inv.reward_quota)}
                                                    </Badge>
                                                    {inv.reward_plan_name && (
                                                        <div className='text-muted-foreground mt-1 text-xs'>
                                                            {inv.reward_plan_name}
                                                        </div>
                                                    )}
                                                </div>
                                            ) : (
                                                <span className='text-muted-foreground text-sm'>
                                                    {t('Pending')}
                                                </span>
                                            )}
                                        </TableCell>
                                    </TableRow>
                                ))}
                            </TableBody>
                        </Table>
                    )}
                </CardContent>
            </Card>
        </div>
    )
}
