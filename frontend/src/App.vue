<template>
  <div :class="theme" class="app-container">
    <button v-if="!isLoginPage" class="mobile-menu-toggle" @click="toggleMobileMenu">☰</button>
    <div v-if="isMobileMenuOpen && !isLoginPage" class="mobile-overlay" @click="closeMobileMenu"></div>

    <aside v-if="!isLoginPage" class="sidebar" :class="{ 'mobile-open': isMobileMenuOpen }">
      <div class="brand-container">
        <div class="brand-title">
          <span class="brand-gradient">FinTrack</span> <span class="brand-pro">Pro</span>
        </div>

        <div class="user-info-box">
          <div class="user-avatar">{{ userInitial }}</div>
          <div class="user-details">
            <span class="user-name">{{ currentUser }}</span>
            <span class="user-role">Premium</span>
          </div>
        </div>
      </div>

      <nav class="menu">
        <div class="menu-item" :class="{ active: currentRoute.includes('/dashboard') }" @click="navigate('/dashboard')">
          <span class="nav-icon">📊</span> Portföy
        </div>
        <div class="menu-item" :class="{ active: currentRoute.includes('/calendar') }" @click="navigate('/calendar')">
          <span class="nav-icon">📅</span> Takvim
        </div>
        <div class="menu-item" :class="{ active: currentRoute.includes('/settings') }" @click="navigate('/settings')">
          <span class="nav-icon">⚙️</span> Ayarlar
        </div>
      </nav>

      <div class="logout-wrapper">
        <div class="menu-item logout-btn" @click="logout"> Çıkış Yap
        </div>
      </div>
    </aside>

    <main :class="{ 'content-shifted': !isLoginPage }" class="main-content">
      <router-view />
    </main>
  </div>
</template>

<script setup>
  import { ref, computed, onMounted, provide } from 'vue';
  import { useRoute, useRouter } from 'vue-router';

  const route = useRoute();
  const router = useRouter();
  const theme = ref('dark');
  const currentUser = ref(localStorage.getItem('username') || 'Kullanıcı');
  const isMobileMenuOpen = ref(false);

  const isLoginPage = computed(() => route.name === 'login');
  const currentRoute = computed(() => route.path);

  const userInitial = computed(() => {
    return currentUser.value ? currentUser.value.charAt(0).toUpperCase() : 'U';
  });

  const toggleTheme = () => {
    theme.value = theme.value === 'dark' ? 'light' : 'dark';
    localStorage.setItem('theme', theme.value);
    document.documentElement.setAttribute('data-theme', theme.value);
  };

  const toggleMobileMenu = () => { isMobileMenuOpen.value = !isMobileMenuOpen.value; };
  const closeMobileMenu = () => { isMobileMenuOpen.value = false; };

  const navigate = (path) => { router.push(path); closeMobileMenu(); };

  const logout = () => {
    localStorage.clear();
    currentUser.value = '';
    router.push('/login');
    closeMobileMenu();
  };

  onMounted(() => {
    const savedTheme = localStorage.getItem('theme');
    if (savedTheme) theme.value = savedTheme;
    document.documentElement.setAttribute('data-theme', theme.value);

    window.addEventListener('storage', () => {
      currentUser.value = localStorage.getItem('username') || 'Kullanıcı';
    });
  });

  provide('theme', { theme, toggleTheme });
</script>

