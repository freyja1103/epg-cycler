# Caching Program Titles with SQLite

## Overview

This document describes the implementation plan for caching program information fetched from Syobocal API using SQLite database. The cache will store program titles and metadata to reduce API calls and improve performance.

## Database Schema

### Programs Table
```sql
CREATE TABLE IF NOT EXISTS programs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    tid INTEGER UNIQUE NOT NULL,
    sid INTEGER DEFAULT 0,           -- EDCB Service ID (0 if not set)
    ch_id INTEGER DEFAULT 0,         -- Syobocal Channel ID (0 if not set)
    title TEXT NOT NULL,
    short_title TEXT,
    title_yomi TEXT,
    last_update TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### Program Cache Table
```sql
CREATE TABLE IF NOT EXISTS program_cache (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    service_id INTEGER NOT NULL,     -- EDCB Service ID
    start_time DATETIME NOT NULL,
    duration INTEGER NOT NULL,
    tid INTEGER NOT NULL,
    title TEXT NOT NULL,
    expires_at DATETIME NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (tid) REFERENCES programs(tid),
    UNIQUE(service_id, start_time, duration)
);
```

## Data Structures

### Program
```go
type Program struct {
    ID          int        `db:"id"`
    TID         int        `db:"tid"`
    SID         int        `db:"sid"`        // EDCB Service ID
    ChID        int        `db:"ch_id"`      // Syobocal Channel ID
    Title       string     `db:"title"`
    ShortTitle  string     `db:"short_title"`
    TitleYomi   string     `db:"title_yomi"`
    LastUpdate  string     `db:"last_update"`
    CreatedAt   time.Time  `db:"created_at"`
    UpdatedAt   time.Time  `db:"updated_at"`
}
```

### Input Structures
```go
type SaveProgramInput struct {
    TID        int
    SID        int        // EDCB Service ID (0 if not applicable)
    ChID       int        // Syobocal Channel ID (0 if not applicable)
    Title      string
    ShortTitle string
    TitleYomi  string
    LastUpdate string
}

type SaveProgramCacheInput struct {
    ServiceID  int        // EDCB Service ID
    StartTime  time.Time
    Duration   int
    TID        int
    Title      string
    ExpiresAt  time.Time
}

type GetCachedProgramInput struct {
    ServiceID int        // EDCB Service ID
    StartTime time.Time
    Duration  int
}
```

## Database Functions

1. `InitializeDB(dbPath string) (*sql.DB, error)` - Initialize database connection and create tables
2. `GetCachedProgram(input *GetCachedProgramInput) (*Program, error)` - Retrieve cached program info
3. `SaveProgram(input *SaveProgramInput) error` - Save program title info
4. `SaveProgramCache(input *SaveProgramCacheInput) error` - Save cache entry
5. `CleanupExpiredCache() (int, error)` - Remove expired cache entries

## Integration with epgCycler

1. Modify the epgCycler struct to include a database connection
2. Update the NewEPGCycler function to accept and initialize the database
3. Modify SimpleTidy to:
   - Check cache before calling Syobocal API
   - Save to cache after successful API calls
   - Handle database errors by returning them (fallback handled by epgCycler)

## Configuration

Add cache expiry configuration:
```go
type Config struct {
    ReserveCutoffHour int
    File              *ProgramFile
    CacheExpiry       time.Duration // Default: 90 days
}
```

## Default Settings

- Database file name: `programs.db`
- Cache expiry duration: 90 days (3 months)
- Database location: Same directory as executable
- Error handling: Return errors directly (fallback handled by epgCycler)