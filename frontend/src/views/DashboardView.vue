<template>
  <div class="dashboard-container">
    <header class="header-section">
      <div class="title-area">
        <h1>Portföy Özeti</h1>
      </div>

      <div class="controls-area">
        <div v-if="baseCurrency === 'GOLD'" class="gold-bar">
           <span
                   v-for="k in [24, 22, 18, 14, 8]"
                   :key="k"
                   class="gold-pill"
                   :class="{ active: targetAyar === k }"
                   @click="selectGold(k)"
           >
             {{ k }}K
           </span>
        </div>

        <div class="currency-dropdown">
          <div class="selected-currency" @click="toggleMenu">
            <span>{{ baseCurrency }}</span>
            <span class="arrow" :class="{ rotated: isMenuOpen }">▼</span>
          </div>

          <div v-if="isMenuOpen" class="dropdown-list">
            <div
                    v-for="item in currencyItems"
                    :key="item.code"
                    class="dropdown-item"
                    :class="{ active: baseCurrency === item.code }"
                    @click="selectCurrency(item)"
            >
              {{ item.label }}
            </div>
          </div>
        </div>
      </div>
    </header>

    <div class="main-content">
      <div class="chart-wrapper">
        <div class="chart-container">
          <Doughnut v-if="hasData" :data="chartData" :options="chartOptions" />
          <div v-else class="no-data-circle">
            <div class="dashed-circle"></div>
            <span>Veri Yok</span>
          </div>
        </div>
      </div>
    </div>

    <div class="cards-grid">
      <div class="asset-card card-try">
        <div class="card-icon">₺</div>
        <div class="card-info">
          <h3>TRY</h3>
          <p>0 ₺</p>
          <span class="rate">%0.0</span>
        </div>
      </div>

      <div class="asset-card card-usd">
        <div class="card-icon">$</div>
        <div class="card-info">
          <h3>USD</h3>
          <p>0 $</p>
          <span class="rate">%0.0</span>
        </div>
      </div>

      <div class="asset-card card-eur">
        <div class="card-icon">€</div>
        <div class="card-info">
          <h3>EUR</h3>
          <p>0 €</p>
          <span class="rate">%0.0</span>
        </div>
      </div>

      <div class="asset-card card-btc">
        <div class="card-icon">₿</div>
        <div class="card-info">
          <h3>BTC</h3>
          <p>0 BTC</p>
          <span class="rate">%0.0</span>
        </div>
      </div>

      <div class="asset-card card-silver">
        <div class="card-icon">⚔️</div>
        <div class="card-info">
          <h3>SILVER</h3>
          <p>0 Gr</p>
          <span class="rate">%0.0</span>
        </div>
      </div>

      <div class="asset-card card-gold">
        <div class="card-icon">👑</div>
        <div class="card-info">
          <h3>GOLD (Gr)</h3>
          <p>0 Gr</p>
          <span class="rate">%0.0</span>
        </div>
      </div>
    </div>

    <div class="fab-container">
      <button class="fab-btn chart-btn">📊</button>
      <button class="fab-btn file-btn">📄</button>
    </div>
  </div>
</template>

<script setup>
  import { ref, computed, onMounted } from 'vue'
  import { Doughnut } from 'vue-chartjs'
  import { Chart as ChartJS, ArcElement, Tooltip, Legend } from 'chart.js'
  import api from '../services/api'

  ChartJS.register(ArcElement, Tooltip, Legend)

  const baseCurrency = ref('TRY')
  const targetAyar = ref(24) // Varsayılan 24K
  const summaryData = ref(null)
  const isMenuOpen = ref(false)

  const currencyItems = [
    { code: 'TRY', label: 'TRY' },
    { code: 'USD', label: 'USD' },
    { code: 'EUR', label: 'EUR' },
    { code: 'BTC', label: 'BTC' },
    { code: 'SILVER', label: 'SILVER' },
    { code: 'GOLD', label: 'GOLD (Gr)' }
  ]

  const toggleMenu = () => {
    isMenuOpen.value = !isMenuOpen.value
  }

  const selectCurrency = (item) => {
    baseCurrency.value = item.code
    if (item.code !== 'GOLD') targetAyar.value = 0
    else targetAyar.value = 24
    isMenuOpen.value = false
    fetchData()
  }

  const selectGold = (k) => {
    targetAyar.value = k
    fetchData()
  }

  const fetchData = async () => {
    try {
      const res = await api.get('/assets/summary', {
        headers: {
          'X-Currency': baseCurrency.value,
          'X-Ayar': targetAyar.value.toString()
        }
      })
      summaryData.value = res.data
    } catch (e) {
      console.log("Veri çekilemedi, mock data kullanılıyor")
    }
  }

  const hasData = computed(() => summaryData.value?.total_value > 0)

  const chartData = computed(() => ({
    labels: ['TRY','USD','EUR','BTC','SILVER','GOLD'],
    datasets: [{
      data: [10, 20, 15, 5, 25, 25],
      backgroundColor: ['#ef4444', '#10b981', '#8b5cf6', '#111111', '#9ca3af', '#facc15'],
      borderWidth: 0,
      hoverOffset: 15
    }]
  }))

  const chartOptions = {
    responsive: true,
    cutout: '80%',
    plugins: { legend: { display: false } }
  }

  onMounted(fetchData)
