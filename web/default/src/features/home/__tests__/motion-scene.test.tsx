import { fireEvent, render, screen } from '@testing-library/react'
import i18n from 'i18next'
import { afterEach, beforeEach, expect, it, vi } from 'vitest'

import MotionScene from '../motion-scene'

beforeEach(async () => {
  await i18n.changeLanguage('en')
  vi.spyOn(HTMLMediaElement.prototype, 'play').mockResolvedValue()
  vi.spyOn(HTMLMediaElement.prototype, 'pause').mockImplementation(() => {})
})
afterEach(() => vi.restoreAllMocks())
it('starts muted inline looping playback without exposing a pause button', () => {
  render(<MotionScene />)
  const video = screen.getByLabelText('White 3D sphere animation')
  expect(video).toHaveAttribute('playsinline')
  expect(video).toHaveAttribute('loop')
  expect(HTMLMediaElement.prototype.play).toHaveBeenCalled()
  expect(screen.queryByRole('button')).toBeNull()
})
it('respects reduced motion on first load instead of starting playback', () => {
  const original = window.matchMedia
  vi.spyOn(window, 'matchMedia').mockImplementation((query) => ({
    ...original(query),
    matches: true,
  }))
  render(<MotionScene />)
  expect(HTMLMediaElement.prototype.play).not.toHaveBeenCalled()
  expect(HTMLMediaElement.prototype.pause).toHaveBeenCalled()
})
it('keeps the poster visible and reports a media loading failure', () => {
  render(<MotionScene />)
  const video = screen.getByLabelText('White 3D sphere animation')
  fireEvent.error(video)
  expect(video).toHaveAttribute(
    'poster',
    '/scaling-home/media/scaling-poster.jpg'
  )
  expect(screen.getByRole('status')).toHaveTextContent(
    'Animation could not load. Showing the preview image.'
  )
})
