package testutil

import (
	"fmt"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"testing"
	"time"

	"final-whistle/backend/internal/db"
)

func CreateTestDatabase(t *testing.T) (string, string, func()) {
	t.Helper()

	adminURL := os.Getenv("TEST_DATABASE_ADMIN_URL")
	if adminURL == "" {
		currentUser, err := user.Current()
		if err != nil {
			t.Fatalf("get current user: %v", err)
		}
		adminURL = fmt.Sprintf("postgres://%s@localhost:5432/postgres?sslmode=disable", currentUser.Username)
	}

	adminDB, err := db.NewConnection(adminURL)
	if err != nil {
		t.Fatalf("connect admin database: %v", err)
	}

	dbName := fmt.Sprintf("final_whistle_test_%d", time.Now().UnixNano())
	if err := adminDB.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName)).Error; err != nil {
		_ = adminDB.Close()
		t.Fatalf("create test database: %v", err)
	}

	testURL := replaceDatabaseName(t, adminURL, dbName)
	testDB, err := db.NewConnection(testURL)
	if err != nil {
		_ = adminDB.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName)).Error
		_ = adminDB.Close()
		t.Fatalf("connect test database: %v", err)
	}

	for _, statement := range migrationStatements(t) {
		if strings.TrimSpace(statement) == "" {
			continue
		}
		if err := testDB.Exec(statement).Error; err != nil {
			_ = testDB.Close()
			_ = adminDB.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName)).Error
			_ = adminDB.Close()
			t.Fatalf("apply migration: %v\nSQL:\n%s", err, statement)
		}
	}
	_ = testDB.Close()

	cleanup := func() {
		_ = adminDB.Exec(
			"SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = ? AND pid <> pg_backend_pid()",
			dbName,
		).Error
		_ = adminDB.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName)).Error
		_ = adminDB.Close()
	}

	return dbName, testURL, cleanup
}

func replaceDatabaseName(t *testing.T, rawURL, dbName string) string {
	t.Helper()
	parsed, err := url.Parse(rawURL)
	if err != nil {
		t.Fatalf("parse database url: %v", err)
	}
	parsed.Path = "/" + dbName
	return parsed.String()
}

func migrationStatements(t *testing.T) []string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("resolve runtime caller failed")
	}
	migrationsDir := filepath.Join(filepath.Dir(file), "..", "..", "migrations")
	files, err := filepath.Glob(filepath.Join(migrationsDir, "*.sql"))
	if err != nil {
		t.Fatalf("glob migrations: %v", err)
	}
	sort.Strings(files)
	if len(files) == 0 {
		t.Fatalf("no migrations found in %s", migrationsDir)
	}

	statements := make([]string, 0, len(files))
	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			t.Fatalf("read migration %s: %v", file, err)
		}
		statements = append(statements, string(content))
	}
	return statements
}
