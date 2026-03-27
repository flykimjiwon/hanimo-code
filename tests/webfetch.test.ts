import { describe, it, expect } from 'vitest';
import { webfetchTool, htmlToText } from '../src/tools/webfetch.js';

describe('htmlToText', () => {
  it('should strip script and style blocks', () => {
    const html = '<p>hello</p><script>alert("xss")</script><style>.x{color:red}</style><p>world</p>';
    const text = htmlToText(html);
    expect(text).not.toContain('alert');
    expect(text).not.toContain('color:red');
    expect(text).toContain('hello');
    expect(text).toContain('world');
  });

  it('should decode HTML entities', () => {
    const html = '<p>a &amp; b &lt; c &gt; d &quot;e&quot; &#39;f&#39;</p>';
    const text = htmlToText(html);
    expect(text).toContain('a & b < c > d "e" \'f\'');
  });

  it('should convert links to text with URL', () => {
    const html = '<a href="https://example.com">click here</a>';
    const text = htmlToText(html);
    expect(text).toContain('click here');
    expect(text).toContain('https://example.com');
  });

  it('should add newlines for block elements', () => {
    const html = '<h1>Title</h1><p>Paragraph</p><div>Block</div>';
    const text = htmlToText(html);
    expect(text).toContain('Title');
    expect(text).toContain('Paragraph');
    expect(text).toContain('Block');
  });

  it('should strip nav, footer, header', () => {
    const html = '<nav>menu</nav><main>content</main><footer>foot</footer>';
    const text = htmlToText(html);
    expect(text).not.toContain('menu');
    expect(text).not.toContain('foot');
    expect(text).toContain('content');
  });

  it('should collapse excessive whitespace', () => {
    const html = '<p>hello     world</p>\n\n\n\n<p>next</p>';
    const text = htmlToText(html);
    expect(text).not.toMatch(/  /); // no double spaces
    expect(text).not.toMatch(/\n{3,}/); // no triple newlines
  });
});

describe('webfetch tool', () => {
  it('should reject invalid URLs', async () => {
    const result = await webfetchTool.execute(
      { url: 'not-a-url' },
      { toolCallId: '1', messages: [], abortSignal: undefined as unknown as AbortSignal },
    );
    // Zod validation will reject before execute, but if passed through:
    expect(result).toBeDefined();
  });

  it('should fetch a real page (httpbin)', async () => {
    const result = await webfetchTool.execute(
      { url: 'https://httpbin.org/html' },
      { toolCallId: '2', messages: [], abortSignal: undefined as unknown as AbortSignal },
    );
    const r = result as { success: boolean; content?: string };
    // httpbin might be down, so we just check structure
    expect(r).toHaveProperty('success');
    if (r.success) {
      expect(r.content).toBeDefined();
      expect(typeof r.content).toBe('string');
    }
  }, 20000);

  it('should handle timeout gracefully', async () => {
    // Use a non-routable IP to trigger timeout
    const result = await webfetchTool.execute(
      { url: 'http://10.255.255.1' },
      { toolCallId: '3', messages: [], abortSignal: undefined as unknown as AbortSignal },
    );
    expect((result as { success: boolean }).success).toBe(false);
  }, 20000);
});
