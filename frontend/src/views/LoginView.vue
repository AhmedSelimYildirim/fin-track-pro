<template>
  <div class="login-container">
    <div class="background-shapes">
      <div class="shape shape-1"></div>
      <div class="shape shape-2"></div>
      <div class="shape shape-3"></div>
    </div>

    <div class="glass-card">
      <div class="card-header">
        <h1 class="brand-title">FinTrack Pro 🚀</h1>
        <p class="subtitle">Finansal Özgürlüğe Açılan Kapı</p>
      </div>

      <div class="toggle-container">
        <button :class="{ active: activeTab === 'login' }" @click="activeTab = 'login'">Giriş</button>
        <button :class="{ active: activeTab === 'register' }" @click="activeTab = 'register'">Kayıt</button>
      </div>

      <div class="form-container">
        <transition name="fade" mode="out-in">
          <div v-if="activeTab === 'login'" key="login" class="form-group">
            <div class="input-wrapper">
              <span class="input-icon">📧</span>
              <input v-model="loginData.email" type="email" placeholder="E-Posta Adresi" />
            </div>
            <div class="input-wrapper">
              <span class="input-icon">🔒</span>
              <input v-model="loginData.password" type="password" placeholder="Şifre" />
            </div>
            <button class="action-btn" @click="handleLogin">
              Giriş Yap ➔
            </button>
          </div>

          <div v-else key="register" class="form-group">
            <div class="input-wrapper">
              <span class="input-icon">👤</span>
              <input v-model="registerData.username" type="text" placeholder="Kullanıcı Adı" />
            </div>
            <div class="input-wrapper">
              <span class="input-icon">📧</span>
              <input v-model="registerData.email" type="email" placeholder="E-Posta Adresi" />
            </div>
            <div class="input-wrapper">
              <span class="input-icon">🔑</span>
              <input v-model="registerData.password" type="password" placeholder="Şifre" />
            </div>
            <button class="action-btn register-btn" @click="handleRegister">
              Hesap Oluştur ✨
            </button>
          </div>
        </transition>
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
  const registerData = reactive({ username: '', email: '', password: '' });

  const handleLogin = async () => {
    if (!loginData.email || !loginData.password) return alert("Lütfen tüm alanları doldurun.");
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
    if (!registerData.username || !registerData.email || !registerData.password) return alert("Lütfen tüm alanları doldurun.");
    try {
      await api.post('/auth/register', {
        username: registerData.username,
        email: registerData.email,
        password: registerData.password
      });
      activeTab.value = 'login';
      registerData.username = '';
      registerData.email = '';
      registerData.password = '';
    } catch (e) {
      alert('Hata: ' + (e.response?.data?.error || 'Kayıt yapılamadı.'));
    }
  };
</script>

<style scoped>
  .login-container {
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 100vh;
    width: 100vw;
    background-color: #0F172A;
    position: relative;
    overflow: hidden;
    font-family: 'Inter', sans-serif;
  }

  .background-shapes .shape {
    position: absolute;
    filter: blur(90px);
    z-index: 0;
    opacity: 0.6;
    animation: float 10s infinite ease-in-out;
  }

  .shape-1 { width: 300px; height: 300px; background: #3B82F6; top: -50px; left: -50px; }
  .shape-2 { width: 400px; height: 400px; background: #8B5CF6; bottom: -100px; right: -100px; animation-delay: 2s; }
  .shape-3 { width: 200px; height: 200px; background: #10B981; top: 40%; left: 40%; animation-delay: 4s; }

  @keyframes float {
    0%, 100% { transform: translate(0, 0); }
    50% { transform: translate(30px, -30px); }
  }

  .glass-card {
    position: relative;
    z-index: 1;
    background: rgba(30, 41, 59, 0.7);
    backdrop-filter: blur(20px);
    -webkit-backdrop-filter: blur(20px);
    border: 1px solid rgba(255, 255, 255, 0.1);
    padding: 50px 40px;
    border-radius: 24px;
    width: 100%;
    max-width: 420px;
    box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5);
    display: flex;
    flex-direction: column;
    align-items: center;
  }

  .card-header { margin-bottom: 30px; text-align: center; }
  .brand-title { color: #fff; font-size: 2rem; font-weight: 800; margin: 0; background: linear-gradient(to right, #3B82F6, #10B981); -webkit-background-clip: text; -webkit-text-fill-color: transparent; }
  .subtitle { color: #94A3B8; margin-top: 5px; font-size: 0.95rem; }

  .toggle-container {
    display: flex;
    background: rgba(15, 23, 42, 0.6);
    padding: 5px;
    border-radius: 12px;
    width: 100%;
    margin-bottom: 30px;
  }

  .toggle-container button {
    flex: 1;
    padding: 10px;
    border: none;
    background: transparent;
    color: #94A3B8;
    cursor: pointer;
    border-radius: 8px;
    font-weight: 600;
    transition: 0.3s;
  }

  .toggle-container button.active {
    background: #3B82F6;
    color: white;
    box-shadow: 0 4px 6px -1px rgba(59, 130, 246, 0.3);
  }

  .form-container { width: 100%; }
  .form-group { display: flex; flex-direction: column; gap: 15px; }

  .input-wrapper {
    position: relative;
    display: flex;
    align-items: center;
  }

  .input-icon {
    position: absolute;
    left: 15px;
    font-size: 1.2rem;
    z-index: 10;
  }

  input {
    width: 100%;
    padding: 14px 14px 14px 45px;
    background: rgba(15, 23, 42, 0.6);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 12px;
    color: white;
    font-size: 1rem;
    outline: none;
    transition: 0.3s;
    box-sizing: border-box;
  }

  input:focus {
    border-color: #3B82F6;
    background: rgba(15, 23, 42, 0.9);
    box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.2);
  }

  .action-btn {
    width: 100%;
    padding: 14px;
    border-radius: 12px;
    border: none;
    background: linear-gradient(135deg, #3B82F6, #2563EB);
    color: white;
    font-weight: bold;
    font-size: 1rem;
    cursor: pointer;
    margin-top: 10px;
    transition: 0.3s;
    box-shadow: 0 10px 15px -3px rgba(59, 130, 246, 0.4);
  }

  .action-btn:hover {
    transform: translateY(-2px);
    box-shadow: 0 20px 25px -5px rgba(59, 130, 246, 0.5);
  }

  .register-btn { background: linear-gradient(135deg, #10B981, #059669); box-shadow: 0 10px 15px -3px rgba(16, 185, 129, 0.4); }
  .register-btn:hover { box-shadow: 0 20px 25px -5px rgba(16, 185, 129, 0.5); }

  .fade-enter-active, .fade-leave-active { transition: opacity 0.3s ease; }
  .fade-enter-from, .fade-leave-to { opacity: 0; }
</style>