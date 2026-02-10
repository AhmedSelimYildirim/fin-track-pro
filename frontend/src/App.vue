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
          <img src="./assets/image_2b3783.png" class="custom-nav-icon" alt="calendar" />
          {{ t('calendar') }}
        </div>
        <div class="menu-item" :class="{ active: currentRoute.includes('/settings') }" @click="navigate('/settings')">
          <span class="nav-icon">⚙️</span> {{ t('settings') }}
        </div>
      </nav>

      <div class="logout-wrapper">
        <div class="menu-item logout" @click="logout">
          <img src="./assets/image_2b3459.png" class="custom-nav-icon" alt="logout" />
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


<style>
  :root { --bg-color: #0F172A; --card-bg: #1E293B; --text-color: #FFFFFF; --text-muted: #94A3B8; --border-color: rgba(255, 255, 255, 0.1); --input-bg: #0F172A; --sidebar-bg: #1E293B; --hover-bg: #334155; --accent-color: #FFD700; --danger-color: #EF4444; --success-color: #10B981; }
  [data-theme="light"] { --bg-color: #F1F5F9; --card-bg: #FFFFFF; --text-color: #0F172A; --text-muted: #64748B; --border-color: #E2E8F0; --input-bg: #F8FAFC; --sidebar-bg: #FFFFFF; --hover-bg: #E2E8F0; --accent-color: #F59E0B; }
  body { margin: 0; padding: 0; font-family: 'Segoe UI', sans-serif; background-color: var(--bg-color); color: var(--text-color); }
  .app-container { display: flex; min-height: 100vh; position: relative; }
  .sidebar { width: 260px; background: var(--sidebar-bg); display: flex; flex-direction: column; padding: 25px; border-right: 1px solid var(--border-color); position: fixed; height: 100vh; z-index: 1000; transition: transform 0.3s; box-sizing: border-box; }
  .main-content { flex: 1; width: 100%; transition: 0.3s; padding-left: 260px; }
  .brand-container { margin-bottom: 40px; text-align: center; }
  .brand { color: var(--accent-color); font-size: 1.6rem; font-weight: 800; letter-spacing: 1px; }
  .user-badge { margin-top: 5px; color: var(--success-color); font-weight: bold; font-size: 1.1rem; border-top: 1px solid var(--border-color); padding-top: 5px; }
  .menu-item { padding: 15px; margin-bottom: 10px; border-radius: 12px; cursor: pointer; color: var(--text-muted); display: flex; gap: 12px; align-items: center; transition: all 0.3s; font-weight: 500; }
  .menu-item:hover, .menu-item.active { background: var(--hover-bg); color: var(--text-color); transform: translateX(5px); }
  .logout-wrapper { margin-top: auto; padding-top: 20px; border-top: 1px solid rgba(255,255,255,0.1); }
  .logout { border: 2px solid var(--danger-color); border-radius: 12px; color: var(--danger-color); justify-content: center; font-weight: bold; padding: 12px; transition: 0.2s; }
  .logout:hover { background: var(--danger-color); color: white; transform: none; }
  .custom-nav-icon { width: 20px; height: 20px; object-fit: contain; }
  .mobile-menu-toggle { display: none; position: fixed; top: 15px; left: 15px; z-index: 1100; background: var(--accent-color); border: none; color: #000; font-size: 1.5rem; padding: 8px 12px; border-radius: 8px; cursor: pointer; }
  .mobile-overlay { display: none; position: fixed; inset: 0; background: rgba(0,0,0,0.5); z-index: 900; backdrop-filter: blur(3px); }
  @media (max-width: 768px) {
    .mobile-menu-toggle { display: block; }
    .sidebar { transform: translateX(-100%); width: 280px; }
    .sidebar.mobile-open { transform: translateX(0); }
    .mobile-overlay { display: block; }
    .main-content { padding-left: 0; padding-top: 60px; }
  }
</style>