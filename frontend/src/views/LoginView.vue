<template>
  <div class="login-page">
    <div class="glow-bg glow-1"></div>
    <div class="glow-bg glow-2"></div>

    <div class="glass-card">
      <h1 class="brand-title">FinTrack <span class="highlight">Pro</span></h1>
      <p class="brand-subtitle">Finansal Özgürlüğe Hoşgeldiniz</p>

      <div class="tab-container">
        <div class="tab-slider" :class="{ 'slide-right': activeTab === 'register' }"></div>
        <button class="tab-btn" :class="{ active: activeTab === 'login' }" @click="activeTab='login'">Giriş Yap</button>
        <button class="tab-btn" :class="{ active: activeTab === 'register' }" @click="activeTab='register'">Kayıt Ol</button>
      </div>

      <div class="form-wrapper">
        <transition name="fade" mode="out-in">
          <div v-if="activeTab==='login'" key="login" class="form-content">
            <div class="input-box">
              <input v-model="loginData.username" placeholder="Kullanıcı Adı veya E-Posta" />
            </div>
            <div class="input-box">
              <input v-model="loginData.password" type="password" placeholder="Şifre" />
            </div>
            <button class="action-btn" @click="handleLogin">
              <span>Giriş Yap</span>
            </button>
          </div>

          <div v-else key="register" class="form-content">
            <div class="input-box">
              <input v-model="registerData.username" placeholder="Kullanıcı Adı" />
            </div>
            <div class="input-box">
              <input v-model="registerData.email" placeholder="E-Posta Adresi" />
            </div>
            <div class="input-box">
              <input v-model="registerData.password" type="password" placeholder="Şifre Belirle" />
            </div>
            <button class="action-btn register-btn" @click="handleRegister">
              <span>Hesap Oluştur</span>
            </button>
          </div>
        </transition>
      </div>
    </div>
  </div>
</template>

<script setup>
  import { ref } from 'vue'
  import { useRouter } from 'vue-router'
  import api from '../services/api'

  const router = useRouter()
  const activeTab = ref('login')

  const loginData = ref({ username: '', password: '' })
  const registerData = ref({ username: '', email: '', password: '' })

  const handleLogin = async () => {
    try {
      const res = await api.post('/auth/login', loginData.value)

      // Token Kaydı
      localStorage.setItem('token', res.data.token)

      // Backend'den gelen kullanıcı verisini al
      const userData = res.data.user || {}

      // Username ve Email'i ayrı ayrı kaydet (Settings için kritik)
      if (userData.username) {
        localStorage.setItem('username', userData.username)
      } else {
        // Eğer backend username dönmezse inputtan al
        localStorage.setItem('username', loginData.value.username)
      }

      if (userData.email) {
        localStorage.setItem('email', userData.email)
      } else {
        // Email yoksa boş string bas, null olmasın
        localStorage.setItem('email', '')
      }

      router.push('/dashboard')
    } catch (e) {
      alert('Giriş başarısız: ' + (e.response?.data?.message || e.message))
    }
  }

  const handleRegister = async () => {
    try {
      await api.post('/auth/register', registerData.value)
      alert('Kayıt başarılı! Lütfen giriş yapın.')
      activeTab.value = 'login'
    } catch (e) {
      alert('Kayıt hatası: ' + e.message)
    }
  }
</script>

<style scoped>
  /* Fixed Layout */
  .login-page {
    position: fixed;
    inset: 0;
    width: 100vw;
    height: 100vh;
    background: #020617;
    display: flex;
    justify-content: center;
    align-items: center;
    overflow: hidden;
    z-index: 1000;
  }

  /* Glow Effects */
  .glow-bg {
    position: absolute;
    border-radius: 50%;
    filter: blur(80px);
    opacity: 0.4;
  }
  .glow-1 { width: 300px; height: 300px; background: #22c55e; top: -50px; left: -50px; }
  .glow-2 { width: 400px; height: 400px; background: #3b82f6; bottom: -100px; right: -100px; }

  /* Glass Card */
  .glass-card {
    position: relative;
    width: 400px;
    background: rgba(30, 41, 59, 0.4);
    backdrop-filter: blur(24px);
    -webkit-backdrop-filter: blur(24px);
    border: 1px solid rgba(255, 255, 255, 0.08);
    border-radius: 24px;
    padding: 40px;
    box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5);
    z-index: 10;
  }

  .brand-title {
    text-align: center;
    color: white;
    font-size: 2rem;
    font-weight: 800;
    margin: 0;
    letter-spacing: -1px;
  }

  .highlight { color: #facc15; }
  .brand-subtitle { text-align: center; color: #94a3b8; font-size: 0.9rem; margin-top: 5px; margin-bottom: 30px; }

  /* Tab Container */
  .tab-container {
    background: rgba(15, 23, 42, 0.6);
    padding: 4px;
    border-radius: 12px;
    display: flex;
    position: relative;
    margin-bottom: 25px;
  }

  .tab-slider {
    position: absolute;
    width: 50%;
    height: calc(100% - 8px);
    background: #334155;
    border-radius: 8px;
    top: 4px;
    left: 4px;
    transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    z-index: 1;
  }

  .tab-slider.slide-right { transform: translateX(96%); }

  .tab-btn {
    flex: 1;
    background: none;
    border: none;
    padding: 12px;
    color: #94a3b8;
    font-weight: 600;
    cursor: pointer;
    z-index: 2;
    transition: color 0.3s;
    font-size: 0.95rem;
  }

  .tab-btn.active { color: white; }

  /* Inputs */
  .input-box { margin-bottom: 15px; }

  input {
    width: 100%;
    padding: 16px;
    background: rgba(2, 6, 23, 0.4);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 12px;
    color: white;
    font-size: 1rem;
    outline: none;
    transition: all 0.3s;
  }

  input:focus {
    border-color: #22c55e;
    box-shadow: 0 0 0 4px rgba(34, 197, 94, 0.1);
    background: rgba(2, 6, 23, 0.8);
  }

  /* Buttons */
  .action-btn {
    width: 100%;
    padding: 16px;
    background: linear-gradient(135deg, #22c55e 0%, #15803d 100%);
    border: none;
    border-radius: 12px;
    color: white;
    font-weight: 700;
    font-size: 1rem;
    cursor: pointer;
    margin-top: 10px;
    transition: transform 0.2s, box-shadow 0.2s;
  }

  .action-btn:hover { transform: translateY(-2px); box-shadow: 0 10px 20px -5px rgba(34, 197, 94, 0.4); }
  .register-btn { background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%); }
  .register-btn:hover { box-shadow: 0 10px 20px -5px rgba(59, 130, 246, 0.4); }

  /* Transitions */
  .fade-enter-active, .fade-leave-active { transition: opacity 0.3s ease; }
  .fade-enter-from, .fade-leave-to { opacity: 0; }
</style>