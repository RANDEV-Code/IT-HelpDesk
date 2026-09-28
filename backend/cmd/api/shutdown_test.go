package main

import (
	"bytes"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

// TestGracefulShutdown membuild binary api, menjalankannya dengan konfigurasi
// aman (DB tidak perlu hidup), mengirim sinyal berhenti (SIGTERM di Unix,
// CTRL_BREAK→SIGTERM di Windows), lalu memastikan proses keluar bersih
// (exit code 0) dalam ≤15 detik tanpa panic (acceptance criteria TASK-002).
func TestGracefulShutdown(t *testing.T) {
	if testing.Short() {
		t.Skip("tes shutdown dilewati pada mode -short")
	}

	dir := t.TempDir()
	exe := filepath.Join(dir, "api")
	if runtime.GOOS == "windows" {
		exe += ".exe"
	}

	build := exec.Command("go", "build", "-o", exe, ".")
	build.Stderr = os.Stderr
	if err := build.Run(); err != nil {
		t.Fatalf("go build: %v", err)
	}

	var stdout, stderr bytes.Buffer
	cmd := exec.Command(exe)
	// cwd di temp dir → tidak menemukan .env repo; seluruh konfigurasi eksplisit.
	cmd.Dir = dir
	cmd.Env = []string{
		"APP_ENV=development",
		"APP_ORIGIN=http://localhost:5173",
		"HTTP_ADDR=127.0.0.1:18099",
		"DATABASE_URL=postgres://shutdown:shutdown@127.0.0.1:1/shutdown?sslmode=disable",
		"STORAGE_ROOT=" + filepath.Join(dir, "storage"),
		"LOG_LEVEL=info",
		"PATH=" + os.Getenv("PATH"),
		"HOME=" + os.Getenv("HOME"),
		"USERPROFILE=" + os.Getenv("USERPROFILE"),
		"GOMAXPROCS=" + os.Getenv("GOMAXPROCS"),
	}
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	prepareCmd(cmd) // SysProcAttr spesifik platform

	if err := cmd.Start(); err != nil {
		t.Fatalf("start api: %v", err)
	}

	waited := make(chan error, 1)
	go func() { waited <- cmd.Wait() }()

	// tunggu /health/live menjawab (maks 10 detik)
	liveURL := "http://127.0.0.1:18099/health/live"
	deadline := time.Now().Add(10 * time.Second)
	var lastErr error
	up := false
	for time.Now().Before(deadline) {
		resp, err := http.Get(liveURL)
		if err == nil {
			resp.Body.Close()
			up = true
			break
		}
		lastErr = err
		time.Sleep(250 * time.Millisecond)
	}
	if !up {
		t.Fatalf("api tidak pernah siap: %v\nstderr: %s", lastErr, stderr.String())
	}

	start := time.Now()
	if err := sendShutdownSignal(cmd); err != nil {
		t.Fatalf("kirim sinyal shutdown: %v", err)
	}

	select {
	case err := <-waited:
		elapsed := time.Since(start)
		if err != nil {
			t.Fatalf("proses keluar dengan error (elapsed %v): %v\nstdout: %s\nstderr: %s",
				elapsed.Round(time.Millisecond), err, stdout.String(), stderr.String())
		}
		if elapsed > 15*time.Second {
			t.Errorf("shutdown memakan %v, ingin ≤15 detik", elapsed)
		}
		if out := stdout.String(); strings.Contains(out, "panic") {
			t.Errorf("stdout memuat panic:\n%s", out)
		}
		if out := stdout.String(); !strings.Contains(out, "api berhenti bersih") {
			t.Errorf("log penutup bersih tidak ditemukan; stdout:\n%s", out)
		}
		t.Logf("graceful shutdown OK dalam %v", elapsed.Round(time.Millisecond))
	case <-time.After(20 * time.Second):
		_ = cmd.Process.Kill()
		t.Fatal("proses tidak berhenti dalam 20 detik; dibunuh paksa")
	}
}
