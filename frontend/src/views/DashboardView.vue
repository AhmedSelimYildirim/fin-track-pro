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
            <span class="curr-label">{{ totalValue }} {{ displayCurrency }}</span>
            <span class="arrow-icon">▼</span>
          </div>

          <transition name="dropdown-anim">
            <div v-if="showSelector" class="currency-dropdown">
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
                <span class="symbol text-icon">Gr</span> SILVER
              </div>

              <div class="c-item has-submenu">
                <div class="gold-trigger">
                  <div style="display:flex; align-items:center; gap:12px;">
                    <span class="symbol text-icon">Gr</span> GOLD
                  </div>
                  <span class="arrow-right">▶</span>
                </div>
                <div class="submenu">
                  <div class="sub-item" style="--i:1" @click="changeCurrency('GOLD', 'GRAM_24', 'GOLD 24K')">24K</div>
                  <div class="sub-item" style="--i:2" @click="changeCurrency('GOLD', 'GRAM_22', 'GOLD 22K')">22K</div>
                  <div class="sub-item" style="--i:3" @click="changeCurrency('GOLD', 'GRAM_18', 'GOLD 18K')">18K</div>
                  <div class="sub-item" style="--i:4" @click="changeCurrency('GOLD', 'CEYREK', 'ÇEYREK')">Çeyrek</div>
                  <div class="sub-item" style="--i:5" @click="changeCurrency('GOLD', 'YARIM', 'YARIM')">Yarım</div>
                  <div class="sub-item" style="--i:6" @click="changeCurrency('GOLD', 'TAM', 'TAM')">Tam</div>
                  <div class="sub-item" style="--i:7" @click="changeCurrency('GOLD', 'CUMHURIYET', 'CUMHURIYET')">Cumhuriyet</div>
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
        <div class="chart-line"></div>
      </div>

      <div class="assets-grid">
        <div v-for="asset in cardConfigs" :key="asset.type"
             class="asset-card" :class="'card-' + asset.type.toLowerCase()"
             @click="openModal(asset.type)">

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

      <div v-if="showModal" class="modal-overlay" @click.self="showModal = false">
        <div class="modal-content large-modal">
          <div class="modal-header">
            <h3>{{ activeAsset === 'GOLD' ? selectedVariantLabel : activeAsset }} İşlemleri</h3>
            <button class="close-btn" @click="showModal = false">✕</button>
          </div>

          <div class="modal-body-split">
            <div class="transaction-form">
              <h4>{{ t('newTransaction') }}</h4>

              <div v-if="activeAsset === 'GOLD'" class="ayar-wrapper">
                <label>Altın Tipi</label>
                <select v-model="selectedVariant" class="ayar-select" @change="updateVariantLabel">
                  <option value="GRAM_24">Gram Altın (24 Ayar)</option>
                  <option value="GRAM_22">Gram Altın (22 Ayar)</option>
                  <option value="GRAM_18">Gram Altın (18 Ayar)</option>
                  <option value="CEYREK">Çeyrek Altın</option>
                  <option value="YARIM">Yarım Altın</option>
                  <option value="TAM">Tam Altın</option>
                  <option value="CUMHURIYET">Cumhuriyet Altını</option>
                  <option value="GREMSE">Gremse Altın</option>
                </select>
              </div>

              <input v-model="amount" type="number" :placeholder="activeAsset === 'GOLD' ? 'Adet / Gram' : t('amount')" class="big-input" />
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
                    <span class="tx-type" :class="tx.type">
                      {{ tx.type === 'add' ? '+' : '-' }}
                      <span v-if="tx.asset_type === 'GOLD'" class="tx-variant">({{ formatVariant(tx.variant) }})</span>
                    </span>
                    <span class="tx-amount">{{ tx.amount }}</span>
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
  const selectedVariant = ref('GRAM_24');
  const selectedVariantLabel = ref('Gram Altın');
  const summaryData = ref(null);
  const allTransactions = ref([]);
  const baseCurrency = ref('TRY');
  const baseCurrencyLabel = ref('TRY');
  const targetVariant = ref('STANDARD');

  const cardConfigs = [
    { type: 'TRY', label: 'TRY', icon: '₺' },
    { type: 'USD', label: 'USD', icon: '$' },
    { type: 'EUR', label: 'EUR', icon: '€' },
    { type: 'BTC', label: 'BTC', icon: '₿' },
    { type: 'SILVER', label: 'SILVER', icon: 'Gr' },
    { type: 'GOLD', label: 'GOLD', icon: 'Gr' }
  ];

  const variantLabels = {
    'GRAM_24': '24K', 'GRAM_22': '22K', 'GRAM_18': '18K',
    'CEYREK': 'Çeyrek', 'YARIM': 'Yarım', 'TAM': 'Tam',
    'CUMHURIYET': 'Cumhuriyet', 'GREMSE': 'Gremse'
  };

  const toggleDropdown = () => { showSelector.value = !showSelector.value; };
  window.addEventListener('click', () => { if(showSelector.value) showSelector.value = false; });

  const displayCurrency = computed(() => baseCurrencyLabel.value);

  const totalValue = computed(() => {
    if (summaryData.value && summaryData.value.total_value) {
      // Eğer seçili birim Altın, Gümüş veya BTC ise 3 basamak, değilse 2 basamak
      const decimals = ['GOLD', 'SILVER', 'BTC'].includes(baseCurrency.value) ? 3 : 2;
      return summaryData.value.total_value.toLocaleString('tr-TR', { maximumFractionDigits: decimals, minimumFractionDigits: decimals });
    }
    return '0,00';
  });

  const hasData = computed(() => {
    return summaryData.value && summaryData.value.assets && summaryData.value.assets.length > 0 && summaryData.value.total_value > 0;
  });

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

  const fetchTransactions = async () => {
    try {
      const res = await api.get('/assets/transactions', {
        params: {
          currency: baseCurrency.value
        }
      });
      allTransactions.value = res.data || [];
    } catch (e) { console.error(e); }
  };

  const filteredTransactions = computed(() => {
    if (!activeAsset.value) return [];
    return allTransactions.value
            .filter(tx => tx.asset_type === activeAsset.value)
            .sort((a, b) => new Date(b.transaction_date) - new Date(a.transaction_date));
  });

  const getAmount = (type) => {
    const assets = summaryData.value?.assets || [];
    // O varlık tipindeki toplam miktar (Örn: Toplam Dolar)
    // Not: Altın için burası toplam gramı ifade etmeyebilir, backend dönüşüne bağlı.
    // Ancak orijinal tasarımda miktar göstermek istiyordun.
    const total = assets.filter(a => a.type === type).reduce((sum, curr) => sum + curr.amount, 0);
    return total.toLocaleString('tr-TR', { maximumFractionDigits: 3 });
  };

  const getAllocation = (type) => {
    const assets = summaryData.value?.assets || [];
    const total = assets.filter(a => a.type === type).reduce((sum, curr) => sum + (curr.allocation || 0), 0);
    return total.toFixed(1);
  };

  const openModal = (asset) => {
    activeAsset.value = asset;
    amount.value = '';
    selectedVariant.value = 'GRAM_24';
    updateVariantLabel();
    transactionDate.value = new Date().toISOString().split('T')[0];
    showModal.value = true;
    fetchTransactions();
  };

  const updateVariantLabel = () => {
    if (activeAsset.value === 'GOLD') {
      selectedVariantLabel.value = variantLabels[selectedVariant.value] || selectedVariant.value;
    } else {
      selectedVariantLabel.value = activeAsset.value;
    }
  };

  const formatVariant = (v) => {
    if (v === 'STANDARD') return '';
    return variantLabels[v] || v;
  }

  const handleTransaction = async (action) => {
    if (!amount.value) return alert(t('noData'));
    try {
      await api.post('/assets/update', {
        type: activeAsset.value,
        variant: activeAsset.value === 'GOLD' ? selectedVariant.value : 'STANDARD',
        amount: parseFloat(amount.value),
        action,
        transaction_date: new Date(transactionDate.value).toISOString()
      });
      showModal.value = false;
      fetchData();
      fetchTransactions();
      amount.value = '';
    } catch(e) { alert("Hata: " + (e.response?.data?.error || e.message)); }
  };

  const changeCurrency = (code, variant, label) => {
    baseCurrency.value = code;
    targetVariant.value = variant;
    baseCurrencyLabel.value = label;
    showSelector.value = false;
    fetchData();
    fetchTransactions();
  };

  const downloadReceipt = async () => {
    try {
      const res = await api.get('/assets/receipt/full', {
        responseType: 'blob',
        params: { currency: baseCurrency.value }
      });
      const url = window.URL.createObjectURL(new Blob([res.data]));
      const link = document.createElement('a');
      link.href = url;
      link.setAttribute('download', 'Genel_Rapor.pdf');
      document.body.appendChild(link);
      link.click();
    } catch (e) { alert('PDF İndirilemedi. Sunucu hatası.'); }
  };

  const downloadExcel = async () => {
    try {
      const res = await api.get('/assets/export/excel', {
        responseType: 'blob',
        params: { currency: baseCurrency.value }
      });
      const url = window.URL.createObjectURL(new Blob([res.data]));
      const link = document.createElement('a');
      link.href = url;
      link.setAttribute('download', `FinTrack_Export_${new Date().toLocaleDateString()}.xlsx`);
      document.body.appendChild(link);
      link.click();
    } catch (e) { alert('Excel İndirilemedi. Sunucu hatası.'); }
  };

  const downloadSingleReceipt = async (id) => {
    try {
      const res = await api.get(`/assets/receipt/${id}`, {
        responseType: 'blob',
        params: { currency: baseCurrency.value }
      });
      const url = window.URL.createObjectURL(new Blob([res.data]));
      const link = document.createElement('a');
      link.href = url;
      link.setAttribute('download', `Islem_Dekontu_${id}.pdf`);
      document.body.appendChild(link);
      link.click();
    } catch (e) { alert('PDF Hatası. İşlem bulunamadı.'); }
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

  .top-bar { display: flex; justify-content: space-between; align-items: center; margin-bottom: 30px; flex-wrap: wrap; gap: 20px; }
  .page-title h2 { color: var(--text-color); margin: 0; font-size: 1.8rem; font-weight: 700; }

  /* SAĞ ÜST KÖŞE: DROP DOWN STİLİ */
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
    min-width: 200px; /* Genişlettim ki değer sığsın */
    box-shadow: 0 4px 10px rgba(0,0,0,0.2);
  }
  .currency-btn:hover { border-color: var(--accent-color); transform: translateY(-2px); }
  .currency-btn.active .arrow-icon { transform: rotate(180deg); }
  .arrow-icon { font-size: 0.8rem; transition: 0.3s; }

  .currency-dropdown {
    position: absolute;
    top: calc(100% + 10px);
    inset-inline-end: 0;
    background: var(--sidebar-bg);
    border: 1px solid var(--border-color);
    border-radius: 12px;
    width: 220px;
    box-shadow: 0 10px 30px rgba(0,0,0,0.5);
    z-index: 101;
    padding: 8px 0;
    overflow: visible;
  }

  .c-item {
    padding: 12px 20px;
    cursor: pointer;
    position: relative;
    color: var(--text-color);
    transition: 0.2s;
    font-weight: 600;
    display: flex;
    align-items: center;
    gap: 12px;
    border-bottom: 1px solid var(--border-color);
  }
  .c-item:last-child { border-bottom: none; }
  .c-item:hover { background: var(--hover-bg); color: var(--accent-color); }

  .symbol { width: 24px; text-align: center; font-weight: bold; font-size: 1.1rem; }
  .text-icon { font-size: 0.85rem; font-weight: 800; opacity: 0.9; }

  .has-submenu { position: relative; overflow: visible; }
  .gold-trigger { display: flex; align-items: center; justify-content: space-between; width: 100%; }

  .submenu {
    position: absolute;
    top: 0;
    inset-inline-end: 100%;
    margin-inline-end: -1px;
    background: var(--sidebar-bg);
    border: 1px solid var(--border-color);
    border-radius: 12px;
    min-width: 140px;
    padding: 8px;
    box-shadow: 0 10px 30px rgba(0,0,0,0.5);
    opacity: 0;
    visibility: hidden;
    transform: translateX(10px);
    transition: all 0.2s ease;
  }

  html[dir="ltr"] .submenu { transform: translateX(-10px); }
  html[dir="rtl"] .submenu { transform: translateX(10px); margin-inline-end: -1px; }

  .has-submenu::after { content: ''; position: absolute; top: 0; bottom: 0; width: 20px; inset-inline-end: 100%; z-index: 102; }
  .has-submenu:hover .submenu { opacity: 1; visibility: visible; transform: translateX(0); }

  .sub-item {
    padding: 10px 15px;
    font-size: 0.9rem;
    color: var(--text-color);
    cursor: pointer;
    transition: 0.2s;
    border-bottom: 1px solid var(--border-color);
    opacity: 0;
    animation: slideIn 0.2s forwards;
    animation-delay: calc(var(--i) * 0.03s);
  }
  .sub-item:last-child { border-bottom: none; }
  .sub-item:hover { background: var(--hover-bg); color: var(--accent-color); }

  @keyframes slideIn { from { opacity: 0; transform: translateX(-5px); } to { opacity: 1; transform: translateX(0); } }

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

  /* FLOATING ACTIONS */
  .floating-actions { position: fixed; bottom: 30px; inset-inline-end: 30px; display: flex; flex-direction: column; gap: 15px; z-index: 110; }
  .f-btn { width: 60px; height: 60px; border-radius: 50%; border: none; font-size: 24px; cursor: pointer; box-shadow: 0 5px 15px rgba(0,0,0,0.3); transition: 0.3s; color: white; display: flex; align-items: center; justify-content: center; }
  .f-btn:hover { transform: scale(1.1); }
  .excel { background: linear-gradient(135deg, #10B981, #059669); }
  .pdf { background: linear-gradient(135deg, #3B82F6, #2563EB); }

  /* MODAL STYLES */
  .modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.85); display: flex; justify-content: center; align-items: center; z-index: 200; backdrop-filter: blur(5px); }
  .large-modal { width: 800px !important; max-width: 95%; padding: 40px; border-radius: 20px; }
  .modal-content { background: var(--card-bg); padding: 30px; border-radius: 20px; border: 1px solid var(--border-color); box-shadow: 0 20px 50px rgba(0,0,0,0.5); }
  .modal-header { display: flex; justify-content: space-between; margin-bottom: 25px; border-bottom: 1px solid var(--border-color); padding-bottom: 15px; }
  .modal-header h3 { color: var(--text-color); margin: 0; font-size: 1.5rem; }
  /* ÇARPI BUTONU İÇİN STİL */
  .close-btn { background: none; border: none; color: var(--text-muted); font-size: 1.5rem; cursor: pointer; transition: 0.2s; }
  .close-btn:hover { color: var(--danger-color); transform: rotate(90deg); }

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
  .tx-variant { font-size: 0.8em; opacity: 0.7; margin-left: 5px; } /* Yeni eklenen variant stili */
  .receipt-download-btn { background: transparent; border: 1px solid var(--border-color); border-radius: 8px; cursor: pointer; padding: 8px; transition: 0.2s; font-size: 1.2rem; }
  .receipt-download-btn:hover { background: var(--hover-bg); }

  .big-input, .ayar-select, .date-input { width: 100%; padding: 15px; background: var(--input-bg); border: 1px solid var(--border-color); color: var(--text-color); border-radius: 12px; margin-bottom: 15px; font-size: 1rem; box-sizing: border-box; }
  .ayar-select { cursor: pointer; }
  .ayar-wrapper label { display: block; font-size: 0.85rem; color: var(--text-muted); margin-bottom: 5px; }

  .actions { display: flex; gap: 15px; margin-top: 10px; }
  .actions button { flex: 1; padding: 15px; border-radius: 12px; border: none; font-weight: bold; cursor: pointer; color: white; transition: 0.2s; font-size: 1rem; }
  .add { background: linear-gradient(135deg, #10B981, #047857); }
  .sub { background: linear-gradient(135deg, #EF4444, #B91C1C); }
  .add:hover, .sub:hover { opacity: 0.9; transform: translateY(-2px); }

  @media (max-width: 1024px) { .assets-grid { grid-template-columns: repeat(2, 1fr); max-width: 100%; } }
  @media (max-width: 768px) {
    .assets-grid { grid-template-columns: 1fr; }
    .modal-body-split { grid-template-columns: 1fr; }
    .transaction-history { border-inline-start: none; border-top: 1px solid var(--border-color); padding-inline-start: 0; padding-top: 30px; }
    .chart-wrapper { width: 250px; height: 250px; }
    .center-balance h3 { font-size: 1.5rem; }
  }
</style>