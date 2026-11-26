package health

import (
	"io"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Handler reports whether PostgreSQL accepts connections.
func Handler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := pool.Ping(r.Context()); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = io.WriteString(w, `{"code":14,"message":"database unavailable","details":[]}`)
			return
		}
		_, _ = io.WriteString(w, `{"status":"ok"}`)
	}
}
