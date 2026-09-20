import { Link } from '@tanstack/react-router'
import { ArrowRight, ArrowUpRight, ChevronDown } from 'lucide-react'
import { useTranslation } from 'react-i18next'

import { useAuthStore } from '@/stores/auth-store'

import MotionScene from './motion-scene'
import { ScalingBrand, ScalingHeader } from './scaling-header'
import { ScalingModels } from './scaling-models'

import '@/styles/scaling-home.css'

export function ScalingHome() {
  const { t, i18n } = useTranslation()
  const user = useAuthStore((s) => s.auth.user)
  return (
    <div
      className='scaling-home'
      lang={i18n.language.startsWith('zh') ? 'zh' : i18n.language}
    >
      <a className='skip-link' href='#main-content'>
        {t('Skip to main content')}
      </a>
      <main id='main-content'>
        <section
          id='overview'
          className='landing-hero'
          aria-labelledby='hero-title'
        >
          <MotionScene />
          <ScalingHeader />
          <div className='hero-copy'>
            <p className='eyeline hero-eyeline'>
              <span />
              ONE API. MORE POSSIBILITIES.
            </p>
            <h1 id='hero-title'>
              <span>{t('One interface')}</span>
              <span>{t('Connect multiple')}</span>
              <span className='hero-accent'>{t('AI models')}</span>
            </h1>
            <p className='hero-description'>
              GPT、Claude、Grok、Gemini。
              <br />
              {t('One connection. Turn great ideas into great products.')}
            </p>
            <div className='hero-actions flex flex-wrap items-center gap-7'>
              <Link
                className='accent-button'
                to={user ? '/dashboard' : '/sign-in'}
                viewTransition={!user}
              >
                {t('Start building')} <ArrowUpRight />
              </Link>
              <a className='text-link' href='#model-plaza'>
                {t('Explore models')} <ArrowRight />
              </a>
            </div>
            <p className='hero-footnote'>
              {t('One key · One interface · Switch freely')}
            </p>
          </div>
          <a className='hero-scroll' href='#model-plaza'>
            <span>{t('Scroll to explore')}</span>
            <ChevronDown />
          </a>
        </section>
        <div
          className='ecosystem-strip section-shell'
          aria-label={t('Model plaza')}
        >
          <p>
            {t('Connect the models you know')}
            <br />
            <strong>{t('Model plaza')}</strong>
          </p>
          <div className='ecosystem-wordmarks'>
            <span>OpenAI</span>
            <span className='claude-wordmark'>Claude</span>
            <span className='grok-wordmark'>grok</span>
            <span>
              Gemini
              <span className='gemini-plus' aria-hidden='true'>
                +
              </span>
            </span>
          </div>
          <span className='ecosystem-note'>
            MULTIPLE MODELS.
            <br />
            ONE CONNECTION.
          </span>
        </div>
        <ScalingModels />
      </main>
      <footer className='site-footer section-shell'>
        <ScalingBrand />
        <p>{t('Leave complexity to the interface. Keep creating.')}</p>
        <a href='#overview'>
          {t('Back to top')} <ArrowUpRight />
        </a>
        <div className='footer-meta'>
          <span>NEWAPI / SCALING POSSIBILITIES</span>
          <a
            href='https://github.com/QuantumNous/new-api'
            target='_blank'
            rel='noopener noreferrer'
          >
            New API / QuantumNous
          </a>
          <a
            href='https://github.com/ljyoukong-cpu/goodapinew'
            target='_blank'
            rel='noopener noreferrer'
          >
            NewAPI · AGPL-3.0
          </a>
        </div>
      </footer>
    </div>
  )
}
