<script setup lang="ts">
import { computed, shallowRef, watch } from 'vue';
import { $t } from '@/locales';
import { systemServiceClient } from '@/grpc';
import { MenuNode } from '@/typings/proto/model/system/menu_pb';

defineOptions({
  name: 'MenuAuthModal'
});

interface Props {
  /** the roleId */
  roleId: string;
}

const props = defineProps<Props>();

const visible = defineModel<boolean>('visible', {
  default: false
});

function closeModal() {
  visible.value = false;
}

const title = computed(() => $t('common.edit') + $t('page.manage.role.menuAuth'));

const pages = shallowRef<string[]>([]);

async function getPages() {
  const { records } = await systemServiceClient.getPageList({});
  pages.value = records;
}

const pageSelectOptions = computed(() => {
  const opts: CommonType.Option[] = pages.value.map(page => ({
    label: page,
    value: page
  }));

  return opts;
});

const tree = shallowRef<MenuNode[]>([]);

async function getTree() {
  const { records } = await systemServiceClient.getMenuTree({});
  tree.value = records;
}

const home = shallowRef<string>('');
const checks = shallowRef<bigint[]>([]);

async function getChecks() {
  const info  = await systemServiceClient.getRoleInfo({
    id: BigInt(props.roleId)
  })

  checks.value = info.permission?.menuIds ?? [];
  home.value = info.permission?.home ?? '';
}

async function updateHome(val: string) {
  home.value = val;
}


async function handleSubmit() {
  try {
    await systemServiceClient.updateRolePermission({
      id: BigInt(props.roleId),
      home: home.value,
      menuIds: checks.value,
      updateMask: {
        paths: ['home', 'menu_ids']
      }
    })

    window.$message?.success?.($t('common.modifySuccess'));
    closeModal();
  } catch {
    window.$message?.error?.($t('common.modifyFailed'));
  }
}

async function init() {
  await getPages();
  await getTree();
  await getChecks();
}

watch(visible, val => {
  if (val) {
    init();
  }
});
</script>

<template>
  <NModal v-model:show="visible" :title="title" preset="card" class="w-480px">
    <div class="flex-y-center gap-16px pb-12px">
      <div>{{ $t('page.manage.menu.home') }}</div>
      <NSelect :value="home" :options="pageSelectOptions" size="small" class="w-160px" @update:value="updateHome" />
    </div>
    <NTree
      v-model:checked-keys="checks"
      :data="tree"
      key-field="id"
      checkable
      cascade
      expand-on-click
      virtual-scroll
      block-line
      class="h-280px"
    />
    <template #footer>
      <NSpace justify="end">
        <NButton size="small" class="mt-16px" @click="closeModal">
          {{ $t('common.cancel') }}
        </NButton>
        <NButton type="primary" size="small" class="mt-16px" @click="handleSubmit">
          {{ $t('common.confirm') }}
        </NButton>
      </NSpace>
    </template>
  </NModal>
</template>

<style scoped></style>
