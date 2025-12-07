# IMPHNEN

**Ingin Menjadi Pengusaha Handal Namun Enggan Ngebuku**

A comprehensive platform designed to help Indonesian UMKM (Micro, Small, and Medium Enterprises) automate financial record-keeping through AI-powered receipt scanning and streamline order management via Telegram bot integration.

## Demo Video

https://github.com/user-attachments/assets/b56591bb-3107-44c5-9c46-c8373de620b5

## Overview

IMPHNEN addresses two critical pain points for small business owners:

1. **Automated Expense Tracking**: Upload receipt photos and let AI extract transaction data automatically using Kolosal OCR technology
2. **Telegram-Based Ordering**: Enable customers to browse products and place orders directly through Telegram without needing a separate e-commerce platform

The platform provides real-time financial analytics, order management, and product catalog capabilities through an intuitive web dashboard.

## Architecture

### Tech Stack

**Backend**
- Go 1.25.4
- Chi router for HTTP routing
- PostgreSQL database
- JWT authentication
- Cloudinary for image storage
- Kolosal AI for OCR processing

**Frontend**
- SvelteKit 2 with Svelte 5
- TypeScript for type safety
- Tailwind CSS 4 with Flowbite components
- ApexCharts for data visualization
- Vercel deployment

**Chat Automation**
- Python 3.14 with python-telegram-bot
- Async conversation handlers
- Kolosal AI for customer service chatbot
- RESTful API integration

### System Components

```
┌─────────────────┐
│  Web Dashboard  │ (SvelteKit + Flowbite)
│  (Frontend)     │
└────────┬────────┘
         │
         │ HTTP/REST
         │
┌────────▼────────┐      ┌──────────────┐
│   Backend API   │◄────►│  PostgreSQL  │
│   (Go + Chi)    │      └──────────────┘
└────────┬────────┘
         │
         ├──────► Kolosal OCR (Receipt scanning)
         ├──────► Cloudinary (Image storage)
         │
         │ REST API
         │
┌────────▼────────┐
│  Telegram Bot   │ (Python)
│  (Chat Auto)    │
└────────┬────────┘
         │
         └──────► Kolosal AI (Customer service)
```

## Features

### Financial Management

**Receipt Scanning**
- Upload receipt images (PNG/JPG, up to 10MB)
- Automatic data extraction via Kolosal OCR
- Extract merchant name, date, items, prices, and totals
- Review and edit extracted data before saving
- Store original receipt images

**Transaction Analytics**
- Real-time income and expense tracking
- Daily, monthly, and custom date range reports
- Interactive cashflow charts
- Transaction filtering by type, source, and date
- Automated transaction creation from receipts and orders

**Financial Dashboard**
- Today's income and expense summary
- Recent transaction history
- Quick action buttons for manual entry
- Statistics: total income, total expense, net amount, transaction count

### Product Management

- Create, read, update, and delete products
- Product images with Cloudinary integration
- Stock tracking and inventory management
- Price management in Rupiah
- Product catalog accessible via Telegram bot

### Order Management

**Web Dashboard**
- View all orders with pagination
- Filter by status (pending, confirmed, cancelled)
- Customer information display
- Order details with line items
- Status management workflow

**Telegram Bot**
- Browse merchants and their products
- Add products to cart with quantity selection
- Place orders with delivery information
- View order history and status
- Cancel pending orders
- AI-powered customer service assistance

### Customer Management

- Customer profiles with contact information
- Telegram user ID integration
- Order history tracking
- Update customer details
- Account deletion support

## Database Schema

**Core Tables**
- `users` - UMKM owner accounts
- `receipts` - Scanned receipt records
- `receipt_items` - Line items from receipts
- `products` - Product catalog
- `customers` - Customer profiles
- `orders` - Order records
- `order_items` - Order line items
- `transactions` - Financial transactions (expenses and income)

**Key Features**
- UUID v7 primary keys
- Foreign key constraints with cascade
- Custom ENUM types for order status and transaction types
- Indexed for query performance
- Support for timezone-aware timestamps

## Installation

### Prerequisites

- Go 1.24 or higher
- Node.js 18+ or Bun
- Python 3.8+
- PostgreSQL 14+
- Telegram bot token (from @BotFather)
- Kolosal AI API key
- Cloudinary account

### Backend Setup

```bash
cd backend

# Copy environment variables
cp .env.example .env

# Edit .env with your configuration
# - Database credentials
# - JWT secret
# - Kolosal API key
# - Cloudinary credentials

# Install dependencies
go mod download

# Run database migrations
make migrate-up

# Start development server
make run

# Or build and run
make build
./bin/imphnen
```

