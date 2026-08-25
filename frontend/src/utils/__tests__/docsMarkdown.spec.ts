import { describe, expect, it } from 'vitest'
import { renderDocumentMarkdown } from '@/utils/docsMarkdown'

describe('renderDocumentMarkdown', () => {
  it('keeps public image links and image layout attributes', () => {
    const source = '<img src="/api/v1/docs/guide/images/example.png" alt="Example" data-align="right" data-width="medium" />'

    const rendered = renderDocumentMarkdown(source, 'guide')

    expect(rendered.html).toContain('src="/api/v1/docs/guide/images/example.png"')
    expect(rendered.html).toContain('data-align="right"')
    expect(rendered.html).toContain('data-width="medium"')
  })

  it('continues to resolve legacy image filenames through the public endpoint', () => {
    const rendered = renderDocumentMarkdown('![Legacy](example.png)', 'guide')

    expect(rendered.html).toContain('src="/api/v1/docs/guide/images/example.png"')
  })

  it('removes unsafe image event handlers', () => {
    const rendered = renderDocumentMarkdown('<img src="/safe.png" onerror="alert(1)" />', 'guide')

    expect(rendered.html).toContain('src="/safe.png"')
    expect(rendered.html).not.toContain('onerror')
  })
})
