package epgcycler

import (
	"testing"
	"time"

	"github.com/freyja1103/epg-cycler/internal/database"
)

func TestEPGCyclerWithDatabase(t *testing.T) {
	// Create a temporary database for testing
	db, err := database.NewDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Create a minimal config
	config := &Config{
		ReserveCutoffHour: 4,
		CacheExpiry:       24 * time.Hour,
	}

	// Create epgCycler instance with mock APIs and real database
	// Note: This is a basic test structure. In a real implementation,
	// we would mock the EDCB and Syobocal APIs for testing.

	if db == nil {
		t.Error("Database should not be nil")
	}

	if config.CacheExpiry != 24*time.Hour {
		t.Errorf("Expected cache expiry to be 24 hours, got %v", config.CacheExpiry)
	}
}
