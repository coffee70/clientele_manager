package db

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"clientele_manager/backend/internal/agent"

	"github.com/jackc/pgx/v5/pgxpool"
)

// CreateTables creates the schema if it does not exist.
func CreateTables(ctx context.Context, pool *pgxpool.Pool) error {
	schema := `
	CREATE TABLE IF NOT EXISTS clients (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		external_id TEXT UNIQUE NOT NULL,
		name TEXT,
		raw_json JSONB,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);
	CREATE TABLE IF NOT EXISTS messages (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		client_id UUID NOT NULL REFERENCES clients(id) ON DELETE CASCADE,
		external_id TEXT NOT NULL,
		body TEXT,
		raw_json JSONB,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW(),
		UNIQUE(client_id, external_id)
	);
	CREATE TABLE IF NOT EXISTS sales_opportunities (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		client_id UUID NOT NULL REFERENCES clients(id) ON DELETE CASCADE,
		external_id TEXT NOT NULL,
		raw_json JSONB,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW(),
		UNIQUE(client_id, external_id)
	);
	`
	_, err := pool.Exec(ctx, schema)
	return err
}

// WriteSync upserts clients, messages, and opportunities from the fetch result.
func WriteSync(ctx context.Context, pool *pgxpool.Pool, result *agent.FetchResult) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Map external_id -> our client UUID for resolving FKs
	clientIDMap := make(map[string]string)

	for _, c := range result.Clients {
		rawJSON, _ := json.Marshal(c)
		var id string
		err := tx.QueryRow(ctx, `
			INSERT INTO clients (external_id, name, raw_json)
			VALUES ($1, $2, $3)
			ON CONFLICT (external_id) DO UPDATE SET
				name = EXCLUDED.name,
				raw_json = EXCLUDED.raw_json,
				updated_at = NOW()
			RETURNING id::text
		`, c.ID, c.Name, rawJSON).Scan(&id)
		if err != nil {
			return fmt.Errorf("upsert client %s: %w", c.ID, err)
		}
		clientIDMap[c.ID] = id
	}

	for _, m := range result.Messages {
		clientUUID, ok := clientIDMap[m.ClientID]
		if !ok {
			log.Printf("Skipping message %s: client %s not found", m.ID, m.ClientID)
			continue
		}
		rawJSON, _ := json.Marshal(m)
		_, err := tx.Exec(ctx, `
			INSERT INTO messages (client_id, external_id, body, raw_json)
			VALUES ($1::uuid, $2, $3, $4)
			ON CONFLICT (client_id, external_id) DO UPDATE SET
				body = EXCLUDED.body,
				raw_json = EXCLUDED.raw_json,
				updated_at = NOW()
		`, clientUUID, m.ID, m.Body, rawJSON)
		if err != nil {
			return fmt.Errorf("upsert message %s: %w", m.ID, err)
		}
	}

	for _, o := range result.Opportunities {
		clientUUID, ok := clientIDMap[o.ClientID]
		if !ok {
			log.Printf("Skipping opportunity %s: client %s not found", o.ID, o.ClientID)
			continue
		}
		rawJSON, _ := json.Marshal(o)
		_, err := tx.Exec(ctx, `
			INSERT INTO sales_opportunities (client_id, external_id, raw_json)
			VALUES ($1::uuid, $2, $3)
			ON CONFLICT (client_id, external_id) DO UPDATE SET
				raw_json = EXCLUDED.raw_json,
				updated_at = NOW()
		`, clientUUID, o.ID, rawJSON)
		if err != nil {
			return fmt.Errorf("upsert opportunity %s: %w", o.ID, err)
		}
	}

	return tx.Commit(ctx)
}
