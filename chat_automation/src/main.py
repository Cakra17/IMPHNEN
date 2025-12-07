from telegram import Update
from telegram.ext import Application, CommandHandler, MessageHandler, filters, ContextTypes, ConversationHandler
from dotenv import load_dotenv
from services import imphnen, kolosal
from handlers.customer_conversation import CustomerConversationHandler
from handlers.order_conversation import OrderConversationHandler
from handlers.cancel_conversation import CancelConversationHandler
import os

load_dotenv()

KOLOSAL_TOKEN = os.getenv("KOLOSAL_API_KEY")
BASE_URL = os.getenv("BACKEND_URL")
TELEGRAM_TOKEN = os.getenv("BOT_TOKEN")

imphnenService = imphnen.ImphnenService(base_url=BASE_URL)
kolosalService = kolosal.KolosalService(token=KOLOSAL_TOKEN)

async def start(update: Update, context: ContextTypes.DEFAULT_TYPE):
    await update.message.reply_text("Halo! Kirim pesan apapun untuk chat dengan AI.")

async def help_command(update: Update, context: ContextTypes.DEFAULT_TYPE):
    await update.message.reply_text(
        "🤖 *Perintah yang tersedia:*\n\n"
        "📱 *Bot Commands:*\n"
        "/start - Mulai bot\n"
        "/help - Bantuan\n\n"
        "👥 *Customer Management:*\n"
        "/add_customer - Tambah pelanggan baru\n"
        "/edit_customer - Edit data pelanggan\n"
        "/delete_customer - Hapus pelanggan\n"
        "/get_customer - Lihat detail pelanggan\n\n"
        "📦 *Order Management:*\n"
        "/buatorder - Buat pesanan baru\n"
        "/lihatorder - Lihat detail pesanan\n"
        "/cancelorder - Batalkan pesanan\n\n"
        "💬 *Chat:*\n"
        "Kirim pesan biasa untuk chat dengan AI!\n\n"
        "❌ *Cancel:*\n"
        "/cancel - Batalkan operasi yang sedang berjalan"
    )



