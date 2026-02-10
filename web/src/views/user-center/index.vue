<script setup lang="ts">
import { ref } from 'vue';
import { userGenderOptions } from '@/constants/business';
import { useFormRules, useNaiveForm } from '@/hooks/common/form';
import { $t } from '@/locales';
import type { User } from '@/typings/proto/model/system/user_pb';
import { authServiceClient, systemServiceClient } from '@/grpc';

defineOptions({
  name: 'UserCenter'
});

const { formRef, validate } = useNaiveForm();
const { defaultRequiredRule } = useFormRules();

type Model = Pick<User, 'id' | 'username' | 'gender' | 'nickname' | 'phone' | 'email' | 'password'>;

const model = ref(createDefaultModel());

function createDefaultModel(): Model {
  return {
    id: BigInt(0),
    username: '',
    gender: '',
    nickname: '',
    phone: '',
    email: '',
    password: ''
  };
}

async function handleInitModel() {
  model.value = createDefaultModel();
  const { user } = await authServiceClient.getUserInfo({});

  Object.assign(model.value, user);
}

handleInitModel();

type RuleKey = Extract<keyof Model, 'username' | 'nickname' | 'email' | 'phone' | 'gender'>;

const rules: Record<RuleKey, App.Global.FormRule> = {
  username: defaultRequiredRule,
  nickname: defaultRequiredRule,
  email: defaultRequiredRule,
  phone: defaultRequiredRule,
  gender: defaultRequiredRule,
};

async function handleSubmit() {
  await validate();
  try {
    await systemServiceClient.updateUser({
      updateMask: {
        paths: ['nickname', 'phone', 'email', 'gender', 'password']
      },
      user: {
        ...model.value,
      }
    });
    window.$message?.success($t('common.updateSuccess'));
  } catch {
    window.$message?.error($t('common.updateFailed'));
  }
  handleInitModel();
}
</script>

<template>
  <NSpace vertical :size="16">
    <NCard :title="$t('page.userCenter.title')" :bordered="false" size="small" segmented class="card-wrapper">
      <NForm
        ref="formRef"
        :model="model"
        class="w-1/3"
        :rules="rules"
        size="large"
        :show-label="true"
        label-placement="left"
        label-align="left"
        label-width="80"
        @keyup.enter="handleSubmit"
      >
        <NFormItem :label="$t('page.manage.user.username')" path="username">
          <NInput v-model:value="model.username" :placeholder="$t('page.manage.user.username')" disabled :maxlength="11" />
        </NFormItem>
        <NFormItem :label="$t('page.manage.user.nickname')" path="nickname">
          <NInput v-model:value="model.nickname" :placeholder="$t('page.manage.user.nickname')" :maxlength="11" />
        </NFormItem>
        <NFormItem :label="$t('page.manage.user.email')" path="email">
          <NInput v-model:value="model.email" :placeholder="$t('page.manage.user.email')" :maxlength="11" />
        </NFormItem>
        <NFormItem :label="$t('page.manage.user.phone')" path="phone">
          <NInput v-model:value="model.phone" :placeholder="$t('page.manage.user.phone')" :maxlength="11" />
        </NFormItem>
        <NFormItem :label="$t('page.manage.user.gender')" path="gender">
          <NRadioGroup v-model:value="model.gender">
            <NRadio v-for="item in userGenderOptions" :key="item.value" :value="item.value" :label="$t(item.label)" />
          </NRadioGroup>
        </NFormItem>
        <NFormItem :label="$t('page.manage.user.password')" path="password">
          <NInput
            v-model:value="model.password"
            type="password"
            :placeholder="$t('page.manage.user.password')"
            :maxlength="11"
          />
        </NFormItem>

        <NSpace vertical :size="18" class="w-full">
          <NButton type="primary" size="large" round block :loading="submiting" @click="handleSubmit">
            {{ $t('common.confirm') }}
          </NButton>
        </NSpace>
      </NForm>
    </NCard>
  </NSpace>
</template>

<style scoped></style>
