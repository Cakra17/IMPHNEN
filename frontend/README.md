<a id="english_ver"></a>
# IMPHNEN

**Ingin Menjadi Pengusaha Handal Namun Enggan Ngebuku** - *Want to Be a Skilled Entrepreneur but Too Lazy to Keep Books*

A web application that helps small and medium enterprises (MSMEs/UMKM) easily record receipts and customer orders. Business owners simply need to take photos of receipts, upload them to the app, and let AI scan and automatically add them to the database without manual typing. It also integrates with WhatsApp Business so AI can detect transactions and automatically enter them into the data. This gives UMKM owners more time for what matters: **Improving service quality and developing their business**.

---

[Go to Indonesian Version](#indonesian_ver)

---

## 🌟 Features

### 📸 AI Receipt Scanning
- Upload receipt photos and let AI automatically extract transaction data
- Supports various receipt formats
- Instant data entry without manual typing

### 📊 Financial Dashboard
- Real-time income and expense tracking
- Daily, weekly, and monthly cashflow analytics
- Visual charts and graphs for better insights
- Net profit/loss calculations

### 📦 Product Management
- Add, edit, and manage product catalog
- Track inventory levels
- Set product prices and stock quantities
- Product image uploads

### 🛒 Order Management
- View and manage customer orders
- Track order status (pending, confirmed, cancelled)
- Customer information management
- Order history and analytics

### 🤖 Telegram Bot Integration
- Receive orders directly through Telegram
- Automatic transaction recording
- Real-time notifications

### 💰 Manual Transaction Entry
- Quick manual income/expense entry
- Date selection for backdated entries
- Transaction categorization

---

## 🚀 Deployment

The web application is deployed on Vercel at [imphnen-one.vercel.app](https://imphnen-one.vercel.app/)

---

## 💻 Running Locally

### Prerequisites

- **Node.js** (v20.19.0 or higher recommended) OR **Bun** (v1.0.0 or higher)
- **npm** (comes with Node.js) or **Bun**
- **Git**

### Installation Steps

#### 1. Clone the Repository

```bash
git clone <repository-url>
cd <repository-folder>/frontend
```

#### 2. Install Dependencies

**Using npm:**
```bash
npm install
```

**Using Bun:**
```bash
bun install
```

#### 3. Configure Environment Variables

Create a `.env` file in the frontend directory:

```env
# API Base URL (adjust to your backend server)
VITE_API_BASE_URL=http://localhost:8080/api/v1

# Node environment
NODE_ENV=development
```

#### 4. Run Development Server

**Using npm:**
```bash
npm run dev
```

**Using Bun:**
```bash
bun run dev
```

The application will be available at `http://localhost:5173`

#### 5. Build for Production

**Using npm:**
```bash
npm run build
npm run preview
```

**Using Bun:**
```bash
bun run build
bun run preview
```

### Available Scripts

| Command | npm | Bun | Description |
|---------|-----|-----|-------------|
| Development | `npm run dev` | `bun run dev` | Start development server |
| Build | `npm run build` | `bun run build` | Build for production |
| Preview | `npm run preview` | `bun run preview` | Preview production build |
| Type Check | `npm run check` | `bun run check` | Type checking with TypeScript |
| Format | `npm run format` | `bun run format` | Format code with Prettier |
| Lint | `npm run lint` | `bun run lint` | Lint code with ESLint |

---

## 🏗️ Tech Stack

### Frontend
- **SvelteKit** - Web application framework
- **TypeScript** - Type-safe JavaScript
- **Tailwind CSS v4** - Utility-first CSS framework
- **Flowbite Svelte** - UI component library
- **ApexCharts** - Data visualization
- **Lucide Icons** - Icon library

### Backend Integration
- RESTful API communication
- Cookie-based authentication
- Server-side data fetching

---

## 📁 Project Structure

```
frontend/
├── src/
│   ├── lib/
│   │   ├── components/      # Reusable Svelte components
│   │   ├── server/          # Server-side utilities
│   │   ├── types/           # TypeScript type definitions
│   │   └── utils/           # Helper functions
│   ├── routes/
│   │   ├── auth/            # Authentication pages
│   │   ├── dashboard/       # Dashboard pages
│   │   │   ├── analytics/   # Analytics page
│   │   │   ├── config/      # Settings page
│   │   │   └── orders/      # Orders page
│   │   └── +layout.svelte   # Root layout
│   └── app.html             # HTML template
├── static/                  # Static assets
└── package.json
```

---

## 🔐 Authentication

The application uses cookie-based authentication:

1. **Login** - Users authenticate with email and password
2. **Register** - New users can create an account
3. **Session Management** - Access tokens stored in HTTP-only cookies
4. **Auto-redirect** - Automatic redirect to login when session expires

---

## 🎨 UI/UX Features

- **Responsive Design** - Works on mobile, tablet, and desktop
- **Dark Mode Pattern** - Teal color scheme with pattern backgrounds
- **Loading States** - Visual feedback for async operations
- **Form Validation** - Client-side and server-side validation
- **Toast Notifications** - Success/error feedback
- **Modal Dialogs** - Clean interaction patterns

---

## 📝 License

© 2025 IMPHNEN. All rights reserved.

Powered by [Kolosal.ai](https://kolosal.ai)

---
<a id="indonesian_ver"></a>
# INDONESIAN VERSION

---

# IMPHNEN

**Ingin Menjadi Pengusaha Handal Namun Enggan Ngebuku**

Aplikasi web yang membantu pelaku UMKM (Usaha Mikro Kecil Menengah) mencatat nota dan orderan pelanggan dengan mudah. Pelaku usaha cukup foto struk belanja, upload ke aplikasi, dan biarkan AI memindai serta memasukkan data secara otomatis tanpa perlu mengetik manual. Terintegrasi juga dengan WhatsApp Business sehingga AI dapat mendeteksi transaksi dan langsung memasukkannya ke data. Pelaku UMKM memiliki lebih banyak waktu untuk hal yang penting: **Meningkatkan kualitas pelayanan dan mengembangkan UMKM-nya**.

---

[Go to English Version](#english_ver)

---

## 🌟 Fitur

### 📸 Pemindaian Struk dengan AI
- Upload foto struk dan biarkan AI mengekstrak data transaksi secara otomatis
- Mendukung berbagai format struk
- Entri data instan tanpa mengetik manual

### 📊 Dashboard Keuangan
- Pelacakan pemasukan dan pengeluaran real-time
- Analitik arus kas harian, mingguan, dan bulanan
- Grafik dan diagram visual untuk wawasan lebih baik
- Perhitungan laba/rugi bersih

### 📦 Manajemen Produk
- Tambah, edit, dan kelola katalog produk
- Lacak level inventaris
- Atur harga dan jumlah stok produk
- Upload gambar produk

### 🛒 Manajemen Pesanan
- Lihat dan kelola pesanan pelanggan
- Lacak status pesanan (pending, dikonfirmasi, dibatalkan)
- Manajemen informasi pelanggan
- Riwayat dan analitik pesanan

### 🤖 Integrasi Bot Telegram
- Terima pesanan langsung melalui Telegram
- Pencatatan transaksi otomatis
- Notifikasi real-time

### 💰 Entri Transaksi Manual
- Entri pemasukan/pengeluaran manual cepat
- Pemilihan tanggal untuk entri mundur
- Kategorisasi transaksi

---

## 🚀 Deployment

Aplikasi web ini di-deploy dengan Vercel pada domain [imphnen-one.vercel.app](https://imphnen-one.vercel.app/)

---

## 💻 Menjalankan Secara Lokal

### Prasyarat

- **Node.js** (v20.19.0 atau lebih tinggi direkomendasikan) ATAU **Bun** (v1.0.0 atau lebih tinggi)
- **npm** (sudah termasuk dengan Node.js) atau **Bun**
- **Git**

### Langkah Instalasi

#### 1. Clone Repository

```bash
git clone <url-repository>
cd <folder-repository>/frontend
```

#### 2. Install Dependencies

**Menggunakan npm:**
```bash
npm install
```

**Menggunakan Bun:**
```bash
bun install
```

#### 3. Konfigurasi Environment Variables

Buat file `.env` di direktori frontend:

```env
# URL Base API (sesuaikan dengan server backend Anda)
VITE_API_BASE_URL=http://localhost:8080/api/v1

# Node environment
NODE_ENV=development
```

#### 4. Jalankan Development Server

**Menggunakan npm:**
```bash
npm run dev
```

**Menggunakan Bun:**
```bash
bun run dev
```

Aplikasi akan tersedia di `http://localhost:5173`

#### 5. Build untuk Production

**Menggunakan npm:**
```bash
npm run build
npm run preview
```

**Menggunakan Bun:**
```bash
bun run build
bun run preview
```

### Script yang Tersedia

| Perintah | npm | Bun | Deskripsi |
|---------|-----|-----|-------------|
| Development | `npm run dev` | `bun run dev` | Jalankan development server |
| Build | `npm run build` | `bun run build` | Build untuk production |
| Preview | `npm run preview` | `bun run preview` | Preview production build |
| Type Check | `npm run check` | `bun run check` | Type checking dengan TypeScript |
| Format | `npm run format` | `bun run format` | Format kode dengan Prettier |
| Lint | `npm run lint` | `bun run lint` | Lint kode dengan ESLint |

---

## 🏗️ Tech Stack

### Frontend
- **SvelteKit** - Framework aplikasi web
- **TypeScript** - JavaScript dengan type-safe
- **Tailwind CSS v4** - Framework CSS utility-first
- **Flowbite Svelte** - Library komponen UI
- **ApexCharts** - Visualisasi data
- **Lucide Icons** - Library ikon

### Integrasi Backend
- Komunikasi RESTful API
- Autentikasi berbasis cookie
- Server-side data fetching

---

## 📁 Struktur Project

```
frontend/
├── src/
│   ├── lib/
│   │   ├── components/      # Komponen Svelte yang dapat digunakan kembali
│   │   ├── server/          # Utilitas server-side
│   │   ├── types/           # Definisi tipe TypeScript
│   │   └── utils/           # Fungsi helper
│   ├── routes/
│   │   ├── auth/            # Halaman autentikasi
│   │   ├── dashboard/       # Halaman dashboard
│   │   │   ├── analytics/   # Halaman analitik
│   │   │   ├── config/      # Halaman pengaturan
│   │   │   └── orders/      # Halaman pesanan
│   │   └── +layout.svelte   # Layout root
│   └── app.html             # Template HTML
├── static/                  # Aset statis
└── package.json
```

---

## 🔐 Autentikasi

Aplikasi menggunakan autentikasi berbasis cookie:

1. **Login** - Pengguna autentikasi dengan email dan password
2. **Register** - Pengguna baru dapat membuat akun
3. **Manajemen Sesi** - Access token disimpan dalam HTTP-only cookie
4. **Auto-redirect** - Redirect otomatis ke login saat sesi berakhir

---

## 🎨 Fitur UI/UX

- **Responsive Design** - Berfungsi di mobile, tablet, dan desktop
- **Dark Mode Pattern** - Skema warna teal dengan latar belakang pattern
- **Loading States** - Feedback visual untuk operasi asinkron
- **Form Validation** - Validasi client-side dan server-side
- **Toast Notifications** - Feedback sukses/error
- **Modal Dialogs** - Pola interaksi yang bersih

---

## 📝 Lisensi

© 2025 IMPHNEN. Hak cipta dilindungi undang-undang.

Didukung oleh [Kolosal.ai](https://kolosal.ai)
