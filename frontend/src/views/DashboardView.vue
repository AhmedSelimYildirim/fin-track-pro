<template>
  <div class="dashboard-page">
    <aside class="sidebar">
      <div class="brand">FinTrack Pro 🚀</div>

      <nav class="menu">
        <div class="menu-item active">
          <span>📊</span> Dashboard
        </div>
        <div class="menu-item" @click="router.push('/calendar')">
          <span>📅</span> Takvim & Notlar
        </div>
      </nav>

      <div class="logout-wrapper">
        <div class="menu-item logout" @click="logout">
          <span>🚪</span> Çıkış Yap
        </div>
      </div>
    </aside>

    <main class="content">
      <header class="top-bar">
        <div class="user-info">
          <h2>Portföy Özeti</h2>
          <span class="subtitle">Hoşgeldin, {{ currentUser }}</span>
        </div>

        <div class="currency-selector">
          <button @click="showSelector = !showSelector">
            {{ displayCurrency }} ▼
          </button>
          <div v-if="showSelector" class="dropdown">
            <div v-for="opt in currencyOptions" :key="opt.label" @click="changeCurrency(opt)">
              {{ opt.label }}
            </div>
          </div>
        </div>
      </header>

      <div class="chart-section">
        <div class="chart-wrapper">
          <Doughnut v-if="chartData" :data="chartData" :options="chartOptions" />
          <div class="center-balance">
            <h3>{{ totalValue }}</h3>
            <small>{{ baseCurrency }}</small>
          </div>
        </div>
      </div>

      <div class="assets-grid">
        <div class="asset-card btc" @click="openModal('BTC')">
          <div class="icon">₿</div>
          <div>
            <div class="name">Bitcoin</div>
            <div class="val">%{{ getAllocation('BTC') }}</div>
          </div>
        </div>
        <div class="asset-card gold" @click="openModal('GOLD')">
          <div class="icon">👑</div>
          <div>
            <div class="name">Gold</div>
            <div class="val">%{{ getAllocation('GOLD') }}</div>
          </div>
        </div>
        <div class="asset-card usd" @click="openModal('USD')">
          <div class="icon">$</div>
          <div>
            <div class="name">USD</div>
            <div class="val">%{{ getAllocation('USD') }}</div>
          </div>
        </div>
        <div class="asset-card eur" @click="openModal('EUR')">
          <div class="icon">€</div>
          <div>
            <div class="name">Euro</div>
            <div class="val">%{{ getAllocation('EUR') }}</div>
          </div>
        </div>
        <div class="asset-card silver" @click="openModal('SILVER')">
          <div class="icon">⚔️</div>
          <div>
            <div class="name">Silver</div>
            <div class="val">%{{ getAllocation('SILVER') }}</div>
          </div>
        </div>
        <div class="asset-card try" @click="openModal('TRY')">
          <div class="icon">₺</div>
          <div>
            <div class="name">TL Nakit</div>
            <div class="val">%{{ getAllocation('TRY') }}</div>
          </div>
        </div>
      </div>
    </main>

    <div v-if="showModal" class="modal-overlay">
      <div class="modal-content">
        <div class="modal-header">
          <h3>{{ activeAsset }} İşlemi</h3>
          <button @click="showModal = false">✕</button>
        </div>
        <div class="modal-body">
          <input v-model="amount" type="number" placeholder="Miktar" class="big-input" />
          <input v-model="transactionDate" type="date" class="date-input" />
          <div class="actions">
            <button class="add" @click="handleTransaction('add')">EKLE (+)</button>
            <button class="sub" @click="handleTransaction('subtract')">ÇIKAR (-)</button>
          </div>
        </div>
      </div>
    </div>

    <button class="receipt-btn" @click="downloadReceipt" title="Genel Rapor Al">📄</button>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { Chart as ChartJS, ArcElement, Tooltip, Legend } from 'chart.js';
import { Doughnut } from 'vue-chartjs';
import api from '../services/api';

ChartJS.register(ArcElement, Tooltip, Legend);
const router = useRouter();
const currentUser = localStorage.getItem('username') || 'Yatırımcı';

const showSelector = ref(false);
const showModal = ref(false);
const activeAsset = ref('');
const amount = ref('');
const transactionDate = ref(new Date().toISOString().split('T')[0]);
const summaryData = ref(null);
const baseCurrency = ref('TRY');
const targetAyar = ref(0);

const currencyOptions = [
  { label: 'Türk Lirası (TL)', val: 'TRY', ayar: 0 },
  { label: 'Dolar ($)', val: 'USD', ayar: 0 },
  { label: 'Euro (€)', val: 'EUR', ayar: 0 },
  { label: 'Bitcoin', val: 'BTC', ayar: 0 },
  { label: 'Altın (Has)', val: 'GOLD', ayar: 24 },
];

const displayCurrency = computed(() => {
  const found = currencyOptions.find(c => c.val === baseCurrency.value && c.ayar === targetAyar.value);
  return found ? found.label : 'Para Birimi';
});

const totalValue = computed(() => summaryData.value ? summaryData.value.total_value.toLocaleString('tr-TR', { maximumFractionDigits: 2 }) : '0');

