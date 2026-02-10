<script setup lang="tsx">
import { NButton, NPopconfirm } from 'naive-ui';
import { useAppStore } from '@/store/modules/app';
import { useTable, useTableOperate } from '@/hooks/common/table';
import { $t } from '@/locales';
import { systemServiceClient } from '@/grpc';
import type { Permission } from '@/typings/proto/model/system/permission_pb';
import PermissionOperateModal from './modules/permission-operate-modal.vue';
import PermissionSearch from './modules/permission-search.vue';
import { useAuth } from '@/hooks/business/auth';

const appStore = useAppStore();
const { hasAuth } = useAuth();

const {
  columns,
  columnChecks,
  data,
  loading,
  getData,
  getDataByPage,
  mobilePagination,
  searchParams,
  resetSearchParams,
  handleSorterChange
} = useTable({
  apiFn: systemServiceClient.getPermissionList,
  apiParams: {
    current: 1,
    size: 15,
    orderBy: '-id',
    // if you want to use the searchParams in Form, you need to define the following properties, and the value is null
    // the value can not be undefined, otherwise the property in Form will not be reactive
    name: '',
    code: '',
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
      width: 64
    },
    {
      key: 'name',
      title: $t('page.manage.permission.name'),
      align: 'left',
      width: 100
    },
    {
      key: 'code',
      title: $t('page.manage.permission.code'),
      align: 'left',
      width: 100,
      sorter: true
    },
    {
      key: 'data',
      title: $t('page.manage.permission.data'),
      align: 'left',
      width: 240,
      render: (row: Permission) => row.data?.map(item => `${item.path}`).join('、') || '-'
    },
    {
      key: 'desc',
      title: $t('page.manage.permission.desc'),
      align: 'left',
      width: 80
    },
    {
      key: 'operate',
      title: $t('common.operate'),
      align: 'center',
      width: 130,
      render: (row: Permission) => {
        const canEdit = hasAuth('system:permission:edit')
        const canDelete = hasAuth('system:permission:delete')

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

const {
  drawerVisible,
  operateType,
  editingData,
  handleAdd,
  handleEdit,
  checkedRowKeys,
  onBatchDeleted,
  onDeleted
} = useTableOperate<Permission>(data, getData);

async function handleBatchDelete() {
  try {
    await systemServiceClient.deletePermissions({ ids: checkedRowKeys.value });
    onBatchDeleted();
  } catch {
    window.$message?.error($t('common.deleteFailed'));
  }
}

async function handleDelete(id: bigint) {
  try {
    await systemServiceClient.deletePermissions({ ids: [id] });
    onDeleted();
  } catch {
    window.$message?.error($t('common.deleteFailed'));
  }
}

function edit(id: bigint) {
  handleEdit(id);
}

function handleSubmitted() {
  drawerVisible.value = false;
  getData();
}
</script>

<template>
  <div class="min-h-500px flex-col-stretch gap-16px overflow-hidden lt-sm:overflow-auto">
    <PermissionSearch v-model:model="searchParams" @reset="resetSearchParams" @search="getDataByPage" />
    <NCard :title="$t('page.manage.permission.title')" :bordered="false" size="small" class="sm:flex-1-hidden card-wrapper">
      <template #header-extra>
        <TableHeaderOperation
          v-model:columns="columnChecks"
          :disabled-delete="checkedRowKeys.length === 0"
          :show-batch-delete="hasAuth('system:permission:delete')"
          :loading="loading"
          :show-add="hasAuth('system:permission:edit')"
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
        :scroll-x="782"
        :loading="loading"
        remote
        :row-key="(row: Permission) => row.id.toString()"
        :pagination="mobilePagination"
        class="sm:h-full"
        @update:sorter="handleSorterChange"
      />
      <PermissionOperateModal
        v-model:visible="drawerVisible"
        :operate-type="operateType"
        :row-data="editingData"
        @submitted="handleSubmitted"
      />
    </NCard>
  </div>
</template>

<style scoped></style>
