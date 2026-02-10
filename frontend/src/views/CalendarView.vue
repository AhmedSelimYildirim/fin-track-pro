<template>
  <div class="calendar-page">
    <aside class="sidebar">
      <div class="brand">FinTrack Pro</div>

      <nav class="menu">
        <div class="menu-item" @click="router.push('/dashboard')">
          <span>📊</span> Dashboard
        </div>
        <div class="menu-item active">
          <span>📅</span> Takvim & Notlar
        </div>
      </nav>

      <div class="logout-wrapper">
        <div class="menu-item logout" @click="logout">
          <span>🚪</span> Çıkış Yap
        </div>
      </div>
    </aside>

    <main class="content">
      <header class="header">
        <h1>Finansal Takvim & Hatırlatıcılar</h1>
        <p>Ödemelerini ve hedeflerini buradan takip et.</p>
      </header>

      <div class="add-note-card">
        <div class="input-group">
          <input v-model="newNote.title" type="text" placeholder="Örn: Kredi Kartı Ödemesi" />
          <input v-model="newNote.date" type="datetime-local" />
          <button @click="addNote" :disabled="loading">
            {{ loading ? 'Ekleniyor...' : 'Not Ekle (+)' }}
          </button>
        </div>
      </div>

      <div class="notes-grid">
        <div v-if="notes.length === 0" class="empty-state">
          Henüz hiç notun yok.
        </div>

        <div v-for="note in notes" :key="note.id" class="note-card">
          <div class="note-header">
            <span class="note-date">{{ formatDate(note.target_date) }}</span>
            <button class="delete-btn" @click="deleteNote(note.id)">🗑️</button>
          </div>
          <div class="note-body">
            {{ note.title }}
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import api from '../services/api';

const router = useRouter();
const notes = ref([]);
const loading = ref(false);
const newNote = ref({ title: '', date: '' });

// Notları Getir (Backend: GET /calendar/list)
const fetchNotes = async () => {
  try {
    const res = await api.get('/calendar/list');
    // Backend'den gelen veriyi tarihe göre sırala
    notes.value = res.data.sort((a, b) => new Date(a.target_date) - new Date(b.target_date));
  } catch (e) {
    console.error("Notlar çekilemedi", e);
  }
};

// Not Ekle (Backend: POST /calendar/remind)
const addNote = async () => {
  if (!newNote.value.title || !newNote.value.date) return alert("Başlık ve Tarih giriniz.");

  loading.value = true;
  try {
    // Backend ISO formatı istiyor: 2027-02-16T10:00:00Z
    const isoDate = new Date(newNote.value.date).toISOString();

    await api.post('/calendar/remind', {
      title: newNote.value.title,
      target_date: isoDate
    });

    newNote.value = { title: '', date: '' };
    fetchNotes(); // Listeyi yenile
  } catch (e) {
    alert("Hata: " + e.message);
  } finally {
    loading.value = false;
  }
};

// Not Sil (Backend: DELETE /calendar/:id)
const deleteNote = async (id) => {
  if(!confirm("Bu notu silmek istiyor musun?")) return;
  try {
    await api.delete(`/calendar/${id}`);
    fetchNotes();
  } catch (e) {
    alert("Silinemedi.");
  }
};

// Tarih Formatlayıcı (Görüntü için)
const formatDate = (isoStr) => {
  const date = new Date(isoStr);
  return date.toLocaleDateString('tr-TR', { day: 'numeric', month: 'long', hour: '2-digit', minute:'2-digit' });
};

const logout = () => {
  localStorage.clear();
  router.push('/login');
};

onMounted(() => {
  fetchNotes();
});
</script>

<style scoped>
.calendar-page { display: flex; min-height: 100vh; background: #0F172A; color: white; font-family: 'Segoe UI', sans-serif; }

/* Sidebar */
.sidebar { width: 250px; background: #1E293B; display: flex; flex-direction: column; padding: 20px; border-right: 1px solid rgba(255,255,255,0.05); }
.brand { color: #FFD700; font-size: 1.5rem; font-weight: bold; margin-bottom: 40px; }
.menu-item { padding: 15px; margin-bottom: 10px; border-radius: 10px; cursor: pointer; color: #94A3B8; display: flex; gap: 10px; align-items: center; transition: 0.2s; }
.menu-item:hover, .menu-item.active { background: #334155; color: white; }
.logout-wrapper { margin-top: auto; }
.logout { color: #EF4444; } .logout:hover { background: rgba(239, 68, 68, 0.1); }

/* Content */
.content { flex: 1; padding: 40px; }
.header h1 { margin: 0; font-size: 2rem; }
.header p { color: #94A3B8; margin-top: 5px; }

/* Form */
.add-note-card { background: #1E293B; padding: 20px; border-radius: 15px; margin: 30px 0; border: 1px solid rgba(255,255,255,0.05); }
.input-group { display: flex; gap: 15px; flex-wrap: wrap; }
.input-group input { flex: 1; padding: 12px; background: #0F172A; border: 1px solid #334155; color: white; border-radius: 8px; min-width: 200px; }
.input-group button { padding: 12px 25px; background: #10B981; border: none; border-radius: 8px; color: white; font-weight: bold; cursor: pointer; }

/* Liste */
.notes-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(250px, 1fr)); gap: 20px; }
.note-card { background: #1E293B; padding: 20px; border-radius: 12px; border-left: 4px solid #3B82F6; position: relative; }
.note-header { display: flex; justify-content: space-between; margin-bottom: 10px; font-size: 0.85rem; color: #94A3B8; }
.delete-btn { background: none; border: none; cursor: pointer; opacity: 0.5; transition: 0.2s; }
.delete-btn:hover { opacity: 1; transform: scale(1.1); }
.note-body { font-size: 1.1rem; font-weight: 500; }
.empty-state { color: #94A3B8; font-style: italic; }
</style>