const chartData = computed(() => {
  if (!summaryData.value) return null;
  const labels = ['BTC', 'GOLD', 'USD', 'EUR', 'SILVER', 'TRY'];
  const colors = ['#F7931A', '#FFD700', '#10B981', '#3B82F6', '#94A3B8', '#EF4444'];
  const data = labels.map(t => summaryData.value.assets.find(a => a.type === t)?.allocation || 0);
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

const getAllocation = (type) => summaryData.value?.assets.find(a => a.type === type)?.allocation.toFixed(1) || '0.0';

const openModal = (asset) => {
  activeAsset.value = asset;
  amount.value = '';
  transactionDate.value = new Date().toISOString().split('T')[0];
  showModal.value = true;
};

const handleTransaction = async (action) => {
  if (!amount.value) return alert("Miktar girin");
  try {
    await api.post('/assets/balance', {
      type: activeAsset.value,
      amount: parseFloat(amount.value),
      action,
      transaction_date: new Date(transactionDate.value).toISOString(),
      ayar: activeAsset.value === 'GOLD' ? 24 : 0
    });
    showModal.value = false;
    fetchData();
  } catch(e) { alert(e.message); }
};

const changeCurrency = (opt) => {
  baseCurrency.value = opt.val;
  targetAyar.value = opt.ayar;
  showSelector.value = false;
  fetchData();
};

const downloadReceipt = async () => {
  try {
    const res = await api.get('/assets/receipt/full', {
        responseType: 'blob',
        headers: { 'X-Currency': baseCurrency.value, 'X-Ayar': targetAyar.value.toString() }
    });
    const url = window.URL.createObjectURL(new Blob([res.data]));
    const link = document.createElement('a');
    link.href = url;
    link.setAttribute('download', 'Genel_Rapor.pdf');
    document.body.appendChild(link);
    link.click();
  } catch (e) { alert('Rapor hatası'); }
};

const logout = () => { localStorage.clear(); router.push('/login'); };

onMounted(fetchData);
</script>

<style scoped>
.dashboard-page { display: flex; min-height: 100vh; background: #0F172A; color: white; font-family: sans-serif; }
.sidebar { width: 250px; background: #1E293B; display: flex; flex-direction: column; padding: 20px; border-right: 1px solid rgba(255,255,255,0.05); }
.brand { color: #FFD700; font-size: 1.5rem; font-weight: bold; margin-bottom: 40px; }
.menu-item { padding: 15px; margin-bottom: 10px; border-radius: 10px; cursor: pointer; color: #94A3B8; display: flex; gap: 10px; align-items: center; transition: 0.2s; }
.menu-item:hover, .menu-item.active { background: #334155; color: white; }
.logout-wrapper { margin-top: auto; }
.logout { color: #EF4444; } .logout:hover { background: rgba(239, 68, 68, 0.1); }

.content { flex: 1; padding: 30px; }
.top-bar { display: flex; justify-content: space-between; align-items: center; margin-bottom: 30px; }
.subtitle { color: #94A3B8; font-size: 0.9rem; }
.currency-selector button { background: #334155; color: white; border: none; padding: 10px 20px; border-radius: 20px; cursor: pointer; }
.dropdown { position: absolute; background: #1E293B; border: 1px solid #475569; padding: 10px; border-radius: 10px; z-index: 100; }

.chart-section { display: flex; justify-content: center; margin-bottom: 40px; }
.chart-wrapper { width: 250px; height: 250px; position: relative; }
.center-balance { position: absolute; top: 50%; left: 50%; transform: translate(-50%, -50%); text-align: center; }

.assets-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(150px, 1fr)); gap: 15px; }
.asset-card { background: #1E293B; padding: 20px; border-radius: 15px; cursor: pointer; display: flex; flex-direction: column; gap: 10px; border: 1px solid rgba(255,255,255,0.05); transition: 0.2s; }
.asset-card:hover { transform: translateY(-5px); }
.btc { border-top: 3px solid #F7931A; } .gold { border-top: 3px solid #FFD700; }
.usd { border-top: 3px solid #10B981; } .eur { border-top: 3px solid #3B82F6; }
.silver { border-top: 3px solid #94A3B8; } .try { border-top: 3px solid #EF4444; }
.icon { font-size: 1.5rem; } .val { font-size: 1.2rem; font-weight: bold; }

.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.8); display: flex; justify-content: center; align-items: center; z-index: 200; }
.modal-content { background: #1E293B; padding: 25px; border-radius: 20px; width: 300px; border: 1px solid #334155; }
.modal-header { display: flex; justify-content: space-between; margin-bottom: 20px; }
.big-input { width: 100%; padding: 10px; background: #0F172A; border: 1px solid #334155; color: white; margin-bottom: 10px; font-size: 1.2rem; text-align: center; }
.date-input { width: 100%; padding: 10px; background: #0F172A; border: 1px solid #334155; color: white; margin-bottom: 15px; color-scheme: dark; }
.actions { display: flex; gap: 10px; }
.actions button { flex: 1; padding: 10px; border: none; border-radius: 8px; color: white; font-weight: bold; cursor: pointer; }
.add { background: #10B981; } .sub { background: #EF4444; }

.receipt-btn { position: fixed; bottom: 30px; right: 30px; background: #3B82F6; width: 60px; height: 60px; border-radius: 50%; border: none; font-size: 24px; cursor: pointer; box-shadow: 0 5px 15px rgba(0,0,0,0.5); }
</style>