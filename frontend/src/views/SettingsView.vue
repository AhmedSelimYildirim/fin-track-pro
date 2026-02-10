<template>
    <div class="settings-content">
        <div class="settings-header">
            <h2>{{ t('settings') }}</h2>
        </div>

        <div class="settings-grid">
            <div class="section-card">
                <h3>{{ t('appearance') }}</h3>

                <div class="setting-item">
                    <div class="item-info">
                        <span class="icon">🌓</span>
                        <div>
                            <span class="label">{{ t('theme') }}</span>
                            <span class="sub-label">{{ theme === 'dark' ? t('darkMode') : t('lightMode') }}</span>
                        </div>
                    </div>
                    <div class="theme-switch" :class="{ 'is-light': theme === 'light' }" @click="toggleTheme">
                        <div class="switch-handle">
                            <span v-if="theme === 'dark'">🌙</span>
                            <span v-else>☀️</span>
                        </div>
                    </div>
                </div>

                <div class="setting-item mt-20">
                    <div class="item-info">
                        <span class="icon">🌍</span>
                        <span class="label">{{ t('language') }}</span>
                    </div>
                    <div class="lang-buttons">
                        <button v-for="lang in languages" :key="lang.code"
                                :class="{ active: selectedLang === lang.code }"
                                @click="changeLang(lang.code)">
                            {{ lang.flag }} {{ lang.name }}
                        </button>
                    </div>
                </div>
            </div>

            <div class="section-card">
                <h3>{{ t('profileSettings') }}</h3>
                <div class="input-group">
                    <label>{{ t('username') }}</label>
                    <input v-model="username" type="text" />
                </div>
                <div class="input-group">
                    <label>{{ t('email') }}</label>
                    <input v-model="email" type="email" />
                </div>
                <div class="input-group">
                    <label>{{ t('password') }}</label>
                    <input v-model="password" type="password" :placeholder="t('passwordHint')" />
                </div>
                <button class="save-btn" :disabled="isLoading" @click="updateProfile">
                    <span v-if="!isLoading">{{ t('update') }}</span>
                    <span v-else>...</span>
                </button>
            </div>

            <div class="section-card danger">
                <h3>{{ t('dangerZone') }}</h3>
                <p>{{ t('deleteWarning') }}</p>
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
    const username = ref('');
    const email = ref('');
    const password = ref('');
    const selectedLang = ref(currentLang.value);
    const isLoading = ref(false);

    const languages = [
        { code: 'tr', name: 'Türkçe', flag: '🇹🇷' },
        { code: 'en', name: 'English', flag: '🇬🇧' },
        { code: 'de', name: 'Deutsch', flag: '🇩🇪' },
        { code: 'fr', name: 'Français', flag: '🇫🇷' },
        { code: 'ar', name: 'العربية', flag: '🇸🇦' }
    ];

    onMounted(() => {
        username.value = localStorage.getItem('username') || '';
    });

    const changeLang = (code) => {
        selectedLang.value = code;
        currentLang.value = code;
        localStorage.setItem('lang', code);
        window.location.reload();
    };

    const updateProfile = async () => {
        isLoading.value = true;
        try {
            const payload = {
                username: username.value,
                email: email.value
            };
            if (password.value) payload.password = password.value;

            await api.put('/user/update', payload);
            localStorage.setItem('username', username.value);
            password.value = '';
            alert("Profil güncellendi!");
        } catch (e) {
            alert('Hata: ' + (e.response?.data?.error || e.message));
        } finally {
            isLoading.value = false;
        }
    };

    const deleteAccount = async () => {
        if (confirm(t('deleteWarning'))) {
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
    .section-card h3 { margin-top: 0; margin-bottom: 20px; color: var(--accent-color); font-size: 1.1rem; border-bottom: 1px solid var(--border-color); padding-bottom: 10px; }
    .setting-item { display: flex; justify-content: space-between; align-items: center; padding: 10px 0; }
    .item-info { display: flex; align-items: center; gap: 15px; }
    .icon { font-size: 1.4rem; }
    .label { display: block; font-weight: bold; font-size: 1rem; }
    .sub-label { font-size: 0.8rem; color: var(--text-muted); }
    .theme-switch { width: 60px; height: 30px; background: #334155; border-radius: 20px; cursor: pointer; position: relative; transition: 0.3s; }
    .theme-switch.is-light { background: #FFD700; }
    .switch-handle { width: 24px; height: 24px; background: white; border-radius: 50%; position: absolute; top: 3px; left: 3px; transition: 0.3s; display: flex; align-items: center; justify-content: center; font-size: 0.8rem; }
    .is-light .switch-handle { left: 33px; }
    .lang-buttons { display: flex; gap: 8px; flex-wrap: wrap; margin-top: 10px; }
    .lang-buttons button { background: var(--input-bg); border: 1px solid var(--border-color); color: var(--text-color); padding: 8px 12px; border-radius: 8px; cursor: pointer; transition: 0.2s; font-size: 0.9rem; }
    .lang-buttons button.active { background: var(--accent-color); color: #000; font-weight: bold; border-color: var(--accent-color); }
    .input-group { margin-bottom: 15px; }
    .input-group label { display: block; margin-bottom: 8px; font-size: 0.9rem; color: var(--text-muted); }
    .input-group input { width: 100%; padding: 12px; background: var(--input-bg); border: 1px solid var(--border-color); color: var(--text-color); border-radius: 8px; box-sizing: border-box; }
    .save-btn { width: 100%; padding: 12px; background: var(--success-color); color: white; border: none; border-radius: 8px; font-weight: bold; cursor: pointer; margin-top: 10px; transition: 0.2s; display: flex; justify-content: center; }
    .save-btn:hover { opacity: 0.9; transform: translateY(-1px); }
    .danger { border-color: var(--danger-color); }
    .danger h3 { color: var(--danger-color); border-color: rgba(239, 68, 68, 0.2); }
    .danger p { color: var(--text-muted); font-size: 0.9rem; margin-bottom: 15px; }
    .delete-btn { background: var(--danger-color); color: white; padding: 10px 20px; border: none; border-radius: 8px; cursor: pointer; font-weight: bold; transition: 0.2s; }
    .delete-btn:hover { background: #DC2626; }
</style>