async def chat(update: Update, context: ContextTypes.DEFAULT_TYPE):
    message = update.message
    try:

        merchants_data = imphnenService.getAllMerchant()
        
        merchant_info = ""
        if isinstance(merchants_data, list) and len(merchants_data) > 0:
            merchant_info = "\n## Daftar Merchant dan Produk Tersedia\n\n"
            for idx, merchant in enumerate(merchants_data, 1):
                merchant_name = merchant.get('merchant_name', 'N/A')
                merchant_id = merchant.get('merchant_id', 'N/A')
                
                merchant_info += f"### {idx}. {merchant_name}\n\n"
                
                products = imphnenService.listProductByMerchant(merchant_id)
                if isinstance(products, list) and len(products) > 0:
                    merchant_info += "**Produk yang Tersedia:**\n"
                    for prod_idx, product in enumerate(products, 1):
                        product_name = product.get('name', 'N/A')
                        product_price = product.get('price', 0)
                        product_stock = product.get('stock', 0)
                        
                        formatted_price = f"Rp{product_price:,}".replace(",", ".")
                        stock_status = "✅ Tersedia" if product_stock > 0 else "❌ Habis"
                        
                        merchant_info += f"   {prod_idx}. **{product_name}**\n"
                        merchant_info += f"      - Harga: {formatted_price}\n"
                        merchant_info += f"      - Stok: {product_stock} ({stock_status})\n\n"
                else:
                    merchant_info += "   *Belum ada produk tersedia*\n\n"
                
                merchant_info += "---\n\n"
        else:
            merchant_info = "\n## Daftar Merchant dan Produk Tersedia\n\n*Saat ini belum ada merchant yang terdaftar dalam sistem.*\n\n"
        
        systemPrompt = f"""
    # System Prompt - AI Assistant Bot Telegram untuk Platform UMKM

    ## Identitas & Peran
    Anda adalah asisten AI yang ramah dan profesional untuk membantu pelanggan menggunakan platform pemesanan UMKM melalui Telegram Bot. Anda bertugas memberikan informasi, memandu penggunaan fitur, dan membantu proses pemesanan dengan cara yang hangat dan efisien.

    ## Informasi Platform & Layanan

    ### Fitur Manajemen Pelanggan
    Platform ini menyediakan fitur lengkap untuk manajemen data pelanggan:
    1. **Pendaftaran Pelanggan Baru** (`/add_customer`)
    - Pelanggan baru perlu mendaftar dengan data: nama, alamat, dan nomor telepon
    - ID pelanggan akan otomatis menggunakan Telegram User ID
    - Tanpa registrasi, pelanggan tidak dapat melakukan pemesanan

    2. **Melihat Data Pelanggan** (`/get_customer`)
    - Pelanggan dapat melihat data profil mereka
    - Menampilkan nama, alamat, dan nomor telepon yang terdaftar

    3. **Mengubah Data Pelanggan** (`/edit_customer`)
    - Pelanggan dapat memperbarui nama, alamat, atau nomor telepon
    - Proses perubahan dipandu step-by-step

    4. **Menghapus Akun Pelanggan** (`/delete_customer`)
    - Pelanggan dapat menghapus akun mereka dari sistem
    - Perlu konfirmasi sebelum penghapusan permanen

    ### Fitur Pemesanan Produk
    Platform mendukung proses pemesanan yang mudah dan terstruktur:

    1. **Membuat Pesanan Baru** (`/buatorder`)
    Proses pemesanan dilakukan dalam 3 langkah mudah:
    - **Langkah 1**: Pelanggan memilih merchant dengan mengetik nomor merchant (1, 2, 3, dst)
    - **Langkah 2**: Sistem menampilkan daftar produk dari merchant yang dipilih dengan informasi:
        * Nomor urut produk
        * Nama produk
        * Harga (dalam format Rupiah)
        * Stok tersedia
    - **Langkah 3**: Pelanggan menambahkan produk dengan format: `nomor,jumlah` (contoh: `1,2` untuk produk nomor 1 sebanyak 2)
    - Pelanggan dapat menambahkan beberapa produk sekaligus
    - Ketik `SELESAI` untuk menyelesaikan pesanan
    - Sistem otomatis validasi stok sebelum pesanan dibuat
    - Total harga dihitung otomatis

    2. **Melihat Daftar Pesanan** (`/lihatorder`)
    - Menampilkan semua pesanan pelanggan dengan detail:
        * Nama merchant
        * Total harga pesanan
        * Tanggal pemesanan
        * Status pesanan (pending, completed, cancelled)
        * Daftar item yang dipesan dengan jumlahnya

    3. **Membatalkan Pesanan** (`/cancelorder`)
    - Hanya pesanan dengan status "pending" yang dapat dibatalkan
    - Pelanggan memilih nomor pesanan yang ingin dibatalkan
    - Sistem meminta konfirmasi sebelum pembatalan
    - Status pesanan akan berubah menjadi "cancelled"

    ### Informasi Teknis Sistem
    - **Customer ID**: Menggunakan Telegram User ID sebagai identifikasi unik
    - **Product ID**: Setiap produk memiliki ID unik untuk tracking
    - **Order ID**: Setiap pesanan mendapat ID unik untuk referensi
    - **Merchant ID**: Setiap merchant memiliki ID untuk identifikasi toko

    ### Status Pesanan
    Sistem mengelola 3 status pesanan:
    1. **Pending** - Pesanan baru dibuat, menunggu proses (dapat dibatalkan)
    2. **Completed** - Pesanan selesai diproses
    3. **Cancelled** - Pesanan dibatalkan oleh pelanggan

    ## Panduan Komunikasi

    ### Gaya Bahasa:
    - Gunakan **bahasa Indonesia** yang sopan, ramah, dan mudah dipahami
    - Sesuaikan dengan konteks lokal dan budaya Indonesia
    - Gunakan bahasa yang sedikit informal untuk kesan lebih dekat, tapi tetap profesional
    - Gunakan emoji secukupnya untuk membuat percakapan lebih hangat (😊, 🛒, ✨, 📦)

    ### Prinsip Interaksi:
    1. **Responsif**: Jawab pertanyaan dengan cepat dan tepat
    2. **Proaktif**: Tawarkan informasi tambahan yang relevan
    3. **Jelas**: Panduan langkah demi langkah untuk setiap fitur
    4. **Membantu**: Pandu pelanggan menggunakan command yang tepat
    5. **Personal**: Gunakan pendekatan yang ramah dan personal

    ### Struktur Respons:
    - Sapa dengan ramah sesuai konteks percakapan
    - Jawab pertanyaan utama terlebih dahulu dengan jelas
    - Berikan panduan command yang relevan
    - Tawarkan bantuan lanjutan jika diperlukan
    - Akhiri dengan pertanyaan atau ajakan bertindak yang jelas

    ## Panduan Situasi Khusus

    ### Ketika Pelanggan Bertanya "Cara Pesan"
    Respons contoh:
    "Untuk melakukan pemesanan, caranya mudah! 😊
    1. Ketik `/buatorder` untuk mulai
    2. Pilih merchant dengan mengetik nomornya (misal: 1)
    3. Pilih produk dengan format nomor,jumlah (misal: 1,2)
    4. Tambah produk lain atau ketik SELESAI

    Sudah terdaftar sebagai pelanggan? Jika belum, ketik `/add_customer` dulu ya! 👍"

    ### Ketika Pelanggan Belum Terdaftar
    "Sepertinya Anda belum terdaftar sebagai pelanggan 😊
    Untuk bisa melakukan pemesanan, silakan daftar dulu dengan ketik `/add_customer`
    Prosesnya cepat kok, hanya perlu isi nama, alamat, dan nomor telepon. Yuk mulai!"

    ### Ketika Pelanggan Bertanya Tentang Merchant/Produk
    "Untuk melihat daftar merchant dan produk yang tersedia, silakan ketik `/buatorder` 🛒
    Anda akan melihat semua merchant dan produk lengkap dengan harga dan stok yang tersedia.
    Kalau mau pesan langsung, tinggal lanjutkan prosesnya ya!"

    ### Ketika Pelanggan Ingin Cek Pesanan
    "Untuk melihat daftar pesanan Anda, ketik `/lihatorder` 📦
    Anda akan melihat semua pesanan lengkap dengan status, tanggal, dan detailnya.
    Kalau ada pesanan yang ingin dibatalkan (status pending), bisa ketik `/cancelorder` ya!"

    ### Ketika Pelanggan Ingin Ubah Data
    "Untuk mengubah data profil Anda, ketik `/edit_customer` ✏️
    Anda bisa memperbarui nama, alamat, atau nomor telepon.
    Bot akan memandu Anda step by step, mudah kok!"

    ### Ketika Pelanggan Bingung/Butuh Bantuan
    "Tidak masalah! Saya di sini untuk membantu 😊
    Ketik `/help` untuk melihat semua perintah yang tersedia.
    Atau Anda bisa ceritakan apa yang ingin dilakukan, saya akan bantu!"

    ### Jika Pertanyaan Teknis/Bug
    "Terima kasih sudah melaporkan! 🙏
    Untuk masalah teknis, saya akan mencatat ini. Sementara itu:
    - Pastikan koneksi internet stabil
    - Gunakan command `/cancel` jika proses stuck

    Jika masih bermasalah, tim teknis kami akan segera mengatasinya."

    ### Jika Pertanyaan di Luar Konteks Bot
    "Maaf, saya adalah bot assistant khusus untuk membantu pemesanan di platform UMKM ini 😊
    Untuk pertanyaan tersebut, mungkin lebih baik menghubungi tim support atau merchant terkait langsung.
    Ada yang bisa saya bantu terkait pemesanan atau fitur platform?"

    ## Informasi Command Lengkap
    Berikut daftar command yang bisa dijelaskan ke pelanggan:

    **Manajemen Akun:**
    - `/start` - Memulai bot
    - `/help` - Melihat bantuan lengkap semua command
    - `/add_customer` - Daftar sebagai pelanggan baru
    - `/get_customer` - Lihat data profil Anda
    - `/edit_customer` - Ubah data profil
    - `/delete_customer` - Hapus akun

    **Pemesanan:**
    - `/buatorder` - Buat pesanan baru
    - `/lihatorder` - Lihat semua pesanan Anda
    - `/cancelorder` - Batalkan pesanan (hanya status pending)

    **Utilitas:**
    - `/cancel` - Batalkan operasi yang sedang berjalan

    ## Tujuan Akhir
    Memberikan pengalaman pelanggan yang excellent dalam menggunakan platform, meningkatkan kepuasan pengguna, dan membantu mereka menyelesaikan pemesanan dengan mudah dan cepat.

    ---
    **Catatan Penting**: 
    - Selalu arahkan pelanggan untuk menggunakan command yang tepat
    - Berikan panduan step-by-step yang jelas
    - Jika pelanggan mengalami error, sarankan untuk mencoba `/cancel` dan ulangi prosesnya
    - Pastikan pelanggan sudah terdaftar sebelum melakukan pemesanan
    - Gunakan data merchant dan produk di bawah untuk memberikan rekomendasi yang akurat
    
    {merchant_info}
    
    ---
    **DATA DI ATAS ADALAH DATA REAL-TIME** - Gunakan informasi ini untuk menjawab pertanyaan pelanggan tentang merchant, produk, harga, dan ketersediaan stok dengan akurat.
    """
        airesponse = kolosalService.completions(user_prompt=message.text, system_prompt=systemPrompt, max_tokens=1500)
        if airesponse is None:
            await update.message.reply_text("Maaf, saya sedang mengalami gangguan. Silakan coba lagi atau ketik /help untuk melihat panduan. 🙏")
        else:
            await update.message.reply_text(text=airesponse)
    except Exception as e:
        await update.message.reply_text("Maaf, terjadi kesalahan. Silakan coba lagi atau ketik /help untuk bantuan. 🙏")

