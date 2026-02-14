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
  const gramMenuOpen = ref(false); // YENİ: Gram menüsü kontrolü
  const summaryData = ref(null);

  const baseCurrency = ref('TRY');
  const targetVariant = ref('STANDARD');
  const baseCurrencyLabel = ref('TRY');

  const karats = [24, 22, 18, 14, 8, 4]; // YENİ: Ayar listesi

  const cardConfigs = [
    { type: 'TRY', label: 'TRY', icon: '₺' },
    { type: 'USD', label: 'USD', icon: '$' },
    { type: 'EUR', label: 'EUR', icon: '€' },
    { type: 'BTC', label: 'BTC', icon: '₿' },
    { type: 'SILVER', label: 'SILVER', icon: 'Gr' },
    { type: 'GOLD', label: 'GOLD', icon: '🥇' }
  ];

  const toggleDropdown = () => { showSelector.value = !showSelector.value; };
  const toggleGramMenu = () => { gramMenuOpen.value = !gramMenuOpen.value; }; // YENİ

  const closeDropdown = () => {
    showSelector.value = false;
    gramMenuOpen.value = false;
  };

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
    gramMenuOpen.value = false; // Menü kapanınca gramı da kapat
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

  /* GÜNCELLENMİŞ DROPDOWN STİLİ (Ağaç Yapısı İçin) */
  .currency-dropdown {
    position: absolute;
    top: calc(100% + 10px);
    inset-inline-end: 0;
    background: #1E293B; /* Koyu arka plan */
    border: 1px solid #334155;
    border-radius: 12px;
    width: 240px;
    box-shadow: 0 10px 40px rgba(0,0,0,0.5);
    z-index: 101;
    padding: 8px 0;
    overflow-y: auto;
    max-height: 400px;
  }
  .currency-dropdown::-webkit-scrollbar { width: 6px; }
  .currency-dropdown::-webkit-scrollbar-thumb { background: #475569; border-radius: 3px; }

  .dd-section-title { font-size: 0.75rem; color: #94A3B8; padding: 10px 15px 5px; font-weight: bold; letter-spacing: 1px; }
  .dd-divider { height: 1px; background: rgba(255,255,255,0.1); margin: 5px 0; }

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
  .c-item:hover { background: #334155; color: #FFD700; padding-left: 25px; }

  /* GRAM ALTIN ÖZEL STİLİ */
  .sub-menu { background: #0F172A; border-left: 2px solid #FFD700; }
  .sub-item { padding-left: 35px; font-size: 0.9rem; color: #94A3B8; }
  .sub-item:hover { color: #fff; background: #1E293B; }
  .group-trigger { color: #FFD700; background: #1E293B; }
  .arrow-right { font-size: 0.8rem; transition: 0.3s; margin-left: auto; }
  .arrow-right.rotated { transform: rotate(90deg); }

  .symbol { width: 24px; text-align: center; font-weight: bold; font-size: 1.1rem; }
  .text-icon { font-size: 0.85rem; font-weight: 800; opacity: 0.9; }

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

  @media (max-width: 1024px) { .assets-grid { grid-template-columns: repeat(2, 1fr); max-width: 100%; } }
  @media (max-width: 768px) {
    .assets-grid { grid-template-columns: 1fr; }
    .chart-wrapper { width: 250px; height: 250px; }
    .center-balance h3 { font-size: 1.5rem; }
  }
</style>