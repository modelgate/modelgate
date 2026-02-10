<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { currencyOptions, enableStatusOptions } from '@/constants/business';
import { useFormRules, useNaiveForm } from '@/hooks/common/form';
import { relayServiceClient } from '@/grpc';
import { $t } from '@/locales';
import type { ModelPricing } from '@/typings/proto/model/relay/model_pricing_pb';
import { NaiveUI } from '@/typings/naive-ui';
import { msToProto, protoToMs } from '@/utils/common';


defineOptions({
  name: 'ModelPricingOperateModal'
});

interface Props {
  /** the type of operation */
  operateType: NaiveUI.TableOperateType;
  /** the edit row data */
  rowData?: ModelPricing | null;
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
    add: $t('page.relay.modelPricing.addModelPricing'),
    edit: $t('page.relay.modelPricing.editModelPricing')
  };
  return titles[props.operateType];
});

interface ModelPricingForm {
  id: number;
  modelCode: string;
  currency: string;
  pointsPerCurrency: number;
  tokenNum: number;
  inputPrice: number;
  inputCachePrice: number;
  outputPrice: number;
  effectiveFrom: number | null;
  effectiveTo: number | null;
  status: string;
}

const model = ref(createDefaultModel());

function createDefaultModel(): ModelPricingForm {
  return {
    id: 0,
    modelCode: '',
    currency: 'CNY',
    pointsPerCurrency: 1000000, // CNY default
    tokenNum: 1000000,
    inputPrice: 0,
    inputCachePrice: 0,
    outputPrice: 0,
    effectiveFrom: null,
    effectiveTo: null,
    status: 'enabled'
  };
}

function getDefaultPointsPerCurrency(currency: string): number {
  return currency === 'USD' ? 7200000 : (currency==='CNY' ? 1000000 : 1);
}

type RuleKey = Extract<keyof ModelPricingForm, 'modelCode' | 'currency' | 'effectiveFrom' | 'effectiveTo' | 'status'>;

const rules: Record<RuleKey, App.Global.FormRule> = {
  modelCode: defaultRequiredRule,
  currency: defaultRequiredRule,
  effectiveFrom: defaultRequiredRule,
  effectiveTo: defaultRequiredRule,
  status: defaultRequiredRule
};

/** the model options */
interface ModelOption extends CommonType.Option<string> {
  providerCode: string;
}
const modelOptions = ref<ModelOption[]>([]);

async function getModelOptions() {
  const resp = await relayServiceClient.getModelList({ current: 1, size: 1000, status: 'enabled', orderBy:"-updated_at" });
  const options = resp.records.map(item => ({
    label: item.name +" ( "+item.providerCode +" / "+(item.actualCode!=''?item.actualCode: item.code)+" )",
    value: item.actualCode!=''?item.actualCode: item.code,
    providerCode: item.providerCode,
  }));

  modelOptions.value = [...options];
}
function handleInitModel() {
  model.value = createDefaultModel();

  if (props.operateType === 'edit' && props.rowData) {
    const row = props.rowData;
    model.value = {
      id: Number(row.id),
      modelCode: row.modelCode,
      currency: row.currency,
      pointsPerCurrency: Number(row.pointsPerCurrency),
      tokenNum: Number(row.tokenNum),
      inputPrice: row.inputPrice,
      inputCachePrice: row.inputCachePrice,
      outputPrice: row.outputPrice,
      effectiveFrom: protoToMs(row.effectiveFrom),
      effectiveTo: protoToMs(row.effectiveTo),
      status: row.status
    };
  }
  getModelOptions();
}

function handleCurrencyUpdate(value: string) {
  model.value.currency = value;
  // Only auto-update pointsPerCurrency in add mode
  if (props.operateType === 'add') {
    model.value.pointsPerCurrency = getDefaultPointsPerCurrency(value);
  }
}

function closeDrawer() {
  visible.value = false;
}

