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
          <div class="lcd-label">TOPLAM VARLIK</div>
          <div class="lcd-value">
            {{ totalValue }} <span class="unit">{{ baseCurrencyLabel }}</span>
          </div>
        </div>
      </header>

      <div class="currency-ribbon">
        <div
                v-for="curr in currencies"
                :key="curr.code"
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
            <div class="no-data-content">
              <span>{{ t('noData') }}</span>
            </div>
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
        <button class="f-btn excel" @click="downloadExcel" title="Excel">📊</button>
        <button class="f-btn pdf" @click="downloadReceipt" title="PDF">📄</button>
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
    { code: 'GOLD', variant: 'GRAM_24', label: 'Gr GOLD', short: 'Gr 24K', icon: '🥇' },
    { code: 'GOLD', variant: 'CEYREK', label: 'Çeyrek', short: 'Çeyrek', icon: '🔸' },
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
    if (summaryData.value && summaryData.value.total_value) {
      const decimals = baseCurrency.value === 'GOLD' || baseCurrency.value === 'BTC' ? 3 : 2;
      return summaryData.value.total_value.toLocaleString('tr-TR', { minimumFractionDigits: decimals, maximumFractionDigits: decimals });
    }
    return '0,00';
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
      const res = await api.get('/assets/summary', {
        params: { currency: baseCurrency.value }
      });
      summaryData.value = res.data;
    } catch(e) { console.error(e); }
  };

  const getAmount = (type) => {
    const total = summaryData.value?.assets?.filter(a => a.type === type).reduce((sum, curr) => sum + curr.amount, 0) || 0;
    const decimals = type === 'GOLD' || type === 'SILVER' || type === 'BTC' ? 3 : 2;
    return total.toLocaleString('tr-TR', { maximumFractionDigits: decimals });
  };

  const getAllocation = (type) => {
    const total = summaryData.value?.assets?.filter(a => a.type === type).reduce((sum, curr) => sum + (curr.allocation || 0), 0) || 0;
    return total.toFixed(1);
  };

  const changeCurrency = (code, variant, label) => {
    baseCurrency.value = code;
    targetVariant.value = variant;
    baseCurrencyLabel.value = label;
    fetchData();
  };

  const navigateToDetail = (type) => {
    router.push({ name: 'asset-detail', params: { type } });
  };

  const downloadReceipt = async () => {
    const res = await api.get('/assets/receipt/full', { responseType: 'blob', params: { currency: baseCurrency.value } });
    const url = window.URL.createObjectURL(new Blob([res.data]));
    const link = document.createElement('a');
    link.href = url; link.setAttribute('download', 'Genel_Rapor.pdf'); link.click();
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
  .animated-background { position: fixed; inset: 0; overflow: hidden; z-index: 0; pointer-events: none; }
  .orb { position: absolute; border-radius: 50%; filter: blur(100px); opacity: 0.3; animation: floatOrb 15s infinite alternate ease-in-out; }
  .orb-1 { width: 60vw; height: 60vw; background: #4f46e5; top: -20%; inset-inline-start: -10%; }
  .orb-2 { width: 50vw; height: 50vw; background: #ec4899; bottom: -20%; inset-inline-end: -10%; }
  .orb-3 { width: 40vw; height: 40vw; background: #10b981; top: 40%; inset-inline-start: 40%; }
  @keyframes floatOrb { 0% { transform: translate(0, 0) scale(1); } 100% { transform: translate(50px, 50px) scale(1.1); } }

  .dashboard-content { position: relative; z-index: 10; padding: 30px; }

  .top-bar { display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 25px; }
  .page-title h2 { color: var(--text-color); margin: 0; font-size: 1.8rem; font-weight: 700; }

  .live-currency-display {
    background: rgba(255,255,255,0.05);
    padding: 15px 25px;
    border-radius: 16px;
    border: 1px solid rgba(255,255,255,0.1);
    text-align: right;
    backdrop-filter: blur(10px);
    min-width: 200px;
    transition: 0.3s;
  }
  .live-currency-display:hover { transform: translateY(-2px); box-shadow: 0 10px 30px rgba(0,0,0,0.3); }
  .live-currency-display.gold { border-color: #FFD700; background: linear-gradient(135deg, rgba(255,215,0,0.1), transparent); }
  .lcd-label { font-size: 0.75rem; color: var(--text-muted); text-transform: uppercase; letter-spacing: 1px; font-weight: 700; margin-bottom: 5px; }
  .lcd-value { font-size: 1.8rem; font-weight: 800; color: var(--text-color); }
  .lcd-value .unit { font-size: 1rem; color: var(--accent-color); margin-left: 5px; }

  .currency-ribbon {
    display: flex;
    gap: 12px;
    overflow-x: auto;
    padding-bottom: 10px;
    margin-bottom: 30px;
    scrollbar-width: none;
  }
  .currency-ribbon::-webkit-scrollbar { display: none; }

  .ribbon-item {
    background: rgba(255,255,255,0.05);
    border: 1px solid rgba(255,255,255,0.1);
    padding: 10px 20px;
    border-radius: 30px;
    cursor: pointer;
    display: flex;
    align-items: center;
    gap: 8px;
    color: var(--text-muted);
    font-weight: 600;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    white-space: nowrap;
  }
  .ribbon-item:hover { background: rgba(255,255,255,0.1); color: white; }
  .ribbon-item.active {
    background: var(--accent-color);
    color: #000;
    font-weight: 800;
    border-color: var(--accent-color);
    transform: scale(1.05);
    box-shadow: 0 5px 20px rgba(0,0,0,0.3);
  }
  .ribbon-icon { font-size: 1.1rem; }

  .chart-section { display: flex; flex-direction: column; align-items: center; margin-bottom: 40px; position: relative; width: 100%; }
  .chart-wrapper { width: 280px; height: 280px; position: relative; }
  .center-balance { position: absolute; top: 50%; left: 50%; transform: translate(-50%, -50%); text-align: center; pointer-events: none; }
  .center-balance h3 { font-size: 1.8rem; margin: 0; font-weight: 800; color: var(--text-color); white-space: nowrap; }
  .center-balance small { color: var(--text-muted); font-size: 0.9rem; font-weight: bold; letter-spacing: 2px; }

  .no-data-circle { width: 100%; height: 100%; border-radius: 50%; border: 6px dashed var(--border-color); display: flex; align-items: center; justify-content: center; }
  .no-data-content { text-align: center; color: var(--text-muted); font-weight: bold; }

  .assets-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 25px; max-width: 900px; margin: 0 auto; }
  .asset-card { padding: 20px; border-radius: 16px; cursor: pointer; transition: 0.3s; display: flex; justify-content: space-between; color: white !important; box-shadow: 0 4px 10px rgba(0,0,0,0.2); position: relative; overflow: hidden; }
  .asset-card:hover { transform: translateY(-5px); box-shadow: 0 10px 30px rgba(0,0,0,0.4); }
  .asset-card * { z-index: 2; position: relative; }

  .card-try { background: linear-gradient(135deg, #DC2626, #7F1D1D); }
  .card-usd { background: linear-gradient(135deg, #059669, #064E3B); }
  .card-eur { background: linear-gradient(135deg, #2563EB, #1E3A8A); }
  .card-btc { background: linear-gradient(135deg, #F59E0B, #B45309); }
  .card-silver { background: linear-gradient(135deg, #64748B, #334155); }
  .card-gold { background: linear-gradient(135deg, #FACC15, #854D0E); }

  .card-left { display: flex; flex-direction: column; gap: 5px; }
  .card-symbol { font-size: 1.5rem; font-weight: 800; opacity: 0.9; }
  .card-name { font-size: 0.8rem; text-transform: uppercase; font-weight: 700; opacity: 0.8; }
  .card-right { display: flex; flex-direction: column; align-items: flex-end; }
  .card-amount { font-size: 1.4rem; font-weight: 800; }
  .card-val { font-size: 0.8rem; font-weight: 600; background: rgba(255,255,255,0.2); padding: 3px 8px; border-radius: 6px; margin-top: 5px; }

  .floating-actions { position: fixed; bottom: 30px; right: 30px; display: flex; flex-direction: column; gap: 15px; z-index: 100; }
  .f-btn { width: 55px; height: 55px; border-radius: 50%; border: none; cursor: pointer; color: white; font-size: 1.2rem; box-shadow: 0 5px 15px rgba(0,0,0,0.3); transition: 0.2s; display: flex; align-items: center; justify-content: center; }
  .f-btn:hover { transform: scale(1.1); }
  .excel { background: #10B981; } .pdf { background: #3B82F6; }

  @media (max-width: 768px) {
    .assets-grid { grid-template-columns: 1fr; }
    .chart-wrapper { width: 220px; height: 220px; }
    .top-bar { flex-direction: column; gap: 20px; }
    .live-currency-display { width: 100%; text-align: center; }
  }
</style>