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
            <span class="curr-label">{{ totalValue }} {{ displayCurrency }}</span>
            <span class="arrow-icon">▼</span>
          </div>

          <transition name="dropdown-anim">
            <div v-if="showSelector" class="currency-dropdown">
              <div class="c-item" @click="changeCurrency('TRY', 'TRY')">₺ TRY</div>
              <div class="c-item" @click="changeCurrency('USD', 'USD')">$ USD</div>
              <div class="c-item" @click="changeCurrency('EUR', 'EUR')">€ EUR</div>
              <div class="c-item" @click="changeCurrency('BTC', 'BTC')">₿ BTC</div>
              <div class="c-item" @click="changeCurrency('SILVER', 'SILVER')">Gr Gümüş</div>
              <div class="c-item" @click="changeCurrency('GOLD', 'GOLD 24K')">Gr Altın</div>
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
    { type: 'GOLD', label: 'GOLD', icon: '🥇' }
  ];

  const toggleDropdown = () => { showSelector.value = !showSelector.value; };
  window.addEventListener('click', () => { if(showSelector.value) showSelector.value = false; });

  const displayCurrency = computed(() => baseCurrencyLabel.value);

  const totalValue = computed(() => {
    if (summaryData.value && summaryData.value.total_value) {
      const decimals = ['GOLD', 'SILVER', 'BTC'].includes(baseCurrency.value) ? 3 : 2;
      return summaryData.value.total_value.toLocaleString('tr-TR', { maximumFractionDigits: decimals, minimumFractionDigits: decimals });
    }
    return '0,00';
  });

  const hasData = computed(() => summaryData.value?.assets?.length > 0);

  const chartData = computed(() => {
    const labels = ['TRY', 'USD', 'EUR', 'BTC', 'SILVER', 'GOLD'];
    const colors = ['#DC2626', '#059669', '#2563EB', '#F59E0B', '#64748B', '#FACC15'];
    const data = labels.map(label => {
      const assets = summaryData.value?.assets || [];
      return assets.filter(a => a.type === label).reduce((sum, curr) => sum + (curr.allocation || 0), 0);
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
    const assets = summaryData.value?.assets || [];
    const total = assets.filter(a => a.type === type).reduce((sum, curr) => sum + curr.value_in_base, 0);
    const decimals = 2;
    return total.toLocaleString('tr-TR', { maximumFractionDigits: decimals }) + ' ' + baseCurrency.value;
  };

  const getAllocation = (type) => {
    const assets = summaryData.value?.assets || [];
    const total = assets.filter(a => a.type === type).reduce((sum, curr) => sum + (curr.allocation || 0), 0);
    return total.toFixed(1);
  };

  const changeCurrency = (code, label) => {
    baseCurrency.value = code;
    baseCurrencyLabel.value = label;
    showSelector.value = false;
    fetchData();
  };

  const navigateToDetail = (type) => {
    router.push({ name: 'asset-detail', params: { type } });
  };

  onMounted(() => { fetchData(); });
</script>

<style scoped>
  .dashboard-wrapper { position: relative; min-height: 100%; width: 100%; overflow: hidden; }
  .animated-background { position: fixed; inset: 0; overflow: hidden; z-index: 0; pointer-events: none; }
  .orb { position: absolute; border-radius: 50%; filter: blur(100px); opacity: 0.3; animation: floatOrb 15s infinite alternate ease-in-out; }
  .orb-1 { width: 60vw; height: 60vw; background: #4f46e5; top: -20%; inset-inline-start: -10%; }
  .orb-2 { width: 50vw; height: 50vw; background: #ec4899; bottom: -20%; inset-inline-end: -10%; }
  .orb-3 { width: 40vw; height: 40vw; background: #10b981; top: 40%; inset-inline-start: 40%; }
  @keyframes floatOrb { 0% { transform: translate(0, 0) scale(1); } 100% { transform: translate(50px, 50px) scale(1.1); } }

  .dashboard-content { position: relative; z-index: 10; width: 100%; height: 100%; overflow-y: auto; padding: 30px; box-sizing: border-box; }
  .top-bar { display: flex; justify-content: space-between; align-items: center; margin-bottom: 30px; }
  .page-title h2 { color: var(--text-color); margin: 0; font-size: 1.8rem; font-weight: 700; }

  .currency-btn { background: var(--sidebar-bg); border: 1px solid var(--border-color); padding: 12px 25px; border-radius: 25px; cursor: pointer; font-weight: bold; color: var(--accent-color); display: flex; align-items: center; gap: 10px; min-width: 180px; justify-content: space-between; position: relative; z-index: 100; }
  .currency-dropdown { position: absolute; top: 100%; right: 0; background: var(--sidebar-bg); border: 1px solid var(--border-color); border-radius: 12px; width: 100%; z-index: 101; margin-top: 10px; overflow: hidden; }
  .c-item { padding: 12px; cursor: pointer; color: var(--text-color); transition: 0.2s; border-bottom: 1px solid rgba(255,255,255,0.05); }
  .c-item:hover { background: var(--hover-bg); }

  .chart-section { display: flex; flex-direction: column; align-items: center; margin-bottom: 40px; }
  .chart-wrapper { width: 280px; height: 280px; position: relative; }
  .center-balance { position: absolute; top: 50%; left: 50%; transform: translate(-50%, -50%); text-align: center; }
  .center-balance h3 { font-size: 1.8rem; margin: 0; color: var(--text-color); }
  .center-balance small { color: var(--text-muted); }

  .assets-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 25px; max-width: 900px; margin: 0 auto; }
  .asset-card { padding: 20px; border-radius: 16px; cursor: pointer; transition: 0.3s; display: flex; justify-content: space-between; align-items: center; color: white !important; box-shadow: 0 4px 15px rgba(0,0,0,0.3); }
  .asset-card:hover { transform: translateY(-5px); box-shadow: 0 10px 30px rgba(0,0,0,0.5); }

  .card-try { background: linear-gradient(135deg, #DC2626, #7F1D1D); }
  .card-usd { background: linear-gradient(135deg, #059669, #064E3B); }
  .card-eur { background: linear-gradient(135deg, #2563EB, #1E3A8A); }
  .card-btc { background: linear-gradient(135deg, #F59E0B, #B45309); }
  .card-silver { background: linear-gradient(135deg, #64748B, #334155); }
  .card-gold { background: linear-gradient(135deg, #FACC15, #854D0E); }

  .card-amount { font-size: 1.2rem; font-weight: 800; }
  .card-val { background: rgba(255,255,255,0.2); padding: 2px 8px; border-radius: 6px; font-size: 0.8rem; }

  @media (max-width: 768px) { .assets-grid { grid-template-columns: 1fr; } }
</style>