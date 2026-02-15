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
                            <div class="dd-item" @click="changeViewUnit('TRY', 'STANDARD', 'TRY')">₺ TRY</div>
                            <div class="dd-item" @click="changeViewUnit('USD', 'STANDARD', 'USD')">$ USD</div>
                            <div class="dd-item" @click="changeViewUnit('EUR', 'STANDARD', 'EUR')">€ EUR</div>
                            <div class="dd-item" @click="changeViewUnit('BTC', 'STANDARD', 'BTC')">₿ BTC</div>
                            <div class="dd-item" @click="changeViewUnit('SILVER', 'STANDARD', 'SILVER')">Ag SILVER</div>
                            <div class="dd-divider"></div>
                            <div class="dd-section-title">Altın Birimleri</div>
                            <div class="dd-group">
                                <div class="dd-item sub-trigger" @click.stop="toggleGramMenu">⚖️ Gram Altın <span class="arrow-right" :class="{ rotated: gramMenuOpen }">▶</span></div>
                                <div v-if="gramMenuOpen" class="deep-menu">
                                    <div class="dd-item deep-item" v-for="k in karats" :key="k" @click="changeViewUnit('GOLD', `GRAM_${k}`, `Gram (${k}K)`)">{{ k }} Ayar</div>
                                </div>
                            </div>
                            <div class="dd-item" @click="changeViewUnit('GOLD', 'CUMHURIYET', 'Cumhuriyet')">🔴 Cumhuriyet</div>
                            <div class="dd-item" @click="changeViewUnit('GOLD', 'TAM', 'Tam')">🌕 Tam</div>
                            <div class="dd-item" @click="changeViewUnit('GOLD', 'YARIM', 'Yarım')">🌓 Yarım</div>
                            <div class="dd-item" @click="changeViewUnit('GOLD', 'CEYREK', 'Çeyrek')">🔸 Çeyrek</div>
                            <div class="dd-item" @click="changeViewUnit('GOLD', 'GREMSE', 'Gremse')">🔶 Gremse</div>
                        </div>
                    </transition>
                </div>

                <div class="actions">
                    <button class="icon-btn" @click="downloadExcel" :title="type + ' Excel'">📊</button>
                    <button class="icon-btn" @click="downloadFullPDF" :title="type + ' Rapor'">📄</button>
                </div>
            </div>
        </header>

        <div class="hero-container">
            <div class="hero-card">
                <div class="hero-label">TOPLAM {{ type }} DEĞERİ</div>
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
                                <option value="GREMSE">Gremse</option>
                            </select>
                        </div>
                        <div v-if="form.goldType === 'GRAM'" class="input-group">
                            <label>Ayar</label>
                            <select v-model="form.goldKarat" class="custom-input">
                                <option v-for="k in karats" :key="k" :value="k">{{ k }}K</option>
                            </select>
                        </div>
                    </div>
                    <div class="input-row">
                        <div class="input-group"><label>Miktar</label><input v-model="form.amount" type="number" class="custom-input"></div>
                        <div class="input-group"><label>Tarih</label><input v-model="form.date" type="date" class="custom-input"></div>
                    </div>
                    <div class="action-buttons">
                        <button class="btn-add" @click="submitTx('add')">+ EKLE</button>
                        <button class="btn-sub" @click="submitTx('subtract')">- ÇIKAR</button>
                    </div>
                </div>
            </section>

            <section class="history-card">
                <div class="card-header"><h3>İşlem Geçmişi</h3></div>
                <div class="history-list">
                    <div v-for="tx in transactions" :key="tx.id" class="tx-item">
                        <div class="tx-left">
                            <div class="tx-date">{{ new Date(tx.transaction_date).toLocaleDateString('tr-TR') }}</div>
                            <div class="tx-desc">● {{ formatTxName(tx) }}</div>
                        </div>
                        <div class="tx-right">
                            <span class="tx-amount" :class="tx.type">{{ tx.type === 'add' ? '+' : '-' }}{{ tx.amount }}</span>
                            <button class="receipt-btn" @click="downloadSingleReceipt(tx.id)">📄</button>
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

    const themeMap = { TRY: 'linear-gradient(135deg, #1e0505 0%, #450a0a 100%)', USD: 'linear-gradient(135deg, #051e0f 0%, #064e3b 100%)', GOLD: 'linear-gradient(135deg, #2A1C05 0%, #713F12 100%)', EUR: 'linear-gradient(135deg, #3E2723 0%, #5D4037 100%)', BTC: 'linear-gradient(135deg, #1A1A1A 0%, #000000 100%)', SILVER: 'linear-gradient(135deg, #111827 0%, #374151 100%)' };
    const pageTheme = computed(() => ({ background: themeMap[props.type] || '#000' }));

    const totalInSelectedUnit = computed(() => assets.value.reduce((acc, curr) => acc + curr.value_in_base, 0).toLocaleString('tr-TR', { minimumFractionDigits: 2 }));

    const fetchData = async () => {
        try {
            const [resSum, resTx] = await Promise.all([
                api.get('/assets/summary', { params: { currency: selectedCurrency.value, variant: selectedVariant.value } }),
                api.get('/assets/transactions', { params: { currency: selectedCurrency.value } })
            ]);
            assets.value = (resSum.data.assets || []).filter(a => a.type === props.type);
            transactions.value = (resTx.data || []).filter(t => t.asset_type === props.type).sort((a, b) => {
                const d = new Date(b.transaction_date) - new Date(a.transaction_date);
                return d !== 0 ? d : b.id - a.id; // SIRALAMA FIX
            });
        } catch (err) { console.error(err); }
    };

    const changeViewUnit = (c, v, l) => { selectedCurrency.value = c; selectedVariant.value = v; displayLabel.value = l; showDd.value = false; gramMenuOpen.value = false; fetchData(); };
    const toggleGramMenu = () => { gramMenuOpen.value = !gramMenuOpen.value; };
    const closeDropdown = () => { showDd.value = false; };
    const formatTxName = (tx) => tx.variant.startsWith('GRAM') ? `${tx.variant.split('_')[1]} Ayar Gram` : (tx.variant === 'STANDARD' ? props.type : tx.variant);

    const submitTx = async (action) => {
        let v = props.type === 'GOLD' ? (form.value.goldType === 'GRAM' ? `GRAM_${form.value.goldKarat}` : form.value.goldType) : 'STANDARD';
        try {
            await api.post('/assets/update', { type: props.type, action, variant: v, amount: parseFloat(form.value.amount), transaction_date: new Date(form.value.date).toISOString() });
            form.value.amount = ''; fetchData();
        } catch (err) { alert(err.response?.data?.error || err.message); }
    };

    const downloadSingleReceipt = async (id) => {
        const res = await api.get(`/assets/receipt/${id}`, { responseType: 'blob', params: { currency: selectedCurrency.value } });
        const url = window.URL.createObjectURL(new Blob([res.data]));
        const link = document.createElement('a'); link.href = url; link.setAttribute('download', `Dekont_${id}.pdf`); link.click();
    };

    const downloadFullPDF = async () => {
        const res = await api.get(`/assets/receipt/full`, {
            responseType: 'blob',
            params: {
                currency: selectedCurrency.value,
                variant: selectedVariant.value,
                asset_type: props.type // FILTRE EKLEDIK
            }
        });
        const url = window.URL.createObjectURL(new Blob([res.data]));
        const link = document.createElement('a'); link.href = url; link.setAttribute('download', `Rapor_${props.type}.pdf`); link.click();
    };

    const downloadExcel = async () => {
        const res = await api.get(`/assets/export/excel`, {
            responseType: 'blob',
            params: {
                currency: selectedCurrency.value,
                variant: selectedVariant.value,
                asset_type: props.type // FILTRE EKLEDIK
            }
        });
        const url = window.URL.createObjectURL(new Blob([res.data]));
        const link = document.createElement('a'); link.href = url; link.setAttribute('download', `Rapor_${props.type}.xlsx`); link.click();
    };

    const vClickOutside = { mounted(el, b) { el.clickEvent = (e) => { if (!(el === event.target || el.contains(event.target))) b.value(e); }; document.body.addEventListener('click', el.clickEvent); }, unmounted(el) { document.body.removeEventListener('click', el.clickEvent); } };

    watch(() => props.type, () => { selectedCurrency.value = 'TRY'; selectedVariant.value = 'STANDARD'; displayLabel.value = 'TRY'; fetchData(); });
    onMounted(fetchData);
