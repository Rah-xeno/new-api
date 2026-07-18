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
import { Plus, Pencil, Trash2, Power, PowerOff } from 'lucide-react'
import { useCallback, useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'

import { Button } from '@/components/ui/button'
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
    DialogFooter,
} from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Switch } from '@/components/ui/switch'
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from '@/components/ui/select'
import { useToast } from '@/hooks/use-toast'
import { api } from '@/lib/api'
import { formatQuota } from '@/lib/format'

interface InvitePlan {
    id: number
    name: string
    remark: string
    enabled: boolean
    priority: number
    trigger_type: string
    trigger_topup_quota: number
    reward_type: string
    reward_percent: number
    reward_max_quota: number
    max_invitees_per_inviter: number
    effective_days: number
    expires_at: number
    created_at: number
    updated_at: number
}

const emptyPlan = (): Partial<InvitePlan> => ({
    name: '',
    remark: '',
    enabled: true,
    priority: 0,
    trigger_type: 'redemption',
    trigger_topup_quota: 0,
    reward_type: 'quota',
    reward_percent: 10,
    reward_max_quota: 0,
    max_invitees_per_inviter: 0,
    effective_days: 0,
    expires_at: 0,
})

export function InvitePlansPage() {
    const { t } = useTranslation()
    const { toast } = useToast()
    const [plans, setPlans] = useState<InvitePlan[]>([])
    const [loading, setLoading] = useState(true)
    const [dialogOpen, setDialogOpen] = useState(false)
    const [editing, setEditing] = useState<Partial<InvitePlan>>(emptyPlan())
    const [saving, setSaving] = useState(false)

    const loadPlans = useCallback(async () => {
        setLoading(true)
        try {
            const res = await api.get('/api/invite-plan/admin/plans')
            setPlans(res.data?.data || [])
        } catch {
            toast({ title: t('Failed to load'), variant: 'destructive' })
        } finally {
            setLoading(false)
        }
    }, [t, toast])

    useEffect(() => {
        loadPlans()
    }, [loadPlans])

    const openCreate = () => {
        setEditing(emptyPlan())
        setDialogOpen(true)
    }

    const openEdit = (plan: InvitePlan) => {
        setEditing({ ...plan })
        setDialogOpen(true)
    }

    const handleSave = async () => {
        if (!editing.name?.trim()) return
        setSaving(true)
        try {
            if (editing.id) {
                await api.put(`/api/invite-plan/admin/plans/${editing.id}`, {
                    plan: editing,
                })
                toast({ title: t('Updated') })
            } else {
                await api.post('/api/invite-plan/admin/plans', { plan: editing })
                toast({ title: t('Created') })
            }
            setDialogOpen(false)
            loadPlans()
        } catch {
            toast({ title: t('Failed to save'), variant: 'destructive' })
        } finally {
            setSaving(false)
        }
    }

    const handleDelete = async (id: number) => {
        if (!confirm(t('Delete this plan?'))) return
        try {
            await api.delete(`/api/invite-plan/admin/plans/${id}`)
            toast({ title: t('Deleted') })
            loadPlans()
        } catch {
            toast({ title: t('Failed to delete'), variant: 'destructive' })
        }
    }

    const toggleEnabled = async (plan: InvitePlan) => {
        try {
            await api.patch(`/api/invite-plan/admin/plans/${plan.id}`, {
                enabled: !plan.enabled,
            })
            toast({ title: plan.enabled ? t('Disabled') : t('Enabled') })
            loadPlans()
        } catch {
            toast({ title: t('Failed'), variant: 'destructive' })
        }
    }

    const updateField = (field: string, value: unknown) => {
        setEditing((prev) => ({ ...prev, [field]: value }))
    }

    return (
        <div className='flex flex-col gap-6 p-6'>
            <div className='flex items-center justify-between'>
                <div>
                    <h1 className='text-2xl font-bold'>{t('Invite Plans')}</h1>
                    <p className='text-muted-foreground text-sm'>
                        {t('Configure reward rules for invitees who top up or redeem codes.')}
                    </p>
                </div>
                <Button onClick={openCreate} size='sm'>
                    <Plus className='mr-2 h-4 w-4' />
                    {t('Create Plan')}
                </Button>
            </div>

            {loading ? (
                <div className='text-muted-foreground text-sm'>{t('Loading...')}</div>
            ) : plans.length === 0 ? (
                <div className='text-muted-foreground py-12 text-center text-sm'>
                    {t('No invite plans yet.')}
                </div>
            ) : (
                <div className='rounded-md border'>
                    <Table>
                        <TableHeader>
                            <TableRow>
                                <TableHead>{t('Name')}</TableHead>
                                <TableHead>{t('Trigger')}</TableHead>
                                <TableHead>{t('Reward')}</TableHead>
                                <TableHead>{t('Status')}</TableHead>
                                <TableHead>{t('Actions')}</TableHead>
                            </TableRow>
                        </TableHeader>
                        <TableBody>
                            {plans.map((plan) => (
                                <TableRow key={plan.id}>
                                    <TableCell>
                                        <div className='font-medium'>{plan.name}</div>
                                        {plan.remark && (
                                            <div className='text-muted-foreground text-xs'>
                                                {plan.remark}
                                            </div>
                                        )}
                                    </TableCell>
                                    <TableCell>
                                        <Badge variant='outline'>
                                            {plan.trigger_type === 'redemption'
                                                ? t('Redemption')
                                                : plan.trigger_type === 'topup'
                                                    ? t('Top-up')
                                                    : plan.trigger_type}
                                        </Badge>
                                        {plan.trigger_topup_quota > 0 && (
                                            <div className='text-muted-foreground mt-1 text-xs'>
                                                {t('Min')}: {formatQuota(plan.trigger_topup_quota)}
                                            </div>
                                        )}
                                    </TableCell>
                                    <TableCell>
                                        <span className='font-semibold text-green-600'>
                                            {plan.reward_percent}%
                                        </span>
                                        {plan.reward_max_quota > 0 && (
                                            <div className='text-muted-foreground text-xs'>
                                                {t('Max')}: {formatQuota(plan.reward_max_quota)}
                                            </div>
                                        )}
                                    </TableCell>
                                    <TableCell>
                                        {plan.enabled ? (
                                            <Badge variant='default' className='bg-green-600'>
                                                {t('Enabled')}
                                            </Badge>
                                        ) : (
                                            <Badge variant='secondary'>{t('Disabled')}</Badge>
                                        )}
                                    </TableCell>
                                    <TableCell>
                                        <div className='flex gap-1'>
                                            <Button
                                                size='icon'
                                                variant='ghost'
                                                onClick={() => toggleEnabled(plan)}
                                            >
                                                {plan.enabled ? (
                                                    <PowerOff className='h-4 w-4' />
                                                ) : (
                                                    <Power className='h-4 w-4' />
                                                )}
                                            </Button>
                                            <Button
                                                size='icon'
                                                variant='ghost'
                                                onClick={() => openEdit(plan)}
                                            >
                                                <Pencil className='h-4 w-4' />
                                            </Button>
                                            <Button
                                                size='icon'
                                                variant='ghost'
                                                onClick={() => handleDelete(plan.id)}
                                            >
                                                <Trash2 className='h-4 w-4 text-destructive' />
                                            </Button>
                                        </div>
                                    </TableCell>
                                </TableRow>
                            ))}
                        </TableBody>
                    </Table>
                </div>
            )}

            <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
                <DialogContent className='max-h-[90vh] overflow-y-auto sm:max-w-lg'>
                    <DialogHeader>
                        <DialogTitle>
                            {editing.id ? t('Edit Plan') : t('Create Plan')}
                        </DialogTitle>
                    </DialogHeader>
                    <div className='flex flex-col gap-4'>
                        <div className='grid gap-2'>
                            <Label htmlFor='plan-name'>{t('Name')} *</Label>
                            <Input
                                id='plan-name'
                                value={editing.name || ''}
                                onChange={(e) => updateField('name', e.target.value)}
                                placeholder={t('e.g. Default Rebate')}
                            />
                        </div>
                        <div className='grid gap-2'>
                            <Label htmlFor='plan-remark'>{t('Remark')}</Label>
                            <Input
                                id='plan-remark'
                                value={editing.remark || ''}
                                onChange={(e) => updateField('remark', e.target.value)}
                            />
                        </div>
                        <div className='grid gap-2'>
                            <Label>{t('Trigger Type')}</Label>
                            <Select
                                value={editing.trigger_type || 'redemption'}
                                onValueChange={(v) => updateField('trigger_type', v)}
                            >
                                <SelectTrigger>
                                    <SelectValue />
                                </SelectTrigger>
                                <SelectContent>
                                    <SelectItem value='redemption'>{t('Redemption')}</SelectItem>
                                    <SelectItem value='topup'>{t('Top-up')}</SelectItem>
                                </SelectContent>
                            </Select>
                        </div>
                        <div className='grid gap-2'>
                            <Label>{t('Min Trigger Quota')} (0 = {t('no limit')})</Label>
                            <Input
                                type='number'
                                value={editing.trigger_topup_quota || 0}
                                onChange={(e) =>
                                    updateField('trigger_topup_quota', Number(e.target.value))
                                }
                            />
                        </div>
                        <div className='grid gap-2'>
                            <Label>{t('Reward Percent')} (%)</Label>
                            <Input
                                type='number'
                                value={editing.reward_percent || 0}
                                onChange={(e) =>
                                    updateField('reward_percent', Number(e.target.value))
                                }
                            />
                        </div>
                        <div className='grid gap-2'>
                            <Label>{t('Max Reward Quota')} (0 = {t('no limit')})</Label>
                            <Input
                                type='number'
                                value={editing.reward_max_quota || 0}
                                onChange={(e) =>
                                    updateField('reward_max_quota', Number(e.target.value))
                                }
                            />
                        </div>
                        <div className='grid gap-2'>
                            <Label>{t('Max Invitees Per Inviter')} (0 = {t('unlimited')})</Label>
                            <Input
                                type='number'
                                value={editing.max_invitees_per_inviter || 0}
                                onChange={(e) =>
                                    updateField(
                                        'max_invitees_per_inviter',
                                        Number(e.target.value),
                                    )
                                }
                            />
                        </div>
                        <div className='grid gap-2'>
                            <Label>{t('Effective Days')} (0 = {t('no limit')})</Label>
                            <Input
                                type='number'
                                value={editing.effective_days || 0}
                                onChange={(e) =>
                                    updateField('effective_days', Number(e.target.value))
                                }
                            />
                            <p className='text-muted-foreground text-xs'>
                                {t('Inviter receives reward only within N days of invitee registration. 0 means no time limit.')}
                            </p>
                        </div>
                        <div className='grid gap-2'>
                            <Label>{t('Priority')}</Label>
                            <Input
                                type='number'
                                value={editing.priority || 0}
                                onChange={(e) =>
                                    updateField('priority', Number(e.target.value))
                                }
                            />
                        </div>
                        <div className='flex items-center gap-2'>
                            <Switch
                                checked={editing.enabled ?? true}
                                onCheckedChange={(v) => updateField('enabled', v)}
                            />
                            <Label>{t('Enabled')}</Label>
                        </div>
                    </div>
                    <DialogFooter>
                        <Button variant='outline' onClick={() => setDialogOpen(false)}>
                            {t('Cancel')}
                        </Button>
                        <Button onClick={handleSave} disabled={saving}>
                            {saving ? t('Saving...') : t('Save')}
                        </Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>
        </div>
    )
}
