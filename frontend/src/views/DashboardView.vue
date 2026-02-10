<template>
  <div class="dashboard-wrapper">
    <header class="top-bar">
      <h2>{{ t('portfolioSummary') }}</h2>

      <div class="currency-bar">
        <div
                v-for="item in currencyItems"
                :key="item.code"
                class="currency-pill"
                :class="{ active: baseCurrency === item.code }"
                @click="selectCurrency(item)"
        >
          {{ item.label }}

          <div v-if="item.code === 'GOLD'" class="gold-ayars">
            <span
                    v-for="k in [24,22,18,14,8,4]"
                    :key="k"
                    :class="{ active: targetAyar === k }"
                    @click.stop="selectGold(k)"
            >
              {{ k }}K
            </span>
          </div>
        </div>
      </div>
    </header>

    <div class="chart-section">
      <Doughnut v-if="hasData" :data="chartData" :options="chartOptions" />
      <div v-else class="no-data-circle">{{ t('noData') }}</div>
    </div>
  </div>
</template>

<script setup>
  import { ref, computed, onMounted } from 'vue'
  import { Doughnut } from 'vue-chartjs'
  import { Chart as ChartJS, ArcElement, Tooltip, Legend } from 'chart.js'
  import api from '../services/api'
  import { t } from '../utils/translations'

  ChartJS.register(ArcElement, Tooltip, Legend)

  const baseCurrency = ref('TRY')
  const targetAyar = ref(0)
  const summaryData = ref(null)

  const currencyItems = [
    { code: 'TRY', label: 'TRY' },
    { code: 'USD', label: 'USD' },
    { code: 'EUR', label: 'EUR' },
    { code: 'BTC', label: 'BTC' },
    { code: 'SILVER', label: 'SILVER' },
    { code: 'GOLD', label: 'GOLD (Gr)' }
  ]

  const selectCurrency = (item) => {
    baseCurrency.value = item.code
    if (item.code !== 'GOLD') targetAyar.value = 0
    fetchData()
  }

  const selectGold = (k) => {
    baseCurrency.value = 'GOLD'
    targetAyar.value = k
    fetchData()
  }

  const fetchData = async () => {
    const res = await api.get('/assets/summary', {
      headers: {
        'X-Currency': baseCurrency.value,
        'X-Ayar': targetAyar.value.toString()
      }
    })
    summaryData.value = res.data
  }

  const hasData = computed(() => summaryData.value?.total_value > 0)

  const chartData = computed(() => ({
    labels: ['TRY','USD','EUR','BTC','SILVER','GOLD'],
    datasets: [{
      data: ['TRY','USD','EUR','BTC','SILVER','GOLD'].map(t =>
              summaryData.value?.assets
                      ?.filter(a => a.type === t)
                      .reduce((s,c)=>s+c.allocation,0) || 0
      ),
      backgroundColor: [
        '#ef4444',
        '#10b981',
        '#8b5cf6',
        '#111111',
        '#9ca3af',
        '#facc15'
      ],
      borderWidth: 0
    }]
  }))

  const chartOptions = {
    cutout: '70%',
    plugins: { legend: { display: false } }
  }

  onMounted(fetchData)
</script>

<style scoped>
  .currency-bar {
    display: flex;
    gap: 12px;
    flex-wrap: wrap;
  }

  .currency-pill {
    padding: 10px 18px;
    border-radius: 999px;
    background: #1e293b;
    cursor: pointer;
    font-weight: 700;
    color: #cbd5e1;
    display: flex;
    align-items: center;
    gap: 10px;
  }

  .currency-pill.active {
    background: #22c55e;
    color: #020617;
  }

  .gold-ayars {
    display: flex;
    gap: 6px;
  }

  .gold-ayars span {
    padding: 4px 8px;
    border-radius: 6px;
    background: #020617;
    font-size: 12px;
    cursor: pointer;
  }

  .gold-ayars span.active {
    background: #facc15;
    color: #000;
    font-weight: 800;
  }

  .chart-section {
    margin-top: 40px;
    display: flex;
    justify-content: center;
  }

  .no-data-circle {
    color: #94a3b8;
  }
</style>
