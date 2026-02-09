<template>
  <div class="dashboard">
    <header class="top-bar">
      <div class="left-group">
        <button class="home-btn" @click="toggleSidebar">
          <component :is="Home" />
        </button>
        <div class="user-info">
          <span class="u-name">{{ currentUser }}</span>
          <span class="u-email">{{ currentEmail }}</span>
        </div>
      </div>

      <div class="currency-selector">
        <button class="select-trigger" @click="showSelector = !showSelector">
          {{ displayCurrency }} ▼
        </button>

        <div v-if="showSelector" class="scroll-dropdown">
          <div v-for="opt in currencyOptions" :key="opt.label"
               class="dropdown-item" @click="changeCurrency(opt)">
            {{ opt.label }}
          </div>
        </div>
      </div>
    </header>

    <transition name="slide">
      <div v-if="isSidebarOpen" class="sidebar">
        <div class="menu-item"><component :is="Calendar" size="18"/> Add Note</div>
        <div class="menu-item"><component :is="User" size="18"/> Update Profile</div>
        <div class="menu-item red"><component :is="Trash2" size="18"/> Delete Account</div>
        <div class="menu-item red" @click="logout"><component :is="LogOut" size="18"/> Logout</div>
      </div>
    </transition>

    <main class="cockpit">
      <div class="chart-container">
        <Doughnut v-if="chartData" :data="chartData" :options="chartOptions" />

        <div class="center-balance">
          <h2>{{ totalValue }}</h2>
          <small>{{ baseCurrency === 'GOLD' ? 'Gram' : baseCurrency }}</small>
        </div>
      </div>

      <div class="asset-buttons">
        <button class="asset-btn btc-btn" @click="openModal('BTC')">BTC</button>
        <button class="asset-btn gold-btn" @click="openModal('GOLD')">GOLD</button>
        <button class="asset-btn usd-btn" @click="openModal('USD')">USD</button>
        <button class="asset-btn eur-btn" @click="openModal('EUR')">EUR</button>
        <button class="asset-btn silver-btn" @click="openModal('SILVER')">SILVER</button>
        <button class="asset-btn try-btn" @click="openModal('TRY')">TRY</button>
      </div>
    </main>

    <button class="receipt-btn" @click="downloadReceipt">
      <component :is="FileText" size="24" />
    </button>

    <div v-if="showModal" class="modal-overlay">
      <div class="modal-content">
        <h3>Manage {{ activeAsset }}</h3>
        <input v-model="amount" type="number" placeholder="Enter Amount" class="modal-input" />
        <div class="modal-actions">
          <button class="action-btn add" @click="handleTransaction('add')">ADD (+)</button>
          <button class="action-btn sub" @click="handleTransaction('subtract')">REMOVE (-)</button>
        </div>
        <button class="close-modal" @click="showModal = false">Cancel</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { Chart as ChartJS, ArcElement, Tooltip, Legend } from 'chart.js';
import { Doughnut } from 'vue-chartjs';
import { Home, Calendar, User, Trash2, LogOut, FileText } from 'lucide-vue-next';
import api from '../services/api';

ChartJS.register(ArcElement, Tooltip, Legend);

const router = useRouter();
const currentUser = localStorage.getItem('username') || 'User';
const currentEmail = localStorage.getItem('email') || 'user@mail.com';

const isSidebarOpen = ref(false);
const showSelector = ref(false);
const showModal = ref(false);
const activeAsset = ref('');
const amount = ref('');
const summaryData = ref(null);

const baseCurrency = ref('TRY');
const targetAyar = ref(0);

const currencyOptions = [
  { label: 'Turkish Lira (TRY)', val: 'TRY', ayar: 0 },
  { label: 'US Dollar (USD)', val: 'USD', ayar: 0 },
  { label: 'Euro (EUR)', val: 'EUR', ayar: 0 },
  { label: 'Bitcoin (BTC)', val: 'BTC', ayar: 0 },
  { label: 'Gold (Standard)', val: 'GOLD', ayar: 0 },
  { label: 'Gold (23 Karat)', val: 'GOLD', ayar: 23 },
  { label: 'Gold (22 Karat)', val: 'GOLD', ayar: 22 },
  { label: 'Gold (21 Karat)', val: 'GOLD', ayar: 21 },
  { label: 'Gold (18 Karat)', val: 'GOLD', ayar: 18 },
  { label: 'Gold (14 Karat)', val: 'GOLD', ayar: 14 },
];

const displayCurrency = computed(() => {
  const found = currencyOptions.find(c => c.val === baseCurrency.value && c.ayar === targetAyar.value);
  return found ? found.label : 'Select Currency';
});

const totalValue = computed(() => {
  if (!summaryData.value) return '0.00';
  return summaryData.value.total_value.toLocaleString('tr-TR', { maximumFractionDigits: 2 });
});

const chartData = computed(() => {
  if (!summaryData.value) return null;
  return {
    labels: summaryData.value.assets.map(a => a.type),
    datasets: [{
      backgroundColor: ['#1A1A1A', '#FFD700', '#10B981', '#795548', '#94A3B8', '#EF4444'],
      data: summaryData.value.assets.map(a => a.allocation),
      borderWidth: 0
    }]
  };
});

const chartOptions = {
  responsive: true,
  cutout: '70%',
  plugins: { legend: { display: false } }
};

