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
import assert from 'node:assert/strict'
import { describe, test } from 'node:test'

import { formatAnalyticsSeconds } from './lib'

describe('admin analytics formatters', () => {
  test('normalizes the simplified Chinese interface language for Intl', () => {
    assert.equal(
      formatAnalyticsSeconds(1.5, 'zhCN'),
      new Intl.NumberFormat('zh-CN', {
        style: 'unit',
        unit: 'second',
        unitDisplay: 'short',
        maximumFractionDigits: 2,
      }).format(1.5)
    )
  })

  test('normalizes the traditional Chinese interface language for Intl', () => {
    assert.doesNotThrow(() => formatAnalyticsSeconds(12, 'zhTW'))
  })
})
