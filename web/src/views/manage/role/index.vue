<script setup lang="tsx">
import { NButton, NPopconfirm, NTag } from 'naive-ui';
import { enableStatusRecord } from '@/constants/business';
import { useAppStore } from '@/store/modules/app';
import { useTable, useTableOperate } from '@/hooks/common/table';
import { $t } from '@/locales';
import { systemServiceClient } from '@/grpc';
import type { Role } from '@/typings/proto/model/system/role_pb';
import type { NaiveUI } from '@/typings/naive-ui';
import RoleOperateModal from './modules/role-operate-modal.vue';
import RoleSearch from './modules/role-search.vue';
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
  resetSearchParams
} = useTable({
  apiFn: systemServiceClient.getRoleList,
  apiParams: {
    current: 1,
    size: 15,
    orderBy: '-id',
    // if you want to use the searchParams in Form, you need to define the following properties, and the value is null
    // the value can not be undefined, otherwise the property in Form will not be reactive
    status: '',
    name: '',
    code: ''
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
      width: 64,
      align: 'center'
    },
    {
      key: 'name',
      title: $t('page.manage.role.name'),
      align: 'left',
      width: 100
    },
    {
      key: 'code',
      title: $t('page.manage.role.code'),
      align: 'left',
      width: 100
    },
    {
      key: 'description',
      title: $t('page.manage.role.description'),
      align: 'left',
      width: 100
    },
    {
      key: 'status',
      title: $t('page.manage.role.status'),
      align: 'center',
      width: 100,
      render: (row: Role) => {
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
      render: (row: Role) => {
        const canEdit = hasAuth('system:role:edit')
        const canDelete = hasAuth('system:role:delete')

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
  // closeDrawer
} = useTableOperate<Role>(data, getData);

async function handleBatchDelete() {
  try {
    await systemServiceClient.deleteRoles({ ids: checkedRowKeys.value });
    onBatchDeleted();
  } catch {
    window.$message?.error($t('common.deleteFailed'));
  }
}

async function handleDelete(id: bigint) {
  try {
    await systemServiceClient.deleteRoles({ ids: [id] });
    onDeleted();
  } catch {
    window.$message?.error($t('common.deleteFailed'));
  }
}

function edit(id: number) {
  handleEdit(id);
}

function handleSubmitted() {
  drawerVisible.value = false;
  getData();
}
</script>

<template>
  <div class="min-h-500px flex-col-stretch gap-16px overflow-hidden lt-sm:overflow-auto">
    <RoleSearch v-model:model="searchParams" @reset="resetSearchParams" @search="getDataByPage" />
    <NCard :title="$t('page.manage.role.title')" :bordered="false" size="small" class="sm:flex-1-hidden card-wrapper">
      <template #header-extra>
        <TableHeaderOperation
          v-model:columns="columnChecks"
          :disabled-delete="checkedRowKeys.length === 0"
          :show-batch-delete="hasAuth('system:role:delete')"
          :loading="loading"
          :show-add="hasAuth('system:role:edit')"
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
        :scroll-x="702"
        :loading="loading"
        remote
        :row-key="(row: Role) => row.id.toString()"
        :pagination="mobilePagination"
        class="sm:h-full"
      />
      <RoleOperateModal
        v-model:visible="drawerVisible"
        :operate-type="operateType"
        :row-data="editingData"
        @submitted="handleSubmitted"
      />
    </NCard>
  </div>
</template>

<style scoped></style>
