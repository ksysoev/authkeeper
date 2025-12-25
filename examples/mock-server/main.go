package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// MockOAuthServer is a simple mock OAuth2 server for testing
// Run with: go run examples/mock-server/main.go
func main() {
	http.HandleFunc("/oauth/token", handleToken)

	fmt.Println("🔐 Mock OAuth2 Server starting on :8080")
	fmt.Println("Token endpoint: http://localhost:8080/oauth/token")
	fmt.Println("Test credentials:")
	fmt.Println("  Client ID: test-client-id")
	fmt.Println("  Client Secret: test-client-secret")
	fmt.Println()

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func handleToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	grantType := r.FormValue("grant_type")
	clientID := r.FormValue("client_id")
	clientSecret := r.FormValue("client_secret")

	fmt.Printf("📥 Token request: grant_type=%s, client_id=%s\n", grantType, clientID)

	if grantType != "client_credentials" {
		http.Error(w, "Unsupported grant type", http.StatusBadRequest)
		return
	}

	// Simple validation
	if clientID == "" || clientSecret == "" {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// For demo purposes, accept any credentials
	// In real server, you'd validate against database
	scope := r.FormValue("scope")
	if scope == "" {
		scope = "read write"
	}

	response := map[string]interface{}{
		"access_token": generateMockToken(clientID),
		"token_type":   "Bearer",
		"expires_in":   3600,
		"scope":        strings.TrimSpace(scope),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	fmt.Printf("✅ Token issued for client: %s\n", clientID)
}

func generateMockToken(clientID string) string {
	// Simple mock token for demonstration
	return fmt.Sprintf("mock_token_%s_%d", clientID, 12345)
}
