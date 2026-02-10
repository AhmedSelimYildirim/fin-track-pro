<template>
  <div class="login-page">
    <div class="animated-background">
      <div class="orb orb-1"></div>
      <div class="orb orb-2"></div>
      <div class="orb orb-3"></div>
    </div>

    <div class="glass-container">
      <div class="brand-section">
        <h1 class="logo-text">FinTrack Pro</h1>
        <div class="logo-accent"></div>
      </div>

      <div class="switch-wrapper">
        <div class="switch-bg" :class="{ 'right': activeTab === 'register' }"></div>
        <button class="switch-btn" :class="{ active: activeTab === 'login' }" @click="activeTab = 'login'">Giriş</button>
        <button class="switch-btn" :class="{ active: activeTab === 'register' }" @click="activeTab = 'register'">Kayıt</button>
      </div>

      <div class="form-wrapper">
        <transition name="slide-fade" mode="out-in">
          <div v-if="activeTab === 'login'" key="login" class="form-content">
            <div class="input-group">
              <span class="icon">✉️</span>
              <input v-model="loginData.email" type="email" placeholder="E-Posta" />
            </div>
            <div class="input-group">
              <span class="icon">🔒</span>
              <input v-model="loginData.password" type="password" placeholder="Şifre" />
            </div>
            <button class="submit-btn" @click="handleLogin">
              <span>GİRİŞ YAP</span>
              <div class="btn-glow"></div>
            </button>
          </div>

          <div v-else key="register" class="form-content">
            <div class="input-group">
              <span class="icon">👤</span>
              <input v-model="registerData.username" type="text" placeholder="Kullanıcı Adı" />
            </div>
            <div class="input-group">
              <span class="icon">✉️</span>
              <input v-model="registerData.email" type="email" placeholder="E-Posta" />
            </div>
            <div class="input-group">
              <span class="icon">🔑</span>
              <input v-model="registerData.password" type="password" placeholder="Şifre" />
            </div>
            <button class="submit-btn register" @click="handleRegister">
              <span>HESAP OLUŞTUR</span>
              <div class="btn-glow"></div>
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
    if (!loginData.email || !loginData.password) return alert("Lütfen bilgileri giriniz.");
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
    if (!registerData.username || !registerData.email || !registerData.password) return alert("Lütfen bilgileri giriniz.");
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
  .login-page {
    position: relative;
    width: 100vw;
    height: 100vh;
    display: flex;
    justify-content: center;
    align-items: center;
    background: #00040f;
    overflow: hidden;
    font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  }

  .animated-background {
    position: absolute;
    top: 0; left: 0; width: 100%; height: 100%;
    overflow: hidden;
    z-index: 0;
  }

  .orb {
    position: absolute;
    border-radius: 50%;
    filter: blur(100px);
    opacity: 0.5;
    animation: floatOrb 15s infinite alternate ease-in-out;
  }

  .orb-1 { width: 60vw; height: 60vw; background: #4f46e5; top: -20%; left: -10%; animation-delay: 0s; }
  .orb-2 { width: 50vw; height: 50vw; background: #ec4899; bottom: -20%; right: -10%; animation-delay: -5s; }
  .orb-3 { width: 40vw; height: 40vw; background: #10b981; top: 40%; left: 40%; animation-delay: -10s; }

  @keyframes floatOrb {
    0% { transform: translate(0, 0) scale(1); }
    100% { transform: translate(50px, 50px) scale(1.1); }
  }

  .glass-container {
    position: relative;
    z-index: 10;
    width: 380px;
    padding: 40px;
    background: rgba(255, 255, 255, 0.03);
    backdrop-filter: blur(25px);
    border-radius: 24px;
    border: 1px solid rgba(255, 255, 255, 0.08);
    box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5);
    display: flex;
    flex-direction: column;
    align-items: center;
  }

  .brand-section { text-align: center; margin-bottom: 30px; position: relative; }
  .logo-text { font-size: 2.2rem; font-weight: 800; background: linear-gradient(135deg, #fff 0%, #cbd5e1 100%); -webkit-background-clip: text; -webkit-text-fill-color: transparent; letter-spacing: -1px; margin: 0; }
  .logo-accent { width: 60px; height: 4px; background: linear-gradient(90deg, #4f46e5, #ec4899); border-radius: 2px; margin: 10px auto 0; }

  .switch-wrapper {
    position: relative;
    display: flex;
    background: rgba(0, 0, 0, 0.3);
    border-radius: 12px;
    padding: 4px;
    width: 100%;
    margin-bottom: 30px;
    border: 1px solid rgba(255, 255, 255, 0.05);
  }

  .switch-bg {
    position: absolute;
    top: 4px; left: 4px;
    width: calc(50% - 4px);
    height: calc(100% - 8px);
    background: #2dd4bf;
    border-radius: 10px;
    transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    box-shadow: 0 4px 6px -1px rgba(45, 212, 191, 0.3);
  }
  .switch-bg.right { transform: translateX(100%); background: #818cf8; box-shadow: 0 4px 6px -1px rgba(129, 140, 248, 0.3); }

  .switch-btn {
    flex: 1;
    position: relative;
    z-index: 1;
    background: none;
    border: none;
    padding: 10px;
    color: #94a3b8;
    font-weight: 600;
    cursor: pointer;
    transition: color 0.3s;
  }
  .switch-btn.active { color: #0f172a; }

  .form-wrapper { width: 100%; min-height: 250px; }
  .form-content { display: flex; flex-direction: column; gap: 16px; }

  .input-group { position: relative; width: 100%; }
  .input-group .icon { position: absolute; left: 16px; top: 50%; transform: translateY(-50%); font-size: 1.2rem; z-index: 2; opacity: 0.7; }
  .input-group input {
    width: 100%;
    padding: 16px 16px 16px 50px;
    background: rgba(0, 0, 0, 0.3);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 14px;
    color: white;
    font-size: 1rem;
    outline: none;
    transition: 0.3s;
    box-sizing: border-box;
  }
  .input-group input:focus {
    border-color: #2dd4bf;
    background: rgba(0, 0, 0, 0.5);
    box-shadow: 0 0 0 4px rgba(45, 212, 191, 0.1);
  }

  .submit-btn {
    position: relative;
    width: 100%;
    padding: 16px;
    margin-top: 10px;
    background: linear-gradient(135deg, #2dd4bf 0%, #0f766e 100%);
    border: none;
    border-radius: 14px;
    color: white;
    font-weight: 800;
    font-size: 1rem;
    cursor: pointer;
    overflow: hidden;
    transition: transform 0.2s;
    letter-spacing: 1px;
  }
  .submit-btn.register { background: linear-gradient(135deg, #818cf8 0%, #4338ca 100%); }
  .submit-btn:hover { transform: translateY(-2px); box-shadow: 0 10px 20px -5px rgba(45, 212, 191, 0.4); }
  .submit-btn.register:hover { box-shadow: 0 10px 20px -5px rgba(129, 140, 248, 0.4); }

  .btn-glow {
    position: absolute;
    top: 0; left: -100%;
    width: 100%; height: 100%;
    background: linear-gradient(90deg, transparent, rgba(255,255,255,0.2), transparent);
    animation: glow 3s infinite;
  }

  @keyframes glow { 0% { left: -100%; } 20% { left: 100%; } 100% { left: 100%; } }

  .slide-fade-enter-active { transition: all 0.3s ease-out; }
  .slide-fade-leave-active { transition: all 0.2s cubic-bezier(1, 0.5, 0.8, 1); }
  .slide-fade-enter-from, .slide-fade-leave-to { transform: translateX(10px); opacity: 0; }
</style>