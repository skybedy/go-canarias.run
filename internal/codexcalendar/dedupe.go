package calendar

import (
	"sort"
	"strings"
	"time"
)

func DedupeKey(r Race) string {
	return strings.Join([]string{
		strings.TrimSpace(r.DateLocal),
		NormalizeName(r.Name),
		NormalizeLocation(r.City, r.Island),
	}, "|")
}

func MergeRaces(existing Race, incoming Race) Race {
	out := existing

	if out.ID == "" {
		out.ID = incoming.ID
	}
	if out.Name == "" || len(incoming.Name) > len(out.Name) {
		out.Name = incoming.Name
	}
	if out.NormalizedName == "" {
		out.NormalizedName = NormalizeName(out.Name)
	}
	if out.DateLocal == "" {
		out.DateLocal = incoming.DateLocal
	}
	if out.City == "" {
		out.City = incoming.City
	}
	if out.Island == "" {
		out.Island = incoming.Island
	}
	if out.URL == "" {
		out.URL = incoming.URL
	}

	if incoming.Confidence > out.Confidence {
		out.Confidence = incoming.Confidence
		if incoming.URL != "" {
			out.URL = incoming.URL
		}
		if incoming.City != "" {
			out.City = incoming.City
		}
		if incoming.Island != "" {
			out.Island = incoming.Island
		}
	}

	if out.FirstSeen.IsZero() || (!incoming.FirstSeen.IsZero() && incoming.FirstSeen.Before(out.FirstSeen)) {
		out.FirstSeen = incoming.FirstSeen
	}
	if out.LastSeen.IsZero() || (!incoming.LastSeen.IsZero() && incoming.LastSeen.After(out.LastSeen)) {
		out.LastSeen = incoming.LastSeen
	}

	out.Sources = mergeSources(out.Sources, incoming.Sources)

	if out.NormalizedName == "" {
		out.NormalizedName = NormalizeName(out.Name)
	}

	return out
}

func MergeCalendar(in []Race) []Race {
	m := make(map[string]Race, len(in))
	for _, r := range in {
		r.NormalizedName = NormalizeName(r.Name)
		if r.ID == "" {
			r.ID = DedupeKey(r)
		}
		k := DedupeKey(r)
		if prev, ok := m[k]; ok {
			m[k] = MergeRaces(prev, r)
			continue
		}
		m[k] = r
	}

	out := make([]Race, 0, len(m))
	for _, r := range m {
		if r.ID == "" {
			r.ID = DedupeKey(r)
		}
		out = append(out, r)
	}

	sort.Slice(out, func(i, j int) bool {
		if out[i].DateLocal == out[j].DateLocal {
			return out[i].Name < out[j].Name
		}
		return out[i].DateLocal < out[j].DateLocal
	})

	return out
}

func mergeSources(a, b []SourceAttribution) []SourceAttribution {
	if len(a) == 0 {
		return append([]SourceAttribution(nil), b...)
	}
	out := append([]SourceAttribution(nil), a...)
	idx := make(map[string]int, len(out))
	for i, s := range out {
		idx[s.Source+"|"+s.URL] = i
	}
	for _, s := range b {
		key := s.Source + "|" + s.URL
		if i, ok := idx[key]; ok {
			if s.LastSeen.After(out[i].LastSeen) {
				out[i].LastSeen = s.LastSeen
			}
			if s.Confidence > out[i].Confidence {
				out[i].Confidence = s.Confidence
			}
			if out[i].RawTitle == "" {
				out[i].RawTitle = s.RawTitle
			}
			continue
		}
		idx[key] = len(out)
		out = append(out, s)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Source == out[j].Source {
			return out[i].URL < out[j].URL
		}
		return out[i].Source < out[j].Source
	})
	return out
}

func nowUTC() time.Time {
	return time.Now().UTC()
}
