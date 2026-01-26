package database

import (
	"testing"
	"time"
)

func TestDatabaseInitialization(t *testing.T) {
	// Create a temporary database for testing
	db, err := NewDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Test that we can save and retrieve a program
	input := &SaveProgramInput{
		TID:        12345,
		SID:        1024,
		ChID:       1,
		Title:      "Test Program",
		ShortTitle: "Test",
		TitleYomi:  "テストプログラム",
		LastUpdate: "2023-01-01 00:00:00",
	}

	if err := db.SaveProgram(input); err != nil {
		t.Fatalf("Failed to save program: %v", err)
	}

	// Test that we can save and retrieve a program cache entry
	cacheExpiry := time.Now().Add(24 * time.Hour)
	cacheInput := &SaveProgramCacheInput{
		ServiceID: 1024,
		StartTime: time.Now(),
		Duration:  3600,
		TID:       12345,
		Title:     "Test Program",
		ExpiresAt: cacheExpiry,
	}

	if err := db.SaveProgramCache(cacheInput); err != nil {
		t.Fatalf("Failed to save program cache: %v", err)
	}

	// Test retrieving cached program
	getInput := &GetCachedProgramInput{
		ServiceID: 1024,
		StartTime: cacheInput.StartTime,
		Duration:  3600,
	}

	program, err := db.GetCachedProgram(getInput)
	if err != nil {
		t.Fatalf("Failed to get cached program: %v", err)
	}

	if program == nil {
		t.Fatal("Expected program, got nil")
	}

	if program.Title != "Test Program" {
		t.Errorf("Expected title 'Test Program', got '%s'", program.Title)
	}

	if program.TID != 12345 {
		t.Errorf("Expected TID 12345, got %d", program.TID)
	}

	// Test cleanup of expired cache
	expiredCount, err := db.CleanupExpiredCache()
	if err != nil {
		t.Fatalf("Failed to cleanup expired cache: %v", err)
	}

	// No expired entries should be cleaned up
	if expiredCount != 0 {
		t.Errorf("Expected 0 expired entries, got %d", expiredCount)
	}
}
