<template>
  <div class="login-page">
    <div class="login-card">
      <div class="brand">FinTrack Pro</div>

      <div class="tabs">
        <div :class="{ active: activeTab === 'login' }" @click="activeTab = 'login'">Giriş Yap</div>
        <div :class="{ active: activeTab === 'register' }" @click="activeTab = 'register'">Kayıt Ol</div>
      </div>

      <div v-if="activeTab === 'login'" class="form-section">
        <input v-model="loginData.email" type="email" placeholder="E-Posta" />
        <input v-model="loginData.password" type="password" placeholder="Şifre" />
        <button @click="handleLogin">GİRİŞ YAP</button>
      </div>

      <div v-else class="form-section">
        <input v-model="registerData.fullName" type="text" placeholder="Ad Soyad" />
        <input v-model="registerData.username" type="text" placeholder="Kullanıcı Adı" />
        <input v-model="registerData.email" type="email" placeholder="E-Posta" />
        <input v-model="registerData.password" type="password" placeholder="Şifre" />
        <button @click="handleRegister">KAYIT OL</button>
      </div>
    </div>
  </div>
</template>

<script setup>
  import { ref, reactive } from 'vue';
  import { useRouter } from 'vue-router';
  import api from '../services/api';

  const router = useRouter();
  const activeTab = ref('login');

  const loginData = reactive({ email: '', password: '' });
  const registerData = reactive({ fullName: '', username: '', email: '', password: '' });

  const handleLogin = async () => {
    try {
      const res = await api.post('/auth/login', {
        email: loginData.email,
        password: loginData.password
      });

      localStorage.setItem('token', res.data.token);
      localStorage.setItem('username', res.data.username || 'Kullanıcı');

      router.push('/dashboard');
    } catch (e) {
      alert('Hata: ' + (e.response?.data?.error || 'Giriş yapılamadı.'));
    }
  };

  const handleRegister = async () => {
    try {
      await api.post('/auth/register', {
        full_name: registerData.fullName,
        username: registerData.username,
        email: registerData.email,
        password: registerData.password
      });

      activeTab.value = 'login';
      registerData.fullName = '';
      registerData.username = '';
      registerData.email = '';
      registerData.password = '';
    } catch (e) {
      alert('Hata: ' + (e.response?.data?.error || 'Kayıt yapılamadı.'));
    }
  };
</script>

<style scoped>
  .login-page { display: flex; justify-content: center; align-items: center; min-height: 100vh; background: #0F172A; font-family: 'Segoe UI', sans-serif; }
  .login-card { background: #1E293B; padding: 40px; border-radius: 20px; width: 350px; border: 1px solid rgba(255,255,255,0.1); box-shadow: 0 20px 50px rgba(0,0,0,0.5); }
  .brand { color: #FFD700; font-size: 1.8rem; font-weight: 800; text-align: center; margin-bottom: 30px; letter-spacing: 1px; }
  .tabs { display: flex; margin-bottom: 25px; border-bottom: 1px solid rgba(255,255,255,0.1); }
  .tabs div { flex: 1; text-align: center; padding: 10px; cursor: pointer; color: #94A3B8; transition: 0.3s; font-weight: bold; }
  .tabs div.active { color: #FFD700; border-bottom: 2px solid #FFD700; }
  .form-section { display: flex; flex-direction: column; gap: 15px; }
  input { padding: 12px; border-radius: 10px; border: 1px solid rgba(255,255,255,0.1); background: #0F172A; color: white; font-size: 1rem; outline: none; transition: 0.3s; }
  input:focus { border-color: #FFD700; }
  button { padding: 15px; border-radius: 10px; border: none; background: #FFD700; color: #000; font-weight: bold; font-size: 1rem; cursor: pointer; transition: 0.3s; margin-top: 10px; }
  button:hover { transform: translateY(-2px); box-shadow: 0 5px 15px rgba(255, 215, 0, 0.3); }
</style>