async function handleSubmit() {
  await validate();

  const selectedModel = modelOptions.value.find(opt => opt.value === model.value.modelCode);
  const providerCode = selectedModel?.providerCode || '';

  const submissionData = {
    ...model.value,
    id: BigInt(model.value.id),
    pointsPerCurrency: BigInt(model.value.pointsPerCurrency),
    tokenNum: BigInt(model.value.tokenNum),
    providerCode,
    effectiveFrom: msToProto(model.value.effectiveFrom),
    effectiveTo: msToProto(model.value.effectiveTo),
   
  };

  if (props.operateType === 'edit') {
    try {
      await relayServiceClient.updateModelPricing({
        updateMask: {
          paths: ['provider_code', 'model_code', 'currency', 'points_per_currency', 'token_num', 'input_price', 'input_cache_price', 'output_price', 'effective_from', 'effective_to', 'status']
        },
        modelPricing: submissionData as any
      });
      window.$message?.success($t('common.updateSuccess'));
    } catch {
      window.$message?.error($t('common.updateFailed'));
    }
  } else {
    try {
      await relayServiceClient.createModelPricing({
        modelPricing: submissionData as any
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
  }
});
</script>

<template>
  <NModal v-model:show="visible" :title="title" preset="card" class="w-600px">
    <NScrollbar class="h-500px pr-20px">
      <NForm ref="formRef" :model="model" :rules="rules">
        <NFormItem :label="$t('page.relay.modelPricing.modelCode')" path="modelCode">
          <NSelect
            v-model:value="model.modelCode"
            :options="modelOptions"
            :placeholder="$t('page.relay.modelPricing.form.modelCode')"
            filterable
            clearable
          />
        </NFormItem>
        <NGrid :cols="2" :x-gap="16">
          <NFormItemGi :label="$t('page.relay.modelPricing.currency')" path="currency">
            <NRadioGroup v-model:value="model.currency" @update:value="handleCurrencyUpdate">
              <NRadio v-for="item in currencyOptions" :key="item.value" :value="item.value" :label="$t(item.label)" />
            </NRadioGroup>
          </NFormItemGi>
          <NFormItemGi :label="$t('page.relay.modelPricing.pointsPerCurrency')" path="pointsPerCurrency">
            <NInputNumber v-model:value="model.pointsPerCurrency" :placeholder="$t('page.relay.modelPricing.form.pointsPerCurrency')" class="w-full" :min="0" />
          </NFormItemGi>
        </NGrid>
        <NGrid :cols="2" :x-gap="16">
          <NFormItemGi :label="$t('page.relay.modelPricing.tokenNum')" path="tokenNum">
            <NInputNumber v-model:value="model.tokenNum" :placeholder="$t('page.relay.modelPricing.form.tokenNum')" class="w-full" :min="0"/>
          </NFormItemGi>
        </NGrid>
        <NGrid :cols="2" :x-gap="16">
          <NFormItemGi :label="$t('page.relay.modelPricing.inputPrice')" path="inputPrice">
            <NInputNumber v-model:value="model.inputPrice" :placeholder="$t('page.relay.modelPricing.form.inputPrice')" class="w-full" :min="0" :precision="4"/>
          </NFormItemGi>
          <NFormItemGi :label="$t('page.relay.modelPricing.inputCachePrice')" path="inputCachePrice">
            <NInputNumber v-model:value="model.inputCachePrice" :placeholder="$t('page.relay.modelPricing.form.inputCachePrice')" class="w-full" :min="0" :precision="4"/>
          </NFormItemGi>
        </NGrid>
        <NGrid :cols="2" :x-gap="16">
          <NFormItemGi :label="$t('page.relay.modelPricing.outputPrice')" path="outputPrice">
            <NInputNumber v-model:value="model.outputPrice" :placeholder="$t('page.relay.modelPricing.form.outputPrice')" class="w-full" :min="0" :precision="4"/>
          </NFormItemGi>
        </NGrid>
        <NGrid :cols="2" :x-gap="16">
          <NFormItemGi :label="$t('page.relay.modelPricing.effectiveFrom')" path="effectiveFrom">
             <NDatePicker v-model:value="model.effectiveFrom" type="datetime" clearable class="w-full"/>
          </NFormItemGi>
          <NFormItemGi :label="$t('page.relay.modelPricing.effectiveTo')" path="effectiveTo">
             <NDatePicker v-model:value="model.effectiveTo" type="datetime" clearable class="w-full"/>
          </NFormItemGi>
        </NGrid>
        <NFormItem :label="$t('page.relay.modelPricing.status')" path="status">
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
