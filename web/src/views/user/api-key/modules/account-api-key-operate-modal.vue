<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { apiKeyStatusOptions } from '@/constants/business';
import { useFormRules, useNaiveForm } from '@/hooks/common/form';
import { relayServiceClient } from '@/grpc';
import { $t } from '@/locales';
import type { AccountApiKey } from '@/typings/proto/model/relay/account_api_key_pb';
import { NaiveUI } from '@/typings/naive-ui';
import { msToProto, protoToMs } from '@/utils/common';

/** Key display modal */
const showKeyModal = ref(false);
const generatedKey = ref('');
const keyCopied = ref(false);

defineOptions({
  name: 'UserApiKeyOperateModal'
});

interface Props {
  /** the type of operation */
  operateType: NaiveUI.TableOperateType;
  /** the edit row data */
  rowData?: AccountApiKey | null;
}

const props = defineProps<Props>();

interface Emits {
  (e: 'submitted'): void;
}

const emit = defineEmits<Emits>();

const visible = defineModel<boolean>('visible', {
  default: false
});

const { formRef, validate, restoreValidation } = useNaiveForm();
const { defaultRequiredRule } = useFormRules();

const title = computed(() => {
  const titles: Record<NaiveUI.TableOperateType, string> = {
    add: $t('page.user.apiKey.addApiKey'),
    edit: $t('page.user.apiKey.editApiKey')
  };
  return titles[props.operateType];
});

interface AccountApiKeyForm {
  id: number;
  accountId: number | null;
  keyName: string;
  scope: string;
  quoteLimit: number;
  rateLimit: number;
  expiredAt: number | null;
  status: string;
  remark: string;
}

const model = ref(createDefaultModel());

function createDefaultModel(): AccountApiKeyForm {
  return {
    id: 0,
    accountId: null,
    keyName: '',
    scope: '',
    quoteLimit: 0,
    rateLimit: 0,
    expiredAt:  null,
    status: 'enabled',
    remark: ''
  };
}

type RuleKey = Extract<keyof AccountApiKeyForm, 'accountId' | 'keyName' | 'scope' | 'status'>;

const rules: Record<RuleKey, App.Global.FormRule> = {
  accountId: defaultRequiredRule,
  keyName: defaultRequiredRule,
  scope: defaultRequiredRule,
  status: defaultRequiredRule
};

/** scope options */
const scopeOptions = ref<CommonType.Option<string>[]>( [
  { label: 'Chat', value: 'chat' },
  { label: 'Image', value: 'image' },
  { label: 'All', value: 'all' }
]);

/** user options */
const accountOptions = ref<CommonType.Option<Number>[]>([]);

async function getAccountOptions() {
  const resp = await relayServiceClient.getAccountList({
    current: 1,
    size: 10000,
    status: 'enabled'
  });
  const options = resp.records.map(item => ({
    label: `${item.nickname} (${item.name})`,
    value: Number(item.id)
  }));

  accountOptions.value = [...options];
}

function handleInitModel() {
  model.value = createDefaultModel();

  if (props.operateType === 'edit' && props.rowData) {
    const row = props.rowData;
    model.value = {
      id: Number(row.id),
      accountId: Number(row.accountId),
      keyName: row.keyName,
      scope: row.scope,
      quoteLimit: Number(row.quoteLimit),
      rateLimit: Number(row.rateLimit),
      expiredAt: protoToMs(row.expiredAt),
      status: row.status,
      remark: row.remark
    };
  }
}

function closeDrawer() {
  visible.value = false;
}

async function handleSubmit() {
  await validate();

  const submissionData = {
    keyName: model.value.keyName,
    scope: model.value.scope,
    status: model.value.status,
    remark: model.value.remark,
    id: BigInt(model.value.id ?? 0),
    accountId: BigInt(model.value.accountId ?? 0),
    quoteLimit: BigInt(model.value.quoteLimit ?? 0),
    rateLimit: BigInt(model.value.rateLimit ?? 0),
    expiredAt: msToProto(model.value.expiredAt)
  };

  if (props.operateType === 'edit') {
    const paths = ['account_id', 'key_name', 'scope', 'quote_limit', 'rate_limit', 'expired_at', 'status', 'remark'];

    try {
      await relayServiceClient.updateAccountApiKey({
        updateMask: {
          paths
        },
        accountApiKey: submissionData as any
      });
      window.$message?.success($t('common.updateSuccess'));
    } catch {
      window.$message?.error($t('common.updateFailed'));
    }

    emit('submitted');
  } else {
    try {
      const result = await relayServiceClient.createAccountApiKey({
        accountApiKey: submissionData as any
      });

      generatedKey.value = result.key;
      keyCopied.value = false;

      // Close the form drawer first
      visible.value = false;

      // Then show the key modal
      showKeyModal.value = true;
    } catch {
      window.$message?.error($t('common.addFailed'));
    }
  }
}

