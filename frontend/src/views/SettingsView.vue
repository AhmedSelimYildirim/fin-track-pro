<template>
    <div class="settings-content">
        <div class="settings-header">
            <h2>Ayarlar</h2>
        </div>

        <div class="settings-grid">
            <div class="section-card">
                <h3>Görünüm & Dil</h3>

                <div class="setting-row" @click="toggleTheme">
                    <div class="label">
                        <span class="icon">🌓</span>
                        <span>Tema</span>
                    </div>
                    <div class="value">{{ theme === 'dark' ? 'Karanlık' : 'Aydınlık' }}</div>
                </div>

                <div class="setting-row">
                    <div class="label">
                        <span class="icon">🌍</span>
                        <span>Dil</span>
                    </div>
                    <select v-model="selectedLang" @change="changeLang" class="lang-select">
                        <option value="tr">Türkçe</option>
                        <option value="en">English</option>
                        <option value="de">Deutsch</option>
                        <option value="fr">Français</option>
                    </select>
                </div>
            </div>

            <div class="section-card">
                <h3>Profil Bilgileri</h3>
                <div class="input-group">
                    <label>Ad Soyad</label>
                    <input v-model="fullName" type="text" />
                </div>
                <div class="input-group">
                    <label>E-Posta</label>
                    <input v-model="email" type="email" />
                </div>
                <div class="input-group">
                    <label>Yeni Şifre</label>
                    <input v-model="password" type="password" placeholder="Değiştirmek istemiyorsanız boş bırakın" />
                </div>
                <button class="save-btn" @click="updateProfile">Bilgileri Güncelle</button>
            </div>

            <div class="section-card danger">
                <h3>Hesap İşlemleri</h3>
                <p>Hesabınızı silerseniz tüm verileriniz kalıcı olarak kaybolur.</p>
                <button class="delete-btn" @click="deleteAccount">Hesabımı Sil</button>
            </div>
        </div>
    </div>
</template>

<script setup>
    import { ref, inject, onMounted } from 'vue';
    import { useRouter } from 'vue-router';
    import api from '../services/api';
    import { currentLang } from '../utils/translations';

    const router = useRouter();
    const { theme, toggleTheme } = inject('theme');
    const fullName = ref('');
    const email = ref('');
    const password = ref('');
    const selectedLang = ref(currentLang.value);

    onMounted(() => {
        fullName.value = localStorage.getItem('username') || '';
    });

    const changeLang = () => {
        currentLang.value = selectedLang.value;
        localStorage.setItem('lang', selectedLang.value);
        window.location.reload();
    };

    const updateProfile = async () => {
        try {
            const payload = { full_name: fullName.value, email: email.value };
            if (password.value) payload.password = password.value;
            await api.put('/user/update', payload);
            localStorage.setItem('username', fullName.value);
            alert('Profil güncellendi!');
            password.value = '';
        } catch (e) {
            alert('Hata: ' + (e.response?.data?.error || e.message));
        }
    };

    const deleteAccount = async () => {
        if (confirm('Hesabınızı silmek istediğinize emin misiniz?')) {
            try {
                await api.delete('/user/delete');
                localStorage.clear();
                router.push('/login');
            } catch (e) {
                alert('Silme işlemi başarısız.');
            }
        }
    };
</script>

<style scoped>
    .settings-content { padding: 30px; width: 100%; max-width: 800px; margin: 0 auto; color: var(--text-color); }
    .settings-header h2 { margin-bottom: 30px; font-size: 1.8rem; }
    .settings-grid { display: grid; gap: 20px; }
    .section-card { background: var(--card-bg); padding: 25px; border-radius: 15px; border: 1px solid var(--border-color); }
    .section-card h3 { margin-top: 0; margin-bottom: 20px; color: var(--accent-color); font-size: 1.1rem; }
    .setting-row { display: flex; justify-content: space-between; align-items: center; padding: 15px 0; border-bottom: 1px solid var(--border-color); cursor: pointer; }
    .setting-row:last-child { border-bottom: none; }
    .label { display: flex; align-items: center; gap: 10px; font-weight: 500; }
    .value { color: var(--text-muted); }
    .lang-select { background: var(--input-bg); color: var(--text-color); border: 1px solid var(--border-color); padding: 5px 10px; border-radius: 5px; }
    .input-group { margin-bottom: 15px; }
    .input-group label { display: block; margin-bottom: 8px; font-size: 0.9rem; color: var(--text-muted); }
    .input-group input { width: 100%; padding: 12px; background: var(--input-bg); border: 1px solid var(--border-color); color: var(--text-color); border-radius: 8px; box-sizing: border-box; }
    .save-btn { width: 100%; padding: 12px; background: var(--success-color); color: white; border: none; border-radius: 8px; font-weight: bold; cursor: pointer; margin-top: 10px; }
    .danger { border-color: var(--danger-color); }
    .danger h3 { color: var(--danger-color); }
    .danger p { color: var(--text-muted); font-size: 0.9rem; margin-bottom: 15px; }
    .delete-btn { background: var(--danger-color); color: white; padding: 10px 20px; border: none; border-radius: 8px; cursor: pointer; font-weight: bold; }
</style>