const fetchData = async () => {
  try {
    const res = await api.get('/assets/summary', {
      headers: { 'X-Currency': baseCurrency.value, 'X-Ayar': targetAyar.value.toString() }
    });
    summaryData.value = res.data;
  } catch (e) {
    console.error(e);
  }
};

const changeCurrency = (opt) => {
  baseCurrency.value = opt.val;
  targetAyar.value = opt.ayar;
  showSelector.value = false;
  fetchData();
};

const toggleSidebar = () => {
  if (showModal.value) {
    showModal.value = false;
  } else {
    isSidebarOpen.value = !isSidebarOpen.value;
  }
};

const openModal = (asset) => {
  activeAsset.value = asset;
  showModal.value = true;
  isSidebarOpen.value = false;
};

const handleTransaction = async (action) => {
  try {
    await api.post('/assets/balance', {
      type: activeAsset.value,
      amount: parseFloat(amount.value),
      action: action,
      ayar: activeAsset.value === 'GOLD' ? 24 : 0
    });
    amount.value = '';
    showModal.value = false;
    fetchData();
  } catch (e) {
    alert('Error: ' + e.message);
  }
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
    link.setAttribute('download', 'Portfolio_Receipt.pdf');
    document.body.appendChild(link);
    link.click();
  } catch (e) {
    alert('Download failed');
  }
};

const logout = () => {
  localStorage.clear();
  router.push('/login');
};

onMounted(() => {
  fetchData();
});
</script>

<style scoped>
.dashboard { min-height: 100vh; position: relative; padding: 20px; color: white; }
.top-bar { display: flex; justify-content: space-between; align-items: center; margin-bottom: 30px; }
.left-group { display: flex; align-items: center; gap: 15px; }
.home-btn { background: #334155; border: none; color: white; padding: 10px; border-radius: 50%; display: flex; align-items: center; justify-content: center; }
.u-name { display: block; font-weight: bold; font-size: 1.1rem; }
.u-email { font-size: 0.8rem; color: #94A3B8; }
.currency-selector { position: relative; }
.select-trigger { background: #334155; color: white; border: none; padding: 10px 20px; border-radius: 8px; min-width: 160px; text-align: left; }
.scroll-dropdown { position: absolute; top: 100%; right: 0; background: var(--card-bg); border: 1px solid #475569; border-radius: 10px; width: 200px; max-height: 200px; overflow-y: auto; z-index: 100; margin-top: 5px; box-shadow: 0 10px 30px rgba(0,0,0,0.5); }
.dropdown-item { padding: 12px; cursor: pointer; border-bottom: 1px solid #334155; font-size: 0.9rem; }
.dropdown-item:hover { background: #334155; }
.cockpit { display: flex; flex-direction: column; align-items: center; margin-top: 20px; }
.chart-container { position: relative; width: 280px; height: 280px; margin-bottom: 40px; }
.center-balance { position: absolute; top: 50%; left: 50%; transform: translate(-50%, -50%); text-align: center; pointer-events: none; }
.center-balance h2 { margin: 0; font-size: 1.5rem; }
.center-balance small { color: #94A3B8; }
.asset-buttons { display: flex; flex-wrap: wrap; justify-content: center; gap: 15px; max-width: 600px; }
.asset-btn { border: none; border-radius: 50px; padding: 15px 30px; font-weight: bold; color: white; font-size: 1rem; box-shadow: 0 4px 10px rgba(0,0,0,0.3); }
.btc-btn { background: var(--btc-color); }
.gold-btn { background: var(--gold-color); color: black; }
.usd-btn { background: var(--usd-color); }
.eur-btn { background: var(--eur-color); }
.silver-btn { background: var(--silver-color); color: black; }
.try-btn { background: var(--try-color); }
.sidebar { position: absolute; top: 80px; left: 20px; background: var(--card-bg); width: 220px; padding: 15px; border-radius: 15px; box-shadow: 0 10px 30px rgba(0,0,0,0.5); z-index: 50; }
.menu-item { display: flex; align-items: center; gap: 10px; padding: 12px; cursor: pointer; border-bottom: 1px solid #334155; }
.menu-item:hover { background: #334155; border-radius: 8px; }
.red { color: #EF4444; }
.receipt-btn { position: fixed; bottom: 30px; right: 30px; background: #3B82F6; width: 60px; height: 60px; border-radius: 50%; border: none; color: white; display: flex; align-items: center; justify-content: center; box-shadow: 0 5px 20px rgba(0,0,0,0.5); }
.modal-overlay { position: fixed; top: 0; left: 0; width: 100%; height: 100%; background: rgba(0,0,0,0.8); display: flex; justify-content: center; align-items: center; z-index: 200; }
.modal-content { background: var(--card-bg); padding: 30px; border-radius: 20px; width: 300px; text-align: center; }
.modal-input { width: 100%; padding: 15px; margin: 20px 0; background: #0F172A; border: 1px solid #475569; color: white; border-radius: 8px; box-sizing: border-box; }
.modal-actions { display: flex; gap: 10px; margin-bottom: 15px; }
.action-btn { flex: 1; padding: 10px; border: none; border-radius: 8px; font-weight: bold; color: white; }
.add { background: var(--usd-color); }
.sub { background: var(--try-color); }
.close-modal { background: transparent; border: 1px solid #475569; color: #94A3B8; padding: 8px 20px; border-radius: 8px; }
.slide-enter-active, .slide-leave-active { transition: transform 0.3s ease, opacity 0.3s ease; }
.slide-enter-from, .slide-leave-to { transform: translateX(-20px); opacity: 0; }
</style>