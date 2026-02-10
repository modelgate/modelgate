<script setup lang="tsx">
import { NButton, NPopconfirm, NTag } from 'naive-ui';
import { modelStatusRecord } from '@/constants/business';
import { useAppStore } from '@/store/modules/app';
import { useTable, useTableOperate } from '@/hooks/common/table';
import type { NaiveUI } from '@/typings/naive-ui';
import { relayServiceClient } from '@/grpc';
import { $t } from '@/locales';
import type { Model } from '@/typings/proto/model/relay/model_pb';
import ModelOperateModal from './modules/model-operate-modal.vue';
import ModelSearch from './modules/model-search.vue';
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
  apiFn: relayServiceClient.getModelList,
  showTotal: true,
  apiParams: {
    current: 1,
    size: 15,
    orderBy: 'priority,-id',
    // if you want to use the searchParams in Form, you need to define the following properties, and the value is null
    // the value can not be undefined, otherwise the property in Form will not be reactive
    // providerId: '',
    providerCode: '',
    code :'',
    name:'',
    // priority:0,
    // weight:0,
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
      title: $t('page.relay.model.name'),
      align: 'left',
      width: 100
    },
    {
      key: 'code',
      title: $t('page.relay.model.code'),
      align: 'left',
      width: 100
    },
    {
      key: 'actualCode',
      title: $t('page.relay.model.actualCode'),
      align: 'left',
      width: 100
    },
    {
      key: 'providerCode',
      title: $t('page.relay.model.providerCode'),
      align: 'left',
      width: 100,
    },
    {
      key: 'priority',
      title: $t('page.relay.model.priority'),
      align: 'right',
      width: 50,
      sorter: true,
    },
    {
      key: 'weight',
      title: $t('page.relay.model.weight'),
      align: 'right',
      width: 50
    },
    {
      key: 'status',
      title: $t('page.relay.model.status'),
      align: 'center',
      width: 50,
      sorter: true,
      render: (row: Model) => {
        if (row.status === null) {
          return null;
        }

        const tagMap: Record<string, NaiveUI.ThemeColor> = {
          'enabled': 'success',
          'disabled': 'warning'
        };

        const label = $t(modelStatusRecord[row.status]);

        return <NTag type={tagMap[row.status]}>{label}</NTag>;
      }
    },
    {
      key: 'operate',
      title: $t('common.operate'),
      align: 'center',
      width: 130,
      render: (row: Model) => {
        const canEdit = hasAuth('relay:model:edit')
        const canDelete = hasAuth('relay:model:delete')

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
  useTableOperate<Model>(data, getData);

async function handleBatchDelete() {
  try {
    await relayServiceClient.deleteModels({ ids: checkedRowKeys.value });
  } catch {
    window.$message?.error($t('common.deleteFailed'));
  }

  onBatchDeleted();
}

async function handleDelete(id: bigint) {
  try {
    await relayServiceClient.deleteModels({ ids: [id] });
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
    <ModelSearch v-model:model="searchParams" @reset="resetSearchParams" @search="getDataByPage" />
    <NCard :title="$t('page.relay.model.title')" :bordered="false" size="small" class="sm:flex-1-hidden card-wrapper">
      <template #header-extra>
        <TableHeaderOperation
          v-model:columns="columnChecks"
          :disabled-delete="checkedRowKeys.length === 0"
          :show-batch-delete="hasAuth('relay:model:delete')"
          :loading="loading"
          :show-add="hasAuth('relay:model:edit')"
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
        :row-key="(row: Model) => row.id.toString()"
        :pagination="mobilePagination"
        @update:sorter="handleSorterChange"
        class="sm:h-full"
      />
      <ModelOperateModal
        v-model:visible="drawerVisible"
        :operate-type="operateType"
        :row-data="editingData"
        @submitted="handleSubmitted"
      />
    </NCard>
  </div>
</template>

<style scoped></style>
