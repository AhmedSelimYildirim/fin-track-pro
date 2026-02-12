<template>
  <div class="dashboard-wrapper">
    <div class="animated-background">
      <div class="orb orb-1"></div>
      <div class="orb orb-2"></div>
      <div class="orb orb-3"></div>
    </div>

    <div class="dashboard-content">
      <header class="top-bar">
        <div class="page-title">
          <h2>{{ t('portfolioSummary') }}</h2>
        </div>
        <div class="currency-wrapper">
          <div class="currency-btn" @click.stop="toggleDropdown" :class="{ 'active': showSelector }">
            <span class="curr-label">{{ displayCurrency }}</span>
            <span class="arrow-icon">▼</span>
          </div>

          <transition name="dropdown-anim">
            <div v-if="showSelector" class="currency-dropdown">
              <div class="c-item" @click="changeCurrency('TRY', 'TRY')">₺ TRY</div>
              <div class="c-item" @click="changeCurrency('USD', 'USD')">$ USD</div>
              <div class="c-item" @click="changeCurrency('EUR', 'EUR')">€ EUR</div>
              <div class="c-item" @click="changeCurrency('BTC', 'BTC')">₿ BTC</div>
              <div class="c-item" @click="changeCurrency('SILVER', 'SILVER')">Gr SILVER</div>
              <div class="c-item" @click="changeCurrency('GOLD', 'GOLD 24K')">Gr GOLD</div>
            </div>
          </transition>
        </div>
      </header>

      <div class="chart-section">
        <div class="chart-wrapper">
          <template v-if="hasData">
            <Doughnut :data="chartData" :options="chartOptions" />
            <div class="center-balance">
              <h3>{{ totalValue }}</h3>
              <small>{{ baseCurrencyLabel }}</small>
            </div>
          </template>
          <div v-else class="no-data-circle">
            <div class="no-data-content"><span>{{ t('noData') }}</span></div>
          </div>
        </div>
      </div>

      <div class="assets-grid">
        <div v-for="asset in cardConfigs" :key="asset.type"
             class="asset-card" :class="'card-' + asset.type.toLowerCase()"
             @click="navigateToDetail(asset.type)">
          <div class="card-left">
            <div class="card-symbol">{{ asset.icon }}</div>
            <div class="card-name">{{ asset.label }}</div>
          </div>
          <div class="card-right">
            <span class="card-amount">{{ getAmount(asset.type) }}</span>
            <span class="card-val">%{{ getAllocation(asset.type) }}</span>
          </div>
        </div>
      </div>

      <div class="floating-actions">
        <button class="f-btn excel" @click="downloadExcel">📊</button>
        <button class="f-btn pdf" @click="downloadReceipt">📄</button>
      </div>
    </div>
  </div>
</template>

