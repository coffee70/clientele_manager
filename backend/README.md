# Backend

Go backend with two entry points: an HTTP API server and a sync agent that fetches data from Clientbook and writes to PostgreSQL.

## Structure

| Path | Description |
|------|--------------|
| `cmd/api` | HTTP API server (port 8081) |
| `cmd/agent` | Browser automation agent using chromedp |
| `internal/agent` | Config, runner (login flow), fetcher (API fetch via page JS) |
| `internal/db` | PostgreSQL client, schema, writer (upserts) |

## How to Run

**API server** (minimal placeholder):

```bash
go run ./cmd/api
```

**Sync agent** (opens visible Chrome window; requires manual login to Clientbook):

```bash
DATABASE_URL="postgres://clientele:clientele@localhost:5432/clientele?sslmode=disable" go run ./cmd/agent
```

## Agent Flow

1. Open Chrome (visible, non-headless)
2. Navigate to Clientbook login page
3. Wait for user to log in manually
4. Poll until URL changes away from `/login`
5. Fetch clients, messages, and opportunities via JavaScript in the page context (session cookies are sent automatically)
6. Write results to PostgreSQL (upsert into `clients`, `messages`, `sales_opportunities`)

## Database Schema

See [internal/db/schema.sql](internal/db/schema.sql). Tables:

- `clients` – External client records with `external_id`, `name`, `raw_json`
- `messages` – Messages linked to clients
- `sales_opportunities` – Opportunities linked to clients

## Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `DATABASE_URL` | For agent | PostgreSQL connection string. If unset, agent skips database writes. |
| `CLIENTBOOK_API_BASE` | No | Base URL (default: `https://dashboard.clientbook.com`) |
| `CLIENTBOOK_API_CLIENTS` | No | Clients endpoint (default: `{base}/api/clients`) |
| `CLIENTBOOK_API_MESSAGES` | No | Messages endpoint (default: `{base}/api/messages`) |
| `CLIENTBOOK_API_OPPORTUNITIES` | No | Opportunities endpoint (default: `{base}/api/opportunities`) |
