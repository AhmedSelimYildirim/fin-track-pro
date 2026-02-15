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

        <div class="right-menu">
          <div class="currency-dropdown-wrapper" v-click-outside="closeDropdown">
            <button class="dd-btn" @click="showSelector = !showSelector">
              <span class="unit">{{ displayCurrency }}</span>
              <span class="arrow">▼</span>
            </button>

            <transition name="dd-anim">
              <div v-if="showSelector" class="dd-menu">
                <div class="dd-section-title">Para Birimleri</div>
                <div class="dd-item" @click="changeCurrency('TRY', 'STANDARD', 'TRY')"><span class="sym">₺</span> TRY</div>
                <div class="dd-item" @click="changeCurrency('USD', 'STANDARD', 'USD')"><span class="sym">$</span> USD</div>
                <div class="dd-item" @click="changeCurrency('EUR', 'STANDARD', 'EUR')"><span class="sym">€</span> EUR</div>
                <div class="dd-item" @click="changeCurrency('BTC', 'STANDARD', 'BTC')"><span class="sym">₿</span> BTC</div>
                <div class="dd-item" @click="changeCurrency('SILVER', 'STANDARD', 'SILVER')"><span class="sym">Ag</span> SILVER</div>

                <div class="dd-divider"></div>
                <div class="dd-section-title">Altın Birimleri</div>

                <div class="dd-group">
                  <div class="dd-item sub-trigger" @click.stop="toggleGramMenu">
                    <span class="sym">⚖️</span> Gram Altın
                    <span class="arrow-right" :class="{ rotated: gramMenuOpen }">▶</span>
                  </div>
                  <div v-if="gramMenuOpen" class="deep-menu">
                    <div class="dd-item deep-item" v-for="k in karats" :key="k" @click="changeCurrency('GOLD', `GRAM_${k}`, `Gram (${k}K)`)">{{ k }} Ayar</div>
                  </div>
                </div>

                <div class="dd-item" @click="changeCurrency('GOLD', 'CUMHURIYET', 'Cumhuriyet')"><span class="sym">🔴</span> Cumhuriyet</div>
                <div class="dd-item" @click="changeCurrency('GOLD', 'TAM', 'Tam')"><span class="sym">🌕</span> Tam</div>
                <div class="dd-item" @click="changeCurrency('GOLD', 'YARIM', 'Yarım')"><span class="sym">🌓</span> Yarım</div>
                <div class="dd-item" @click="changeCurrency('GOLD', 'CEYREK', 'Çeyrek')"><span class="sym">🔸</span> Çeyrek</div>
                <div class="dd-item" @click="changeCurrency('GOLD', 'GREMSE', 'Gremse')"><span class="sym">🔶</span> Gremse</div>
              </div>
            </transition>
          </div>

          <div class="actions">
            <button class="icon-btn" @click="downloadExcel" title="Spesifik Excel">📊</button>
            <button class="icon-btn" @click="downloadFullPDF" title="Spesifik PDF">📄</button>
          </div>
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
        <div v-for="asset in cardConfigs" :key="asset.type" class="asset-card" :class="'card-' + asset.type.toLowerCase()" @click="navigateToDetail(asset.type)">
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
  const cardConfigs = [{ type: 'TRY', label: 'TRY', icon: '₺' }, { type: 'USD', label: 'USD', icon: '$' }, { type: 'EUR', label: 'EUR', icon: '€' }, { type: 'BTC', label: 'BTC', icon: '₿' }, { type: 'SILVER', label: 'SILVER', icon: 'Ag' }, { type: 'GOLD', label: 'GOLD', icon: '🥇' }];

  const toggleDropdown = () => { showSelector.value = !showSelector.value; };
  const toggleGramMenu = () => { gramMenuOpen.value = !gramMenuOpen.value; };
  const closeDropdown = () => { showSelector.value = false; gramMenuOpen.value = false; };

  const displayCurrency = computed(() => baseCurrencyLabel.value);
  const totalValue = computed(() => summaryData.value?.total_value?.toLocaleString('tr-TR', { minimumFractionDigits: 2 }) || '0,00');
  const hasData = computed(() => summaryData.value?.assets?.length > 0);

  const chartData = computed(() => {
    const labels = ['TRY', 'USD', 'EUR', 'BTC', 'SILVER', 'GOLD'];
    const colors = ['#E63946', '#2A9D8F', '#A0522D', '#1B1B1B', '#A9A9A9', '#FFD700'];
    const data = labels.map(label => summaryData.value?.assets?.filter(a => a.type === label).reduce((sum, curr) => sum + (curr.allocation || 0), 0) || 0);
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
    const total = summaryData.value?.assets?.filter(a => a.type === type).reduce((sum, curr) => sum + curr.value_in_base, 0) || 0;
    return total.toLocaleString('tr-TR', { minimumFractionDigits: 2 });
  };

  const getAllocation = (type) => (summaryData.value?.assets?.filter(a => a.type === type).reduce((sum, curr) => sum + (curr.allocation || 0), 0) || 0).toFixed(1);

  const changeCurrency = (code, variant, label) => {
    baseCurrency.value = code; targetVariant.value = variant; baseCurrencyLabel.value = label;
    showSelector.value = false; gramMenuOpen.value = false; fetchData();
  };

  const downloadFullPDF = async () => {
    // BURADA asset_type GÖNDERİYORUZ Kİ SADECE O GELSİN
    const res = await api.get('/assets/receipt/full', {
      responseType: 'blob',
      params: { currency: baseCurrency.value, variant: targetVariant.value, asset_type: baseCurrency.value }
    });
    const url = window.URL.createObjectURL(new Blob([res.data]));
    const link = document.createElement('a'); link.href = url; link.setAttribute('download', 'Ozel_Rapor.pdf'); link.click();
  };

  const downloadExcel = async () => {
    // BURADA DA asset_type GÖNDERİYORUZ
    const res = await api.get('/assets/export/excel', {
      responseType: 'blob',
      params: { currency: baseCurrency.value, variant: targetVariant.value, asset_type: baseCurrency.value }
    });
    const url = window.URL.createObjectURL(new Blob([res.data]));
    const link = document.createElement('a'); link.href = url; link.setAttribute('download', 'Ozel_Rapor.xlsx'); link.click();
  };

  const navigateToDetail = (type) => router.push({ name: 'asset-detail', params: { type } });

  const vClickOutside = {
    mounted(el, binding) {
      el.clickOutsideEvent = function(event) { if (!(el === event.target || el.contains(event.target))) { binding.value(event, el); } };
      document.body.addEventListener('click', el.clickOutsideEvent);
    },
    unmounted(el) { document.body.removeEventListener('click', el.clickOutsideEvent); }
  };

  onMounted(() => fetchData());
</script>

<style scoped>
  .dashboard-wrapper { position: relative; min-height: 100%; width: 100%; overflow: hidden; }
  .animated-background { position: fixed; inset: 0; z-index: 0; pointer-events: none; }
  .orb { position: absolute; border-radius: 50%; filter: blur(100px); opacity: 0.3; animation: floatOrb 15s infinite alternate ease-in-out; }
  .orb-1 { width: 60vw; height: 60vw; background: #4f46e5; top: -20%; inset-inline-start: -10%; }
  .orb-2 { width: 50vw; height: 50vw; background: #ec4899; bottom: -20%; inset-inline-end: -10%; }
  .orb-3 { width: 40vw; height: 40vw; background: #10b981; top: 40%; inset-inline-start: 40%; }
  @keyframes floatOrb { 0% { transform: translate(0, 0) scale(1); } 100% { transform: translate(50px, 50px) scale(1.1); } }
  .dashboard-content { position: relative; z-index: 10; width: 100%; padding: 30px; box-sizing: border-box; }
  .top-bar { display: flex; justify-content: space-between; align-items: center; margin-bottom: 30px; }

  .right-menu { display: flex; align-items: center; gap: 15px; }
  .currency-dropdown-wrapper { position: relative; }
  .dd-btn { background: rgba(255,255,255,0.1); border: 1px solid rgba(255,255,255,0.2); color: white; padding: 10px 20px; border-radius: 12px; cursor: pointer; font-weight: 700; font-size: 1rem; min-width: 160px; display: flex; justify-content: space-between; align-items: center; }
  .dd-btn:hover { background: rgba(255,255,255,0.15); border-color: #FFD700; }
  .dd-menu { position: absolute; top: 110%; right: 0; background: #1E293B; border: 1px solid #334155; border-radius: 12px; width: 220px; z-index: 100; max-height: 400px; overflow-y: auto; overflow-x: hidden; box-shadow: 0 10px 40px rgba(0,0,0,0.5); }
  .dd-menu::-webkit-scrollbar { width: 6px; }
  .dd-menu::-webkit-scrollbar-thumb { background: #475569; border-radius: 3px; }
  .dd-section-title { font-size: 0.75rem; color: #94A3B8; padding: 10px 15px 5px; font-weight: bold; letter-spacing: 1px; }
  .dd-divider { height: 1px; background: rgba(255,255,255,0.1); margin: 5px 0; }
  .dd-item { padding: 12px 15px; cursor: pointer; border-bottom: 1px solid rgba(255,255,255,0.05); transition: 0.2s; display: flex; align-items: center; gap: 10px; font-size: 0.95rem; }
  .dd-item:hover { background: #334155; color: #FFD700; padding-left: 20px; }

  .sub-menu { background: #0F172A; border-left: 2px solid #FFD700; }
  .sub-trigger { font-weight: bold; color: #FFD700; background: #1E293B; }
  .deep-menu { background: #020617; }
  .deep-item { padding-left: 35px; font-size: 0.9rem; color: #94A3B8; }
  .deep-item:hover { color: #fff; background: #1E293B; }
  .arrow-right { font-size: 0.8rem; transition: 0.3s; margin-left: auto; }
  .arrow-right.rotated { transform: rotate(90deg); }

  .actions { display: flex; gap: 10px; }
  .icon-btn { background: rgba(255,255,255,0.1); border: none; width: 45px; height: 45px; border-radius: 10px; cursor: pointer; font-size: 1.2rem; display: flex; align-items: center; justify-content: center; transition: 0.3s; color: white; }
  .icon-btn:hover { background: rgba(255,255,255,0.2); transform: scale(1.1); }

  .chart-section { display: flex; flex-direction: column; align-items: center; margin-bottom: 40px; }
  .chart-wrapper { width: 280px; height: 280px; position: relative; }
  .center-balance { position: absolute; top: 50%; left: 50%; transform: translate(-50%, -50%); text-align: center; }
  .center-balance h3 { font-size: 1.8rem; margin: 0; color: var(--text-color); }
  .center-balance small { color: var(--text-muted); }
  .chart-line { width: 50px; height: 3px; background: var(--border-color); border-radius: 2px; margin-top: 25px; opacity: 0.7; }

  .no-data-circle { width: 100%; height: 100%; border-radius: 50%; border: 6px dashed var(--border-color); display: flex; align-items: center; justify-content: center; }
  .no-data-content { text-align: center; color: var(--text-muted); font-weight: bold; }

  .assets-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 25px; max-width: 900px; margin: 0 auto; }
  .asset-card { padding: 20px; border-radius: 16px; cursor: pointer; transition: 0.3s; display: flex; justify-content: space-between; align-items: center; color: white !important; box-shadow: 0 4px 15px rgba(0,0,0,0.3); }
  .asset-card:hover { transform: translateY(-5px); box-shadow: 0 10px 30px rgba(0,0,0,0.5); }

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
  .card-val { background: rgba(255,255,255,0.2); padding: 2px 8px; border-radius: 6px; font-size: 0.8rem; }

  @media (max-width: 1024px) { .assets-grid { grid-template-columns: repeat(2, 1fr); max-width: 100%; } }
  @media (max-width: 768px) { .assets-grid { grid-template-columns: 1fr; } }
</style>