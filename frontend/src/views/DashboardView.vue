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
              <div class="c-item" @click="changeCurrency('TRY', 0, 'TRY')">TRY</div>
              <div class="c-item" @click="changeCurrency('USD', 0, 'USD')">USD</div>
              <div class="c-item" @click="changeCurrency('EUR', 0, 'EUR')">EUR</div>
              <div class="c-item" @click="changeCurrency('BTC', 0, 'BTC')">BTC</div>
              <div class="c-item" @click="changeCurrency('SILVER', 0, 'SILVER')">SILVER (Gr)</div>

              <div class="c-item has-submenu">
                <div class="gold-trigger">
                  GOLD (Gr) <span class="arrow-right">▶</span>
                </div>
                <div class="submenu">
                  <div class="sub-item" style="--i:1" @click="changeCurrency('GOLD', 24, 'GOLD 24K')">24K</div>
                  <div class="sub-item" style="--i:2" @click="changeCurrency('GOLD', 22, 'GOLD 22K')">22K</div>
                  <div class="sub-item" style="--i:3" @click="changeCurrency('GOLD', 18, 'GOLD 18K')">18K</div>
                  <div class="sub-item" style="--i:4" @click="changeCurrency('GOLD', 14, 'GOLD 14K')">14K</div>
                  <div class="sub-item" style="--i:5" @click="changeCurrency('GOLD', 8, 'GOLD 8K')">8K</div>
                  <div class="sub-item" style="--i:6" @click="changeCurrency('GOLD', 4, 'GOLD 4K')">4K</div>
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
      </div>

      <div class="assets-grid">
        <div v-for="asset in cardConfigs" :key="asset.type"
             class="asset-card" :class="'card-' + asset.type.toLowerCase()"
             @click="openModal(asset.type)">
          <div class="card-icon">{{ asset.icon }}</div>
          <div class="card-info">
            <span class="card-name">{{ asset.label }}</span>
            <span class="card-amount">{{ getAmount(asset.type) }} {{ asset.unit }}</span>
            <span class="card-val">%{{ getAllocation(asset.type) }}</span>
          </div>
        </div>
      </div>

      <div v-if="showModal" class="modal-overlay" @click.self="showModal = false">
        <div class="modal-content large-modal">
          <div class="modal-header">
            <h3>{{ activeAsset }} {{ t('newTransaction') }}</h3>
            <button @click="showModal = false">✕</button>
          </div>

          <div class="modal-body-split">
            <div class="transaction-form">
              <h4>{{ t('newTransaction') }}</h4>
              <input v-model="amount" type="number" :placeholder="t('amount')" class="big-input" />
              <div v-if="activeAsset === 'GOLD'" class="ayar-wrapper">
                <select v-model="modalAyar" class="ayar-select">
                  <option :value="24">24K</option>
                  <option :value="22">22K</option>
                  <option :value="18">18K</option>
                  <option :value="14">14K</option>
                  <option :value="8">8K</option>
                  <option :value="4">4K</option>
                </select>
              </div>
              <input v-model="transactionDate" type="date" class="date-input" />
              <div class="actions">
                <button class="add" @click="handleTransaction('add')">{{ t('add') }}</button>
                <button class="sub" @click="handleTransaction('subtract')">{{ t('subtract') }}</button>
              </div>
            </div>

            <div class="transaction-history">
              <h4>{{ t('history') }} & {{ t('receipt') }}</h4>
              <div class="history-list">
                <div v-if="filteredTransactions.length === 0" class="no-history">{{ t('noData') }}</div>
                <div v-for="tx in filteredTransactions" :key="tx.id" class="history-item">
                  <div class="tx-info">
                    <span class="tx-date">{{ formatDate(tx.transaction_date) }}</span>
                    <span class="tx-type" :class="tx.type">{{ tx.type === 'add' ? '+' : '-' }}</span>
                    <span class="tx-amount">{{ tx.amount }} <span v-if="tx.ayar > 0">({{ tx.ayar }}K)</span></span>
                  </div>
                  <button class="receipt-download-btn" @click="downloadSingleReceipt(tx.id)">📄</button>
                </div>
              </div>
            </div>
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
  import { Chart as ChartJS, ArcElement, Tooltip, Legend } from 'chart.js';
  import { Doughnut } from 'vue-chartjs';
  import api from '../services/api';
  import { t } from '../utils/translations';

  ChartJS.register(ArcElement, Tooltip, Legend);

  const showSelector = ref(false);
  const showModal = ref(false);
  const activeAsset = ref('');
  const amount = ref('');
  const transactionDate = ref(new Date().toISOString().split('T')[0]);
  const modalAyar = ref(24);
  const summaryData = ref(null);
  const allTransactions = ref([]);
  const baseCurrency = ref('TRY');
  const baseCurrencyLabel = ref('TRY');
  const targetAyar = ref(0);

  const cardConfigs = [
    { type: 'TRY', label: 'TRY', icon: '₺', unit: '₺' },
    { type: 'USD', label: 'USD', icon: '$', unit: '$' },
    { type: 'EUR', label: 'EUR', icon: '€', unit: '€' },
    { type: 'BTC', label: 'BTC', icon: '₿', unit: 'BTC' },
    { type: 'SILVER', label: 'SILVER (Gr)', icon: '⚔️', unit: 'Gr' },
    { type: 'GOLD', label: 'GOLD (Gr)', icon: '👑', unit: 'Gr' }
  ];

  const toggleDropdown = () => { showSelector.value = !showSelector.value; };
  window.addEventListener('click', () => { if(showSelector.value) showSelector.value = false; });

  const displayCurrency = computed(() => baseCurrencyLabel.value);

  const totalValue = computed(() => {
    if (summaryData.value && summaryData.value.total_value) {
      return summaryData.value.total_value.toLocaleString('tr-TR', { maximumFractionDigits: 2 });
    }
    return '0,00';
  });

  const hasData = computed(() => {
    return summaryData.value && summaryData.value.assets && summaryData.value.assets.length > 0 && summaryData.value.total_value > 0;
  });

  const chartData = computed(() => {
    const labels = ['TRY', 'USD', 'EUR', 'BTC', 'SILVER', 'GOLD'];
    const colors = ['#EF4444', '#10B981', '#8B4513', '#1a1a1a', '#A0A0A0', '#FFD700'];
    const data = labels.map(label => {
      const assets = summaryData.value?.assets || [];
      return assets.filter(a => a.type === label).reduce((sum, curr) => sum + (curr.allocation || 0), 0);
    });
    return { labels, datasets: [{ backgroundColor: colors, data, borderWidth: 0 }] };
  });

  const chartOptions = { responsive: true, cutout: '75%', plugins: { legend: { display: false } } };

  const fetchData = async () => {
    try {
      const res = await api.get('/assets/summary', {
        headers: { 'X-Currency': baseCurrency.value, 'X-Ayar': targetAyar.value.toString() }
      });
      summaryData.value = res.data;
    } catch(e) { console.error(e); }
  };

  const fetchTransactions = async () => {
    try {
      const res = await api.get('/assets/transactions', { headers: { 'X-Currency': baseCurrency.value } });
      allTransactions.value = res.data || [];
    } catch (e) { console.error(e); }
  };

  const filteredTransactions = computed(() => {
    if (!activeAsset.value) return [];
    return allTransactions.value.filter(tx => tx.asset_type === activeAsset.value).sort((a, b) => new Date(b.transaction_date) - new Date(a.transaction_date));
  });

  const getAmount = (type) => {
    const assets = summaryData.value?.assets || [];
    const total = assets.filter(a => a.type === type).reduce((sum, curr) => sum + curr.amount, 0);
    return total.toLocaleString('tr-TR', { maximumFractionDigits: 4 });
  };

  const getAllocation = (type) => {
    const assets = summaryData.value?.assets || [];
    const total = assets.filter(a => a.type === type).reduce((sum, curr) => sum + (curr.allocation || 0), 0);
    return total.toFixed(1);
  };

  const openModal = (asset) => {
    activeAsset.value = asset;
    amount.value = '';
    modalAyar.value = 24;
    transactionDate.value = new Date().toISOString().split('T')[0];
    showModal.value = true;
    fetchTransactions();
  };

  const handleTransaction = async (action) => {
    if (!amount.value) return alert(t('noData'));
    try {
      await api.post('/assets/update', {
        type: activeAsset.value,
        amount: parseFloat(amount.value),
        action,
        transaction_date: new Date(transactionDate.value).toISOString(),
        ayar: activeAsset.value === 'GOLD' ? parseInt(modalAyar.value) : 0
      });
      showModal.value = false;
      fetchData();
      fetchTransactions();
      amount.value = '';
    } catch(e) { alert("Error: " + (e.response?.data?.error || e.message)); }
  };

  const changeCurrency = (code, ayar, label) => {
    baseCurrency.value = code;
    targetAyar.value = ayar;
    baseCurrencyLabel.value = label;
    showSelector.value = false;
    fetchData();
  };

  const downloadReceipt = async () => {
    try {
      const res = await api.get('/assets/receipt/full', { responseType: 'blob', headers: { 'X-Currency': baseCurrency.value, 'X-Ayar': targetAyar.value.toString() } });
      const url = window.URL.createObjectURL(new Blob([res.data]));
      const link = document.createElement('a');
      link.href = url;
      link.setAttribute('download', 'Genel_Rapor.pdf');
      document.body.appendChild(link);
      link.click();
    } catch (e) { alert('PDF Error'); }
  };

  const downloadExcel = async () => {
    try {
      const res = await api.get('/assets/export/excel', { responseType: 'blob' });
      const url = window.URL.createObjectURL(new Blob([res.data]));
      const link = document.createElement('a');
      link.href = url;
      link.setAttribute('download', `FinTrack_Export_${new Date().toLocaleDateString()}.xlsx`);
      document.body.appendChild(link);
      link.click();
    } catch (e) { alert('Excel Error'); }
  };

  const downloadSingleReceipt = async (id) => {
    try {
      const res = await api.get(`/assets/receipt/${id}`, { responseType: 'blob', headers: { 'X-Currency': baseCurrency.value } });
      const url = window.URL.createObjectURL(new Blob([res.data]));
      const link = document.createElement('a');
      link.href = url;
      link.setAttribute('download', `Islem_Dekontu_${id}.pdf`);
      document.body.appendChild(link);
      link.click();
    } catch (e) { alert('PDF Error'); }
  };

  const formatDate = (dateString) => {
    if (!dateString) return '';
    return new Date(dateString).toLocaleDateString('tr-TR');
  };

  onMounted(() => { fetchData(); fetchTransactions(); });
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

  .top-bar { display: flex; justify-content: space-between; align-items: center; margin-bottom: 40px; flex-wrap: wrap; gap: 20px; }
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
    transition: all 0.3s ease;
    display: flex;
    align-items: center;
    justify-content: space-between;
    min-width: 160px;
    box-shadow: 0 4px 10px rgba(0,0,0,0.2);
  }
  .currency-btn:hover, .currency-btn.active { border-color: var(--accent-color); transform: translateY(-2px); box-shadow: 0 8px 20px rgba(0,0,0,0.3); }
  .arrow-icon { font-size: 0.8rem; transition: 0.3s; }
  .currency-btn.active .arrow-icon { transform: rotate(180deg); }

  .currency-dropdown {
    position: absolute;
    top: calc(100% + 10px);
    inset-inline-end: 0;
    background: var(--sidebar-bg);
    border: 1px solid var(--border-color);
    border-radius: 15px;
    width: 220px;
    box-shadow: 0 10px 30px rgba(0,0,0,0.5);
    z-index: 101;
    padding: 8px;
    overflow: visible;
  }

  .c-item {
    padding: 12px 20px;
    cursor: pointer;
    border-radius: 10px;
    margin-bottom: 4px;
    position: relative;
    color: var(--text-color);
    transition: 0.2s;
    font-weight: 600;
  }
  .c-item:last-child { margin-bottom: 0; }
  .c-item:hover { background: var(--hover-bg); color: var(--accent-color); transform: translateX(4px); }

  /* GOLD SUBMENU STYLE */
  .has-submenu { position: relative; overflow: visible; }
  .gold-trigger { display: flex; align-items: center; justify-content: space-between; width: 100%; }

  .submenu {
    position: absolute;
    top: -8px;
    inset-inline-end: 100%;
    margin-inline-end: -2px; /* NO GAP */
    background: var(--sidebar-bg);
    border: 1px solid var(--border-color);
    border-radius: 15px;
    min-width: 140px;
    padding: 8px;
    box-shadow: 0 10px 30px rgba(0,0,0,0.5);
    opacity: 0;
    visibility: hidden;
    transform: translateX(20px);
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  }

  /* RTL specific fix for transform */
  html[dir="ltr"] .submenu { transform: translateX(-20px); }
  html[dir="rtl"] .submenu { transform: translateX(20px); }

  /* Bridge for gapless hover */
  .has-submenu::after {
    content: '';
    position: absolute;
    top: 0;
    bottom: 0;
    width: 20px;
    inset-inline-end: 100%;
    z-index: 102;
  }

  .has-submenu:hover .submenu {
    opacity: 1;
    visibility: visible;
    transform: translateX(0);
  }

  .sub-item {
    padding: 10px 15px;
    font-size: 0.9rem;
    border-radius: 8px;
    color: var(--text-color);
    cursor: pointer;
    transition: 0.2s;
    opacity: 0;
    animation: fadeSlideIn 0.3s forwards;
    animation-delay: calc(var(--i) * 0.05s);
  }
  .sub-item:hover { background: var(--hover-bg); color: var(--accent-color); padding-inline-start: 20px; }

  @keyframes fadeSlideIn {
    from { opacity: 0; transform: translateY(5px); }
    to { opacity: 1; transform: translateY(0); }
  }

  .dropdown-anim-enter-active, .dropdown-anim-leave-active { transition: all 0.3s ease; }
  .dropdown-anim-enter-from, .dropdown-anim-leave-to { opacity: 0; transform: translateY(-10px); }

  .chart-section { display: flex; flex-direction: column; align-items: center; margin-bottom: 50px; position: relative; width: 100%; }
  .chart-wrapper { width: 320px; height: 320px; position: relative; }
  .center-balance { position: absolute; top: 50%; left: 50%; transform: translate(-50%, -50%); text-align: center; pointer-events: none; }
  .center-balance h3 { font-size: 2rem; margin: 0; font-weight: 800; color: var(--text-color); white-space: nowrap; }
  .center-balance small { color: var(--text-muted); font-size: 1rem; font-weight: bold; }

  .no-data-circle { width: 100%; height: 100%; border-radius: 50%; border: 6px dashed var(--border-color); display: flex; align-items: center; justify-content: center; }
  .no-data-content { text-align: center; color: var(--text-muted); font-weight: bold; }

  .assets-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(240px, 1fr)); gap: 25px; width: 100%; }
  .asset-card {
    padding: 25px;
    border-radius: 20px;
    cursor: pointer;
    transition: all 0.3s ease;
    border: 1px solid var(--border-color);
    display: flex;
    align-items: center;
    gap: 20px;
    position: relative;
    overflow: hidden;
    background: rgba(30, 41, 59, 0.6);
    backdrop-filter: blur(10px);
  }
  .asset-card:hover { transform: translateY(-5px); box-shadow: 0 15px 30px rgba(0,0,0,0.3); border-color: var(--accent-color); }
  .asset-card * { color: white !important; }

  .card-btc { background: linear-gradient(135deg, rgba(26,26,26,0.8), rgba(68,68,68,0.8)); }
  .card-gold { background: linear-gradient(135deg, rgba(218,165,32,0.8), rgba(255,215,0,0.8)); }
  .card-usd { background: linear-gradient(135deg, rgba(5,150,105,0.8), rgba(16,185,129,0.8)); }
  .card-eur { background: linear-gradient(135deg, rgba(93,64,55,0.8), rgba(141,110,99,0.8)); }
  .card-silver { background: linear-gradient(135deg, rgba(117,117,117,0.8), rgba(158,158,158,0.8)); }
  .card-try { background: linear-gradient(135deg, rgba(153,27,27,0.8), rgba(239,68,68,0.8)); }

  .card-icon { font-size: 2.5rem; }
  .card-info { display: flex; flex-direction: column; }
  .card-name { font-size: 0.9rem; opacity: 0.9; text-transform: uppercase; font-weight: bold; letter-spacing: 1px; }
  .card-amount { font-size: 1.5rem; font-weight: 800; margin: 2px 0; }
  .card-val { font-size: 0.85rem; opacity: 0.8; font-weight: 600; }

  .floating-actions { position: fixed; bottom: 30px; inset-inline-end: 30px; display: flex; flex-direction: column; gap: 15px; z-index: 110; }
  .f-btn { width: 60px; height: 60px; border-radius: 50%; border: none; font-size: 24px; cursor: pointer; box-shadow: 0 5px 15px rgba(0,0,0,0.3); transition: 0.3s; color: white; display: flex; align-items: center; justify-content: center; }
  .f-btn:hover { transform: scale(1.1); }
  .excel { background: #10B981; }
  .pdf { background: #3B82F6; }

  .modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.85); display: flex; justify-content: center; align-items: center; z-index: 200; backdrop-filter: blur(5px); }
  .large-modal { width: 800px !important; max-width: 95%; }
  .modal-content { background: var(--card-bg); padding: 35px; border-radius: 20px; border: 1px solid var(--border-color); box-shadow: 0 20px 50px rgba(0,0,0,0.5); }
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
  .add { background: var(--success-color); }
  .sub { background: var(--danger-color); }
  .add:hover, .sub:hover { opacity: 0.9; transform: translateY(-2px); }

  @media (max-width: 768px) {
    .modal-body-split { grid-template-columns: 1fr; }
    .transaction-history { border-inline-start: none; border-top: 1px solid var(--border-color); padding-inline-start: 0; padding-top: 30px; }
    .chart-wrapper { width: 280px; height: 280px; }
    .center-balance h3 { font-size: 1.5rem; }
  }
</style>