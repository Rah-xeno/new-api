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
import type { SVGProps } from 'react'

type IconSuanliProps = SVGProps<SVGSVGElement> & {
  size?: number
}

/** Crystal-core currency mark (variant A), shared by all quota/balance displays. */
export function IconSuanli({ size = 20, ...props }: IconSuanliProps) {
  return (
    <svg
      xmlns='http://www.w3.org/2000/svg'
      viewBox='0 0 24 24'
      width={size}
      height={size}
      fill='none'
      {...props}
    >
      <circle
        cx='12'
        cy='12'
        r='9.25'
        fill='none'
        stroke='currentColor'
        strokeWidth='2.2'
      />
      <rect
        x='8.4'
        y='8.4'
        width='7.2'
        height='7.2'
        rx='1.2'
        transform='rotate(45 12 12)'
        fill='currentColor'
      />
    </svg>
  )
}
