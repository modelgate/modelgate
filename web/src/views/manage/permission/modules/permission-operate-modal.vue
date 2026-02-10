<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { useFormRules, useNaiveForm } from '@/hooks/common/form';
import { $t } from '@/locales';
import type { Permission, ApiPerm } from '@/typings/proto/model/system/permission_pb';
import { ApiPermSchema } from '@/typings/proto/model/system/permission_pb';
import { systemServiceClient } from '@/grpc';
import { NaiveUI } from '@/typings/naive-ui';
import { create } from '@bufbuild/protobuf';

defineOptions({
  name: 'PermissionOperateModal'
});

interface Props {
  /** the type of operation */
  operateType: NaiveUI.TableOperateType;
  /** the edit row data */
  rowData?: Permission | null;
}

const props = defineProps<Props>();

interface Emits {
  (e: 'submitted'): void;
}

const emit = defineEmits<Emits>();

const visible = defineModel<boolean>('visible', {
  default: false
});

const { formRef, validate, restoreValidation } = useNaiveForm();
const { defaultRequiredRule } = useFormRules();

const title = computed(() => {
  const titles: Record<NaiveUI.TableOperateType, string> = {
    add: $t('page.manage.permission.addPermission'),
    edit: $t('page.manage.permission.editPermission')
  };
  return titles[props.operateType];
});

type Model = Pick<Permission, 'id' | 'name' | 'code' | 'data' | 'desc'> & {
  data: ApiPerm[];
};

const model = ref(createDefaultModel());

function createDefaultModel(): Model {
  return {
    id: 0n,
    name: '',
    code: '',
    data: [],
    desc: ''
  };
}

type RuleKey = Extract<keyof Model, 'name' | 'code'>;

const rules: Record<RuleKey, App.Global.FormRule> = {
  name: defaultRequiredRule,
  code: defaultRequiredRule,
};

const isEdit = computed(() => props.operateType === 'edit');

function handleInitModel() {
  model.value = createDefaultModel();

  if (props.operateType === 'edit' && props.rowData) {
    const { data, ...rest } = props.rowData;
    Object.assign(model.value, rest);
    // Deep copy data array to avoid modifying the original
    if (data) {
      model.value.data = data.map(item => create(ApiPermSchema, {
        path: item.path,
        method: item.method
      }));
    }
  }
}

/** the enabled role options */
const apiListOptions = ref<CommonType.Option<string>[]>([]);

async function getApiListOptions() {
  const { records } = await systemServiceClient.getApiList({});
  const options = records.map(item => ({
    label: item,
    value: item
  }));

  apiListOptions.value = [...options];
}

const methodOptions = [
  { label: 'GET', value: 'GET' },
  { label: 'POST', value: 'POST' },
  { label: 'PUT', value: 'PUT' },
  { label: 'DELETE', value: 'DELETE' },
  { label: 'PATCH', value: 'PATCH' }
];

function addApiPerm() {
  model.value.data.push(create(ApiPermSchema, {
    path: '',
    method: 'POST'
  }));
}

function removeApiPerm(index: number) {
  model.value.data.splice(index, 1);
}

function closeDrawer() {
  visible.value = false;
}

async function handleSubmit() {
  await validate();

  // Convert data array to protobuf format
  const validData = model.value.data
    .filter(item => item.path && item.path.trim() !== '')
    .map(item => create(ApiPermSchema, {
      path: item.path,
      method: item.method
    }));

  const submitData = { ...model.value, data: validData };

  // request
  if (props.operateType === 'edit') {
    try {
      await systemServiceClient.updatePermission({
        updateMask: {
          paths: ['name', 'code', 'data', 'desc']
        },
        permission: submitData
      });
      window.$message?.success($t('common.updateSuccess'));
    } catch {
      window.$message?.error($t('common.updateFailed'));
    }
  } else {
    try {
      await systemServiceClient.createPermission({
        permission: submitData
      });
      window.$message?.success($t('common.addSuccess'));
    } catch {
      window.$message?.error($t('common.addFailed'));
    }
  }

  closeDrawer();
  emit('submitted');
}

watch(visible, () => {
  if (visible.value) {
    handleInitModel();
    restoreValidation();
    getApiListOptions();
  }
});
</script>

<template>
  <NModal v-model:show="visible" :title="title" preset="card" class="w-600px">
    <NScrollbar class="h-360px pr-20px">
      <NForm ref="formRef" :model="model" :rules="rules">
        <NFormItem :label="$t('page.manage.permission.name')" path="name">
          <NInput v-model:value="model.name" :placeholder="$t('page.manage.permission.form.name')" />
        </NFormItem>
        <NFormItem :label="$t('page.manage.permission.code')" path="code">
          <NInput v-model:value="model.code" :placeholder="$t('page.manage.permission.form.code')" />
        </NFormItem>
        <NFormItem :label="$t('page.manage.permission.desc')" path="desc">
          <NInput v-model:value="model.desc" :placeholder="$t('page.manage.permission.form.desc')" />
        </NFormItem>
        <NFormItem :label="$t('page.manage.permission.data')">
          <NSpace vertical :size="12" class="w-full">
            <div v-for="(item, index) in model.data" :key="index" class="flex items-center gap-2">
              <NSelect
                v-model:value="item.path"
                :options="apiListOptions"
                :placeholder="$t('page.manage.permission.form.path')"
                filterable
                class="flex-1"
              />
              <NSelect
                v-model:value="item.method"
                :placeholder="$t('page.manage.permission.form.method')"
                :options="methodOptions"
                class="w-120px"
              />
              <NButton text type="error" @click="removeApiPerm(index)">
                <template #icon>
                  <icon-mdi-close />
                </template>
              </NButton>
            </div>
            <NButton dashed @click="addApiPerm">
              <template #icon>
                <icon-mdi-plus />
              </template>
              {{ $t('page.manage.permission.addApiPerm') }}
            </NButton>
          </NSpace>
        </NFormItem>
      </NForm>
    </NScrollbar>
    <template #footer>
      <NSpace justify="end" :size="16">
        <NButton @click="closeDrawer">{{ $t('common.cancel') }}</NButton>
        <NButton type="primary" @click="handleSubmit">{{ $t('common.confirm') }}</NButton>
      </NSpace>
    </template>
  </NModal>
</template>

<style scoped></style>
