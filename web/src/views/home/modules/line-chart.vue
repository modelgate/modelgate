<script setup lang="ts">
import { watch, ref, computed } from 'vue';
import { useAppStore } from '@/store/modules/app';
import { useEcharts } from '@/hooks/common/echarts';
import { $t } from '@/locales';
import { relayServiceClient } from '@/grpc';
import { create } from '@bufbuild/protobuf';
import { GetRelayUsageRequestSchema } from '@/typings/proto/admin/v1/relay_pb';

defineOptions({
  name: 'LineChart'
});

const appStore = useAppStore();

const colors = ['#8e9dff', '#26deca', '#ffa502', '#ff6b6b', '#a55eea', '#45aaf2', '#fd9644', '#fc5c65'];

// Chart type tabs
type ChartTabType = 'request' | 'point';
const chartTab = ref<ChartTabType>('request');

// Time filter options
type TimeFilterType = 'today' | '7days' | '30days' | 'custom';
const timeFilter = ref<TimeFilterType>('today');
const customDateRange = ref<[number, number] | null>(null);

// Chart tab options
const chartTabOptions = computed(() => [
  { label: $t('page.home.chartTab.request'), value: 'request' },
  { label: $t('page.home.chartTab.point'), value: 'point' }
]);

const { domRef, updateOptions, showLoading, hideLoading } = useEcharts(() => ({
  tooltip: {
    trigger: 'axis',
    axisPointer: {
      type: 'cross',
      label: {
        backgroundColor: '#6a7985'
      }
    }
  },
  legend: {
    data: [] as string[]
  },
  grid: {
    left: '3%',
    right: '4%',
    bottom: '3%',
    containLabel: true
  },
  xAxis: {
    type: 'category',
    boundaryGap: false,
    data: [] as string[]
  },
  yAxis: {
    type: 'value',
    name: $t('page.home.requestVolume')
  },
  series: [] as any[]
}));

// Store series names for locale updates
const seriesNames = ref<string[]>([]);

// Get Y axis name based on chart tab
function getYAxisName() {
  switch (chartTab.value) {
    case 'request':
      return $t('page.home.requestVolume');
    case 'point':
      return $t('page.home.pointVolume');
    default:
      return $t('page.home.requestVolume');
  }
}

// Flag to prevent duplicate requests
const isSyncingDateRange = ref(false);

// Sync customDateRange with timeFilter
function syncDateRange() {
  isSyncingDateRange.value = true;
  const now = new Date();
  let startTime: Date;

  switch (timeFilter.value) {
    case 'today':
      startTime = new Date(now);
      startTime.setHours(0, 0, 0, 0);
      break;
    case '7days':
      startTime = new Date(now.getTime() - 7 * 24 * 60 * 60 * 1000);
      break;
    case '30days':
      startTime = new Date(now.getTime() - 30 * 24 * 60 * 60 * 1000);
      break;
    default:
      isSyncingDateRange.value = false;
      return;
  }

  customDateRange.value = [startTime.getTime(), now.getTime()];
  // Reset flag after a tick to ensure the watch has processed the change
  setTimeout(() => {
    isSyncingDateRange.value = false;
  }, 0);
}

// Calculate time range based on filter
function getTimeRange() {
  const now = new Date();
  let startTime: Date;

  switch (timeFilter.value) {
    case 'today':
      startTime = new Date(now);
      startTime.setHours(0, 0, 0, 0);
      break;
    case '7days':
      startTime = new Date(now.getTime() - 7 * 24 * 60 * 60 * 1000);
      break;
    case '30days':
      startTime = new Date(now.getTime() - 30 * 24 * 60 * 60 * 1000);
      break;
    case 'custom':
      if (customDateRange.value) {
        startTime = new Date(customDateRange.value[0]);
        now.setTime(customDateRange.value[1]);
      } else {
        startTime = new Date(now);
        startTime.setHours(0, 0, 0, 0);
      }
      break;
    default:
      startTime = new Date(now);
      startTime.setHours(0, 0, 0, 0);
  }

  return {
    startTime: {
      seconds: BigInt(Math.floor(startTime.getTime() / 1000)),
      nanos: 0
    },
    endTime: {
      seconds: BigInt(Math.floor(now.getTime() / 1000)),
      nanos: 0
    }
  };
}

