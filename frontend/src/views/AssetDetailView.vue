<template>
    <div class="detail-container" :style="pageTheme">
        <header class="detail-header">
            <div class="left">
                <button class="back-btn" @click="$router.push('/dashboard')">← Geri</button>
                <h1>{{ type === 'GOLD' ? 'GOLD' : type }} Portföyü</h1>
            </div>

            <div class="right-menu">
                <div class="currency-dropdown-wrapper" v-click-outside="closeDropdown">
                    <button class="dd-btn" @click="showDd = !showDd">
                        <span class="unit">{{ displayLabel }}</span>
                        <span class="arrow">▼</span>
                    </button>

                    <transition name="dd-anim">
                        <div v-if="showDd" class="dd-menu">
                            <template v-if="type === 'GOLD'">
                                <div class="dd-category">Gram Altınlar</div>
                                <div class="dd-item" @click="changeViewUnit('GOLD', 'GRAM_24', 'Has (24K)')">Has (24K)</div>
                                <div class="dd-item" @click="changeViewUnit('GOLD', 'GRAM_22', '22 Ayar')">22 Ayar</div>
                                <div class="dd-item" @click="changeViewUnit('GOLD', 'GRAM_18', '18 Ayar')">18 Ayar</div>
                                <div class="dd-item" @click="changeViewUnit('GOLD', 'GRAM_14', '14 Ayar')">14 Ayar</div>
                                <div class="dd-item" @click="changeViewUnit('GOLD', 'GRAM_8', '8 Ayar')">8 Ayar</div>
                                <div class="dd-item" @click="changeViewUnit('GOLD', 'GRAM_4', '4 Ayar')">4 Ayar</div>

                                <div class="dd-divider"></div>
                                <div class="dd-category">Ziynet Altınlar</div>
                                <div class="dd-item" @click="changeViewUnit('GOLD', 'CEYREK', 'Çeyrek')">Çeyrek</div>
                                <div class="dd-item" @click="changeViewUnit('GOLD', 'YARIM', 'Yarım')">Yarım</div>
                                <div class="dd-item" @click="changeViewUnit('GOLD', 'TAM', 'Tam')">Tam</div>
                                <div class="dd-item" @click="changeViewUnit('GOLD', 'CUMHURIYET', 'Cumhuriyet')">Cumhuriyet</div>
                                <div class="dd-item" @click="changeViewUnit('GOLD', 'GREMSE', 'Gremse')">Gremse</div>
                            </template>

                            <template v-else>
                                <div class="dd-item" @click="changeViewUnit('TRY', 'STANDARD', 'TRY')">₺ TRY</div>
                                <div class="dd-item" @click="changeViewUnit('USD', 'STANDARD', 'USD')">$ USD</div>
                                <div class="dd-item" @click="changeViewUnit('EUR', 'STANDARD', 'EUR')">€ EUR</div>
                            </template>
                        </div>
                    </transition>
                </div>

                <div class="actions">
                    <button class="icon-btn" @click="downloadExcel" title="Excel Raporu">📊</button>
                    <button class="icon-btn" @click="downloadFullPDF" title="PDF Raporu">📄</button>
                </div>
            </div>
        </header>

        <div class="hero-section">
            <div class="hero-value">{{ totalInSelectedUnit }}</div>
            <div class="hero-unit">{{ displayLabel }}</div>
            <div class="hero-line"></div>
        </div>

        <div class="detail-grid">
            <section class="action-card">
                <div class="card-header">
                    <h3>Hızlı İşlem</h3>
                </div>

                <div class="form-body">
                    <div v-if="type === 'GOLD'" class="gold-options">
                        <div class="input-group">
                            <label>Altın Türü</label>
                            <select v-model="form.goldType" class="custom-input">
                                <option value="GRAM">Gram Altın</option>
                                <option value="CEYREK">Çeyrek Altın</option>
                                <option value="YARIM">Yarım Altın</option>
                                <option value="TAM">Tam Altın</option>
                                <option value="CUMHURIYET">Cumhuriyet</option>
                                <option value="GREMSE">Gremse</option>
                            </select>
                        </div>

                        <div v-if="form.goldType === 'GRAM'" class="input-group fade-in">
                            <label>Ayar (Karat)</label>
                            <select v-model="form.goldKarat" class="custom-input">
                                <option :value="24">24 Ayar (Has)</option>
                                <option :value="22">22 Ayar</option>
                                <option :value="18">18 Ayar</option>
                                <option :value="14">14 Ayar</option>
                                <option :value="8">8 Ayar</option>
                                <option :value="4">4 Ayar</option>
                            </select>
                        </div>
                    </div>

                    <div class="input-row">
                        <div class="input-group">
                            <label>Miktar</label>
                            <input v-model="form.amount" type="number" placeholder="0.00" class="custom-input">
                        </div>
                        <div class="input-group">
                            <label>Tarih</label>
                            <input v-model="form.date" type="date" class="custom-input">
                        </div>
                    </div>

                    <div class="action-buttons">
                        <button class="btn-add" @click="submitTx('add')">
                            <span>+</span> EKLE
                        </button>
                        <button class="btn-sub" @click="submitTx('subtract')">
                            <span>-</span> ÇIKAR
                        </button>
                    </div>
                </div>
            </section>

            <section class="history-card">
                <div class="card-header">
                    <h3>İşlem Geçmişi</h3>
                </div>
                <div class="history-list">
                    <div v-if="transactions.length === 0" class="no-data">Henüz işlem yapılmadı.</div>
                    <div v-for="tx in transactions" :key="tx.id" class="tx-item">
                        <div class="tx-left">
                            <div class="tx-date">{{ new Date(tx.transaction_date).toLocaleDateString('tr-TR') }}</div>
                            <div class="tx-desc">
                                <span class="dot" :class="tx.type">●</span>
                                <span class="tx-name">{{ formatTxName(tx) }}</span>
                            </div>
                        </div>
                        <div class="tx-right">
              <span class="tx-amount" :class="tx.type">
                {{ tx.type === 'add' ? '+' : '-' }}{{ tx.amount }}
              </span>
                            <button class="receipt-btn" @click="downloadSingleReceipt(tx.id)" title="Dekont">📄</button>
                        </div>
                    </div>
                </div>
            </section>
        </div>
    </div>
