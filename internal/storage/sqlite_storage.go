package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"

	"canarias.run/internal/models"

	_ "modernc.org/sqlite"
)

type SQLiteStorage struct {
	db *sql.DB
	mu sync.Mutex
}

func NewSQLiteStorage(dsn string) (*SQLiteStorage, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	storage := &SQLiteStorage{db: db}
	if err := storage.initSchema(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return storage, nil
}

func (s *SQLiteStorage) initSchema() error {
	const query = `
CREATE TABLE IF NOT EXISTS races (
	id TEXT,
	name TEXT NOT NULL,
	date_raw TEXT,
	date_parsed TEXT,
	month TEXT,
	island TEXT,
	location TEXT,
	distances_json TEXT,
	source TEXT,
	status TEXT,
	url TEXT,
	type TEXT,
	description TEXT
);
CREATE INDEX IF NOT EXISTS idx_races_date ON races(date_parsed);
CREATE INDEX IF NOT EXISTS idx_races_name ON races(name);
`
	_, err := s.db.Exec(query)
	if err != nil {
		return fmt.Errorf("init sqlite schema: %w", err)
	}
	return nil
}

func (s *SQLiteStorage) SaveRaces(ctx context.Context, races []models.Race) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if _, err := tx.ExecContext(ctx, `DELETE FROM races`); err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, `
INSERT INTO races (
	id, name, date_raw, date_parsed, month, island, location,
	distances_json, source, status, url, type, description
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, race := range races {
		distancesJSON, err := json.Marshal(race.Distances)
		if err != nil {
			return err
		}
		if _, err := stmt.ExecContext(
			ctx,
			race.ID,
			race.Name,
			race.DateRaw,
			race.DateParsed,
			race.Month,
			race.Island,
			race.Location,
			string(distancesJSON),
			race.Source,
			race.Status,
			race.URL,
			race.Type,
			race.Description,
		); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (s *SQLiteStorage) GetAllRaces(ctx context.Context) ([]models.Race, error) {
	rows, err := s.db.QueryContext(ctx, `
SELECT id, name, date_raw, date_parsed, month, island, location, distances_json, source, status, url, type, description
FROM races
ORDER BY date_parsed, name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	races := make([]models.Race, 0)
	for rows.Next() {
		var race models.Race
		var distancesJSON string
		if err := rows.Scan(
			&race.ID,
			&race.Name,
			&race.DateRaw,
			&race.DateParsed,
			&race.Month,
			&race.Island,
			&race.Location,
			&distancesJSON,
			&race.Source,
			&race.Status,
			&race.URL,
			&race.Type,
			&race.Description,
		); err != nil {
			return nil, err
		}
		if distancesJSON != "" {
			if err := json.Unmarshal([]byte(distancesJSON), &race.Distances); err != nil {
				return nil, err
			}
		}
		races = append(races, race)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return races, nil
}
