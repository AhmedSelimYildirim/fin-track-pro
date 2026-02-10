<template>
    <div class="settings-page">
        <div class="settings-container">
            <div class="header">
                <button class="back-btn" @click="router.push('/dashboard')">← Geri</button>
                <h2>Ayarlar</h2>
            </div>

            <div class="section">
                <h3>Görünüm</h3>
                <div class="theme-toggle" @click="toggleTheme">
                    <span>{{ theme === 'dark' ? '🌙 Karanlık Mod' : '☀️ Aydınlık Mod' }}</span>
                    <div class="switch" :class="{ active: theme === 'light' }">
                        <div class="slider"></div>
                    </div>
                </div>
            </div>

            <div class="section">
                <h3>Profil Ayarları</h3>
                <div class="form-group">
                    <label>Ad Soyad</label>
                    <input v-model="fullName" type="text" placeholder="Ad Soyad" />
                </div>
                <div class="form-group">
                    <label>Email</label>
                    <input v-model="email" type="email" placeholder="Email" />
                </div>
                <button class="save-btn" @click="updateProfile">Bilgileri Güncelle</button>
            </div>

            <div class="section danger-zone">
                <h3>Tehlikeli Bölge</h3>
                <p>Hesabını silersen tüm verilerin kaybolur ve geri getirilemez.</p>
                <button class="delete-btn" @click="deleteAccount">Hesabımı Sil</button>
            </div>
        </div>
    </div>
</template>

<script setup>
    import { ref, inject, onMounted } from 'vue';
    import { useRouter } from 'vue-router';
    import api from '../services/api';

    const router = useRouter();
    const { theme, toggleTheme } = inject('theme');
    const fullName = ref('');
    const email = ref('');

    onMounted(() => {
        fullName.value = localStorage.getItem('username') || '';
    });

    const updateProfile = async () => {
        try {
            await api.put('/user/update', { full_name: fullName.value, email: email.value });
            localStorage.setItem('username', fullName.value);
            alert('Profil güncellendi!');
        } catch (e) {
            alert('Güncelleme henüz aktif değil veya hata oluştu.');
        }
    };

    const deleteAccount = async () => {
        if (confirm('Hesabını silmek istediğine emin misin? Bu işlem geri alınamaz!')) {
            try {
                await api.delete('/user/delete');
                localStorage.clear();
                router.push('/login');
            } catch (e) {
                alert('Silme işlemi sırasında hata oluştu.');
            }
        }
    };
</script>

<style scoped>
    .settings-page { min-height: 100vh; background-color: var(--bg-color); display: flex; justify-content: center; padding: 40px; }
    .settings-container { width: 100%; max-width: 600px; }
    .header { display: flex; align-items: center; gap: 20px; margin-bottom: 30px; }
    .back-btn { background: none; border: none; color: var(--text-color); cursor: pointer; font-size: 1rem; }
    h2, h3 { color: var(--text-color); }
    .section { background: var(--card-bg); padding: 25px; border-radius: 15px; margin-bottom: 20px; border: 1px solid var(--border-color); }

    .theme-toggle { display: flex; justify-content: space-between; align-items: center; cursor: pointer; padding: 10px; border-radius: 8px; background: var(--input-bg); color: var(--text-color); }
    .switch { width: 50px; height: 26px; background: #334155; border-radius: 20px; position: relative; transition: 0.3s; }
    .switch.active { background: var(--success-color); }
    .slider { width: 20px; height: 20px; background: white; border-radius: 50%; position: absolute; top: 3px; left: 3px; transition: 0.3s; }
    .switch.active .slider { left: 27px; }

    .form-group { margin-bottom: 15px; }
    label { display: block; margin-bottom: 5px; color: var(--text-muted); }
    input { width: 100%; padding: 12px; background: var(--input-bg); border: 1px solid var(--border-color); color: var(--text-color); border-radius: 8px; box-sizing: border-box; }

    .save-btn { width: 100%; padding: 12px; background: #3B82F6; color: white; border: none; border-radius: 8px; font-weight: bold; cursor: pointer; }
    .danger-zone { border: 1px solid var(--danger-color); }
    .danger-zone h3 { color: var(--danger-color); }
    .danger-zone p { color: var(--text-muted); margin-bottom: 15px; }
    .delete-btn { width: 100%; padding: 12px; background: var(--danger-color); color: white; border: none; border-radius: 8px; font-weight: bold; cursor: pointer; }
</style>