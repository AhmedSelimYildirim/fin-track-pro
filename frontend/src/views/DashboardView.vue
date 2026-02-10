<template>
  <div class="dashboard">
    <header class="top-bar">
      <div class="left-group">
        <button class="menu-btn" @click="isSidebarOpen = !isSidebarOpen">
          <component :is="MenuIcon" />
        </button>
        <div class="user-info">
          <span class="welcome-text">Hoşgeldin,</span>
          <span class="u-name">{{ currentUser }}</span>
        </div>
      </div>

      <div class="currency-selector">
        <button class="select-trigger" @click="showSelector = !showSelector">
          {{ displayCurrency }} ▼
        </button>
        <div v-if="showSelector" class="scroll-dropdown">
          <div v-for="opt in currencyOptions" :key="opt.label" class="dropdown-item" @click="changeCurrency(opt)">
            {{ opt.label }}
          </div>
        </div>
      </div>
    </header>

    <transition name="slide">
      <div v-if="isSidebarOpen" class="sidebar">
        <div class="sidebar-header">FinTrack Pro</div>
        <div class="menu-item active"><component :is="LayoutDashboard" size="20"/> Dashboard</div>
        <div class="menu-item"><component :is="Calendar" size="20"/> Calendar / Notes</div>
        <div class="menu-item"><component :is="User" size="20"/> Profile</div>
        <div class="menu-spacer"></div>
        <div class="menu-item red" @click="logout"><component :is="LogOut" size="20"/> Logout</div>
      </div>
    </transition>
    <div v-if="isSidebarOpen" class="sidebar-overlay" @click="isSidebarOpen = false"></div>

    <main class="cockpit">

      <div class="chart-section">
        <div class="chart-wrapper">
          <Doughnut v-if="chartData" :data="chartData" :options="chartOptions" />
          <div class="center-balance">
            <h2>{{ totalValue }}</h2>
            <small>{{ baseCurrency === 'GOLD' ? 'Gram' : baseCurrency }}</small>
          </div>
        </div>
      </div>

      <div class="assets-grid">
        <div class="asset-card btc-card" @click="openModal('BTC')">
          <div class="card-icon">₿</div>
          <div class="card-info">
            <span class="asset-name">Bitcoin</span>
            <span class="asset-percent">%{{ getAllocation('BTC') }}</span>
          </div>
        </div>

        <div class="asset-card gold-card" @click="openModal('GOLD')">
          <div class="card-icon">👑</div>
          <div class="card-info">
            <span class="asset-name">Gold</span>
            <span class="asset-percent">%{{ getAllocation('GOLD') }}</span>
          </div>
        </div>

        <div class="asset-card usd-card" @click="openModal('USD')">
          <div class="card-icon">$</div>
          <div class="card-info">
            <span class="asset-name">USD</span>
            <span class="asset-percent">%{{ getAllocation('USD') }}</span>
          </div>
        </div>

        <div class="asset-card eur-card" @click="openModal('EUR')">
          <div class="card-icon">€</div>
          <div class="card-info">
            <span class="asset-name">Euro</span>
            <span class="asset-percent">%{{ getAllocation('EUR') }}</span>
          </div>
        </div>

        <div class="asset-card silver-card" @click="openModal('SILVER')">
          <div class="card-icon">⚔️</div>
          <div class="card-info">
            <span class="asset-name">Silver</span>
            <span class="asset-percent">%{{ getAllocation('SILVER') }}</span>
          </div>
        </div>

        <div class="asset-card try-card" @click="openModal('TRY')">
          <div class="card-icon">₺</div>
          <div class="card-info">
            <span class="asset-name">TL</span>
            <span class="asset-percent">%{{ getAllocation('TRY') }}</span>
          </div>
        </div>
      </div>
    </main>

    <div v-if="showModal" class="modal-overlay">
      <div class="modal-content">
        <div class="modal-header">
          <h3>Manage {{ activeAsset }}</h3>
          <button class="close-icon" @click="showModal = false">✕</button>
        </div>

        <div class="modal-body">
          <label>Miktar</label>
          <input v-model="amount" type="number" placeholder="0.00" class="modal-input big-input" />

          <label>İşlem Tarihi</label>
          <input v-model="transactionDate" type="date" class="modal-input date-input" />

          <div class="modal-actions">
            <button class="action-btn add" @click="handleTransaction('add')">
              <span>GİRİŞ (+)</span>
            </button>
            <button class="action-btn sub" @click="handleTransaction('subtract')">
              <span>ÇIKIŞ (-)</span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <button class="receipt-btn" @click="downloadReceipt" title="Rapor Al">
      <component :is="FileText" size="24" />
    </button>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { Chart as ChartJS, ArcElement, Tooltip, Legend } from 'chart.js';