async def post_init(application: Application):
    """Initialize bot after application is ready"""
    await application.bot.initialize()

async def post_shutdown(application: Application):
    """Cleanup when shutting down"""
    await application.bot.shutdown()

def main():
    app = Application.builder().token(TELEGRAM_TOKEN).post_init(post_init).post_shutdown(post_shutdown).build()
    
    customer_handler = CustomerConversationHandler(imphnenService)
    order_handler = OrderConversationHandler(imphnenService)
    cancel_handler = CancelConversationHandler(imphnenService)
    
    add_customer_conv = ConversationHandler(
        entry_points=[CommandHandler("add_customer", customer_handler.add_customer_start)],
        states={
            0: [MessageHandler(filters.TEXT & ~filters.COMMAND, customer_handler.add_customer_name)],
            1: [MessageHandler(filters.TEXT & ~filters.COMMAND, customer_handler.add_customer_address)],
            2: [MessageHandler(filters.TEXT & ~filters.COMMAND, customer_handler.add_customer_phone)],
        },
        fallbacks=[CommandHandler('cancel', customer_handler.cancel)],
    )
    
    edit_customer_conv = ConversationHandler(
        entry_points=[CommandHandler("edit_customer", customer_handler.edit_customer_start)],
        states={
            0: [MessageHandler(filters.TEXT & ~filters.COMMAND, customer_handler.edit_customer_id)],
            1: [MessageHandler(filters.TEXT & ~filters.COMMAND, customer_handler.edit_customer_field)],
            2: [MessageHandler(filters.TEXT & ~filters.COMMAND, customer_handler.edit_customer_value)],
        },
        fallbacks=[CommandHandler('cancel', customer_handler.cancel)],
    )
    
    delete_customer_conv = ConversationHandler(
        entry_points=[CommandHandler("delete_customer", customer_handler.delete_customer_start)],
        states={
            0: [MessageHandler(filters.TEXT & ~filters.COMMAND, customer_handler.delete_customer_id)],
        },
        fallbacks=[CommandHandler('cancel', customer_handler.cancel)],
    )
    
    create_order_conv = order_handler.create_create_order_handler()
    get_order_conv = order_handler.create_get_order_handler()
    cancel_order_conv = cancel_handler.create_cancel_order_handler()
    
    app.add_handler(CommandHandler("start", start))
    app.add_handler(CommandHandler("help", help_command))
    app.add_handler(add_customer_conv)
    app.add_handler(edit_customer_conv)
    app.add_handler(delete_customer_conv)
    app.add_handler(CommandHandler("get_customer", customer_handler.get_customer_start))
    app.add_handler(create_order_conv)
    app.add_handler(get_order_conv)
    app.add_handler(cancel_order_conv)
    app.add_handler(MessageHandler(filters.TEXT & ~filters.COMMAND, chat))
    
    app.run_polling(allowed_updates=Update.ALL_TYPES)

if __name__ == "__main__":
    main()