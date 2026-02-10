<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { modelStatusOptions } from '@/constants/business';
import { useFormRules, useNaiveForm } from '@/hooks/common/form';
import { relayServiceClient } from '@/grpc';
import { $t } from '@/locales';
import type { Model } from '@/typings/proto/model/relay/model_pb';
import { NaiveUI } from '@/typings/naive-ui';

defineOptions({
  name: 'ModelOperateModal'
});

interface Props {
  /** the type of operation */
  operateType: NaiveUI.TableOperateType;
  /** the edit row data */
  rowData?: Model | null;
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
    add: $t('page.relay.model.addModel'),
    edit: $t('page.relay.model.editModel')
  };
  return titles[props.operateType];
});

// Use number for IDs/Integers in the form to work with Naive UI components
interface ModelForm {
  id: number;
  providerId: number | null;
  providerCode: string;
  name: string;
  code: string;
  actualCode: string;
  priority: number;
  weight: number;
  status: string;
}

const model = ref(createDefaultModel());

function createDefaultModel(): ModelForm {
  return {
    id: 0,
    providerId: null,
    providerCode: '',
    name: '',
    code: '',
    actualCode: '',
    priority: 1,
    weight: 100,
    status: 'enabled'
  };
}

type RuleKey = Extract<keyof ModelForm, 'providerId' | 'name' | 'code' | 'status'>;

const rules: Record<RuleKey, App.Global.FormRule> = {
  providerId: defaultRequiredRule,
  name: defaultRequiredRule,
  code: defaultRequiredRule,
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

  if (props.operateType === 'edit' && props.rowData) {
    const row = props.rowData;
    model.value = {
      id: Number(row.id),
      providerId: Number(row.providerId),
      providerCode: row.providerCode,
      name: row.name,
      code: row.code,
      actualCode: row.actualCode,
      priority: Number(row.priority),
      weight: Number(row.weight),
      status: row.status
    };
  }
}

function handleProviderUpdate(value: number, option: ProviderOption) {
  model.value.providerId = value;
  model.value.providerCode = option.code;
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
    priority: BigInt(model.value.priority),
    weight: BigInt(model.value.weight),
    providerCode: model.value.providerCode
  };

  if (props.operateType === 'edit') {
    try {
      await relayServiceClient.updateModel({
        updateMask: {
          paths: ['provider_id', 'provider_code', 'name', 'code', 'actual_code','priority', 'weight', 'status']
        },
        model: submissionData as any // Cast to any or Model to bypass exact type match issues
      });
      window.$message?.success($t('common.updateSuccess'));
    } catch {
      window.$message?.error($t('common.updateFailed'));
    }
  } else {
    try {
      await relayServiceClient.createModel({
        model: submissionData as any
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
    getProviderOptions();
  }
});
</script>

<template>
  <NModal v-model:show="visible" :title="title" preset="card" class="w-600px">
    <NScrollbar class="h-520px pr-20px">
      <NForm ref="formRef" :model="model" :rules="rules">
        <NFormItem :label="$t('page.relay.model.providerId')" path="providerId">
           <NSelect
            v-model:value="model.providerId"
            :options="providerOptions"
            :placeholder="$t('page.relay.model.form.providerId')"
            @update:value="handleProviderUpdate"
            filterable
          />
        </NFormItem>
        <NFormItem :label="$t('page.relay.model.name')" path="name">
          <NInput v-model:value="model.name" :placeholder="$t('page.relay.model.form.name')" />
        </NFormItem>
        <NFormItem :label="$t('page.relay.model.code')" path="code">
          <NInput v-model:value="model.code" :placeholder="$t('page.relay.model.form.code')" />
        </NFormItem>
        <NFormItem :label="$t('page.relay.model.actualCode')" path="actualCode">
          <NInput v-model:value="model.actualCode" :placeholder="$t('page.relay.model.form.actualCode')" />
        </NFormItem>
        <NFormItem :label="$t('page.relay.model.priority')" path="priority">
          <NInputNumber v-model:value="model.priority" :placeholder="$t('page.relay.model.form.priority')" class="w-full" />
        </NFormItem>
        <NFormItem :label="$t('page.relay.model.weight')" path="weight">
           <NInputNumber v-model:value="model.weight" :placeholder="$t('page.relay.model.form.weight')" class="w-full" />
        </NFormItem>
        <NFormItem :label="$t('page.relay.model.status')" path="status">
          <NRadioGroup v-model:value="model.status">
            <NRadio v-for="item in modelStatusOptions" :key="item.value" :value="item.value" :label="$t(item.label)" />
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
