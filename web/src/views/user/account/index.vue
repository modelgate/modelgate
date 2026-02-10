<script setup lang="tsx">
import { useAppStore } from '@/store/modules/app';
import { NButton, NPopconfirm, NTag } from 'naive-ui';
import { useTable, useTableOperate } from '@/hooks/common/table';
import { enableStatusRecord } from '@/constants/business';
import { relayServiceClient } from '@/grpc';
import { $t } from '@/locales';
import { Account } from '@/typings/proto/model/relay/accout_pb';
import { formatProtoTime } from '@/utils/common';
import { NaiveUI } from '@/typings/naive-ui';
import AccountOperateModal from './modules/account-operate-modal.vue';
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
  apiFn: relayServiceClient.getAccountList,
  showTotal: true,
  apiParams: {
    current: 1,
    size: 15,
    orderBy: '-id',
    // if you want to use the searchParams in Form, you need to define the following properties, and the value is null
    // the value can not be undefined, otherwise the property in Form will not be reactive
    accountId: 0n,
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
      title: $t('page.user.account.name'),
      align: 'left',
      width: 100,
    },
    {
      key: 'nickname',
      title: $t('page.user.account.nickname'),
      align: 'left',
      width: 100,
    },
    {
      key: 'balance',
      title: $t('page.user.account.balance'),
      align: 'right',
      width: 100,
      render: (row: Account) => row.balance.toLocaleString(),
    },
    {
      key: 'status',
      title: $t('page.user.account.status'),
      align: 'center',
      width: 80,
      sorter: true,
      render: (row: Account) => {
        if (!row.status) {
          return null;
        }

        const tagMap: Record<string, NaiveUI.ThemeColor> = {
          'enabled': 'success',
          'disabled': 'warning',
        };

        const label = $t(enableStatusRecord[row.status]);

        return <NTag type={tagMap[row.status]}>{label}</NTag>;
      }
    },
    {
      key: 'createdAt',
      title: $t('page.user.account.createdAt'),
      align: 'center',
      width: 100,
      render: (row: Account) => formatProtoTime(row.createdAt),
    },
    {
      key: 'updatedAt',
      title: $t('page.user.account.updatedAt'),
      align: 'center',
      width: 100,
      render: (row: Account) => formatProtoTime(row.updatedAt),
    },
    {
      key: 'operate',
      title: $t('common.operate'),
      align: 'center',
      width: 130,
      render: (row: Account) => {
        const canEdit = hasAuth('user:account:edit')
        const canDelete = hasAuth('user:account:delete')

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
  useTableOperate<Account>(data, getData);

async function handleBatchDelete() {
  try {
    await relayServiceClient.deleteAccounts({ ids: checkedRowKeys.value });
  } catch {
    window.$message?.error($t('common.deleteFailed'));
  }

  onBatchDeleted();
}

async function handleDelete(id: bigint) {
  try {
    await relayServiceClient.deleteAccounts({ ids: [id] });
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
    <NCard :title="$t('page.user.account.title')" :bordered="false" size="small" class="sm:flex-1-hidden card-wrapper">
      <template #header-extra>
        <TableHeaderOperation
          v-model:columns="columnChecks"
          :disabled-delete="checkedRowKeys.length === 0"
          :loading="loading"
          :show-add="hasAuth('user:account:edit')"
          :show-batch-delete="hasAuth('user:account:delete')"
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
        :row-key="(row: Account) => row.id.toString()"
        :pagination="mobilePagination"
        @update:sorter="handleSorterChange"
        class="sm:h-full"
      />
      <AccountOperateModal
        v-model:visible="drawerVisible"
        :operate-type="operateType"
        :row-data="editingData"
        @submitted="handleSubmitted"
      />
    </NCard>
  </div>
</template>

<style scoped></style>
