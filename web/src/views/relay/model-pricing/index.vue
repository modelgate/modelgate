<script setup lang="tsx">
import { NButton, NPopconfirm, NTag } from 'naive-ui';
import { useAppStore } from '@/store/modules/app';
import { useTable, useTableOperate } from '@/hooks/common/table';
import { relayServiceClient } from '@/grpc';
import { $t } from '@/locales';
import ModelPricingOperateModal from './modules/model-pricing-operate-modal.vue';
import ModelPricingSearch from './modules/model-pricing-search.vue';
import { formatProtoTime,formatCurrency } from '@/utils/common';
import { ModelPricing } from '@/typings/proto/model/relay/model_pricing_pb';
import { enableStatusRecord } from '@/constants/business';
import { NaiveUI } from '@/typings/naive-ui';
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
  apiFn: relayServiceClient.getModelPricingList,
  showTotal: true,
  apiParams: {
    current: 1,
    size: 15,
    orderBy: '-id',
    // if you want to use the searchParams in Form, you need to define the following properties, and the value is null
    // the value can not be undefined, otherwise the property in Form will not be reactive
    providerCode: '',
    modelCode :'',
    currency:'',
    effectiveFrom: null,
    effectiveTo: null,
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
      width: 64
    },
    {
      key: 'providerCode',
      title: $t('page.relay.modelPricing.providerCode'),
      align: 'left',
      width: 100,
      render: (row: ModelPricing) => row.providerCode,
    },
    {
      key: 'modelCode',
      title: $t('page.relay.modelPricing.modelCode'),
      align: 'left',
      width: 100,
      render: (row: ModelPricing) => row.modelCode,
    },
    {
      key: 'pointsPerCurrency',
      title: $t('page.relay.modelPricing.pointsPerCurrency'),
      align: 'right',
      width: 150,
      render: (row: ModelPricing) => Number(row.pointsPerCurrency).toLocaleString() + ' / ' + row.currency,
    },
    {
      key: 'tokenNum',
      title: $t('page.relay.modelPricing.tokenNum'),
      align: 'right',
      width: 100,
      render: (row: ModelPricing) => Number(row.tokenNum).toLocaleString(),
    },
    {
      key: 'inputPrice',
      title: $t('page.relay.modelPricing.inputPrice'),
      align: 'right',
      width: 100,
      render: (row: ModelPricing) => formatCurrency(row.inputPrice, row.currency),
    },
    {
      key: 'inputCachePrice',
      title: $t('page.relay.modelPricing.inputCachePrice'),
      align: 'right',
      width: 100,
      render: (row: ModelPricing) => formatCurrency(row.inputCachePrice, row.currency),
    },
    {
      key: 'outputPrice',
      title: $t('page.relay.modelPricing.outputPrice'),
      align: 'right',
      width: 100,
      render: (row: ModelPricing) => formatCurrency(row.outputPrice, row.currency),
    },
    {
      key: 'effectiveFrom',
      title: $t('page.relay.modelPricing.effectiveFrom'),
      align: 'center',
      width: 100,
      render: (row: ModelPricing) => formatProtoTime(row.effectiveFrom),
    },
    {
      key: 'effectiveTo',
      title: $t('page.relay.modelPricing.effectiveTo'),
      align: 'center',
      width: 100,
      render: (row: ModelPricing) => formatProtoTime(row.effectiveTo),
    },
    {
      key: 'status',
      title: $t('page.relay.modelPricing.status'),
      align: 'center',
      width: 100,
      render: (row: ModelPricing) => {
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
      render: (row: ModelPricing) => {
        const canEdit = hasAuth('relay:model_pricing:edit')
        const canDelete = hasAuth('relay:model_pricing:delete')

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
  useTableOperate<ModelPricing>(data, getData);

async function handleBatchDelete() {
  try {
    await relayServiceClient.deleteModelPricings({ ids: checkedRowKeys.value });
  } catch {
    window.$message?.error($t('common.deleteFailed'));
  }

  onBatchDeleted();
}

async function handleDelete(id: bigint) {
  try {
    await relayServiceClient.deleteModelPricings({ ids: [id] });
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
    <ModelPricingSearch v-model:model="searchParams" @reset="resetSearchParams" @search="getDataByPage" />
    <NCard :title="$t('page.relay.model.title')" :bordered="false" size="small" class="sm:flex-1-hidden card-wrapper">
      <template #header-extra>
        <TableHeaderOperation
          v-model:columns="columnChecks"
          :disabled-delete="checkedRowKeys.length === 0"
          :show-batch-delete="hasAuth('relay:model_pricing:delete')"
          :loading="loading"
          :show-add="hasAuth('relay:model_pricing:edit')"
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
        :row-key="(row: ModelPricing) => row.id.toString()"
        :pagination="mobilePagination"
        @update:sorter="handleSorterChange"
        class="sm:h-full"
      />
      <ModelPricingOperateModal
        v-model:visible="drawerVisible"
        :operate-type="operateType"
        :row-data="editingData"
        @submitted="handleSubmitted"
      />
    </NCard>
  </div>
</template>

<style scoped></style>
