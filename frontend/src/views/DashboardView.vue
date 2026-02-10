<template>
  <div class="dashboard-page">
    <aside class="sidebar">
      <div class="brand-container">
        <div class="brand">FinTrack Pro 🚀</div>
        <div class="user-badge">{{ currentUser }}</div>
      </div>

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
        <div class="page-title">
          <h2>Portföy Özeti</h2>
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
              Altın (Gram) ▶
              <div class="submenu">
                <div @click="changeCurrency('GOLD', 24, 'Altın (24 Ayar)')">24 Ayar (Has)</div>
                <div @click="changeCurrency('GOLD', 22, 'Altın (22 Ayar)')">22 Ayar</div>
                <div @click="changeCurrency('GOLD', 18, 'Altın (18 Ayar)')">18 Ayar</div>
                <div @click="changeCurrency('GOLD', 14, 'Altın (14 Ayar)')">14 Ayar</div>
                <div @click="changeCurrency('GOLD', 8, 'Altın (8 Ayar)')">8 Ayar</div>
                <div @click="changeCurrency('GOLD', 4, 'Altın (4 Ayar)')">4 Ayar</div>
              </div>
            </div>
          </div>
        </div>
      </header>

      <div class="chart-section">
        <div class="chart-wrapper">
          <Doughnut v-if="chartData" :data="chartData" :options="chartOptions" />
          <div class="center-balance">
            <h3>{{ totalValue }}</h3>
            <small>{{ baseCurrencyLabel }}</small>
          </div>
        </div>
        <div class="total-underline"></div>
      </div>

      <div class="assets-grid">
        <div class="asset-card card-btc" @click="openModal('BTC')">
          <div class="card-icon">₿</div>
          <div class="card-info">
            <span class="card-name">Bitcoin</span>
            <span class="card-amount">{{ getAmount('BTC') }} BTC</span>
            <span class="card-val">%{{ getAllocation('BTC') }}</span>
          </div>
        </div>

        <div class="asset-card card-gold" @click="openModal('GOLD')">
          <div class="card-icon">👑</div>
          <div class="card-info">
            <span class="card-name">Altın</span>
            <span class="card-amount">{{ getAmount('GOLD') }} Gr</span>
            <span class="card-val">%{{ getAllocation('GOLD') }}</span>
          </div>
        </div>

        <div class="asset-card card-usd" @click="openModal('USD')">
          <div class="card-icon">$</div>
          <div class="card-info">
            <span class="card-name">Dolar</span>
            <span class="card-amount">{{ getAmount('USD') }} $</span>
            <span class="card-val">%{{ getAllocation('USD') }}</span>
          </div>
        </div>

        <div class="asset-card card-eur" @click="openModal('EUR')">
          <div class="card-icon">€</div>
          <div class="card-info">
            <span class="card-name">Euro</span>
            <span class="card-amount">{{ getAmount('EUR') }} €</span>
            <span class="card-val">%{{ getAllocation('EUR') }}</span>
          </div>
        </div>

        <div class="asset-card card-silver" @click="openModal('SILVER')">
          <div class="card-icon">⚔️</div>
          <div class="card-info">
            <span class="card-name">Gümüş</span>
            <span class="card-amount">{{ getAmount('SILVER') }} Gr</span>
            <span class="card-val">%{{ getAllocation('SILVER') }}</span>
          </div>
        </div>

        <div class="asset-card card-try" @click="openModal('TRY')">
          <div class="card-icon">₺</div>
          <div class="card-info">
            <span class="card-name">TL Nakit</span>
            <span class="card-amount">{{ getAmount('TRY') }} ₺</span>
            <span class="card-val">%{{ getAllocation('TRY') }}</span>
          </div>
        </div>
      </div>
    </main>

    <div v-if="showModal" class="modal-overlay" @click.self="showModal = false">
      <div class="modal-content">
        <div class="modal-header">
          <h3>{{ activeAsset }} İşlemi</h3>
          <button @click="showModal = false">✕</button>
        </div>
        <div class="modal-body">
          <input v-model="amount" type="number" placeholder="Miktar Giriniz" class="big-input" />

          <div v-if="activeAsset === 'GOLD'" class="ayar-wrapper">
            <label>Altın Ayarı:</label>
            <select v-model="modalAyar" class="ayar-select">
              <option :value="24">24 Ayar (Has)</option>
              <option :value="22">22 Ayar</option>
              <option :value="18">18 Ayar</option>
              <option :value="14">14 Ayar</option>
              <option :value="8">8 Ayar</option>
              <option :value="4">4 Ayar</option>
            </select>
          </div>

          <div class="date-wrapper">
            <label>İşlem Tarihi:</label>
            <input v-model="transactionDate" type="date" class="date-input" />
          </div>

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

  const currentUser = ref(localStorage.getItem('username') || 'Yatırımcı');
  const showSelector = ref(false);
  const showModal = ref(false);
  const activeAsset = ref('');
  const amount = ref('');
  const transactionDate = ref(new Date().toISOString().split('T')[0]);
  const modalAyar = ref(24);
  const summaryData = ref(null);
  const baseCurrency = ref('TRY');
  const baseCurrencyLabel = ref('Türk Lirası (TL)');
  const targetAyar = ref(0);

  const toggleDropdown = () => {
    showSelector.value = !showSelector.value;
  };

  window.addEventListener('click', () => {
    if(showSelector.value) showSelector.value = false;
  });

  const displayCurrency = computed(() => baseCurrencyLabel.value);

  const totalValue = computed(() => summaryData.value ? summaryData.value.total_value.toLocaleString('tr-TR', { maximumFractionDigits: 2 }) : '0');

  const chartData = computed(() => {
    if (!summaryData.value) return null;
    const labels = ['BTC', 'GOLD', 'USD', 'EUR', 'SILVER', 'TRY'];
    const colors = ['#1a1a1a', '#FFD700', '#10B981', '#8B4513', '#A0A0A0', '#EF4444'];
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
    } catch(e) {
      if(e.response && e.response.status === 401) {
        router.push('/login');
      }
      console.error(e);
    }
  };

  const getAmount = (type) => {
    const asset = summaryData.value?.assets.find(a => a.type === type);
    return asset ? asset.amount.toLocaleString('tr-TR', { maximumFractionDigits: 4 }) : '0';
  };

  const getAllocation = (type) => summaryData.value?.assets.find(a => a.type === type)?.allocation.toFixed(1) || '0.0';

  const openModal = (asset) => {
    activeAsset.value = asset;
    amount.value = '';
    modalAyar.value = 24;
    transactionDate.value = new Date().toISOString().split('T')[0];
    showModal.value = true;
  };

  const handleTransaction = async (action) => {
    if (!amount.value) return alert("Lütfen miktar giriniz.");
    try {
      await api.post('/assets/balance', {
        type: activeAsset.value,
        amount: parseFloat(amount.value),
        action,
        transaction_date: new Date(transactionDate.value).toISOString(),
        ayar: activeAsset.value === 'GOLD' ? parseInt(modalAyar.value) : 0
      });
      showModal.value = false;
      fetchData();
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
    } catch (e) { alert('Rapor alınırken hata oluştu.'); }
  };

  const logout = () => { localStorage.clear(); router.push('/login'); };

  onMounted(fetchData);
