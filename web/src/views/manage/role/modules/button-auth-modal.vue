<script setup lang="ts">
import { computed, shallowRef } from 'vue';
import { $t } from '@/locales';
import { systemServiceClient } from '@/grpc';
import { ButtonNode } from '@/typings/proto/model/system/menu_pb';

defineOptions({
  name: 'ButtonAuthModal'
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

const title = computed(() => $t('common.edit') + $t('page.manage.role.buttonAuth'));

const tree = shallowRef<ButtonNode[]>([]);

async function getAllButtons() {
  const {records} = await systemServiceClient.getButtonList({});
  tree.value = records;
}

const checks = shallowRef<number[]>([]);

async function getChecks() {
  const info = await systemServiceClient.getRoleInfo({
    id: BigInt(props.roleId)
  });
  
  const buttons = info.permission?.buttons ?? [];
  tree.value.map(item => {
    item.children?.map(child => {
      if (buttons.includes(child.code)) {
        checks.value.push(child.id);
      }
    })
  })
}

async function handleSubmit() {
  const buttons: string[] = [];
  tree.value.map(item => {
    item.children?.map(child => {
      if (checks.value.includes(child.id)) {
        buttons.push(child.code);
      }
    })
  })

  try {
    await systemServiceClient.updateRolePermission({
      id: BigInt(props.roleId),
      buttons: buttons,
      updateMask: {
        paths: ['buttons']
      }
    })

    window.$message?.success?.($t('common.modifySuccess'));
    closeModal();
  } catch {
    window.$message?.error?.($t('common.modifyFailed'));
  }
}

async function init() {
  await getAllButtons();
  await getChecks();
}

// init
init();
</script>

<template>
  <NModal v-model:show="visible" :title="title" preset="card" class="w-480px">
    <NTree
      v-model:checked-keys="checks"
      :data="tree"
      key-field="id"
      checkable
      cascade
      block-line
      expand-on-click
      virtual-scroll
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
