package health

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestDatabaseUnavailable(t *testing.T) {
	pool, err := pgxpool.New(t.Context(), "postgres://booking:booking@127.0.0.1:1/booking?sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	resp := httptest.NewRecorder()
	Handler(pool)(resp, httptest.NewRequest(http.MethodGet, "/api/health", nil))
	if resp.Code != http.StatusServiceUnavailable {
		t.Fatalf("got %d %s, want HTTP 503", resp.Code, resp.Body)
	}
	if got := resp.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("Content-Type = %q, want application/json", got)
	}
	var body struct {
		Code    int
		Message string
		Details []any
	}
	if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Code != 14 || body.Message != "database unavailable" || body.Details == nil || len(body.Details) != 0 {
		t.Fatalf("unexpected error response: %s", resp.Body)
	}
}
