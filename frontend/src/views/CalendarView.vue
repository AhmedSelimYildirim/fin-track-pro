<template>
  <div class="calendar-page">
    <div class="calendar-container">

      <div class="calendar-header">
        <button class="nav-btn" @click="changeMonth(-1)">&#10094;</button>
        <h2>{{ currentMonthName }} {{ currentYear }}</h2>
        <button class="nav-btn" @click="changeMonth(1)">&#10095;</button>
      </div>

      <div class="weekdays-grid">
        <div v-for="day in weekDays" :key="day" class="weekday">{{ day }}</div>
      </div>

      <div class="days-grid">
        <div
                v-for="(day, index) in calendarDays"
                :key="index"
                class="day-cell"
                :class="{ 'dimmed': !day.isCurrentMonth, 'today': isToday(day.date) }"
                @click="openAddModal(day.date)"
        >
          <div class="day-number">{{ day.dayNumber }}</div>

          <div class="events-stack">
            <div
                    v-for="event in getEventsForDate(day.date)"
                    :key="event.id"
                    class="event-pill"
                    @click.stop="deleteEvent(event.id)"
                    :title="event.title"
            >
              {{ event.title }}
            </div>
          </div>
        </div>
      </div>
    </div>

    <div v-if="showModal" class="modal-overlay" @click.self="showModal = false">
      <div class="modal-content">
        <div class="modal-header">
          <h3>{{ t('addAsset') }} - {{ formatDate(selectedDate) }}</h3>
          <button @click="showModal = false">✕</button>
        </div>
        <div class="modal-body">
          <input v-model="newEventText" type="text" :placeholder="t('calendar') + '...'" class="big-input" />
          <div class="actions">
            <button class="add" @click="addEvent">KAYDET</button>
          </div>
        </div>
      </div>
    </div>

  </div>
</template>

<script setup>
  import { ref, computed, onMounted, watch } from 'vue';
  import api from '../services/api';
  import { t, currentLang } from '../utils/translations';

  // --- State ---
  const currentDate = ref(new Date());
  const events = ref([]);
  const showModal = ref(false);
  const selectedDate = ref(null);
  const newEventText = ref('');

  // --- Takvim Hesaplamaları ---
  const currentYear = computed(() => currentDate.value.getFullYear());
  const currentMonth = computed(() => currentDate.value.getMonth());

  // Ay İsimleri (Dil desteğiyle)
  const currentMonthName = computed(() => {
    return currentDate.value.toLocaleString(currentLang.value, { month: 'long' });
  });

  // Gün İsimleri (Pzt, Sal...)
  const weekDays = computed(() => {
    const days = [];
    const d = new Date(2024, 0, 1); // Pazartesi ile başlayan bir tarih seçtik
    for (let i = 0; i < 7; i++) {
      days.push(d.toLocaleString(currentLang.value, { weekday: 'short' }));
      d.setDate(d.getDate() + 1);
    }
    return days;
  });

  // Takvim Izgarasını Oluşturma Mantığı
  const calendarDays = computed(() => {
    const year = currentYear.value;
    const month = currentMonth.value;

    const firstDayOfMonth = new Date(year, month, 1);
    const lastDayOfMonth = new Date(year, month + 1, 0);

    const startDay = firstDayOfMonth.getDay() === 0 ? 6 : firstDayOfMonth.getDay() - 1; // Pzt=0 yapıyoruz
    const daysInMonth = lastDayOfMonth.getDate();

    const days = [];

    // Önceki ayın günleri
    const prevMonthLastDay = new Date(year, month, 0).getDate();
    for (let i = startDay - 1; i >= 0; i--) {
      days.push({
        dayNumber: prevMonthLastDay - i,
        isCurrentMonth: false,
        date: new Date(year, month - 1, prevMonthLastDay - i)
      });
    }

    // Bu ayın günleri
    for (let i = 1; i <= daysInMonth; i++) {
      days.push({
        dayNumber: i,
        isCurrentMonth: true,
        date: new Date(year, month, i)
      });
    }

    // Sonraki ayın günleri (42 kareye tamamla - 6 satır)
    const remainingCells = 42 - days.length;
    for (let i = 1; i <= remainingCells; i++) {
      days.push({
        dayNumber: i,
        isCurrentMonth: false,
        date: new Date(year, month + 1, i)
      });
    }

    return days;
  });

  // --- Fonksiyonlar ---

  // Ay Değiştirme
  const changeMonth = (step) => {
    currentDate.value = new Date(currentYear.value, currentMonth.value + step, 1);
  };

  // Tarih Formatlama (YYYY-MM-DD)
  const formatISODate = (date) => {
    const offset = date.getTimezoneOffset();
    const adjustedDate = new Date(date.getTime() - (offset*60*1000));
    return adjustedDate.toISOString().split('T')[0];
  };

  // Görsel Formatlama (10 Şubat 2026)
  const formatDate = (date) => {
    if(!date) return '';
    return date.toLocaleDateString(currentLang.value, { day: 'numeric', month: 'long', year: 'numeric' });
  };

  // Bugün mü?
  const isToday = (date) => {
    const today = new Date();
    return date.getDate() === today.getDate() &&
            date.getMonth() === today.getMonth() &&
            date.getFullYear() === today.getFullYear();
  };

  // API'den Verileri Çek
  const fetchEvents = async () => {
    try {
      const res = await api.get('/calendar/list');
      events.value = res.data || [];
    } catch (e) {
      console.error("Takvim verisi çekilemedi", e);
    }
  };

  // Belirli bir günün notlarını filtrele
  const getEventsForDate = (date) => {
    const dateStr = formatISODate(date);
    return events.value.filter(e => e.target_date && e.target_date.split('T')[0] === dateStr);
  };

  // Modal Aç
  const openAddModal = (date) => {
    selectedDate.value = date;
    newEventText.value = '';
    showModal.value = true;
  };

  // Not Ekle
  const addEvent = async () => {
    if (!newEventText.value || !selectedDate.value) return;

    try {
      // Saat bilgisini de ekleyerek tam ISO formatı (Backend buna göre çalışıyor)
      const targetDate = new Date(selectedDate.value);
      targetDate.setHours(12, 0, 0, 0); // Varsayılan öğlen saati

      await api.post('/calendar/remind', {
        title: newEventText.value,
        target_date: targetDate.toISOString()
      });

      await fetchEvents(); // Listeyi yenile
      showModal.value = false;
    } catch (e) {
      alert("Hata: " + (e.response?.data?.error || "Ekleme başarısız"));
    }
  };

  // Not Sil
  const deleteEvent = async (id) => {
    if(confirm("Bu notu silmek istiyor musun?")) {
      try {
        await api.delete(`/calendar/${id}`);
        await fetchEvents();
      } catch (e) {
        alert("Silinemedi.");
      }
    }
  };

  // --- Lifecycle ---
  onMounted(fetchEvents);
  // Dil değişirse takvimi yeniden render etmesi için izle
  watch(currentLang, () => {
    // Tetikleyici olarak kullanılabilir
  });
