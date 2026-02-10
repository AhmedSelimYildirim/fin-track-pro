<template>
  <div class="dashboard-content">
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

    <div class="floating-actions">
      <button class="f-btn excel" @click="downloadExcel" title="Excel Al">📊</button>
      <button class="f-btn pdf" @click="downloadReceipt" title="PDF Al">📄</button>
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
      // BAŞARILI MESAJI KALDIRILDI (Sessiz mod)
    } catch(e) {
      // Sadece hata varsa uyar
      alert("Hata: " + (e.response?.data?.error || e.message));
    }
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

  onMounted(() => { fetchData(); fetchTransactions(); });
</script>

<style scoped>
  .dashboard-content { width: 100%; height: 100%; overflow-y: auto; padding: 30px; box-sizing: border-box; }
  .top-bar { display: flex; justify-content: space-between; align-items: center; margin-bottom: 40px; flex-wrap: wrap; gap: 20px; }
  .page-title h2 { color: var(--text-color); margin: 0; font-size: 1.8rem; }
  .currency-wrapper { position: relative; z-index: 100; }
  .currency-btn { background: var(--sidebar-bg); border: 1px solid var(--border-color); padding: 12px 25px; border-radius: 25px; cursor: pointer; font-weight: bold; color: var(--accent-color); transition: 0.3s; display: flex; align-items: center; justify-content: center; min-width: 160px; }
  .currency-btn:hover { border-color: var(--accent-color); }

  .currency-dropdown { position: absolute; top: 110%; right: 0; background: var(--sidebar-bg); border: 1px solid var(--border-color); border-radius: 15px; width: 220px; box-shadow: 0 10px 30px rgba(0,0,0,0.5); z-index: 101; }
  .c-item { padding: 12px 20px; cursor: pointer; border-bottom: 1px solid var(--border-color); position: relative; color: var(--text-color); transition: 0.2s; }
  .c-item:hover { background: var(--hover-bg); color: var(--accent-color); }
  .has-submenu:hover .submenu { display: block; }
  .submenu { display: none; position: absolute; top: 0; right: 100%; background: var(--sidebar-bg); border: 1px solid var(--border-color); border-radius: 15px; width: 180px; box-shadow: 0 10px 30px rgba(0,0,0,0.5); margin-right: 5px; }
  .submenu div { padding: 12px 20px; cursor: pointer; border-bottom: 1px solid var(--border-color); color: var(--text-color); }
  .submenu div:hover { background: var(--hover-bg); color: var(--accent-color); }

  .chart-section { display: flex; flex-direction: column; align-items: center; margin-bottom: 50px; position: relative; width: 100%; }
  .chart-wrapper { width: 300px; height: 300px; position: relative; }
  .center-balance { position: absolute; top: 50%; left: 50%; transform: translate(-50%, -50%); text-align: center; pointer-events: none; }
  .center-balance h3 { font-size: 1.8rem; margin: 0; font-weight: 800; color: var(--text-color); white-space: nowrap; }
  .center-balance small { color: var(--text-muted); }
  .total-underline { width: 150px; height: 4px; background: linear-gradient(90deg, transparent, var(--accent-color), transparent); margin-top: 20px; border-radius: 2px; }

  .no-data-circle { width: 100%; height: 100%; border-radius: 50%; border: 8px dashed var(--border-color); display: flex; align-items: center; justify-content: center; }
  .no-data-content { text-align: center; color: var(--text-muted); font-weight: bold; }

  .assets-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(220px, 1fr)); gap: 25px; width: 100%; }
  .asset-card { padding: 25px; border-radius: 20px; cursor: pointer; transition: all 0.3s ease; border: 1px solid var(--border-color); display: flex; align-items: center; gap: 15px; position: relative; overflow: hidden; }
  .asset-card:hover { transform: translateY(-5px); box-shadow: 0 15px 30px rgba(0,0,0,0.3); }
  .asset-card * { color: white !important; }

  .card-btc { background: linear-gradient(135deg, #1a1a1a, #444); }
  .card-gold { background: linear-gradient(135deg, #DAA520, #FFD700); }
  .card-usd { background: linear-gradient(135deg, #059669, #10B981); }
  .card-eur { background: linear-gradient(135deg, #5D4037, #8D6E63); }
  .card-silver { background: linear-gradient(135deg, #757575, #9E9E9E); }
  .card-try { background: linear-gradient(135deg, #991B1B, #EF4444); }

  .card-icon { font-size: 2.2rem; }
  .card-info { display: flex; flex-direction: column; }
  .card-name { font-size: 0.9rem; opacity: 0.9; text-transform: uppercase; font-weight: bold; letter-spacing: 1px; }
  .card-amount { font-size: 1.4rem; font-weight: bold; margin: 2px 0; }
  .card-val { font-size: 0.85rem; opacity: 0.8; }

  .floating-actions { position: fixed; bottom: 30px; right: 30px; display: flex; flex-direction: column; gap: 15px; z-index: 110; }
  .f-btn { width: 60px; height: 60px; border-radius: 50%; border: none; font-size: 24px; cursor: pointer; box-shadow: 0 5px 15px rgba(0,0,0,0.3); transition: 0.3s; color: white; display: flex; align-items: center; justify-content: center; }
  .f-btn:hover { transform: scale(1.1); }
  .excel { background: #10B981; }
  .pdf { background: #3B82F6; }

  .modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.85); display: flex; justify-content: center; align-items: center; z-index: 200; backdrop-filter: blur(5px); }
  .large-modal { width: 750px !important; max-width: 95%; }
  .modal-content { background: var(--card-bg); padding: 30px; border-radius: 20px; border: 1px solid var(--border-color); box-shadow: 0 20px 50px rgba(0,0,0,0.5); }
  .modal-header { display: flex; justify-content: space-between; margin-bottom: 20px; border-bottom: 1px solid var(--border-color); padding-bottom: 10px; }
  .modal-header h3 { color: var(--text-color); margin: 0; }
  .modal-header button { background: none; border: none; color: var(--text-color); font-size: 1.5rem; cursor: pointer; }

  .modal-body-split { display: grid; grid-template-columns: 1fr 1fr; gap: 30px; margin-top: 20px; }
  .transaction-form, .transaction-history { display: flex; flex-direction: column; gap: 10px; }
  .transaction-history { border-left: 1px solid var(--border-color); padding-left: 20px; }
  h4 { color: var(--text-muted); margin: 0 0 10px 0; }
  .history-list { max-height: 300px; overflow-y: auto; padding-right: 5px; }
  .history-item { background: var(--input-bg); padding: 12px; border-radius: 10px; margin-bottom: 10px; display: flex; justify-content: space-between; align-items: center; border: 1px solid var(--border-color); }
  .tx-info { display: flex; flex-direction: column; gap: 4px; }
  .tx-date { color: var(--text-muted); font-size: 0.8rem; }
  .tx-type.add { color: var(--success-color) !important; font-weight: bold; }
  .tx-type.subtract { color: var(--danger-color) !important; font-weight: bold; }
  .tx-amount { color: var(--text-color); font-weight: bold; }
  .receipt-download-btn { background: transparent; border: 1px solid var(--border-color); border-radius: 5px; cursor: pointer; padding: 5px; transition: 0.2s; font-size: 1rem; }
  .receipt-download-btn:hover { background: var(--hover-bg); }

  .big-input, .ayar-select, .date-input { width: 100%; padding: 12px; background: var(--input-bg); border: 1px solid var(--border-color); color: var(--text-color); border-radius: 10px; margin-bottom: 15px; font-size: 1rem; box-sizing: border-box; }
  .actions { display: flex; gap: 10px; }
  .actions button { flex: 1; padding: 15px; border-radius: 10px; border: none; font-weight: bold; cursor: pointer; color: white; transition: 0.2s; }
  .add { background: var(--success-color); }
  .sub { background: var(--danger-color); }
  .add:hover, .sub:hover { opacity: 0.9; }

  @media (max-width: 768px) {
    .modal-body-split { grid-template-columns: 1fr; }
    .transaction-history { border-left: none; border-top: 1px solid var(--border-color); padding-left: 0; padding-top: 20px; }
  }
</style>