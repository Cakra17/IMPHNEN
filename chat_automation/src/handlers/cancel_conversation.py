from telegram import Update
from telegram.ext import ContextTypes, ConversationHandler
from telegram.ext import MessageHandler, CommandHandler
from telegram.ext import filters
from services import imphnen
from datetime import datetime

SELECT_ORDER_NUMBER, CANCEL_ORDER_CONFIRM = range(2)

class CancelConversationHandler:
    def __init__(self, imphnen_service: imphnen.ImphnenService):
        self.imphnen_service = imphnen_service

    async def cancel_order_start(self, update: Update, context: ContextTypes.DEFAULT_TYPE):
        customer_id = update.effective_user.id
        
        orders_response = self.imphnen_service.listCustomerOrders(customer_id)
        
        if isinstance(orders_response, str):
            await update.message.reply_text(f"❌ {orders_response}")
            return ConversationHandler.END
        

        pending_orders = []
        if isinstance(orders_response, dict) and "data" in orders_response:
            orders_data = orders_response["data"]
            for order in orders_data:
                if order.get("status", "").lower() == "pending":
                    pending_orders.append(order)
        
        if not pending_orders:
            await update.message.reply_text("📭 Anda tidak memiliki pesanan dengan status pending yang dapat dibatalkan.")
            return ConversationHandler.END
        
        merchants = self.imphnen_service.getAllMerchant()
        if isinstance(merchants, str):
            await update.message.reply_text(f"❌ {merchants}")
            return ConversationHandler.END
        
        order_list = "❌ *Batalkan Pesanan*\n\nPilih nomor pesanan yang ingin dibatalkan:\n\n"
        
        for i, order in enumerate(pending_orders, 1):
            merchant_name = "N/A"
            for merchant in merchants:
                if merchant.get("merchant_id") == order.get("user_id"):
                    merchant_name = merchant.get("merchant_name", "N/A")
                    break
            
            formatted_date = order.get("order_date", "N/A")
            if formatted_date != "N/A":
                try:
                    date_obj = datetime.fromisoformat(formatted_date.replace('Z', '+00:00'))
                    formatted_date = date_obj.strftime("%Y-%m-%d")
                except:
                    pass
            
            total_price = order.get("total_price", 0)
            formatted_price = f"Rp{total_price:,}".replace(",", ".")
            
            items_text = ""
            for order_item in order.get("order_items", []):
                product = order_item.get("product", {})
                product_name = product.get("name", "N/A")
                quantity = order_item.get("quantity", 0)
                items_text += f"      • {product_name} (x{quantity})\n"
            
            order_list += (
                f"{i}. 🏪 Merchant: {merchant_name}\n"
                f"   💰 Total: {formatted_price}\n"
                f"   📅 Tanggal: {formatted_date}\n"
                f"   📦 Items:\n{items_text}\n"
            )
        
        order_list += "Masukkan nomor pesanan (1-{}):".format(len(pending_orders))
        
        await update.message.reply_text(order_list)
        
        context.user_data['pending_orders'] = pending_orders
        context.user_data['merchants'] = merchants
        
        return SELECT_ORDER_NUMBER

    async def select_order_number(self, update: Update, context: ContextTypes.DEFAULT_TYPE):
        try:
            order_number = int(update.message.text.strip())
            pending_orders = context.user_data.get('pending_orders', [])
            
            if order_number < 1 or order_number > len(pending_orders):
                await update.message.reply_text(
                    f"❌ Nomor tidak valid. Silakan masukkan nomor antara 1-{len(pending_orders)}:"
                )
                return SELECT_ORDER_NUMBER
            
            selected_order = pending_orders[order_number - 1]
            order_id = selected_order.get("id")
            
            merchant_name = "N/A"
            merchants = context.user_data.get('merchants', [])
            for merchant in merchants:
                if merchant.get("merchant_id") == selected_order.get("user_id"):
                    merchant_name = merchant.get("merchant_name", "N/A")
                    break
            
            total_price = selected_order.get("total_price", 0)
            formatted_price = f"Rp{total_price:,}".replace(",", ".")
            
            items_text = ""
            for order_item in selected_order.get("order_items", []):
                product = order_item.get("product", {})
                product_name = product.get("name", "N/A")
                quantity = order_item.get("quantity", 0)
                items_text += f"   • {product_name} (x{quantity})\n"
            
            await update.message.reply_text(
                f"Apakah Anda yakin ingin membatalkan pesanan berikut?\n\n"
                f"🆔 Order ID: {order_id}\n"
                f"🏪 Merchant: {merchant_name}\n"
                f"💰 Total: {formatted_price}\n"
                f"📦 Items:\n{items_text}\n"
                f"Ketik 'YA' untuk konfirmasi atau 'TIDAK' untuk batal:"
            )
            
            context.user_data['cancel_order_id'] = order_id
            return CANCEL_ORDER_CONFIRM
            
        except ValueError:
            await update.message.reply_text("❌ Input tidak valid. Silakan masukkan nomor:")
            return SELECT_ORDER_NUMBER

    async def cancel_order_confirm(self, update: Update, context: ContextTypes.DEFAULT_TYPE):
        confirmation = update.message.text.upper()
        
        if confirmation == 'YA':
            order_id = context.user_data['cancel_order_id']
            result = self.imphnen_service.cancelOrder(order_id)
            
            if isinstance(result, dict) and result.get("message"):
                await update.message.reply_text(
                    f"✅ Pesanan berhasil dibatalkan!\n\n"
                    f"🆔 Order ID: {order_id}\n"
                    f"📅 Status: Dibatalkan"
                )
            else:
                await update.message.reply_text(f"❌ Gagal membatalkan pesanan: {result}")
        elif confirmation == 'TIDAK':
            await update.message.reply_text("❌ Pembatalan dibatalkan.")
        else:
            await update.message.reply_text("❌ Ketik 'YA' atau 'TIDAK':")
            return CANCEL_ORDER_CONFIRM
        
        context.user_data.clear()
        return ConversationHandler.END

    async def cancel(self, update: Update, context: ContextTypes.DEFAULT_TYPE):
        await update.message.reply_text("❌ Operasi dibatalkan.")
        context.user_data.clear()
        return ConversationHandler.END

    def create_cancel_order_handler(self):
        return ConversationHandler(
            entry_points=[CommandHandler('cancelorder', self.cancel_order_start)],
            states={
                SELECT_ORDER_NUMBER: [MessageHandler(filters.TEXT & ~filters.COMMAND, self.select_order_number)],
                CANCEL_ORDER_CONFIRM: [MessageHandler(filters.TEXT & ~filters.COMMAND, self.cancel_order_confirm)],
            },
            fallbacks=[CommandHandler('cancel', self.cancel)],
        )