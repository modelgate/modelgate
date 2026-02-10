<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { enableStatusOptions } from '@/constants/business';
import { useFormRules, useNaiveForm } from '@/hooks/common/form';
import { relayServiceClient } from '@/grpc';
import { $t } from '@/locales';
import type { Provider } from '@/typings/proto/model/relay/provider_pb';
import { NaiveUI } from '@/typings/naive-ui';

defineOptions({
  name: 'ProviderOperateModal'
});

interface Props {
  /** the type of operation */
  operateType: NaiveUI.TableOperateType;
  /** the edit row data */
  rowData?: Provider | null;
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
    add: $t('page.relay.provider.addProvider'),
    edit: $t('page.relay.provider.editProvider')
  };
  return titles[props.operateType];
});

type Model = Pick<Provider, 'id' | 'name' | 'code' | 'baseUrl'|'status'>;

const model = ref(createDefaultModel());

function createDefaultModel(): Model {
  return {
    id: 0n,
    name: '',
    code: '',
    baseUrl: '',
    status: '',
  };
}

type RuleKey = Extract<keyof Model, 'name' | 'code' | 'baseUrl' | 'status'>;

const rules: Record<RuleKey, App.Global.FormRule> = {
  name: defaultRequiredRule,
  code: defaultRequiredRule,
  baseUrl: defaultRequiredRule,
  status: defaultRequiredRule
};

/** the provider options */
interface ProviderCodeOption extends CommonType.Option<string> {

}
const providerCodeOptions = ref<ProviderCodeOption[]>([]);

async function getProviderCodeOptions() {
  const resp = await relayServiceClient.getProviderCodeList({});
  const options = resp.records.map(item => ({
    label: item,
    value: item,
  }));

  providerCodeOptions.value = [...options];
}

function handleInitModel() {
  model.value = createDefaultModel();

  if (props.operateType === 'edit' && props.rowData) {
    Object.assign(model.value, props.rowData);
  }
}

function closeDrawer() {
  visible.value = false;
}

async function handleSubmit() {
  await validate();

  if (props.operateType === 'edit') {
    try {
      await relayServiceClient.updateProvider({
        updateMask: {
          paths: ['name', 'code', 'base_url', 'status']
        },
        provider: { ...model.value }
      });
      window.$message?.success($t('common.updateSuccess'));
    } catch {
      window.$message?.error($t('common.updateFailed'));
    }
  } else {
    try {
      await relayServiceClient.createProvider({
        provider: { ...model.value }
      });
      window.$message?.success($t('common.addSuccess'));
    } catch {
      window.$message?.error($t('common.addFailed'));
    }
  }

  // closeDrawer();
  emit('submitted');
}

watch(visible, () => {
  if (visible.value) {
    handleInitModel();
    restoreValidation();
    getProviderCodeOptions();
  }
});
</script>

<template>
  <NModal v-model:show="visible" :title="title" preset="card" class="w-600px">
    <NScrollbar class="h-350px pr-20px">
      <NForm ref="formRef" :model="model" :rules="rules">
        <NFormItem :label="$t('page.relay.provider.name')" path="name">
          <NInput v-model:value="model.name" :placeholder="$t('page.relay.provider.form.name')" />
        </NFormItem>
        <NFormItem :label="$t('page.relay.provider.code')" path="code">
          <NSelect v-model:value="model.code" :options="providerCodeOptions" />
        </NFormItem>
        <NFormItem :label="$t('page.relay.provider.baseUrl')" path="baseUrl">
          <NInput v-model:value="model.baseUrl" :placeholder="$t('page.relay.provider.form.baseUrl')" />
        </NFormItem>
        <NFormItem :label="$t('page.relay.provider.status')" path="status">
          <NRadioGroup v-model:value="model.status">
            <NRadio v-for="item in enableStatusOptions" :key="item.value" :value="item.value" :label="$t(item.label)" />
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
