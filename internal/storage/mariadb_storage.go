package storage

import (
	"context"
	"database/sql"
	"encoding/json"

	"canarias.run/internal/models"

	_ "github.com/go-sql-driver/mysql"
)

type MariaDBStorage struct {
	db *sql.DB
}

func NewMariaDBStorage(dsn string) (*MariaDBStorage, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return &MariaDBStorage{db: db}, nil
}

func (s *MariaDBStorage) SaveRaces(ctx context.Context, races []models.Race) error {
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

		var dateParsed any
		if race.DateParsed != "" && race.DateParsed != "0001-01-01" {
			dateParsed = race.DateParsed
		}

		if _, err := stmt.ExecContext(
			ctx,
			race.ID,
			race.Name,
			race.DateRaw,
			dateParsed,
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

	return tx.Commit()
}

func (s *MariaDBStorage) GetAllRaces(ctx context.Context) ([]models.Race, error) {
	rows, err := s.db.QueryContext(ctx, `
SELECT id, name, date_raw, COALESCE(DATE_FORMAT(date_parsed, '%Y-%m-%d'), ''), month, island, location,
       distances_json, source, status, url, type, description
FROM races
ORDER BY date_parsed ASC, name ASC`)
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
