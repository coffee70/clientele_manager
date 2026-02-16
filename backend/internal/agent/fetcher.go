package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

// Client represents a client from Clientbook (placeholder shape).
type Client struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// Add other fields when API response is known
}

// Message represents a message from user-client interaction (placeholder shape).
type Message struct {
	ID       string `json:"id"`
	ClientID string `json:"client_id"`
	Body     string `json:"body"`
	// Add other fields when API response is known
}

// Opportunity represents a sales opportunity (placeholder shape).
type Opportunity struct {
	ID       string `json:"id"`
	ClientID string `json:"client_id"`
	// Add other fields when API response is known
}

// FetchResult holds the fetched data and any per-endpoint errors.
type FetchResult struct {
	Clients      []Client
	Messages     []Message
	Opportunities []Opportunity
	Errors       []string
}

// FetchData executes JavaScript in the page context to fetch from the configured API URLs.
// Session cookies are sent automatically since we run in the page context.
func FetchData(ctx context.Context, cfg *Config) (*FetchResult, error) {
	result := &FetchResult{}

	// Fetch clients
	clientsJSON, err := fetchURL(ctx, cfg.APIClients)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("clients: %v", err))
		log.Printf("Failed to fetch clients: %v", err)
	} else if clientsJSON != "" {
		// API may return { data: [...] } or direct array - try both
		var clients []Client
		if err := json.Unmarshal([]byte(clientsJSON), &clients); err != nil {
			var wrapper struct {
				Data []Client `json:"data"`
			}
			if err2 := json.Unmarshal([]byte(clientsJSON), &wrapper); err2 != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("clients parse: %v", err))
			} else {
				result.Clients = wrapper.Data
			}
		} else {
			result.Clients = clients
		}
	}

	// Fetch messages
	messagesJSON, err := fetchURL(ctx, cfg.APIMessages)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("messages: %v", err))
		log.Printf("Failed to fetch messages: %v", err)
	} else if messagesJSON != "" {
		var messages []Message
		if err := json.Unmarshal([]byte(messagesJSON), &messages); err != nil {
			var wrapper struct {
				Data []Message `json:"data"`
			}
			if err2 := json.Unmarshal([]byte(messagesJSON), &wrapper); err2 != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("messages parse: %v", err))
			} else {
				result.Messages = wrapper.Data
			}
		} else {
			result.Messages = messages
		}
	}

	// Fetch opportunities
	oppsJSON, err := fetchURL(ctx, cfg.APIOpportunities)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("opportunities: %v", err))
		log.Printf("Failed to fetch opportunities: %v", err)
	} else if oppsJSON != "" {
		var opps []Opportunity
		if err := json.Unmarshal([]byte(oppsJSON), &opps); err != nil {
			var wrapper struct {
				Data []Opportunity `json:"data"`
			}
			if err2 := json.Unmarshal([]byte(oppsJSON), &wrapper); err2 != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("opportunities parse: %v", err))
			} else {
				result.Opportunities = wrapper.Data
			}
		} else {
			result.Opportunities = opps
		}
	}

	return result, nil
}

// fetchURL runs fetch in the page context (cookies included) and returns the JSON string.
func fetchURL(ctx context.Context, url string) (string, error) {
	// JavaScript: fetch from page context so session cookies are sent
	js := fmt.Sprintf(`(async () => {
		try {
			const r = await fetch(%q);
			if (!r.ok) {
				return JSON.stringify({ error: r.status, statusText: r.statusText });
			}
			const d = await r.json();
			return JSON.stringify(d);
		} catch (e) {
			return JSON.stringify({ error: e.message });
		}
	})()`, url)

	var res string
	err := chromedp.Run(ctx,
		chromedp.Evaluate(js, &res, func(p *runtime.EvaluateParams) *runtime.EvaluateParams {
			return p.WithAwaitPromise(true)
		}),
	)
	if err != nil {
		return "", err
	}

	// Check for fetch error in response (error can be string or number)
	var errCheck struct {
		Error interface{} `json:"error"`
	}
	if err := json.Unmarshal([]byte(res), &errCheck); err == nil && errCheck.Error != nil {
		return "", fmt.Errorf("API error: %v", errCheck.Error)
	}

	return res, nil
}
