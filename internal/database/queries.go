package database

import (
	"database/sql"
	"time"
)

// GetCachedProgram retrieves a cached program from the database
func (db *DB) GetCachedProgram(input *GetCachedProgramInput) (*Program, error) {
	query := `
		SELECT p.id, p.tid, p.sid, p.ch_id, p.title, p.short_title, p.title_yomi, p.last_update, p.created_at, p.updated_at
		FROM program_cache pc
		JOIN programs p ON pc.tid = p.tid
		WHERE pc.service_id = ? AND pc.start_time = ? AND pc.duration = ? AND pc.expires_at > ?
	`

	now := time.Now()
	row := db.conn.QueryRow(query, input.ServiceID, input.StartTime, input.Duration, now)

	var program Program
	err := row.Scan(
		&program.ID,
		&program.TID,
		&program.SID,
		&program.ChID,
		&program.Title,
		&program.ShortTitle,
		&program.TitleYomi,
		&program.LastUpdate,
		&program.CreatedAt,
		&program.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No cached program found
		}
		return nil, err
	}

	return &program, nil
}

// SaveProgram saves a program to the database
func (db *DB) SaveProgram(input *SaveProgramInput) error {
	query := `
		INSERT INTO programs (tid, sid, ch_id, title, short_title, title_yomi, last_update, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(tid) DO UPDATE SET
			sid = excluded.sid,
			ch_id = excluded.ch_id,
			title = excluded.title,
			short_title = excluded.short_title,
			title_yomi = excluded.title_yomi,
			last_update = excluded.last_update,
			updated_at = excluded.updated_at
	`

	now := time.Now()
	_, err := db.conn.Exec(
		query,
		input.TID,
		input.SID,
		input.ChID,
		input.Title,
		input.ShortTitle,
		input.TitleYomi,
		input.LastUpdate,
		now,
	)
	return err
}

// SaveProgramCache saves a program cache entry to the database
func (db *DB) SaveProgramCache(input *SaveProgramCacheInput) error {
	// First ensure the program exists
	programQuery := `SELECT id FROM programs WHERE tid = ?`
	var programID int
	err := db.conn.QueryRow(programQuery, input.TID).Scan(&programID)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
		// Program doesn't exist, this shouldn't happen but we'll handle it gracefully
		return nil
	}

	query := `
		INSERT INTO program_cache (service_id, start_time, duration, tid, title, expires_at)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(service_id, start_time, duration) DO UPDATE SET
			tid = excluded.tid,
			title = excluded.title,
			expires_at = excluded.expires_at
	`

	_, err = db.conn.Exec(
		query,
		input.ServiceID,
		input.StartTime,
		input.Duration,
		input.TID,
		input.Title,
		input.ExpiresAt,
	)
	return err
}

// CleanupExpiredCache removes expired cache entries from the database
func (db *DB) CleanupExpiredCache() (int, error) {
	query := `DELETE FROM program_cache WHERE expires_at < ?`
	result, err := db.conn.Exec(query, time.Now())
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(rowsAffected), nil
}
