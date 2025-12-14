package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/0xjuanma/golazo/internal/data"
)

func main() {
	// Get API key
	apiKey, err := data.FootballDataAPIKey()
	if err != nil {
		fmt.Printf("❌ Error getting API key: %v\n", err)
		fmt.Println("\nMake sure FOOTBALL_DATA_API_KEY is set:")
		fmt.Println("  export FOOTBALL_DATA_API_KEY=\"your-key-here\"")
		os.Exit(1)
	}

	fmt.Printf("✓ API key found (length: %d)\n\n", len(apiKey))

	// Test the API endpoint directly
	baseURL := "https://v3.football.api-sports.io"

	// Test 1: Get fixtures for today
	fmt.Println("Test 1: Fetching fixtures for today...")
	today := time.Now().Format("2006-01-02")
	url1 := fmt.Sprintf("%s/fixtures?date=%s", baseURL, today)
	testEndpoint(url1, apiKey)

	// Test 2: Get fixtures with date range (last 7 days)
	fmt.Println("\nTest 2: Fetching fixtures from last 7 days...")
	dateFrom := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	dateTo := time.Now().Format("2006-01-02")
	url2 := fmt.Sprintf("%s/fixtures?from=%s&to=%s", baseURL, dateFrom, dateTo)
	testEndpoint(url2, apiKey)

	// Test 3: Get finished matches only
	fmt.Println("\nTest 3: Fetching finished matches from last 7 days...")
	url3 := fmt.Sprintf("%s/fixtures?from=%s&to=%s&status=FT", baseURL, dateFrom, dateTo)
	testEndpoint(url3, apiKey)
}

func testEndpoint(url string, apiKey string) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("  ❌ Failed to create request: %v\n", err)
		return
	}

	req.Header.Set("x-apisports-key", apiKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("  ❌ Request failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	bodyStr := string(bodyBytes)

	fmt.Printf("  Status: %d\n", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("  ❌ Error response: %s\n", bodyStr[:min(500, len(bodyStr))])
		return
	}

	// Try to parse JSON
	var result map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		fmt.Printf("  ❌ Failed to parse JSON: %v\n", err)
		fmt.Printf("  Response: %s\n", bodyStr[:min(500, len(bodyStr))])
		return
	}

	// Check response structure
	if response, ok := result["response"].([]interface{}); ok {
		fmt.Printf("  ✓ Found %d matches\n", len(response))
		if len(response) > 0 {
			// Show first match
			if firstMatch, ok := response[0].(map[string]interface{}); ok {
				fmt.Printf("  Sample match: %+v\n", firstMatch)
			}
		}
	} else {
		fmt.Printf("  ⚠ Unexpected response structure: %+v\n", result)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

