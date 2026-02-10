<template>
  <div class="calendar-page">
    <div class="calendar-container">

      <div class="calendar-header">
        <button class="nav-btn" @click="changeMonth(-1)">&#10094;</button>

        <div class="header-title-wrapper">
          <h2 @click="toggleDatePicker" class="clickable-header">
            {{ currentMonthName }} {{ currentYear }}
            <span class="edit-icon">✎</span>
          </h2>

          <div v-if="showDatePicker" class="date-picker-popover">
            <select v-model="targetMonth" class="picker-select">
              <option v-for="(m, index) in allMonths" :key="index" :value="index">{{ m }}</option>
            </select>
            <input type="number" v-model="targetYear" class="picker-input" placeholder="Yıl" />
            <button @click="jumpToDate" class="jump-btn">GİT</button>
          </div>
        </div>

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
          <div class="day-header">
            <span class="day-number">{{ day.dayNumber }}</span>
          </div>

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
            <button class="add" @click="addEvent">{{ t('save') }}</button>
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

  const currentDate = ref(new Date());
  const events = ref([]);
  const showModal = ref(false);
  const showDatePicker = ref(false);
  const selectedDate = ref(null);
  const newEventText = ref('');

  const targetMonth = ref(new Date().getMonth());
  const targetYear = ref(new Date().getFullYear());

  const currentYear = computed(() => currentDate.value.getFullYear());
  const currentMonth = computed(() => currentDate.value.getMonth());

  const currentMonthName = computed(() => {
    return currentDate.value.toLocaleString(currentLang.value, { month: 'long' });
  });

  const allMonths = computed(() => {
    const months = [];
    for(let i=0; i<12; i++) {
      const d = new Date(2024, i, 1);
      months.push(d.toLocaleString(currentLang.value, { month: 'long' }));
    }
    return months;
  });

  const weekDays = computed(() => {
    const days = [];
    const d = new Date(2024, 0, 1);
    for (let i = 0; i < 7; i++) {
      days.push(d.toLocaleString(currentLang.value, { weekday: 'short' }));
      d.setDate(d.getDate() + 1);
    }
    return days;
  });

  const calendarDays = computed(() => {
    const year = currentYear.value;
    const month = currentMonth.value;

    const firstDayOfMonth = new Date(year, month, 1);
    const lastDayOfMonth = new Date(year, month + 1, 0);

    const startDay = firstDayOfMonth.getDay() === 0 ? 6 : firstDayOfMonth.getDay() - 1;
    const daysInMonth = lastDayOfMonth.getDate();

    const days = [];

    const prevMonthLastDay = new Date(year, month, 0).getDate();
    for (let i = startDay - 1; i >= 0; i--) {
      days.push({
        dayNumber: prevMonthLastDay - i,
        isCurrentMonth: false,
        date: new Date(year, month - 1, prevMonthLastDay - i)
      });
    }

    for (let i = 1; i <= daysInMonth; i++) {
      days.push({
        dayNumber: i,
        isCurrentMonth: true,
        date: new Date(year, month, i)
      });
    }

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

  const changeMonth = (step) => {
    currentDate.value = new Date(currentYear.value, currentMonth.value + step, 1);
    syncPickerValues();
  };

  const toggleDatePicker = () => {
    syncPickerValues();
    showDatePicker.value = !showDatePicker.value;
  };

  const syncPickerValues = () => {
    targetMonth.value = currentMonth.value;
    targetYear.value = currentYear.value;
  };

  const jumpToDate = () => {
    currentDate.value = new Date(targetYear.value, targetMonth.value, 1);
    showDatePicker.value = false;
  };

  const formatISODate = (date) => {
    const offset = date.getTimezoneOffset();
    const adjustedDate = new Date(date.getTime() - (offset*60*1000));
    return adjustedDate.toISOString().split('T')[0];
  };

  const formatDate = (date) => {
    if(!date) return '';
    return date.toLocaleDateString(currentLang.value, { day: 'numeric', month: 'long', year: 'numeric' });
  };

  const isToday = (date) => {
    const today = new Date();
    return date.getDate() === today.getDate() &&
            date.getMonth() === today.getMonth() &&
            date.getFullYear() === today.getFullYear();
  };

  const fetchEvents = async () => {
    try {
      const res = await api.get('/calendar/list');
      events.value = res.data || [];
    } catch (e) {
      console.error("Takvim verisi çekilemedi", e);
    }
  };

  const getEventsForDate = (date) => {
    const dateStr = formatISODate(date);
    return events.value.filter(e => e.target_date && e.target_date.split('T')[0] === dateStr);
  };

  const openAddModal = (date) => {
    selectedDate.value = date;
    newEventText.value = '';
    showModal.value = true;
  };

  const addEvent = async () => {
    if (!newEventText.value || !selectedDate.value) return;

    try {
      const targetDate = new Date(selectedDate.value);
      targetDate.setHours(12, 0, 0, 0);

      await api.post('/calendar/remind', {
        title: newEventText.value,
        target_date: targetDate.toISOString()
      });

      await fetchEvents();
      showModal.value = false;
    } catch (e) {
      alert("Hata: " + (e.response?.data?.error || "Ekleme başarısız"));
    }
  };

  const deleteEvent = async (id) => {
    if(confirm("Bu not silinsin mi?")) {
      try {
        await api.delete(`/calendar/${id}`);
        await fetchEvents();
      } catch (e) {
        alert("Silinemedi.");
      }
    }
  };

  onMounted(fetchEvents);
  watch(currentLang, () => {});
</script>

<style scoped>
  .calendar-page { padding: 30px; display: flex; justify-content: center; width: 100%; box-sizing: border-box; }
  .calendar-container { width: 100%; max-width: 1200px; background: var(--card-bg); border-radius: 20px; padding: 25px; box-shadow: 0 10px 30px rgba(0,0,0,0.3); border: 1px solid var(--border-color); }

  .calendar-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; position: relative; }
  .header-title-wrapper { position: relative; }
  .clickable-header { color: var(--text-color); margin: 0; font-size: 1.8rem; text-transform: capitalize; cursor: pointer; display: flex; align-items: center; gap: 10px; }
  .clickable-header:hover { opacity: 0.8; }
  .edit-icon { font-size: 1rem; color: var(--accent-color); opacity: 0.7; }

  .date-picker-popover { position: absolute; top: 100%; left: 50%; transform: translateX(-50%); background: var(--input-bg); border: 1px solid var(--border-color); padding: 15px; border-radius: 12px; box-shadow: 0 5px 20px rgba(0,0,0,0.5); z-index: 50; display: flex; gap: 10px; align-items: center; }
  .picker-select, .picker-input { padding: 8px; border-radius: 6px; border: 1px solid var(--border-color); background: var(--card-bg); color: var(--text-color); }
  .jump-btn { padding: 8px 15px; background: var(--success-color); color: white; border: none; border-radius: 6px; cursor: pointer; font-weight: bold; }

  .nav-btn { background: var(--hover-bg); border: none; color: var(--text-color); font-size: 1.2rem; width: 40px; height: 40px; border-radius: 50%; cursor: pointer; transition: 0.3s; }
  .nav-btn:hover { background: var(--accent-color); color: #000; }

  .weekdays-grid { display: grid; grid-template-columns: repeat(7, 1fr); margin-bottom: 10px; text-align: center; }
  .weekday { color: var(--text-muted); font-weight: bold; text-transform: uppercase; font-size: 0.9rem; letter-spacing: 1px; }

  .days-grid { display: grid; grid-template-columns: repeat(7, 1fr); gap: 8px; }
  .day-cell {
    background: var(--input-bg);
    border-radius: 8px;
    min-height: 120px;
    padding: 8px;
    cursor: pointer;
    border: 1px solid var(--border-color);
    transition: 0.2s;
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }
  .day-cell:hover { border-color: var(--accent-color); transform: translateY(-2px); }

  .day-header { display: flex; justify-content: flex-end; margin-bottom: 5px; }
  .day-number { font-weight: bold; font-size: 1rem; color: var(--text-color); }

  .dimmed { opacity: 0.3; background: transparent; }
  .dimmed .day-number { color: var(--text-muted); }

  .today { border: 2px solid var(--accent-color); background: rgba(255, 215, 0, 0.05); }
  .today .day-number { color: var(--accent-color); }

  /* --- İŞTE DÜZELTİLEN KISIM: TAŞMAYI ÖNLEYEN STİL --- */
  .events-stack {
    display: flex;
    flex-direction: column;
    gap: 3px;
    overflow-y: auto; /* Çok not varsa kaydır */
    flex: 1;
    max-height: 85px; /* Hücrenin dışına taşmasını engeller */
  }

  /* Scrollbar'ı ince yap */
  .events-stack::-webkit-scrollbar { width: 4px; }
  .events-stack::-webkit-scrollbar-thumb { background: var(--border-color); border-radius: 4px; }

  .event-pill {
    background: var(--accent-color);
    color: #000;
    font-size: 0.75rem;
    padding: 4px 6px;
    border-radius: 4px;
    font-weight: bold;

    /* Taşmayı engelle ve ... koy */
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;

    max-width: 100%;
    transition: 0.2s;
    border-left: 3px solid #000;
  }
  .event-pill:hover { opacity: 0.8; text-decoration: line-through; cursor: pointer; }

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