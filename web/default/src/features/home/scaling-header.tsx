import { Link } from '@tanstack/react-router'
import { ArrowUpRight, Menu, X } from 'lucide-react'
import { useRef, useState } from 'react'
import { useTranslation } from 'react-i18next'

import { Button } from '@/components/ui/button'
import { useAuthStore } from '@/stores/auth-store'

export function ScalingBrand() {
  const { t } = useTranslation()
  return (
    <a href='#overview' className='brand-link' aria-label={t('NewAPI home')}>
      <span className='brand-symbol' aria-hidden='true'>
        <i />
        <i />
        <i />
      </span>
      <span>
        new<span className='brand-suffix'>api</span>
      </span>
    </a>
  )
}

export function ScalingHeader() {
  const { t } = useTranslation()
  const user = useAuthStore((s) => s.auth.user)
  const [open, setOpen] = useState(false)
  const toggleRef = useRef<HTMLButtonElement>(null)
  const navigation = [{ label: t('Model plaza'), to: '/pricing' as const }]
  return (
    <header
      className='site-header'
      onKeyDown={(event) => {
        if (event.key === 'Escape' && open) {
          setOpen(false)
          toggleRef.current?.focus()
        }
      }}
    >
      <ScalingBrand />
      <nav className='desktop-nav' aria-label={t('Main navigation')}>
        {navigation.map((item) => (
          <Link key={item.to} to={item.to}>
            {item.label}
          </Link>
        ))}
      </nav>
      <Link
        to={user ? '/dashboard' : '/sign-in'}
        viewTransition={!user}
        className='header-link'
      >
        {t(user ? 'Console' : 'Sign In')} <ArrowUpRight />
      </Link>
      <Button
        ref={toggleRef}
        variant='ghost'
        className='menu-toggle'
        type='button'
        aria-label={t(open ? 'Close navigation' : 'Open navigation')}
        aria-expanded={open}
        aria-controls='mobile-navigation'
        onClick={() => setOpen(!open)}
      >
        {open ? <X /> : <Menu />}
      </Button>
      {open && (
        <nav
          id='mobile-navigation'
          className='mobile-navigation'
          aria-label={t('Mobile navigation')}
        >
          {navigation.map((item, index) => (
            <Link key={item.to} to={item.to} onClick={() => setOpen(false)}>
              <span>0{index + 1}</span>
              {item.label}
              <ArrowUpRight />
            </Link>
          ))}
        </nav>
      )}
    </header>
  )
}
