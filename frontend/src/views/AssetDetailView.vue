<template>
    <div class="detail-container" :style="pageTheme">
        <header class="detail-header">
            <div class="left">
                <button class="back-btn" @click="$router.push('/dashboard')">← {{ t('home') }}</button>
                <h1>{{ type }} {{ t('details') }}</h1>
            </div>
            <div class="right">
                <div class="action-group">
                    <button class="icon-btn excel" @click="downloadExcel" title="Excel">📊</button>
                    <button class="icon-btn pdf" @click="downloadFullPDF" title="PDF">📄</button>
                </div>
                <div class="total-box">
                    <span class="label">Toplam Değer</span>
                    <span class="value">{{ totalInSelectedCurrency }} {{ displayCurrency }}</span>
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
                    <div v-if="filteredAssets.length === 0" class="no-data">Veri bulunamadı.</div>
                    <div v-for="asset in filteredAssets" :key="asset.variant" class="h-item">
                        <div class="h-name">
                            <span class="variant-tag">{{ formatVariant(asset.variant) }}</span>
                            <span class="amount">{{ asset.amount }} Adet/Gr</span>
                        </div>
                        <div class="h-price">
                            <span class="val">{{ asset.value_in_base.toLocaleString('tr-TR', {maximumFractionDigits: 2}) }} {{ displayCurrency }}</span>
                        </div>
                    </div>
                </div>
            </section>

            <section class="history-card">
                <h3>Son İşlemler</h3>
                <div class="history-list">
                    <div v-if="filteredTransactions.length === 0" class="no-data">İşlem geçmişi boş.</div>
                    <div v-for="tx in filteredTransactions" :key="tx.id" class="tx-item">
                        <div class="tx-date">{{ new Date(tx.transaction_date).toLocaleDateString('tr-TR') }}</div>
                        <div class="tx-main">
                            <span :class="tx.type" class="type-dot">●</span>
                            <span class="variant">{{ formatVariant(tx.variant) }}</span>
                            <span class="amt">{{ tx.type === 'add' ? '+' : '-' }}{{ tx.amount }}</span>
                        </div>
                        <button class="mini-btn" @click="downloadSingleReceipt(tx.id)">📄</button>
                    </div>
                </div>
            </section>
        </div>

        <div v-if="showModal" class="modal-overlay" @click.self="showModal = false">
            <div class="modal-content">
                <h3>{{ type }} İşlemi</h3>
                <div v-if="type === 'GOLD'" class="input-group">
                    <label>Birim Seçin</label>
                    <select v-model="form.variant" class="modal-input">
                        <option value="GRAM_24">Gram 24K</option>
                        <option value="CEYREK">Çeyrek</option>
                        <option value="TAM">Tam</option>
                        <option value="CUMHURIYET">Cumhuriyet</option>
                    </select>
                </div>
                <div class="input-group">
                    <label>Miktar (Adet/Gr)</label>
                    <input v-model="form.amount" type="number" placeholder="0.00" class="modal-input">
                </div>
                <div class="modal-actions">
                    <button class="btn-add" @click="submitTx('add')">Ekle (+)</button>
                    <button class="btn-sub" @click="submitTx('subtract')">Çıkar (-)</button>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
    import { ref, computed, onMounted } from 'vue';
    import api from '../services/api';
    import { t } from '../utils/translations';

    const props = defineProps(['type']);
    const assets = ref([]);
    const transactions = ref([]);
    const showModal = ref(false);
    const displayCurrency = ref(localStorage.getItem('preferredCurrency') || 'TRY');
    const form = ref({ amount: '', variant: 'STANDARD' });

    const themeMap = {
        TRY: 'linear-gradient(135deg, #1e0505 0%, #450a0a 100%)',
        USD: 'linear-gradient(135deg, #051e0f 0%, #064e3b 100%)',
        GOLD: 'linear-gradient(135deg, #1e1b05 0%, #451a03 100%)',
        EUR: 'linear-gradient(135deg, #050b1e 0%, #1e3a8a 100%)',
        BTC: 'linear-gradient(135deg, #1e1305 0%, #78350f 100%)',
        SILVER: 'linear-gradient(135deg, #111827 0%, #374151 100%)'
    };

    const pageTheme = computed(() => ({ background: themeMap[props.type] || '#000' }));
    const filteredAssets = computed(() => assets.value.filter(a => a.type === props.type));
    const filteredTransactions = computed(() => transactions.value.filter(tx => tx.asset_type === props.type));

    const totalInSelectedCurrency = computed(() => {
        const total = filteredAssets.value.reduce((sum, a) => sum + a.value_in_base, 0);
        return total.toLocaleString('tr-TR', { minimumFractionDigits: 2 });
    });

    const formatVariant = (v) => v === 'STANDARD' ? props.type : v.replace('_', ' ');

    const fetchData = async () => {
        try {
            const [resSum, resTx] = await Promise.all([
                api.get('/assets/summary', { params: { currency: displayCurrency.value } }),
                api.get('/assets/transactions', { params: { currency: displayCurrency.value } })
            ]);
            assets.value = resSum.data.assets || [];
            transactions.value = resTx.data || [];
        } catch (err) {
            console.error("Veri çekme hatası:", err);
        }
    };

    const openModal = () => {
        form.value.variant = props.type === 'GOLD' ? 'GRAM_24' : 'STANDARD';
        form.value.amount = '';
        showModal.value = true;
    };

    const submitTx = async (action) => {
        if (!form.value.amount || form.value.amount <= 0) return alert("Geçerli bir miktar girin.");
        try {
            await api.post('/assets/update', {
                ...form.value,
                type: props.type,
                action,
                amount: parseFloat(form.value.amount),
                transaction_date: new Date().toISOString()
            });
            showModal.value = false;
            await fetchData();
        } catch (err) {
            alert("İşlem başarısız: " + (err.response?.data?.error || "Sunucu hatası"));
        }
    };

    const downloadSingleReceipt = async (id) => {
        const res = await api.get(`/assets/receipt/${id}`, { responseType: 'blob', params: { currency: displayCurrency.value } });
        const url = window.URL.createObjectURL(new Blob([res.data]));
        const link = document.createElement('a');
        link.href = url; link.setAttribute('download', `Dekont_${id}.pdf`); link.click();
    };

    const downloadFullPDF = async () => {
        const res = await api.get('/assets/receipt/full', { responseType: 'blob', params: { currency: displayCurrency.value } });
        const url = window.URL.createObjectURL(new Blob([res.data]));
        const link = document.createElement('a');
        link.href = url; link.setAttribute('download', `${props.type}_Raporu.pdf`); link.click();
    };

    const downloadExcel = async () => {
        const res = await api.get('/assets/export/excel', { responseType: 'blob', params: { currency: displayCurrency.value } });
        const url = window.URL.createObjectURL(new Blob([res.data]));
        const link = document.createElement('a');
        link.href = url; link.setAttribute('download', `${props.type}_Export.xlsx`); link.click();
    };

    onMounted(fetchData);
