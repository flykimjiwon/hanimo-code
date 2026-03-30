import * as esbuild from 'esbuild';
import { readFileSync } from 'node:fs';

const pkg = JSON.parse(readFileSync('package.json', 'utf-8'));
const deps = Object.keys(pkg.dependencies ?? {});

await esbuild.build({
  entryPoints: ['src/cli.ts'],
  bundle: true,
  platform: 'node',
  target: 'node20',
  format: 'esm',
  outfile: 'dist/modol-bundle.mjs',
  sourcemap: false,
  minify: true,
  // Externalize native addons, React/Ink TUI (optional runtime dep),
  // and packages that don't bundle cleanly
  external: [
    'yoga-layout',
    'better-sqlite3',
    'fsevents',
    'ink',
    'react',
    'react-devtools-core',
  ],
  banner: {
    js: [
      // Polyfill __dirname and __filename for ESM
      'import { fileURLToPath as __fileURLToPath } from "node:url";',
      'import { dirname as __dirname_fn } from "node:path";',
      'const __filename = __fileURLToPath(import.meta.url);',
      'const __dirname = __dirname_fn(__filename);',
    ].join('\n'),
  },
});

console.log('Bundle built: dist/modol-bundle.mjs');
