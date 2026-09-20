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

import { buildSeedanceExamples, getSeedanceModel } from '../lib/seedance-docs'

export function SeedanceDocsButton(props: { modelName: string }) {
  const { t } = useTranslation()
  const [open, setOpen] = useState(false)
  if (!getSeedanceModel(props.modelName)) return null
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
      {open && <SeedanceDocs modelName={props.modelName} />}
    </Dialog>
  )
}

export function SeedanceDocs(props: { modelName: string }) {
  const { t } = useTranslation()
  const model = getSeedanceModel(props.modelName)
  if (!model) return null
  const baseUrl = window.location.origin
  const parameters = model.parameters
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
        {model.description}{' '}
        {t(
          'Seedance videos are billed per request. The model marketplace and selected group show the current price.'
        )}
      </p>
      <p>
        {t(
          'Set NEW_API_KEY to your site API key. Use the id returned by creation as TASK_ID; replace example.com media URLs with your own accessible HTTPS files.'
        )}
      </p>
      <StaticDataTable
        data={[...parameters]}
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
      {buildSeedanceExamples(props.modelName).map((example) => (
        <section key={example.title}>
          <h3 className='mb-2 font-semibold'>{example.title}</h3>
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