</script>

<style scoped>
  .dashboard-container {
    min-height: 100vh;
    background: radial-gradient(circle at top left, #1e1b4b, #020617);
    padding: 40px;
    color: white;
    position: relative;
    font-family: 'Inter', sans-serif;
  }

  .header-section {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
    flex-wrap: wrap;
    gap: 20px;
  }

  .header-section h1 {
    font-size: 2rem;
    font-weight: 800;
    margin: 0;
    letter-spacing: -0.02em;
  }

  .controls-area {
    display: flex;
    align-items: center;
    gap: 15px;
  }

  /* GOLD BAR AYARLARI */
  .gold-bar {
    display: flex;
    background: rgba(15, 23, 42, 0.6);
    padding: 4px;
    border-radius: 12px;
    border: 1px solid rgba(255, 255, 255, 0.1);
  }

  .gold-pill {
    padding: 8px 12px;
    border-radius: 8px;
    font-size: 0.85rem;
    font-weight: 700;
    cursor: pointer;
    color: #94a3b8;
    transition: all 0.2s;
  }

  .gold-pill.active {
    background: #facc15;
    color: #000;
    box-shadow: 0 0 10px rgba(250, 204, 21, 0.4);
  }

  /* DROPDOWN AYARLARI */
  .currency-dropdown {
    position: relative;
    width: 140px;
  }

  .selected-currency {
    background: rgba(15, 23, 42, 0.6);
    backdrop-filter: blur(10px);
    border: 1px solid rgba(255, 255, 255, 0.1);
    padding: 12px 16px;
    border-radius: 12px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    cursor: pointer;
    font-weight: 700;
    transition: border 0.2s;
  }

  .selected-currency:hover {
    border-color: rgba(255,255,255,0.3);
  }

  .arrow {
    font-size: 0.8rem;
    transition: transform 0.2s;
  }

  .arrow.rotated {
    transform: rotate(180deg);
  }

  .dropdown-list {
    position: absolute;
    top: 110%;
    right: 0;
    width: 100%;
    background: #0f172a;
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 12px;
    overflow: hidden;
    z-index: 50;
    box-shadow: 0 10px 30px rgba(0,0,0,0.5);
  }

  .dropdown-item {
    padding: 10px 16px;
    cursor: pointer;
    font-size: 0.9rem;
    color: #cbd5e1;
    transition: background 0.2s;
  }

  .dropdown-item:hover, .dropdown-item.active {
    background: rgba(255,255,255,0.1);
    color: white;
  }

  /* CHART VE KARTLAR */
  .main-content {
    display: flex;
    justify-content: center;
    height: 300px;
    margin-bottom: 40px;
  }

  .chart-container {
    width: 280px;
    height: 280px;
    position: relative;
  }

  .no-data-circle {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    height: 100%;
    color: #64748b;
  }

  .dashed-circle {
    width: 200px;
    height: 200px;
    border: 4px dashed #334155;
    border-radius: 50%;
    margin-bottom: 10px;
  }

  .cards-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
    gap: 20px;
    padding-bottom: 80px;
  }

  .asset-card {
    padding: 20px;
    border-radius: 20px;
    color: white;
    position: relative;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    min-height: 120px;
    transition: transform 0.2s;
    box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  }

  .asset-card:hover { transform: translateY(-5px); }
  .card-icon { font-size: 1.5rem; margin-bottom: 10px; }
  .card-info h3 { font-size: 0.8rem; opacity: 0.8; margin: 0; font-weight: 600; }
  .card-info p { font-size: 1.25rem; font-weight: 800; margin: 4px 0; }
  .rate { font-size: 0.8rem; opacity: 0.7; }

  /* RENKLER */
  .card-try { background-color: #ef4444; }
  .card-usd { background-color: #10b981; }
  .card-eur { background-color: #57534e; }
  .card-btc { background-color: #1c1917; }
  .card-silver { background-color: #64748b; }
  .card-gold { background-color: #ca8a04; color: white; }

  .fab-container {
    position: fixed; bottom: 30px; right: 30px;
    display: flex; flex-direction: column; gap: 15px; z-index: 40;
  }
  .fab-btn {
    width: 50px; height: 50px; border-radius: 50%; border: none;
    cursor: pointer; font-size: 1.2rem;
    box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.5);
    transition: transform 0.2s;
  }
  .fab-btn:hover { transform: scale(1.1); }
  .chart-btn { background: #10b981; color: white; }
  .file-btn { background: #3b82f6; color: white; }

  @media (max-width: 768px) {
    .header-section { flex-direction: column; align-items: flex-start; }
    .controls-area { width: 100%; justify-content: space-between; }
  }
</style>