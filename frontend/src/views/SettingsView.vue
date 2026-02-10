<template>
    <div class="settings-wrapper">
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
                </div>

                <div class="section-card">
                    <h3>{{ t('profileSettings') }}</h3>
                    <div class="input-group">
                        <label>{{ t('username') }}</label>
                        <input v-model="username" type="text" placeholder="Kullanıcı Adı" />
                    </div>
                    <div class="input-group">
                        <label>{{ t('email') }}</label>
                        <input v-model="email" type="email" placeholder="E-Posta (Zorunlu)" />
                    </div>
                    <div class="input-group">
                        <label>{{ t('password') }}</label>
                        <input v-model="password" type="password" :placeholder="t('passwordHint')" />
                    </div>
                    <button class="save-btn" :disabled="isLoading" @click="updateProfile">
                        <span v-if="!isLoading">{{ t('update') }}</span>
                        <span v-else class="spinner"></span>
                    </button>
                </div>

                <div class="section-card danger">
                    <h3>{{ t('dangerZone') }}</h3>
                    <button class="delete-btn" @click="deleteAccount">{{ t('deleteAccount') }}</button>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
    import { ref, inject, onMounted } from 'vue'
    import { useRouter } from 'vue-router'
    import api from '../services/api'
    import { t } from '../utils/translations'

    const router = useRouter()
    const { theme, toggleTheme } = inject('theme')

    const username = ref('')
    const email = ref('')
    const password = ref('')
    const isLoading = ref(false)

    // Backend'den en güncel veriyi çekmeye çalış, yoksa localStorage'a bak
    const loadUserData = async () => {
        try {
            // Önce API'den çekmeyi dene (Endpoint tahminidir, kendi endpointine göre düzenle)
            const res = await api.get('/user/me') // veya /auth/me
            if (res.data) {
                username.value = res.data.username
                email.value = res.data.email
            }
        } catch (e) {
            console.warn("API'den veri çekilemedi, LocalStorage kullanılıyor.")
            // API hatası olursa LocalStorage'dan yükle
            username.value = localStorage.getItem('username') || ''
            email.value = localStorage.getItem('email') || ''
        }

        // --- CRITICAL FALLBACK LOGIC ---
        // Kullanıcı adı boşsa ve Email varsa, Username'i Email'den türet
        if (!username.value && email.value) {
            username.value = email.value.split('@')[0]
        }
    }

    onMounted(loadUserData)

    const updateProfile = async () => {
        if (!email.value) {
            alert("E-Posta alanı zorunludur!")
            return
        }

        isLoading.value = true
        try {
            const payload = {
                username: username.value,
                email: email.value
            }

            if (password.value) payload.password = password.value

            // API çağrısı
            await api.put('/user/update', payload)

            // Başarılı olursa LocalStorage'ı güncelle
            localStorage.setItem('username', username.value)
            localStorage.setItem('email', email.value)

            password.value = ''
            alert("Profil başarıyla güncellendi!")
        } catch (e) {
            alert('Güncelleme hatası: ' + (e.response?.data?.error || e.message))
        } finally {
            isLoading.value = false
        }
    }

    const deleteAccount = async () => {
        if (confirm(t('deleteWarning'))) {
            try {
                await api.delete('/user/delete')
                localStorage.clear()
                router.push('/login')
            } catch {
                alert('Silme işlemi başarısız.')
            }
        }
    }
</script>

<style scoped>
    .settings-wrapper { width: 100%; min-height: 100%; padding: 40px; box-sizing: border-box; }
    .settings-content { max-width: 800px; margin: 0 auto; }

    .settings-header h2 {
        font-size: 2rem; margin-bottom: 30px; color: white; font-weight: 800;
    }

    .settings-grid { display: grid; gap: 24px; }

    .section-card {
        background: rgba(30, 41, 59, 0.5);
        padding: 24px;
        border-radius: 16px;
        border: 1px solid rgba(255, 255, 255, 0.1);
        backdrop-filter: blur(10px);
    }

    .section-card h3 {
        margin-top: 0; margin-bottom: 20px; color: #facc15; font-size: 1.1rem;
        border-bottom: 1px solid rgba(255,255,255,0.1); padding-bottom: 12px;
    }

    .setting-item { display: flex; justify-content: space-between; align-items: center; }
    .item-info { display: flex; align-items: center; gap: 15px; color: white; }
    .icon { font-size: 1.5rem; }

    .theme-switch {
        width: 56px; height: 32px; background: #334155; border-radius: 99px;
        cursor: pointer; position: relative; transition: 0.3s;
    }
    .theme-switch.is-light { background: #facc15; }

    .switch-handle {
        width: 26px; height: 26px; background: white; border-radius: 50%;
        position: absolute; top: 3px; left: 3px; transition: 0.3s;
        display: flex; align-items: center; justify-content: center; font-size: 0.8rem;
    }
    .is-light .switch-handle { transform: translateX(24px); }

    .input-group { margin-bottom: 16px; }
    .input-group label { display: block; margin-bottom: 8px; font-size: 0.9rem; color: #94a3b8; }
    .input-group input {
        width: 100%; padding: 12px; background: rgba(2, 6, 23, 0.5);
        border: 1px solid rgba(255, 255, 255, 0.1); color: white; border-radius: 8px; box-sizing: border-box;
    }
    .input-group input:focus { border-color: #22c55e; outline: none; }

    .save-btn {
        width: 100%; padding: 14px; background: #22c55e; color: white;
        border: none; border-radius: 8px; font-weight: bold; cursor: pointer; transition: 0.2s;
    }
    .save-btn:hover:not(:disabled) { background: #16a34a; }
    .save-btn:disabled { opacity: 0.6; cursor: wait; }

    .danger { border: 1px solid rgba(239, 68, 68, 0.3); }
    .danger h3 { color: #ef4444; border-color: rgba(239, 68, 68, 0.2); }
    .delete-btn {
        background: rgba(239, 68, 68, 0.2); color: #ef4444; padding: 12px 20px;
        border: 1px solid #ef4444; border-radius: 8px; cursor: pointer; font-weight: bold; transition: 0.2s;
    }
    .delete-btn:hover { background: #ef4444; color: white; }

    .spinner {
        display: inline-block; width: 20px; height: 20px;
        border: 3px solid rgba(255,255,255,0.3); border-top-color: white;
        border-radius: 50%; animation: spin 1s linear infinite;
    }
    @keyframes spin { to { transform: rotate(360deg); } }
</style>