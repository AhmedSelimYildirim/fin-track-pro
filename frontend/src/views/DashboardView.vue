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

        <div class="currency-wrapper" v-click-outside="closeDropdown">
          <div class="currency-btn" @click.stop="toggleDropdown" :class="{ 'active': showSelector }">
            <span class="curr-label">{{ totalValue }} {{ displayCurrency }}</span>
            <span class="arrow-icon">▼</span>
          </div>

          <transition name="dropdown-anim">
            <div v-if="showSelector" class="currency-dropdown">

              <div class="dd-section-title">Para Birimleri</div>
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
                <span class="symbol text-icon">Ag</span> SILVER
              </div>

              <div class="dd-divider"></div>
              <div class="dd-section-title">Altın Birimleri</div>

              <div class="dd-group">
                <div class="c-item group-trigger" @click.stop="toggleGramMenu">
                  <span class="symbol text-icon">⚖️</span> Gram Altın
                  <span class="arrow-right" :class="{ rotated: gramMenuOpen }">▶</span>
                </div>

                <div v-if="gramMenuOpen" class="sub-menu">
                  <div class="c-item sub-item" v-for="k in karats" :key="k"
                       @click="changeCurrency('GOLD', `GRAM_${k}`, `Gram (${k}K)`)">
                    {{ k }} Ayar
                  </div>
                </div>
              </div>

              <div class="c-item" @click="changeCurrency('GOLD', 'CUMHURIYET', 'Cumhuriyet')">
                <span class="symbol text-icon">🔴</span> Cumhuriyet
              </div>
              <div class="c-item" @click="changeCurrency('GOLD', 'TAM', 'Tam')">
                <span class="symbol text-icon">🌕</span> Tam
              </div>
              <div class="c-item" @click="changeCurrency('GOLD', 'YARIM', 'Yarım')">
                <span class="symbol text-icon">🌓</span> Yarım
              </div>
              <div class="c-item" @click="changeCurrency('GOLD', 'CEYREK', 'Çeyrek')">
                <span class="symbol text-icon">🔸</span> Çeyrek
              </div>
              <div class="c-item" @click="changeCurrency('GOLD', 'GREMSE', 'Gremse')">
                <span class="symbol text-icon">🔶</span> Gremse
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
  const gramMenuOpen = ref(false);
  const summaryData = ref(null);

  const baseCurrency = ref('TRY');
  const targetVariant = ref('STANDARD');
  const baseCurrencyLabel = ref('TRY');

  const karats = [24, 22, 18, 14, 8, 4];

  const cardConfigs = [
    { type: 'TRY', label: 'TRY', icon: '₺' },
    { type: 'USD', label: 'USD', icon: '$' },
    { type: 'EUR', label: 'EUR', icon: '€' },
    { type: 'BTC', label: 'BTC', icon: '₿' },
    { type: 'SILVER', label: 'SILVER', icon: 'Ag' },
    { type: 'GOLD', label: 'GOLD', icon: '🥇' }
  ];

  const toggleDropdown = () => { showSelector.value = !showSelector.value; };
  const toggleGramMenu = () => { gramMenuOpen.value = !gramMenuOpen.value; };
  const closeDropdown = () => { showSelector.value = false; gramMenuOpen.value = false; };

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
    gramMenuOpen.value = false;
    fetchData();
  };

  const navigateToDetail = (type) => {
    router.push({ name: 'asset-detail', params: { type } });
  };

  const vClickOutside = {
    mounted(el, binding) {
      el.clickOutsideEvent = function(event) {
        if (!(el === event.target || el.contains(event.target))) { binding.value(event, el); }
      };
      document.body.addEventListener('click', el.clickOutsideEvent);
    },
    unmounted(el) { document.body.removeEventListener('click', el.clickOutsideEvent); }
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
    background: #1E293B;
    border: 1px solid #334155;
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
    color: #FFFFFF;
    transition: 0.2s;
    font-weight: 600;
    display: flex;
    align-items: center;
    gap: 12px;
    border-bottom: 1px solid rgba(255,255,255,0.05);
  }
  .c-item:last-child { border-bottom: none; }
  .c-item:hover { background: #334155; color: #FFD700; padding-left: 20px; }

  .symbol { width: 24px; text-align: center; font-weight: bold; font-size: 1.1rem; }
  .text-icon { font-size: 0.85rem; font-weight: 800; opacity: 0.9; }

  .has-submenu { position: relative; overflow: visible; }
  .gold-trigger { display: flex; align-items: center; justify-content: space-between; width: 100%; }

  .sub-menu { background: #0F172A; border-left: 2px solid #FFD700; }
  .sub-item { padding-left: 35px; font-size: 0.9rem; color: #94A3B8; }
  .sub-item:hover { color: #fff; background: #1E293B; }
  .group-trigger { color: #FFD700; background: #1E293B; }
  .arrow-right { font-size: 0.8rem; transition: 0.3s; margin-left: auto; }
  .arrow-right.rotated { transform: rotate(90deg); }

  .dropdown-anim-enter-active, .dropdown-anim-leave-active { transition: all 0.2s ease; }
  .dropdown-anim-enter-from, .dropdown-anim-leave-to { opacity: 0; transform: translateY(-10px); }

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

  .card-left { display: flex; flex-direction: column; justify-content: center; gap: 4px; }
  .card-symbol { font-size: 1.6rem; font-weight: 800; opacity: 0.9; line-height: 1; }
  .card-name { font-size: 0.75rem; text-transform: uppercase; letter-spacing: 1px; opacity: 0.8; font-weight: bold; }

  .card-right { display: flex; flex-direction: column; align-items: flex-end; justify-content: center; }
  .card-amount { font-size: 1.2rem; font-weight: 800; }
  .card-val { background: rgba(255,255,255,0.2); padding: 2px 8px; border-radius: 6px; font-size: 0.8rem; }

  @media (max-width: 768px) { .assets-grid { grid-template-columns: 1fr; } }
</style>