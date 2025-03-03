package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"astralis/internal/core/domain"
)

const (
	meteorShowerArt = `
    *    *    *    *    *
  *   *   *   *   *   *
    *    *    *    *    *
 *   *   *   *   *   *
   *    *    *    *    *
`
	eclipseArt = `
      @@@@@@@@
    @@@@@@@@@@@@
  @@@@@@@@@@@@@@@@
 @@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@
`
)

func main() {
	baseURL := flag.String("api", "http://localhost:8080", "Base URL of the Astralis API")
	flag.Parse()

	// Get current time in RFC3339 format
	now := time.Now().Format(time.RFC3339)
	endTime := time.Now().AddDate(0, 1, 0).Format(time.RFC3339)

	// Fetch events from the API
	resp, err := http.Get(fmt.Sprintf("%s/events?start=%s&end=%s", *baseURL, now, endTime))
	if err != nil {
		fmt.Printf("Error fetching events: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("API returned error: %s\n", resp.Status)
		os.Exit(1)
	}

	var events []domain.Event
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		fmt.Printf("Error decoding response: %v\n", err)
		os.Exit(1)
	}

	// Display events with ASCII art
	for _, event := range events {
		fmt.Printf("\n=== %s ===\n", event.Title)
		fmt.Printf("Date: %s\n", event.StartTime.Format("2006-01-02"))
		fmt.Printf("Type: %s\n", event.Type)
		fmt.Printf("Description: %s\n", event.Description)

		// Display ASCII art based on event type
		switch event.Type {
		case domain.MeteorShower:
			fmt.Print(meteorShowerArt)
		case domain.Eclipse:
			fmt.Print(eclipseArt)
		default:
			fmt.Println("*    *    *")
		}

		fmt.Print(strings.Repeat("-", 50) + "\n")
	}

	if len(events) == 0 {
		fmt.Println("No upcoming astronomical events found.")
	}
}