</template>

<script setup>
    import { ref, computed, onMounted, watch } from 'vue';
    import api from '../services/api';

    const props = defineProps(['type']);
    const assets = ref([]);
    const transactions = ref([]);
    const showDd = ref(false);

    const selectedCurrency = ref(props.type === 'GOLD' ? 'GOLD' : 'TRY');
    const selectedVariant = ref(props.type === 'GOLD' ? 'GRAM_24' : 'STANDARD');
    const displayLabel = ref(props.type === 'GOLD' ? 'Has (24K)' : 'TRY');

    const form = ref({
        amount: '',
        date: new Date().toISOString().split('T')[0],
        goldType: 'GRAM',
        goldKarat: 24
    });

    const themeMap = {
        TRY: 'linear-gradient(135deg, #1e0505 0%, #450a0a 100%)',
        USD: 'linear-gradient(135deg, #051e0f 0%, #064e3b 100%)',
        GOLD: 'linear-gradient(135deg, #2A1C05 0%, #713F12 100%)',
        EUR: 'linear-gradient(135deg, #050b1e 0%, #1e3a8a 100%)',
        BTC: 'linear-gradient(135deg, #1e1305 0%, #78350f 100%)',
        SILVER: 'linear-gradient(135deg, #111827 0%, #374151 100%)'
    };

    const pageTheme = computed(() => ({ background: themeMap[props.type] || '#000' }));

    const totalInSelectedUnit = computed(() => {
        const sum = assets.value.reduce((acc, curr) => acc + curr.value_in_base, 0);
        return sum.toLocaleString('tr-TR', { minimumFractionDigits: 2, maximumFractionDigits: 2 });
    });

    const fetchData = async () => {
        try {
            const [resSum, resTx] = await Promise.all([
                api.get('/assets/summary', { params: { currency: selectedCurrency.value, variant: selectedVariant.value } }),
                api.get('/assets/transactions', { params: { currency: selectedCurrency.value } })
            ]);
            assets.value = (resSum.data.assets || []).filter(a => a.type === props.type);
            transactions.value = (resTx.data || []).filter(t => t.asset_type === props.type);
        } catch (err) { console.error(err); }
    };

    const changeViewUnit = (code, variant, label) => {
        selectedCurrency.value = code;
        selectedVariant.value = variant;
        displayLabel.value = label;
        showDd.value = false;
        fetchData();
    };

    const formatTxName = (tx) => {
        if (tx.asset_type !== 'GOLD') return props.type;

        if (tx.variant.startsWith('GRAM')) {
            const k = tx.variant.split('_')[1];
            return `${k} Ayar Gram`;
        }

        const map = {
            'CEYREK': 'Çeyrek Altın',
            'YARIM': 'Yarım Altın',
            'TAM': 'Tam Altın',
            'CUMHURIYET': 'Cumhuriyet',
            'GREMSE': 'Gremse'
        };
        return map[tx.variant] || tx.variant;
    };

    const submitTx = async (action) => {
        if (!form.value.amount || form.value.amount <= 0) return alert("Geçerli miktar giriniz");

        let variantToSend = 'STANDARD';
        if (props.type === 'GOLD') {
            if (form.value.goldType === 'GRAM') {
                variantToSend = `GRAM_${form.value.goldKarat}`;
            } else {
                variantToSend = form.value.goldType;
            }
        }

        try {
            await api.post('/assets/update', {
                type: props.type,
                action,
                variant: variantToSend,
                amount: parseFloat(form.value.amount),
                transaction_date: new Date(form.value.date).toISOString()
            });
            form.value.amount = '';
            await fetchData();
        } catch (err) { alert(err.response?.data?.error || err.message); }
    };

    const downloadSingleReceipt = async (id) => {
        const res = await api.get(`/assets/receipt/${id}`, { responseType: 'blob', params: { currency: 'TRY' } });
        const url = window.URL.createObjectURL(new Blob([res.data]));
        const link = document.createElement('a'); link.href = url; link.setAttribute('download', `Dekont_${id}.pdf`); link.click();
    };

    const downloadFullPDF = async () => {
        const res = await api.get(`/assets/receipt/full`, { responseType: 'blob', params: { currency: 'TRY' } });
        const url = window.URL.createObjectURL(new Blob([res.data]));
        const link = document.createElement('a'); link.href = url; link.setAttribute('download', 'Rapor.pdf'); link.click();
    };

    const downloadExcel = async () => {
        const res = await api.get(`/assets/export/excel`, { responseType: 'blob', params: { currency: 'TRY' } });
        const url = window.URL.createObjectURL(new Blob([res.data]));
        const link = document.createElement('a'); link.href = url; link.setAttribute('download', 'Export.xlsx'); link.click();
    };

    const closeDropdown = () => { showDd.value = false; };
    const vClickOutside = {
        mounted(el, binding) {
            el.clickOutsideEvent = function(event) {
                if (!(el === event.target || el.contains(event.target))) { binding.value(event, el); }
            };
            document.body.addEventListener('click', el.clickOutsideEvent);
        },
        unmounted(el) { document.body.removeEventListener('click', el.clickOutsideEvent); }
    };

    watch(() => props.type, () => {
        selectedCurrency.value = props.type === 'GOLD' ? 'GOLD' : 'TRY';
        selectedVariant.value = props.type === 'GOLD' ? 'GRAM_24' : 'STANDARD';
        displayLabel.value = props.type === 'GOLD' ? 'Has (24K)' : 'TRY';
        fetchData();
    });

    onMounted(fetchData);
