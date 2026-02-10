<template>
  <div class="dashboard-page">
    <aside class="sidebar">
      <div class="brand-container">
        <div class="brand">FinTrack Pro</div>
        <div class="user-badge">{{ currentUser }}</div>
      </div>

      <nav class="menu">
        <div class="menu-item" :class="{ active: currentTab === 'dashboard' }" @click="currentTab = 'dashboard'">
          <span>📊</span> {{ t('home') }}
        </div>
        <div class="menu-item" :class="{ active: currentTab === 'calendar' }" @click="router.push('/calendar')">
          <span>📅</span> {{ t('calendar') }}
        </div>
        <div class="menu-item" :class="{ active: currentTab === 'settings' }" @click="currentTab = 'settings'">
          <span>⚙️</span> {{ t('settings') }}
        </div>
      </nav>

      <div class="logout-wrapper">
        <div class="menu-item logout" @click="logout">
          <span>🚪</span> {{ t('logout') }}
        </div>
      </div>
    </aside>

    <main class="content">
      <SettingsView v-if="currentTab === 'settings'" />

      <div v-else>
        <header class="top-bar">
          <div class="page-title">
            <h2>{{ t('portfolioSummary') }}</h2>
          </div>
          <div class="currency-wrapper">
            <div class="currency-btn" @click.stop="toggleDropdown">
              {{ displayCurrency }} ▼
            </div>
            <div v-if="showSelector" class="currency-dropdown">
              <div class="c-item" @click="changeCurrency('TRY', 0, 'Türk Lirası (TL)')">Türk Lirası (TL)</div>
              <div class="c-item" @click="changeCurrency('USD', 0, 'Dolar ($)')">Dolar ($)</div>
              <div class="c-item" @click="changeCurrency('EUR', 0, 'Euro (€)')">Euro (€)</div>
              <div class="c-item" @click="changeCurrency('BTC', 0, 'Bitcoin')">Bitcoin</div>
              <div class="c-item" @click="changeCurrency('SILVER', 0, 'Gümüş (Gram)')">Gümüş (Gram)</div>
              <div class="c-item has-submenu">
                {{ t('addAsset') }} (Altın) ▶
                <div class="submenu">
                  <div @click="changeCurrency('GOLD', 24, '24 Ayar (Has)')">24 Ayar</div>
                  <div @click="changeCurrency('GOLD', 22, '22 Ayar')">22 Ayar</div>
                  <div @click="changeCurrency('GOLD', 18, '18 Ayar')">18 Ayar</div>
                  <div @click="changeCurrency('GOLD', 14, '14 Ayar')">14 Ayar</div>
                </div>
              </div>
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
          <div class="total-underline"></div>
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
      </div>
    </main>

    <div v-if="showModal" class="modal-overlay" @click.self="showModal = false">
      <div class="modal-content large-modal">
        <div class="modal-header">
          <h3>{{ activeAsset }} {{ t('newTransaction') }}</h3>
          <button @click="showModal = false">✕</button>
        </div>

        <div class="modal-body-split">
          <div class="transaction-form">
            <h4>{{ t('newTransaction') }}</h4>
            <input v-model="amount" type="number" placeholder="Miktar" class="big-input" />
            <div v-if="activeAsset === 'GOLD'" class="ayar-wrapper">
              <select v-model="modalAyar" class="ayar-select">
                <option :value="24">24 Ayar (Has)</option>
                <option :value="22">22 Ayar</option>
                <option :value="18">18 Ayar</option>
                <option :value="14">14 Ayar</option>
              </select>
            </div>
            <input v-model="transactionDate" type="date" class="date-input" />
            <div class="actions">
              <button class="add" @click="handleTransaction('add')">EKLE (+)</button>
              <button class="sub" @click="handleTransaction('subtract')">ÇIKAR (-)</button>
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

    <div class="floating-actions" v-if="currentTab === 'dashboard'">
      <button class="f-btn excel" @click="downloadExcel" title="Excel Al">📊</button>
      <button class="f-btn pdf" @click="downloadReceipt" title="PDF Al">📄</button>
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
  import SettingsView from './SettingsView.vue';

  ChartJS.register(ArcElement, Tooltip, Legend);
  const router = useRouter();

  const currentUser = ref(localStorage.getItem('username') || 'Yatırımcı');
  const currentTab = ref('dashboard');
  const showSelector = ref(false);
  const showModal = ref(false);
  const activeAsset = ref('');
  const amount = ref('');
  const transactionDate = ref(new Date().toISOString().split('T')[0]);
  const modalAyar = ref(24);
  const summaryData = ref(null);
  const allTransactions = ref([]);
  const baseCurrency = ref('TRY');
  const baseCurrencyLabel = ref('Türk Lirası (TL)');
  const targetAyar = ref(0);

  const cardConfigs = [
    { type: 'BTC', label: 'Bitcoin', icon: '₿', unit: 'BTC' },
    { type: 'GOLD', label: 'Altın', icon: '👑', unit: 'Gr' },
    { type: 'USD', label: 'USD', icon: '$', unit: '$' },
    { type: 'EUR', label: 'EURO', icon: '€', unit: '€' },
    { type: 'SILVER', label: 'Gümüş', icon: '⚔️', unit: 'Gr' },
    { type: 'TRY', label: 'TL Nakit', icon: '₺', unit: '₺' }
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
    const labels = ['BTC', 'GOLD', 'USD', 'EUR', 'SILVER', 'TRY'];
    const colors = ['#1a1a1a', '#FFD700', '#10B981', '#8B4513', '#A0A0A0', '#EF4444'];
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
    if (!amount.value) return alert("Miktar girin");
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
      alert("İşlem Başarılı! ✅");
    } catch(e) { alert("Hata: " + (e.response?.data?.error || e.message)); }
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
    } catch (e) { alert('Rapor hatası'); }
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
    } catch (e) { alert('Excel indirilemedi'); }
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
    } catch (e) { alert('Dekont indirilemedi'); }
  };

  const formatDate = (dateString) => {
    if (!dateString) return '';
    return new Date(dateString).toLocaleDateString('tr-TR');
  };

  const logout = () => { localStorage.clear(); router.push('/login'); };

  onMounted(() => { fetchData(); fetchTransactions(); });
