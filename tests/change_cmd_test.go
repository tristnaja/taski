package tests

import (
	"strings"
	"testing"

	"github.com/tristnaja/taski/internal/io"
)

func TestRunChange(t *testing.T) {
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
			name: "change task successfully",
			args: []string{"-i", "0", "-t", "Changed Task", "-d", "Changed description"},
			initialDB: io.Database{
				Size:  1,
				Tasks: []io.Task{{ID: 0, Title: "Original Task", Description: "Original Description"}},
			},
			expectedStdout:   "Changed Task:",
			expectedInDB:     "Changed Task",
			expectedExitCode: 0,
		},
		{
			name: "missing index",
			args: []string{"-t", "some title"},
			initialDB: io.Database{
				Size:  1,
				Tasks: []io.Task{{ID: 0, Title: "Original Task", Description: "Original Description"}},
			},
			expectedStderr:   "Usage of change:|unfilled arguments",
			expectedExitCode: 1,
		},
		{
			name: "no new values",
			args: []string{"-i", "0"},
			initialDB: io.Database{
				Size:  1,
				Tasks: []io.Task{{ID: 0, Title: "Original Task", Description: "Original Description"}},
			},
			expectedStdout:   "",
			expectedStderr:   "changing task: No value is changed",
			expectedExitCode: 1,
		},
		{
			name: "invalid index",
			args: []string{"-i", "99", "-t", "some title"},
			initialDB: io.Database{
				Size:  1,
				Tasks: []io.Task{{ID: 0, Title: "Original Task", Description: "Original Description"}},
			},
			expectedStderr:   "changing task: invalid index 99: out of bounds",
			expectedExitCode: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dbFile := setupTestDB(t, tc.initialDB)

			stdout, stderr, exitCode := runTestCommand(t, "RunChange", tc.args, dbFile)

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