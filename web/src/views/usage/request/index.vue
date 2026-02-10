<script setup lang="tsx">
import { NButton, NPopconfirm, NTag } from 'naive-ui';
import { requestStatusRecord } from '@/constants/business';
import { useAppStore } from '@/store/modules/app';
import { useTable, useTableOperate } from '@/hooks/common/table';
import type { NaiveUI } from '@/typings/naive-ui';
import { relayServiceClient } from '@/grpc';
import { $t } from '@/locales';
import { Request } from '@/typings/proto/model/relay/request_pb';
import RequestSearch from './modules/request-search.vue';
import { formatProtoTime } from '@/utils/common';
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
  apiFn: relayServiceClient.getRequestList,
  showTotal: true,
  apiParams: {
    current: 1,
    size: 15,
    orderBy: '-id',
    // if you want to use the searchParams in Form, you need to define the following properties, and the value is null
    // the value can not be undefined, otherwise the property in Form will not be reactive
    accountId: null,
    providerCode: '',
    modelCode: '',
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
      key: 'requestUuid',
      title: $t('page.usage.request.requestUuid'),
      align: 'center',
      width: 180
    },
    {
      key: 'createdAt',
      title: $t('page.usage.request.createdAt'),
      align: 'center',
      width: 100,
      render: (row: Request) => formatProtoTime(row.createdAt),
    },
    {
      key: 'accountName',
      title: $t('page.usage.request.accountName'),
      align: 'left',
      width: 100
    },
    {
      key: 'providerCode',
      title: $t('page.usage.request.providerCode'),
      align: 'left',
      width: 100
    },
    {
      key: 'actualModel',
      title: $t('page.usage.request.modelCode'),
      align: 'left',
      width: 100,
    },
    {
      key: 'promptTokens',
      title: $t('page.usage.request.promptTokens'),
      align: 'right',
      width: 80,
    },
    {
      key: 'completionTokens',
      title: $t('page.usage.request.completionTokens'),
      align: 'right',
      width: 80,
    },
    {
      key: 'totalTokens',
      title: $t('page.usage.request.totalTokens'),
      align: 'right',
      width: 80,
    },
    {
      key: 'status',
      title: $t('page.usage.request.status'),
      align: 'center',
      width: 100,
      sorter: true,
      render: (row: Request) => {
        if (row.status === null) {
          return null;
        }

        const tagMap: Record<string, NaiveUI.ThemeColor> = {
          'pending': 'warning',
          'success': 'success',
          'failed': 'error',
          'cancelled': 'warning',
        };

        const label = $t(requestStatusRecord[row.status]);

        return <NTag type={tagMap[row.status]}>{label}</NTag>;
      }
    },
    {
      key: 'elapsedTime',
      title: $t('page.usage.request.elapsedTime'),
      align: 'right',
      width: 50,
      render: (row: Request) => `${Number(row.elapsedTime) / 1000}`,
    },
    {
      key: 'operate',
      title: $t('common.operate'),
      align: 'center',
      width: 130,
      render: (row: Request) => {
        const canDelete = hasAuth('usage:request:delete')

        if (!canDelete) {
          return null
        }
        return (
        <div class="flex-center gap-8px">
          <NPopconfirm onPositiveClick={() => handleDelete(row.id)}>
            {{
              default: () => $t('common.confirmDelete'),
              trigger: () => (
                <NButton type="error" ghost size="small">
                  {$t('common.delete')}
                </NButton>
              )
            }}
          </NPopconfirm>
        </div>
        )
      }
    }
  ]
});

const { drawerVisible, operateType, editingData, checkedRowKeys, onBatchDeleted, onDeleted } =
  useTableOperate<Request>(data, getData);

async function handleBatchDelete() {
  try {
    await relayServiceClient.deleteRequests({ ids: checkedRowKeys.value });
  } catch {
    window.$message?.error($t('common.deleteFailed'));
  }

  onBatchDeleted();
}

async function handleDelete(id: bigint) {
  try {
    await relayServiceClient.deleteRequests({ ids: [id] });
  } catch {
    window.$message?.error($t('common.deleteFailed'));
  }

  onDeleted();
}

</script>

<template>
  <div class="min-h-500px flex-col-stretch gap-16px overflow-hidden lt-sm:overflow-auto">
    <RequestSearch v-model:model="searchParams" @reset="resetSearchParams" @search="getDataByPage" />
    <NCard :title="$t('page.usage.request.title')" :bordered="false" size="small" class="sm:flex-1-hidden card-wrapper">
      <template #header-extra>
        <TableHeaderOperation
          v-model:columns="columnChecks"
          :show-add="false"
          :disabled-delete="checkedRowKeys.length === 0"
          :show-batch-delete="hasAuth('usage:request:delete')"
          :loading="loading"
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
        :row-key="(row: Request) => row.id.toString()"
        :pagination="mobilePagination"
        @update:sorter="handleSorterChange"
        class="sm:h-full"
      />
    </NCard>
  </div>
</template>

<style scoped></style>
