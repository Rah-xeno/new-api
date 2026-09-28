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
import { zodResolver } from '@hookform/resolvers/zod'
import { useQuery } from '@tanstack/react-query'
import { useMemo, useState } from 'react'
import { useForm } from 'react-hook-form'
import { useTranslation } from 'react-i18next'
import { z } from 'zod'

import { StaticDataTable } from '@/components/data-table/static/static-data-table'
import { Button } from '@/components/ui/button'
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { api } from '@/lib/api'

import { SettingsForm } from '../components/settings-form-layout'
import { SettingsPageFormActions } from '../components/settings-page-context'
import { SettingsSection } from '../components/settings-section'
import { useUpdateOption } from '../hooks/use-update-option'

type AgentOverviewItem = {
  agent_id: number
  username: string
  display_name: string
  email: string
  aff_code: string
  invitee_total: number
  paid_invitee_total: number
  total_topup_amount: number
  commission_total: number
  commission_balance: number
  commission_withdrawn: number
}

type PagedData<T> = {
  page: number
  page_size: number
  total: number
  items: T[]
}

const PAGE_SIZE = 10

function formatMoney(cents: number | undefined) {
  return `¥${((cents ?? 0) / 100).toFixed(2)}`
}

const schema = z.object({
  firstRate: z.number().min(0).max(100),
  repeatRate: z.number().min(0).max(100),
})

type Values = z.infer<typeof schema>

export function AgentDistributionSettingsSection(props: {
  defaultValues: Values
}) {
  const { t } = useTranslation()
  const updateOption = useUpdateOption()
  const form = useForm<Values>({
    resolver: zodResolver(schema),
    defaultValues: props.defaultValues,
  })

  async function onSubmit(values: Values) {
    const updates = [
      {
        key: 'agent_distribution_setting.default_first_topup_rate',
        value: String(values.firstRate),
      },
      {
        key: 'agent_distribution_setting.default_repeat_topup_rate',
        value: String(values.repeatRate),
      },
    ]
    for (const update of updates) {
      const response = await updateOption.mutateAsync(update)
      if (!response.success) return
    }
    form.reset(values)
  }

  return (
    <SettingsSection title={t('Agent Distribution')}>
      <Form {...form}>
        <SettingsForm onSubmit={form.handleSubmit(onSubmit)}>
          <SettingsPageFormActions
            onSave={form.handleSubmit(onSubmit)}
            isSaving={updateOption.isPending || form.formState.isSubmitting}
            isSaveDisabled={!form.formState.isDirty}
            saveLabel='Save agent distribution settings'
          />
          <div className='grid gap-6 sm:grid-cols-2'>
            <FormField
              control={form.control}
              name='firstRate'
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t('Default first top-up rate (%)')}</FormLabel>
                  <FormControl>
                    <Input
                      type='number'
                      min={0}
                      max={100}
                      step='0.01'
                      {...field}
                      onChange={(event) =>
                        field.onChange(event.target.valueAsNumber)
                      }
                    />
                  </FormControl>
                  <FormDescription>
                    {t(
                      'Commission rate for the invited customer’s first successful payment.'
                    )}
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name='repeatRate'
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t('Default repeat top-up rate (%)')}</FormLabel>
                  <FormControl>
                    <Input
                      type='number'
                      min={0}
                      max={100}
                      step='0.01'
                      {...field}
                      onChange={(event) =>
                        field.onChange(event.target.valueAsNumber)
                      }
                    />
                  </FormControl>
                  <FormDescription>
                    {t(
                      'Commission rate for later successful payments from the same customer.'
                    )}
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
          </div>
        </SettingsForm>
      </Form>

      <AgentOverviewList />
    </SettingsSection>
  )
}

function AgentOverviewList() {
  const { t } = useTranslation()
  const [page, setPage] = useState(1)

  const query = useQuery({
    queryKey: ['agent-distribution', 'admin', 'overview', page],
    queryFn: async () => {
      const response = await api.get<{ success: boolean; data: PagedData<AgentOverviewItem> }>(
        '/api/agent-distribution/admin/overview',
        { params: { p: page, page_size: PAGE_SIZE } }
      )
      return response.data.data
    },
  })

  const data = query.data
  const totalPages = useMemo(
    () => (data ? Math.max(1, Math.ceil(data.total / PAGE_SIZE)) : 1),
    [data]
  )

  const columns = useMemo(
    () => [
      {
        id: 'agent',
        header: t('Agent'),
        cell: (row: AgentOverviewItem) =>
          row.display_name || row.username || row.email || `#${row.agent_id}`,
      },
      {
        id: 'aff_code',
        header: t('Affiliate code'),
        cell: (row: AgentOverviewItem) => row.aff_code || '-',
      },
      {
        id: 'invitee_total',
        header: t('Invited customers'),
        cell: (row: AgentOverviewItem) => row.invitee_total,
      },
      {
        id: 'paid_invitee_total',
        header: t('Paid customers'),
        cell: (row: AgentOverviewItem) => row.paid_invitee_total,
      },
      {
        id: 'total_topup_amount',
        header: t('Total top-up'),
        cell: (row: AgentOverviewItem) => formatMoney(row.total_topup_amount),
      },
      {
        id: 'commission_total',
        header: t('Total commission'),
        cell: (row: AgentOverviewItem) => formatMoney(row.commission_total),
      },
      {
        id: 'commission_balance',
        header: t('Withdrawable balance'),
        cell: (row: AgentOverviewItem) => formatMoney(row.commission_balance),
      },
      {
        id: 'commission_withdrawn',
        header: t('Withdrawn'),
        cell: (row: AgentOverviewItem) => formatMoney(row.commission_withdrawn),
      },
    ],
    [t]
  )

  return (
    <div className='flex flex-col gap-3'>
      <div className='flex items-center justify-between'>
        <h3 className='text-sm font-medium'>{t('Agent overview')}</h3>
      </div>
      <StaticDataTable
        columns={columns}
        data={data?.items ?? []}
        getRowKey={(row) => row.agent_id}
        empty={!data || data.items.length === 0}
        emptyContent={
          <div className='py-8 text-center text-sm text-muted-foreground'>
            {query.isLoading ? t('Loading...') : t('No agents yet')}
          </div>
        }
      />
      {totalPages > 1 && (
        <div className='flex items-center justify-end gap-2'>
          <Button
            variant='outline'
            size='sm'
            disabled={page <= 1}
            onClick={() => setPage((p) => Math.max(1, p - 1))}
          >
            {t('Previous')}
          </Button>
          <span className='text-sm text-muted-foreground'>
            {page} / {totalPages}
          </span>
          <Button
            variant='outline'
            size='sm'
            disabled={page >= totalPages}
            onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
          >
            {t('Next')}
          </Button>
        </div>
      )}
    </div>
  )
}