function copyKey() {
  navigator.clipboard.writeText(generatedKey.value);
  keyCopied.value = true;
  window.$message?.success($t('common.copySuccess'));
}

function closeKeyModal() {
  showKeyModal.value = false;
  emit('submitted');
}

watch(visible, () => {
  if (visible.value) {
    handleInitModel();
    restoreValidation();
    getAccountOptions();
  }
});
</script>

<template>
  <NModal v-model:show="visible" :title="title" preset="card" class="w-600px">
    <NScrollbar class="h-500px pr-20px">
      <NForm ref="formRef" :model="model" :rules="rules">
        <NFormItem :label="$t('page.user.apiKey.accountName')" path="accountId">
          <NSelect
            v-model:value="model.accountId"
            :options="accountOptions"
            filterable
            :placeholder="$t('page.user.apiKey.form.accountId')"
            :disabled="operateType === 'edit'"
          />
        </NFormItem>
        <NFormItem :label="$t('page.user.apiKey.keyName')" path="keyName">
          <NInput
            v-model:value="model.keyName"
            :placeholder="$t('page.user.apiKey.form.keyName')"
            :disabled="operateType === 'edit'"
          />
        </NFormItem>
        <NFormItem :label="$t('page.user.apiKey.scope')" path="scope">
          <NSelect
            v-model:value="model.scope"
            :options="scopeOptions"
            :placeholder="$t('page.user.apiKey.form.scope')"
          />
        </NFormItem>
        <NGrid :cols="2" :x-gap="16">
          <NFormItemGi :label="$t('page.user.apiKey.expiredAt')" path="expiredAt">
          <NDatePicker
            v-model:value="model.expiredAt"
            type="datetime"
            :placeholder="$t('page.user.apiKey.form.expiredAt')"
            class="w-full"
            clearable
          />
          </NFormItemGi>
        </NGrid>
        <NGrid :cols="2" :x-gap="16">
          <NFormItemGi :label="$t('page.user.apiKey.quoteLimit')" path="quoteLimit">
            <NInputNumber
              v-model:value="model.quoteLimit"
              :placeholder="$t('page.user.apiKey.form.quoteLimit')"
              :min="0"
            class="w-full"
          />
        </NFormItemGi>
        <NFormItemGi :label="$t('page.user.apiKey.rateLimit')" path="rateLimit">
          <NInputNumber
            v-model:value="model.rateLimit"
            :placeholder="$t('page.user.apiKey.form.rateLimit')"
            :min="0"
            class="w-full"
          />
        </NFormItemGi>
        </NGrid>
        <NFormItem :label="$t('page.user.apiKey.status')" path="status">
          <NRadioGroup v-model:value="model.status">
            <NRadio v-for="item in apiKeyStatusOptions" :key="item.value" :value="item.value" :label="$t(item.label)" />
          </NRadioGroup>
        </NFormItem>
      </NForm>
    </NScrollbar>
    <template #footer>
      <NSpace :size="16">
        <NButton @click="closeDrawer">{{ $t('common.cancel') }}</NButton>
        <NButton type="primary" @click="handleSubmit">{{ $t('common.confirm') }}</NButton>
      </NSpace>
    </template>
  </NModal>

  <!-- Key Display Modal -->
  <NModal v-model:show="showKeyModal" preset="card" :title="$t('page.user.apiKey.keyModal.title')" class="w-500px">
    <NSpace vertical :size="20">
      <NAlert type="warning" :closable="false">
        <template #header>
          <span class="font-bold">{{ $t('page.user.apiKey.keyModal.warningTitle') }}</span>
        </template>
        {{ $t('page.user.apiKey.keyModal.warningMessage') }}
      </NAlert>

      <div>
        <div class="mb-2 flex items-center justify-between">
          <span class="text-sm text-gray-600">{{ $t('page.user.apiKey.keyModal.yourKey') }}</span>
          <NButton
            tertiary
            size="small"
            @click="copyKey"
          >
            <template #icon>
              <icon-mdi-content-copy v-if="!keyCopied" />
              <icon-mdi-check v-else class="text-green-500" />
            </template>
            {{ keyCopied ? $t('page.user.apiKey.keyModal.copied') : $t('page.user.apiKey.keyModal.copy') }}
          </NButton>
        </div>
        <NInput
          :value="generatedKey"
          readonly
          type="textarea"
          :rows="3"
          class="font-mono text-sm"
        />
      </div>
    </NSpace>
    <template #footer>
      <NSpace :size="16">
        <NButton type="primary" @click="closeKeyModal">{{ $t('page.user.apiKey.keyModal.confirm') }}</NButton>
      </NSpace>
    </template>
  </NModal>
</template>

<style scoped></style>