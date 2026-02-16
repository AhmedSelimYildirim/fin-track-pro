<template>
    <div class="detail-container" :style="pageTheme">
        <header class="detail-header">
            <div class="left">
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
                            <div class="dd-section-title">Para Birimleri</div>
                            <div class="dd-item" @click="changeViewUnit('TRY', 'STANDARD', 'TRY')"><span class="sym">₺</span> TRY</div>
                            <div class="dd-item" @click="changeViewUnit('USD', 'STANDARD', 'USD')"><span class="sym">$</span> USD</div>
                            <div class="dd-item" @click="changeViewUnit('EUR', 'STANDARD', 'EUR')"><span class="sym">€</span> EUR</div>
                            <div class="dd-item" @click="changeViewUnit('BTC', 'STANDARD', 'BTC')"><span class="sym">₿</span> BTC</div>
                            <div class="dd-item" @click="changeViewUnit('SILVER', 'STANDARD', 'SILVER')"><span class="sym">Ag</span> SILVER</div>
                            <div class="dd-divider"></div>
                            <div class="dd-section-title">Altın Birimleri</div>
                            <div class="dd-group">
                                <div class="dd-item sub-trigger" @click.stop="toggleGramMenu">
                                    <span class="sym">⚖️</span> Gram Altın
                                    <span class="arrow-right" :class="{ rotated: gramMenuOpen }">▶</span>
                                </div>
                                <div v-if="gramMenuOpen" class="deep-menu">
                                    <div class="dd-item deep-item" v-for="k in karats" :key="k" @click="changeViewUnit('GOLD', `GRAM_${k}`, `Gram (${k}K)`)">{{ k }} Ayar</div>
                                </div>
                            </div>
                            <div class="dd-item" @click="changeViewUnit('GOLD', 'CUMHURIYET', 'Cumhuriyet')"><span class="sym">🔴</span> Cumhuriyet</div>
                            <div class="dd-item" @click="changeViewUnit('GOLD', 'TAM', 'Tam')"><span class="sym">🌕</span> Tam</div>
                            <div class="dd-item" @click="changeViewUnit('GOLD', 'YARIM', 'Yarım')"><span class="sym">🌓</span> Yarım</div>
                            <div class="dd-item" @click="changeViewUnit('GOLD', 'CEYREK', 'Çeyrek')"><span class="sym">🔸</span> Çeyrek</div>
                            <div class="dd-item" @click="changeViewUnit('GOLD', 'GREMSE', 'Gremse')"><span class="sym">🔶</span> Gremse</div>
                        </div>
                    </transition>
                </div>

                <div class="actions">
                    <button class="icon-btn" @click="downloadExcel" :title="type + ' Excel İndir'">📊</button>
                    <button class="icon-btn" @click="downloadFullPDF" :title="type + ' Rapor İndir'">📄</button>
                </div>
            </div>
        </header>

        <div class="hero-container">
            <div class="hero-card">
                <div class="hero-label">TOPLAM VARLIK DEĞERİ</div>
                <div class="hero-amount">{{ totalInSelectedUnit }} <span class="hero-curr">{{ displayLabel }}</span></div>
                <div class="hero-line"></div>
            </div>
        </div>

        <div class="detail-grid">
            <section class="action-card">
                <div class="card-header"><h3>Hızlı İşlem</h3></div>
                <div class="form-body">
                    <div v-if="type === 'GOLD'" class="gold-options">
                        <div class="input-group">
                            <label>Altın Türü</label>
                            <select v-model="form.goldType" class="custom-input">
                                <option value="GRAM">Gram Altın</option>
                                <option value="CUMHURIYET">Cumhuriyet</option>
                                <option value="TAM">Tam Altın</option>
                                <option value="YARIM">Yarım Altın</option>
                                <option value="CEYREK">Çeyrek Altın</option>
                                <option value="GREMSE">Gremse Altın</option>
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
                        <button class="btn-add" @click="submitTx('add')"><span>+</span> EKLE</button>
                        <button class="btn-sub" @click="submitTx('subtract')"><span>-</span> ÇIKAR</button>
                    </div>
                </div>
            </section>

            <section class="history-card">
                <div class="card-header"><h3>İşlem Geçmişi</h3></div>
                <div class="history-list">
                    <div v-if="transactions.length === 0" class="no-data">İşlem geçmişi boş.</div>
                    <div v-for="tx in transactions" :key="tx.id" class="tx-item">
                        <div class="tx-left">
                            <div class="tx-date">{{ new Date(tx.transaction_date).toLocaleDateString('tr-TR') }}</div>
                            <div class="tx-desc">
                                <span class="dot" :class="tx.type">●</span>
                                <span class="tx-name">{{ formatTxName(tx) }}</span>
                            </div>
                        </div>
                        <div class="tx-right">
                            <span class="tx-amount" :class="tx.type">{{ tx.type === 'add' ? '+' : '-' }}{{ tx.amount }}</span>
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
    const gramMenuOpen = ref(false);

    const selectedCurrency = ref('TRY');
    const selectedVariant = ref('STANDARD');
    const displayLabel = ref('TRY');

    const form = ref({ amount: '', date: new Date().toISOString().split('T')[0], goldType: 'GRAM', goldKarat: 24 });
    const karats = [24, 22, 18, 14, 8, 4];

    const themeMap = {
        TRY: 'linear-gradient(135deg, #1e0505 0%, #450a0a 100%)',
        USD: 'linear-gradient(135deg, #051e0f 0%, #064e3b 100%)',
        GOLD: 'linear-gradient(135deg, #2A1C05 0%, #713F12 100%)',
        EUR: 'linear-gradient(135deg, #3E2723 0%, #5D4037 100%)',
        BTC: 'linear-gradient(135deg, #1A1A1A 0%, #000000 100%)',
        SILVER: 'linear-gradient(135deg, #111827 0%, #374151 100%)'
    };

    const pageTheme = computed(() => ({ background: themeMap[props.type] || '#0f172a' }));

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
            transactions.value = (resTx.data || [])
                .filter(t => t.asset_type === props.type)
                .sort((a, b) => {
                    const dateA = new Date(a.transaction_date);
                    const dateB = new Date(b.transaction_date);
                    if (dateA.getTime() === dateB.getTime()) { return b.id - a.id; }
                    return dateB - dateA;
                });
        } catch (err) { console.error(err); }
    };

    const changeViewUnit = (code, variant, label) => {
        selectedCurrency.value = code;
        selectedVariant.value = variant;
        displayLabel.value = label;
        showDd.value = false;
        gramMenuOpen.value = false;
        fetchData();
    };

    const toggleGramMenu = () => { gramMenuOpen.value = !gramMenuOpen.value; };
    const closeDropdown = () => { showDd.value = false; };

    const formatTxName = (tx) => {
        if (tx.asset_type !== 'GOLD') return props.type;
        if (tx.variant.startsWith('GRAM')) { return `${tx.variant.split('_')[1]} Ayar Gram`; }
        const map = { 'CEYREK': 'Çeyrek', 'YARIM': 'Yarım', 'TAM': 'Tam', 'CUMHURIYET': 'Cumhuriyet', 'GREMSE': 'Gremse' };
        return map[tx.variant] || tx.variant;
    };

    const submitTx = async (action) => {
        if (!form.value.amount || form.value.amount <= 0) return alert("Geçerli miktar giriniz");
        let variantToSend = 'STANDARD';
        if (props.type === 'GOLD') {
            variantToSend = form.value.goldType === 'GRAM' ? `GRAM_${form.value.goldKarat}` : form.value.goldType;
        }
        try {
            await api.post('/assets/update', {
                type: props.type, action, variant: variantToSend,
                amount: parseFloat(form.value.amount), transaction_date: new Date(form.value.date).toISOString()
            });
            form.value.amount = '';
            await fetchData();
        } catch (err) { alert(err.response?.data?.error || err.message); }
    };

    const downloadSingleReceipt = async (id) => {
        const res = await api.get(`/assets/receipt/${id}`, { responseType: 'blob', params: { currency: selectedCurrency.value } });
        const url = window.URL.createObjectURL(new Blob([res.data]));
        const link = document.createElement('a'); link.href = url; link.setAttribute('download', `Dekont_${id}.pdf`); link.click();
    };

    const downloadFullPDF = async () => {
        const res = await api.get(`/assets/receipt/full`, { responseType: 'blob', params: { currency: selectedCurrency.value, variant: selectedVariant.value, asset_type: props.type } });
        const url = window.URL.createObjectURL(new Blob([res.data]));
        const link = document.createElement('a'); link.href = url; link.setAttribute('download', `Rapor_${props.type}.pdf`); link.click();
    };

    const downloadExcel = async () => {
        const res = await api.get(`/assets/export/excel`, { responseType: 'blob', params: { currency: selectedCurrency.value, variant: selectedVariant.value, asset_type: props.type } });
        const url = window.URL.createObjectURL(new Blob([res.data]));
        const link = document.createElement('a'); link.href = url; link.setAttribute('download', `Rapor_${props.type}.xlsx`); link.click();
    };

    const vClickOutside = {
        mounted(el, binding) { el.clickOutsideEvent = function(event) { if (!(el === event.target || el.contains(event.target))) { binding.value(event, el); } }; document.body.addEventListener('click', el.clickOutsideEvent); },
        unmounted(el) { document.body.removeEventListener('click', el.clickOutsideEvent); }
    };

    watch(() => props.type, () => { selectedCurrency.value = 'TRY'; selectedVariant.value = 'STANDARD'; displayLabel.value = 'TRY'; fetchData(); });
    onMounted(fetchData);
