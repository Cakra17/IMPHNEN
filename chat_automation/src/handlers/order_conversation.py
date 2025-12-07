from telegram import Update
from telegram.ext import ContextTypes, ConversationHandler
from telegram.ext import MessageHandler, CommandHandler
from telegram.ext import filters
from services import imphnen
from services.imphnen_data import CreateOrdersData, OrderItems, ListCustomerOrder, OrderItemDetail

CHOOSE_MERCHANT, CHOOSE_PRODUCTS, ADD_PRODUCTS = range(3)

class OrderConversationHandler:
    def __init__(self, imphnen_service: imphnen.ImphnenService):
        self.imphnen_service = imphnen_service

    async def create_order_start(self, update: Update, context: ContextTypes.DEFAULT_TYPE):
        customer_id = update.effective_user.id
        context.user_data['customer_id'] = customer_id
        
        customer = self.imphnen_service.getCustomerByID(customer_id)
        if isinstance(customer, str):
            await update.message.reply_text(f"❌ {customer}\n\nAnda belum terdaftar sebagai pelanggan. Silakan daftar terlebih dahulu.")
            return ConversationHandler.END
        
        merchants = self.imphnen_service.getAllMerchant()
        if isinstance(merchants, str):
            await update.message.reply_text(f"❌ {merchants}")
            return ConversationHandler.END
        
        if not merchants or len(merchants) == 0:
            await update.message.reply_text("❌ Tidak ada merchant yang tersedia saat ini.")
            return ConversationHandler.END
        
        merchant_list = "\n".join([f"{i+1}. {merchant['merchant_name']}" for i, merchant in enumerate(merchants)])
        
        await update.message.reply_text(
            f"🆕 *Buat Pesanan Baru*\n\n"
            f"Pelanggan: {customer['data']['name']}\n\n"
            f"Daftar Merchant:\n{merchant_list}\n\n"
            f"Masukkan nomor merchant (1-{len(merchants)}):"
        )
        
        context.user_data['merchants'] = merchants
        context.user_data['items'] = []
        
        return CHOOSE_MERCHANT
    
    async def choose_merchant(self, update: Update, context: ContextTypes.DEFAULT_TYPE):
        try:
            merchant_number = int(update.message.text.strip())
            merchants = context.user_data.get('merchants', [])
            
            if merchant_number < 1 or merchant_number > len(merchants):
                await update.message.reply_text(
                    f"❌ Nomor tidak valid. Silakan masukkan nomor antara 1-{len(merchants)}:"
                )
                return CHOOSE_MERCHANT
            
            selected_merchant = merchants[merchant_number - 1]
            merchant_id = selected_merchant['merchant_id']
            context.user_data['merchant_id'] = merchant_id
            
            products = self.imphnen_service.listProductByMerchant(merchant_id)
            if isinstance(products, str):
                await update.message.reply_text(f"❌ {products}\n\nSilakan masukkan nomor merchant yang valid:")
                return CHOOSE_MERCHANT
            
            if not products or len(products) == 0:
                await update.message.reply_text("❌ Merchant ini tidak memiliki produk. Silakan pilih merchant lain:")
                return CHOOSE_MERCHANT
            
            product_list = "\n".join([
                f"{i+1}. {product['name']} - Rp{product['price']:,}".replace(",", ".") + f" (Stok: {product['stock']})"
                for i, product in enumerate(products)
            ])
            
            await update.message.reply_text(
                f"📦 *Daftar Produk*\n\n"
                f"🏪 Merchant: {selected_merchant['merchant_name']}\n\n"
                f"{product_list}\n\n"
                f"Masukkan produk yang ingin dipesan (format: nomor,jumlah)\n"
                f"Contoh: 1,2 (untuk produk nomor 1 sebanyak 2)\n"
                f"Ketik 'SELESAI' jika sudah selesai memilih produk:"
            )
            
            context.user_data['products'] = products
            
            return ADD_PRODUCTS
            
        except ValueError:
            await update.message.reply_text("❌ Input tidak valid. Silakan masukkan nomor:")
            return CHOOSE_MERCHANT
    
    async def add_products(self, update: Update, context: ContextTypes.DEFAULT_TYPE):
        text = update.message.text.strip().upper()
        
        if text == 'SELESAI':
            if not context.user_data['items']:
                await update.message.reply_text("❌ Belum ada item yang dipilih. Silakan tambahkan item terlebih dahulu:")
                return ADD_PRODUCTS
            
            order_data = CreateOrdersData(
                customer_id=context.user_data['customer_id'],
                items=context.user_data['items'],
                merchant_id=context.user_data['merchant_id']
            )
            
            result = self.imphnen_service.createOrder(order_data)
            
            if isinstance(result, str):
                await update.message.reply_text(f"❌ {result}")
            else:
                items_text = ""
                products = context.user_data.get('products', [])
                for item in context.user_data['items']:
                    product = next((p for p in products if p['id'] == item.product_id), None)
                    product_name = product['name'] if product else item.product_id
                    items_text += f"   • {product_name} (x{item.quantity})\n"
                
                merchant_name = [merchant['merchant_name'] for merchant in context.user_data['merchants'] if merchant['merchant_id'] == context.user_data['merchant_id']]
                formatted_price = f"Rp{result['total_price']:,}".replace(",", ".")
                
                await update.message.reply_text(
                    f"✅ Pesanan berhasil dibuat!\n\n"
                    f"📋 Detail Pesanan:\n"
                    f"🆔 Order ID: {result['id']}\n"
                    f"👤 Pelanggan: {result['customer']['name']}\n"
                    f"🏪 Merchant: {merchant_name[0]}\n"
                    f"📦 Items:\n{items_text}"
                    f"💰 Total: {formatted_price}"
                )
                await update.message.reply_text(f"https://imphnen-one.vercel.app/payment/{result["id"]}")
                 
            context.user_data.clear()
            return ConversationHandler.END
        
        try:
            parts = text.split(',')
            if len(parts) != 2:
                raise ValueError("Format tidak valid")
            
            product_number = int(parts[0].strip())
            quantity = int(parts[1].strip())
            
            products = context.user_data.get('products', [])
            
            if product_number < 1 or product_number > len(products):
                await update.message.reply_text(
                    f"❌ Nomor produk tidak valid. Silakan masukkan nomor antara 1-{len(products)}\n"
                    f"Format: nomor,jumlah (contoh: 1,2)\n\n"
                    f"Item saat ini: {len(context.user_data['items'])}\n"
                    f"Tambah item lagi atau ketik 'SELESAI':"
                )
                return ADD_PRODUCTS
            
            if quantity <= 0:
                raise ValueError("Jumlah harus lebih dari 0")
            
            selected_product = products[product_number - 1]
            product_id = selected_product['id']
            product_name = selected_product['name']
            
            if quantity > selected_product['stock']:
                await update.message.reply_text(
                    f"❌ Stok tidak mencukupi. Stok tersedia: {selected_product['stock']}\n"
                    f"Silakan masukkan jumlah yang lebih kecil:"
                )
                return ADD_PRODUCTS
            
            context.user_data['items'].append(OrderItems(product_id=product_id, quantity=quantity))
            
            await update.message.reply_text(
                f"✅ Item ditambahkan: {product_name} (qty: {quantity})\n"
                f"Total items: {len(context.user_data['items'])}\n\n"
                f"Tambah item lagi (format: nomor,jumlah) atau ketik 'SELESAI':"
            )
            return ADD_PRODUCTS
            
        except (ValueError, IndexError) as e:
            await update.message.reply_text(
                "❌ Format tidak valid. Gunakan format: nomor,jumlah\n"
                f"Contoh: 1,2 (untuk produk nomor 1 sebanyak 2)\n\n"
                f"Item saat ini: {len(context.user_data['items'])}\n"
                f"Tambah item lagi atau ketik 'SELESAI':"
            )
            return ADD_PRODUCTS

    async def get_order_start(self, update: Update, context: ContextTypes.DEFAULT_TYPE):
        customer_id = update.effective_user.id
        
        orders_response = self.imphnen_service.listCustomerOrders(customer_id)
        
        if isinstance(orders_response, str):
            await update.message.reply_text(f"❌ {orders_response}")
            return ConversationHandler.END
        
        orders = []
        if isinstance(orders_response, dict) and "data" in orders_response:
            orders_data = orders_response["data"]
            for order in orders_data:

                items = []
                for order_item in order.get("order_items", []):
                    product = order_item.get("product", {})
                    items.append(OrderItemDetail(
                        product_name=product.get("name", "N/A"),
                        quantity=order_item.get("quantity", 0)
                    ))
                
                orders.append(ListCustomerOrder(
                    status=order.get("status", "N/A"),
                    total_price=order.get("total_price", 0),
                    user_id=order.get("user_id", ""),
                    order_date=order.get("order_date", "N/A"),
                    items=items
                ))
        
        if not orders:
            await update.message.reply_text("📭 Anda belum memiliki pesanan.")
            return ConversationHandler.END
        
        merchants = self.imphnen_service.getAllMerchant()
        if isinstance(merchants, str):
            await update.message.reply_text(f"❌ {merchants}")
            return ConversationHandler.END
        
        order_list = "📋 *Daftar Pesanan Anda*\n\n"
        for i, order in enumerate(orders, 1):

            merchant_name = "N/A"
            for merchant in merchants:
                if merchant.get("merchant_id") == order.user_id:
                    merchant_name = merchant.get("merchant_name", "N/A")
                    break
            
            formatted_date = order.order_date
            if order.order_date != "N/A":
                try:
                    from datetime import datetime
                    date_obj = datetime.fromisoformat(order.order_date.replace('Z', '+00:00'))
                    formatted_date = date_obj.strftime("%Y-%m-%d")
                except:
                    formatted_date = order.order_date
            
            formatted_price = f"Rp{order.total_price:,}".replace(",", ".")
            
            items_list = ""
            for item in order.items:
                items_list += f"      • {item.product_name} (x{item.quantity})\n"
            
            order_list += (
                f"{i}. 🏪 Merchant: {merchant_name}\n"
                f"   💰 Total: {formatted_price}\n"
                f"   📅 Tanggal: {formatted_date}\n"
                f"   📊 Status: {order.status}\n"
                f"   📦 Items:\n{items_list}\n"
            )
        
        await update.message.reply_text(order_list)
        return ConversationHandler.END


    async def cancel(self, update: Update, context: ContextTypes.DEFAULT_TYPE):
        await update.message.reply_text("❌ Operasi dibatalkan.")
        context.user_data.clear()
        return ConversationHandler.END

    def create_create_order_handler(self):
        return ConversationHandler(
            entry_points=[CommandHandler('buatorder', self.create_order_start)],
            states={
                CHOOSE_MERCHANT: [MessageHandler(filters.TEXT & ~filters.COMMAND, self.choose_merchant)],
                ADD_PRODUCTS: [MessageHandler(filters.TEXT & ~filters.COMMAND, self.add_products)],
            },
            fallbacks=[CommandHandler('cancel', self.cancel)],
        )

    def create_get_order_handler(self):
        return ConversationHandler(
            entry_points=[CommandHandler('lihatorder', self.get_order_start)],
            states={},
            fallbacks=[CommandHandler('cancel', self.cancel)],
        )