<script setup lang="tsx">
import { NButton, NPopconfirm, NTag } from 'naive-ui';
import { enableStatusRecord } from '@/constants/business';
import { useAppStore } from '@/store/modules/app';
import { useTable, useTableOperate } from '@/hooks/common/table';
import type { NaiveUI } from '@/typings/naive-ui';
import { relayServiceClient } from '@/grpc';
import { $t } from '@/locales';
import ProviderApiKeyOperateModal from './modules/provider-api-key-operate-modal.vue';
import ProviderApiKeySearch from './modules/provider-api-key-search.vue';
import {formatProtoTime} from '@/utils/common';
import { ProviderApiKey } from '@/typings/proto/model/relay/provider_api_key_pb';
import { useAuth } from '@/hooks/business/auth';

const appStore = useAppStore();
const { hasAuth } = useAuth();

const {
  columns,
  columnChecks,
  data,
  getData,
  getDataByPage,
  loading,
  mobilePagination,
  searchParams,
  resetSearchParams,
  handleSorterChange
} = useTable({
  apiFn: relayServiceClient.getProviderApiKeyList,
  showTotal: true,
  apiParams: {
    current: 1,
    size: 15,
    orderBy: '-id',
    // if you want to use the searchParams in Form, you need to define the following properties, and the value is null
    // the value can not be undefined, otherwise the property in Form will not be reactive
    // id: '',
    name: '',
    providerCode : '',
    status: '',
  },
  columns: () => [
    {
      type: 'selection',
      align: 'center',
      width: 48
    },
    {
      key: 'id',
      title: $t('common.index'),
      align: 'center',
      width: 64,
      sorter: true,
    },
    {
      key: 'name',
      title: $t('page.relay.providerApiKey.name'),
      align: 'left',
      width: 100,
    },
    {
      key: 'providerCode',
      title: $t('page.relay.providerApiKey.providerCode'),
      align: 'left',
      width: 100
    },
    {
      key: 'key',
      title: $t('page.relay.providerApiKey.key'),
      align: 'left',
      width: 100
    },
    {
      key: 'weight',
      title: $t('page.relay.providerApiKey.weight'),
      align: 'center',
      width: 50,
    },
    {
      key: 'lastUsedAt',
      title: $t('page.relay.providerApiKey.lastUsedAt'),
      align: 'center',
      width: 100,
      render: (row: ProviderApiKey) =>  formatProtoTime(row.lastUsedAt),
    },
    {
      key: 'status',
      title: $t('page.relay.providerApiKey.status'),
      align: 'center',
      width: 50,
      sorter: true,
      render: (row: ProviderApiKey) => {
        if (row.status === null) {
          return null;
        }

        const tagMap: Record<string, NaiveUI.ThemeColor> = {
          'enabled': 'success',
          'disabled': 'warning'
        };

        const label = $t(enableStatusRecord[row.status]);

        return <NTag type={tagMap[row.status]}>{label}</NTag>;
      }
    },
    {
      key: 'operate',
      title: $t('common.operate'),
      align: 'center',
      width: 130,
      render: (row: ProviderApiKey) => {
        const canEdit = hasAuth('relay:provider_api_key:edit')
        const canDelete = hasAuth('relay:provider_api_key:delete')

        if (!canEdit && !canDelete) {
          return null
        }
        return (
        <div class="flex-center gap-8px">
          {canEdit &&
            (<NButton type="primary" ghost size="small" onClick={() => edit(row.id)}>
              {$t('common.edit')}
            </NButton>)}
          {canDelete &&
            (<NPopconfirm onPositiveClick={() => handleDelete(row.id)}>
              {{
              default: () => $t('common.confirmDelete'),
              trigger: () => (
                <NButton type="error" ghost size="small">
                  {$t('common.delete')}
                </NButton>
              )
            }}
          </NPopconfirm>)}
        </div>
        )
      }
    }
  ]
});

const { drawerVisible, operateType, editingData, handleAdd, handleEdit, checkedRowKeys, onBatchDeleted, onDeleted } =
  useTableOperate<ProviderApiKey>(data, getData);

async function handleBatchDelete() {
  try {
    await relayServiceClient.deleteProviderApiKeys({ ids: checkedRowKeys.value });
  } catch {
    window.$message?.error($t('common.deleteFailed'));
  }

  onBatchDeleted();
}

async function handleDelete(id: bigint) {
  try {
    await relayServiceClient.deleteProviderApiKeys({ ids: [id] });
  } catch {
    window.$message?.error($t('common.deleteFailed'));
  }

  onDeleted();
}

function handleSubmitted() {
  drawerVisible.value = false;
  getData();
}

function edit(id: bigint) {
  handleEdit(id);
}
</script>

<template>
  <div class="min-h-500px flex-col-stretch gap-16px overflow-hidden lt-sm:overflow-auto">
    <ProviderApiKeySearch v-model:model="searchParams" @reset="resetSearchParams" @search="getDataByPage" />
    <NCard :title="$t('page.relay.providerApiKey.title')" :bordered="false" size="small" class="sm:flex-1-hidden card-wrapper">
      <template #header-extra>
        <TableHeaderOperation
          v-model:columns="columnChecks"
          :disabled-delete="checkedRowKeys.length === 0"
          :show-batch-delete="hasAuth('relay:provider_api_key:delete')"
          :loading="loading"
          :show-add="hasAuth('relay:provider_api_key:edit')"
          @add="handleAdd"
          @delete="handleBatchDelete"
          @refresh="getData"
        />
      </template>
      <NDataTable
        v-model:checked-row-keys="checkedRowKeys"
        :columns="columns"
        :data="data"
        size="small"
        :flex-height="!appStore.isMobile"
        :scroll-x="962"
        :loading="loading"
        remote
        @update:sorter="handleSorterChange"
        :row-key="(row: ProviderApiKey) => row.id.toString()"
        :pagination="mobilePagination"
        class="sm:h-full"
      />
      <ProviderApiKeyOperateModal
        v-model:visible="drawerVisible"
        :operate-type="operateType"
        :row-data="editingData"
        @submitted="handleSubmitted"
      />
    </NCard>
  </div>
</template>

<style scoped></style>
