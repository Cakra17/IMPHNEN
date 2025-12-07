# Chat Automation Bot

A Telegram bot for UMKM (Micro, Small, and Medium Enterprises) ordering platform. The bot provides AI-powered customer service and automated order management with integrated merchant and product catalog.

## Features

### Customer Management
- Register new customers with personal information
- View customer profile details
- Update customer information (name, address, phone)
- Delete customer accounts

### Order Management
- Create new orders with multiple products
- View order history and details
- Cancel pending orders
- Real-time stock validation
- Automatic price calculation

### AI Assistant
- Natural language interaction in Indonesian
- Product and merchant recommendations
- Contextual help and guidance
- Real-time merchant and product information

## Requirements

- Python 3.8+
- Telegram Bot Token
- Kolosal API Key
- Backend API URL

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd chat_automation
```

2. Create and activate virtual environment:
```bash
python -m venv .venv
source .venv/bin/activate  # On Linux/Mac
# OR
.venv\Scripts\activate  # On Windows
```

3. Install dependencies:
```bash
pip install -r requirements.txt
```

4. Configure environment variables:
```bash
cp .env.example .env
```

Edit `.env` file with your credentials:
```
BOT_TOKEN=your_telegram_bot_token
KOLOSAL_API_KEY=your_kolosal_api_key
BACKEND_URL=your_backend_api_url
```

## Usage

Start the bot:
```bash
python src/main.py
```

## Bot Commands

### General
- `/start` - Start the bot
- `/help` - Display all available commands

### Customer Management
- `/add_customer` - Register as a new customer
- `/get_customer` - View your profile information
- `/edit_customer` - Update profile information
- `/delete_customer` - Delete your account

### Order Management
- `/buatorder` - Create a new order
- `/lihatorder` - View your order history
- `/cancelorder` - Cancel a pending order

### Utilities
- `/cancel` - Cancel the current operation

## Project Structure

```
chat_automation/
├── src/
│   ├── main.py                      # Main application entry point
│   ├── handlers/                    # Command handlers
│   │   ├── customer_conversation.py # Customer management handlers
│   │   ├── order_conversation.py    # Order management handlers
│   │   └── cancel_conversation.py   # Order cancellation handlers
│   ├── services/                    # External service integrations
│   │   ├── imphnen.py              # Backend API service
│   │   ├── imphnen_data.py         # Data models
│   │   └── kolosal.py              # AI service integration
│   └── utils/                       # Utility functions
├── requirements.txt                 # Python dependencies
└── .env                            # Environment variables (not tracked)
```

## Order Process Flow

1. Customer selects a merchant by number
2. System displays available products with pricing and stock
3. Customer adds products using format: `number,quantity`
4. Customer can add multiple products or type `SELESAI` to complete
5. System validates stock and calculates total price
6. Order is created with status `pending`

## Order Status

- **pending** - Order created, awaiting processing (can be cancelled)
- **completed** - Order successfully processed
- **cancelled** - Order cancelled by customer


## License

This project is proprietary software. All rights reserved.

## Support

For technical issues or feature requests, contact the development team.