<style>
  :root {
    --bg-color: #0F172A;
    --card-bg: #1E293B;
    --text-color: #FFFFFF;
    --text-muted: #94A3B8;
    --border-color: rgba(255, 255, 255, 0.1);
    --input-bg: #0F172A;
    --sidebar-bg: #1E293B;
    --hover-bg: #334155;
    --accent-color: #FFD700;
    --danger-color: #EF4444;
    --success-color: #10B981;
  }

  [data-theme="light"] {
    --bg-color: #F1F5F9;
    --card-bg: #FFFFFF;
    --text-color: #0F172A;
    --text-muted: #64748B;
    --border-color: #E2E8F0;
    --input-bg: #F8FAFC;
    --sidebar-bg: #FFFFFF;
    --hover-bg: #E2E8F0;
    --accent-color: #F59E0B;
  }

  body { margin: 0; padding: 0; font-family: 'Inter', 'Segoe UI', sans-serif; background-color: var(--bg-color); color: var(--text-color); overflow-x: hidden; }

  .app-container { display: flex; min-height: 100vh; position: relative; }

  .sidebar {
    width: 280px;
    background: var(--sidebar-bg);
    display: flex;
    flex-direction: column;
    padding: 30px 20px;
    border-inline-end: 1px solid var(--border-color);
    position: fixed;
    height: 100vh;
    z-index: 1000;
    transition: transform 0.3s ease-in-out;
    box-sizing: border-box;
    inset-inline-start: 0;
  }

  .main-content {
    flex: 1;
    width: 100%;
    transition: 0.3s;
    padding-inline-start: 280px;
    box-sizing: border-box;
  }

  .brand-container { margin-bottom: 40px; text-align: center; }
  .brand-title { font-size: 1.8rem; font-weight: 800; letter-spacing: -0.5px; margin-bottom: 20px; }
  .brand-gradient {
    background: linear-gradient(135deg, #fff 30%, #cbd5e1 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
  }
  [data-theme="light"] .brand-gradient {
    background: linear-gradient(135deg, #1e293b 30%, #475569 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
  }
  .brand-pro { color: #FACC15; margin-inline-start: 4px; text-shadow: 0 0 15px rgba(250, 204, 21, 0.4); }

  .user-info-box {
    display: flex;
    align-items: center;
    gap: 12px;
    background: var(--input-bg);
    padding: 12px;
    border-radius: 12px;
    border: 1px solid var(--border-color);
  }

  .user-avatar {
    width: 40px;
    height: 40px;
    background: linear-gradient(135deg, var(--accent-color), #f59e0b);
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #000;
    font-weight: 800;
    font-size: 1.2rem;
  }

  .user-details { display: flex; flex-direction: column; align-items: flex-start; overflow: hidden; }
  .user-name { font-weight: 700; font-size: 0.95rem; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; max-width: 140px; color: var(--text-color); }
  .user-role { font-size: 0.75rem; color: var(--accent-color); font-weight: 600; letter-spacing: 0.5px; text-transform: uppercase; }

  .menu { display: flex; flex-direction: column; gap: 8px; margin-top: 20px; }

  .menu-item {
    padding: 14px 16px;
    border-radius: 12px;
    cursor: pointer;
    color: var(--text-muted);
    display: flex;
    gap: 12px;
    align-items: center;
    transition: all 0.2s ease;
    font-weight: 600;
    font-size: 0.95rem;
  }

  .menu-item:hover, .menu-item.active {
    background: var(--hover-bg);
    color: var(--text-color);
    transform: translateX(5px);
  }

  .menu-item.active {
    background: rgba(250, 204, 21, 0.1);
    color: var(--accent-color);
    border-inline-start: 3px solid var(--accent-color);
  }

  .nav-icon { font-size: 1.2rem; min-width: 24px; text-align: center; }

  .logout-wrapper { margin-top: auto; padding-top: 20px; border-top: 1px solid var(--border-color); }
  .logout-btn {
    background: rgba(239, 68, 68, 0.1);
    color: var(--danger-color);
    justify-content: center;
  }
  .logout-btn:hover {
    background: var(--danger-color);
    color: white;
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(239, 68, 68, 0.3);
  }

  .mobile-menu-toggle { display: none; position: fixed; top: 15px; inset-inline-start: 15px; z-index: 1100; background: var(--accent-color); border: none; color: #000; font-size: 1.5rem; padding: 8px 12px; border-radius: 8px; cursor: pointer; box-shadow: 0 4px 10px rgba(0,0,0,0.3); }
  .mobile-overlay { display: none; position: fixed; inset: 0; background: rgba(0,0,0,0.6); z-index: 900; backdrop-filter: blur(4px); }

  @media (max-width: 768px) {
    .mobile-menu-toggle { display: block; }
    .sidebar { width: 280px; transform: translateX(-100%); transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1); }
    .sidebar.mobile-open { transform: translateX(0); }
    .mobile-overlay { display: block; }
    .main-content { padding-inline-start: 0; padding-top: 70px; }
  }
</style>