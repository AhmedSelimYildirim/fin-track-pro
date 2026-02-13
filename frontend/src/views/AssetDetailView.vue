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
                        <span class="val">{{ totalInSelectedUnit }}</span>
                        <span class="unit">{{ displayLabel }}</span>
                        <span class="arrow">▼</span>
                    </button>

                    <transition name="dd-anim">
                        <div v-if="showDd" class="dd-menu">

                            <div class="dd-item" @click="changeViewUnit('TRY', 'STANDARD', 'TRY')">
                                <span class="sym">₺</span> TRY
                            </div>
                            <div class="dd-item" @click="changeViewUnit('USD', 'STANDARD', 'USD')">
                                <span class="sym">$</span> USD
                            </div>
                            <div class="dd-item" @click="changeViewUnit('EUR', 'STANDARD', 'EUR')">
                                <span class="sym">€</span> EUR
                            </div>
                            <div class="dd-item" @click="changeViewUnit('BTC', 'STANDARD', 'BTC')">
                                <span class="sym">₿</span> BTC
                            </div>
                            <div class="dd-item" @click="changeViewUnit('SILVER', 'STANDARD', 'SILVER')">
                                <span class="sym">Ag</span> SILVER
                            </div>

                            <div class="dd-group">
                                <div class="dd-item group-trigger" @click.stop="toggleGoldMenu">
                                    <span class="sym">🥇</span> GOLD
                                    <span class="arrow-right" :class="{ rotated: goldMenuOpen }">▶</span>
                                </div>

                                <div v-if="goldMenuOpen" class="sub-menu">
                                    <div v-for="gType in goldTypes" :key="gType.key" class="sub-group">
                                        <div class="dd-item sub-trigger" @click.stop="toggleSubGold(gType.key)">
                                            {{ gType.label }}
                                            <span class="arrow-right" :class="{ rotated: openSubGold === gType.key }">▶</span>
                                        </div>

                                        <div v-if="openSubGold === gType.key" class="deep-menu">
                                            <div class="dd-item deep-item"
                                                 v-for="k in karats"
                                                 :key="k"
                                                 @click="changeViewUnit('GOLD', `${gType.key}_${k}`, `${gType.label} (${k}K)`)">
                                                {{ k }}K
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </div>

                        </div>
                    </transition>
                </div>

                <div class="actions">
                    <button class="icon-btn" @click="downloadExcel">📊</button>
                    <button class="icon-btn" @click="downloadPDF">📄</button>
                </div>
            </div>
        </header>

        <div class="detail-grid">
            <section class="holdings-card">
                <div class="section-header">
                    <h3>Varlık Dağılımı</h3>
                    <button class="add-btn" @click="openModal">İşlem Yap</button>
                </div>
                <div class="holding-items">
                    <div v-if="assets.length === 0" class="no-data">Veri yok.</div>
                    <div v-for="asset in assets" :key="asset.variant" class="h-item">
                        <div class="h-name">
                            <span class="variant-tag">{{ formatVariant(asset.variant) }}</span>
                            <span class="amount">{{ asset.amount }} Adet/Gr</span>
                        </div>
                        <div class="h-price">
                            {{ asset.value_in_base.toLocaleString('tr-TR', {maximumFractionDigits: 2}) }} {{ displayLabel }}
                        </div>
                    </div>
                </div>
            </section>

            <section class="history-card">
                <h3>İşlem Geçmişi</h3>
                <div class="history-list">
                    <div v-if="transactions.length === 0" class="no-data">İşlem yok.</div>
                    <div v-for="tx in transactions" :key="tx.id" class="tx-item">
                        <div class="tx-date">{{ new Date(tx.transaction_date).toLocaleDateString() }}</div>
                        <div class="tx-main">
                            <span :class="tx.type" class="type-dot">●</span>
                            <span class="variant">{{ formatVariant(tx.variant) }}</span>
                            <span class="amt">{{ tx.type === 'add' ? '+' : '-' }}{{ tx.amount }}</span>
                        </div>
                        <button class="mini-btn" @click="downloadReceipt(tx.id)">📄</button>
                    </div>
                </div>
            </section>
        </div>

        <div v-if="showModal" class="modal-overlay" @click.self="showModal = false">
            <div class="modal-content">
                <div class="modal-header">
                    <h3>{{ type }} İşlemi</h3>
                    <button class="close-icon" @click="showModal = false">✕</button>
                </div>

                <div class="modal-body">
                    <div v-if="type === 'GOLD'" class="gold-selector-container">
                        <label>Altın Türü & Ayarı</label>
                        <div class="gold-select-group">
                            <select v-model="modalGoldType" class="modal-input half">
                                <option v-for="t in goldTypes" :key="t.key" :value="t.key">{{ t.label }}</option>
                            </select>
                            <select v-model="modalGoldKarat" class="modal-input half">
                                <option v-for="k in karats" :key="k" :value="k">{{ k }}K</option>
                            </select>
                        </div>
                    </div>

                    <div class="input-group">
                        <label>Miktar (Adet / Gram)</label>
                        <input v-model="form.amount" type="number" placeholder="0.00" class="modal-input">
                    </div>

                    <div class="input-group">
                        <label>Tarih</label>
                        <input v-model="form.date" type="date" class="modal-input">
                    </div>

                    <div class="modal-actions">
                        <button class="btn-add" @click="submitTx('add')">EKLE (+)</button>
                        <button class="btn-sub" @click="submitTx('subtract')">ÇIKAR (-)</button>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
    import { ref, computed, onMounted, watch } from 'vue';
    import api from '../services/api';

    const props = defineProps(['type']);
    const assets = ref([]);
    const transactions = ref([]);
    const showModal = ref(false);
    const showDd = ref(false);
    const goldMenuOpen = ref(false);
    const openSubGold = ref(null);

    const selectedCurrency = ref(props.type === 'GOLD' ? 'GOLD' : 'TRY');
    const selectedVariant = ref('STANDARD');
    const displayLabel = ref('TRY');

    const modalGoldType = ref('GRAM');
    const modalGoldKarat = ref(24);
    const form = ref({ amount: '', date: new Date().toISOString().split('T')[0] });

    const goldTypes = [
        { key: 'CUMHURIYET', label: 'Cumhuriyet' },
        { key: 'TAM', label: 'Tam' },
        { key: 'YARIM', label: 'Yarım' },
        { key: 'CEYREK', label: 'Çeyrek' },
        { key: 'GRAM', label: 'Gram' },
        { key: 'GREMSE', label: 'Gremse' }
    ];
    const karats = [24, 22, 18, 14, 8, 4];

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
        return sum.toLocaleString('tr-TR', { minimumFractionDigits: 2, maximumFractionDigits: 3 });
    });

    const formatVariant = (v) => {
        if (v === 'STANDARD') return props.type;
        return v.replace('_', ' ');
    }

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
        goldMenuOpen.value = false;
        openSubGold.value = null;
        fetchData();
    };

    const toggleGoldMenu = () => { goldMenuOpen.value = !goldMenuOpen.value; };
    const toggleSubGold = (key) => { openSubGold.value = openSubGold.value === key ? null : key; };
    const closeDropdown = () => { showDd.value = false; };

    const openModal = () => {
        form.value.amount = '';
        modalGoldType.value = 'GRAM';
        modalGoldKarat.value = 24;
        showModal.value = true;
    };

    const submitTx = async (action) => {
        if (!form.value.amount) return alert("Miktar giriniz");
        let variantToSend = 'STANDARD';
        if (props.type === 'GOLD') {
            variantToSend = `${modalGoldType.value}_${modalGoldKarat.value}`;
        }

        try {
            await api.post('/assets/update', {
                type: props.type,
                action,
                variant: variantToSend,
                amount: parseFloat(form.value.amount),
                transaction_date: new Date(form.value.date).toISOString()
            });
            showModal.value = false;
            await fetchData();
        } catch (err) { alert(err.response?.data?.error || err.message); }
    };

    const downloadReceipt = async (id) => {
        const res = await api.get(`/assets/receipt/${id}`, { responseType: 'blob', params: { currency: 'TRY' } });
        const url = window.URL.createObjectURL(new Blob([res.data]));
        const link = document.createElement('a'); link.href = url; link.setAttribute('download', 'Dekont.pdf'); link.click();
    };

    const downloadExcel = async () => {
        const res = await api.get(`/assets/export/excel`, { responseType: 'blob', params: { currency: 'TRY' } });
        const url = window.URL.createObjectURL(new Blob([res.data]));
        const link = document.createElement('a'); link.href = url; link.setAttribute('download', 'Export.xlsx'); link.click();
    };

    const downloadPDF = async () => {
        const res = await api.get(`/assets/receipt/full`, { responseType: 'blob', params: { currency: 'TRY' } });
        const url = window.URL.createObjectURL(new Blob([res.data]));
        const link = document.createElement('a'); link.href = url; link.setAttribute('download', 'Rapor.pdf'); link.click();
    };

    const vClickOutside = {
        mounted(el, binding) {
            el.clickOutsideEvent = function(event) {
                if (!(el === event.target || el.contains(event.target))) {
                    binding.value(event, el);
                }
            };
            document.body.addEventListener('click', el.clickOutsideEvent);
        },
        unmounted(el) {
            document.body.removeEventListener('click', el.clickOutsideEvent);
        }
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
    .detail-container { min-height: 100vh; padding: 40px; color: white; font-family: 'Inter', sans-serif; }
    .detail-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 40px; }
    .back-btn { background: rgba(255,255,255,0.1); border: none; color: white; padding: 10px 25px; border-radius: 30px; cursor: pointer; font-weight: bold; }

    .right-menu { display: flex; align-items: center; gap: 15px; }
    .currency-dropdown-wrapper { position: relative; }
    .dd-btn { background: rgba(255,255,255,0.1); border: 1px solid rgba(255,255,255,0.2); color: white; padding: 12px 20px; border-radius: 12px; cursor: pointer; font-weight: 800; font-size: 1.1rem; min-width: 220px; display: flex; justify-content: space-between; align-items: center; }
    .dd-btn .val { margin-right: 5px; } .dd-btn .unit { opacity: 0.8; font-size: 0.9em; }

    .dd-menu {
        position: absolute; top: 100%; right: 0;
        background: #1E293B; border: 1px solid #334155;
        border-radius: 12px; width: 220px; margin-top: 10px;
        z-index: 100; max-height: 400px; overflow-y: auto; overflow-x: hidden;
    }
    /* Scrollbar */
    .dd-menu::-webkit-scrollbar { width: 6px; }
    .dd-menu::-webkit-scrollbar-thumb { background: #475569; border-radius: 3px; }

    .dd-item { padding: 12px 15px; cursor: pointer; border-bottom: 1px solid rgba(255,255,255,0.05); transition: 0.2s; display: flex; justify-content: space-between; align-items: center; }
    .dd-item:hover { background: #334155; color: #FFD700; }
    .sym { width: 25px; display: inline-block; font-weight: bold; }

    /* Nested Menu Styles */
    .sub-menu { background: #0F172A; border-left: 2px solid #FFD700; }
    .sub-trigger { padding-left: 25px; font-size: 0.95rem; }
    .deep-menu { background: #020617; }
    .deep-item { padding-left: 40px; font-size: 0.85rem; color: #94A3B8; }
    .deep-item:hover { color: #fff; }
    .arrow-right { font-size: 0.8rem; transition: 0.3s; }
    .arrow-right.rotated { transform: rotate(90deg); }

    .actions { display: flex; gap: 10px; }
    .icon-btn { background: rgba(255,255,255,0.1); border: none; width: 45px; height: 45px; border-radius: 10px; cursor: pointer; font-size: 1.2rem; }

    .detail-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 30px; }
    .holdings-card, .history-card { background: rgba(0,0,0,0.3); backdrop-filter: blur(20px); padding: 30px; border-radius: 24px; border: 1px solid rgba(255,255,255,0.1); min-height: 400px; }
    .section-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
    .add-btn { background: #FFD700; border: none; padding: 10px 20px; border-radius: 8px; font-weight: bold; cursor: pointer; color: #000; }

    .h-item, .tx-item { display: flex; justify-content: space-between; align-items: center; padding: 15px 0; border-bottom: 1px solid rgba(255,255,255,0.05); }
    .variant-tag { background: #FFD700; color: black; padding: 4px 10px; border-radius: 6px; font-size: 0.8rem; font-weight: bold; margin-right: 10px; }
    .type-dot { margin-right: 8px; font-size: 1.2rem; }
    .type-dot.add { color: #10B981; } .type-dot.subtract { color: #EF4444; }

    .modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.85); backdrop-filter: blur(5px); display: flex; justify-content: center; align-items: center; z-index: 2000; }
    .modal-content { background: #1E293B; padding: 30px; border-radius: 24px; width: 450px; box-shadow: 0 20px 60px rgba(0,0,0,0.5); border: 1px solid #334155; }
    .modal-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 25px; }
    .close-icon { background: none; border: none; color: #94A3B8; font-size: 1.5rem; cursor: pointer; transition: 0.2s; }
    .close-icon:hover { color: #EF4444; }

    .gold-select-group { display: flex; gap: 10px; margin-bottom: 20px; }
    .input-group { margin-bottom: 20px; }
    .input-group label, .gold-selector-container label { display: block; margin-bottom: 8px; color: #94A3B8; font-size: 0.9rem; }
    .modal-input { width: 100%; padding: 14px; background: #0F172A; border: 1px solid #334155; color: white; border-radius: 10px; font-size: 1rem; box-sizing: border-box; }
    .modal-input.half { width: 50%; }

    .modal-actions { display: flex; gap: 15px; margin-top: 30px; }
    .btn-add { flex: 1; background: #10B981; border: none; padding: 16px; border-radius: 10px; font-weight: 800; cursor: pointer; color: white; }
    .btn-sub { flex: 1; background: #EF4444; border: none; padding: 16px; border-radius: 10px; font-weight: 800; cursor: pointer; color: white; }

    @media (max-width: 768px) { .detail-grid { grid-template-columns: 1fr; } }
</style>