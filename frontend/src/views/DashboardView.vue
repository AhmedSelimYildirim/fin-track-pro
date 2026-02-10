<template>
  <div class="dashboard-page">
    <aside class="sidebar">
      <div class="brand-container">
        <div class="brand">FinTrack Pro 🚀</div>
        <div class="user-badge">{{ currentUser }}</div>
      </div>

      <nav class="menu">
        <div class="menu-item" :class="{ active: currentTab === 'dashboard' }" @click="currentTab = 'dashboard'">
          <span>📊</span> Portföy
        </div>
        <div class="menu-item" :class="{ active: currentTab === 'calendar' }" @click="router.push('/calendar')">
          <span>📅</span> Takvim & Notlar
        </div>
        <div class="menu-item" :class="{ active: currentTab === 'settings' }" @click="currentTab = 'settings'">
          <span>⚙️</span> Ayarlar
        </div>
      </nav>

      <div class="logout-wrapper">
        <div class="menu-item logout" @click="logout">
          <span>🚪</span> Çıkış Yap
        </div>
      </div>
    </aside>

    <main class="content">
      <SettingsView v-if="currentTab === 'settings'" />

      <div v-else>
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
                  <div @click="changeCurrency('GOLD', 24, '24 Ayar (Has)')">24 Ayar</div>
                  <div @click="changeCurrency('GOLD', 22, '22 Ayar')">22 Ayar</div>
                  <div @click="changeCurrency('GOLD', 18, '18 Ayar')">18 Ayar</div>
                  <div @click="changeCurrency('GOLD', 14, '14 Ayar')">14 Ayar</div>
                  <div @click="changeCurrency('GOLD', 8, '8 Ayar')">8 Ayar</div>
                  <div @click="changeCurrency('GOLD', 4, '4 Ayar')">4 Ayar</div>
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
                <span>Veri Yok</span>
                <small>Varlık Ekleyin</small>
              </div>
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
              <span class="card-name">Gold</span>
              <span class="card-amount">{{ getAmount('GOLD') }} Gr</span>
              <span class="card-val">%{{ getAllocation('GOLD') }}</span>
            </div>
          </div>
          <div class="asset-card card-usd" @click="openModal('USD')">
            <div class="card-icon">$</div>
            <div class="card-info">
              <span class="card-name">USD</span>
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
              <span class="card-name">Silver</span>
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
      </div>
    </main>

    <div v-if="showModal" class="modal-overlay" @click.self="showModal = false">
      <div class="modal-content large-modal">
        <div class="modal-header">
          <h3>{{ activeAsset }} İşlemleri</h3>
          <button @click="showModal = false">✕</button>
        </div>

        <div class="modal-body-split">
          <div class="transaction-form">
            <h4>Yeni İşlem</h4>
            <input v-model="amount" type="number" placeholder="Miktar" class="big-input" />

            <div v-if="activeAsset === 'GOLD'" class="ayar-wrapper">
              <select v-model="modalAyar" class="ayar-select">
                <option :value="24">24 Ayar (Has)</option>
                <option :value="22">22 Ayar</option>
                <option :value="18">18 Ayar</option>
                <option :value="14">14 Ayar</option>
                <option :value="8">8 Ayar</option>
                <option :value="4">4 Ayar</option>
              </select>
            </div>

            <input v-model="transactionDate" type="date" class="date-input" />

            <div class="actions">
              <button class="add" @click="handleTransaction('add')">EKLE (+)</button>
              <button class="sub" @click="handleTransaction('subtract')">ÇIKAR (-)</button>
            </div>
          </div>

          <div class="transaction-history">
            <h4>İşlem Geçmişi & Dekont</h4>
            <div class="history-list">
              <div v-if="filteredTransactions.length === 0" class="no-history">Henüz işlem yok.</div>
              <div v-for="tx in filteredTransactions" :key="tx.id" class="history-item">
                <div class="tx-info">
                  <span class="tx-date">{{ formatDate(tx.transaction_date) }}</span>
                  <span class="tx-type" :class="tx.type">{{ tx.type === 'add' ? 'Ekleme' : 'Çıkarma' }}</span>
                  <span class="tx-amount">{{ tx.amount }} <span v-if="tx.ayar > 0">({{ tx.ayar }}K)</span></span>
                </div>
                <button class="receipt-download-btn" @click="downloadSingleReceipt(tx.id)" title="Dekont İndir">
                  📄
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <button v-if="currentTab === 'dashboard'" class="receipt-btn" @click="downloadReceipt" title="Genel Rapor Al">📄</button>
  </div>
