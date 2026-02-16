package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"clientele_manager/backend/internal/agent"
	"clientele_manager/backend/internal/db"

	"github.com/chromedp/chromedp"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := agent.LoadConfig()

	// Init DB connection (optional)
	var pool *pgxpool.Pool
	if cfg.DatabaseURL != "" {
		ctx := context.Background()
		p, err := db.NewPool(ctx, cfg.DatabaseURL)
		if err != nil {
			log.Fatalf("Database connection failed: %v", err)
		}
		defer p.Close()
		pool = p

		if err := db.CreateTables(ctx, pool); err != nil {
			log.Fatalf("Create tables: %v", err)
		}
	} else {
		log.Println("DATABASE_URL not set - skipping database writes")
	}

	// Allocator options: visible window, no headless
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
	)
	allocCtx, cancelAlloc := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancelAlloc()

	ctx, cancelCtx := chromedp.NewContext(allocCtx)
	defer cancelCtx()

	// Run agent flow: navigate, wait for login, fetch data
	result, err := agent.Run(ctx, cfg)
	if err != nil {
		log.Fatalf("Agent run failed: %v", err)
	}

	// Write to database if configured
	if pool != nil {
		if err := db.WriteSync(ctx, pool, result); err != nil {
			log.Printf("Database write failed: %v", err)
		} else {
			fmt.Printf("Synced %d clients, %d messages, %d opportunities\n",
				len(result.Clients), len(result.Messages), len(result.Opportunities))
		}
	}
	if len(result.Errors) > 0 {
		for _, e := range result.Errors {
			log.Printf("Fetch error: %s", e)
		}
	}

	fmt.Println("Sync complete. Close the browser window or press Ctrl+C to exit.")

	// Wait for either: browser closed (ctx.Done) or SIGINT
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		// Browser was closed by user
	case <-sigCh:
		// Ctrl+C or kill - program will exit, defer will cancel context and close Chrome
	}
}
