# XPomodoro Tracker API

A gamified Pomodoro productivity tracker REST API built with Go. This backend service provides user authentication, Pomodoro session tracking, statistics, rankings, and heatmap visualization for productivity data.

## Features

- **User Authentication & Management**

  - User registration and login with JWT tokens
  - Email verification for account updates
  - Password reset functionality via email codes
  - User profile management (country, email)

- **Pomodoro Tracking**

  - Create and track Pomodoro sessions
  - Session duration and completion tracking
  - XP (experience points) system for gamification

- **Statistics & Analytics**

  - User statistics (streaks, longest streak, current streak)
  - Heatmap visualization for daily Pomodoro activity
  - Historical data tracking

- **Ranking System**

  - Global rankings based on XP
  - Country-based local rankings
  - User rank lookup (global and local)

- **Gamification**
  - XP-based ranking system
  - Multiple rank tiers (starting from Wood I)
  - Progress tracking and leveling

## Tech Stack

- **Language**: Go 1.24.4
- **Web Framework**: Gorilla Mux
- **Database**: MySQL
- **Authentication**: JWT (golang-jwt/jwt/v5)
- **Password Hashing**: bcrypt (golang.org/x/crypto)
- **Email**: SMTP (Gmail)
- **API Documentation**: Swagger/OpenAPI
- **Environment Management**: godotenv

## Project Structure

```
backend/
├── cmd/
│   ├── main.go              # Application entry point
│   └── server/
│       └── server.go         # HTTP server setup and routing
├── config/
│   ├── constants.go         # Application constants
│   └── env.go               # Environment configuration
├── database/
│   └── database.go          # Database connection setup
├── docs/                    # Swagger documentation
├── helpers/
│   └── send_email.go        # Email sending utilities
├── middleware/
│   └── middleware.go        # JWT authentication & CORS middleware
├── migrations/              # Database migration files
├── services/
│   ├── auth/                # Authentication utilities (JWT, password hashing)
│   ├── pomodoros/           # Pomodoro session handlers
│   ├── ranking/             # Ranking system handlers
│   ├── stats/               # Statistics handlers
│   └── user/                # User management handlers
├── types/                   # Data models and interfaces
└── utils/                   # Utility functions
```

## Prerequisites

- Go 1.24.4 or higher
- MySQL 5.7+ or MySQL 8.0+
- Gmail account with App Password (for email functionality)

## Setup Instructions

### 1. Clone the Repository

```bash
git clone <repository-url>
cd backend
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Database Setup

Create a MySQL database:

```sql
CREATE DATABASE xpomodoro;
```

Run the migrations in order (they are numbered sequentially):

```bash
# Run migrations using your preferred migration tool
# The migrations are located in the migrations/ directory
```

Migration files:

- `001_create_ranks_table.sql` - Creates the ranks table
- `002_create_users_table.sql` - Creates the users table
- `00003_insert_ranks_data.sql` - Inserts initial rank data
- `00004_create_pomodoros_table.sql` - Creates pomodoros table
- `00005_create_stats_table.sql` - Creates stats table
- `00006_create_heatmap_table.sql` - Creates heatmap table
- `00007_sample_pomodoro_inserts.sql` - Sample data (optional)
- `00008_create_pending_email_updates_table.sql` - Email verification table
- `00009_add_default_country_to_users_table.sql` - Adds default country
- `00010_create_password_reset_table.sql` - Password reset table

### 4. Environment Configuration

Create a `.env` file in the root directory:

```env
# Database Configuration
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_HOST=localhost
DB_PORT=3306
DB_NAME=xpomodoro

# Server Configuration
PUBLIC_HOST=http://localhost
PORT=8000

# JWT Configuration
JWTSecret=your_secret_key_here
JWTExpirationInSeconds=2592000  # 30 days in seconds

# Email Configuration (Gmail)
GMAIL_TOKEN=your_gmail_app_password
```

**Note**: For Gmail, you'll need to:

1. Enable 2-Factor Authentication on your Google account
2. Generate an App Password: https://myaccount.google.com/apppasswords
3. Use the generated app password as `GMAIL_TOKEN`

### 5. Run the Application

```bash
go run cmd/main.go
```

The server will start on `http://127.0.0.1:8000`

## API Documentation

Once the server is running, you can access the Swagger documentation at:

```
http://localhost:8000/api/v1/swagger/
```

## API Endpoints

### Authentication (Public)

| Method | Endpoint                  | Description                 |
| ------ | ------------------------- | --------------------------- |
| POST   | `/api/v1/login`           | User login                  |
| POST   | `/api/v1/register`        | User registration           |
| GET    | `/api/v1/verify`          | Verify email with token     |
| POST   | `/api/v1/password/forgot` | Request password reset code |
| POST   | `/api/v1/password/reset`  | Reset password with code    |

### User Management (Protected)

| Method | Endpoint                     | Description                       |
| ------ | ---------------------------- | --------------------------------- |
| PUT    | `/api/v1/users/email`        | Update email (sends verification) |
| PATCH  | `/api/v1/users/{id}/country` | Update user country               |

### Pomodoros (Protected)

| Method | Endpoint           | Description                   |
| ------ | ------------------ | ----------------------------- |
| POST   | `/api/v1/pomodoro` | Create a new Pomodoro session |

### Statistics (Protected)

| Method | Endpoint                | Description            |
| ------ | ----------------------- | ---------------------- |
| GET    | `/api/v1/stats/{id}`    | Get user statistics    |
| POST   | `/api/v1/stats`         | Add user statistics    |
| PUT    | `/api/v1/stats`         | Update user statistics |
| GET    | `/api/v1/stats/heatmap` | Get user heatmap data  |
| PUT    | `/api/v1/stats/heatmap` | Upsert heatmap entry   |

### Rankings (Protected)

| Method | Endpoint                         | Description               |
| ------ | -------------------------------- | ------------------------- |
| GET    | `/api/v1/ranking/global`         | Get global ranking        |
| GET    | `/api/v1/ranking/global/{id}`    | Get user's global rank    |
| GET    | `/api/v1/ranking/{country}`      | Get country-based ranking |
| GET    | `/api/v1/ranking/{country}/{id}` | Get user's local rank     |

## Authentication

Most endpoints require JWT authentication. Include the token in the Authorization header:

```
Authorization: Bearer <your_jwt_token>
```

To obtain a token:

1. Register a new user via `/api/v1/register`
2. Login via `/api/v1/login` to receive a JWT token

## Database Schema

### Core Tables

- **ranks**: Rank tiers with minimum XP requirements
- **users**: User accounts with XP and rank tracking
- **pomodoros**: Pomodoro session records
- **stats**: User statistics (streaks, etc.)
- **heatmap**: Daily Pomodoro activity counts
- **pending_email_updates**: Email verification tokens
- **password_reset**: Password reset codes

## Development

### Running Tests

```bash
go test ./...
```

### Building

```bash
go build -o bin/server cmd/main.go
```

### Generating Swagger Docs

If you modify API documentation comments, regenerate Swagger docs:

```bash
swag init -g cmd/main.go
```

## Configuration Constants

The application uses the following constants (defined in `config/constants.go`):

- `GrowthFactor`: 0.01 (for XP calculations)
- `BaseMultiplier`: 1.0 (base XP multiplier)

## CORS

The API has CORS enabled for all origins. In production, you should restrict this to specific domains.

## License
This project is licensed for personal and educational use only.
Commercial use, redistribution, or use in any production environment is strictly prohibited without my explicit written permission.

