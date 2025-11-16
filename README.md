# Trading Dashboard

A full-stack trading application built with Go (backend) and React (frontend) featuring real-time price updates, order management, and WebSocket communication.

## Tech Stack

- **Backend**: Go with Gin framework, SQLite database
- **Frontend**: React 18 with Vite, Chart.js for price visualization
- **Communication**: REST API + WebSocket for real-time updates

## Prerequisites

- Go 1.19+ ([Download](https://golang.org/dl/))
- Node.js 16+ ([Download](https://nodejs.org/))
- npm 7+

## Project Structure

```
trading-dashboard/
├── backend/
│   ├── cmd/server/main.go          # Backend entry point
│   ├── internal/
│   │   ├── api/                    # API routes and middleware
│   │   ├── models/                 # Data models
│   │   ├── services/               # Business logic
│   │   └── websocket/              # WebSocket hub and clients
│   ├── go.mod                      # Go dependencies
│   └── Dockerfile
├── frontend/
│   ├── src/
│   │   ├── App.jsx                 # Main app component
│   │   ├── components/             # React components
│   │   ├── api/                    # API client functions
│   │   ├── hooks/                  # Custom React hooks
│   │   └── main.jsx                # React entry point
│   ├── package.json
│   ├── vite.config.js
│   └── Dockerfile
└── docker-compose.yml
```

## Running Locally

### Option 1: Run Backend and Frontend Separately (Recommended for Development)

#### 1. Start Backend Server

```bash
cd backend
go run ./cmd/server/main.go
```

The backend will start on `http://localhost:8080`

**Environment Variables (optional):**
```bash
set JWT_SECRET=your-secret-key
set DB_PATH=./data/trading.db
set GIN_MODE=release
```

Or on Linux/Mac:
```bash
export JWT_SECRET=your-secret-key
export DB_PATH=./data/trading.db
export GIN_MODE=release
go run ./cmd/server/main.go
```

#### 2. Start Frontend Development Server

In a new terminal:

```bash
cd frontend
npm install
npm run dev
```

The frontend will start on `http://localhost:5173`

### Option 2: Run with Docker Compose

```bash
docker-compose up --build
```

- Backend: `http://localhost:8080`
- Frontend: `http://localhost:3000`

## API Endpoints

### Authentication
- `POST /login` - Login with username

**Request:**
```json
{
  "username": "trader1"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Prices
- `GET /prices` - Get current prices for all symbols

### Orders
- `POST /orders` - Place a new order (requires authentication)
- `GET /orders` - Get all orders
- `GET /orders/:id` - Get specific order
- `POST /orders/:id/cancel` - Cancel an order (requires authentication)

### Holdings
- `GET /holdings` - Get user holdings (requires authentication)

### WebSocket
- `GET /ws` - WebSocket connection for real-time price updates

## Features

- ✅ User authentication with JWT tokens
- ✅ Real-time price updates via WebSocket
- ✅ Place buy/sell orders
- ✅ View order history and status
- ✅ Track holdings and portfolio
- ✅ Interactive price charts
- ✅ Order cancellation
- ✅ CORS enabled for development

## Troubleshooting

### Port Already in Use
If port 8080 or 5173 is already in use:

**Windows:**
```powershell
Get-NetTCPConnection -LocalPort 8080 | Select-Object -ExpandProperty OwningProcess | ForEach-Object { Stop-Process -Id $_ -Force }
```

**Linux/Mac:**
```bash
lsof -ti:8080 | xargs kill -9
```

### Token Expired Error
If you get a "401 Unauthorized" error with "token expired":
1. Hard refresh your browser (Ctrl+F5 or Cmd+Shift+R)
2. Log out and log back in to get a fresh token
3. Tokens expire after 24 hours

### CORS Issues
If you see CORS errors:
- Ensure the frontend is running on `http://localhost:5173`
- Ensure the backend is running on `http://localhost:8080`
- The backend has CORS enabled for all origins in development mode

### Database Errors
If you see database errors:
1. Delete the `backend/data/trading.db` file
2. Restart the backend - it will recreate the database
3. Log in again

## Development

### Build Frontend for Production
```bash
cd frontend
npm run build
```

Output will be in `frontend/dist/`

### Build Docker Images
```bash
docker build -t trading-dashboard-backend ./backend
docker build -t trading-dashboard-frontend ./frontend
```

## Default Login Credentials

Username: Any username (auto-creates user on first login)
- Try: `trader1`, `user123`, etc.

## Environment Configuration

### Backend Environment Variables
| Variable | Default | Description |
|----------|---------|-------------|
| JWT_SECRET | dev-secret | Secret key for JWT token signing |
| DB_PATH | ./data/trading.db | Path to SQLite database |
| GIN_MODE | debug | Set to "release" for production |

### Frontend Environment Variables
| Variable | Default | Description |
|----------|---------|-------------|
| VITE_API_URL | http://localhost:8080 | Backend API URL |
| VITE_WS_URL | ws://localhost:8080/ws | WebSocket URL |

## Production Deployment

1. Build Docker images
2. Use docker-compose or Kubernetes for orchestration
3. Set proper environment variables
4. Use a reverse proxy (nginx) for frontend serving
5. Enable HTTPS
6. Use a proper database (PostgreSQL) instead of SQLite

## Support

For issues or questions, please check the troubleshooting section above or review the application logs.