</script>

<style scoped>
  .dashboard-page { display: flex; min-height: 100vh; background: #0F172A; color: white; font-family: 'Segoe UI', sans-serif; }
  .sidebar { width: 260px; background: #1E293B; display: flex; flex-direction: column; padding: 25px; border-right: 1px solid rgba(255,255,255,0.05); }
  .brand-container { margin-bottom: 40px; text-align: center; }
  .brand { color: #FFD700; font-size: 1.6rem; font-weight: 800; letter-spacing: 1px; }
  .user-badge { margin-top: 5px; color: #10B981; font-weight: bold; font-size: 1.1rem; border-top: 1px solid rgba(255,255,255,0.1); padding-top: 5px; }
  .menu-item { padding: 15px; margin-bottom: 10px; border-radius: 12px; cursor: pointer; color: #94A3B8; display: flex; gap: 12px; align-items: center; transition: all 0.3s; font-weight: 500; }
  .menu-item:hover, .menu-item.active { background: linear-gradient(90deg, #334155, transparent); color: white; transform: translateX(5px); }
  .logout-wrapper { margin-top: auto; }
  .logout { color: #EF4444; } .logout:hover { background: rgba(239, 68, 68, 0.1); transform: translateX(5px); }
  .content { flex: 1; padding: 40px; }
  .top-bar { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
  .currency-wrapper { position: relative; z-index: 100; }
  .currency-btn { background: linear-gradient(135deg, #1E293B, #0F172A); border: 1px solid #475569; padding: 12px 25px; border-radius: 25px; cursor: pointer; font-weight: bold; color: #FFD700; min-width: 150px; text-align: center; transition: 0.3s; }
  .currency-btn:hover { border-color: #FFD700; }
  .currency-dropdown { position: absolute; top: 110%; right: 0; background: #1E293B; border: 1px solid #475569; border-radius: 15px; width: 220px; box-shadow: 0 10px 30px rgba(0,0,0,0.5); }
  .c-item { padding: 12px 20px; cursor: pointer; border-bottom: 1px solid rgba(255,255,255,0.05); position: relative; }
  .c-item:hover { background: #334155; color: #FFD700; }
  .has-submenu:hover .submenu { display: block; }
  .submenu { display: none; position: absolute; top: 0; right: 100%; background: #1E293B; border: 1px solid #475569; border-radius: 15px; width: 180px; box-shadow: 0 10px 30px rgba(0,0,0,0.5); margin-right: 5px; }
  .submenu div { padding: 12px 20px; cursor: pointer; border-bottom: 1px solid rgba(255,255,255,0.05); }
  .submenu div:hover { background: #334155; color: #FFD700; }
  .chart-section { display: flex; flex-direction: column; align-items: center; margin: 20px 0 50px 0; position: relative; }
  .chart-wrapper { width: 300px; height: 300px; position: relative; }
  .center-balance { position: absolute; top: 50%; left: 50%; transform: translate(-50%, -50%); text-align: center; pointer-events: none; }
  .center-balance h3 { font-size: 1.8rem; margin: 0; font-weight: 800; }
  .center-balance small { color: #94A3B8; }
  .total-underline { width: 150px; height: 4px; background: linear-gradient(90deg, transparent, #FFD700, transparent); margin-top: 20px; border-radius: 2px; }
  .assets-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 20px; padding: 10px; }
  .asset-card { padding: 20px; border-radius: 20px; display: flex; align-items: center; gap: 15px; cursor: pointer; transition: all 0.3s ease; position: relative; overflow: hidden; border: 1px solid rgba(255,255,255,0.05); box-shadow: 0 10px 20px rgba(0,0,0,0.2); }
  .asset-card:hover { transform: translateY(-7px) scale(1.02); }
  .card-btc { background: linear-gradient(135deg, #2b2b2b 0%, #000000 100%); border-bottom: 4px solid #555; }
  .card-btc .card-icon { color: #e0e0e0; }
  .card-gold { background: linear-gradient(135deg, #DAA520 0%, #FFD700 50%, #B8860B 100%); color: #000; border-bottom: 4px solid #FFF; }
  .card-gold .card-name, .card-gold .card-amount { color: #000; font-weight: 900; }
  .card-usd { background: linear-gradient(135deg, #059669 0%, #10B981 100%); border-bottom: 4px solid #A7F3D0; }
  .card-eur { background: linear-gradient(135deg, #5D4037 0%, #8D6E63 100%); border-bottom: 4px solid #D7CCC8; }
  .card-silver { background: linear-gradient(135deg, #757575 0%, #9E9E9E 100%); border-bottom: 4px solid #E0E0E0; }
  .card-try { background: linear-gradient(135deg, #991B1B 0%, #EF4444 100%); border-bottom: 4px solid #FCA5A5; }
  .card-icon { font-size: 2rem; }
  .card-info { display: flex; flex-direction: column; }
  .card-name { font-size: 0.9rem; opacity: 0.9; text-transform: uppercase; font-weight: bold; }
  .card-amount { font-size: 1.2rem; font-weight: bold; }
  .card-val { font-size: 0.8rem; opacity: 0.7; }
  .modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.85); display: flex; justify-content: center; align-items: center; z-index: 200; backdrop-filter: blur(5px); }
  .modal-content { background: #1E293B; padding: 30px; border-radius: 20px; width: 350px; border: 1px solid rgba(255,255,255,0.1); box-shadow: 0 20px 50px rgba(0,0,0,0.5); }
  .modal-header { display: flex; justify-content: space-between; margin-bottom: 20px; border-bottom: 1px solid rgba(255,255,255,0.1); padding-bottom: 10px; }
  .modal-body label { display: block; margin-bottom: 5px; color: #94A3B8; font-size: 0.9rem; }
  .big-input, .ayar-select, .date-input { width: 100%; padding: 12px; background: #0F172A; border: 1px solid #334155; color: white; margin-bottom: 15px; font-size: 1.1rem; border-radius: 10px; box-sizing: border-box; }
  .actions { display: flex; gap: 15px; margin-top: 10px; }
  .actions button { flex: 1; padding: 15px; border: none; border-radius: 12px; color: white; font-weight: bold; cursor: pointer; transition: 0.2s; }
  .add { background: #10B981; } .sub { background: #EF4444; }
  .receipt-btn { position: fixed; bottom: 30px; right: 30px; background: #3B82F6; width: 60px; height: 60px; border-radius: 50%; border: none; font-size: 24px; cursor: pointer; box-shadow: 0 5px 20px rgba(59, 130, 246, 0.5); z-index: 90; }
</style>