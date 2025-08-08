# Wallet Sync API

A robust and scalable wallet management system built with Go and Fiber framework. This API provides comprehensive wallet functionality including user authentication, wallet operations, transaction management, and reconciliation services.

## Features

- **User Authentication**: JWT-based authentication with secure login/registration
- **Wallet Management**: Fund, withdraw, and transfer operations
- **Transaction History**: Paginated transaction history with filtering
- **Double-Entry Ledger**: Complete audit trail with ledger entries
- **Reconciliation Service**: Automated transaction reconciliation
- **Rate Limiting**: Built-in request rate limiting
- **Caching**: Redis-based caching for improved performance
- **Database Migrations**: Automated database schema management

## Tech Stack

- **Backend**: Go 1.24+ with Fiber framework
- **Database**: MySQL with GORM ORM
- **Cache**: Redis
- **Authentication**: JWT tokens
- **Validation**: Ozzo validation
- **Documentation**: Built-in monitoring dashboard

## Project Structure

```
.
├── dto/                    # Data Transfer Objects
├── handler/               # HTTP request handlers
├── internal/             # Internal packages
│   ├── config/           # Configuration management
│   ├── constants/        # Application constants
│   ├── helper/           # Utility helpers
│   └── seed/             # Database seeders
├── lib/                  # External libraries
│   └── database/         # Database connections
├── middleware/           # HTTP middlewares
├── migrations/           # Database migrations
├── model/                # Database models
├── payload/              # Request/Response structures
│   ├── request/          # Request DTOs
│   └── response/         # Response DTOs
├── repository/           # Data access layer
│   ├── core/             # Core repositories
│   └── user/             # User repositories
├── router/               # Route definitions
├── service/              # Business logic layer
└── validator/            # Input validation
```

## Prerequisites

- Go 1.24 or higher
- MySQL 5.7+
- Redis 6.0+

## Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd wallet-sync-api
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up environment variables**
   ```bash
   cp .env.example .env
   ```
   
   Configure your `.env` file:
   ```env
   PORT=8000
   
   DB_HOST=localhost
   DB_USER=root
   DB_PASSWORD=password
   DB_PORT=3306
   DB_NAME=wallet_sync
   
   REDIS_SERVER=localhost:6379
   ```

4. **Set up database**
   ```bash
   # Create MySQL database
   mysql -u root -p -e "CREATE DATABASE wallet_sync;"
   ```

5. **Run the application**
   ```bash
   go run main.go
   ```

The API will be available at `http://localhost:8000`

## API Endpoints

### Authentication

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/v1/auth/register` | Register a new user |
| POST | `/v1/auth/login` | Login user |

### Wallet Operations

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/v1/wallet/` | Get wallet details | ✅ |
| POST | `/v1/wallet/deposit` | Fund wallet | ✅ |
| POST | `/v1/wallet/withdraw` | Withdraw from wallet | ✅ |

### Transactions

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/v1/transaction/` | Get transaction history | ✅ |

### Monitoring

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/v1/monitor` | System monitoring dashboard |
| GET | `/health` | Health check endpoint |
| GET | `/logs/:key` | Get application logs |

## API Usage Examples

### Register User
```bash
curl -X POST http://localhost:8000/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123"
  }'
```

### Login
```bash
curl -X POST http://localhost:8000/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

### Fund Wallet
```bash
curl -X POST http://localhost:8000/v1/wallet/deposit \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-jwt-token>" \
  -d '{
    "amount": 1000.00
  }'
```

### Get Wallet Details
```bash
curl -X GET http://localhost:8000/v1/wallet/ \
  -H "Authorization: Bearer <your-jwt-token>"
```

### Get Transaction History
```bash
curl -X GET "http://localhost:8000/v1/transaction/?page=1&size=10&type=credit" \
  -H "Authorization: Bearer <your-jwt-token>"
```

## Database Schema

The application uses the following core models:

- **User**: User account information
- **Account**: Wallet accounts with balances
- **Transaction**: Financial transactions
- **LedgerEntry**: Double-entry bookkeeping records
- **ReconciliationLog**: Reconciliation audit logs

## Features in Detail

### Double-Entry Bookkeeping
Every transaction creates corresponding ledger entries ensuring financial accuracy and complete audit trails.

### Reconciliation Service
The [`ReconciliationService`](service/reconciliation_service.go) provides automated reconciliation of transactions and account balances.

### Validation
Input validation is handled by the [`validator`](validator/) package using ozzo-validation for robust data validation.

### Caching
Redis integration provides efficient caching through the [`RedisClientInterface`](lib/database/redis.go).

## Development

### Running Tests
```bash
go test ./...
```

### Building for Production
```bash
go build -o wallet_sync main.go
```

### Database Migrations
Migrations are automatically run on startup through the [`database.Migrate`](lib/database/database.go) function.

## Configuration

The application uses environment-based configuration managed by the [`config`](internal/config/) package. Key configuration options:

- **PORT**: HTTP server port
- **DB_***: Database connection settings
- **REDIS_SERVER**: Redis server address

## Security

- JWT-based authentication
- Password hashing using bcrypt
- Input validation and sanitization
- Rate limiting middleware
- CORS protection

## Monitoring

- Built-in monitoring dashboard at `/v1/monitor`
- Health check endpoint at `/health`
- Request logging with Redis storage
- Performance metrics tracking

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License.

## Support

For support or questions, please