</script>

<style scoped>
  .dashboard-page { display: flex; min-height: 100vh; background: var(--bg-color); color: var(--text-color); }
  .sidebar { width: 260px; background: var(--sidebar-bg); display: flex; flex-direction: column; padding: 25px; border-right: 1px solid var(--border-color); position: fixed; height: 100vh; z-index: 100; }
  .content { flex: 1; margin-left: 260px; padding: 40px; }

  .brand { color: var(--accent-color); font-size: 1.6rem; font-weight: 800; text-align: center; margin-bottom: 40px; }
  .user-badge { color: var(--success-color); font-weight: bold; text-align: center; margin-bottom: 20px; border-bottom: 1px solid var(--border-color); padding-bottom: 10px; }

  .menu-item { padding: 15px; margin-bottom: 10px; border-radius: 12px; cursor: pointer; color: var(--text-muted); display: flex; gap: 12px; transition: 0.3s; }
  .menu-item:hover, .menu-item.active { background: var(--hover-bg); color: var(--text-color); }

  .logout-wrapper { margin-top: auto; padding-top: 20px; border-top: 1px solid var(--border-color); }
  .logout { border: 2px solid var(--danger-color); border-radius: 12px; color: var(--danger-color); justify-content: center; font-weight: bold; padding: 12px; }
  .logout:hover { background: var(--danger-color); color: white; }

  .top-bar { display: flex; justify-content: space-between; align-items: center; margin-bottom: 40px; }
  .currency-btn { background: var(--card-bg); border: 1px solid var(--border-color); padding: 12px 25px; border-radius: 25px; cursor: pointer; font-weight: bold; color: var(--accent-color); }

  .chart-section { margin-bottom: 50px; display: flex; flex-direction: column; align-items: center; }
  .chart-wrapper { width: 300px; height: 300px; position: relative; }
  .center-balance { position: absolute; top: 50%; left: 50%; transform: translate(-50%, -50%); text-align: center; }
  .center-balance h3 { font-size: 1.8rem; margin: 0; }

  .assets-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 20px; }
  .asset-card { padding: 20px; border-radius: 20px; cursor: pointer; transition: 0.3s; border: 1px solid var(--border-color); }
  .asset-card:hover { transform: translateY(-5px); box-shadow: 0 10px 20px rgba(0,0,0,0.2); }
  .asset-card * { color: white !important; }

  .card-btc { background: linear-gradient(135deg, #1a1a1a, #444); }
  .card-gold { background: linear-gradient(135deg, #DAA520, #FFD700); }
  .card-usd { background: linear-gradient(135deg, #059669, #10B981); }
  .card-eur { background: linear-gradient(135deg, #5D4037, #8D6E63); }
  .card-silver { background: linear-gradient(135deg, #757575, #9E9E9E); }
  .card-try { background: linear-gradient(135deg, #991B1B, #EF4444); }

  .floating-actions { position: fixed; bottom: 30px; right: 30px; display: flex; flex-direction: column; gap: 15px; }
  .f-btn { width: 60px; height: 60px; border-radius: 50%; border: none; font-size: 24px; cursor: pointer; box-shadow: 0 5px 15px rgba(0,0,0,0.3); transition: 0.3s; }
  .f-btn:hover { transform: scale(1.1); }
  .excel { background: #10B981; color: white; }
  .pdf { background: #3B82F6; color: white; }

  .modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.85); display: flex; justify-content: center; align-items: center; z-index: 200; backdrop-filter: blur(5px); }
  .large-modal { width: 750px !important; max-width: 95%; }
  .modal-content { background: var(--card-bg); padding: 30px; border-radius: 20px; border: 1px solid var(--border-color); }
  .modal-body-split { display: grid; grid-template-columns: 1fr 1fr; gap: 30px; margin-top: 20px; }
  .transaction-history { border-left: 1px solid var(--border-color); padding-left: 20px; }
  .history-list { max-height: 300px; overflow-y: auto; }
  .history-item { background: var(--input-bg); padding: 12px; border-radius: 10px; margin-bottom: 10px; display: flex; justify-content: space-between; align-items: center; }
  .tx-type.add { color: var(--success-color); }
  .tx-type.subtract { color: var(--danger-color); }

  .big-input, .ayar-select, .date-input { width: 100%; padding: 12px; background: var(--input-bg); border: 1px solid var(--border-color); color: var(--text-color); border-radius: 10px; margin-bottom: 15px; }
  .actions { display: flex; gap: 10px; }
  .actions button { flex: 1; padding: 15px; border-radius: 10px; border: none; font-weight: bold; cursor: pointer; color: white; }
  .add { background: var(--success-color); }
  .sub { background: var(--danger-color); }

  @media (max-width: 768px) {
    .sidebar { display: none; }
    .content { margin-left: 0; }
    .modal-body-split { grid-template-columns: 1fr; }
    .transaction-history { border-left: none; border-top: 1px solid var(--border-color); padding-left: 0; padding-top: 20px; }
  }
</style>