<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { apiKeyStatusOptions } from '@/constants/business';
import { useFormRules, useNaiveForm } from '@/hooks/common/form';
import { relayServiceClient } from '@/grpc';
import { $t } from '@/locales';
import type { ProviderApiKey } from '@/typings/proto/model/relay/provider_api_key_pb';
import { NaiveUI } from '@/typings/naive-ui';

defineOptions({
  name: 'ProviderApiKeyOperateModal'
});

interface Props {
  /** the type of operation */
  operateType: NaiveUI.TableOperateType;
  /** the edit row data */
  rowData?: ProviderApiKey | null;
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
    add: $t('page.relay.providerApiKey.addProviderApiKey'),
    edit: $t('page.relay.providerApiKey.editProviderApiKey')
  };
  return titles[props.operateType];
});

interface ProviderApiKeyForm {
  id: number;
  providerId: number | null;
  name: string;
  key: string;
  weight: number;
  status: string;
}

const model = ref(createDefaultModel());
const hasUserModifiedKey = ref(false);

function createDefaultModel(): ProviderApiKeyForm {
  return {
    id: 0,
    providerId: null,
    name: '',
    key: '',
    weight: 100,
    status: 'enabled'
  };
}

type RuleKey = Extract<keyof ProviderApiKeyForm, 'providerId' | 'name' | 'weight' | 'status'>;

const rules: Record<RuleKey, App.Global.FormRule> = {
  providerId: defaultRequiredRule,
  name: defaultRequiredRule,
  weight: defaultRequiredRule,
  status: defaultRequiredRule
};

/** the provider options */
interface ProviderOption extends CommonType.Option<number> {
  code: string;
}
const providerOptions = ref<ProviderOption[]>([]);

async function getProviderOptions() {
  const resp = await relayServiceClient.getProviderList({ current: 1, size: 1000, status: 'enabled' });
  const options = resp.records.map(item => ({
    label: item.name,
    value: Number(item.id), // Convert bigint to number
    code: item.code
  }));

  providerOptions.value = [...options];
}

function handleInitModel() {
  model.value = createDefaultModel();
  hasUserModifiedKey.value = false;

  if (props.operateType === 'edit' && props.rowData) {
    const row = props.rowData;
    model.value = {
      id: Number(row.id),
      providerId: Number(row.providerId),
      name: row.name,
      key: row.key, // server already returns masked key
      weight: Number(row.weight),
      status: row.status
    };
  }
}

function handleKeyFocus() {
  // Clear the masked value on first focus
  if (!hasUserModifiedKey.value && props.operateType === 'edit') {
    model.value.key = '';
    hasUserModifiedKey.value = true;
  }
}

function handleProviderUpdate(value: number, option: ProviderOption) {
  model.value.providerId = value;
}

function closeDrawer() {
  visible.value = false;
}

async function handleSubmit() {
  await validate();

  const submissionData = {
    ...model.value,
    id: BigInt(model.value.id),
    providerId: BigInt(model.value.providerId!),
    weight: BigInt(model.value.weight),
  };

  if (props.operateType === 'edit') {
    // All fields are always submitted, except key which is only submitted if user modified it
    const paths = ['provider_id', 'provider_code', 'name', 'weight', 'status'];

    // Only include key if user has clicked and entered a value
    if (hasUserModifiedKey.value && model.value.key) {
      paths.push('key');
    }

    try {
      await relayServiceClient.updateProviderApiKey({
        updateMask: {
          paths
        },
        providerApiKey: submissionData as any
      });
      window.$message?.success($t('common.updateSuccess'));
    } catch {
      window.$message?.error($t('common.updateFailed'));
    }
  } else {
    try {
      await relayServiceClient.createProviderApiKey({
        providerApiKey: submissionData as any
      });
      window.$message?.success($t('common.addSuccess'));
    } catch {
      window.$message?.error($t('common.addFailed'));
    }
  }

  emit('submitted');
}

watch(visible, () => {
  if (visible.value) {
    handleInitModel();
    restoreValidation();
    getProviderOptions();
  }
});
</script>

<template>
  <NModal v-model:show="visible" :title="title" preset="card" class="w-600px">
    <NScrollbar class="h-420px pr-20px">
      <NForm ref="formRef" :model="model" :rules="rules">
        <NFormItem :label="$t('page.relay.providerApiKey.providerId')" path="providerId">
           <NSelect
            v-model:value="model.providerId"
            :options="providerOptions"
            :placeholder="$t('page.relay.providerApiKey.form.providerId')"
            @update:value="handleProviderUpdate"
            filterable
          />
        </NFormItem>
        <NFormItem :label="$t('page.relay.providerApiKey.name')" path="name">
          <NInput v-model:value="model.name" :placeholder="$t('page.relay.providerApiKey.form.name')" />
        </NFormItem>
        <NFormItem :label="$t('page.relay.providerApiKey.key')" path="key">
          <NInput
            v-model:value="model.key"
            :placeholder="$t('page.relay.providerApiKey.form.key')"
            @focus="handleKeyFocus"
          />
        </NFormItem>
        <NFormItem :label="$t('page.relay.providerApiKey.weight')" path="weight">
           <NInputNumber v-model:value="model.weight" :placeholder="$t('page.relay.providerApiKey.form.weight')" class="w-full" :min="0" />
        </NFormItem>
        <NFormItem :label="$t('page.relay.providerApiKey.status')" path="status">
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
</template>

<style scoped></style>