</template>

<script setup>
  import { ref, computed, onMounted } from 'vue';
  import { useRouter } from 'vue-router';
  import { Chart as ChartJS, ArcElement, Tooltip, Legend } from 'chart.js';
  import { Doughnut } from 'vue-chartjs';
  import api from '../services/api';
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
  const allTransactions = ref([]); // Tüm işlemleri burada tutacağız
  const baseCurrency = ref('TRY');
  const baseCurrencyLabel = ref('Türk Lirası (TL)');
  const targetAyar = ref(0);

  const toggleDropdown = () => { showSelector.value = !showSelector.value; };
  window.addEventListener('click', () => { if(showSelector.value) showSelector.value = false; });

  const displayCurrency = computed(() => baseCurrencyLabel.value);

  // ALTIN TOPLAMA DÜZELTMESİ: Tüm ayarları topla
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

    // Chart için de verileri topluyoruz
    const data = labels.map(label => {
      const assets = summaryData.value?.assets || [];
      // O tipe ait tüm varlıkların yüzdelerini topla (Örn: 24K Gold + 14K Gold)
      return assets
              .filter(a => a.type === label)
              .reduce((sum, curr) => sum + (curr.allocation || 0), 0);
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

  // İŞLEMLERİ ÇEKME
  const fetchTransactions = async () => {
    try {
      const res = await api.get('/assets/transactions', {
        headers: { 'X-Currency': baseCurrency.value }
      });
      allTransactions.value = res.data || [];
    } catch (e) { console.error(e); }
  };

  // AKTİF VARLIĞA GÖRE FİLTRELEME
  const filteredTransactions = computed(() => {
    if (!activeAsset.value) return [];
    return allTransactions.value
            .filter(tx => tx.asset_type === activeAsset.value)
            .sort((a, b) => new Date(b.transaction_date) - new Date(a.transaction_date));
  });

  // DÜZELTME: Belirli bir tipin (örn: GOLD) tüm miktarlarını topla
  const getAmount = (type) => {
    const assets = summaryData.value?.assets || [];
    const total = assets
            .filter(a => a.type === type)
            .reduce((sum, curr) => sum + curr.amount, 0);

    return total.toLocaleString('tr-TR', { maximumFractionDigits: 4 });
  };

  const getAllocation = (type) => {
    const assets = summaryData.value?.assets || [];
    const total = assets
            .filter(a => a.type === type)
            .reduce((sum, curr) => sum + (curr.allocation || 0), 0);

    return total.toFixed(1);
  };

  const openModal = (asset) => {
    activeAsset.value = asset;
    amount.value = '';
    modalAyar.value = 24;
    transactionDate.value = new Date().toISOString().split('T')[0];
    showModal.value = true;
    fetchTransactions(); // Modal açılınca geçmişi çek
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
      // Modalı kapatma, verileri yenile ki kullanıcı görsün
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

  // TEKİL DEKONT İNDİRME
  const downloadSingleReceipt = async (id) => {
    try {
      const res = await api.get(`/assets/receipt/${id}`, {
        responseType: 'blob',
        headers: { 'X-Currency': baseCurrency.value }
      });
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

  onMounted(() => {
    fetchData();
    fetchTransactions();
  });
</script>

<style scoped>
  .dashboard-page { display: flex; min-height: 100vh; background: var(--bg-color); color: var(--text-color); font-family: 'Segoe UI', sans-serif; }
  .sidebar { width: 260px; background: var(--sidebar-bg); display: flex; flex-direction: column; padding: 25px; border-right: 1px solid var(--border-color); }
  .brand-container { margin-bottom: 40px; text-align: center; }
  .brand { color: var(--accent-color); font-size: 1.6rem; font-weight: 800; letter-spacing: 1px; }
  .user-badge { margin-top: 5px; color: var(--success-color); font-weight: bold; font-size: 1.1rem; border-top: 1px solid var(--border-color); padding-top: 5px; }
  .menu-item { padding: 15px; margin-bottom: 10px; border-radius: 12px; cursor: pointer; color: var(--text-muted); display: flex; gap: 12px; align-items: center; transition: all 0.3s; font-weight: 500; }
  .menu-item:hover, .menu-item.active { background: var(--hover-bg); color: var(--text-color); transform: translateX(5px); }
  .logout-wrapper { margin-top: auto; padding-top: 20px; border-top: 1px solid rgba(255,255,255,0.1); }
  .logout { border: 2px solid var(--danger-color); border-radius: 12px; color: var(--danger-color); justify-content: center; font-weight: bold; padding: 12px; }
  .logout:hover { background: var(--danger-color); color: white; transform: none; }
  .content { flex: 1; padding: 40px; overflow-y: auto; }
  .top-bar { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
  .page-title h2 { color: var(--text-color); margin: 0; }
  .currency-wrapper { position: relative; z-index: 100; }
  .currency-btn { background: var(--sidebar-bg); border: 1px solid var(--border-color); padding: 12px 25px; border-radius: 25px; cursor: pointer; font-weight: bold; color: var(--accent-color); min-width: 150px; text-align: center; transition: 0.3s; }
  .currency-btn:hover { border-color: var(--accent-color); }
  .currency-dropdown { position: absolute; top: 110%; right: 0; background: var(--sidebar-bg); border: 1px solid var(--border-color); border-radius: 15px; width: 220px; box-shadow: 0 10px 30px rgba(0,0,0,0.5); }
  .c-item { padding: 12px 20px; cursor: pointer; border-bottom: 1px solid var(--border-color); position: relative; color: var(--text-color); }
  .c-item:hover { background: var(--hover-bg); color: var(--accent-color); }
  .has-submenu:hover .submenu { display: block; }
  .submenu { display: none; position: absolute; top: 0; right: 100%; background: var(--sidebar-bg); border: 1px solid var(--border-color); border-radius: 15px; width: 180px; box-shadow: 0 10px 30px rgba(0,0,0,0.5); margin-right: 5px; }
  .submenu div { padding: 12px 20px; cursor: pointer; border-bottom: 1px solid var(--border-color); color: var(--text-color); }
  .submenu div:hover { background: var(--hover-bg); color: var(--accent-color); }
  .chart-section { display: flex; flex-direction: column; align-items: center; margin: 20px 0 50px 0; position: relative; }
  .chart-wrapper { width: 300px; height: 300px; position: relative; }
  .center-balance { position: absolute; top: 50%; left: 50%; transform: translate(-50%, -50%); text-align: center; pointer-events: none; }
  .center-balance h3 { font-size: 1.8rem; margin: 0; font-weight: 800; color: var(--text-color); }
  .center-balance small { color: var(--text-muted); }
  .total-underline { width: 150px; height: 4px; background: linear-gradient(90deg, transparent, var(--accent-color), transparent); margin-top: 20px; border-radius: 2px; }
  .no-data-circle { width: 100%; height: 100%; border-radius: 50%; border: 8px dashed var(--border-color); display: flex; align-items: center; justify-content: center; }
  .no-data-content { text-align: center; color: var(--text-muted); font-weight: bold; }
  .assets-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 20px; padding: 10px; }
  .asset-card { padding: 20px; border-radius: 20px; display: flex; align-items: center; gap: 15px; cursor: pointer; transition: all 0.3s ease; position: relative; overflow: hidden; border: 1px solid var(--border-color); box-shadow: 0 10px 20px rgba(0,0,0,0.2); }
  .asset-card:hover { transform: translateY(-7px) scale(1.02); }
  .asset-card, .asset-card * { color: white !important; }
  .card-btc { background: linear-gradient(135deg, #2b2b2b 0%, #000000 100%); border-bottom: 4px solid #555; }
  .card-gold { background: linear-gradient(135deg, #DAA520 0%, #FFD700 50%, #B8860B 100%); border-bottom: 4px solid #FFF; }
  .card-usd { background: linear-gradient(135deg, #059669 0%, #10B981 100%); border-bottom: 4px solid #A7F3D0; }
  .card-eur { background: linear-gradient(135deg, #5D4037 0%, #8D6E63 100%); border-bottom: 4px solid #D7CCC8; }
  .card-silver { background: linear-gradient(135deg, #757575 0%, #9E9E9E 100%); border-bottom: 4px solid #E0E0E0; }
  .card-try { background: linear-gradient(135deg, #991B1B 0%, #EF4444 100%); border-bottom: 4px solid #FCA5A5; }
  .card-icon { font-size: 2rem; }
  .card-info { display: flex; flex-direction: column; }
  .card-name { font-size: 0.9rem; opacity: 0.9; text-transform: uppercase; font-weight: bold; }
  .card-amount { font-size: 1.2rem; font-weight: bold; }
  .card-val { font-size: 0.8rem; opacity: 0.7; }

  /* MODAL DÜZENLEMELERİ (Split View) */
  .modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.85); display: flex; justify-content: center; align-items: center; z-index: 200; backdrop-filter: blur(5px); }
  .large-modal { width: 700px !important; max-width: 95%; }
  .modal-content { background: var(--card-bg); padding: 30px; border-radius: 20px; border: 1px solid var(--border-color); box-shadow: 0 20px 50px rgba(0,0,0,0.5); }
  .modal-header { display: flex; justify-content: space-between; margin-bottom: 20px; border-bottom: 1px solid var(--border-color); padding-bottom: 10px; }
  .modal-header h3 { color: var(--text-color); margin: 0; }
  .modal-header button { background: none; border: none; color: var(--text-color); font-size: 1.5rem; cursor: pointer; }

  .modal-body-split { display: grid; grid-template-columns: 1fr 1fr; gap: 30px; }
  .transaction-form, .transaction-history { display: flex; flex-direction: column; gap: 10px; }
  .transaction-history { border-left: 1px solid var(--border-color); padding-left: 20px; }
  h4 { color: var(--text-muted); margin: 0 0 10px 0; }

  .history-list { max-height: 250px; overflow-y: auto; display: flex; flex-direction: column; gap: 8px; }
  .history-item { background: var(--input-bg); padding: 10px; border-radius: 8px; display: flex; justify-content: space-between; align-items: center; font-size: 0.9rem; }
  .tx-info { display: flex; flex-direction: column; }
  .tx-date { color: var(--text-muted); font-size: 0.8rem; }
  .tx-type.add { color: var(--success-color); }
  .tx-type.subtract { color: var(--danger-color); }
  .tx-amount { color: var(--text-color); font-weight: bold; }
  .receipt-download-btn { background: transparent; border: 1px solid var(--border-color); border-radius: 5px; cursor: pointer; padding: 5px; transition: 0.2s; }
  .receipt-download-btn:hover { background: var(--hover-bg); }

  .big-input, .ayar-select, .date-input { width: 100%; padding: 12px; background: var(--input-bg); border: 1px solid var(--border-color); color: var(--text-color); margin-bottom: 15px; font-size: 1.1rem; border-radius: 10px; box-sizing: border-box; }
  .actions { display: flex; gap: 15px; margin-top: 10px; }
  .actions button { flex: 1; padding: 15px; border: none; border-radius: 12px; color: white; font-weight: bold; cursor: pointer; transition: 0.2s; }
  .add { background: var(--success-color); } .sub { background: var(--danger-color); }
  .receipt-btn { position: fixed; bottom: 30px; right: 30px; background: #3B82F6; width: 60px; height: 60px; border-radius: 50%; border: none; font-size: 24px; cursor: pointer; box-shadow: 0 5px 20px rgba(59, 130, 246, 0.5); z-index: 90; }

  @media (max-width: 768px) {
    .modal-body-split { grid-template-columns: 1fr; }
    .transaction-history { border-left: none; border-top: 1px solid var(--border-color); padding-left: 0; padding-top: 20px; }
  }
</style>