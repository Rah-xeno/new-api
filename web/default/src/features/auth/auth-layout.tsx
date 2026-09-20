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
import { Link } from '@tanstack/react-router'
import { ArrowLeft } from 'lucide-react'
import { useTranslation } from 'react-i18next'

import { Skeleton } from '@/components/ui/skeleton'
import MotionScene from '@/features/home/motion-scene'
import { useSystemConfig } from '@/hooks/use-system-config'
import { cn } from '@/lib/utils'

type AuthLayoutProps = {
  children: React.ReactNode
  variant?: 'default' | 'scaling'
}

export function AuthLayout({ children, variant = 'default' }: AuthLayoutProps) {
  const { t } = useTranslation()
  const { systemName, logo, loading } = useSystemConfig()
  const isScaling = variant === 'scaling'

  return (
    <div
      className={cn(
        'relative grid h-svh max-w-none',
        isScaling && 'scaling-auth-layout',
        isScaling &&
          typeof document.startViewTransition !== 'function' &&
          'scaling-auth-fallback'
      )}
    >
      {isScaling && (
        <aside className='scaling-auth-visual' aria-hidden='true'>
          <MotionScene />
          <div className='scaling-auth-visual-foot'>
            <span>NEWAPI / SCALING POSSIBILITIES</span>
            <span>01—04</span>
          </div>
        </aside>
      )}
      <Link
        to='/'
        viewTransition={isScaling}
        className={cn(
          'absolute top-4 left-4 z-10 flex items-center gap-2 transition-opacity hover:opacity-80 sm:top-8 sm:left-8',
          isScaling && 'scaling-auth-brand'
        )}
      >
        <div className='relative h-8 w-8'>
          {loading ? (
            <Skeleton className='absolute inset-0 rounded-full' />
          ) : (
            <img
              src={logo}
              alt={t('Logo')}
              className='h-8 w-8 rounded-full object-cover'
            />
          )}
        </div>
        {loading ? (
          <Skeleton className='h-6 w-24' />
        ) : (
          <h1 className='text-xl font-medium'>{systemName}</h1>
        )}
      </Link>
      {isScaling && (
        <Link to='/' viewTransition className='scaling-auth-back'>
          <ArrowLeft size={14} aria-hidden='true' />
          {t('Back to Home')}
        </Link>
      )}
      <div
        className={cn(
          'container flex items-center pt-16 sm:pt-0',
          isScaling && 'scaling-auth-content'
        )}
      >
        <div
          className={cn(
            'mx-auto flex w-full flex-col justify-center space-y-2 px-4 py-8 sm:w-[480px] sm:p-8',
            isScaling && 'scaling-auth-form-shell'
          )}
        >
          {children}
        </div>
      </div>
    </div>
  )
}
