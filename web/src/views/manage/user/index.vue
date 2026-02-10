<script setup lang="tsx">
import { NButton, NPopconfirm, NTag } from 'naive-ui';
import { enableStatusRecord, genderRecord } from '@/constants/business';
import { useAppStore } from '@/store/modules/app';
import { useTable, useTableOperate } from '@/hooks/common/table';
import type { NaiveUI } from '@/typings/naive-ui';
import { systemServiceClient } from '@/grpc';
import { $t } from '@/locales';
import type { User } from '@/typings/proto/model/system/user_pb';
import UserOperateModal from './modules/user-operate-modal.vue';
import UserSearch from './modules/user-search.vue';
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
  resetSearchParams
} = useTable({
  apiFn: systemServiceClient.getUserList,
  showTotal: true,
  apiParams: {
    current: 1,
    size: 15,
    orderBy: '-id',
    // if you want to use the searchParams in Form, you need to define the following properties, and the value is null
    // the value can not be undefined, otherwise the property in Form will not be reactive
    status: '',
    username: '',
    gender: '',
    nickname: '',
    phone: '',
    email: ''
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
      key: 'username',
      title: $t('page.manage.user.username'),
      align: 'left',
      width: 100
    },
    {
      key: 'nickname',
      title: $t('page.manage.user.nickname'),
      align: 'left',
      width: 100
    },
    {
      key: 'phone',
      title: $t('page.manage.user.phone'),
      align: 'left',
      width: 100
    },
    {
      key: 'email',
      title: $t('page.manage.user.email'),
      align: 'left',
      width: 100
    },
    {
      key: 'gender',
      title: $t('page.manage.user.gender'),
      align: 'center',
      width: 50,
      render: (row: User) => {
        if (row.gender === null) {
          return null;
        }
        const tagMap: Record<string, NaiveUI.ThemeColor> = {
          male: 'primary',
          female: 'error',
          unknown: 'warning'
        };

        const label = $t(genderRecord[row.gender]);
        return <NTag type={tagMap[row.gender]}>{label}</NTag>;
      }
    },
    {
      key: 'status',
      title: $t('page.manage.user.status'),
      align: 'center',
      width: 50,
      render: (row: User) => {
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
      render: (row: User) => {
        const canEdit = hasAuth('system:user:edit')
        const canDelete = hasAuth('system:user:delete')

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
  useTableOperate<User>(data, getData);

async function handleBatchDelete() {
  try {
    await systemServiceClient.deleteUsers({ ids: checkedRowKeys.value });
  } catch {
    window.$message?.error($t('common.deleteFailed'));
  }

  onBatchDeleted();
}

async function handleDelete(id: bigint) {
  try {
    await systemServiceClient.deleteUsers({ ids: [id] });
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
    <UserSearch v-model:model="searchParams" @reset="resetSearchParams" @search="getDataByPage" />
    <NCard :title="$t('page.manage.user.title')" :bordered="false" size="small" class="sm:flex-1-hidden card-wrapper">
      <template #header-extra>
        <TableHeaderOperation
          v-model:columns="columnChecks"
          :disabled-delete="checkedRowKeys.length === 0"
          :show-batch-delete="hasAuth('system:user:delete')"
          :loading="loading"
          :show-add="hasAuth('system:user:edit')"
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
        :row-key="(row: User) => row.id.toString()"
        :pagination="mobilePagination"
        class="sm:h-full"
      />
      <UserOperateModal
        v-model:visible="drawerVisible"
        :operate-type="operateType"
        :row-data="editingData"
        @submitted="handleSubmitted"
      />
    </NCard>
  </div>
</template>

<style scoped></style>
