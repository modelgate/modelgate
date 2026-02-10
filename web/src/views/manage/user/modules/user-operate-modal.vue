<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { enableStatusOptions, genderOptions } from '@/constants/business';
import { useFormRules, useNaiveForm } from '@/hooks/common/form';
import { systemServiceClient } from '@/grpc';
import { $t } from '@/locales';
import type { User } from '@/typings/proto/model/system/user_pb';

defineOptions({
  name: 'UserOperateModal'
});

interface Props {
  /** the type of operation */
  operateType: NaiveUI.TableOperateType;
  /** the edit row data */
  rowData?: User | null;
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
    add: $t('page.manage.user.addUser'),
    edit: $t('page.manage.user.editUser')
  };
  return titles[props.operateType];
});

type Model = Pick<User, 'id' | 'username' | 'gender' | 'nickname' | 'phone' | 'email' | 'roles' | 'status' | 'password'>;

const model = ref(createDefaultModel());
const usernameInputRef = ref();
const passwordInputRef = ref();

function createDefaultModel(): Model {
  return {
    id: 0n,
    username: '',
    gender: '',
    nickname: '',
    phone: '',
    email: '',
    roles: [],
    status: '',
    password: ''
  };
}

type RuleKey = Extract<keyof Model, 'username' | 'password' | 'status'>;

const rules: Record<RuleKey, App.Global.FormRule> = {
  username: defaultRequiredRule,
  password: {
    required: false,
    trigger: 'blur',
    validator: (_rule, value) => {
      if (props.operateType === 'add' && !value) {
        return new Error($t('form.pwd.required'));
      }
      return true;
    }
  },
  status: defaultRequiredRule
};

/** the enabled role options */
const roleOptions = ref<CommonType.Option<string>[]>([]);

async function getRoleOptions() {
  const resp = await systemServiceClient.getRoleList({ current: 1, size: 10000, status: 'enabled' });
  const options = resp.records.map(item => ({
    label: item.name,
    value: item.id.toString()
  }));

  roleOptions.value = [...options];
}

function handleInitModel() {
  model.value = createDefaultModel();

  if (props.operateType === 'edit' && props.rowData) {
    Object.assign(model.value, props.rowData);
    model.value.password = '';
  }

  // Force clear input values in add mode to prevent browser autofill
  if (props.operateType === 'add') {
    setTimeout(() => {
      if (usernameInputRef.value) {
        usernameInputRef.value.value = '';
      }
      if (passwordInputRef.value) {
        passwordInputRef.value.value = '';
      }
    }, 0);
  }
}

function closeDrawer() {
  visible.value = false;
}

async function handleSubmit() {
  await validate();

  if (props.operateType === 'edit') {
    const paths = ['username', 'gender', 'nickname', 'phone', 'email', 'status', 'roles'];
    const submitData = { ...model.value };

    // Only include password if it's provided
    if (model.value.password) {
      paths.push('password');
    }

    try {
      await systemServiceClient.updateUser({
        updateMask: {
          paths
        },
        user: submitData
      });
      window.$message?.success($t('common.updateSuccess'));
    } catch {
      window.$message?.error($t('common.updateFailed'));
    }
  } else {
    try {
      await systemServiceClient.createUser({
        user: { ...model.value }
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
    getRoleOptions();
  }
});
</script>

<template>
  <NModal v-model:show="visible" :title="title" preset="card" class="w-600px">
    <NScrollbar class="h-480px pr-20px">
      <NForm ref="formRef" :model="model" :rules="rules">
        <NFormItem :label="$t('page.manage.user.nickname')" path="nickname">
          <NInput v-model:value="model.nickname" :placeholder="$t('page.manage.user.form.nickname')" />
        </NFormItem>
        <NFormItem :label="$t('page.manage.user.username')" path="username">
          <NInput
            ref="usernameInputRef"
            v-model:value="model.username"
            :placeholder="$t('page.manage.user.form.username')"
            :disabled="operateType === 'edit'"
            :readonly="operateType === 'add'"
            @focus="() => { if (operateType === 'add') { (usernameInputRef as any)?.$el?.querySelector('input')?.removeAttribute('readonly'); } }"
          />
        </NFormItem>
        <NFormItem :label="$t('page.manage.user.password')" path="password">
          <NInput
            ref="passwordInputRef"
            v-model:value="model.password"
            type="password"
            show-password-on="click"
            :placeholder="operateType === 'add' ? $t('page.manage.user.form.password') : $t('page.manage.user.form.password') + ' (' + $t('common.optional') + ')'"
            :readonly="operateType === 'add'"
            @focus="() => { if (operateType === 'add') { (passwordInputRef as any)?.$el?.querySelector('input')?.removeAttribute('readonly'); } }"
          />
        </NFormItem>
        <NFormItem :label="$t('page.manage.user.gender')" path="gender">
          <NRadioGroup v-model:value="model.gender">
            <NRadio v-for="item in genderOptions" :key="item.value" :value="item.value" :label="$t(item.label)" />
          </NRadioGroup>
        </NFormItem>
        <NFormItem :label="$t('page.manage.user.phone')" path="phone">
          <NInput v-model:value="model.phone" :placeholder="$t('page.manage.user.form.phone')" />
        </NFormItem>
        <NFormItem :label="$t('page.manage.user.email')" path="email">
          <NInput v-model:value="model.email" :placeholder="$t('page.manage.user.form.email')" />
        </NFormItem>
        <NFormItem :label="$t('page.manage.user.role')" path="roles">
          <NSelect
            v-model:value="model.roles"
            multiple
            filterable
            :options="roleOptions"
            :placeholder="$t('page.manage.user.form.role')"
          />
        </NFormItem>
        <NFormItem :label="$t('page.manage.user.status')" path="status">
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
