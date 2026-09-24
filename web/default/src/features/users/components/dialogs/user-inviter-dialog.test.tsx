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
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it, vi } from 'vitest'

import { bindUserInviter } from '../../api'
import type { User } from '../../types'
import { UserInviterDialog } from './user-inviter-dialog'

vi.mock('../../api', () => ({ bindUserInviter: vi.fn() }))

const invitee: User = {
  id: 42,
  username: 'alice',
  display_name: 'Alice',
  role: 1,
  status: 1,
  quota: 100,
  used_quota: 0,
  request_count: 0,
  group: 'default',
  inviter_id: 0,
}

function renderDialog(user: User = invitee) {
  const onOpenChange = vi.fn()
  const onSuccess = vi.fn()
  const client = new QueryClient({
    defaultOptions: { mutations: { retry: false } },
  })
  render(
    <QueryClientProvider client={client}>
      <UserInviterDialog
        user={user}
        onOpenChange={onOpenChange}
        onSuccess={onSuccess}
      />
    </QueryClientProvider>
  )
  return { onOpenChange, onSuccess }
}

describe('UserInviterDialog', () => {
  it('requires a code and submits the trimmed code for the selected user', async () => {
    vi.mocked(bindUserInviter).mockResolvedValue({ success: true })
    const user = userEvent.setup()
    const callbacks = renderDialog()
    const confirm = screen.getByRole('button', { name: 'Confirm binding' })

    await user.click(confirm)
    expect(await screen.findByText('Invitation code is required')).toBeVisible()
    expect(bindUserInviter).not.toHaveBeenCalled()

    await user.type(screen.getByLabelText('Invitation code'), ' code ')
    await user.click(confirm)

    await waitFor(() => expect(callbacks.onSuccess).toHaveBeenCalledOnce())
    expect(bindUserInviter).toHaveBeenCalledExactlyOnceWith(42, 'code')
    expect(callbacks.onOpenChange).toHaveBeenCalledWith(false)
  })

  it('keeps the dialog open and shows a server rejection', async () => {
    vi.mocked(bindUserInviter).mockResolvedValue({
      success: false,
      message: 'This user already has an inviter and it cannot be changed',
    })
    const user = userEvent.setup()
    const callbacks = renderDialog()
    await user.type(screen.getByLabelText('Invitation code'), 'code')
    await user.click(screen.getByRole('button', { name: 'Confirm binding' }))

    expect(
      await screen.findByText(
        'This user already has an inviter and it cannot be changed'
      )
    ).toBeVisible()
    expect(callbacks.onSuccess).not.toHaveBeenCalled()
    expect(callbacks.onOpenChange).not.toHaveBeenCalled()
  })

  it('disables binding for a user who already has an inviter', () => {
    renderDialog({ ...invitee, inviter_id: 7 })
    expect(screen.getByLabelText('Invitation code')).toBeDisabled()
    expect(
      screen.getByRole('button', { name: 'Confirm binding' })
    ).toBeDisabled()
    expect(bindUserInviter).not.toHaveBeenCalled()
  })
})
