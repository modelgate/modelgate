<script setup lang="tsx">
import { NButton, NPopconfirm, NTag } from 'naive-ui';
import { ledgerTypeRecord } from '@/constants/business';
import { useAppStore } from '@/store/modules/app';
import { useTable, useTableOperate } from '@/hooks/common/table';
import type { NaiveUI } from '@/typings/naive-ui';
import { relayServiceClient } from '@/grpc';
import { $t } from '@/locales';
import { Ledger } from '@/typings/proto/model/relay/ledger_pb';
import LedgerSearch from './modules/ledger-search.vue';
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
  apiFn: relayServiceClient.getLedgerList,
  showTotal: true,
  apiParams: {
    current: 1,
    size: 15,
    orderBy: '-id',
    // if you want to use the searchParams in Form, you need to define the following properties, and the value is null
    // the value can not be undefined, otherwise the property in Form will not be reactive
    accountId: null,
    type: '',
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
      key: 'accountName',
      title: $t('page.usage.ledger.accountName'),
      align: 'center',
      width: 100
    },
    {
      key: 'type',
      title: $t('page.usage.ledger.type'),
      align: 'center',
      width: 80,
      render: (row: Ledger) => {
        if (row.type === null) {
          return null;
        }

        const tagMap: Record<string, NaiveUI.ThemeColor> = {
          'consume': 'primary',
          'refund': 'warning',
          'charge': 'success',
          'adjust': 'error',
        };

        const label = $t(ledgerTypeRecord[row.type]);

        return <NTag type={tagMap[row.type]}>{label}</NTag>;
      }
    },
    {
      key: 'amount',
      title: $t('page.usage.ledger.amount'),
      align: 'right',
      width: 100,
      render: (row: Ledger) => row.amount.toLocaleString()
    },
    {
      key: 'balanceAfter',
      title: $t('page.usage.ledger.balanceAfter'),
      align: 'right',
      width: 100,
      render: (row: Ledger) => row.balanceAfter.toLocaleString()
    },
    {
      key: 'requestId',
      title: $t('page.usage.ledger.requestId'),
      align: 'center',
      width: 100,
    },
    {
      key: 'createdAt',
      title: $t('page.usage.ledger.createdAt'),
      align: 'center',
      width: 100,
      render: (row: Ledger) => formatProtoTime(row.createdAt),
    },
    {
      key: 'operate',
      title: $t('common.operate'),
      align: 'center',
      width: 130,
      render: (row: Ledger) => {
        const canDelete = hasAuth('usage:ledger:delete')

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
  useTableOperate<Ledger>(data, getData);

async function handleBatchDelete() {
  try {
    await relayServiceClient.deleteLedgers({ ids: checkedRowKeys.value });
  } catch {
    window.$message?.error($t('common.deleteFailed'));
  }

  onBatchDeleted();
}

async function handleDelete(id: bigint) {
  try {
    await relayServiceClient.deleteLedgers({ ids: [id] });
  } catch {
    window.$message?.error($t('common.deleteFailed'));
  }

  onDeleted();
}


</script>

<template>
  <div class="min-h-500px flex-col-stretch gap-16px overflow-hidden lt-sm:overflow-auto">
    <LedgerSearch v-model:model="searchParams" @reset="resetSearchParams" @search="getDataByPage" />
    <NCard :title="$t('page.usage.ledger.title')" :bordered="false" size="small" class="sm:flex-1-hidden card-wrapper">
      <template #header-extra>
        <TableHeaderOperation
          v-model:columns="columnChecks"
          :show-add="false"
          :disabled-delete="checkedRowKeys.length === 0"
          :show-batch-delete="hasAuth('usage:ledger:delete')"
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
        :row-key="(row: Ledger) => row.id.toString()"
        :pagination="mobilePagination"
        @update:sorter="handleSorterChange"
        class="sm:h-full"
      />
    </NCard>
  </div>
</template>

<style scoped></style>
