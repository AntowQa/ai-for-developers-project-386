package app

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	bookingv1 "github.com/AntowQa/ai-for-developers-project-386/api/gen/booking/v1"
	"github.com/AntowQa/ai-for-developers-project-386/api/internal/health"
)

func TestRunExitsWhenDatabaseDown(t *testing.T) {
	if os.Getenv("HEALTH_EXIT") == "1" {
		err := Run(context.Background())
		if err != nil {
			os.Exit(1)
		}
		os.Exit(0)
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestRunExitsWhenDatabaseDown")
	cmd.Env = replaceEnv(os.Environ(), "HEALTH_EXIT", "1")
	cmd.Env = replaceEnv(cmd.Env, "DATABASE_URL", "postgres://booking:booking@127.0.0.1:1/booking?sslmode=disable")
	err := cmd.Run()
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) || exitErr.ExitCode() == 0 {
		t.Fatalf("expected process to exit with an error, got %v", err)
	}
}

func TestHealth(t *testing.T) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Fatal("DATABASE_URL is required")
	}

	t.Run("ok", func(t *testing.T) {
		status, body := getHealth(t, dsn)
		if status != http.StatusOK {
			t.Fatalf("status %d body %s", status, body)
		}
		var payload struct {
			Status string `json:"status"`
		}
		if err := json.Unmarshal([]byte(body), &payload); err != nil {
			t.Fatal(err)
		}
		if payload.Status != "ok" {
			t.Fatalf("status field %q", payload.Status)
		}
	})

	t.Run("database unavailable", func(t *testing.T) {
		status, _ := getHealth(t, dsnWithDatabase(t, dsn, "missing_health_db"))
		if status != http.StatusServiceUnavailable {
			t.Fatalf("status %d, want 503", status)
		}
	})
}

func getHealth(t *testing.T, dsn string) (int, string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	t.Cleanup(cancel)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	grpcServer := grpc.NewServer()
	bookingv1.RegisterHealthServiceServer(grpcServer, health.New(pool))
	go func() {
		_ = grpcServer.Serve(lis)
	}()
	t.Cleanup(grpcServer.Stop)

	mux := runtime.NewServeMux()
	err = bookingv1.RegisterHealthServiceHandlerFromEndpoint(ctx, mux, lis.Addr().String(), []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	})
	if err != nil {
		t.Fatal(err)
	}

	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, server.URL+"/api/health", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("close body: %v", err)
		}
	}()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	return resp.StatusCode, string(body)
}

func dsnWithDatabase(t *testing.T, dsn, name string) string {
	t.Helper()
	parsed, err := url.Parse(dsn)
	if err != nil {
		t.Fatal(err)
	}
	parsed.Path = "/" + name
	return parsed.String()
}

func replaceEnv(env []string, key, value string) []string {
	prefix := key + "="
	out := make([]string, 0, len(env)+1)
	for _, item := range env {
		if !strings.HasPrefix(item, prefix) {
			out = append(out, item)
		}
	}
	return append(out, prefix+value)
}
