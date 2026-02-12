<template>
    <div class="detail-container" :style="pageTheme">
        <header class="detail-header">
            <div class="left">
                <button class="back-btn" @click="$router.push('/dashboard')">← Geri Dön</button>
                <h1>{{ type }} Portföyü</h1>
            </div>
            <div class="right">
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
                    <div v-if="filteredAssets.length === 0" class="no-data">Henüz varlık eklenmedi.</div>
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
                <h3>İşlem Geçmişi</h3>
                <div class="history-list">
                    <div v-if="filteredTransactions.length === 0" class="no-data">İşlem yok.</div>
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
                            <option value="CUMHURIYET">Cumhuriyet</option>
                            <option value="GREMSE">Gremse</option>
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
    import { ref, computed, onMounted } from 'vue';
    import api from '../services/api';

    const props = defineProps(['type']);
    const assets = ref([]);
    const transactions = ref([]);
    const showModal = ref(false);
    const displayCurrency = ref('TRY'); // Bu detay sayfasında varsayılan görüntüleme birimi
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
    const filteredAssets = computed(() => assets.value.filter(a => a.type === props.type));
    const filteredTransactions = computed(() => transactions.value.filter(tx => tx.asset_type === props.type));

    const totalInSelectedCurrency = computed(() => {
        const total = filteredAssets.value.reduce((sum, a) => sum + a.value_in_base, 0);
        return total.toLocaleString('tr-TR', { minimumFractionDigits: 2 });
    });

    const formatVariant = (v) => {
        if (v === 'STANDARD') return props.type;
        return v.replace('GRAM_', 'Gr ').replace('_', ' ');
    }

    const fetchData = async () => {
        try {
            const [resSum, resTx] = await Promise.all([
                api.get('/assets/summary', { params: { currency: displayCurrency.value } }),
                api.get('/assets/transactions', { params: { currency: displayCurrency.value } })
            ]);
            assets.value = resSum.data.assets || [];
            transactions.value = resTx.data || [];
        } catch (err) { console.error(err); }
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
        } catch (err) { alert(err.message); }
    };

    const downloadSingleReceipt = async (id) => {
        const res = await api.get(`/assets/receipt/${id}`, { responseType: 'blob', params: { currency: displayCurrency.value } });
        const url = window.URL.createObjectURL(new Blob([res.data]));
        const link = document.createElement('a');
        link.href = url; link.setAttribute('download', `Dekont_${id}.pdf`); link.click();
    };

    onMounted(fetchData);
</script>

<style scoped>
    .detail-container { min-height: 100vh; padding: 40px; color: white; font-family: 'Inter', sans-serif; }
    .detail-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 40px; }
    .back-btn { background: rgba(255,255,255,0.1); border: none; color: white; padding: 10px 25px; border-radius: 30px; cursor: pointer; font-weight: bold; }
    .total-box { text-align: right; background: rgba(255,255,255,0.1); padding: 15px 30px; border-radius: 16px; backdrop-filter: blur(10px); }
    .total-box .value { display: block; font-size: 2.2rem; font-weight: 800; color: #FFD700; }

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