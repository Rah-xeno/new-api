import { useTranslation } from 'react-i18next'

export function ScalingModels() {
  const { t } = useTranslation()
  const families = [
    {
      name: 'OpenAI',
      model: 'GPT',
      index: '01',
      summary: t('From everyday questions to your next product.'),
      tags: [t('General conversation'), t('App development')],
      glyph: 'O',
    },
    {
      name: 'Anthropic',
      model: 'Claude',
      index: '02',
      summary: t('Bring words, reasoning, and code together.'),
      tags: [t('Content writing'), t('Coding collaboration')],
      glyph: 'C',
    },
    {
      name: 'xAI',
      model: 'Grok',
      index: '03',
      summary: t('Another choice for exploration and creation.'),
      tags: [t('Creative exploration'), t('Conversation assistant')],
      glyph: 'X',
    },
    {
      name: 'Google',
      model: 'Gemini',
      index: '04',
      summary: t('Bring multimodal ideas into your workflow.'),
      tags: [t('Multimodal'), t('Workflows')],
      glyph: 'G',
    },
  ]

  return (
    <section
      id='model-plaza'
      className='section-shell model-section'
      aria-labelledby='models-title'
    >
      <div className='section-heading'>
        <div>
          <p className='eyeline'>01 / MODEL PLAZA</p>
          <h2 id='models-title'>{t('Great models. More than one choice.')}</h2>
        </div>
        <p className='section-description'>
          {t('Familiar model families. One integration.')}
          <br />
          {t('Choose for your business, not another integration.')}
        </p>
      </div>
      <div className='model-toolbar'>
        <span className='model-count'>
          {t('Model families')} <span>04</span>
        </span>
      </div>
      <div className='model-grid'>
        {families.map((family) => (
          <article className='model-family' key={family.model}>
            <div className='flex items-center justify-between'>
              <span className='family-symbol' aria-hidden='true'>
                {family.glyph}
              </span>
              <span className='family-index'>/{family.index}</span>
            </div>
            <p className='family-provider'>{family.name}</p>
            <h3>{family.model}</h3>
            <p className='family-summary'>{family.summary}</p>
            <div className='flex flex-wrap gap-2'>
              {family.tags.map((tag) => (
                <span className='family-tag' key={tag}>
                  {tag}
                </span>
              ))}
            </div>
          </article>
        ))}
      </div>
      <p className='catalog-note'>
        {t(
          'Model families shown here. Check the console for available models and pricing.'
        )}
      </p>
    </section>
  )
}
