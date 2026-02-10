<script setup lang="tsx">
import { ref } from 'vue';
import type { Ref } from 'vue';
import { NButton, NPopconfirm, NTag } from 'naive-ui';
import { useBoolean } from '@sa/hooks';
import { yesOrNoRecord } from '@/constants/common';
import { enableStatusRecord, menuTypeRecord } from '@/constants/business';
import { useAppStore } from '@/store/modules/app';
import { useTable, useTableOperate } from '@/hooks/common/table';
import type { NaiveUI } from '@/typings/naive-ui';
import { $t } from '@/locales';
import SvgIcon from '@/components/custom/svg-icon.vue';
import { systemServiceClient } from '@/grpc';
import { IconType, type Menu, MenuType } from '@/typings/proto/model/system/menu_pb';
import MenuOperateModal, { type OperateType } from './modules/menu-operate-modal.vue';
import { useAuth } from '@/hooks/business/auth';


const { hasAuth } = useAuth();

const appStore = useAppStore();

const { bool: visible, setTrue: openModal } = useBoolean();

const wrapperRef = ref<HTMLElement | null>(null);

const { columns, columnChecks, data, loading, pagination, getData, getDataByPage } = useTable({
  apiFn: systemServiceClient.getMenuList,
  apiParams: {
    orderBy: "order"
  },
  columns: () => [
    {
      type: 'selection',
      align: 'center',
      width: 48
    },
    {
      key: 'id',
      title: $t('page.manage.menu.id'),
      align: 'center'
    },
    {
      key: 'type',
      title: $t('page.manage.menu.menuType'),
      align: 'center',
      width: 80,
      render: (row: Menu) => {
        const tagMap: Record<Exclude<MenuType, MenuType.UNSPECIFIED>, NaiveUI.ThemeColor> = {
          [MenuType.DIRECTORY]: 'default',
          [MenuType.MENU]: 'primary'
        };

        const label = $t(menuTypeRecord[row.type]);
        return <NTag type={tagMap[row.type]}>{label}</NTag>;
      }
    },
    {
      key: 'name',
      title: $t('page.manage.menu.menuName'),
      align: 'left',
      minWidth: 100,
      render: (row: Menu) => {
        const { i18nKey, name } = row;
        const label = i18nKey ? $t(i18nKey as App.I18n.I18nKey) : name;

        return <span>{label}</span>;
      }
    },
    {
      key: 'icon',
      title: $t('page.manage.menu.icon'),
      align: 'center',
      width: 60,
      render: (row: Menu) => {
        const icon = row.iconType === IconType.ICONIFY ? row.icon : undefined;
        const localIcon = row.iconType === IconType.LOCAL ? row.icon : undefined;

        return (
          <div class="flex-center">
            <SvgIcon icon={icon} localIcon={localIcon} class="text-icon" />
          </div>
        );
      }
    },
    {
      key: 'routeName',
      title: $t('page.manage.menu.routeName'),
      align: 'left',
      minWidth: 100
    },
    {
      key: 'routePath',
      title: $t('page.manage.menu.routePath'),
      align: 'left',
      minWidth: 120
    },
    {
      key: 'status',
      title: $t('page.manage.menu.menuStatus'),
      align: 'center',
      width: 80,
      render: (row: Menu) => {
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
      key: 'hideInMenu',
      title: $t('page.manage.menu.hideInMenu'),
      align: 'center',
      width: 80,
      render: (row: Menu) => {
        const hide: CommonType.YesOrNo = row.hideInMenu ? 'Y' : 'N';

        const tagMap: Record<CommonType.YesOrNo, NaiveUI.ThemeColor> = {
          Y: 'error',
          N: 'default'
        };

        const label = $t(yesOrNoRecord[hide]);

        return <NTag type={tagMap[hide]}>{label}</NTag>;
      }
    },
    {
      key: 'pid',
      title: $t('page.manage.menu.parentId'),
      width: 90,
      align: 'center'
    },
    {
      key: 'order',
      title: $t('page.manage.menu.order'),
      align: 'center',
      width: 60
    },
    {
      key: 'operate',
      title: $t('common.operate'),
      align: 'center',
      width: 230,
      render: (row: Menu) => {
        const canEdit = hasAuth('system:menu:edit')
        const canDelete = hasAuth('system:menu:delete')

        if (!canEdit && !canDelete) {
          return null
        }
        return (
        <div class="flex-center gap-8px">
          {canEdit && row.type === MenuType.DIRECTORY && 
            (<NButton type="primary" ghost size="small" onClick={() => handleAddChildMenu(row)}>
              {$t('page.manage.menu.addChildMenu')}
            </NButton>)}
          {canEdit && 
            (<NButton type="primary" ghost size="small" onClick={() => handleEdit(row)}>
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

const { checkedRowKeys, onBatchDeleted, onDeleted } = useTableOperate(data, getData);

const operateType = ref<OperateType>('add');

function handleAdd() {
  operateType.value = 'add';
  openModal();
}

async function handleBatchDelete() {
  try {
    await systemServiceClient.deleteMenus({ ids: checkedRowKeys.value });
    onBatchDeleted();
  } catch {
    window.$message?.error($t('common.deleteFailed'));
  }
}

async function handleDelete(id: bigint) {
  try {
    await systemServiceClient.deleteMenus({ ids: [id] });
    onDeleted();
  } catch {
    window.$message?.error($t('common.deleteFailed'));
  }
}

/** the edit menu data or the parent menu data when adding a child menu */
const editingData: Ref<Menu | null> = ref(null);

function handleEdit(item: Menu) {
  operateType.value = 'edit';
  editingData.value = { ...item };

  openModal();
}

function handleAddChildMenu(item: Menu) {
  operateType.value = 'addChild';

  editingData.value = { ...item };

  openModal();
}

const allPages = ref<string[]>([]);

async function getAllPages() {
   const { records } = await systemServiceClient.getPageList({});
  allPages.value = records;
}

function init() {
  getAllPages();
}

// init
init();
</script>

<template>
  <div ref="wrapperRef" class="flex-col-stretch gap-16px overflow-hidden lt-sm:overflow-auto">
    <NCard :title="$t('page.manage.menu.title')" :bordered="false" size="small" class="sm:flex-1-hidden card-wrapper">
      <template #header-extra>
        <TableHeaderOperation
          v-model:columns="columnChecks"
          :disabled-delete="checkedRowKeys.length === 0"
          :show-batch-delete="hasAuth('system:menu:delete')"
          :loading="loading"
          :show-add="hasAuth('system:menu:edit')"
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
        :scroll-x="1088"
        :loading="loading"
        :row-key="row => row.id.toString()"
        remote
        :pagination="pagination"
        class="sm:h-full"
      />
      <MenuOperateModal
        v-model:visible="visible"
        :operate-type="operateType"
        :row-data="editingData"
        :all-pages="allPages"
        @submitted="getDataByPage"
      />
    </NCard>
  </div>
</template>

<style scoped></style>