<script setup>
  import { ref, computed, onMounted } from 'vue';
  import { useRouter } from 'vue-router';
  import { Chart as ChartJS, ArcElement, Tooltip, Legend } from 'chart.js';
  import { Doughnut } from 'vue-chartjs';
  import api from '../services/api';
  import { t } from '../utils/translations';

  ChartJS.register(ArcElement, Tooltip, Legend);
  const router = useRouter();
  const showSelector = ref(false);
  const summaryData = ref(null);
  const baseCurrency = ref('TRY');
  const baseCurrencyLabel = ref('TRY');

  const cardConfigs = [
    { type: 'TRY', label: 'TRY', icon: '₺' },
    { type: 'USD', label: 'USD', icon: '$' },
    { type: 'EUR', label: 'EUR', icon: '€' },
    { type: 'BTC', label: 'BTC', icon: '₿' },
    { type: 'SILVER', label: 'SILVER', icon: 'Gr' },
    { type: 'GOLD', label: 'GOLD', icon: 'Gr' }
  ];

  const displayCurrency = computed(() => baseCurrencyLabel.value);
  const toggleDropdown = () => { showSelector.value = !showSelector.value; };

  const totalValue = computed(() => {
    return summaryData.value?.total_value?.toLocaleString('tr-TR', { maximumFractionDigits: 2 }) || '0,00';
  });

  const hasData = computed(() => summaryData.value?.assets?.length > 0);

  const chartData = computed(() => {
    const labels = ['TRY', 'USD', 'EUR', 'BTC', 'SILVER', 'GOLD'];
    const colors = ['#E63946', '#2A9D8F', '#A0522D', '#1B1B1B', '#A9A9A9', '#FFD700'];
    const data = labels.map(label => {
      return summaryData.value?.assets?.filter(a => a.type === label).reduce((sum, curr) => sum + (curr.allocation || 0), 0) || 0;
    });
    return { labels, datasets: [{ backgroundColor: colors, data, borderWidth: 0 }] };
  });

  const chartOptions = { responsive: true, cutout: '85%', plugins: { legend: { display: false } } };

  const fetchData = async () => {
    try {
      const res = await api.get('/assets/summary', { params: { currency: baseCurrency.value } });
      summaryData.value = res.data;
    } catch(e) { console.error(e); }
  };

  const getAmount = (type) => {
    const total = summaryData.value?.assets?.filter(a => a.type === type).reduce((sum, curr) => sum + curr.amount, 0) || 0;
    return total.toLocaleString('tr-TR', { maximumFractionDigits: 2 });
  };

  const getAllocation = (type) => {
    const total = summaryData.value?.assets?.filter(a => a.type === type).reduce((sum, curr) => sum + (curr.allocation || 0), 0) || 0;
    return total.toFixed(1);
  };

  const changeCurrency = (code, label) => {
    baseCurrency.value = code; baseCurrencyLabel.value = label;
    showSelector.value = false; fetchData();
  };

  const navigateToDetail = (type) => {
    router.push({ name: 'asset-detail', params: { type } });
  };

  const downloadReceipt = async () => {
    const res = await api.get('/assets/receipt/full', { responseType: 'blob', params: { currency: baseCurrency.value } });
    const url = window.URL.createObjectURL(new Blob([res.data]));
    const link = document.createElement('a');
    link.href = url; link.setAttribute('download', 'Rapor.pdf'); link.click();
  };

  const downloadExcel = async () => {
    const res = await api.get('/assets/export/excel', { responseType: 'blob', params: { currency: baseCurrency.value } });
    const url = window.URL.createObjectURL(new Blob([res.data]));
    const link = document.createElement('a');
    link.href = url; link.setAttribute('download', 'Export.xlsx'); link.click();
  };

  onMounted(fetchData);
</script>

<style scoped>
  .dashboard-wrapper { position: relative; min-height: 100%; width: 100%; overflow: hidden; }
  .dashboard-content { position: relative; z-index: 10; padding: 30px; }
  .assets-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 25px; max-width: 900px; margin: 0 auto; }
  .asset-card { padding: 20px; border-radius: 16px; cursor: pointer; transition: 0.3s; display: flex; justify-content: space-between; color: white !important; }
  .asset-card:hover { transform: translateY(-5px); }
  .card-try { background: linear-gradient(135deg, #FF4500, #8B0000); }
  .card-usd { background: linear-gradient(135deg, #32CD32, #006400); }
  .card-eur { background: linear-gradient(135deg, #8B4513, #5D4037); }
  .card-btc { background: linear-gradient(135deg, #4F4F4F, #000000); }
  .card-silver { background: linear-gradient(135deg, #D3D3D3, #708090); }
  .card-gold { background: linear-gradient(135deg, #FFD700, #B8860B); }
  .floating-actions { position: fixed; bottom: 30px; right: 30px; display: flex; flex-direction: column; gap: 15px; }
  .f-btn { width: 50px; height: 50px; border-radius: 50%; border: none; cursor: pointer; color: white; }
  .excel { background: #10B981; } .pdf { background: #3B82F6; }
</style>