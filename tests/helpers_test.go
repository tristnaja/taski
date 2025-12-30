package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tristnaja/taski/app/cmd"
	"github.com/tristnaja/taski/internal/io"
)

// setupTestDB creates a temporary database file for testing.
func setupTestDB(t *testing.T, initialData io.Database) string {
	t.Helper()

	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test_db.json")

	// Create the file
	file, err := os.Create(tempFile)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	// Write initial data if provided
	if len(initialData.Tasks) > 0 {
		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "	")
		if err := encoder.Encode(initialData); err != nil {
			t.Fatalf("Failed to write initial data to temp file: %v", err)
		}
	}
	file.Close()

	return tempFile
}

// readTestDB reads the content of a test database file.
func readTestDB(t *testing.T, filePath string) io.Database {
	t.Helper()
	var db io.Database

	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return io.Database{Tasks: []io.Task{}}
		}
		t.Fatalf("Failed to open test db file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&db); err != nil {
		// EOF is fine for an empty file
		if err.Error() == "EOF" {
			return io.Database{Tasks: []io.Task{}}
		}
		t.Fatalf("Failed to decode test db file: %v", err)
	}

	return db
}

// runTestCommand executes the target Run function in a subprocess to test code that calls os.Exit.
// It returns stdout, stderr, and the exit code.
func runTestCommand(t *testing.T, funcName string, args []string, dbFile string) (string, string, int) {
	// Find the test function name to re-run in the subprocess.
	testName := t.Name()

	cmd := exec.Command(os.Args[0], "-test.run", "^"+testName+"$")
	cmd.Env = append(os.Environ(),
		"GO_TEST_SUBPROCESS=1",
		"TEST_FUNC_NAME="+funcName,
		"TEST_DB_FILE="+dbFile,
		"TEST_ARGS="+strings.Join(args, "|"),
	)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			t.Fatalf("cmd.Run() failed with non-exit error: %v", err)
		}
	}

	return stdout.String(), stderr.String(), exitCode
}

// TestMain is a test wrapper that checks for a subprocess indicator.
// If the indicator is present, it runs the target function and exits.
func TestMain(m *testing.M) {
	if os.Getenv("GO_TEST_SUBPROCESS") == "1" {
		funcName := os.Getenv("TEST_FUNC_NAME")
		dbFile := os.Getenv("TEST_DB_FILE")
		args := strings.Split(os.Getenv("TEST_ARGS"), "|")
		if len(args) == 1 && args[0] == "" {
			args = []string{}
		}

		var err error
		switch funcName {
		case "RunAdd":
			err = cmd.RunAdd(args, dbFile)
		case "RunChange":
			err = cmd.RunChange(args, dbFile)
		case "RunDelete":
			err = cmd.RunDelete(args, dbFile)
		case "RunRestore":
			err = cmd.RunRestore(args, dbFile)
		case "RunView":
			err = cmd.RunView(args, dbFile)
		}

		if err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}

		os.Exit(0)
	}

	// Run the normal tests
	os.Exit(m.Run())
}
