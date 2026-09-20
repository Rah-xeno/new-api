import {
  createMemoryHistory,
  createRootRoute,
  createRoute,
  createRouter,
  RouterProvider,
} from '@tanstack/react-router'
import { fireEvent, render, screen } from '@testing-library/react'
import i18n from 'i18next'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import { AuthLayout } from '../auth-layout'

vi.mock('@/hooks/use-system-config', () => ({
  useSystemConfig: () => ({
    systemName: 'NewAPI',
    logo: '/logo.png',
    loading: false,
  }),
}))

async function renderLayout(variant?: 'default' | 'scaling') {
  const root = createRootRoute()
  const home = createRoute({
    getParentRoute: () => root,
    path: '/',
    component: () => <h1>Homepage</h1>,
  })
  const signIn = createRoute({
    getParentRoute: () => root,
    path: '/sign-in',
    component: () => (
      <AuthLayout variant={variant}>
        <p>Form content</p>
      </AuthLayout>
    ),
  })
  const router = createRouter({
    routeTree: root.addChildren([home, signIn]),
    history: createMemoryHistory({ initialEntries: ['/sign-in'] }),
  })
  render(<RouterProvider router={router} />)
  await screen.findByText('Form content')
}

beforeEach(async () => {
  await i18n.changeLanguage('en')
  vi.spyOn(window, 'scrollTo').mockImplementation(() => {})
  vi.spyOn(HTMLMediaElement.prototype, 'play').mockResolvedValue()
  vi.spyOn(HTMLMediaElement.prototype, 'pause').mockImplementation(() => {})
})

afterEach(() => {
  Reflect.deleteProperty(document, 'startViewTransition')
  vi.restoreAllMocks()
})

describe('AuthLayout', () => {
  it('shows the homepage sphere and preview image alongside the sign-in content', async () => {
    await renderLayout('scaling')

    const video = screen.getByLabelText('White 3D sphere animation')
    expect(video).toHaveAttribute(
      'src',
      '/scaling-home/media/scaling-sphere.avc.mp4'
    )
    expect(video).toHaveAttribute(
      'poster',
      '/scaling-home/media/scaling-poster.jpg'
    )
    expect(screen.getByText('Form content')).toBeVisible()
  })

  it.each(['NewAPI', 'Back to Home'])(
    'returns home through %s without a document reload and requests a transition',
    async (label) => {
      await renderLayout('scaling')
      const startViewTransition = vi.fn((update: () => Promise<void>) => {
        void update()
      })
      Object.defineProperty(document, 'startViewTransition', {
        configurable: true,
        value: startViewTransition,
      })

      fireEvent.click(screen.getByRole('link', { name: new RegExp(label) }))

      expect(
        await screen.findByRole('heading', { name: 'Homepage' })
      ).toBeVisible()
      expect(startViewTransition).toHaveBeenCalledOnce()
    }
  )

  it('keeps homepage navigation working when the transition API is unavailable', async () => {
    await renderLayout('scaling')

    expect(document.querySelector('.scaling-auth-fallback')).toBeInTheDocument()

    fireEvent.click(screen.getByRole('link', { name: 'Back to Home' }))

    expect(
      await screen.findByRole('heading', { name: 'Homepage' })
    ).toBeVisible()
  })

  it('keeps the default auth layout free of the scaling visual', async () => {
    await renderLayout()

    expect(document.querySelector('.scaling-auth-layout')).toBeNull()
    expect(document.querySelector('video')).toBeNull()
    expect(screen.getByText('Form content')).toBeVisible()
  })
})
