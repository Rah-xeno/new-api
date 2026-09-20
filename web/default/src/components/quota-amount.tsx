import { Fragment } from 'react'

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
import { IconSuanli } from '@/assets/custom/icon-suanli'
import {
  formatCurrencyFromUSD,
  formatLocalCurrencyAmount,
  formatQuotaWithCurrency,
  type CurrencyFormatOptions,
} from '@/lib/currency'
import { cn } from '@/lib/utils'
import {
  DEFAULT_CURRENCY_CONFIG,
  useSystemConfigStore,
} from '@/stores/system-config-store'

/** Currency unit for table headings and quota details. */
export function QuotaUnit(props: { fallback: string }) {
  const currency = useSystemConfigStore((state) => state.config.currency)
  if (currency.quotaDisplayType !== 'CUSTOM') return props.fallback
  return (
    <IconSuanli
      className='text-success inline-block size-[1em] align-[-0.125em]'
      role='img'
      aria-label={currency.customCurrencySymbol}
    />
  )
}

/** Preserve formatted pricing precision, ranges and units; replace only the credit mark. */
export function QuotaText(props: { children: string }) {
  const currency = useSystemConfigStore((state) => state.config.currency)
  if (currency.quotaDisplayType !== 'CUSTOM') return props.children

  const symbol =
    currency.customCurrencySymbol.trim() ||
    DEFAULT_CURRENCY_CONFIG.customCurrencySymbol
  const escapedSymbol = symbol.replaceAll(/[.*+?^${}()|[\]\\]/g, '\\$&')
  const matches = [
    ...props.children.matchAll(
      new RegExp(`${escapedSymbol}\\s*(?=[+-]?\\d)`, 'g')
    ),
  ]
  if (matches.length === 0) return props.children

  return (
    <span aria-label={props.children}>
      {props.children.slice(0, matches[0].index)}
      {matches.map((match, index) => (
        <Fragment key={match.index}>
          <IconSuanli
            className='text-success mr-1.5 inline-block size-[1em] align-[-0.125em]'
            aria-hidden='true'
          />
          {props.children.slice(
            match.index + match[0].length,
            matches[index + 1]?.index
          )}
        </Fragment>
      ))}
    </span>
  )
}

interface QuotaAmountProps {
  value: number
  unit?: 'quota' | 'usd' | 'display'
  options?: CurrencyFormatOptions
  className?: string
  iconClassName?: string
}

/** Display site credit amounts without substituting the actual payment currency. */
export function QuotaAmount(props: QuotaAmountProps) {
  const currency = useSystemConfigStore((state) => state.config.currency)
  const iconCurrency = currency.quotaDisplayType === 'CUSTOM'
  let formatter = formatQuotaWithCurrency
  if (props.unit === 'usd') formatter = formatCurrencyFromUSD
  if (props.unit === 'display') formatter = formatLocalCurrencyAmount
  const value = formatter(props.value, {
    digitsLarge: 2,
    digitsSmall: 4,
    abbreviate: true,
    ...props.options,
    ...(iconCurrency ? { showSymbol: false } : {}),
  })

  if (!iconCurrency) return value

  return (
    <span
      className={cn(
        'inline-flex items-center gap-1.5 align-middle',
        props.className
      )}
      aria-label={`${currency.customCurrencySymbol} ${value}`}
    >
      <IconSuanli
        className={cn('text-success size-[1em] shrink-0', props.iconClassName)}
        aria-hidden='true'
      />
      <span>{value}</span>
    </span>
  )
}
