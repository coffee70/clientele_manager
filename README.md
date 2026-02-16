# Clientele Manager

Clientele Manager syncs [Clientbook](https://dashboard.clientbook.com) CRM data into a local PostgreSQL database and provides an iOS/macOS app for access.

## Components

- **Backend (Go)**: An agent that automates browser login to Clientbook, fetches clients/messages/opportunities, and writes to PostgreSQL; plus a minimal API server
- **Frontend (Swift/SwiftUI)**: An iOS/macOS app with login UI (currently wired to a placeholder auth endpoint)
- **Scripts (Python)**: Utility scripts (e.g., app icon generation)
- **Infrastructure**: Docker Compose for PostgreSQL

## Prerequisites

- Go 1.24+
- Python 3 (for scripts)
- Docker
- Xcode (for frontend)

## Quick Start

1. **Start PostgreSQL**

   ```bash
   docker compose up -d
   ```

2. **Run the sync agent** (opens a visible Chrome window; log in to Clientbook when prompted)

   ```bash
   cd backend
   DATABASE_URL="postgres://clientele:clientele@localhost:5432/clientele?sslmode=disable" go run ./cmd/agent
   ```

3. **Run the API server** (listens on port 8081)

   ```bash
   cd backend
   go run ./cmd/api
   ```

4. **Build and run the frontend** – Open `frontend/Clientele Manager/Clientele Manager.xcodeproj` in Xcode, then build and run.

## Project Structure

```
clientele_manager/
├── backend/          # Go backend (API + sync agent)
├── frontend/         # SwiftUI iOS/macOS app
├── scripts/          # Python utility scripts
├── docker-compose.yml
└── requirements.txt  # Python dependencies for scripts
```

- [backend/README.md](backend/README.md) – Backend architecture and usage
- [frontend/README.md](frontend/README.md) – Frontend app structure and build
- [scripts/README.md](scripts/README.md) – Scripts overview and usage

## Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `DATABASE_URL` | For agent | PostgreSQL connection string (e.g. `postgres://clientele:clientele@localhost:5432/clientele?sslmode=disable`) |
| `CLIENTBOOK_API_BASE` | No | Base URL for Clientbook API (default: `https://dashboard.clientbook.com`) |
| `CLIENTBOOK_API_CLIENTS` | No | Clients API endpoint |
| `CLIENTBOOK_API_MESSAGES` | No | Messages API endpoint |
| `CLIENTBOOK_API_OPPORTUNITIES` | No | Opportunities API endpoint |
