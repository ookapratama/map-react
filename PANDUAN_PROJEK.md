# Panduan Lengkap Membangun Sistem Pengelolaan Parkiran (Vite + React + Tailwind CSS + Konva.js)

Panduan ini dirancang untuk membantu Anda membangun ulang proyek ini dari nol secara sistematis dan mudah dipahami.

## Daftar Isi

1. [Setup Proyek & Instalasi](#1-setup-proyek--instalasi)
2. [Konfigurasi Styling (Tailwind)](#2-konfigurasi-styling-tailwind)
3. [Persiapan Data & Utilitas](#3-persiapan-data--utilitas)
4. [Membangun Komponen Denah (Canvas)](#4-membangun-komponen-denah-canvas)
5. [Membangun Komponen UI (Form & Details)](#5-membangun-komponen-ui-form--details)
6. [Implementasi Logika Utama (App.jsx)](#6-implementasi-logika-utama-appjsx)

---

## 1. Setup Proyek & Instalasi

Pertama, buat proyek React baru menggunakan Vite:

```bash
# Buat proyek baru
npm create vite@latest parking-system -- --template react

# Masuk ke folder proyek
cd parking-system

# Instal dependensi dasar
npm install

# Instal library tambahan yang diperlukan
npm install konva react-konva lucide-react framer-motion date-fns
```

## 2. Konfigurasi Styling (Tailwind)

Instal Tailwind CSS versi terbaru (v4) atau versi v3 sesuai preferensi. Berikut adalah konfigurasi yang kita gunakan di proyek ini:

### 2.1 Konfigurasi `tailwind.config.js`

Tambahkan palette warna custom agar tampilan terlihat premium:

```javascript
/** @type {import('tailwindcss').Config} */
export default {
  content: ["./index.html", "./src/**/*.{js,ts,jsx,tsx}"],
  theme: {
    extend: {
      colors: {
        primary: {
          DEFAULT: "#4f46e5", // Indigo
          hover: "#4338ca",
          light: "#eef2ff",
        },
        success: "#10b981", // Emerald
        danger: "#ef4444", // Red
      },
      fontFamily: {
        sans: ["Inter", "sans-serif"],
        heading: ["Outfit", "sans-serif"],
      },
      borderRadius: {
        xl: "12px",
      },
    },
  },
  plugins: [],
};
```

### 2.2 Setup `src/index.css`

Tambahkan font dari Google Fonts (Inter & Outfit) dan layer Tailwind:

```css
@import url("https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&family=Outfit:wght@700;800;900&display=swap");

@tailwind base;
@tailwind components;
@tailwind utilities;

@layer base {
  body {
    @apply text-slate-900 antialiased;
  }
  h1,
  h2,
  h3,
  h4 {
    @apply font-heading font-black tracking-tight;
  }
}
```

## 3. Persiapan Data & Utilitas

### 3.1 Data Titik Parkir (`src/constants/spots.js`)

Kita buat data titik parkir secara otomatis menggunakan loop:

```javascript
export const INITIAL_SPOTS = Array.from({ length: 20 }, (_, i) => {
  const row = Math.floor(i / 10);
  const col = i % 10;
  return {
    id: `spot-${i + 1}`,
    label: `A-${i + 1}`,
    x: 50 + col * 70, // Jarak horizontal
    y: 50 + row * 150, // Jarak vertikal (antar baris)
    width: 60,
    height: 100,
  };
});
```

### 3.2 Hook Local Storage (`src/hooks/useLocalStorage.js`)

Agar data tidak hilang saat refresh:

```javascript
import { useState, useEffect } from "react";

function useLocalStorage(key, initialValue) {
  const [value, setValue] = useState(() => {
    const jsonValue = localStorage.getItem(key);
    if (jsonValue != null) return JSON.parse(jsonValue);
    return initialValue;
  });

  useEffect(() => {
    localStorage.setItem(key, JSON.stringify(value));
  }, [key, value]);

  return [value, setValue];
}

export default useLocalStorage;
```

## 4. Membangun Komponen Denah (Canvas)

Gunakan `react-konva` untuk menggambar denah. Komponen ini (`src/components/ParkingMap.jsx`) akan menerima data `spots` dan `bookings`.

**Poin Penting:**

- Gunakan `Stage` dan `Layer` sebagai wadah utama.
- Gunakan `Rect` untuk kotak parkir dan mobil.
- Gunakan `Group` untuk mengelompokkan elemen per titik parkir agar mudah diklik.

## 5. Membangun Komponen UI (Form & Details)

### 5.1 Modal Form (`src/components/BookingForm.jsx`)

Gunakan `framer-motion` untuk animasi modal yang halus. Form ini digunakan untuk menginput Nama, Plat Nomor, dan Durasi.

### 5.2 Rincian Pemesanan (`src/components/BookingDetails.jsx`)

Gunakan `lucide-react` untuk ikon-ikon seperti jam, user, dan plat nomor. Tambahkan fitur hitung mundur (countdown) menggunakan `date-fns`.

## 6. Implementasi Logika Utama (App.jsx)

Di `App.jsx`, kita menyatukan semua logika:

1. **State Management**: Gunakan `useLocalStorage` untuk data booking.
2. **Kalkulasi**: Hitung jumlah parkiran yang terisi (`Terisi`) dan sisa (`Tersedia`).
3. **Event Handler**:
   - `handleSpotClick`: Cek apakah titik tersebut sudah dipesan atau belum.
   - `handleBook`: Menambahkan booking baru ke array.
   - `handleEndSession`: Mengubah status booking menjadi 'completed'.

---

## 🎨 Tips Agar Desain Terlihat Premium

1. **Gunakan Gradients & Soft Shadows**: Hindari warna solid yang terlalu tajam. Gunakan shadow halus (`shadow-sm` atau `shadow-xl`).
2. **Micro-interactions**: Gunakan animasi exit/enter pada modal dan notifikasi agar aplikasi terasa "hidup".
3. **Typography yang Kuat**: Gunakan font sans-serif yang modern (seperti Inter) untuk teks biasa, dan font yang lebih berkarakter (seperti Outfit) untuk heading.
4. **Empty State**: Selalu buat tampilan "Pilih Tempat Parkir" jika belum ada area yang dipilih, jangan biarkan sidebar kosong.

## 🚀 Cara Menjalankan

1. Terakhir, jalankan `npm run dev`.
2. Buka browser pada alamat yang tertera (biasanya `localhost:5173`).

---

_Dibuat untuk mempermudah Anda membangun ulang sistem parkir yang interaktif dan modern._
