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
import { VChart } from '@visactor/react-vchart'
import type { ComponentProps, ReactNode } from 'react'
import { useTranslation } from 'react-i18next'

import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import {
  Empty,
  EmptyDescription,
  EmptyHeader,
  EmptyTitle,
} from '@/components/ui/empty'
import { Skeleton } from '@/components/ui/skeleton'
import { useChartTheme } from '@/lib/use-chart-theme'
import { VCHART_OPTION } from '@/lib/vchart'

type AnalyticsChartCardProps = {
  title: string
  description: string
  loading: boolean
  hasData: boolean
  spec: ComponentProps<typeof VChart>['spec']
  className?: string
}

export function AnalyticsChartCard(props: AnalyticsChartCardProps) {
  const { t } = useTranslation()
  const { resolvedTheme, themeReady } = useChartTheme()
  let content: ReactNode
  if (props.loading) {
    content = <Skeleton className='h-80 w-full' />
  } else if (!props.hasData) {
    content = (
      <Empty className='h-80'>
        <EmptyHeader>
          <EmptyTitle>{t('No analytics data')}</EmptyTitle>
          <EmptyDescription>
            {t('There is no data for the selected period.')}
          </EmptyDescription>
        </EmptyHeader>
      </Empty>
    )
  } else if (!themeReady) {
    content = <Skeleton className='h-80 w-full' />
  } else {
    content = (
      <div className='h-80 w-full'>
        <VChart
          spec={{
            ...(props.spec as Record<string, unknown>),
            theme: resolvedTheme === 'dark' ? 'dark' : 'light',
            background: 'transparent',
          }}
          option={VCHART_OPTION}
        />
      </div>
    )
  }

  return (
    <Card className={props.className}>
      <CardHeader>
        <CardTitle>{props.title}</CardTitle>
        <CardDescription>{props.description}</CardDescription>
      </CardHeader>
      <CardContent>{content}</CardContent>
    </Card>
  )
}
