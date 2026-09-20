import {
  createMemoryHistory,
  createRootRoute,
  createRoute,
  createRouter,
  RouterProvider,
} from '@tanstack/react-router'
import {
  act,
  fireEvent,
  render,
  screen,
  waitFor,
  within,
} from '@testing-library/react'
import i18n from 'i18next'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import en from '@/i18n/locales/en.json'
import zh from '@/i18n/locales/zh.json'
import { useAuthStore } from '@/stores/auth-store'

import { ScalingHome } from '../scaling-home'

async function showHome() {
  const root = createRootRoute()
  const home = createRoute({
    getParentRoute: () => root,
    path: '/',
    component: ScalingHome,
  })
  const signin = createRoute({
    getParentRoute: () => root,
    path: '/sign-in',
    component: () => <h1>Native sign in route</h1>,
  })
  const dashboard = createRoute({
    getParentRoute: () => root,
    path: '/dashboard',
    component: () => <h1>Native dashboard route</h1>,
  })
  const pricing = createRoute({
    getParentRoute: () => root,
    path: '/pricing',
    component: () => <h1>Native model square route</h1>,
  })
  const router = createRouter({
    routeTree: root.addChildren([home, signin, dashboard, pricing]),
    history: createMemoryHistory({ initialEntries: ['/'] }),
  })
  render(<RouterProvider router={router} />)
  await screen.findByRole('heading', { level: 1 })
}

beforeEach(async () => {
  vi.spyOn(window, 'scrollTo').mockImplementation(() => {})
  vi.spyOn(HTMLMediaElement.prototype, 'play').mockResolvedValue()
  vi.spyOn(HTMLMediaElement.prototype, 'pause').mockImplementation(() => {})
  useAuthStore.setState((s) => ({ auth: { ...s.auth, user: null } }))
  i18n.addResourceBundle('zhCN', 'translation', zh.translation, true, true)
  i18n.addResourceBundle('en', 'translation', en.translation, true, true)
  await i18n.changeLanguage('zhCN')
})
afterEach(() => {
  Reflect.deleteProperty(document, 'startViewTransition')
  vi.restoreAllMocks()
  useAuthStore.setState((s) => ({ auth: { ...s.auth, user: null } }))
})

describe('Scaling homepage integration', () => {
  it('keeps the three Chinese title lines without punctuation and omits removed sections', async () => {
    await showHome()
    expect(screen.getByRole('heading', { level: 1 }).textContent).toBe(
      '一个接口连接多种AI 模型'
    )
    expect(document.querySelector('#quickstart')).toBeNull()
    expect(
      screen.queryByRole('button', { name: /暂停动画/ })
    ).not.toBeInTheDocument()
    expect(screen.queryByText(/尚未连接业务后端/)).not.toBeInTheDocument()
  })
  it('routes signed-out visitors to native sign-in instead of the landing-page fallback', async () => {
    await showHome()
    const login = screen.getByRole('link', { name: '登录' })
    expect(login).toHaveAttribute('href', '/sign-in')
    fireEvent.click(login)
    expect(
      await screen.findByRole('heading', { name: 'Native sign in route' })
    ).toBeVisible()
  })
  it.each(['登录', '开始构建'])(
    'navigates %s within the router and starts the shared-element transition when supported',
    async (label) => {
      await showHome()
      const startViewTransition = vi.fn((update: () => Promise<void>) => {
        void update()
      })
      Object.defineProperty(document, 'startViewTransition', {
        configurable: true,
        value: startViewTransition,
      })

      fireEvent.click(screen.getByRole('link', { name: label }))

      expect(
        await screen.findByRole('heading', { name: 'Native sign in route' })
      ).toBeVisible()
      expect(startViewTransition).toHaveBeenCalledOnce()
    }
  )
  it('sends an authenticated user to the native console', async () => {
    useAuthStore.setState((s) => ({
      auth: { ...s.auth, user: { id: 1, username: 'tester', role: 1 } },
    }))
    await showHome()
    expect(screen.getByRole('link', { name: '开始构建' })).toHaveAttribute(
      'href',
      '/dashboard'
    )
  })
  it('shows all model families in the model plaza without a search control', async () => {
    await showHome()
    expect(screen.queryByRole('textbox')).not.toBeInTheDocument()
    expect(screen.getAllByRole('article')).toHaveLength(4)
    expect(screen.getByText('模型家族')).toBeVisible()
    expect(screen.getByRole('link', { name: '模型广场' })).toBeVisible()
  })
  it('opens the native model square from desktop navigation', async () => {
    await showHome()
    const modelPlaza = screen.getByRole('link', { name: '模型广场' })
    expect(modelPlaza).toHaveAttribute('href', '/pricing')
    expect(
      screen.queryByRole('link', { name: '平台概览' })
    ).not.toBeInTheDocument()
    fireEvent.click(modelPlaza)
    expect(
      await screen.findByRole('heading', { name: 'Native model square route' })
    ).toBeVisible()
  })
  it('opens the native model square from mobile navigation without the removed overview link', async () => {
    await showHome()
    fireEvent.click(screen.getByRole('button', { name: '打开导航' }))
    const nav = within(screen.getByRole('navigation', { name: '移动导航' }))
    expect(
      nav.queryByRole('link', { name: '平台概览' })
    ).not.toBeInTheDocument()
    const modelPlaza = nav.getByRole('link', { name: /模型广场/ })
    expect(modelPlaza).toHaveAttribute('href', '/pricing')
    fireEvent.click(modelPlaza)
    expect(
      await screen.findByRole('heading', { name: 'Native model square route' })
    ).toBeVisible()
  })
  it('closes mobile navigation with Escape and restores focus to the toggle', async () => {
    await showHome()
    const toggle = screen.getByRole('button', { name: '打开导航' })
    fireEvent.click(toggle)
    expect(toggle).toHaveAttribute('aria-expanded', 'true')
    const nav = screen.getByRole('navigation', { name: '移动导航' })
    fireEvent.keyDown(nav, { key: 'Escape' })
    expect(toggle).toHaveAttribute('aria-expanded', 'false')
    expect(toggle).toHaveFocus()
    expect(
      screen.queryByRole('navigation', { name: '移动导航' })
    ).not.toBeInTheDocument()
  })
  it('updates the homepage copy when the native language changes', async () => {
    await showHome()
    await act(() => i18n.changeLanguage('en'))
    await waitFor(() =>
      expect(screen.getByRole('heading', { level: 1 }).textContent).toBe(
        'One interfaceConnect multipleAI models'
      )
    )
    expect(screen.getByRole('link', { name: 'Model plaza' })).toBeVisible()
  })
})
