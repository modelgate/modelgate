<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { apiKeyStatusOptions, enableStatusOptions } from '@/constants/business';
import { useFormRules, useNaiveForm } from '@/hooks/common/form';
import { relayServiceClient } from '@/grpc';
import { $t } from '@/locales';
import type { Account } from '@/typings/proto/model/relay/accout_pb';
import { NaiveUI } from '@/typings/naive-ui';

defineOptions({
  name: 'AccountOperateModal'
});

interface Props {
  /** the type of operation */
  operateType: NaiveUI.TableOperateType;
  /** the edit row data */
  rowData?: Account | null;
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
    add: $t('page.user.account.addAccount'),
    edit: $t('page.user.account.editAccount')
  };
  return titles[props.operateType];
});

type Model = Omit<Pick<Account, 'id' | 'name' | 'nickname' | 'balance' | 'status'>, 'balance'> & {
  balance: number;
};

const model = ref(createDefaultModel());

function createDefaultModel(): Model {
  return {
    id: 0n,
    name: '',
    nickname: '',
    balance: 0,
    status: '',
  };
}

type RuleKey = Extract<keyof Model, 'name' | 'nickname' | 'status'>;

const rules: Record<RuleKey, App.Global.FormRule> = {
  name: defaultRequiredRule,
  nickname: defaultRequiredRule,
  status: defaultRequiredRule
};

function handleInitModel() {
  model.value = createDefaultModel();

  if (props.operateType === 'edit' && props.rowData) {
    Object.assign(model.value, {
      ...props.rowData,
      balance: Number(props.rowData.balance)
    });
  }
}

function closeDrawer() {
  visible.value = false;
}

async function handleSubmit() {
  await validate();

  if (props.operateType === 'edit') {
    try {
      await relayServiceClient.updateAccount({
        updateMask: {
          paths: ['name', 'nickname', 'balance', 'status']
        },
        account: {
          ...model.value,
          balance: BigInt(model.value.balance)
        }
      });
      window.$message?.success($t('common.updateSuccess'));
    } catch {
      window.$message?.error($t('common.updateFailed'));
    }
  } else {
    try {
      await relayServiceClient.createAccount({
        account: { 
          ...model.value,
          balance: BigInt(model.value.balance),
         }
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
    <NScrollbar class="h-420px pr-20px">
      <NForm ref="formRef" :model="model" :rules="rules">
        <NFormItem :label="$t('page.user.account.name')" path="name">
          <NInput v-model:value="model.name" :placeholder="$t('page.user.account.form.name')" />
        </NFormItem>
        <NFormItem :label="$t('page.user.account.nickname')" path="nickname">
          <NInput v-model:value="model.nickname" :placeholder="$t('page.user.account.form.nickname')" />
        </NFormItem>
        <NFormItem :label="$t('page.user.account.balance')" path="balance">
           <NInputNumber v-model:value="model.balance" :placeholder="$t('page.user.account.form.balance')" class="w-full" :min="0" />
        </NFormItem>
        <NFormItem :label="$t('page.user.account.status')" path="status">
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
