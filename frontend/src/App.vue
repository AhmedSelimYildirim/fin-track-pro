<template>
  <div :class="theme">
    <router-view />
  </div>
</template>

<script setup>
  import { ref, onMounted, provide } from 'vue';

  const theme = ref('dark');

  const toggleTheme = () => {
    theme.value = theme.value === 'dark' ? 'light' : 'dark';
    localStorage.setItem('theme', theme.value);
    document.documentElement.setAttribute('data-theme', theme.value);
  };

  onMounted(() => {
    const savedTheme = localStorage.getItem('theme');
    if (savedTheme) {
      theme.value = savedTheme;
    }
    document.documentElement.setAttribute('data-theme', theme.value);
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

  body {
    margin: 0;
    padding: 0;
    background-color: var(--bg-color);
    color: var(--text-color);
    font-family: 'Segoe UI', sans-serif;
    transition: background-color 0.3s, color 0.3s;
  }
</style>