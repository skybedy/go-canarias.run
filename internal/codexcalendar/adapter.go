package calendar

import (
	"context"
	"net/http"
)

type Adapter interface {
	Name() string
	Fetch(ctx context.Context, client *http.Client) ([]Race, error)
}