</script>

<style scoped>
    .detail-container { min-height: 100vh; padding: 40px; color: white; font-family: 'Inter', sans-serif; }
    .detail-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 40px; }
    .back-btn { background: rgba(255,255,255,0.1); border: none; color: white; padding: 10px 25px; border-radius: 30px; cursor: pointer; transition: 0.3s; }
    .back-btn:hover { background: rgba(255,255,255,0.2); }
    .action-group { display: flex; gap: 10px; margin-bottom: 10px; justify-content: flex-end; }
    .icon-btn { background: rgba(255,255,255,0.1); border: none; border-radius: 10px; width: 40px; height: 40px; cursor: pointer; color: white; transition: 0.2s; }
    .icon-btn:hover { transform: scale(1.1); background: rgba(255,255,255,0.2); }
    .total-box { text-align: right; background: rgba(255,255,255,0.1); padding: 15px 25px; border-radius: 15px; backdrop-filter: blur(10px); }
    .total-box .value { display: block; font-size: 2.2rem; font-weight: 800; color: #FFD700; }
    .detail-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 30px; }
    .holdings-card, .history-card { background: rgba(0,0,0,0.3); backdrop-filter: blur(100px); padding: 30px; border-radius: 24px; border: 1px solid rgba(255,255,255,0.1); box-shadow: 0 20px 50px rgba(0,0,0,0.5); }
    .h-item, .tx-item { display: flex; justify-content: space-between; align-items: center; padding: 18px 0; border-bottom: 1px solid rgba(255,255,255,0.05); }
    .variant-tag { background: #FFD700; color: black; padding: 5px 12px; border-radius: 8px; font-size: 0.75rem; font-weight: 900; text-transform: uppercase; margin-right: 15px; }
    .add-btn { background: #FFD700; border: none; padding: 12px 24px; border-radius: 12px; font-weight: 800; cursor: pointer; color: black; transition: 0.3s; }
    .add-btn:hover { transform: translateY(-3px); box-shadow: 0 5px 15px rgba(255,215,0,0.4); }
    .type-dot { margin-right: 8px; font-size: 0.8rem; }
    .type-dot.add { color: #10B981; }
    .type-dot.subtract { color: #EF4444; }
    .no-data { color: var(--text-muted); text-align: center; padding: 40px; opacity: 0.5; }
    .modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.85); backdrop-filter: blur(8px); display: flex; justify-content: center; align-items: center; z-index: 2000; }
    .modal-content { background: #1e293b; padding: 40px; border-radius: 24px; width: 420px; border: 1px solid #334155; }
    .input-group { margin-bottom: 20px; }
    .input-group label { display: block; font-size: 0.85rem; color: #94A3B8; margin-bottom: 8px; }
    .modal-input { width: 100%; padding: 15px; background: #0f172a; border: 1px solid #334155; color: white; border-radius: 12px; font-size: 1rem; box-sizing: border-box; }
    .modal-actions { display: flex; gap: 15px; margin-top: 30px; }
    .btn-add { flex: 1; background: #10B981; border: none; padding: 16px; border-radius: 12px; color: white; font-weight: 800; cursor: pointer; transition: 0.2s; }
    .btn-sub { flex: 1; background: #EF4444; border: none; padding: 16px; border-radius: 12px; color: white; font-weight: 800; cursor: pointer; transition: 0.2s; }
    .btn-add:hover, .btn-sub:hover { opacity: 0.9; transform: scale(1.02); }
</style>