</script>

<style scoped>
    * { box-sizing: border-box; }
    .detail-container { min-height: 100vh; padding: 40px; color: white; font-family: 'Inter', sans-serif; transition: background 0.3s ease; }
    .detail-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 40px; position: relative; z-index: 50; }
    .left h1 { margin: 0; font-size: 2rem; font-weight: 700; letter-spacing: -0.5px; }
    .right-menu { display: flex; align-items: center; gap: 20px; }
    .dd-btn { background: #1e293b; color: white; border: 1px solid #334155; padding: 12px 24px; border-radius: 12px; cursor: pointer; font-weight: 700; min-width: 180px; display: flex; justify-content: space-between; align-items: center; box-shadow: 0 4px 10px rgba(0,0,0,0.2); }
    .dd-btn:hover { background: #334155; border-color: #fbbf24; }
    .dd-menu { position: absolute; top: calc(100% + 10px); right: 0; background: #1E293B; border: 1px solid #334155; border-radius: 12px; width: 240px; z-index: 100; max-height: 400px; overflow-y: auto; box-shadow: 0 10px 40px rgba(0,0,0,0.5); }
    .dd-section-title { font-size: 0.75rem; color: #94A3B8; padding: 12px 16px 4px; font-weight: bold; letter-spacing: 1px; }
    .dd-divider { height: 1px; background: rgba(255,255,255,0.1); margin: 6px 0; }
    .dd-item { padding: 12px 16px; cursor: pointer; border-bottom: 1px solid rgba(255,255,255,0.05); color: white; display: flex; align-items: center; gap: 10px; }
    .dd-item:hover { background: #334155; color: #fbbf24; padding-left: 20px; }
    .deep-menu { background: #020617; }
    .deep-item { padding-left: 35px; font-size: 0.9rem; color: #94A3B8; }
    .arrow-right { font-size: 0.8rem; transition: 0.3s; margin-left: auto; }
    .arrow-right.rotated { transform: rotate(90deg); }
    .actions { display: flex; gap: 15px; }
    .icon-btn { background: #1e293b; color: white; border: 1px solid #334155; width: 50px; height: 50px; border-radius: 12px; cursor: pointer; font-size: 1.4rem; display: flex; align-items: center; justify-content: center; transition: 0.3s; }
    .icon-btn:hover { background: #334155; transform: translateY(-2px); }
    .hero-container { display: flex; justify-content: center; margin-bottom: 60px; }
    .hero-card { background: rgba(255, 255, 255, 0.05); backdrop-filter: blur(20px); border: 1px solid rgba(255, 255, 255, 0.1); padding: 40px 80px; border-radius: 30px; text-align: center; box-shadow: 0 20px 60px rgba(0,0,0,0.4); min-width: 400px; }
    .hero-amount { font-size: 4rem; font-weight: 800; letter-spacing: -2px; }
    .hero-curr { font-size: 1.8rem; color: #fbbf24; margin-left: 12px; }
    .hero-line { width: 120px; height: 5px; background: #fbbf24; margin: 25px auto 0; border-radius: 3px; }
    .detail-grid { display: grid; grid-template-columns: 1fr 1.5fr; gap: 40px; max-width: 1400px; margin: 0 auto; }
    .action-card, .history-card { background: rgba(0,0,0,0.3); backdrop-filter: blur(20px); border-radius: 24px; border: 1px solid rgba(255,255,255,0.1); overflow: hidden; }
    .card-header { padding: 25px 30px; border-bottom: 1px solid rgba(255,255,255,0.1); }
    .form-body { padding: 30px; display: flex; flex-direction: column; gap: 25px; }
    .custom-input { width: 100%; padding: 18px; background: rgba(0,0,0,0.4); border: 1px solid rgba(255,255,255,0.15); color: white; border-radius: 14px; font-size: 1.1rem; }
    .action-buttons { display: flex; gap: 20px; }
    .btn-add, .btn-sub { flex: 1; padding: 20px; border-radius: 14px; font-weight: 800; cursor: pointer; border: none; color: white; transition: 0.3s; }
    .btn-add { background: linear-gradient(135deg, #10B981, #059669); }
    .btn-sub { background: linear-gradient(135deg, #EF4444, #B91C1C); }
    .history-list { padding: 0 30px; max-height: 600px; overflow-y: auto; }
    .tx-item { display: flex; justify-content: space-between; align-items: center; padding: 25px 0; border-bottom: 1px solid rgba(255,255,255,0.05); }
    .tx-amount { font-weight: 800; font-size: 1.2rem; }
    .tx-amount.add { color: #10B981; }
    .tx-amount.subtract { color: #EF4444; }
    .receipt-btn { background: rgba(255,255,255,0.05); border: 1px solid rgba(255,255,255,0.1); width: 45px; height: 45px; border-radius: 12px; cursor: pointer; color: white; }
</style>