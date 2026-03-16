# Panduan Implementasi Sistem Pengelolaan Parkiran (v2 - Tailwind CSS)

Dokumen ini berisi urutan langkah-langkah, penjelasan teknologi, dan referensi untuk membangun sistem pengelolaan parkiran berbasis React, Konva.js, dan Tailwind CSS.

## 1. Urutan Pengembangan (Project Roadmap)

Untuk membangun proyek serupa, ikuti urutan berikut agar terstruktur:

1.  **Inisialisasi Proyek**:
    - Gunakan Vite atau Next.js untuk setup cepat.
    - Instal dependensi utama: `react-konva`, `konva`, `lucide-react` (ikon), `framer-motion` (animasi), `date-fns` (manipulasi waktu).
    - **Penting**: Instal dan konfigurasi **Tailwind CSS**. Ini akan sangat memudahkan styling tanpa harus membuat file CSS yang besar manual.
2.  **Perancangan Data (Data Schema)**:
    - Tentukan struktur data untuk `spots` (koordinat dan ID) dan `bookings` (siapa, di mana, berapa lama).
3.  **Pembuatan Denah Interaktif (Canvas)**:
    - Gunakan Konva.js untuk menggambar denah. Mengapa Konva? Karena lebih efisien daripada ribuan elemen DOM jika denahnya kompleks.
    - Gunakan konteks React untuk meneruskan status booking ke elemen canvas.
4.  **UI Layout dengan Tailwind CSS**:
    - Bangun sidebar, header, dan grid dashboard menggunakan utility classes Tailwind.
    - Manfaatkan plugin seperti `@tailwindcss/forms` jika perlu form yang lebih kompleks.
5.  **Sistem Persistensi**:
    - Implementasikan custom hook `useLocalStorage` agar data tidak hilang saat halaman di-refresh.
6.  **Fitur Pemesanan (Booking Logic)**:
    - Buat formulir modal dengan animasi transisi (menggunakan Framer Motion).
    - Validasi input dan simpan data ke state utama.
7.  **Real-time Monitoring & Timer**:
    - Implementasikan logika countdown menggunakan `setInterval` di dalam `useEffect`.
    - Gunakan library `date-fns` untuk menghitung sisa waktu dan status overtime.
8.  **Responsivitas**:
    - Gunakan prefix `sm:`, `md:`, `lg:`, `xl:` dari Tailwind untuk menyesuaikan layout di berbagai ukuran layar.

## 2. Mengapa Menggunakan Framework CSS (Tailwind)?

- **Kecepatan**: Tidak perlu berpindah antar file `.css` dan `.jsx`.
- **Konsistensi**: Menggunakan sistem spacing dan warna yang terstandarisasi.
- **Maintainability**: Komponen menjadi lebih modular karena gaya terikat langsung pada struktur HTML.
- **Responsivitas Mudah**: Jauh lebih intuitif daripada menulis media queries manual.

## 3. Library Komponen yang Direkomendasikan

Jika proyek ini ingin dikembangkan lebih lanjut, pertimbangkan menggunakan:

- **Shadcn UI**: Untuk komponen UI (Button, Dialog, Badge) yang sangat fleksibel.
- **Headless UI / Radix UI**: Jika butuh komponen interaktif yang _accessible_ (seperti Modal, Dropdown) tanpa opini styling.
- **Framer Motion**: Untuk animasi mikro yang membuat aplikasi terasa premium.

## 4. Cara Menjalankan Proyek

1. `npm install`
2. `npm run dev`
3. Akses `http://localhost:5173`

---

_Dibuat oleh Antigravity (AI Coding Assistant) - Refactored for Tailwind CSS._