async function fetchData() {
  showLoading();
  try {
    const { startTime, endTime } = getTimeRange();

    const request = create(GetRelayUsageRequestSchema, {
      startTime,
      endTime,
      chartType: chartTab.value
    });

    const response = await relayServiceClient.getRelayUsage(request);

    // Build x-axis data (use labels from first series)
    const xAxisData: string[] = [];
    if (response.series.length > 0) {
      response.series[0].data.forEach(item => {
        xAxisData.push(item.label);
      });
    }

    // Build series dynamically
    const newSeries = response.series.map((serie, index) => ({
      color: colors[index % colors.length],
      name: serie.name,
      type: 'line',
      smooth: true,
      data: serie.data.map(item => Number(item.value))
    }));

    // Store series names for locale update
    seriesNames.value = response.series.map(s => s.name);

    updateOptions(opts => {
      opts.legend.data = seriesNames.value;
      opts.xAxis.data = xAxisData;
      opts.series = newSeries;
      opts.yAxis.name = getYAxisName();

      return opts;
    });
  } catch (error) {
    console.error('Failed to fetch chart data:', error);
    window.$message?.error($t('common.requestFailed'));
  } finally {
    hideLoading();
  }
}

function updateLocale() {
  updateOptions((opts, factory) => {
    const originOpts = factory();
    // Re-translate legend data if needed
    opts.legend.data = seriesNames.value;
    opts.series = opts.series.map((s: any) => ({
      ...s,
      name: s.name
    }));
    // Update Y axis name
    opts.yAxis.name = getYAxisName();
    return opts;
  });
}

async function init() {
  syncDateRange();
  fetchData();
}

// Watch for chart tab changes
watch(chartTab, () => {
  fetchData();
});

// Watch for time filter changes
watch(timeFilter, () => {
  if (timeFilter.value !== 'custom') {
    syncDateRange();
  }
  fetchData();
});

// When user manually changes date range, set filter to custom
watch(customDateRange, () => {
  // Skip if the change was triggered by syncDateRange
  if (isSyncingDateRange.value) {
    return;
  }
  if (customDateRange.value) {
    timeFilter.value = 'custom';
  }
  fetchData();
});

watch(
  () => appStore.locale,
  () => {
    updateLocale();
  }
);

// init
init();
</script>

<template>
  <NCard :bordered="false" class="card-wrapper">
    <NSpace vertical :size="16">
      <!-- Chart Tabs and Time Filter -->
      <div class="header-container">
        <div class="compact-tabs">
          <NTabs v-model:value="chartTab" type="segment" size="small">
            <NTab v-for="tab in chartTabOptions" :key="tab.value" :name="tab.value" class="px-2!">
              {{ tab.label }}
            </NTab>
          </NTabs>
        </div>
        <NSpace :size="12" align="center">
          <NRadioGroup v-model:value="timeFilter">
            <NRadioButton value="today">{{ $t('page.home.today') }}</NRadioButton>
            <NRadioButton value="7days">{{ $t('page.home.last7Days') }}</NRadioButton>
            <NRadioButton value="30days">{{ $t('page.home.last30Days') }}</NRadioButton>
          </NRadioGroup>
          <NDatePicker
            v-model:value="customDateRange"
            type="daterange"
            clearable
            class="w-280px"
          />
        </NSpace>
      </div>

      <!-- Chart -->
      <div ref="domRef" class="h-360px overflow-hidden"></div>
    </NSpace>
  </NCard>
</template>

<style>
/* 定义容器查询上下文 */
.card-wrapper {
  container-type: inline-size;
  container-name: chart-card;
}

.header-container {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  flex-wrap: wrap;
}

.compact-tabs {
  flex-shrink: 0;
  max-width: 150px;
}

/* 使用容器查询,基于卡片宽度而不是视口宽度 */
@container chart-card (max-width: 768px) {
  .header-container {
    flex-direction: column;
    align-items: stretch;
  }
  
  .compact-tabs {
    max-width: none !important;
    width: 100% !important;
  }
} 
</style>