</script>

<style scoped>
    .detail-container { min-height: 100vh; padding: 40px; color: white; font-family: 'Inter', sans-serif; }
    .detail-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 30px; }
    .right-menu { display: flex; align-items: center; gap: 15px; }
    .dd-btn { background: rgba(255,255,255,0.1); border: 1px solid rgba(255,255,255,0.2); color: white; padding: 10px 20px; border-radius: 12px; cursor: pointer; min-width: 160px; display: flex; justify-content: space-between; }
    .dd-menu { position: absolute; top: 110%; right: 0; background: #1E293B; border-radius: 12px; width: 220px; z-index: 100; box-shadow: 0 10px 40px rgba(0,0,0,0.5); }
    .dd-item { padding: 12px 15px; cursor: pointer; border-bottom: 1px solid rgba(255,255,255,0.05); }
    .dd-item:hover { background: #334155; color: #FFD700; }
    .sub-menu { background: #0F172A; border-left: 2px solid #FFD700; }
    .icon-btn { background: rgba(255,255,255,0.1); border: none; width: 45px; height: 45px; border-radius: 10px; cursor: pointer; font-size: 1.2rem; color: white; }
    .hero-container { display: flex; justify-content: center; margin-bottom: 50px; }
    .hero-card { background: rgba(255, 255, 255, 0.05); backdrop-filter: blur(15px); border: 1px solid rgba(255, 255, 255, 0.1); padding: 30px 60px; border-radius: 24px; text-align: center; }
    .hero-amount { font-size: 3.5rem; font-weight: 800; }
    .hero-curr { font-size: 1.5rem; color: #FFD700; }
    .hero-line { width: 100px; height: 4px; background: #FFD700; margin: 20px auto 0; border-radius: 2px; }
    .detail-grid { display: grid; grid-template-columns: 1fr 1.5fr; gap: 40px; max-width: 1200px; margin: 0 auto; }
    .action-card, .history-card { background: rgba(0,0,0,0.3); backdrop-filter: blur(20px); border-radius: 24px; border: 1px solid rgba(255,255,255,0.1); padding: 20px; }
    .custom-input { width: 100%; padding: 12px; background: rgba(0,0,0,0.4); border: 1px solid rgba(255,255,255,0.15); color: white; border-radius: 10px; margin-bottom: 15px; }
    .btn-add, .btn-sub { flex: 1; padding: 15px; border-radius: 12px; font-weight: 800; cursor: pointer; border: none; color: white; margin: 5px; }
    .btn-add { background: #10B981; } .btn-sub { background: #EF4444; }
    .tx-item { display: flex; justify-content: space-between; align-items: center; padding: 15px 0; border-bottom: 1px solid rgba(255,255,255,0.05); }
    .tx-amount.add { color: #10B981; } .tx-amount.subtract { color: #EF4444; }
    .receipt-btn { background: none; border: 1px solid rgba(255,255,255,0.1); border-radius: 8px; color: white; cursor: pointer; padding: 5px 10px; }
</style>