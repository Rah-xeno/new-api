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
// This is the resolution-suffixed video API, not the official MiniMax-H3 API.
export const MINIMAX_H3_RESOLUTIONS: Readonly<Record<string, string>> = {
  'Minimax-h3-480p': '480p',
  'Minimax-h3-720p': '720p',
  'Minimax-h3-1080p': '1080p',
  'Minimax-h3-2k': '2k',
}

export function getMinimaxH3Resolution(modelName: string): string | undefined {
  return Object.hasOwn(MINIMAX_H3_RESOLUTIONS, modelName)
    ? MINIMAX_H3_RESOLUTIONS[modelName]
    : undefined
}

export function buildMinimaxH3Examples(modelName: string) {
  const resolution = getMinimaxH3Resolution(modelName)
  if (!resolution) return []
  const common = { model: modelName, duration: 5, resolution }
  return [
    {
      title: 'Text to video',
      body: {
        ...common,
        prompt: 'A slow camera move over a mountain lake at sunrise.',
        aspect_ratio: '16:9',
        generate_audio: true,
        mode: 't2v',
      },
    },
    {
      title: 'Reference images and audio',
      body: {
        ...common,
        prompt: 'Animate the subject with gentle movement and matching sound.',
        mode: 'ref2v',
        reference_images: [
          { url: 'https://example.com/reference.png', role: 'reference_image' },
        ],
        reference_audios: ['https://example.com/reference.mp3'],
      },
    },
    {
      title: 'First and last frames',
      body: {
        ...common,
        prompt: 'Create a smooth transition between the two frames.',
        mode: 'fl2v',
        generate_audio: false,
        reference_images: [
          { url: 'https://example.com/first.png', role: 'first_frame' },
          { url: 'https://example.com/last.png', role: 'last_frame' },
        ],
      },
    },
  ]
}
