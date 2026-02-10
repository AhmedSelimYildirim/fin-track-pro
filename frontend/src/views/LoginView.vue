<template>
  <div class="login-container">
    <div class="login-card">
      <h2 class="title">{{ isRegister ? 'Register' : 'Login' }}</h2>

      <form @submit.prevent="handleAuth">
        <div class="input-group" v-if="isRegister">
          <input v-model="fullName" type="text" placeholder="Ad Soyad" required />
        </div>

        <div class="input-group">
          <input v-model="email" type="email" placeholder="Email" required />
        </div>

        <div class="input-group">
          <input v-model="password" type="password" placeholder="Şifre" required />
        </div>

        <button type="submit" :disabled="loading" class="auth-btn">
          {{ loading ? 'İşleniyor...' : (isRegister ? 'Kayıt Ol' : 'Giriş Yap') }}
        </button>
      </form>

      <p class="toggle-text" @click="isRegister = !isRegister">
        {{ isRegister ? 'Zaten hesabın var mı? Giriş Yap' : "Hesabın yok mu? Kayıt Ol" }}
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import api from '../services/api'; // <--- İŞTE KRİTİK NOKTA BURASI! "axios" değil "api" import ediyoruz.

const router = useRouter();
const isRegister = ref(false);
const loading = ref(false);
const fullName = ref('');
const email = ref('');
const password = ref('');

const handleAuth = async () => {
  loading.value = true;
  try {
    if (isRegister.value) {
      // DİKKAT: Artık "axios.post" değil "api.post" kullanıyoruz.
      // api.js içindeki baseURL sayesinde linkin başını yazmana gerek yok.
      await api.post('/auth/register', {
        full_name: fullName.value,
        email: email.value,
        password: password.value
      });
      alert("Kayıt başarılı! Şimdi giriş yapabilirsin.");
      isRegister.value = false;
    } else {
      const res = await api.post('/auth/login', {
        email: email.value,
        password: password.value
      });

      localStorage.setItem('token', res.data.token);
      localStorage.setItem('username', res.data.username || 'Kullanıcı');

      // Giriş başarılı, Dashboard'a git
      router.push('/dashboard');
    }
  } catch (e) {
    console.error(e);
    alert(e.response?.data?.error || "İşlem başarısız! Bilgileri kontrol et.");
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
.login-container { display: flex; justify-content: center; align-items: center; height: 100vh; background-color: #0F172A; }
.login-card { background: #1E293B; padding: 2rem; border-radius: 15px; width: 100%; max-width: 400px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); text-align: center; }
.title { color: white; margin-bottom: 1.5rem; font-size: 2rem; font-weight: bold; }
.input-group { margin-bottom: 1rem; }
.input-group input { width: 100%; padding: 12px; border-radius: 8px; border: 1px solid #334155; background: #0F172A; color: white; box-sizing: border-box; }
.auth-btn { width: 100%; padding: 12px; background: #EF4444; color: white; border: none; border-radius: 8px; font-weight: bold; cursor: pointer; transition: 0.3s; }
.auth-btn:hover { background: #DC2626; }
.toggle-text { color: #94A3B8; margin-top: 1rem; cursor: pointer; font-size: 0.9rem; }
.toggle-text:hover { color: white; text-decoration: underline; }
</style>