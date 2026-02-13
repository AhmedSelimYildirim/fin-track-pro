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
              <div class="c-item" @click="changeCurrency('TRY', 'STANDARD', 'TRY')">
                <span class="symbol">₺</span> TRY
              </div>
              <div class="c-item" @click="changeCurrency('USD', 'STANDARD', 'USD')">
                <span class="symbol">$</span> USD
              </div>
              <div class="c-item" @click="changeCurrency('EUR', 'STANDARD', 'EUR')">
                <span class="symbol">€</span> EUR
              </div>
              <div class="c-item" @click="changeCurrency('BTC', 'STANDARD', 'BTC')">
                <span class="symbol">₿</span> BTC
              </div>
              <div class="c-item" @click="changeCurrency('SILVER', 'STANDARD', 'SILVER')">
                <span class="symbol text-icon">Gr</span> SILVER
              </div>

              <div class="c-item has-submenu">
                <div class="gold-trigger">
                  <div style="display:flex; align-items:center; gap:12px;">
                    <span class="symbol text-icon"> </span> GOLD
                  </div>
                  <span class="arrow-right">▶</span>
                </div>
                <div class="submenu">
                  <div class="sub-item" style="--i:1" @click="changeCurrency('GOLD', 'GRAM_24', 'Has (24K)')">Has (24K)</div>
                  <div class="sub-item" style="--i:2" @click="changeCurrency('GOLD', 'GRAM_22', '22 Ayar')">22 Ayar</div>
                  <div class="sub-item" style="--i:3" @click="changeCurrency('GOLD', 'GRAM_18', '18 Ayar')">18 Ayar</div>
                  <div class="sub-item" style="--i:4" @click="changeCurrency('GOLD', 'GRAM_14', '14 Ayar')">14 Ayar</div>
                  <div class="sub-item" style="--i:5" @click="changeCurrency('GOLD', 'GRAM_8', '8 Ayar')">8 Ayar</div>
                  <div class="sub-item" style="--i:6" @click="changeCurrency('GOLD', 'GRAM_4', '4 Ayar')">4 Ayar</div>
                  <div class="sub-item" style="--i:7" @click="changeCurrency('GOLD', 'CEYREK', 'Çeyrek')">Çeyrek</div>
                  <div class="sub-item" style="--i:8" @click="changeCurrency('GOLD', 'YARIM', 'Yarım')">Yarım</div>
                  <div class="sub-item" style="--i:9" @click="changeCurrency('GOLD', 'TAM', 'Tam')">Tam</div>
                  <div class="sub-item" style="--i:10" @click="changeCurrency('GOLD', 'CUMHURIYET', 'Cumhuriyet')">Cumhuriyet</div>
                  <div class="sub-item" style="--i:11" @click="changeCurrency('GOLD', 'GREMSE', 'Gremse')">Gremse</div>
                </div>
              </div>
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
            <div class="no-data-content">
              <span>{{ t('noData') }}</span>
            </div>
          </div>
        </div>
        <div class="chart-line"></div>
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
  const targetVariant = ref('STANDARD');
  const baseCurrencyLabel = ref('TRY');

  const cardConfigs = [
    { type: 'TRY', label: 'TRY', icon: '₺' },
    { type: 'USD', label: 'USD', icon: '$' },
    { type: 'EUR', label: 'EUR', icon: '€' },
    { type: 'BTC', label: 'BTC', icon: '₿' },
    { type: 'SILVER', label: 'SILVER', icon: 'Gr' },
    { type: 'GOLD', label: 'GOLD', icon: ' ' }
  ];

  const toggleDropdown = () => { showSelector.value = !showSelector.value; };
  window.addEventListener('click', () => { if(showSelector.value) showSelector.value = false; });

  const displayCurrency = computed(() => baseCurrencyLabel.value);

  const totalValue = computed(() => {
    if (summaryData.value && summaryData.value.total_value) {
      const decimals = ['GOLD', 'SILVER', 'BTC'].includes(baseCurrency.value) ? 3 : 2;
      return summaryData.value.total_value.toLocaleString('tr-TR', { minimumFractionDigits: decimals, maximumFractionDigits: decimals });
    }
    return '0,00';
  });

  const hasData = computed(() => summaryData.value?.assets?.length > 0);

  const chartData = computed(() => {
    const labels = ['TRY', 'USD', 'EUR', 'BTC', 'SILVER', 'GOLD'];
    const colors = ['#E63946', '#2A9D8F', '#A0522D', '#1B1B1B', '#A9A9A9', '#FFD700'];
    const data = labels.map(label => {
      const assets = summaryData.value?.assets || [];
      return assets.filter(a => a.type === label).reduce((sum, curr) => sum + (curr.allocation || 0), 0);
    });
    return { labels, datasets: [{ backgroundColor: colors, data, borderWidth: 0 }] };
  });

  const chartOptions = { responsive: true, cutout: '85%', plugins: { legend: { display: false } } };

  const fetchData = async () => {
    try {
      const res = await api.get('/assets/summary', {
        params: {
          currency: baseCurrency.value,
          variant: targetVariant.value
        }
      });
      summaryData.value = res.data;
    } catch(e) { console.error(e); }
  };

  const getAmount = (type) => {
    const assets = summaryData.value?.assets || [];
    const total = assets.filter(a => a.type === type).reduce((sum, curr) => sum + curr.value_in_base, 0);
    if (total === 0) return '0,00';
    const decimals = 2;
    return total.toLocaleString('tr-TR', { maximumFractionDigits: decimals, minimumFractionDigits: decimals });
  };

  const getAllocation = (type) => {
    const assets = summaryData.value?.assets || [];
    const total = assets.filter(a => a.type === type).reduce((sum, curr) => sum + (curr.allocation || 0), 0);
    return total.toFixed(1);
  };

  const changeCurrency = (code, variant, label) => {
    baseCurrency.value = code;
    targetVariant.value = variant;
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

  .top-bar { display: flex; justify-content: space-between; align-items: center; margin-bottom: 30px; flex-wrap: wrap; gap: 20px; }
  .page-title h2 { color: var(--text-color); margin: 0; font-size: 1.8rem; font-weight: 700; }

  .currency-wrapper { position: relative; z-index: 100; }

  .currency-btn {
    background: var(--sidebar-bg);
    border: 1px solid var(--border-color);
    padding: 12px 25px;
    border-radius: 25px;
    cursor: pointer;
    font-weight: bold;
    color: var(--accent-color);
    transition: 0.3s;
    display: flex;
    align-items: center;
    justify-content: space-between;
    min-width: 220px;
    box-shadow: 0 4px 10px rgba(0,0,0,0.2);
  }
  .currency-btn:hover { border-color: var(--accent-color); transform: translateY(-2px); }
  .currency-btn.active .arrow-icon { transform: rotate(180deg); }
  .arrow-icon { font-size: 0.8rem; transition: 0.3s; }

  .currency-dropdown {
    position: absolute;
    top: calc(100% + 10px);
    inset-inline-end: 0;
    background: var(--sidebar-bg);
    border: 1px solid var(--border-color);
    border-radius: 12px;
    width: 240px;
    box-shadow: 0 10px 30px rgba(0,0,0,0.5);
    z-index: 101;
    padding: 8px 0;
    overflow: visible;
  }

  .c-item {
    padding: 12px 20px;
    cursor: pointer;
    position: relative;
    color: var(--text-color);
    transition: 0.2s;
    font-weight: 600;
    display: flex;
    align-items: center;
    gap: 12px;
    border-bottom: 1px solid var(--border-color);
  }
  .c-item:last-child { border-bottom: none; }
  .c-item:hover { background: var(--hover-bg); color: var(--accent-color); }

  .symbol { width: 24px; text-align: center; font-weight: bold; font-size: 1.1rem; }
  .text-icon { font-size: 0.85rem; font-weight: 800; opacity: 0.9; }

  .has-submenu { position: relative; overflow: visible; }
  .gold-trigger { display: flex; align-items: center; justify-content: space-between; width: 100%; }

  .submenu {
    position: absolute;
    top: 0;
    inset-inline-end: 100%;
    margin-inline-end: -1px;
    background: var(--sidebar-bg);
    border: 1px solid var(--border-color);
    border-radius: 12px;
    min-width: 160px;
    padding: 8px;
    box-shadow: 0 10px 30px rgba(0,0,0,0.5);
    opacity: 0;
    visibility: hidden;
    transform: translateX(10px);
    transition: all 0.2s ease;
  }

  html[dir="ltr"] .submenu { transform: translateX(-10px); }
  html[dir="rtl"] .submenu { transform: translateX(10px); margin-inline-end: -1px; }

  .has-submenu::after { content: ''; position: absolute; top: 0; bottom: 0; width: 20px; inset-inline-end: 100%; z-index: 102; }
  .has-submenu:hover .submenu { opacity: 1; visibility: visible; transform: translateX(0); }

  .sub-item {
    padding: 10px 15px;
    font-size: 0.9rem;
    color: var(--text-color);
    cursor: pointer;
    transition: 0.2s;
    border-bottom: 1px solid var(--border-color);
  }
  .sub-item:last-child { border-bottom: none; }
  .sub-item:hover { background: var(--hover-bg); color: var(--accent-color); }

  .dropdown-anim-enter-active, .dropdown-anim-leave-active { transition: all 0.2s ease; }
  .dropdown-anim-enter-from, .dropdown-anim-leave-to { opacity: 0; transform: translateY(-10px); }

  .chart-section { display: flex; flex-direction: column; align-items: center; margin-bottom: 40px; position: relative; width: 100%; }
  .chart-wrapper { width: 280px; height: 280px; position: relative; }
  .center-balance { position: absolute; top: 50%; left: 50%; transform: translate(-50%, -50%); text-align: center; pointer-events: none; }
  .center-balance h3 { font-size: 1.8rem; margin: 0; font-weight: 800; color: var(--text-color); white-space: nowrap; }
  .center-balance small { color: var(--text-muted); font-size: 0.9rem; font-weight: bold; letter-spacing: 2px; }
  .chart-line { width: 50px; height: 3px; background: var(--border-color); border-radius: 2px; margin-top: 25px; opacity: 0.7; }

  .no-data-circle { width: 100%; height: 100%; border-radius: 50%; border: 6px dashed var(--border-color); display: flex; align-items: center; justify-content: center; }
  .no-data-content { text-align: center; color: var(--text-muted); font-weight: bold; }

  .assets-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 25px;
    width: 100%;
    max-width: 900px;
    margin: 0 auto;
    padding-bottom: 80px;
  }

  .asset-card {
    padding: 20px;
    border-radius: 16px;
    cursor: pointer;
    transition: all 0.3s ease;
    display: flex;
    justify-content: space-between;
    align-items: center;
    position: relative;
    overflow: hidden;
    box-shadow: 0 4px 10px rgba(0,0,0,0.2);
    min-height: 95px;
    border: 1px solid rgba(255,255,255,0.05);
  }
  .asset-card:hover { transform: translateY(-5px); box-shadow: 0 10px 25px rgba(0,0,0,0.3); }
  .asset-card * { color: white !important; z-index: 2; position: relative; }

  .card-try { background: linear-gradient(135deg, #FF4500, #8B0000); }
  .card-usd { background: linear-gradient(135deg, #32CD32, #006400); }
  .card-eur { background: linear-gradient(135deg, #8B4513, #5D4037); }
  .card-btc { background: linear-gradient(135deg, #4F4F4F, #000000); }
  .card-silver { background: linear-gradient(135deg, #D3D3D3, #708090); }
  .card-gold { background: linear-gradient(135deg, #FFD700, #B8860B); }

  .card-left { display: flex; flex-direction: column; justify-content: center; gap: 4px; }
  .card-symbol { font-size: 1.6rem; font-weight: 800; opacity: 0.9; line-height: 1; }
  .card-name { font-size: 0.75rem; text-transform: uppercase; letter-spacing: 1px; opacity: 0.8; font-weight: bold; }

  .card-right { display: flex; flex-direction: column; align-items: flex-end; justify-content: center; }
  .card-amount { font-size: 1.2rem; font-weight: 800; }
  .card-val { font-size: 0.75rem; font-weight: 600; background: rgba(255,255,255,0.2); padding: 2px 8px; border-radius: 6px; margin-top: 4px; }

  [data-theme="light"] .asset-card { border: 1px solid #475569; }

  .floating-actions { position: fixed; bottom: 30px; inset-inline-end: 30px; display: flex; flex-direction: column; gap: 15px; z-index: 110; }
  .f-btn { width: 60px; height: 60px; border-radius: 50%; border: none; font-size: 24px; cursor: pointer; box-shadow: 0 5px 15px rgba(0,0,0,0.3); transition: 0.3s; color: white; display: flex; align-items: center; justify-content: center; }
  .f-btn:hover { transform: scale(1.1); }
  .excel { background: linear-gradient(135deg, #10B981, #059669); }
  .pdf { background: linear-gradient(135deg, #3B82F6, #2563EB); }

  .modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.85); display: flex; justify-content: center; align-items: center; z-index: 200; backdrop-filter: blur(5px); }
  .large-modal { width: 800px !important; max-width: 95%; padding: 40px; border-radius: 20px; }
  .modal-content { background: var(--card-bg); padding: 30px; border-radius: 20px; border: 1px solid var(--border-color); box-shadow: 0 20px 50px rgba(0,0,0,0.5); }
  .modal-header { display: flex; justify-content: space-between; margin-bottom: 25px; border-bottom: 1px solid var(--border-color); padding-bottom: 15px; }
  .modal-header h3 { color: var(--text-color); margin: 0; font-size: 1.5rem; }
  .modal-header button { background: none; border: none; color: var(--text-color); font-size: 1.5rem; cursor: pointer; transition: 0.2s; }
  .modal-header button:hover { color: var(--danger-color); transform: rotate(90deg); }

  .modal-body-split { display: grid; grid-template-columns: 1fr 1.2fr; gap: 40px; margin-top: 20px; }
  .transaction-form, .transaction-history { display: flex; flex-direction: column; gap: 15px; }
  .transaction-history { border-inline-start: 1px solid var(--border-color); padding-inline-start: 30px; }
  h4 { color: var(--text-muted); margin: 0 0 10px 0; font-size: 1.1rem; }
  .history-list { max-height: 350px; overflow-y: auto; padding-inline-end: 5px; }
  .history-item { background: var(--input-bg); padding: 15px; border-radius: 12px; margin-bottom: 12px; display: flex; justify-content: space-between; align-items: center; border: 1px solid var(--border-color); }
  .tx-info { display: flex; flex-direction: column; gap: 4px; }
  .tx-date { color: var(--text-muted); font-size: 0.8rem; }
  .tx-type.add { color: var(--success-color) !important; font-weight: bold; }
  .tx-type.subtract { color: var(--danger-color) !important; font-weight: bold; }
  .tx-amount { color: var(--text-color); font-weight: bold; font-size: 1.1rem; }
  .receipt-download-btn { background: transparent; border: 1px solid var(--border-color); border-radius: 8px; cursor: pointer; padding: 8px; transition: 0.2s; font-size: 1.2rem; }
  .receipt-download-btn:hover { background: var(--hover-bg); }

  .big-input, .ayar-select, .date-input { width: 100%; padding: 15px; background: var(--input-bg); border: 1px solid var(--border-color); color: var(--text-color); border-radius: 12px; margin-bottom: 15px; font-size: 1rem; box-sizing: border-box; }
  .actions { display: flex; gap: 15px; margin-top: 10px; }
  .actions button { flex: 1; padding: 15px; border-radius: 12px; border: none; font-weight: bold; cursor: pointer; color: white; transition: 0.2s; font-size: 1rem; }
  .add { background: linear-gradient(135deg, #10B981, #047857); }
  .sub { background: linear-gradient(135deg, #EF4444, #B91C1C); }
  .add:hover, .sub:hover { opacity: 0.9; transform: translateY(-2px); }

  @media (max-width: 1024px) { .assets-grid { grid-template-columns: repeat(2, 1fr); max-width: 100%; } }
  @media (max-width: 768px) {
    .assets-grid { grid-template-columns: 1fr; }
    .modal-body-split { grid-template-columns: 1fr; }
    .transaction-history { border-inline-start: none; border-top: 1px solid var(--border-color); padding-inline-start: 0; padding-top: 30px; }
    .chart-wrapper { width: 250px; height: 250px; }
    .center-balance h3 { font-size: 1.5rem; }
  }
</style>