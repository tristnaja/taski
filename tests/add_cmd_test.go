package tests

// NOTE: AI generated because im lazy

import (
	"strings"
	"testing"

	"github.com/tristnaja/taski/internal/io"
)

func TestRunAdd(t *testing.T) {
	testCases := []struct {
		name             string
		args             []string
		initialDB        io.Database
		expectedStdout   string
		expectedStderr   string
		expectedInDB     string
		expectedExitCode int
	}{
		{
			name:             "add task successfully",
			args:             []string{"-t", "New Task", "-d", "A description"},
			initialDB:        io.Database{Tasks: []io.Task{}},
			expectedStdout:   "Added New Task:",
			expectedInDB:     "New Task",
			expectedExitCode: 0,
		},
		{
			name:             "missing title",
			args:             []string{"-d", "A description"},
			initialDB:        io.Database{Tasks: []io.Task{}},
			expectedStdout:   "",
			expectedStderr:   "Usage of add:|unfilled arguments",
			expectedExitCode: 1,
		},
		{
			name:             "missing description",
			args:             []string{"-t", "New Task"},
			initialDB:        io.Database{Tasks: []io.Task{}},
			expectedStdout:   "",
			expectedStderr:   "Usage of add:|unfilled arguments",
			expectedExitCode: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dbFile := setupTestDB(t, tc.initialDB)

			stdout, stderr, exitCode := runTestCommand(t, "RunAdd", tc.args, dbFile)

			if tc.expectedStdout != "" && !strings.Contains(stdout, tc.expectedStdout) {
				t.Errorf("expected stdout to contain %q, got %q", tc.expectedStdout, stdout)
			}

			if tc.expectedStderr != "" {
				for _, expected := range strings.Split(tc.expectedStderr, "|") {
					if !strings.Contains(stderr, expected) {
						t.Errorf("expected stderr to contain %q, got %q", expected, stderr)
					}
				}
			}

			if exitCode != tc.expectedExitCode {
				t.Errorf("expected exit code %d, got %d", tc.expectedExitCode, exitCode)
			}

			if tc.expectedInDB != "" {
				db := readTestDB(t, dbFile)
				found := false
				for _, task := range db.Tasks {
					if task.Title == tc.expectedInDB {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected task '%s' was not found in the database", tc.expectedInDB)
				}
			}
		})
	}
}