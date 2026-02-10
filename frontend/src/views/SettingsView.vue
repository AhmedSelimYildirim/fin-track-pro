<template>
    <div class="settings-page">
        <div class="settings-container">
            <div class="header">
                <h2>{{ t('settings') }}</h2>
            </div>

            <div class="section">
                <h3>{{ t('appearance') }}</h3>

                <div class="setting-item" @click="toggleTheme">
                    <span>{{ theme === 'dark' ? t('darkMode') : t('lightMode') }}</span>
                    <div class="switch" :class="{ active: theme === 'light' }"><div class="slider"></div></div>
                </div>

                <div class="setting-item mt-20">
                    <span>{{ t('language') }}</span>
                    <select v-model="selectedLang" @change="changeLang" class="lang-select">
                        <option value="tr">Türkçe 🇹🇷</option>
                        <option value="en">English 🇬🇧</option>
                        <option value="de">Deutsch 🇩🇪</option>
                        <option value="fr">Français 🇫🇷</option>
                    </select>
                </div>
            </div>

            <div class="section">
                <h3>{{ t('profileSettings') }}</h3>
                <div class="form-group">
                    <input v-model="fullName" type="text" placeholder="Ad Soyad" />
                </div>
                <div class="form-group">
                    <input v-model="email" type="email" placeholder="Email" />
                </div>
                <div class="form-group">
                    <input v-model="password" type="password" placeholder="Yeni Şifre (İsteğe Bağlı)" />
                </div>
                <button class="save-btn" @click="updateProfile">{{ t('update') }}</button>
            </div>

            <div class="section danger-zone">
                <h3>{{ t('dangerZone') }}</h3>
                <button class="delete-btn" @click="deleteAccount">{{ t('deleteAccount') }}</button>
            </div>
        </div>
    </div>
</template>

<script setup>
    import { ref, inject, onMounted } from 'vue';
    import { useRouter } from 'vue-router';
    import api from '../services/api';
    import { t, currentLang } from '../utils/translations';

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
        // Sayfayı yenile ki menüler de güncellensin
        window.location.reload();
    };

    const updateProfile = async () => {
        try {
            const payload = { full_name: fullName.value, email: email.value };
            if (password.value) payload.password = password.value;
            await api.put('/user/update', payload);
            localStorage.setItem('username', fullName.value);
            alert('Güncellendi!');
            window.location.reload(); // İsim güncellensin diye
        } catch (e) { alert('Hata oluştu.'); }
    };

    const deleteAccount = async () => {
        if (confirm('Emin misiniz?')) {
            try {
                await api.delete('/user/delete');
                localStorage.clear();
                router.push('/login');
            } catch (e) { alert('Hata.'); }
        }
    };
</script>

<style scoped>
    .settings-page { padding: 40px; display: flex; justify-content: center; }
    .settings-container { width: 100%; max-width: 600px; }
    h2, h3 { color: var(--text-color); margin-bottom: 20px; }
    .section { background: var(--card-bg); padding: 25px; border-radius: 15px; margin-bottom: 20px; border: 1px solid var(--border-color); }
    .setting-item { display: flex; justify-content: space-between; align-items: center; cursor: pointer; padding: 10px; border-radius: 8px; background: var(--input-bg); border: 1px solid var(--border-color); color: var(--text-color); }
    .mt-20 { margin-top: 15px; }
    .lang-select { background: transparent; border: none; color: var(--text-color); font-size: 1rem; outline: none; }
    .switch { width: 50px; height: 26px; background: #334155; border-radius: 20px; position: relative; transition: 0.3s; }
    .switch.active { background: var(--success-color); }
    .slider { width: 20px; height: 20px; background: white; border-radius: 50%; position: absolute; top: 3px; left: 3px; transition: 0.3s; }
    .switch.active .slider { left: 27px; }
    .form-group { margin-bottom: 15px; }
    input { width: 100%; padding: 12px; background: var(--input-bg); border: 1px solid var(--border-color); color: var(--text-color); border-radius: 8px; box-sizing: border-box; }
    .save-btn { width: 100%; padding: 12px; background: #3B82F6; color: white; border: none; border-radius: 8px; font-weight: bold; cursor: pointer; }
    .danger-zone { border: 1px solid var(--danger-color); }
    .danger-zone h3 { color: var(--danger-color); }
    .delete-btn { width: 100%; padding: 12px; background: var(--danger-color); color: white; border: none; border-radius: 8px; font-weight: bold; cursor: pointer; }
</style>