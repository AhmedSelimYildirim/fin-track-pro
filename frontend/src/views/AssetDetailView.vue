<template>
    <div class="detail-container" :style="pageTheme">
        <header class="detail-header">
            <div class="left">
                <button class="back-btn" @click="$router.push('/dashboard')">← {{ t('home') }}</button>
                <h1>{{ type }} {{ t('details') }}</h1>
            </div>

            <div class="right-menu">
                <div class="currency-dropdown-wrapper">
                    <button class="dd-btn" @click="showDd = !showDd">
                        {{ totalInSelectedUnit }} {{ displayLabel }} ▼
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
                    <div v-if="type === 'GOLD'" class="input-group">
                        <label>Altın Tipi</label>
                        <select v-model="form.variant" class="modal-input">
                            <option value="GRAM_24">Gram (24 Ayar)</option>
                            <option value="GRAM_22">Gram (22 Ayar)</option>
                            <option value="GRAM_18">Gram (18 Ayar)</option>
                            <option value="GRAM_14">Gram (14 Ayar)</option>
                            <option value="GRAM_8">Gram (8 Ayar)</option>
                            <option value="GRAM_4">Gram (4 Ayar)</option>
                            <option value="CEYREK">Çeyrek Altın</option>
                            <option value="YARIM">Yarım Altın</option>
                            <option value="TAM">Tam Altın</option>
                            <option value="CUMHURIYET">Cumhuriyet Altını</option>
                            <option value="GREMSE">Gremse Altın</option>
                        </select>
                    </div>

                    <div class="input-group">
                        <label>Miktar</label>
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
    import { t } from '../utils/translations';

    const props = defineProps(['type']);
    const assets = ref([]);
    const transactions = ref([]);
    const showModal = ref(false);
    const showDd = ref(false);

    const selectedCurrency = ref(props.type === 'GOLD' ? 'GOLD' : 'TRY');
    const selectedVariant = ref(props.type === 'GOLD' ? 'GRAM_24' : 'STANDARD');
    const displayLabel = ref(props.type === 'GOLD' ? 'Has (24K)' : 'TRY');

    const form = ref({ amount: '', variant: 'STANDARD', date: new Date().toISOString().split('T')[0] });

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
        return v.replace('GRAM_', 'Gr ').replace('_', ' ');
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
        fetchData();
    };

    const openModal = () => {
        form.value.variant = props.type === 'GOLD' ? 'GRAM_24' : 'STANDARD';
        form.value.amount = '';
        showModal.value = true;
    };

    const submitTx = async (action) => {
        if (!form.value.amount) return alert("Miktar giriniz");
        try {
            await api.post('/assets/update', {
                type: props.type,
                action,
                variant: props.type === 'GOLD' ? form.value.variant : 'STANDARD',
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
    .dd-btn { background: rgba(255,255,255,0.1); border: 1px solid rgba(255,255,255,0.2); color: white; padding: 12px 20px; border-radius: 12px; cursor: pointer; font-weight: 800; font-size: 1.1rem; min-width: 220px; text-align: right; }
    .dd-menu { position: absolute; top: 100%; right: 0; background: #1E293B; border: 1px solid #334155; border-radius: 12px; width: 220px; margin-top: 10px; z-index: 100; max-height: 400px; overflow-y: auto; padding: 10px 0; }
    .dd-category { font-size: 0.8rem; color: #94A3B8; padding: 8px 15px; font-weight: bold; text-transform: uppercase; }
    .dd-item { padding: 10px 15px; cursor: pointer; border-bottom: 1px solid rgba(255,255,255,0.05); transition: 0.2s; }
    .dd-item:hover { background: #334155; color: #FFD700; padding-left: 20px; }
    .dd-divider { height: 1px; background: rgba(255,255,255,0.1); margin: 5px 0; }
    .dd-anim-enter-active, .dd-anim-leave-active { transition: all 0.2s ease; }
    .dd-anim-enter-from, .dd-anim-leave-to { opacity: 0; transform: translateY(-10px); }

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

    .input-group { margin-bottom: 20px; }
    .input-group label { display: block; margin-bottom: 8px; color: #94A3B8; font-size: 0.9rem; }
    .modal-input { width: 100%; padding: 14px; background: #0F172A; border: 1px solid #334155; color: white; border-radius: 10px; font-size: 1rem; box-sizing: border-box; }
    .modal-actions { display: flex; gap: 15px; margin-top: 30px; }
    .btn-add { flex: 1; background: #10B981; border: none; padding: 16px; border-radius: 10px; font-weight: 800; cursor: pointer; color: white; }
    .btn-sub { flex: 1; background: #EF4444; border: none; padding: 16px; border-radius: 10px; font-weight: 800; cursor: pointer; color: white; }

    @media (max-width: 768px) { .detail-grid { grid-template-columns: 1fr; } }
</style>