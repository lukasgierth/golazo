package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/0xjuanma/golazo/internal/fotmob"
)

func main() {
	client := fotmob.NewClient()
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	fmt.Println("Fetching matches for today...")
	fmt.Printf("Supported leagues: %v\n\n", fotmob.SupportedLeagues)
	
	// First, get all matches for today to verify API is working
	allMatches, err := client.MatchesByDate(ctx, time.Now())
	if err != nil {
		fmt.Printf("❌ Error fetching all matches: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Printf("✓ Found %d total matches for today\n", len(allMatches))
	
	// Count by status
	liveCount := 0
	finishedCount := 0
	notStartedCount := 0
	for _, match := range allMatches {
		switch match.Status {
		case "live":
			liveCount++
		case "finished":
			finishedCount++
		case "not_started":
			notStartedCount++
		}
	}
	
	fmt.Printf("  - Live: %d\n", liveCount)
	fmt.Printf("  - Finished: %d\n", finishedCount)
	fmt.Printf("  - Not Started: %d\n\n", notStartedCount)
	
	// Now get only live matches
	liveMatches, err := client.LiveMatches(ctx)
	if err != nil {
		fmt.Printf("❌ Error fetching live matches: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Printf("✓ Found %d live matches (started but not finished)\n\n", len(liveMatches))
	
	if len(liveMatches) == 0 {
		fmt.Println("No live matches found for today.")
		fmt.Println("This could be normal if there are no matches currently in progress.")
		if len(allMatches) > 0 {
			fmt.Println("\nSample matches from today:")
			for i, match := range allMatches {
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
				fmt.Printf("  %s %s-%s %s [%s] - Status: %s\n",
					match.HomeTeam.ShortName,
					homeScore,
					awayScore,
					match.AwayTeam.ShortName,
					match.League.Name,
					match.Status,
				)
			}
		}
	} else {
		fmt.Println("Live matches:")
		for i, match := range liveMatches {
			homeScore := "?"
			awayScore := "?"
			if match.HomeScore != nil {
				homeScore = fmt.Sprintf("%d", *match.HomeScore)
			}
			if match.AwayScore != nil {
				awayScore = fmt.Sprintf("%d", *match.AwayScore)
			}
			
			liveTime := ""
			if match.LiveTime != nil {
				liveTime = fmt.Sprintf(" (%s)", *match.LiveTime)
			}
			
			fmt.Printf("  %d. %s %s-%s %s [%s]%s\n",
				i+1,
				match.HomeTeam.ShortName,
				homeScore,
				awayScore,
				match.AwayTeam.ShortName,
				match.League.Name,
				liveTime,
			)
		}
	}
}

