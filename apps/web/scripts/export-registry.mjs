import { execFileSync } from 'node:child_process';
import { writeFileSync } from 'node:fs';

// Export public endpoint contracts only. This command does not authenticate or invoke APIs.
const root = new URL('../../../', import.meta.url);
const { data } = JSON.parse(execFileSync('go', ['run', './cmd/agc', 'endpoints'], {
  cwd: root,
  encoding: 'utf8',
  maxBuffer: 5 * 1024 * 1024,
}));
writeFileSync(new URL('../src/endpoint-registry.json', import.meta.url), `${JSON.stringify(data, null, 2)}\n`);
console.log(`Exported ${data.length} endpoint contracts.`);
