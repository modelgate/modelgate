<script setup lang="tsx">
import { NButton, NPopconfirm, NTag } from 'naive-ui';
import { apiKeyStatusRecord } from '@/constants/business';
import { useAppStore } from '@/store/modules/app';
import { useTable, useTableOperate } from '@/hooks/common/table';
import type { NaiveUI } from '@/typings/naive-ui';
import { relayServiceClient } from '@/grpc';
import { $t } from '@/locales';
import AccountApiKeySearch from './modules/account-api-key-search.vue';
import AccountApiKeyOperateModal from './modules/account-api-key-operate-modal.vue';
import { formatProtoTime } from '@/utils/common';
import { AccountApiKey } from '@/typings/proto/model/relay/account_api_key_pb';
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
  apiFn: relayServiceClient.getAccountApiKeyList,
  showTotal: true,
  apiParams: {
    current: 1,
    size: 15,
    orderBy: '-id',
    accountId: 0n,
    keyword: '',
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
      key: 'accountName',
      title: $t('page.user.apiKey.accountName'),
      align: 'left',
      width: 100,
    },
    {
      key: 'keyName',
      title: $t('page.user.apiKey.keyName'),
      align: 'left',
      width: 100,
    },
    {
      key: 'key',
      title: $t('page.user.apiKey.key'),
      align: 'left',
      width: 200,
    },
    {
      key: 'scope',
      title: $t('page.user.apiKey.scope'),
      align: 'left',
      width: 100
    },
    {
      key: 'quoteLimit',
      title: $t('page.user.apiKey.quoteLimit'),
      align: 'right',
      width: 80,
      render: (row: AccountApiKey) => row.quoteLimit.toLocaleString(),
    },
    {
      key: 'quoteUsed',
      title: $t('page.user.apiKey.quoteUsed'),
      align: 'right',
      width: 80,
      render: (row: AccountApiKey) => row.quoteUsed.toLocaleString(),
    },
    {
      key: 'rateLimit',
      title: $t('page.user.apiKey.rateLimit'),
      align: 'center',
      width: 80,
      render: (row: AccountApiKey) => row.rateLimit.toLocaleString(),
    },
    {
      key: 'lastUsedAt',
      title: $t('page.user.apiKey.lastUsedAt'),
      align: 'center',
      width: 100,
      render: (row: AccountApiKey) => formatProtoTime(row.lastUsedAt),
    },
    {
      key: 'expiredAt',
      title: $t('page.user.apiKey.expiredAt'),
      align: 'center',
      width: 100,
      render: (row: AccountApiKey) => formatProtoTime(row.expiredAt),
    },
    {
      key: 'status',
      title: $t('page.user.apiKey.status'),
      align: 'center',
      width: 80,
      sorter: true,
      render: (row: AccountApiKey) => {
        if (!row.status) {
          return null;
        }

        const tagMap: Record<string, NaiveUI.ThemeColor> = {
          'enabled': 'success',
          'disabled': 'warning',
          'revoked': 'error'
        };

        const label = $t(apiKeyStatusRecord[row.status]);

        return <NTag type={tagMap[row.status]}>{label}</NTag>;
      }
    },
    {
      key: 'operate',
      title: $t('common.operate'),
      align: 'center',
      width: 130,
      render: (row: AccountApiKey) => {
        const canEdit = hasAuth('relay:account_api_key:edit')
        const canDelete = hasAuth('relay:account_api_key:delete')

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
  useTableOperate<AccountApiKey>(data, getData);

async function handleBatchDelete() {
  try {
    await relayServiceClient.deleteAccountApiKeys({ ids: checkedRowKeys.value });
  } catch {
    window.$message?.error($t('common.deleteFailed'));
  }

  onBatchDeleted();
}

async function handleDelete(id: bigint) {
  try {
    await relayServiceClient.deleteAccountApiKeys({ ids: [id] });
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
    <AccountApiKeySearch v-model:model="searchParams" @reset="resetSearchParams" @search="getDataByPage" />
    <NCard :title="$t('page.user.apiKey.title')" :bordered="false" size="small" class="sm:flex-1-hidden card-wrapper">
      <template #header-extra>
        <TableHeaderOperation
          v-model:columns="columnChecks"
          :disabled-delete="checkedRowKeys.length === 0"
          :loading="loading"
          :show-add="hasAuth('relay:account_api_key:edit')"
          :show-batch-delete="hasAuth('relay:account_api_key:delete')"
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
        :scroll-x="1400"
        :loading="loading"
        remote
        @update:sorter="handleSorterChange"
        :row-key="(row: AccountApiKey) => row.id.toString()"
        :pagination="mobilePagination"
        class="sm:h-full"
      />
      <AccountApiKeyOperateModal
        v-model:visible="drawerVisible"
        :operate-type="operateType"
        :row-data="editingData"
        @submitted="handleSubmitted"
      />
    </NCard>
  </div>
</template>

<style scoped></style>
