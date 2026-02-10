import { reactive } from 'vue';

export const currentLang = reactive({ value: localStorage.getItem('lang') || 'tr' });

export const translations = {
    tr: {
        home: 'Anasayfa',
        calendar: 'Takvim & Notlar',
        settings: 'Ayarlar',
        logout: 'Çıkış Yap',
        portfolioSummary: 'Portföy Özeti',
        welcome: 'Hoşgeldin',
        totalValue: 'Toplam Değer',
        noData: 'Veri Yok',
        addAsset: 'Varlık Ekleyin',
        appearance: 'Görünüm',
        profileSettings: 'Profil Ayarları',
        dangerZone: 'Tehlikeli Bölge',
        update: 'Güncelle',
        deleteAccount: 'Hesabımı Sil',
        language: 'Dil Seçeneği',
        darkMode: 'Karanlık Mod',
        lightMode: 'Aydınlık Mod'
    },
    en: {
        home: 'Home',
        calendar: 'Calendar & Notes',
        settings: 'Settings',
        logout: 'Logout',
        portfolioSummary: 'Portfolio Summary',
        welcome: 'Welcome',
        totalValue: 'Total Value',
        noData: 'No Data',
        addAsset: 'Add Asset',
        appearance: 'Appearance',
        profileSettings: 'Profile Settings',
        dangerZone: 'Danger Zone',
        update: 'Update',
        deleteAccount: 'Delete Account',
        language: 'Language',
        darkMode: 'Dark Mode',
        lightMode: 'Light Mode'
    },
    de: {
        home: 'Startseite',
        calendar: 'Kalender & Notizen',
        settings: 'Einstellungen',
        logout: 'Abmelden',
        portfolioSummary: 'Portfolio Übersicht',
        welcome: 'Willkommen',
        totalValue: 'Gesamtwert',
        noData: 'Keine Daten',
        addAsset: 'Asset Hinzufügen',
        appearance: 'Aussehen',
        profileSettings: 'Profileinstellungen',
        dangerZone: 'Gefahrenzone',
        update: 'Aktualisieren',
        deleteAccount: 'Konto Löschen',
        language: 'Sprache',
        darkMode: 'Dunkelmodus',
        lightMode: 'Lichtmodus'
    },
    fr: {
        home: 'Accueil',
        calendar: 'Calendrier & Notes',
        settings: 'Paramètres',
        logout: 'Se Déconnecter',
        portfolioSummary: 'Résumé du Portfolio',
        welcome: 'Bienvenue',
        totalValue: 'Valeur Totale',
        noData: 'Pas de Données',
        addAsset: 'Ajouter un Actif',
        appearance: 'Apparence',
        profileSettings: 'Paramètres de Profil',
        dangerZone: 'Zone de Danger',
        update: 'Mettre à jour',
        deleteAccount: 'Supprimer le Compte',
        language: 'Langue',
        darkMode: 'Mode Sombre',
        lightMode: 'Mode Clair'
    }
};

export const t = (key) => {
    return translations[currentLang.value][key] || key;
};
