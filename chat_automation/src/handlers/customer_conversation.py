from telegram import Update
from telegram.ext import ContextTypes, ConversationHandler, MessageHandler, CommandHandler, filters
from services import imphnen
from services.imphnen_data import AddCustomerData, UpdateCustomerData
import re

ADD_CUSTOMER_NAME, ADD_CUSTOMER_ADDRESS, ADD_CUSTOMER_PHONE = range(3)
EDIT_CUSTOMER_VALIDATION, EDIT_CUSTOMER_FIELD, EDIT_CUSTOMER_VALUE = range(3)
DELETE_CUSTOMER_VALIDATION = 0
GET_CUSTOMER_ID = range(1)

class CustomerConversationHandler:
    def __init__(self, imphnen_service: imphnen.ImphnenService):
        self.imphnen_service = imphnen_service

    async def add_customer_start(self, update: Update, context: ContextTypes.DEFAULT_TYPE):
        context.user_data['customer_id'] = update.effective_user.id
        await update.message.reply_text("🆕 *Tambah Pelanggan Baru*\n\nMasukkan nama:")
        return ADD_CUSTOMER_NAME

    async def add_customer_name(self, update: Update, context: ContextTypes.DEFAULT_TYPE):
        context.user_data['customer_name'] = update.message.text
        await update.message.reply_text("Masukkan alamat:")
        return ADD_CUSTOMER_ADDRESS

    async def add_customer_address(self, update: Update, context: ContextTypes.DEFAULT_TYPE):
        context.user_data['customer_address'] = update.message.text
        await update.message.reply_text("Masukkan nomor telepon:")
        return ADD_CUSTOMER_PHONE

    async def add_customer_phone(self, update: Update, context: ContextTypes.DEFAULT_TYPE):
        phone = update.message.text
        if not re.match(r'^\+?[\d\s-]+$', phone):
            await update.message.reply_text("❌ Format nomor telepon tidak valid. Silakan coba lagi:")
            return ADD_CUSTOMER_PHONE

        customer_data = AddCustomerData(
            id=context.user_data['customer_id'],
            name=context.user_data['customer_name'],
            address=context.user_data['customer_address'],
            phone=phone
        )

        result = self.imphnen_service.addCustomer(customer_data)
        await update.message.reply_text(f"✅ {result}")
        
        context.user_data.clear()
        return ConversationHandler.END

    async def edit_customer_start(self, update: Update, context: ContextTypes.DEFAULT_TYPE):
        await update.message.reply_text("✏️ *Edit Pelanggan*\n\nApakah anda yakin ingin merubah data? (Y/N)")
        return EDIT_CUSTOMER_VALIDATION

    async def edit_customer_id(self, update: Update, context: ContextTypes.DEFAULT_TYPE):
        try:
            valid = update.message.text
            if valid.lower() == "n":
                await update.message.reply_text(" ✏️ *Edit Pelanggan* dibatalkan")
                return ConversationHandler.END
            elif valid.lower() != "n" and valid != "y":
                await update.message.reply_text(" ❌ Pilih antara Y/N")
                return EDIT_CUSTOMER_VALIDATION
            
            context.user_data['edit_customer_id'] = update.effective_user.id
            customer_id = update.effective_user.id
            
            customer = self.imphnen_service.getCustomerByID(customer_id)
            if isinstance(customer, str):
                await update.message.reply_text(f"❌ {customer}")
                return ConversationHandler.END
            
            context.user_data['current_customer'] = customer['data']
            await update.message.reply_text(
                f"Pelanggan ditemukan:\n"
                f"ID: {customer["data"]["id"]}\n"
                f"Nama: {customer["data"]["name"]}\n"
                f"Alamat: {customer["data"]["address"]}\n"
                f"Telepon: {customer["data"]["phone"]}\n\n"
                f"Pilih field yang ingin diedit:\n"
                f"1. Nama\n"
                f"2. Alamat\n"
                f"3. Telepon"
            )
            return EDIT_CUSTOMER_FIELD
        except ValueError:
            await update.message.reply_text("❌ Data anda tidak ditemukan, silahkan mendaftar terlebih dahulu")
            return ConversationHandler.END

    async def edit_customer_field(self, update: Update, context: ContextTypes.DEFAULT_TYPE):
        choice = update.message.text.strip()
        
        if choice in ['1', 'nama', 'name']:
            context.user_data['edit_field'] = 'name'
            await update.message.reply_text("Masukkan nama baru:")
        elif choice in ['2', 'alamat', 'address']:
            context.user_data['edit_field'] = 'address'
            await update.message.reply_text("Masukkan alamat baru:")
        elif choice in ['3', 'telepon', 'phone']:
            context.user_data['edit_field'] = 'phone'
            await update.message.reply_text("Masukkan nomor telepon baru:")
        else:
            await update.message.reply_text("❌ Pilihan tidak valid. Silakan pilih 1-3:")
            return EDIT_CUSTOMER_FIELD
        
        return EDIT_CUSTOMER_VALUE

    async def edit_customer_value(self, update: Update, context: ContextTypes.DEFAULT_TYPE):
        field = context.user_data['edit_field']
        value = update.message.text
        
        if field == 'phone' and not re.match(r'^\+?[\d\s-]+$', value):
            await update.message.reply_text("❌ Format nomor telepon tidak valid. Silakan coba lagi:")
            return EDIT_CUSTOMER_VALUE

        customer_id = context.user_data['edit_customer_id']
        current_customer = context.user_data['current_customer']
        
        update_data = UpdateCustomerData(
            name=current_customer.get('name') if field != 'name' else value,
            address=current_customer.get('address') if field != 'address' else value,
            phone=current_customer.get('phone') if field != 'phone' else value
        )

        result = self.imphnen_service.updateCustomer(customer_id, update_data)
        await update.message.reply_text(f"✅ {result}!")
        
        context.user_data.clear()
        return ConversationHandler.END

    async def delete_customer_start(self, update: Update, context: ContextTypes.DEFAULT_TYPE):
        customer_id = update.effective_user.id
        
        customer = self.imphnen_service.getCustomerByID(customer_id)
        if isinstance(customer, str):
            await update.message.reply_text(f"❌ {customer}")
            return ConversationHandler.END
        
        await update.message.reply_text(
            f"🗑️ *Hapus Pelanggan*\n\n"
            f"Apakah Anda yakin ingin menghapus pelanggan berikut?\n"
            f"ID: {customer["data"]["id"]}\n"
            f"Nama: {customer["data"]["name"]}\n"
            f"Alamat: {customer["data"]["address"]}\n"
            f"Telepon: {customer["data"]["phone"]}\n\n"
            f"Ketik ""YA"" untuk konfirmasi atau ""TIDAK"" untuk batal:"
        )

        context.user_data["delete_customer_id"] = customer_id
        return DELETE_CUSTOMER_VALIDATION

    async def delete_customer_id(self, update: Update, context: ContextTypes.DEFAULT_TYPE):
        confirmation = update.message.text
        
        if confirmation.upper() == 'YA':
            customer_id = context.user_data['delete_customer_id']
            result = self.imphnen_service.deleteCustomer(customer_id)
            await update.message.reply_text(f"✅ Pelanggan berhasil dihapus!")
        elif confirmation.upper() == 'TIDAK':
            await update.message.reply_text("❌ Penghapusan dibatalkan.")
        else:
            await update.message.reply_text("❌ Ketik 'YA' atau 'TIDAK':")
            return DELETE_CUSTOMER_VALIDATION
        
        context.user_data.clear()
        return ConversationHandler.END

    async def get_customer_start(self, update: Update, context: ContextTypes.DEFAULT_TYPE):
        customer_id = update.effective_user.id
        
        customer = self.imphnen_service.getCustomerByID(customer_id)
        
        if isinstance(customer, str):
            await update.message.reply_text(f"❌ {customer}")
        else:
            await update.message.reply_text(
                f"📋 *Detail Pelanggan*\n\n"
                f"🆔 ID: {customer["data"]["id"]}\n"
                f"👤 Nama: {customer["data"]["name"]}\n"
                f"📍 Alamat: {customer["data"]["address"]}\n"
                f"📞 Telepon: {customer["data"]["phone"]}"
            )
        
        context.user_data.clear()
        return ConversationHandler.END


    async def cancel(self, update: Update, context: ContextTypes.DEFAULT_TYPE):
        await update.message.reply_text("❌ Operasi dibatalkan.")
        context.user_data.clear()
        return ConversationHandler.END

    def create_add_customer_handler(self):
        return ConversationHandler(
            entry_points=[],
            states={
                ADD_CUSTOMER_NAME: [MessageHandler(filters.TEXT & ~filters.COMMAND, self.add_customer_name)],
                ADD_CUSTOMER_ADDRESS: [MessageHandler(filters.TEXT & ~filters.COMMAND, self.add_customer_address)],
                ADD_CUSTOMER_PHONE: [MessageHandler(filters.TEXT & ~filters.COMMAND, self.add_customer_phone)],
            },
            fallbacks=[CommandHandler('cancel', self.cancel)],
        )

    def create_edit_customer_handler(self):
        return ConversationHandler(
            entry_points=[],
            states={
                EDIT_CUSTOMER_VALIDATION: [MessageHandler(filters.TEXT & ~filters.COMMAND, self.edit_customer_id)],
                EDIT_CUSTOMER_FIELD: [MessageHandler(filters.TEXT & ~filters.COMMAND, self.edit_customer_field)],
                EDIT_CUSTOMER_VALUE: [MessageHandler(filters.TEXT & ~filters.COMMAND, self.edit_customer_value)],
            },
            fallbacks=[CommandHandler('cancel', self.cancel)],
        )

    def create_delete_customer_handler(self):
        return ConversationHandler(
            entry_points=[],
            states={
                DELETE_CUSTOMER_VALIDATION: [MessageHandler(filters.TEXT & ~filters.COMMAND, self.delete_customer_id)],
            },
            fallbacks=[CommandHandler('cancel', self.cancel)],
        )
        # )

    def create_get_customer_handler(self):
        return ConversationHandler(
            entry_points=[],
            states={
                GET_CUSTOMER_ID: [MessageHandler(filters.TEXT & ~filters.COMMAND, self.get_customer_id)],
            },
            fallbacks=[CommandHandler('cancel', self.cancel)],
        )