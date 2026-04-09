import { describe, it, expect } from 'vitest';
import { parseDDGResults } from '../src/tools/web-search.js';

describe('web-search', () => {
  it('should parse DuckDuckGo HTML results', () => {
    const html = `
      <div class="result results_links results_links_deep web-result">
        <div class="links_main links_deep result__body">
          <a class="result__a" href="https://duckduckgo.com/l/?uddg=https%3A%2F%2Fexample.com%2Fdocs&rut=abc">
            Example Documentation
          </a>
          <a class="result__snippet">This is a snippet about example docs.</a>
        </div>
      </div>
      <div class="result results_links results_links_deep web-result">
        <div class="links_main links_deep result__body">
          <a class="result__a" href="https://duckduckgo.com/l/?uddg=https%3A%2F%2Fother.com&rut=def">
            Other Site
          </a>
          <a class="result__snippet">Another snippet here.</a>
        </div>
      </div>
    `;

    const results = parseDDGResults(html, 5);
    expect(results.length).toBe(2);
    expect(results[0]?.title).toBe('Example Documentation');
    expect(results[0]?.url).toBe('https://example.com/docs');
    expect(results[0]?.snippet).toContain('snippet');
    expect(results[1]?.url).toBe('https://other.com');
  });

  it('should respect maxResults limit', () => {
    const html = `
      <div class="result "><a class="result__a" href="https://duckduckgo.com/l/?uddg=https%3A%2F%2Fa.com">A</a><a class="result__snippet">A</a></div>
      <div class="result "><a class="result__a" href="https://duckduckgo.com/l/?uddg=https%3A%2F%2Fb.com">B</a><a class="result__snippet">B</a></div>
      <div class="result "><a class="result__a" href="https://duckduckgo.com/l/?uddg=https%3A%2F%2Fc.com">C</a><a class="result__snippet">C</a></div>
    `;
    const results = parseDDGResults(html, 2);
    expect(results.length).toBe(2);
  });

  it('should return empty array for invalid HTML', () => {
    const results = parseDDGResults('<html><body>No results</body></html>', 5);
    expect(results).toEqual([]);
  });

  it('should handle HTML with no valid URLs', () => {
    const html = `
      <div class="result ">
        <a class="result__a" href="javascript:void(0)">Bad Link</a>
        <a class="result__snippet">No URL</a>
      </div>
    `;
    const results = parseDDGResults(html, 5);
    expect(results).toEqual([]);
  });
});
