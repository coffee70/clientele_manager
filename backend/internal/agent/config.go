package agent

import "os"

const (
	defaultAPIBase = "https://dashboard.clientbook.com"
)

// Config holds agent configuration from environment variables.
type Config struct {
	// API URLs (placeholders until real endpoints are discovered)
	APIClients      string
	APIMessages     string
	APIOpportunities string

	// Database connection string (e.g. postgres://user:pass@localhost:5432/clientele?sslmode=disable)
	DatabaseURL string

	// Login wait timeout in seconds
	LoginTimeoutSeconds int
}

// LoadConfig loads configuration from environment variables.
func LoadConfig() *Config {
	base := getEnv("CLIENTBOOK_API_BASE", defaultAPIBase)
	return &Config{
		APIClients:       getEnv("CLIENTBOOK_API_CLIENTS", base+"/api/clients"),
		APIMessages:      getEnv("CLIENTBOOK_API_MESSAGES", base+"/api/messages"),
		APIOpportunities: getEnv("CLIENTBOOK_API_OPPORTUNITIES", base+"/api/opportunities"),
		DatabaseURL:      os.Getenv("DATABASE_URL"),
		LoginTimeoutSeconds: 300, // 5 minutes
	}
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
