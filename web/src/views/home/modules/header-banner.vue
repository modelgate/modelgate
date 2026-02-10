<script setup lang="ts">
import { computed, ref, onMounted } from 'vue';
import { useAppStore } from '@/store/modules/app';
import { useAuthStore } from '@/store/modules/auth';
import { $t } from '@/locales';
import { relayServiceClient } from '@/grpc';

defineOptions({
  name: 'HeaderBanner'
});

const appStore = useAppStore();
const authStore = useAuthStore();

const gap = computed(() => (appStore.isMobile ? 0 : 16));

const relayInfo = ref({
  providerCount: 0n,
  modelCount: 0n,
  apiKeyCount: 0n
});

interface StatisticData {
  id: number;
  label: string;
  value: string;
  align: string;
}

const statisticData = computed<StatisticData[]>(() => [
  {
    id: 0,
    label: $t('page.home.providerCount'),
    align: 'center',
    value: relayInfo.value.providerCount.toString()
  },
  {
    id: 1,
    label: $t('page.home.modelCount'),
    align: 'center',
    value: relayInfo.value.modelCount.toString()
  },
  {
    id: 2,
    label: $t('page.home.apiKeyCount'),
    align: 'center',
    value: relayInfo.value.apiKeyCount.toString()
  }
]);

async function fetchRelayInfo() {
  const resp = await relayServiceClient.getRelayInfo({});
  relayInfo.value = {
    providerCount: resp.providerCount,
    modelCount: resp.modelCount,
    apiKeyCount: resp.apiKeyCount
  };
}

onMounted(() => {
  fetchRelayInfo();
});
</script>

<template>
  <NCard :bordered="false" class="card-wrapper">
    <NGrid :x-gap="gap" :y-gap="16" responsive="screen" item-responsive>
      <NGi span="24 s:24 m:18">
        <div class="flex-y-center">
          <div class="size-72px shrink-0 overflow-hidden rd-1/2">
            <img src="@/assets/imgs/soybean.jpg" class="size-full" />
          </div>
          <div class="pl-12px">
            <h3 class="text-18px font-semibold">
              {{ $t('page.home.greeting', { userName: authStore.userInfo.username }) }}
            </h3>
            <p class="text-#999 leading-30px">{{ $t('page.home.weatherDesc') }}</p>
          </div>
        </div>
      </NGi>
      <NGi span="24 s:24 m:6">
        <NSpace :size="24" justify="end">
          <NStatistic v-for="item in statisticData" :key="item.id" class="whitespace-nowrap" v-bind="item" />
        </NSpace>
      </NGi>
    </NGrid>
  </NCard>
</template>

<style scoped></style>
