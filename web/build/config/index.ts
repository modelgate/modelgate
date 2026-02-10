import { readFileSync } from 'node:fs';

export function getAppVersion() {
  const pkgPath = new URL('../../package.json', import.meta.url);
  const pkg = JSON.parse(readFileSync(pkgPath, 'utf-8'));
  return pkg.version;
}

export * from './proxy';
export * from './time';
