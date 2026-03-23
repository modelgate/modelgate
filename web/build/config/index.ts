import { readFileSync } from 'node:fs';

export function getAppVersion() {
  // 优先使用环境变量（GitHub Action/Docker 构建时传入）
  if (process.env.APP_VERSION) {
    return process.env.APP_VERSION;
  }
  // 否则从 package.json 读取
  const pkgPath = new URL('../../package.json', import.meta.url);
  const pkg = JSON.parse(readFileSync(pkgPath, 'utf-8'));
  return pkg.version;
}

export * from './proxy';
export * from './time';
