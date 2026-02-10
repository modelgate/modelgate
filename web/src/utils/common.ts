import { $t } from '@/locales';
import type { Timestamp } from "@bufbuild/protobuf/wkt";
import dayjs from 'dayjs';

/**
 * Transform record to option
 *
 * @example
 *   ```ts
 *   const record = {
 *     key1: 'label1',
 *     key2: 'label2'
 *   };
 *   const options = transformRecordToOption(record);
 *   // [
 *   //   { value: 'key1', label: 'label1' },
 *   //   { value: 'key2', label: 'label2' }
 *   // ]
 *   ```;
 *
 * @param record
 */
export function transformRecordToOption<T extends Record<string | number, string>>(record: T) {
  return Object.entries(record).map(([value, label]) => ({
    value,
    label
  })) as CommonType.Option<keyof T>[];
}

/**
 * 将枚举 + i18n 映射转换为 options 列表
 *
 * @param enumObj 枚举对象
 * @param labelMap 枚举值到 label（i18n key）的映射
 * @param excludeValues 可选，指定要排除的枚举值数组
 */
export function enumToOptions<T extends Record<string, string | number>>(
  enumObj: T,
  labelMap: Record<string | number, string>,
  excludeValues: (string | number)[] = []
): CommonType.Option<string | number>[] {
  const numericKeys = Object.keys(enumObj)
    .filter(k => !Number.isNaN(Number(k)))
    .map(k => Number(k))
    .filter(val => !excludeValues.includes(val));

  return numericKeys.map(value => ({
    value,
    label: labelMap[value] ?? String(value)
  }));
}

/**
 * Translate options
 *
 * @param options
 */
export function translateOptions(options: CommonType.Option<string>[]) {
  return options.map(option => ({
    ...option,
    label: $t(option.label as App.I18n.I18nKey)
  }));
}

/**
 * Toggle html class
 *
 * @param className
 */
export function toggleHtmlClass(className: string) {
  function add() {
    document.documentElement.classList.add(className);
  }

  function remove() {
    document.documentElement.classList.remove(className);
  }

  return {
    add,
    remove
  };
}

/**
 * Get file size
 *
 * @param size
 */
export function formatFileSize(filesize: number) {
  const units = ['B', 'KB', 'MB', 'GB', 'TB'];
  let size = filesize;
  let i = 0;
  while (size >= 1024 && i < units.length - 1) {
    size /= 1024;
    i += 1;
  }
  return `${size.toFixed(2)} ${units[i]}`;
}

/**
 * Format proto timestamp to string
 *
 * @param ts
 * @param format
 */
export function formatProtoTime(ts: Timestamp | undefined | null, format = 'YYYY-MM-DD HH:mm:ss') {
  if (!ts) {
    return '-';
  }
  const ms = protoToMs(ts);
  return ms ? dayjs(ms).format(format) : '-';
}

/**
 * Convert proto timestamp to milliseconds
 *
 * @param ts
 */
export function protoToMs(ts: Timestamp | undefined | null): number | null {
  if (!ts) {
    return null;
  }
  const ms = ts.seconds * BigInt(1000) + BigInt(ts.nanos || 0) / BigInt(1e6);
  return Number(ms);
}

/**
 * Convert milliseconds to proto timestamp
 *
 * @param ms
 */
export function msToProto(ms: number | undefined | null): Timestamp | undefined {
  if (ms === undefined || ms === null) {
    return undefined;
  }
  return {
    $typeName: 'google.protobuf.Timestamp',
    seconds: BigInt(Math.floor(ms / 1000)),
    nanos: (ms % 1000) * 1e6
  } as Timestamp;
}

/**
 * Currency locale map
 */
const currencyLocaleMap: Record<string, string> = {
  CNY: 'zh-CN',
  USD: 'en-US',
  EUR: 'de-DE',
  JPY: 'ja-JP'
};

/**
 * Format currency
 *
 * @param value
 * @param currency
 */
export function formatCurrency(value: number, currency: string) {
  return new Intl.NumberFormat(currencyLocaleMap[currency] || 'en-US', {
    style: 'currency',
    currency
  }).format(value);
}