</script>

<style scoped>
    .detail-container { min-height: 100vh; padding: 40px; color: white; font-family: 'Inter', sans-serif; box-sizing: border-box; }
    .detail-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 30px; }
    .back-btn { background: rgba(255,255,255,0.1); border: none; color: white; padding: 10px 25px; border-radius: 30px; cursor: pointer; font-weight: bold; transition: 0.3s; }
    .back-btn:hover { background: rgba(255,255,255,0.2); }

    .right-menu { display: flex; align-items: center; gap: 15px; }
    .currency-dropdown-wrapper { position: relative; }
    .dd-btn { background: rgba(255,255,255,0.1); border: 1px solid rgba(255,255,255,0.2); color: white; padding: 10px 20px; border-radius: 12px; cursor: pointer; font-weight: 700; font-size: 1rem; min-width: 160px; display: flex; justify-content: space-between; align-items: center; }
    .dd-btn:hover { background: rgba(255,255,255,0.15); border-color: #FFD700; }

    .dd-menu {
        position: absolute; top: 100%; right: 0;
        background: #1E293B; border: 1px solid #334155;
        border-radius: 12px; width: 220px; margin-top: 10px;
        z-index: 100; max-height: 400px; overflow-y: auto; overflow-x: hidden;
        box-shadow: 0 10px 40px rgba(0,0,0,0.5);
    }
    .dd-menu::-webkit-scrollbar { width: 6px; }
    .dd-menu::-webkit-scrollbar-thumb { background: #475569; border-radius: 3px; }
    .dd-category { font-size: 0.75rem; color: #94A3B8; padding: 10px 15px 5px; font-weight: bold; letter-spacing: 1px; }
    .dd-divider { height: 1px; background: rgba(255,255,255,0.1); margin: 5px 0; }
    .dd-item { padding: 12px 15px; cursor: pointer; border-bottom: 1px solid rgba(255,255,255,0.05); transition: 0.2s; display: flex; align-items: center; gap: 10px; font-size: 0.95rem; }
    .dd-item:hover { background: #334155; color: #FFD700; padding-left: 20px; }
    .dd-anim-enter-active, .dd-anim-leave-active { transition: all 0.2s ease; }
    .dd-anim-enter-from, .dd-anim-leave-to { opacity: 0; transform: translateY(-10px); }

    .actions { display: flex; gap: 10px; }
    .icon-btn { background: rgba(255,255,255,0.1); border: none; width: 45px; height: 45px; border-radius: 10px; cursor: pointer; font-size: 1.2rem; display: flex; align-items: center; justify-content: center; transition: 0.3s; }
    .icon-btn:hover { background: rgba(255,255,255,0.2); transform: scale(1.1); }

    /* HERO SECTION */
    .hero-section { text-align: center; margin-bottom: 50px; }
    .hero-value { font-size: 4rem; font-weight: 800; letter-spacing: -2px; line-height: 1; }
    .hero-unit { font-size: 1.2rem; color: #FFD700; font-weight: bold; margin-top: 5px; text-transform: uppercase; letter-spacing: 2px; opacity: 0.9; }
    .hero-line { width: 100px; height: 4px; background: #FFD700; margin: 20px auto 0; border-radius: 2px; opacity: 0.7; }

    /* GRID */
    .detail-grid { display: grid; grid-template-columns: 1fr 1.5fr; gap: 40px; max-width: 1200px; margin: 0 auto; }
    .action-card, .history-card { background: rgba(0,0,0,0.3); backdrop-filter: blur(20px); border-radius: 24px; border: 1px solid rgba(255,255,255,0.1); overflow: hidden; display: flex; flex-direction: column; }
    .card-header { padding: 20px 30px; border-bottom: 1px solid rgba(255,255,255,0.1); background: rgba(255,255,255,0.02); }
    .card-header h3 { margin: 0; font-size: 1.2rem; color: #fff; font-weight: 700; }

    /* FORM */
    .form-body { padding: 30px; display: flex; flex-direction: column; gap: 20px; }
    .input-group label { display: block; margin-bottom: 8px; color: #94A3B8; font-size: 0.9rem; font-weight: 600; }
    .custom-input { width: 100%; padding: 16px; background: rgba(0,0,0,0.4); border: 1px solid rgba(255,255,255,0.15); color: white; border-radius: 12px; font-size: 1.1rem; box-sizing: border-box; transition: 0.3s; }
    .custom-input:focus { border-color: #FFD700; outline: none; background: rgba(0,0,0,0.6); }
    .input-row { display: grid; grid-template-columns: 1fr 1fr; gap: 20px; }
    .gold-options { display: flex; flex-direction: column; gap: 20px; border-bottom: 1px solid rgba(255,255,255,0.1); padding-bottom: 20px; margin-bottom: 10px; }
    .fade-in { animation: fadeIn 0.3s ease; }
    @keyframes fadeIn { from { opacity: 0; transform: translateY(-5px); } to { opacity: 1; transform: translateY(0); } }

    .action-buttons { display: flex; gap: 20px; margin-top: 10px; }
    .btn-add, .btn-sub { flex: 1; padding: 18px; border-radius: 16px; font-weight: 800; cursor: pointer; border: none; font-size: 1rem; display: flex; align-items: center; justify-content: center; gap: 10px; transition: 0.3s; color: white; }
    .btn-add { background: linear-gradient(135deg, #10B981, #059669); box-shadow: 0 4px 15px rgba(16, 185, 129, 0.3); }
    .btn-sub { background: linear-gradient(135deg, #EF4444, #B91C1C); box-shadow: 0 4px 15px rgba(239, 68, 68, 0.3); }
    .btn-add:hover, .btn-sub:hover { transform: translateY(-2px); filter: brightness(1.1); }

    /* HISTORY */
    .history-list { padding: 0 30px; max-height: 500px; overflow-y: auto; }
    .history-list::-webkit-scrollbar { width: 6px; }
    .history-list::-webkit-scrollbar-thumb { background: rgba(255,255,255,0.1); border-radius: 3px; }
    .tx-item { display: flex; justify-content: space-between; align-items: center; padding: 20px 0; border-bottom: 1px solid rgba(255,255,255,0.05); }
    .tx-left { display: flex; flex-direction: column; gap: 5px; }
    .tx-date { font-size: 0.8rem; color: #94A3B8; }
    .tx-desc { display: flex; align-items: center; gap: 10px; font-weight: 600; font-size: 1.05rem; }
    .dot { font-size: 0.8rem; }
    .dot.add { color: #10B981; } .dot.subtract { color: #EF4444; }
    .tx-right { display: flex; align-items: center; gap: 15px; }
    .tx-amount { font-weight: 800; font-size: 1.1rem; }
    .tx-amount.add { color: #10B981; } .tx-amount.subtract { color: #EF4444; }
    .receipt-btn { background: rgba(255,255,255,0.05); border: 1px solid rgba(255,255,255,0.1); width: 40px; height: 40px; border-radius: 10px; cursor: pointer; transition: 0.2s; font-size: 1.2rem; }
    .receipt-btn:hover { background: rgba(255,255,255,0.15); }
    .no-data { text-align: center; padding: 50px; color: #94A3B8; font-style: italic; }

    @media (max-width: 1024px) {
        .detail-grid { grid-template-columns: 1fr; }
        .hero-value { font-size: 3rem; }
    }
</style>