</script>

<style scoped>
  .calendar-page { padding: 30px; display: flex; justify-content: center; }
  .calendar-container { width: 100%; max-width: 1000px; background: var(--card-bg); border-radius: 20px; padding: 25px; box-shadow: 0 10px 30px rgba(0,0,0,0.3); border: 1px solid var(--border-color); }

  /* Header */
  .calendar-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
  .calendar-header h2 { color: var(--text-color); margin: 0; font-size: 1.8rem; text-transform: capitalize; }
  .nav-btn { background: var(--hover-bg); border: none; color: var(--text-color); font-size: 1.2rem; width: 40px; height: 40px; border-radius: 50%; cursor: pointer; transition: 0.3s; }
  .nav-btn:hover { background: var(--accent-color); color: #000; }

  /* Haftanın Günleri */
  .weekdays-grid { display: grid; grid-template-columns: repeat(7, 1fr); margin-bottom: 10px; text-align: center; }
  .weekday { color: var(--text-muted); font-weight: bold; text-transform: uppercase; font-size: 0.9rem; letter-spacing: 1px; }

  /* Günler Izgarası */
  .days-grid { display: grid; grid-template-columns: repeat(7, 1fr); gap: 10px; }
  .day-cell {
    background: var(--input-bg);
    border-radius: 12px;
    min-height: 100px; /* Kutucuk yüksekliği */
    padding: 10px;
    cursor: pointer;
    border: 1px solid var(--border-color);
    transition: 0.2s;
    display: flex;
    flex-direction: column;
  }
  .day-cell:hover { border-color: var(--accent-color); transform: translateY(-2px); }

  .day-number { font-weight: bold; font-size: 1.1rem; color: var(--text-color); margin-bottom: 5px; }

  /* Silik Günler */
  .dimmed { opacity: 0.3; background: transparent; }
  .dimmed .day-number { color: var(--text-muted); }

  /* Bugün */
  .today { border: 2px solid var(--accent-color); background: rgba(255, 215, 0, 0.05); }
  .today .day-number { color: var(--accent-color); }

  /* Etkinlikler (Pill) */
  .events-stack { display: flex; flex-direction: column; gap: 4px; overflow-y: auto; max-height: 80px; }
  .event-pill {
    background: var(--accent-color);
    color: #000;
    font-size: 0.75rem;
    padding: 4px 8px;
    border-radius: 6px;
    font-weight: bold;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    transition: 0.2s;
  }
  .event-pill:hover { background: #fff; color: #000; opacity: 0.8; text-decoration: line-through; cursor:  no-drop;}

  /* Modal Stilleri */
  .modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.85); display: flex; justify-content: center; align-items: center; z-index: 200; backdrop-filter: blur(5px); }
  .modal-content { background: var(--card-bg); padding: 30px; border-radius: 20px; width: 350px; border: 1px solid var(--border-color); box-shadow: 0 20px 50px rgba(0,0,0,0.5); }
  .modal-header { display: flex; justify-content: space-between; margin-bottom: 20px; border-bottom: 1px solid var(--border-color); padding-bottom: 10px; }
  .modal-header h3 { color: var(--text-color); margin: 0; }
  .modal-header button { background: none; border: none; color: var(--text-color); font-size: 1.5rem; cursor: pointer; }
  .big-input { width: 100%; padding: 12px; background: var(--input-bg); border: 1px solid var(--border-color); color: var(--text-color); margin-bottom: 15px; font-size: 1.1rem; border-radius: 10px; box-sizing: border-box; }
  .actions button { width: 100%; padding: 15px; border: none; border-radius: 12px; color: white; font-weight: bold; cursor: pointer; transition: 0.2s; }
  .add { background: var(--success-color); } .add:hover { opacity: 0.9; }
</style>