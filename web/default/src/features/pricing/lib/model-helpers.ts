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
import { EXCLUDED_GROUPS, FILTER_ALL, QUOTA_TYPE_VALUES } from '../constants'
import type { PricingModel, UsableGroup } from '../types'

// ----------------------------------------------------------------------------
// Model Helper Utilities
// ----------------------------------------------------------------------------

type UsableGroupInput =
  | UsableGroup
  | Record<string, { desc?: string; ratio?: number }>

function usableGroupNames(usableGroup?: UsableGroupInput): string[] {
  return Object.keys(usableGroup ?? {}).filter(
    (group) => !EXCLUDED_GROUPS.includes(group)
  )
}

/** Normalize backend model membership to the concrete groups visible to users. */
export function normalizeModelGroups(
  model: PricingModel,
  usableGroup?: UsableGroupInput
): string[] {
  const enabled = Array.isArray(model.enable_groups) ? model.enable_groups : []
  const visible = usableGroupNames(usableGroup)
  const candidates = enabled.includes('all') ? visible : enabled
  const visibleSet = new Set(visible)
  const seen = new Set<string>()

  return candidates.filter((group) => {
    if (EXCLUDED_GROUPS.includes(group) || seen.has(group)) return false
    if (usableGroup && !visibleSet.has(group)) return false
    seen.add(group)
    return true
  })
}

/** Groups a model can show in the catalog, already intersected with access. */
export function getModelDisplayGroups(
  model: PricingModel,
  usableGroup?: UsableGroupInput
): string[] {
  return normalizeModelGroups(model, usableGroup)
}

/** Backwards-compatible name used by the model details drawer. */
export function getAvailableGroups(
  model: PricingModel,
  usableGroup: UsableGroupInput
): string[] {
  return getModelDisplayGroups(model, usableGroup)
}

/**
 * Read a configured group ratio while preserving valid zero ratios.
 */
export function getConfiguredGroupRatio(
  groupRatio: Record<string, number>,
  group: string
): number {
  const ratio = groupRatio[group]
  return typeof ratio === 'number' && Number.isFinite(ratio) && ratio >= 0
    ? ratio
    : 1
}

export type ModelPriceGroup = { group: string | null; ratio: number }

/** Resolve the concrete group and multiplier used for a model's displayed price. */
export function getModelPriceGroup(
  model: PricingModel,
  selectedGroup?: string,
  usableGroup?: UsableGroupInput
): ModelPriceGroup {
  const groups = getModelDisplayGroups(model, usableGroup)
  const ratios = model.group_ratio ?? {}

  if (
    selectedGroup &&
    selectedGroup !== FILTER_ALL &&
    groups.includes(selectedGroup)
  ) {
    return {
      group: selectedGroup,
      ratio: getConfiguredGroupRatio(ratios, selectedGroup),
    }
  }

  let best: ModelPriceGroup = { group: null, ratio: 1 }
  let bestFixedPrice = Number.POSITIVE_INFINITY
  for (const group of groups) {
    const ratio = getConfiguredGroupRatio(ratios, group)
    if (model.quota_type === QUOTA_TYPE_VALUES.REQUEST) {
      const override = model.fixed_price_overrides?.[group]?.price
      const fixedPrice =
        typeof override === 'number' && Number.isFinite(override)
          ? override
          : (model.model_price || 0) * (model.uniform_group_price ? 1 : ratio)
      if (best.group === null || fixedPrice < bestFixedPrice) {
        best = { group, ratio }
        bestFixedPrice = fixedPrice
      }
      continue
    }
    if (best.group === null || ratio < best.ratio) {
      best = { group, ratio }
    }
  }
  return best
}

/**
 * Resolve the group ratio used by model square summary prices.
 *
 * When no specific group is selected, the model square shows the best price
 * available to the viewer. When a group filter is active, it shows that
 * group's price instead.
 */
export function getDisplayGroupRatio(
  model: PricingModel,
  selectedGroup?: string
): number {
  return getModelPriceGroup(model, selectedGroup).ratio
}

/**
 * Replace model placeholder in endpoint path
 */
export function replaceModelInPath(path: string, modelName: string): string {
  return path.replaceAll('{model}', modelName)
}

/**
 * Check if model is token-based pricing
 */
export function isTokenBasedModel(model: PricingModel): boolean {
  return model.quota_type === QUOTA_TYPE_VALUES.TOKEN
}
