<script setup lang="ts">
import { requestStatusOptions } from '@/constants/business';
import { useNaiveForm } from '@/hooks/common/form';
import { translateOptions } from '@/utils/common';
import { $t } from '@/locales';
import type { GetRequestListRequest } from '@/typings/proto/admin/v1/relay';
import { relayServiceClient } from '@/grpc';
import { ref, onMounted } from 'vue';

defineOptions({
  name: 'RequestSearch'
});

interface Emits {
  (e: 'reset'): void;
  (e: 'search'): void;
}

const emit = defineEmits<Emits>();

const { validate, restoreValidation } = useNaiveForm();

const model = defineModel<GetRequestListRequest>('model', { required: true });

/** the provider options */
interface AccountOption extends CommonType.Option<string> {};

const accountOptions = ref<AccountOption[]>([]);

async function getAccountOptions() {
  const resp = await relayServiceClient.getAccountList({ current: 1, size: 1000 });
  accountOptions.value = resp.records.map(item => ({
    label: item.name,
    value: item.id.toString()
  }));
} 

/** the provider options */
interface ProviderCodeOption extends CommonType.Option<string> {}

const providerCodeOptions = ref<ProviderCodeOption[]>([]);

async function getProviderCodeOptions() {
  const resp = await relayServiceClient.getProviderCodeList({});
  const options = resp.records.map(item => ({
    label: item,
    value: item,
  }));

  providerCodeOptions.value = [...options];
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
  getProviderCodeOptions();
  getAccountOptions();
});
</script>

<template>
  <NCard :bordered="false" size="small" class="card-wrapper">
    <NCollapse>
      <NCollapseItem :title="$t('common.search')" name="user-search">
        <NForm ref="formRef" :model="model" label-placement="left" :label-width="80">
          <NGrid responsive="screen" item-responsive>
            <NFormItemGi span="24 s:12 m:6" :label="$t('page.usage.request.providerCode')" path="providerCode" class="pr-24px">
              <NSelect
                v-model:value="model.providerCode"
                :placeholder="$t('page.usage.request.form.providerCode')"
                :options="providerCodeOptions"
                filterable
                clearable
              />
            </NFormItemGi>
            <NFormItemGi span="24 s:12 m:6" :label="$t('page.usage.request.accountId')" path="accountId" class="pr-24px">
              <NSelect
                v-model:value="model.accountId"
                :placeholder="$t('page.usage.request.form.accountId')"
                :options="accountOptions"
                filterable
                clearable
              />
            </NFormItemGi>
            <NFormItemGi span="24 s:12 m:6" :label="$t('page.usage.request.status')" path="status" class="pr-24px">
              <NSelect
                v-model:value="model.status"
                :placeholder="$t('page.usage.request.form.status')"
                :options="translateOptions(requestStatusOptions)"
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
