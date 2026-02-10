<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { useBoolean } from '@sa/hooks';
import { enableStatusOptions, yesOrNoOptions } from '@/constants/business';
import { useFormRules, useNaiveForm } from '@/hooks/common/form';
import { $t } from '@/locales';
import type { Role } from '@/typings/proto/model/system/role_pb';
import { systemServiceClient } from '@/grpc';
import ButtonAuthModal from './button-auth-modal.vue';
import MenuAuthModal from './menu-auth-modal.vue';
import { NaiveUI } from '@/typings/naive-ui';

defineOptions({
  name: 'RoleOperateModal'
});

interface Props {
  /** the type of operation */
  operateType: NaiveUI.TableOperateType;
  /** the edit row data */
  rowData?: Role | null;
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
const { bool: menuAuthVisible, setTrue: openMenuAuthModal } = useBoolean();
const { bool: buttonAuthVisible, setTrue: openButtonAuthModal } = useBoolean();

const title = computed(() => {
  const titles: Record<NaiveUI.TableOperateType, string> = {
    add: $t('page.manage.role.addRole'),
    edit: $t('page.manage.role.editRole')
  };
  return titles[props.operateType];
});

type Model = Pick<Role, 'id' | 'name' | 'code' | 'description' | 'status' | 'isSuperAdmin'>;

const model = ref(createDefaultModel());

function createDefaultModel(): Model {
  return {
    id: 0n,
    name: '',
    code: '',
    isSuperAdmin: false,
    description: '',
    status: ''
  };
}

type RuleKey = Extract<keyof Model, 'name' | 'code' | 'isSuperAdmin' | 'status'>;

const rules: Record<RuleKey, App.Global.FormRule> = {
  name: defaultRequiredRule,
  code: defaultRequiredRule,
  isSuperAdmin: defaultRequiredRule,
  status: defaultRequiredRule
};

const roleId = computed(() => props.rowData?.id.toString() || "-1");

const isEdit = computed(() => props.operateType === 'edit');

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
  // request
  if (props.operateType === 'edit') {
    try {
      await systemServiceClient.updateRole({
        updateMask: {
          paths: ['name', 'code', 'is_super_admin', 'description', 'status']
        },
        role: { ...model.value}
      });
      window.$message?.success($t('common.updateSuccess'));
    } catch {
      window.$message?.error($t('common.updateFailed'));
    }
  } else {
    try {
      await systemServiceClient.createRole({
        role: { ...model.value}
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
    <NScrollbar class="h-360px pr-20px">
      <NForm ref="formRef" :model="model" :rules="rules">
        <NFormItem :label="$t('page.manage.role.name')" path="name">
          <NInput v-model:value="model.name" :placeholder="$t('page.manage.role.form.name')" />
        </NFormItem>
        <NFormItem :label="$t('page.manage.role.code')" path="code">
          <NInput v-model:value="model.code" :placeholder="$t('page.manage.role.form.code')" />
        </NFormItem>
        <NFormItem :label="$t('page.manage.role.isSuperAdmin')" path="isSuperAdmin">
          <NRadioGroup v-model:value="model.isSuperAdmin">
            <NRadio :value="true" :label="$t('common.yesOrNo.yes')" />
            <NRadio :value="false" :label="$t('common.yesOrNo.no')" />
          </NRadioGroup>
        </NFormItem>
        <NFormItem :label="$t('page.manage.role.status')" path="status">
          <NRadioGroup v-model:value="model.status">
            <NRadio v-for="item in enableStatusOptions" :key="item.value" :value="item.value" :label="$t(item.label)" />
          </NRadioGroup>
        </NFormItem>
        <NFormItem :label="$t('page.manage.role.description')" path="description">
          <NInput v-model:value="model.description" :placeholder="$t('page.manage.role.form.description')" />
        </NFormItem>
      </NForm>
      <NSpace v-if="isEdit">
        <NButton @click="openMenuAuthModal">{{ $t('page.manage.role.menuAuth') }}</NButton>
        <MenuAuthModal v-model:visible="menuAuthVisible" :role-id="roleId" />
        <NButton @click="openButtonAuthModal">{{ $t('page.manage.role.buttonAuth') }}</NButton>
        <ButtonAuthModal v-model:visible="buttonAuthVisible" :role-id="roleId" />
      </NSpace>
      <template #footer>
        <NSpace :size="16">
          <NButton @click="closeDrawer">{{ $t('common.cancel') }}</NButton>
          <NButton type="primary" @click="handleSubmit">{{ $t('common.confirm') }}</NButton>
        </NSpace>
      </template>
    </NScrollbar>
    <template #footer>
      <NSpace justify="end" :size="16">
        <NButton @click="closeDrawer">{{ $t('common.cancel') }}</NButton>
        <NButton type="primary" @click="handleSubmit">{{ $t('common.confirm') }}</NButton>
      </NSpace>
    </template>
  </NModal>
</template>

<style scoped></style>
