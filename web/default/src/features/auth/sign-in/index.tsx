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
import { Link, useSearch } from '@tanstack/react-router'
import { useTranslation } from 'react-i18next'

import { useStatus } from '@/hooks/use-status'

import { AuthLayout } from '../auth-layout'
import { TermsFooter } from '../components/terms-footer'
import { UserAuthForm } from './components/user-auth-form'

import '@/styles/scaling-auth.css'

export function SignIn() {
  const { t } = useTranslation()
  const { redirect } = useSearch({ from: '/(auth)/sign-in' })
  const { status } = useStatus()

  return (
    <AuthLayout variant='scaling'>
      <div className='scaling-sign-in w-full'>
        <div className='scaling-auth-heading'>
          <p className='scaling-auth-kicker'>
            <span aria-hidden='true'>01 /</span>
            {t('Secure access')}
          </p>
          <h2 className='scaling-auth-title'>{t('Sign in')}</h2>
          {!status?.self_use_mode_enabled &&
            status?.register_enabled !== false && (
              <p className='scaling-auth-subtitle'>
                {t("Don't have an account?")}{' '}
                <Link
                  to='/sign-up'
                  viewTransition
                  className='scaling-auth-inline-link'
                >
                  {t('Sign up')}
                </Link>
                .
              </p>
            )}
        </div>

        <UserAuthForm
          className='scaling-auth-form-fields'
          redirectTo={redirect}
        />

        <TermsFooter
          variant='sign-in'
          status={status}
          className='scaling-auth-terms'
        />
      </div>
    </AuthLayout>
  )
}
