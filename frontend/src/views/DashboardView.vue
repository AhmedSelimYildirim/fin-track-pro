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

        <div class="live-currency-display" :class="baseCurrency.toLowerCase()">
          <div class="lcd-label">CÜZDAN TOPLAMI</div>
          <div class="lcd-value">
            {{ totalValue }} <span class="unit">{{ baseCurrencyLabel }}</span>
          </div>
        </div>
      </header>

      <div class="currency-ribbon">
        <div
                v-for="curr in currencies"
                :key="curr.label"
                class="ribbon-item"
                :class="{ active: baseCurrency === curr.code && targetVariant === curr.variant }"
                @click="changeCurrency(curr.code, curr.variant, curr.label)"
        >
          <span class="ribbon-icon">{{ curr.icon }}</span>
          <span class="ribbon-label">{{ curr.short }}</span>
        </div>
      </div>

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
  const summaryData = ref(null);
  const baseCurrency = ref('TRY');
  const targetVariant = ref('STANDARD');
  const baseCurrencyLabel = ref('TRY');

  const currencies = [
    { code: 'TRY', variant: 'STANDARD', label: 'TRY', short: 'TL', icon: '₺' },
    { code: 'USD', variant: 'STANDARD', label: 'USD', short: 'USD', icon: '$' },
    { code: 'EUR', variant: 'STANDARD', label: 'EUR', short: 'EUR', icon: '€' },
    { code: 'GOLD', variant: 'GRAM_24', label: 'GRAM_24', short: 'Gr 24K', icon: '🥇' },
    { code: 'GOLD', variant: 'CEYREK', label: 'CEYREK', short: 'Çeyrek', icon: '🔸' },
    { code: 'GOLD', variant: 'TAM', label: 'TAM', short: 'Tam', icon: '🌕' },
    { code: 'BTC', variant: 'STANDARD', label: 'BTC', short: 'BTC', icon: '₿' },
    { code: 'SILVER', variant: 'STANDARD', label: 'SILVER', short: 'Gümüş', icon: '🥈' },
  ];

  const cardConfigs = [
    { type: 'TRY', label: 'TRY', icon: '₺' },
    { type: 'USD', label: 'USD', icon: '$' },
    { type: 'EUR', label: 'EUR', icon: '€' },
    { type: 'BTC', label: 'BTC', icon: '₿' },
    { type: 'SILVER', label: 'SILVER', icon: 'Gr' },
    { type: 'GOLD', label: 'GOLD', icon: 'Gr' }
  ];

  const totalValue = computed(() => {
    if (!summaryData.value?.total_value) return '0,00';
    const decimals = baseCurrency.value === 'GOLD' || baseCurrency.value === 'BTC' ? 3 : 2;
    return summaryData.value.total_value.toLocaleString('tr-TR', { minimumFractionDigits: decimals, maximumFractionDigits: decimals });
  });

  const hasData = computed(() => summaryData.value?.assets?.length > 0);

  const chartData = computed(() => {
    const labels = ['TRY', 'USD', 'EUR', 'BTC', 'SILVER', 'GOLD'];
    const colors = ['#DC2626', '#059669', '#2563EB', '#F59E0B', '#64748B', '#FACC15'];
    const data = labels.map(label => {
      const found = summaryData.value?.assets?.filter(a => a.type === label);
      return found?.reduce((sum, curr) => sum + (curr.allocation || 0), 0) || 0;
    });
    return { labels, datasets: [{ backgroundColor: colors, data, borderWidth: 0 }] };
  });

  const chartOptions = { responsive: true, cutout: '85%', plugins: { legend: { display: false } } };

  const fetchData = async () => {
    try {
      const res = await api.get('/assets/summary', { params: { currency: baseCurrency.value, variant: targetVariant.value } });
      summaryData.value = res.data;
    } catch(e) { console.error(e); }
  };

  const getAmount = (type) => {
    const total = summaryData.value?.assets?.filter(a => a.type === type).reduce((sum, curr) => sum + curr.amount, 0) || 0;
    const decimals = type === 'GOLD' || type === 'BTC' ? 3 : 2;
    return total.toLocaleString('tr-TR', { maximumFractionDigits: decimals });
  };

  const getAllocation = (type) => {
    const total = summaryData.value?.assets?.filter(a => a.type === type).reduce((sum, curr) => sum + (curr.allocation || 0), 0) || 0;
    return total.toFixed(1);
  };

  const changeCurrency = (code, variant, label) => {
    baseCurrency.value = code;
    targetVariant.value = variant;
    baseCurrencyLabel.value = label.replace('_', ' ');
    localStorage.setItem('preferredCurrency', code);
    fetchData();
  };

  const navigateToDetail = (type) => { router.push({ name: 'asset-detail', params: { type } }); };

  const downloadReceipt = async () => {
    const res = await api.get('/assets/receipt/full', { responseType: 'blob', params: { currency: baseCurrency.value } });
    const url = window.URL.createObjectURL(new Blob([res.data]));
    const link = document.createElement('a');
    link.href = url; link.setAttribute('download', 'FinTrack_Rapor.pdf'); link.click();
  };

  const downloadExcel = async () => {
    const res = await api.get('/assets/export/excel', { responseType: 'blob', params: { currency: baseCurrency.value } });
    const url = window.URL.createObjectURL(new Blob([res.data]));
    const link = document.createElement('a');
    link.href = url; link.setAttribute('download', 'FinTrack_Export.xlsx'); link.click();
  };

  onMounted(fetchData);
