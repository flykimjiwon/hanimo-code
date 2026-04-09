import { tool } from 'ai';
import { z } from 'zod';
import { htmlToText } from './webfetch.js';

const SEARCH_TIMEOUT_MS = 10000;
const MAX_RESULTS = 10;

export interface SearchResult {
  title: string;
  url: string;
  snippet: string;
}

/**
 * Parse DuckDuckGo HTML search results into structured data.
 * DuckDuckGo HTML endpoint returns results in <a class="result__a"> + <a class="result__snippet"> pattern.
 */
export function parseDDGResults(html: string, maxResults: number): SearchResult[] {
  const results: SearchResult[] = [];

  // DuckDuckGo HTML uses .result__a for title links and .result__snippet for descriptions
  const resultBlocks = html.split(/class="result\s/g).slice(1);

  for (const block of resultBlocks) {
    if (results.length >= maxResults) break;

    // Extract URL and title from result__a
    const linkMatch = block.match(/class="result__a"[^>]*href="([^"]*)"[^>]*>([\s\S]*?)<\/a>/);
    if (!linkMatch) continue;

    let url = linkMatch[1] ?? '';
    const titleHtml = linkMatch[2] ?? '';

    // DuckDuckGo wraps URLs in a redirect — extract the actual URL
    const uddgMatch = url.match(/uddg=([^&]+)/);
    if (uddgMatch) {
      url = decodeURIComponent(uddgMatch[1] ?? '');
    }

    if (!url || !url.startsWith('http')) continue;

    // Extract snippet
    const snippetMatch = block.match(/class="result__snippet"[^>]*>([\s\S]*?)<\/(?:a|span|div)>/);
    const snippetHtml = snippetMatch?.[1] ?? '';

    const title = htmlToText(titleHtml).trim();
    const snippet = htmlToText(snippetHtml).trim();

    if (title) {
      results.push({ title, url, snippet });
    }
  }

  return results;
}

export const webSearchTool = tool({
  description:
    'Search the web and return relevant results. Use this to find documentation, ' +
    'solutions to errors, package info, API references, etc. Returns titles, URLs, and snippets.',
  parameters: z.object({
    query: z.string().describe('Search query'),
    maxResults: z
      .number()
      .min(1)
      .max(MAX_RESULTS)
      .default(5)
      .describe('Max results to return (1-10)'),
  }),
  execute: async ({ query, maxResults }) => {
    try {
      const encodedQuery = encodeURIComponent(query);
      const searchUrl = `https://html.duckduckgo.com/html/?q=${encodedQuery}`;

      const response = await fetch(searchUrl, {
        headers: {
          'User-Agent': 'hanimo/0.1 (AI coding assistant)',
          Accept: 'text/html',
        },
        signal: AbortSignal.timeout(SEARCH_TIMEOUT_MS),
        redirect: 'follow',
      });

      if (!response.ok) {
        return {
          success: false,
          error: `Search failed: HTTP ${response.status} ${response.statusText}`,
        };
      }

      const html = await response.text();
      const results = parseDDGResults(html, maxResults);

      if (results.length === 0) {
        return {
          success: true,
          query,
          results: [],
          message: 'No results found. Try different search terms.',
        };
      }

      return {
        success: true,
        query,
        resultCount: results.length,
        results,
      };
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : String(err);
      if (message.includes('TimeoutError') || message.includes('aborted')) {
        return { success: false, error: `Search timed out after ${SEARCH_TIMEOUT_MS / 1000}s` };
      }
      return { success: false, error: `Search failed: ${message}` };
    }
  },
});