import { Doughnut } from 'vue-chartjs';
import { Menu as MenuIcon, Calendar, User, LogOut, FileText, LayoutDashboard } from 'lucide-vue-next';
import api from '../services/api';

ChartJS.register(ArcElement, Tooltip, Legend);

const router = useRouter();
const currentUser = localStorage.getItem('username') || 'Yatırımcı';

// Durumlar (State)
const isSidebarOpen = ref(false);
const showSelector = ref(false);
const showModal = ref(false);
const activeAsset = ref('');
const amount = ref('');
const transactionDate = ref(new Date().toISOString().split('T')[0]); // Bugünün tarihi (YYYY-MM-DD)
const summaryData = ref(null);
const baseCurrency = ref('TRY');
const targetAyar = ref(0);

// Para Birimi Seçenekleri
const currencyOptions = [
  { label: 'Turkish Lira (TRY)', val: 'TRY', ayar: 0 },
  { label: 'US Dollar (USD)', val: 'USD', ayar: 0 },
  { label: 'Euro (EUR)', val: 'EUR', ayar: 0 },
  { label: 'Bitcoin (BTC)', val: 'BTC', ayar: 0 },
  { label: 'Gold (Has)', val: 'GOLD', ayar: 24 },
  { label: 'Gold (22 Ayar)', val: 'GOLD', ayar: 22 },
];

const displayCurrency = computed(() => {
  const found = currencyOptions.find(c => c.val === baseCurrency.value && c.ayar === targetAyar.value);
  return found ? found.label : 'Para Birimi Seç';
});

const totalValue = computed(() => {
  if (!summaryData.value) return '0.00';
  return summaryData.value.total_value.toLocaleString('tr-TR', { minimumFractionDigits: 2, maximumFractionDigits: 2 });
});

// Grafik Verileri
const chartData = computed(() => {
  if (!summaryData.value) return null;
  // Backend'den gelen veriyi renklere eşle
  const labels = ['BTC', 'GOLD', 'USD', 'EUR', 'SILVER', 'TRY'];
  const colors = ['#F7931A', '#FFD700', '#10B981', '#3B82F6', '#94A3B8', '#EF4444'];

  // Veriyi sıralı çekmek için map yapıyoruz
  const dataValues = labels.map(type => {
    const asset = summaryData.value.assets.find(a => a.type === type);
    return asset ? asset.allocation : 0;
  });

  return {
    labels: labels,
    datasets: [{
      backgroundColor: colors,
      data: dataValues,
      borderWidth: 0,
      hoverOffset: 10
    }]
  };
});

const chartOptions = {
  responsive: true,
  cutout: '75%',
  plugins: { legend: { display: false }, tooltip: { enabled: true } }
};

// Fonksiyonlar
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

const getAllocation = (type) => {
  if (!summaryData.value) return '0';
  const asset = summaryData.value.assets.find(a => a.type === type);
  return asset ? asset.allocation.toFixed(1) : '0.0';
};

const openModal = (asset) => {
  activeAsset.value = asset;
  amount.value = '';
  transactionDate.value = new Date().toISOString().split('T')[0]; // Tarihi bugüne sıfırla
  showModal.value = true;
  isSidebarOpen.value = false;
};

