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
const RELOAD_KEY = 'newapi_asset_reload_at'
const RELOAD_COOLDOWN_MS = 30_000
const CHUNK_ERROR_PATTERNS = [
  'chunkloaderror',
  'loading chunk',
  'loading css chunk',
  'failed to fetch dynamically imported module',
  'error loading dynamically imported module',
]

function reloadWithFreshDocument(): void {
  const now = Date.now()
  try {
    const lastReload = Number(window.sessionStorage.getItem(RELOAD_KEY) || 0)
    if (lastReload > 0 && now - lastReload <= RELOAD_COOLDOWN_MS) return
    window.sessionStorage.setItem(RELOAD_KEY, String(now))
  } catch {
    // Continue when browser storage is disabled.
  }

  const url = new URL(window.location.href)
  url.searchParams.set('__newapi_asset_reload', String(now))
  window.location.replace(url.toString())
}

function isChunkError(value: unknown): boolean {
  const message = String(
    value instanceof Error ? value.message : (value ?? '')
  ).toLowerCase()
  return CHUNK_ERROR_PATTERNS.some((pattern) => message.includes(pattern))
}

export function installAssetLoadRecovery(): void {
  window.addEventListener(
    'error',
    (event) => {
      const target = event.target
      if (
        target instanceof HTMLScriptElement ||
        target instanceof HTMLLinkElement
      ) {
        const assetUrl =
          target instanceof HTMLScriptElement ? target.src : target.href
        try {
          const url = new URL(assetUrl, window.location.href)
          if (
            url.origin === window.location.origin &&
            (url.pathname.startsWith('/static/') ||
              url.pathname.startsWith('/assets/'))
          ) {
            reloadWithFreshDocument()
            return
          }
        } catch {
          // Fall through to checking the error message.
        }
      }
      if (isChunkError(event.error || event.message)) reloadWithFreshDocument()
    },
    true
  )

  window.addEventListener('unhandledrejection', (event) => {
    if (isChunkError(event.reason)) reloadWithFreshDocument()
  })
}
