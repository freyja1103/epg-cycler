package database

import "time"

// Program represents a program entry in the database
type Program struct {
	ID         int       `db:"id"`
	TID        int       `db:"tid"`
	SID        int       `db:"sid"`   // EDCB Service ID
	ChID       int       `db:"ch_id"` // Syobocal Channel ID
	Title      string    `db:"title"`
	ShortTitle string    `db:"short_title"`
	TitleYomi  string    `db:"title_yomi"`
	LastUpdate string    `db:"last_update"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

// SaveProgramInput represents the input for saving a program
type SaveProgramInput struct {
	TID        int
	SID        int // EDCB Service ID (0 if not applicable)
	ChID       int // Syobocal Channel ID (0 if not applicable)
	Title      string
	ShortTitle string
	TitleYomi  string
	LastUpdate string
}

// SaveProgramCacheInput represents the input for saving a program cache entry
type SaveProgramCacheInput struct {
	ServiceID int // EDCB Service ID
	StartTime time.Time
	Duration  int
	TID       int
	Title     string
	ExpiresAt time.Time
}

// GetCachedProgramInput represents the input for retrieving a cached program
type GetCachedProgramInput struct {
	ServiceID int // EDCB Service ID
	StartTime time.Time
	Duration  int
}
