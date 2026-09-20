import { act, render, screen } from '@testing-library/react'
import { beforeEach, expect, it, vi } from 'vitest'

import { getHomePageContent } from '../api'
import { Home } from '../index'

vi.mock('../api', () => ({ getHomePageContent: vi.fn() }))
vi.mock('@/context/theme-provider', () => ({
  useTheme: () => ({ resolvedTheme: 'light' }),
}))
vi.mock('@/components/layout', () => ({
  PublicLayout: ({ children }: { children: React.ReactNode }) => (
    <div data-testid='legacy-public-layout'>{children}</div>
  ),
}))
vi.mock('@/components/rich-content', () => ({
  RichContent: ({ content }: { content: string }) => (
    <article>{content}</article>
  ),
}))
vi.mock('../scaling-home', () => ({
  ScalingHome: () => <main data-testid='scaling-home'>Scaling homepage</main>,
}))

beforeEach(() => localStorage.clear())

it('shows the scaling homepage, not the old public layout, while configuration loads', () => {
  vi.mocked(getHomePageContent).mockReturnValue(new Promise(() => {}))
  render(<Home />)

  expect(screen.queryByTestId('legacy-public-layout')).not.toBeInTheDocument()
  expect(screen.getByTestId('scaling-home')).toBeVisible()
})

it('retains the same homepage node when the default configuration finishes loading', async () => {
  let finish!: (value: Awaited<ReturnType<typeof getHomePageContent>>) => void
  vi.mocked(getHomePageContent).mockReturnValue(
    new Promise((resolve) => {
      finish = resolve
    })
  )
  render(<Home />)
  const initialHome = screen.getByTestId('scaling-home')

  await act(async () => {
    finish({ success: true, data: '' })
  })

  expect(screen.getByTestId('scaling-home')).toBe(initialHome)
  expect(screen.queryByTestId('legacy-public-layout')).not.toBeInTheDocument()
})

it('still honors administrator-configured homepage content after loading', async () => {
  vi.mocked(getHomePageContent).mockResolvedValue({
    success: true,
    data: '# Custom homepage',
  })
  render(<Home />)

  expect(await screen.findByRole('article')).toHaveTextContent(
    '# Custom homepage'
  )
  expect(screen.queryByTestId('scaling-home')).not.toBeInTheDocument()
})
