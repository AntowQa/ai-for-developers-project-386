package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestAPI(t *testing.T) {
	binary := filepath.Join(t.TempDir(), "api")
	if output, err := exec.CommandContext(t.Context(), "go", "build", "-o", binary, ".").CombinedOutput(); err != nil {
		t.Fatalf("build API: %v\n%s", err, output)
	}

	t.Run("health", func(t *testing.T) {
		t.Setenv("HTTP_ADDR", "127.0.0.1:0")
		ctx, cancel := context.WithTimeout(t.Context(), 15*time.Second)
		defer cancel()

		cmd := exec.CommandContext(ctx, binary)
		cmd.Stderr = os.Stderr
		out, err := cmd.StdoutPipe()
		if err != nil {
			t.Fatal(err)
		}
		if err := cmd.Start(); err != nil {
			t.Fatal(err)
		}
		defer func() {
			if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
				t.Errorf("signal API: %v", err)
			}
			if err := cmd.Wait(); err != nil {
				t.Errorf("stop API: %v", err)
			}
		}()

		line, err := bufio.NewReader(out).ReadString('\n')
		if err != nil {
			t.Fatalf("read API address: %v", err)
		}
		baseURL, ok := strings.CutPrefix(strings.TrimSpace(line), "Listening on ")
		if !ok {
			t.Fatalf("unexpected startup message %q", line)
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/api/health", nil)
		if err != nil {
			t.Fatal(err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Error(err)
			}
		}()
		if resp.StatusCode != http.StatusOK || resp.Header.Get("Content-Type") != "application/json" {
			t.Fatalf("GET /api/health returned %s, %q", resp.Status, resp.Header.Get("Content-Type"))
		}
		var body struct{ Status string }
		if err := json.NewDecoder(resp.Body).Decode(&body); err != nil || body.Status != "ok" {
			t.Fatalf("decode health response: status %q, error %v", body.Status, err)
		}
	})

	t.Run("database unavailable at startup", func(t *testing.T) {
		t.Setenv("DATABASE_URL", "postgres://booking:booking@127.0.0.1:1/booking?sslmode=disable")
		ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
		defer cancel()
		output, err := exec.CommandContext(ctx, binary).CombinedOutput()
		var exitErr *exec.ExitError
		if !errors.As(err, &exitErr) || exitErr.ExitCode() != 1 || !strings.Contains(string(output), "database ping") {
			t.Fatalf("startup with unavailable database: %v\n%s", err, output)
		}
	})
}
