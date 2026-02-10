<script setup lang="ts">
import { ledgerTypeOptions } from '@/constants/business';
import { useNaiveForm } from '@/hooks/common/form';
import { translateOptions } from '@/utils/common';
import { $t } from '@/locales';
import type { GetAccountListRequest } from '@/typings/proto/admin/v1/relay';
import { relayServiceClient } from '@/grpc';
import { ref, onMounted } from 'vue';

defineOptions({
  name: 'LedgerSearch'
});

interface Emits {
  (e: 'reset'): void;
  (e: 'search'): void;
}

const emit = defineEmits<Emits>();

const { validate, restoreValidation } = useNaiveForm();

const model = defineModel<GetAccountListRequest>('model', { required: true });

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
  getAccountOptions();
});
</script>

<template>
  <NCard :bordered="false" size="small" class="card-wrapper">
    <NCollapse>
      <NCollapseItem :title="$t('common.search')" name="user-search">
        <NForm ref="formRef" :model="model" label-placement="left" :label-width="80">
          <NGrid responsive="screen" item-responsive>
            <NFormItemGi span="24 s:12 m:6" :label="$t('page.usage.ledger.accountId')" path="accountId" class="pr-24px">
              <NSelect
                v-model:value="model.accountId"
                :placeholder="$t('page.usage.ledger.form.accountId')"
                :options="accountOptions"
                filterable
                clearable
              />
            </NFormItemGi>

            <NFormItemGi span="24 s:12 m:6" :label="$t('page.usage.ledger.type')" path="type" class="pr-24px">
              <NSelect
                v-model:value="model.type"
                :placeholder="$t('page.usage.ledger.form.type')"
                :options="translateOptions(ledgerTypeOptions)"
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
