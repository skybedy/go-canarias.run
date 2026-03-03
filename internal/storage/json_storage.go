package storage

import (
	"context"
	"encoding/json"
	"os"
	"sync"

	"canarias.run/internal/models"
)

// JSONStorage je jednoduchá implementace Storage,
// která ukládá a čte závody ze souboru na disku (např. data.json).
type JSONStorage struct {
	filePath string
	mu       sync.RWMutex
}

// NewJSONStorage vytvoří novou instanci JSONStorage.
func NewJSONStorage(filePath string) *JSONStorage {
	return &JSONStorage{
		filePath: filePath,
	}
}

// SaveRaces bezpečně uloží závody do JSON souboru (zabrání souběžnému zápisu).
func (s *JSONStorage) SaveRaces(ctx context.Context, races []models.Race) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := json.MarshalIndent(races, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.filePath, data, 0644)
}

// GetAllRaces přečte závody z podkladového JSON souboru.
// Pokud soubor neexistuje, vrátí prázdné pole (nikoliv chybu).
func (s *JSONStorage) GetAllRaces(ctx context.Context) ([]models.Race, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, err := os.ReadFile(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []models.Race{}, nil
		}
		return nil, err
	}

	var races []models.Race
	if err := json.Unmarshal(data, &races); err != nil {
		return nil, err
	}

	return races, nil
}
