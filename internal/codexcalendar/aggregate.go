package calendar

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

// Aggregate fetches all sources and returns a merged unified calendar.
func Aggregate(ctx context.Context, client *http.Client, adapters ...Adapter) ([]Race, error) {
	if client == nil {
		client = http.DefaultClient
	}

	var all []Race
	var errs []error
	for _, a := range adapters {
		races, err := a.Fetch(ctx, client)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", a.Name(), err))
			continue
		}
		all = append(all, races...)
	}

	merged := MergeCalendar(all)
	if len(errs) > 0 {
		return merged, errors.Join(errs...)
	}
	return merged, nil
}