</script>

<style scoped>
  .dashboard-wrapper { position: relative; min-height: 100%; width: 100%; overflow: hidden; }
  .dashboard-content { position: relative; z-index: 10; padding: 30px; }
  .top-bar { display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 25px; }
  .live-currency-display { background: rgba(255,255,255,0.05); padding: 15px 25px; border-radius: 16px; border: 1px solid rgba(255,255,255,0.1); text-align: right; backdrop-filter: blur(10px); min-width: 220px; }
  .live-currency-display.gold { border-color: #FACC15; background: linear-gradient(135deg, rgba(250, 204, 21, 0.1), transparent); }
  .lcd-label { font-size: 0.7rem; color: var(--text-muted); letter-spacing: 1.5px; font-weight: 800; margin-bottom: 5px; }
  .lcd-value { font-size: 1.8rem; font-weight: 900; }
  .lcd-value .unit { font-size: 0.9rem; color: var(--accent-color); margin-left: 5px; }

  .currency-ribbon { display: flex; gap: 12px; overflow-x: auto; padding: 10px 0 20px 0; margin-bottom: 20px; scrollbar-width: none; }
  .currency-ribbon::-webkit-scrollbar { display: none; }
  .ribbon-item { background: rgba(255,255,255,0.05); border: 1px solid rgba(255,255,255,0.1); padding: 12px 24px; border-radius: 40px; cursor: pointer; display: flex; align-items: center; gap: 10px; color: var(--text-muted); font-weight: 700; transition: 0.3s; white-space: nowrap; }
  .ribbon-item:hover { background: rgba(255,255,255,0.1); transform: translateY(-2px); }
  .ribbon-item.active { background: var(--accent-color); color: #000; border-color: var(--accent-color); box-shadow: 0 8px 20px rgba(250, 204, 21, 0.3); }

  .chart-section { display: flex; flex-direction: column; align-items: center; margin-bottom: 40px; }
  .chart-wrapper { width: 260px; height: 260px; position: relative; }
  .center-balance { position: absolute; top: 50%; left: 50%; transform: translate(-50%, -50%); text-align: center; }
  .center-balance h3 { font-size: 1.6rem; margin: 0; font-weight: 900; }

  .assets-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 25px; max-width: 1000px; margin: 0 auto; }
  .asset-card { padding: 22px; border-radius: 20px; cursor: pointer; transition: 0.4s; display: flex; justify-content: space-between; color: white !important; box-shadow: 0 10px 20px rgba(0,0,0,0.2); }
  .asset-card:hover { transform: translateY(-8px); box-shadow: 0 15px 30px rgba(0,0,0,0.4); }
  .card-try { background: linear-gradient(135deg, #DC2626, #7F1D1D); }
  .card-usd { background: linear-gradient(135deg, #059669, #064E3B); }
  .card-eur { background: linear-gradient(135deg, #2563EB, #1E3A8A); }
  .card-btc { background: linear-gradient(135deg, #F59E0B, #B45309); }
  .card-silver { background: linear-gradient(135deg, #64748B, #334155); }
  .card-gold { background: linear-gradient(135deg, #FACC15, #854D0E); }

  .floating-actions { position: fixed; bottom: 30px; right: 30px; display: flex; flex-direction: column; gap: 15px; }
  .f-btn { width: 55px; height: 55px; border-radius: 50%; border: none; cursor: pointer; color: white; font-size: 1.3rem; box-shadow: 0 8px 20px rgba(0,0,0,0.3); transition: 0.3s; }
  .f-btn:hover { transform: scale(1.15) rotate(5deg); }
</style>