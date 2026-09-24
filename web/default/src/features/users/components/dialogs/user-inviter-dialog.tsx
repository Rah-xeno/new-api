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
import { useMutation } from '@tanstack/react-query'
import { useId } from 'react'
import { useForm } from 'react-hook-form'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'

import { Dialog } from '@/components/dialog'
import { Button } from '@/components/ui/button'
import {
  Field,
  FieldDescription,
  FieldError,
  FieldGroup,
  FieldLabel,
} from '@/components/ui/field'
import { Input } from '@/components/ui/input'

import { bindUserInviter } from '../../api'
import {
  userInviterFormSchema,
  type UserInviterFormValues,
} from '../../lib/user-form'
import type { User } from '../../types'

interface UserInviterDialogProps {
  user: User
  onOpenChange: (open: boolean) => void
  onSuccess: () => void
}

export function UserInviterDialog(props: UserInviterDialogProps) {
  const { t } = useTranslation()
  const formId = useId()
  const inputId = useId()
  const form = useForm<UserInviterFormValues>({
    resolver: zodResolver(userInviterFormSchema),
    defaultValues: { aff_code: '' },
  })
  const mutation = useMutation({
    mutationFn: (values: UserInviterFormValues) =>
      bindUserInviter(props.user.id, values.aff_code),
    onSuccess: (result) => {
      if (!result.success) {
        form.setError('aff_code', {
          type: 'server',
          message: result.message || t('Failed to bind inviter'),
        })
        return
      }
      toast.success(t('Inviter bound successfully'))
      props.onOpenChange(false)
      props.onSuccess()
    },
  })
  const error = form.formState.errors.aff_code
  const alreadyBound = (props.user.inviter_id ?? 0) !== 0

  return (
    <Dialog
      open
      onOpenChange={(open) => {
        if (!mutation.isPending) props.onOpenChange(open)
      }}
      title={t('Bind inviter')}
      description={t('Bind an inviter to {{username}} (ID: {{id}}).', {
        username: props.user.username,
        id: props.user.id,
      })}
      contentClassName='sm:max-w-lg'
      showCloseButton={!mutation.isPending}
      footer={
        <>
          <Button
            variant='outline'
            disabled={mutation.isPending}
            onClick={() => props.onOpenChange(false)}
          >
            {t('Cancel')}
          </Button>
          <Button
            type='submit'
            form={formId}
            disabled={mutation.isPending || alreadyBound}
          >
            {mutation.isPending ? t('Processing...') : t('Confirm binding')}
          </Button>
        </>
      }
    >
      <form
        id={formId}
        onSubmit={form.handleSubmit((values) => {
          if (!mutation.isPending && !alreadyBound) mutation.mutate(values)
        })}
      >
        <FieldGroup>
          <Field data-invalid={!!error}>
            <FieldLabel htmlFor={inputId}>{t('Invitation code')}</FieldLabel>
            <Input
              {...form.register('aff_code')}
              id={inputId}
              maxLength={32}
              autoComplete='off'
              disabled={mutation.isPending || alreadyBound}
              aria-invalid={!!error}
              aria-describedby={`${inputId}-description${error ? ` ${inputId}-error` : ''}`}
            />
            <FieldDescription id={`${inputId}-description`}>
              {t(
                'Enter the inviter’s code. Once bound, the inviter cannot be changed. Registration rewards and past payment rewards will not be issued retroactively.'
              )}
            </FieldDescription>
            {error && (
              <FieldError id={`${inputId}-error`}>
                {error.type === 'server'
                  ? error.message
                  : t(error.message || '')}
              </FieldError>
            )}
          </Field>
        </FieldGroup>
      </form>
    </Dialog>
  )
}
