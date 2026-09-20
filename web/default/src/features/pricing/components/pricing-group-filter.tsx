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
import { useTranslation } from 'react-i18next'

import { Button } from '@/components/ui/button'

import { FILTER_ALL } from '../constants'

export interface PricingGroupFilterProps {
  groups: string[]
  groupRatios: Record<string, number>
  selectedGroup: string
  onGroupChange: (group: string) => void
}

export function PricingGroupFilter(props: PricingGroupFilterProps) {
  const { t } = useTranslation()
  if (props.groups.length === 0) return null

  return (
    <section className='pricing-group-panel' aria-label={t('Call groups')}>
      <div className='pricing-group-heading'>
        <h2>{t('Call groups')}</h2>
        <p>{t('Compare by group')}</p>
      </div>
      <div
        className='pricing-group-options'
        role='group'
        aria-label={t('Select a pricing group')}
      >
        {[FILTER_ALL, ...props.groups].map((group) => (
          <Button
            key={group}
            variant={props.selectedGroup === group ? 'default' : 'outline'}
            aria-pressed={props.selectedGroup === group}
            onClick={() => props.onGroupChange(group)}
            title={group === FILTER_ALL ? t('All Groups') : group}
            className='min-h-11 max-w-64 gap-2 px-4'
          >
            <span className='truncate'>
              {group === FILTER_ALL ? t('All Groups') : group}
            </span>
            {group !== FILTER_ALL && props.groupRatios[group] != null && (
              <span className='shrink-0 font-mono text-xs opacity-80'>
                ×{props.groupRatios[group]}
              </span>
            )}
          </Button>
        ))}
      </div>
      <p className='pricing-group-context' role='status'>
        {props.selectedGroup === FILTER_ALL
          ? t(
              'One card per model. Prices start at the lowest available group rate; compare all groups in details.'
            )
          : t(
              'Showing models and prices for {{group}}. This does not change your API key group.',
              { group: props.selectedGroup }
            )}
      </p>
    </section>
  )
}