const handleTransaction = async (action) => {
  try {
    if (!amount.value || amount.value <= 0) {
      alert("Lütfen geçerli bir miktar girin.");
      return;
    }

    // Backend'in istediği tarih formatı: 2023-02-08T21:10:00Z
    // Bizim inputtan gelen: 2023-02-08
    const datePayload = new Date(transactionDate.value);
    const isoDate = datePayload.toISOString(); // Otomatik olarak backend formatına çevirir

    await api.post('/assets/balance', {
      type: activeAsset.value,
      amount: parseFloat(amount.value),
      action: action,
      transaction_date: isoDate, // Tarih eklendi
      ayar: activeAsset.value === 'GOLD' ? 24 : 0
    });

    showModal.value = false;
    fetchData(); // Verileri yenile
  } catch (e) {
    alert('Hata: ' + e.message);
  }
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
    link.setAttribute('download', `Portfolio_Report_${new Date().toISOString().split('T')[0]}.pdf`);
    document.body.appendChild(link);
    link.click();
  } catch (e) {
    alert('Rapor indirilemedi.');
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
/* Genel Yapı */
.dashboard { min-height: 100vh; background: var(--bg-dark); color: white; padding-bottom: 80px; font-family: 'Segoe UI', sans-serif; }

/* Üst Bar */
.top-bar { display: flex; justify-content: space-between; align-items: center; padding: 20px; background: rgba(30, 41, 59, 0.5); backdrop-filter: blur(10px); position: sticky; top: 0; z-index: 40; border-bottom: 1px solid rgba(255,255,255,0.05); }
.left-group { display: flex; align-items: center; gap: 15px; }
.menu-btn { background: transparent; border: none; color: white; cursor: pointer; }
.user-info { display: flex; flex-direction: column; }
.welcome-text { font-size: 0.75rem; color: #94A3B8; }
.u-name { font-weight: bold; font-size: 1rem; }

/* Para Birimi Seçici */
.currency-selector { position: relative; }
.select-trigger { background: #334155; color: white; border: none; padding: 8px 16px; border-radius: 20px; font-size: 0.9rem; cursor: pointer; transition: 0.2s; }
.select-trigger:hover { background: #475569; }
.scroll-dropdown { position: absolute; top: 120%; right: 0; background: var(--card-bg); border: 1px solid #475569; border-radius: 12px; width: 180px; overflow: hidden; z-index: 100; box-shadow: 0 10px 40px rgba(0,0,0,0.5); }
.dropdown-item { padding: 12px; cursor: pointer; font-size: 0.9rem; color: #cbd5e1; border-bottom: 1px solid rgba(255,255,255,0.05); }
.dropdown-item:hover { background: #334155; color: white; }

/* Sidebar */
.sidebar { position: fixed; top: 0; left: 0; height: 100vh; width: 260px; background: var(--card-bg); z-index: 1000; padding: 30px 20px; display: flex; flex-direction: column; box-shadow: 10px 0 30px rgba(0,0,0,0.5); }
.sidebar-header { font-size: 1.5rem; font-weight: bold; margin-bottom: 40px; color: var(--gold-color); letter-spacing: 1px; }
.menu-item { display: flex; align-items: center; gap: 15px; padding: 15px; border-radius: 12px; color: #94A3B8; cursor: pointer; transition: 0.2s; font-weight: 500; }
.menu-item:hover, .menu-item.active { background: #334155; color: white; }
.menu-spacer { flex-grow: 1; }
.sidebar-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.8); z-index: 900; backdrop-filter: blur(2px); }
.red { color: #EF4444; } .red:hover { background: rgba(239, 68, 68, 0.1); color: #EF4444; }

/* Kokpit & Grafik */
.cockpit { padding: 20px; max-width: 1200px; margin: 0 auto; }
.chart-section { display: flex; justify-content: center; margin-bottom: 40px; margin-top: 10px; }
.chart-wrapper { position: relative; width: 260px; height: 260px; }
.center-balance { position: absolute; top: 50%; left: 50%; transform: translate(-50%, -50%); text-align: center; pointer-events: none; }
.center-balance h2 { margin: 0; font-size: 1.8rem; font-weight: 800; }
.center-balance small { color: #94A3B8; font-size: 0.9rem; }

/* Varlık Kartları (Grid) */
.assets-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(150px, 1fr)); gap: 15px; }
.asset-card { background: var(--card-bg); border-radius: 16px; padding: 20px; display: flex; flex-direction: column; align-items: flex-start; gap: 10px; cursor: pointer; transition: transform 0.2s, box-shadow 0.2s; border: 1px solid rgba(255,255,255,0.05); position: relative; overflow: hidden; }
.asset-card:hover { transform: translateY(-5px); box-shadow: 0 10px 20px rgba(0,0,0,0.3); }

/* Kart Renkleri ve İkonlar */
.btc-card { border-top: 4px solid #F7931A; } .btc-card .card-icon { color: #F7931A; }
.gold-card { border-top: 4px solid #FFD700; } .gold-card .card-icon { color: #FFD700; }
.usd-card { border-top: 4px solid #10B981; } .usd-card .card-icon { color: #10B981; }
.eur-card { border-top: 4px solid #3B82F6; } .eur-card .card-icon { color: #3B82F6; }
.silver-card { border-top: 4px solid #94A3B8; } .silver-card .card-icon { color: #94A3B8; }
.try-card { border-top: 4px solid #EF4444; } .try-card .card-icon { color: #EF4444; }

.card-icon { font-size: 1.5rem; background: rgba(255,255,255,0.05); width: 40px; height: 40px; display: flex; align-items: center; justify-content: center; border-radius: 10px; }
.card-info { display: flex; flex-direction: column; }
.asset-name { font-size: 0.9rem; color: #cbd5e1; font-weight: 500; }
.asset-percent { font-size: 1.2rem; font-weight: bold; color: white; margin-top: 2px; }

/* Modal */
.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.85); display: flex; justify-content: center; align-items: center; z-index: 2000; backdrop-filter: blur(5px); }
.modal-content { background: #1e293b; width: 90%; max-width: 400px; border-radius: 20px; overflow: hidden; box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5); border: 1px solid #334155; }
.modal-header { padding: 20px; background: #0f172a; display: flex; justify-content: space-between; align-items: center; border-bottom: 1px solid #334155; }
.modal-header h3 { margin: 0; color: white; }
.close-icon { background: none; border: none; color: #94A3B8; font-size: 1.2rem; cursor: pointer; }
.modal-body { padding: 25px; display: flex; flex-direction: column; gap: 15px; }
.modal-body label { color: #94A3B8; font-size: 0.85rem; margin-bottom: -5px; }
.modal-input { background: #0F172A; border: 1px solid #334155; color: white; padding: 15px; border-radius: 10px; font-size: 1rem; outline: none; transition: 0.2s; width: 100%; box-sizing: border-box; }
.modal-input:focus { border-color: var(--usd-color); }
.big-input { font-size: 1.5rem; font-weight: bold; text-align: center; letter-spacing: 1px; }
.date-input { color-scheme: dark; } /* Takvim ikonunu beyaz yapar */

.modal-actions { display: flex; gap: 10px; margin-top: 10px; }
.action-btn { flex: 1; padding: 15px; border: none; border-radius: 10px; font-weight: bold; color: white; cursor: pointer; display: flex; flex-direction: column; align-items: center; gap: 5px; transition: 0.2s; }
.action-btn:active { transform: scale(0.98); }
.add { background: linear-gradient(135deg, #10B981 0%, #059669 100%); }
.sub { background: linear-gradient(135deg, #EF4444 0%, #B91C1C 100%); }

/* Rapor Butonu */
.receipt-btn { position: fixed; bottom: 30px; right: 30px; background: #3B82F6; width: 60px; height: 60px; border-radius: 50%; border: none; color: white; display: flex; align-items: center; justify-content: center; box-shadow: 0 10px 25px rgba(59, 130, 246, 0.5); cursor: pointer; transition: 0.3s; z-index: 50; }
.receipt-btn:hover { transform: translateY(-5px); background: #2563EB; }

/* Animasyonlar */
.slide-enter-active, .slide-leave-active { transition: transform 0.3s ease; }
.slide-enter-from, .slide-leave-to { transform: translateX(-100%); }
</style>