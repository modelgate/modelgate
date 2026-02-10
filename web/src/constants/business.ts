import { enumToOptions, transformRecordToOption } from '@/utils/common';
import { IconType, MenuType } from '@/typings/proto/model/system/menu_pb';
import type { Api } from '@/typings/api';

export const enableStatusRecord: Record<string, App.I18n.I18nKey> = {
  'enabled': 'common.status.enabled',
  'disabled': 'common.status.disabled'
};

export const enableStatusOptions = transformRecordToOption(enableStatusRecord);

export const yesOrNoRecord: Record<number, App.I18n.I18nKey> = {
  1: 'common.yesOrNo.yes',
  0: 'common.yesOrNo.no'
};

export const apiMethodRecord: Record<string, App.I18n.I18nKey> = {
  'GET': 'common.apiMethod.get',
  'POST': 'common.apiMethod.post',
  'PUT': 'common.apiMethod.put',
  'DELETE': 'common.apiMethod.delete',
  'PATCH': 'common.apiMethod.patch'
};

export const apiMethodOptions = transformRecordToOption(apiMethodRecord);

export const yesOrNoOptions = transformRecordToOption(yesOrNoRecord);

export const modelStatusRecord: Record<string, App.I18n.I18nKey> = {
  'enabled': 'page.relay.common.model.status.enabled',
  'disabled': 'page.relay.common.model.status.disabled',
  'deprecated': 'page.relay.common.model.status.deprecated'
};

export const modelStatusOptions = transformRecordToOption(modelStatusRecord);

export const requestStatusRecord: Record<string, App.I18n.I18nKey> = {
  'pending': 'page.usage.common.request.status.pending',
  'success': 'page.usage.common.request.status.success',
  'failed': 'page.usage.common.request.status.failed',
  'cancelled': 'page.usage.common.request.status.cancelled',
};

export const requestStatusOptions = transformRecordToOption(requestStatusRecord);

export const ledgerTypeRecord: Record<string, App.I18n.I18nKey> = {
  'consume': 'page.usage.common.ledger.type.consume',
  'refund': 'page.usage.common.ledger.type.refund',
  'charge': 'page.usage.common.ledger.type.charge',
  'adjust': 'page.usage.common.ledger.type.adjust',
};

export const ledgerTypeOptions = transformRecordToOption(ledgerTypeRecord);

export const apiKeyStatusRecord: Record<string, App.I18n.I18nKey> = {
  'enabled': 'page.relay.common.apiKey.status.enabled',
  'disabled': 'page.relay.common.apiKey.status.disabled',
  'cooldown': 'page.relay.common.apiKey.status.cooldown',
  'revoked': 'page.relay.common.apiKey.status.revoked'
};

export const apiKeyStatusOptions = transformRecordToOption(apiKeyStatusRecord);

export const currencyRecord: Record<string, App.I18n.I18nKey> = {
  'CNY': 'page.relay.common.currency.cny',
  'USD': 'page.relay.common.currency.usd',
  'POINT': 'page.relay.common.currency.point'
};

export const currencyOptions = transformRecordToOption(currencyRecord);

export const genderRecord: Record<string, App.I18n.I18nKey> = {
  'male': 'page.manage.user.userGender.male',
  'female': 'page.manage.user.userGender.female',
  'unknown': 'page.manage.user.userGender.unknown'
};

export const genderOptions = transformRecordToOption(genderRecord);


export const userGenderRecord: Record<Api.SystemManage.UserGender, App.I18n.I18nKey> = {
  male: 'page.manage.user.userGender.male',
  female: 'page.manage.user.userGender.female',
  unknown: 'page.manage.user.userGender.unknown'
};

export const userGenderOptions = transformRecordToOption(userGenderRecord);

export const menuTypeRecord: Record<Exclude<MenuType, MenuType.UNSPECIFIED>, App.I18n.I18nKey> = {
  [MenuType.DIRECTORY]: 'page.manage.menu.type.directory',
  [MenuType.MENU]: 'page.manage.menu.type.menu'
};

export const menuTypeOptions = enumToOptions(MenuType, menuTypeRecord, [MenuType.UNSPECIFIED]);

export const menuIconTypeRecord: Record<Exclude<IconType, IconType.UNSPECIFIED>, App.I18n.I18nKey> = {
  [IconType.ICONIFY]: 'page.manage.menu.iconType.iconify',
  [IconType.LOCAL]: 'page.manage.menu.iconType.local'
};

export const menuIconTypeOptions = enumToOptions(IconType, menuIconTypeRecord, [MenuType.UNSPECIFIED]);
