<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { enableStatusOptions } from '@/constants/business';
import { useNaiveForm } from '@/hooks/common/form';
import { translateOptions } from '@/utils/common';
import { $t } from '@/locales';
import type { GetProviderListRequest } from '@/typings/proto/admin/v1/relay';
import { relayServiceClient } from '@/grpc';

defineOptions({
  name: 'ProviderSearch'
});

interface Emits {
  (e: 'reset'): void;
  (e: 'search'): void;
}

const emit = defineEmits<Emits>();

const { validate, restoreValidation } = useNaiveForm();

const model = defineModel<GetProviderListRequest>('model', { required: true });

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
});
</script>

<template>
  <NCard :bordered="false" size="small" class="card-wrapper">
    <NCollapse>
      <NCollapseItem :title="$t('common.search')" name="user-search">
        <NForm ref="formRef" :model="model" label-placement="left" :label-width="80">
          <NGrid responsive="screen" item-responsive>
            <NFormItemGi span="24 s:12 m:6" :label="$t('page.relay.provider.name')" path="name" class="pr-24px">
              <NInput v-model:value="model.name" :placeholder="$t('page.relay.provider.form.name')" />
            </NFormItemGi>
            <NFormItemGi span="24 s:12 m:6" :label="$t('page.relay.provider.code')" path="code" class="pr-24px">
              <NSelect
                v-model:value="model.code"
                :placeholder="$t('page.relay.provider.form.code')"
                :options="providerCodeOptions"
                clearable
              />
            </NFormItemGi>
            <NFormItemGi span="24 s:12 m:6" :label="$t('page.relay.provider.status')" path="status" class="pr-24px">
              <NSelect
                v-model:value="model.status"
                :placeholder="$t('page.relay.provider.form.status')"
                :options="translateOptions(enableStatusOptions)"
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
