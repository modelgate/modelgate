<script setup lang="tsx">
import { NButton, NPopconfirm, NTag } from 'naive-ui';
import { enableStatusRecord } from '@/constants/business';
import { useAppStore } from '@/store/modules/app';
import { useTable, useTableOperate } from '@/hooks/common/table';
import type { NaiveUI } from '@/typings/naive-ui';
import { relayServiceClient } from '@/grpc';
import { $t } from '@/locales';
import { Provider } from '@/typings/proto/model/relay/provider_pb';
import ProviderOperateModal from './modules/provider-operate-modal.vue';
import ProviderSearch from './modules/provider-search.vue';
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
  apiFn: relayServiceClient.getProviderList,
  showTotal: true,
  apiParams: {
    current: 1,
    size: 15,
    orderBy: '-id',
    // if you want to use the searchParams in Form, you need to define the following properties, and the value is null
    // the value can not be undefined, otherwise the property in Form will not be reactive
    name: '',
    code: '',
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
      title: $t('page.relay.provider.name'),
      align: 'left',
      width: 150
    },
    {
      key: 'code',
      title: $t('page.relay.provider.code'),
      align: 'left',
      width: 150
    },
    {
      key: 'baseUrl',
      title: $t('page.relay.provider.baseUrl'),
      align: 'left',
      width: 200,
    },
    {
      key: 'status',
      title: $t('page.relay.provider.status'),
      align: 'center',
      width: 100,
      sorter: true,
      render: (row: Provider) => {
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
      render: (row: Provider) => {
        const canEdit = hasAuth('relay:provider:edit')
        const canDelete = hasAuth('relay:provider:delete')

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
  useTableOperate<Provider>(data, getData);

async function handleBatchDelete() {
  try {
    await relayServiceClient.deleteProviders({ ids: checkedRowKeys.value });
  } catch {
    window.$message?.error($t('common.deleteFailed'));
  }

  onBatchDeleted();
}

async function handleDelete(id: bigint) {
  try {
    await relayServiceClient.deleteProviders({ ids: [id] });
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
    <ProviderSearch v-model:model="searchParams" @reset="resetSearchParams" @search="getDataByPage" />
    <NCard :title="$t('page.relay.provider.title')" :bordered="false" size="small" class="sm:flex-1-hidden card-wrapper">
      <template #header-extra>
        <TableHeaderOperation
          v-model:columns="columnChecks"
          :disabled-delete="checkedRowKeys.length === 0"
          :show-batch-delete="hasAuth('relay:provider:delete')"
          :loading="loading"
          :show-add="hasAuth('relay:provider:edit')"
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
        :row-key="(row: Provider) => row.id.toString()"
        :pagination="mobilePagination"
        @update:sorter="handleSorterChange"
        class="sm:h-full"
      />
      <ProviderOperateModal
        v-model:visible="drawerVisible"
        :operate-type="operateType"
        :row-data="editingData"
        @submitted="handleSubmitted"
      />
    </NCard>
  </div>
</template>

<style scoped></style>
