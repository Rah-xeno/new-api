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
import { BookOpen } from 'lucide-react'
import { useState } from 'react'
import { useTranslation } from 'react-i18next'

import {
  CodeBlock,
  CodeBlockCopyButton,
} from '@/components/ai-elements/code-block'
import { StaticDataTable } from '@/components/data-table'
import { Dialog } from '@/components/dialog'
import { Button } from '@/components/ui/button'

import {
  buildMinimaxH3Examples,
  getMinimaxH3Resolution,
} from '../lib/minimax-h3-docs'

export function MinimaxH3DocsButton(props: { modelName: string }) {
  const { t } = useTranslation()
  const [open, setOpen] = useState(false)
  if (!getMinimaxH3Resolution(props.modelName)) return null
  return (
    <Dialog
      open={open}
      onOpenChange={setOpen}
      title={`${props.modelName} · ${t('View document')}`}
      description={t(
        'Create, poll and download a video using your site API key.'
      )}
      contentClassName='sm:max-w-4xl'
      trigger={
        <Button variant='outline' size='sm'>
          <BookOpen aria-hidden className='size-3.5' />
          {t('View document')}
        </Button>
      }
    >
      {open && <MinimaxH3Docs modelName={props.modelName} />}
    </Dialog>
  )
}

export function MinimaxH3Docs(props: { modelName: string }) {
  const { t } = useTranslation()
  const resolution = getMinimaxH3Resolution(props.modelName)
  if (!resolution) return null
  const baseUrl = window.location.origin
  const examples = buildMinimaxH3Examples(props.modelName)
  const parameters = [
    ['model', props.modelName],
    ['prompt', t('Required video description.')],
    [
      'duration',
      t(
        'Positive integer seconds; defaults to 5, gateway limit 3600. Common provider durations are 5–15 seconds; provider limits still apply.'
      ),
    ],
    ['resolution', resolution],
    ['aspect_ratio', '16:9 / 9:16 / 1:1'],
    [
      'generate_audio',
      t('Defaults to true; must be false for first and last frames.'),
    ],
    ['mode', 't2v / ref2v / fl2v'],
    [
      'reference_images',
      t(
        'Up to 5 accessible HTTPS image URLs with roles: reference_image, first_frame or last_frame.'
      ),
    ],
    [
      'reference_audios',
      t('Up to 3 HTTPS audio URLs. Reference videos are not supported.'),
    ],
  ]
  const lifecycle = [
    {
      title: t('Create response'),
      code: JSON.stringify(
        {
          id: 'video_TASK_ID',
          object: 'video',
          model: props.modelName,
          status: 'queued',
          progress: 0,
        },
        null,
        2
      ),
      language: 'json',
    },
    {
      title: t('Poll task status'),
      code: `curl '${baseUrl}/v1/videos/'"$TASK_ID" \\\n  -H "Authorization: Bearer $NEW_API_KEY"`,
      language: 'bash',
    },
    {
      title: t('Completed response'),
      code: JSON.stringify(
        {
          id: 'video_TASK_ID',
          object: 'video',
          model: props.modelName,
          status: 'completed',
          progress: 100,
        },
        null,
        2
      ),
      language: 'json',
    },
    {
      title: t('Download completed video'),
      code: `curl --fail --location '${baseUrl}/v1/videos/'"$TASK_ID"'/content' \\\n  -H "Authorization: Bearer $NEW_API_KEY" \\\n  -o result.mp4`,
      language: 'bash',
    },
    {
      title: t('Failed response'),
      code: JSON.stringify(
        {
          id: 'video_TASK_ID',
          object: 'video',
          model: props.modelName,
          status: 'failed',
          error: {
            code: 'video_generation_failed',
            message: 'The video generation task failed.',
          },
        },
        null,
        2
      ),
      language: 'json',
    },
  ]
  return (
    <div className='min-w-0 space-y-6 text-sm'>
      <p className='text-muted-foreground'>
        {t(
          'H3 videos are billed by duration, not by tokens or by request. Your current model and group prices apply; see the pricing overview.'
        )}
      </p>
      <p>
        {t(
          'Set NEW_API_KEY to your site API key. Use the id returned by creation as TASK_ID; do not create a second task when polling. Replace example.com media URLs with your own accessible files.'
        )}
      </p>
      <StaticDataTable
        data={parameters}
        getRowKey={(row) => row[0]}
        columns={[
          {
            id: 'name',
            header: t('Parameter'),
            cell: (row) => <code>{row[0]}</code>,
          },
          {
            id: 'description',
            header: t('Description'),
            cellClassName: 'whitespace-normal',
            cell: (row) => row[1],
          },
        ]}
      />
      {examples.map((example) => (
        <section key={example.title}>
          <h3 className='mb-2 font-semibold'>{t(example.title)}</h3>
          <CodeBlock
            language='bash'
            code={`curl '${baseUrl}/v1/videos' \\\n  -H "Authorization: Bearer $NEW_API_KEY" \\\n  -H "Content-Type: application/json" \\\n  -d '${JSON.stringify(example.body, null, 2)}'`}
          >
            <CodeBlockCopyButton />
          </CodeBlock>
        </section>
      ))}
      <p className='text-muted-foreground'>
        {t(
          'To use only reference images, omit reference_audios. First/last-frame mode requires both frame roles and generate_audio=false.'
        )}
      </p>
      <p>
        {t(
          'Poll every 5–10 seconds until completed or failed. Only completed tasks can be downloaded. A 4xx/5xx creation response is an error, not a task id.'
        )}
      </p>
      {lifecycle.map((step) => (
        <section key={step.title}>
          <h3 className='mb-2 font-semibold'>{step.title}</h3>
          <CodeBlock language={step.language} code={step.code}>
            <CodeBlockCopyButton />
          </CodeBlock>
        </section>
      ))}
    </div>
  )
}
