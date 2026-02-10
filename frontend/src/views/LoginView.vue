<template>
  <div class="login-wrapper-fixed">
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
        <div class="switch-bg" :class="{ right: activeTab === 'register' }"></div>
        <button class="switch-btn" :class="{ active: activeTab === 'login' }" @click="activeTab = 'login'">Giriş</button>
        <button class="switch-btn" :class="{ active: activeTab === 'register' }" @click="activeTab = 'register'">Kayıt</button>
      </div>

      <div class="form-wrapper">
        <transition name="slide-fade" mode="out-in">
          <div v-if="activeTab === 'login'" key="login" class="form-content">
            <div class="input-group">
              <span class="icon">👤</span>
              <input v-model="loginData.username" type="text" placeholder="Kullanıcı Adı" />
            </div>
            <div class="input-group">
              <span class="icon">🔒</span>
              <input v-model="loginData.password" type="password" placeholder="Şifre" />
            </div>
            <button class="submit-btn" :disabled="isLoading" @click="handleLogin">
              <span v-if="!isLoading">GİRİŞ YAP</span>
              <div v-else class="spinner"></div>
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
            <button class="submit-btn register" :disabled="isLoading" @click="handleRegister">
              <span v-if="!isLoading">HESAP OLUŞTUR</span>
              <div v-else class="spinner"></div>
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
  const isLoading = ref(false)

  const loginData = ref({
    username: '',
    password: ''
  })

  const registerData = ref({
    username: '',
    email: '',
    password: ''
  })

  const handleLogin = async () => {
    isLoading.value = true
    const res = await api.post('/auth/login', loginData.value)
    localStorage.setItem('token', res.data.token)
    localStorage.setItem('username', loginData.value.username)
    router.push('/dashboard')
    isLoading.value = false
  }

  const handleRegister = async () => {
    isLoading.value = true
    await api.post('/auth/register', registerData.value)
    activeTab.value = 'login'
    isLoading.value = false
  }
</script>

<style scoped>
  .login-wrapper-fixed { position: fixed; inset: 0; display: flex; justify-content: center; align-items: center; background: #00040f; }
  .animated-background { position: absolute; inset: 0; overflow: hidden; }
  .orb { position: absolute; border-radius: 50%; filter: blur(100px); opacity: 0.5; animation: floatOrb 15s infinite alternate ease-in-out; }
  .orb-1 { width: 60vw; height: 60vw; background: #4f46e5; top: -20%; left: -10%; }
  .orb-2 { width: 50vw; height: 50vw; background: #ec4899; bottom: -20%; right: -10%; }
  .orb-3 { width: 40vw; height: 40vw; background: #10b981; top: 40%; left: 40%; }
  @keyframes floatOrb { to { transform: translate(50px, 50px) scale(1.1); } }
  .glass-container { width: 380px; padding: 40px; background: rgba(255,255,255,0.03); backdrop-filter: blur(25px); border-radius: 24px; border: 1px solid rgba(255,255,255,0.08); }
  .logo-text { text-align: center; font-size: 2.2rem; font-weight: 800; color: white; }
  .logo-accent { width: 60px; height: 4px; background: linear-gradient(90deg,#4f46e5,#ec4899); margin: 12px auto 30px; border-radius: 2px; }
  .switch-wrapper { position: relative; display: flex; background: rgba(0,0,0,0.3); border-radius: 12px; padding: 4px; margin-bottom: 30px; }
  .switch-bg { position: absolute; inset: 4px auto 4px 4px; width: calc(50% - 4px); background: #2dd4bf; border-radius: 10px; transition: transform .3s; }
  .switch-bg.right { transform: translateX(100%); background: #818cf8; }
  .switch-btn { flex: 1; background: none; border: none; color: #94a3b8; font-weight: 600; cursor: pointer; z-index: 1; }
  .switch-btn.active { color: #0f172a; }
  .input-group { position: relative; margin-bottom: 16px; }
  .input-group .icon { position: absolute; left: 16px; top: 50%; transform: translateY(-50%); }
  .input-group input { width: 100%; padding: 16px 16px 16px 50px; border-radius: 14px; border: 1px solid rgba(255,255,255,0.1); background: rgba(0,0,0,0.3); color: white; }
  .submit-btn { width: 100%; padding: 16px; border-radius: 14px; border: none; font-weight: 800; background: linear-gradient(135deg,#2dd4bf,#0f766e); cursor: pointer; }
  .submit-btn.register { background: linear-gradient(135deg,#818cf8,#4338ca); }
  .spinner { width: 20px; height: 20px; border: 3px solid rgba(255,255,255,0.3); border-top-color: white; border-radius: 50%; animation: spin 1s linear infinite; }
  @keyframes spin { to { transform: rotate(360deg); } }
  .slide-fade-enter-active,.slide-fade-leave-active{transition:.3s}
  .slide-fade-enter-from,.slide-fade-leave-to{opacity:0;transform:translateX(20px)}
</style>