**Backend Environment Variables**
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=imphnen
JWT_SECRET=your_jwt_secret
KOLOSAL_API_KEY=your_kolosal_key
CLOUDINARY_CLOUD_NAME=your_cloud_name
CLOUDINARY_API_KEY=your_api_key
CLOUDINARY_API_SECRET=your_api_secret
```

### Frontend Setup

```bash
cd frontend

# Install dependencies
bun install
# or: npm install

# Start development server
bun run dev

# Build for production
bun run build

# Preview production build
bun run preview
```

**Frontend Configuration**
- Backend API URL configured in `hooks.server.ts`
- Deployed to Vercel: https://imphnen-one.vercel.app/

### Telegram Bot Setup

```bash
cd chat_automation

# Create virtual environment
python -m venv .venv
source .venv/bin/activate  # Linux/Mac
# or: .venv\Scripts\activate  # Windows

# Install dependencies
pip install -r requirements.txt

# Copy environment variables
cp .env.example .env

# Edit .env with your configuration
# - Backend API URL
# - Telegram bot token
# - Kolosal API key

# Start bot
python src/main.py
```

**Bot Environment Variables**
```env
BACKEND_URL=http://localhost:8080/api/v1/telegram
BOT_TOKEN=your_telegram_bot_token
KOLOSAL_API_KEY=your_kolosal_api_key
```

## Usage

### Web Dashboard

1. **Register**: Create UMKM owner account
2. **Login**: Access your dashboard
3. **Upload Receipts**: Add expenses via receipt scanning
4. **Manage Products**: Add products for Telegram bot
5. **View Analytics**: Track income and expenses
6. **Manage Orders**: Process customer orders

### Telegram Bot

**Customer Commands**
- `/start` - Initialize bot
- `/help` - View all commands
- `/add_customer` - Register as customer
- `/get_customer` - View profile
- `/edit_customer` - Update information
- `/buatorder` - Create new order
- `/lihatorder` - View order history
- `/cancelorder` - Cancel pending order
- Chat freely for AI assistance

**Order Flow**
1. Customer registers using `/add_customer`
2. Browse merchants and products with `/buatorder`
3. Select products with quantity
4. Confirm order and receive payment link
5. Track order status with `/lihatorder`

## API Documentation

**Backend API**: Swagger documentation available at `/api/v1/docs/` when running the backend server

**Customer API**: See `docs/CUSTOMER_API.md` for detailed Telegram bot API documentation

**Base URL**: `http://localhost:8080/api/v1`

**Authentication**: Bearer token in Authorization header for protected endpoints

## Development

### Backend

```bash
# Run tests
make test

# Run database migrations
make migrate-up

# Rollback migrations
make migrate-down

# Generate Swagger docs
swag init -g cmd/main.go -o docs/swagger

# Build binary
make build
```

### Frontend

```bash
# Type checking
bun run check

# Format code
bun run format

# Lint code
bun run lint
```

### Testing

```bash
# Backend tests
cd backend && go test ./...

# Chat automation tests
cd chat_automation && pytest
```

## Deployment

**Frontend**: Automatically deployed to Vercel on push to main branch

**Backend**: Docker support available
```bash
cd backend
docker-compose up -d
```

**Bot**: Run as a long-lived process using systemd, supervisor, or container orchestration

## Security Considerations

- JWT-based authentication with token expiration
- Password hashing with bcrypt/Argon2
- HTTPS required for production
- Environment variables for sensitive credentials
- File upload validation and size limits
- SQL injection prevention via parameterized queries
- CORS configuration for frontend access

## Project Structure

```
imphnen/
├── backend/              # Go API server
│   ├── cmd/             # Application entry point
│   ├── internal/        # Internal packages
│   │   ├── handlers/    # HTTP handlers
│   │   ├── store/       # Database repositories
│   │   ├── models/      # Data models
│   │   ├── middleware/  # Auth middleware
│   │   └── config/      # Configuration
│   ├── pkg/             # Public packages
│   │   └── service/     # External services
│   └── db/              # Database migrations
├── frontend/            # SvelteKit web app
│   └── src/
│       ├── routes/      # SvelteKit routes
│       ├── lib/         # Shared components and utilities
│       └── hooks.server.ts  # Server hooks
├── chat_automation/     # Telegram bot
│   └── src/
│       ├── handlers/    # Bot command handlers
│       └── services/    # API integration services
└── docs/               # Documentation
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes with clear commit messages
4. Test thoroughly
5. Submit a pull request

## License

This project is proprietary software developed for Indonesian UMKM businesses.

## Support

For technical support or questions:
- Backend API: Check Swagger documentation
- Frontend: Review SvelteKit documentation
- Telegram Bot: Use `/help` command in bot
- Issues: Open an issue in the repository

## Acknowledgments

- Kolosal AI for OCR and AI chat capabilities
- Telegram Bot API for messaging platform
- Cloudinary for image management
- Flowbite for UI components
