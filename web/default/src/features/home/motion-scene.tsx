import { useCallback, useEffect, useRef, useState } from 'react'
import { useTranslation } from 'react-i18next'

export default function MotionScene() {
  const { t } = useTranslation()
  const videoRef = useRef<HTMLVideoElement>(null)
  const reducedMotion = useRef(false)
  const [failed, setFailed] = useState(false)

  const start = useCallback(async () => {
    const video = videoRef.current
    if (!video || reducedMotion.current || document.hidden) return
    video.muted = true
    try {
      await video.play()
    } catch {
      /* Keep the poster when autoplay is restricted. */
    }
  }, [])

  useEffect(() => {
    const preference = window.matchMedia('(prefers-reduced-motion: reduce)')
    const syncPreference = () => {
      reducedMotion.current = preference.matches
      if (preference.matches) videoRef.current?.pause()
      else void start()
    }
    syncPreference()
    document.addEventListener('visibilitychange', start)
    document.addEventListener('pointerdown', start, { once: true })
    window.addEventListener('focus', start)
    preference.addEventListener('change', syncPreference)
    return () => {
      document.removeEventListener('visibilitychange', start)
      document.removeEventListener('pointerdown', start)
      window.removeEventListener('focus', start)
      preference.removeEventListener('change', syncPreference)
    }
  }, [start])

  return (
    <>
      <div className='hero-atmosphere' aria-hidden='true'>
        <div className='atmosphere-tone' />
        <div className='atmosphere-orbits'>
          <i />
          <i />
          <i />
        </div>
      </div>
      <div className='hero-art'>
        <video
          ref={videoRef}
          className='hero-video'
          src='/scaling-home/media/scaling-sphere.avc.mp4'
          muted
          loop
          playsInline
          preload='auto'
          poster='/scaling-home/media/scaling-poster.jpg'
          aria-label={t('White 3D sphere animation')}
          onLoadedData={() => {
            setFailed(false)
            void start()
          }}
          onError={() => setFailed(true)}
        />
        <div className='art-caption'>
          <span className='art-cross' aria-hidden='true'>
            +
          </span>
          <div>
            ENGINEERED TO CONNECT.
            <br />
            <span>DESIGNED TO SCALE.</span>
          </div>
        </div>
        {failed && (
          <p className='media-error' role='status'>
            {t('Animation could not load. Showing the preview image.')}
          </p>
        )}
      </div>
    </>
  )
}
