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
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'

export type AnalyticsMetric = {
  key: string
  label: string
  value: string
  hint?: string
}

export function SummaryMetricCards(props: {
  metrics: AnalyticsMetric[]
  loading: boolean
}) {
  return (
    <div className='grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-3'>
      {props.metrics.map((metric) => (
        <Card key={metric.key} size='sm'>
          <CardHeader>
            <CardDescription>{metric.label}</CardDescription>
          </CardHeader>
          <CardContent className='flex flex-col gap-2'>
            {props.loading ? (
              <Skeleton className='h-8 w-32' />
            ) : (
              <div className='text-2xl font-semibold tracking-tight'>
                {metric.value}
              </div>
            )}
            <div className='text-muted-foreground min-h-4 text-xs'>
              {metric.hint}
            </div>
          </CardContent>
        </Card>
      ))}
    </div>
  )
}

export function MetricSection(props: {
  title: string
  description: string
  metrics: AnalyticsMetric[]
  loading: boolean
}) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>{props.title}</CardTitle>
        <CardDescription>{props.description}</CardDescription>
      </CardHeader>
      <CardContent className='grid grid-cols-1 gap-3 sm:grid-cols-2'>
        {props.metrics.map((metric) => (
          <div
            key={metric.key}
            className='bg-muted/40 flex min-w-0 flex-col gap-2 rounded-lg border p-3'
          >
            <div className='text-muted-foreground text-xs'>{metric.label}</div>
            {props.loading ? (
              <Skeleton className='h-6 w-28' />
            ) : (
              <div className='truncate text-lg font-semibold'>
                {metric.value}
              </div>
            )}
            {metric.hint ? (
              <div className='text-muted-foreground text-xs'>{metric.hint}</div>
            ) : null}
          </div>
        ))}
      </CardContent>
    </Card>
  )
}
