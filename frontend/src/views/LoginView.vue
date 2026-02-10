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

<template>
  <div :class="theme" class="app-container">
    <button v-if="!isLoginPage" class="mobile-menu-toggle" @click="toggleMobileMenu">☰</button>
    <div v-if="isMobileMenuOpen && !isLoginPage" class="mobile-overlay" @click="closeMobileMenu"></div>

    <aside v-if="!isLoginPage" class="sidebar" :class="{ 'mobile-open': isMobileMenuOpen }">
      <div class="brand-container">
        <div class="brand">FinTrack Pro</div>
        <div class="user-badge">{{ currentUser }}</div>
      </div>

      <nav class="menu">
        <div class="menu-item" :class="{ active: currentRoute.includes('/dashboard') }" @click="navigate('/dashboard')">
          <span class="nav-icon">📊</span> {{ t('home') }}
        </div>
        <div class="menu-item" :class="{ active: currentRoute.includes('/calendar') }" @click="navigate('/calendar')">
          <img :src="calendarIcon" class="custom-nav-icon" />
          {{ t('calendar') }}
        </div>
        <div class="menu-item" :class="{ active: currentRoute.includes('/settings') }" @click="navigate('/settings')">
          <span class="nav-icon">⚙️</span> {{ t('settings') }}
        </div>
      </nav>

      <div class="logout-wrapper">
        <div class="menu-item logout" @click="logout">
          <img :src="logoutIcon" class="custom-nav-icon" />
          {{ t('logout') }}
        </div>
      </div>
    </aside>

    <main :class="{ 'content-shifted': !isLoginPage }" class="main-content">
      <router-view />
    </main>
  </div>
</template>

<script setup>
  import { ref, computed, onMounted, provide } from 'vue'
  import { useRoute, useRouter } from 'vue-router'
  import { t } from './utils/translations'
  import calendarIcon from './assets/calendar-icon.png'
  import logoutIcon from './assets/logout-icon.png'

  const route = useRoute()
  const router = useRouter()

  const theme = ref('dark')
  const currentUser = ref(localStorage.getItem('username') || '')
  const isMobileMenuOpen = ref(false)

  const isLoginPage = computed(() => route.name === 'login')
  const currentRoute = computed(() => route.path)

  const toggleTheme = () => {
    theme.value = theme.value === 'dark' ? 'light' : 'dark'
    localStorage.setItem('theme', theme.value)
    document.documentElement.setAttribute('data-theme', theme.value)
  }

  const toggleMobileMenu = () => { isMobileMenuOpen.value = !isMobileMenuOpen.value }
  const closeMobileMenu = () => { isMobileMenuOpen.value = false }

  const navigate = (path) => {
    router.push(path)
    closeMobileMenu()
  }

  const logout = () => {
    localStorage.clear()
    router.push('/login')
    closeMobileMenu()
  }

  onMounted(() => {
    const savedTheme = localStorage.getItem('theme')
    if (savedTheme) theme.value = savedTheme
    document.documentElement.setAttribute('data-theme', theme.value)

    window.addEventListener('storage', () => {
      currentUser.value = localStorage.getItem('username') || ''
    })
  })

  provide('theme', { theme, toggleTheme })
</script>


<style scoped>
  .login-wrapper-fixed { position: fixed; top: 0; left: 0; width: 100vw; height: 100vh; display: flex; justify-content: center; align-items: center; background: #00040f; z-index: 9999; }
  .animated-background { position: absolute; inset: 0; overflow: hidden; z-index: 0; }
  .orb { position: absolute; border-radius: 50%; filter: blur(100px); opacity: 0.5; animation: floatOrb 15s infinite alternate ease-in-out; }
  .orb-1 { width: 60vw; height: 60vw; background: #4f46e5; top: -20%; left: -10%; }
  .orb-2 { width: 50vw; height: 50vw; background: #ec4899; bottom: -20%; right: -10%; }
  .orb-3 { width: 40vw; height: 40vw; background: #10b981; top: 40%; left: 40%; }
  @keyframes floatOrb { 0% { transform: translate(0, 0) scale(1); } 100% { transform: translate(50px, 50px) scale(1.1); } }
  .glass-container { position: relative; z-index: 10; width: 380px; padding: 40px; background: rgba(255, 255, 255, 0.03); backdrop-filter: blur(25px); border-radius: 24px; border: 1px solid rgba(255, 255, 255, 0.08); box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5); }
  .logo-text { font-size: 2.2rem; font-weight: 800; background: linear-gradient(135deg, #fff 0%, #cbd5e1 100%); -webkit-background-clip: text; -webkit-text-fill-color: transparent; letter-spacing: -1px; text-align: center; }
  .logo-accent { width: 60px; height: 4px; background: linear-gradient(90deg, #4f46e5, #ec4899); border-radius: 2px; margin: 10px auto 30px; }
  .switch-wrapper { position: relative; display: flex; background: rgba(0, 0, 0, 0.3); border-radius: 12px; padding: 4px; width: 100%; margin-bottom: 30px; border: 1px solid rgba(255, 255, 255, 0.05); }
  .switch-bg { position: absolute; top: 4px; left: 4px; width: calc(50% - 4px); height: calc(100% - 8px); background: #2dd4bf; border-radius: 10px; transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1); }
  .switch-bg.right { transform: translateX(100%); background: #818cf8; }
  .switch-btn { flex: 1; position: relative; z-index: 1; background: none; border: none; padding: 10px; color: #94a3b8; font-weight: 600; cursor: pointer; }
  .switch-btn.active { color: #0f172a; }
  .input-group { position: relative; margin-bottom: 16px; }
  .input-group .icon { position: absolute; left: 16px; top: 50%; transform: translateY(-50%); font-size: 1.2rem; z-index: 2; }
  .input-group input { width: 100%; padding: 16px 16px 16px 50px; background: rgba(0, 0, 0, 0.3); border: 1px solid rgba(255, 255, 255, 0.1); border-radius: 14px; color: white; outline: none; box-sizing: border-box; }
  .submit-btn { width: 100%; padding: 16px; background: linear-gradient(135deg, #2dd4bf 0%, #0f766e 100%); border: none; border-radius: 14px; color: white; font-weight: 800; cursor: pointer; position: relative; overflow: hidden; display: flex; justify-content: center; align-items: center; }
  .submit-btn:disabled { opacity: 0.7; cursor: not-allowed; }
  .submit-btn.register { background: linear-gradient(135deg, #818cf8 0%, #4338ca 100%); }
  .spinner { width: 20px; height: 20px; border: 3px solid rgba(255,255,255,0.3); border-top-color: white; border-radius: 50%; animation: spin 1s linear infinite; }
  @keyframes spin { to { transform: rotate(360deg); } }
  .slide-fade-enter-active, .slide-fade-leave-active { transition: all 0.3s ease; }
  .slide-fade-enter-from, .slide-fade-leave-to { opacity: 0; transform: translateX(20px); }
</style>