<script setup lang="ts">
import {  modelStatusOptions } from '@/constants/business';
import {  useNaiveForm } from '@/hooks/common/form';
import { translateOptions } from '@/utils/common';
import { $t } from '@/locales';
import type { GetModelListRequest } from '@/typings/proto/admin/v1/relay';
import { relayServiceClient } from '@/grpc';
import { ref, onMounted } from 'vue';

defineOptions({
  name: 'ModelSearch'
});

interface Emits {
  (e: 'reset'): void;
  (e: 'search'): void;
}

const emit = defineEmits<Emits>();

const { validate, restoreValidation } = useNaiveForm();

const model = defineModel<GetModelListRequest>('model', { required: true });

/** the provider options */
interface ProviderOption extends CommonType.Option<string> {};

const providerOptions = ref<ProviderOption[]>([]);

 async function getProviderOptions() {
  const resp = await relayServiceClient.getProviderList({ current: 1, size: 1000});
  providerOptions.value = resp.records.map(item => ({
    label: item.name,
    value: item.code
  }));
} 

async function reset() {
  await restoreValidation();
  emit('reset');
  emit('search');
}

async function search() {
  await validate();
  emit('search');
}

onMounted(() => {
  getProviderOptions();
});

</script>

<template>
  <NCard :bordered="false" size="small" class="card-wrapper">
    <NCollapse>
      <NCollapseItem :title="$t('common.search')" name="user-search">
        <NForm ref="formRef" :model="model" label-placement="left" :label-width="80">
          <NGrid responsive="screen" item-responsive>
            <NFormItemGi span="24 s:12 m:6" :label="$t('page.relay.model.form.providerCode')" path="providerCode" class="pr-24px">
              <NSelect
                v-model:value="model.providerCode"
                :placeholder="$t('page.relay.model.form.providerCode')"
                :options="providerOptions"
                filterable
                clearable
              />
            </NFormItemGi>
            <NFormItemGi span="24 s:12 m:6" :label="$t('page.relay.model.name')" path="name" class="pr-24px">
              <NInput v-model:value="model.name" :placeholder="$t('page.relay.model.form.name')" />
            </NFormItemGi>
            <NFormItemGi span="24 s:12 m:6" :label="$t('page.relay.model.code')" path="code" class="pr-24px">
              <NInput v-model:value="model.code" :placeholder="$t('page.relay.model.form.code')" />
            </NFormItemGi>
            <NFormItemGi span="24 s:12 m:6" :label="$t('page.relay.model.status')" path="status" class="pr-24px">
              <NSelect
                v-model:value="model.status"
                :placeholder="$t('page.relay.model.form.status')"
                :options="translateOptions(modelStatusOptions)"
                clearable
              />
            </NFormItemGi>
            <NFormItemGi span="24 m:12" class="pr-24px">
              <NSpace class="w-full" justify="end">
                <NButton @click="reset">
                  <template #icon>
                    <icon-ic-round-refresh class="text-icon" />
                  </template>
                  {{ $t('common.reset') }}
                </NButton>
                <NButton type="primary" ghost @click="search">
                  <template #icon>
                    <icon-ic-round-search class="text-icon" />
                  </template>
                  {{ $t('common.search') }}
                </NButton>
              </NSpace>
            </NFormItemGi>
          </NGrid>
        </NForm>
      </NCollapseItem>
    </NCollapse>
  </NCard>
</template>

<style scoped></style>
