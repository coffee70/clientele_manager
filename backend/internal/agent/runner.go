package agent

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

const (
	loginURL = "https://dashboard.clientbook.com/login"
)

// Run executes the agent flow: navigate to login, wait for user to log in, fetch data.
func Run(ctx context.Context, cfg *Config) (*FetchResult, error) {
	// Navigate to login page
	if err := chromedp.Run(ctx, chromedp.Navigate(loginURL)); err != nil {
		return nil, fmt.Errorf("navigate to login: %w", err)
	}

	fmt.Println("Please log in. Waiting for redirect...")

	// Wait for login (poll until URL changes away from /login)
	loginCtx, cancel := context.WithTimeout(ctx, time.Duration(cfg.LoginTimeoutSeconds)*time.Second)
	defer cancel()

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-loginCtx.Done():
			return nil, fmt.Errorf("login timeout after %d seconds", cfg.LoginTimeoutSeconds)
		case <-ticker.C:
			var loc string
			if err := chromedp.Run(ctx, chromedp.Location(&loc)); err != nil {
				continue
			}
			u, err := url.Parse(loc)
			if err != nil {
				continue
			}
			path := strings.TrimSuffix(u.Path, "/")
			if path != "/login" && path != "/Login" {
				// Login succeeded - URL changed
				goto loggedIn
			}
		}
	}

loggedIn:
	// Short sleep to allow dashboard to finish loading
	if err := chromedp.Run(ctx, chromedp.Sleep(2*time.Second)); err != nil {
		return nil, fmt.Errorf("sleep after login: %w", err)
	}

	fmt.Println("Login detected. Fetching data...")

	// Fetch data via JavaScript in page context
	result, err := FetchData(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("fetch data: %w", err)
	}

	return result, nil
}
