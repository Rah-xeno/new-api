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
import { useForm } from 'react-hook-form'
import { useTranslation } from 'react-i18next'
import { z } from 'zod'

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

import { SettingsForm } from '../components/settings-form-layout'
import { SettingsPageFormActions } from '../components/settings-page-context'
import { SettingsSection } from '../components/settings-section'
import { useUpdateOption } from '../hooks/use-update-option'

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
    </SettingsSection>
  )
}
