<template>
  <div class="login-wrapper-fixed">
    <div class="animated-background">
      <div class="orb orb-1"></div>
      <div class="orb orb-2"></div>
      <div class="orb orb-3"></div>
    </div>

    <div class="glass-container">
      <h1 class="logo-text">FinTrack Pro</h1>
      <div class="logo-accent"></div>

      <div class="switch-wrapper">
        <div class="switch-bg" :class="{ right: activeTab === 'register' }"></div>
        <button class="switch-btn" :class="{ active: activeTab === 'login' }" @click="activeTab='login'">Giriş</button>
        <button class="switch-btn" :class="{ active: activeTab === 'register' }" @click="activeTab='register'">Kayıt</button>
      </div>

      <div v-if="activeTab==='login'">
        <input v-model="loginData.username" placeholder="Kullanıcı Adı" />
        <input v-model="loginData.password" type="password" placeholder="Şifre" />
        <button class="submit-btn" @click="handleLogin">Giriş Yap</button>
      </div>

      <div v-else>
        <input v-model="registerData.username" placeholder="Kullanıcı Adı" />
        <input v-model="registerData.email" placeholder="E-Posta" />
        <input v-model="registerData.password" type="password" placeholder="Şifre" />
        <button class="submit-btn register" @click="handleRegister">Hesap Oluştur</button>
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

  const loginData = ref({ username:'', password:'' })
  const registerData = ref({ username:'', email:'', password:'' })

  const handleLogin = async () => {
    const res = await api.post('/auth/login', loginData.value)
    localStorage.setItem('token', res.data.token)
    localStorage.setItem('username', loginData.value.username)
    router.push('/dashboard')
  }

  const handleRegister = async () => {
    await api.post('/auth/register', registerData.value)
    activeTab.value = 'login'
  }
</script>

<style scoped>
  .login-wrapper-fixed {
    position: fixed;
    inset: 0;
    background: #020617;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .glass-container {
    width: 360px;
    padding: 40px;
    background: #020617;
    border-radius: 20px;
  }

  .logo-text {
    text-align: center;
    color: white;
    font-size: 32px;
    font-weight: 800;
  }

  .logo-accent {
    width: 60px;
    height: 4px;
    background: #22c55e;
    margin: 12px auto 30px;
  }

  .switch-wrapper {
    display: flex;
    gap: 10px;
    margin-bottom: 20px;
  }

  .switch-btn {
    flex: 1;
    padding: 10px;
    background: #1e293b;
    border: none;
    color: white;
    cursor: pointer;
  }

  .switch-btn.active {
    background: #22c55e;
    color: #020617;
  }

  input {
    width: 100%;
    padding: 12px;
    margin-bottom: 12px;
    border-radius: 8px;
    border: none;
  }

  .submit-btn {
    width: 100%;
    padding: 12px;
    background: #22c55e;
    border: none;
    font-weight: 800;
    cursor: pointer;
  }

  .submit-btn.register {
    background: #3b82f6;
  }
</style>
