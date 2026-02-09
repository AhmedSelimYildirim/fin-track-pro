<template>
  <div class="login-container">
    <div class="auth-card">
      <h1 class="title">{{ isLogin ? 'Login' : 'Register' }}</h1>

      <form @submit.prevent="handleAuth" class="auth-form">
        <div v-if="!isLogin" class="input-group">
          <input v-model="username" type="text" placeholder="Username" required />
        </div>
        <div class="input-group">
          <input v-model="email" type="email" placeholder="Email" required />
        </div>
        <div class="input-group">
          <input v-model="password" type="password" placeholder="Password" required />
        </div>

        <button type="submit" class="submit-btn">
          {{ isLogin ? 'Sign In' : 'Create Account' }}
        </button>
      </form>

      <p class="toggle-text" @click="isLogin = !isLogin">
        {{ isLogin ? "Don't have an account? Sign Up" : "Already have an account? Login" }}
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import api from '../services/api';

const isLogin = ref(true);
const email = ref('');
const password = ref('');
const username = ref('');
const router = useRouter();

const handleAuth = async () => {
  try {
    const endpoint = isLogin ? '/auth/login' : '/auth/register';
    const payload = isLogin ? { email, password } : { username, email, password };

    const res = await api.post(endpoint, payload);
    localStorage.setItem('token', res.data.token);
    // Kullanıcı adını kaydet, dashboard'da göstereceğiz
    localStorage.setItem('username', res.data.username || username.value);
    localStorage.setItem('email', email.value);

    router.push('/dashboard');
  } catch (error) {
    alert('İşlem başarısız! Bilgileri kontrol et.');
  }
};
</script>

<style scoped>
.login-container { height: 100vh; display: flex; justify-content: center; align-items: center; }
.auth-card { background: var(--card-bg); padding: 40px; border-radius: 20px; width: 350px; text-align: center; box-shadow: 0 10px 30px rgba(0,0,0,0.5); border: 1px solid #334155; }
.title { margin-bottom: 30px; font-size: 2rem; color: var(--text-white); }
.input-group input { width: 100%; padding: 12px; margin: 10px 0; border-radius: 8px; border: 1px solid #475569; background: #0F172A; color: white; box-sizing: border-box; }
.submit-btn { width: 100%; padding: 12px; margin-top: 20px; background-color: var(--try-color); color: white; border: none; border-radius: 8px; font-weight: bold; font-size: 1rem; }
.toggle-text { margin-top: 20px; color: #94A3B8; cursor: pointer; font-size: 0.9rem; }
.toggle-text:hover { color: white; }
</style>