package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/0xjuanma/golazo/internal/fotmob"
)

func main() {
	fmt.Println("Testing FotMob Stats Data Retrieval")
	fmt.Println("===================================\n")

	// Create FotMob client
	client := fotmob.NewClient()

	// Test 1: Recent finished matches and upcoming matches (1 day)
	fmt.Println("Test 1: Fetching finished and upcoming matches for 1 day...")
	today := time.Now().UTC()
	dateFrom1d := today.AddDate(0, 0, -(1 - 1))
	fmt.Printf("  Date range: %s to %s\n", dateFrom1d.Format("2006-01-02"), today.Format("2006-01-02"))
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Fetch finished matches for 1 day
	matches1d, err := client.RecentFinishedMatches(ctx, 1)
	if err != nil {
		fmt.Printf("❌ Error fetching 1-day finished matches: %v\n", err)
	} else {
		fmt.Printf("✓ Found %d finished matches (1 day)\n", len(matches1d))
		if len(matches1d) > 0 {
			fmt.Println("\nSample finished matches (1 day):")
			for i, match := range matches1d {
				if i >= 5 {
					break
				}
				homeScore := "?"
				awayScore := "?"
				if match.HomeScore != nil {
					homeScore = fmt.Sprintf("%d", *match.HomeScore)
				}
				if match.AwayScore != nil {
					awayScore = fmt.Sprintf("%d", *match.AwayScore)
				}
				dateStr := "?"
				if match.MatchTime != nil {
					dateStr = match.MatchTime.Format("Jan 2 15:04")
				}
				fmt.Printf("  %d. %s %s-%s %s (%s) - %s [Status: %s]\n",
					i+1,
					match.HomeTeam.ShortName,
					homeScore,
					awayScore,
					match.AwayTeam.ShortName,
					dateStr,
					match.League.Name,
					match.Status,
				)
			}
		}
	}

	// Fetch upcoming matches for today
	upcoming, err := client.UpcomingMatches(ctx)
	if err != nil {
		fmt.Printf("❌ Error fetching upcoming matches: %v\n", err)
	} else {
		fmt.Printf("\n✓ Found %d upcoming matches (for today)\n", len(upcoming))
		if len(upcoming) > 0 {
			fmt.Println("\nSample upcoming matches:")
			for i, match := range upcoming {
				if i >= 5 {
					break
				}
				dateStr := "?"
				if match.MatchTime != nil {
					dateStr = match.MatchTime.Format("Jan 2 15:04")
				}
				fmt.Printf("  %d. %s vs %s (%s) - %s [Status: %s]\n",
					i+1,
					match.HomeTeam.ShortName,
					match.AwayTeam.ShortName,
					dateStr,
					match.League.Name,
					match.Status,
				)
			}
		}
	}

	fmt.Println("\n" + strings.Repeat("-", 50) + "\n")

	// Test 2: Recent finished matches (3 days)
	fmt.Println("Test 2: Fetching finished matches from last 3 days...")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel2()

	matches3d, err := client.RecentFinishedMatches(ctx2, 3)
	if err != nil {
		fmt.Printf("❌ Error fetching 3-day matches: %v\n", err)
	} else {
		fmt.Printf("✓ Found %d finished matches (3 days)\n", len(matches3d))
	}

	fmt.Println("\n" + strings.Repeat("-", 50) + "\n")

	// Test 3: Test match details if we have matches
	if len(matches1d) > 0 {
		fmt.Printf("Test 3: Fetching match details for first match (ID: %d)...\n", matches1d[0].ID)
		ctx6, cancel6 := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel6()

		details, err := client.MatchDetails(ctx6, matches1d[0].ID)
		if err != nil {
			fmt.Printf("❌ Error fetching match details: %v\n", err)
		} else {
			fmt.Printf("✓ Match details retrieved\n")
			fmt.Printf("  Home: %s, Away: %s\n", details.HomeTeam.Name, details.AwayTeam.Name)
			fmt.Printf("  Events: %d\n", len(details.Events))
			if details.Venue != "" {
				fmt.Printf("  Venue: %s\n", details.Venue)
			}
		}
	} else {
		fmt.Println("Test 3: Skipped (no matches available)")
	}

	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("Summary:")
	fmt.Printf("  1-day matches: %d\n", len(matches1d))
	fmt.Printf("  3-day matches: %d\n", len(matches3d))
	fmt.Printf("  Upcoming matches: %d\n", len(upcoming))
	fmt.Println(strings.Repeat("=", 50))
}
