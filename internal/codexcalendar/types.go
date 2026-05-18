package calendar

import "time"

// Race is a unified race calendar entry merged across sources.
type Race struct {
	ID             string
	Name           string
	NormalizedName string
	DateLocal      string // YYYY-MM-DD
	City           string
	Island         string
	URL            string
	Confidence     float64
	FirstSeen      time.Time
	LastSeen       time.Time
	Sources        []SourceAttribution
}

// SourceAttribution tracks where this race record came from.
type SourceAttribution struct {
	Source     string
	URL        string
	LastSeen   time.Time
	RawTitle   string
	Confidence float64
}
