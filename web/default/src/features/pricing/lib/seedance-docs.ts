/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

For commercial licensing, please contact support@quantumnous.com
*/

export type SeedanceModel = {
  name: string
  description: string
  parameters: ReadonlyArray<readonly [string, string]>
  examples: ReadonlyArray<{
    title: string
    body: Record<string, unknown>
  }>
}

const SIZE_NOTE =
  '16:9=1280x720，9:16=720x1280，1:1=720x720，4:3=960x720，3:4=720x960，21:9=1680x720。'

const COMMON_PARAMETERS: ReadonlyArray<readonly [string, string]> = [
  ['model', '模型名称。'],
  ['prompt', '必填，视频内容描述。'],
  ['seconds', '视频时长，使用字符串（例如 "5"）。'],
  ['size', SIZE_NOTE],
]

const REFERENCE_PARAMETERS: ReadonlyArray<readonly [string, string]> = [
  ['images', '可选，公网 HTTPS 图片 URL 数组。'],
  ['reference_videos', '可选，公网 HTTPS 视频 URL 数组。'],
  ['reference_audios', '可选，公网 HTTPS 音频 URL 数组。'],
]

const pair = (key: string, value: string) => [key, value] as const

function referenceDescription(key: string): readonly [string, string] {
  if (key === 'images') {
    return pair(key, '可选，最多 9 张公网 HTTPS 图片 URL。')
  }
  if (key === 'reference_videos') {
    return pair(key, '可选，最多 3 条公网 HTTPS 视频 URL。')
  }
  return pair(key, '可选，最多 3 条公网 HTTPS 音频 URL。')
}

const MODEL_MAP: Record<string, SeedanceModel> = {
  'sd-2-c4': {
    name: 'sd-2-c4',
    description:
      'Seedance 2：支持最多 9 张图片、3 条视频和 3 条音频参考，4–15 秒，全比例生成。',
    parameters: [...COMMON_PARAMETERS, ...REFERENCE_PARAMETERS],
    examples: [
      {
        title: '图生视频（多参考素材）',
        body: {
          model: 'sd-2-c4',
          prompt: '让人物自然地向镜头挥手，保持原有风格。',
          seconds: '5',
          size: '1280x720',
          images: ['https://example.com/reference.png'],
          reference_videos: ['https://example.com/motion.mp4'],
          reference_audios: ['https://example.com/music.mp3'],
        },
      },
    ],
  },
  'sd10-seedance-2.5': {
    name: 'sd10-seedance-2.5',
    description:
      'Seedance 2.5（SD10 卡脸线路）：固定 30 秒，支持文生、图生和最多 30 张参考图。',
    parameters: [
      ...COMMON_PARAMETERS.filter(([key]) => key !== 'size').map(
        ([key, value]) =>
          key === 'seconds'
            ? pair(key, '固定 30 秒，传字符串 "30"。')
            : pair(key, value)
      ),
      ['images', '可选，最多 30 张公网 HTTPS 参考图。'],
    ],
    examples: [
      {
        title: '固定 30 秒视频',
        body: {
          model: 'sd10-seedance-2.5',
          prompt: '一段电影感的城市夜景运镜。',
          seconds: '30',
          images: ['https://example.com/reference.png'],
        },
      },
    ],
  },
  'seedance-2.0': {
    name: 'seedance-2.0',
    description:
      'Seedance 2.0 视频模型，支持文生、图生以及首尾帧等多模态输入，4–15 秒。',
    parameters: [
      ...COMMON_PARAMETERS,
      ['aspect_ratio', '可选，16:9 / 9:16 / 1:1 等比例。'],
      ['resolution', '可选，480p 或 720p。'],
      ...REFERENCE_PARAMETERS,
    ],
    examples: [
      {
        title: '文生视频',
        body: {
          model: 'seedance-2.0',
          prompt: '日出时分，镜头缓慢掠过山间湖泊。',
          seconds: '5',
          size: '1280x720',
          aspect_ratio: '16:9',
          resolution: '720p',
        },
      },
    ],
  },
  '特价seedance-2.5-720p': {
    name: '特价seedance-2.5-720p',
    description: '特价 Seedance 2.5 720p，OpenAI 兼容视频接口，按次计费。',
    parameters: [
      ...COMMON_PARAMETERS,
      ...REFERENCE_PARAMETERS.map(([key]) => referenceDescription(key)),
    ],
    examples: [
      {
        title: '720p 视频',
        body: {
          model: '特价seedance-2.5-720p',
          prompt: '产品展示视频，镜头平滑推进。',
          seconds: '5',
          size: '1280x720',
          images: ['https://example.com/reference.png'],
        },
      },
    ],
  },
}

export const SEEDANCE_MODEL_NAMES = Object.keys(
  MODEL_MAP
) as ReadonlyArray<string>

export function getSeedanceModel(modelName: string): SeedanceModel | undefined {
  return MODEL_MAP[modelName]
}

export function buildSeedanceExamples(modelName: string) {
  return getSeedanceModel(modelName)?.examples ?? []
}
