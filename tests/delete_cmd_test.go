package tests

import (
	"strings"
	"testing"

	"github.com/tristnaja/taski/internal/io"
)

func TestRunDelete(t *testing.T) {
	testCases := []struct {
		name                string
		args                []string
		initialDB           io.Database
		expectedStdout      string
		expectedStderr      string
		taskShouldBeDeleted bool
		expectedExitCode    int
	}{
		{
			name: "delete task successfully",
			args: []string{"-i", "0"},
			initialDB: io.Database{
				Size:  1,
				Tasks: []io.Task{{ID: 0, Title: "Task to delete"}},
			},
			expectedStdout:      "Deleted Task:",
			taskShouldBeDeleted: true,
			expectedExitCode:    0,
		},
		{
			name:             "missing index",
			args:             []string{},
			initialDB:        io.Database{},
			expectedStderr:   "Usage of delete:|unfilled arguments",
			expectedExitCode: 1,
		},
		{
			name: "invalid index",
			args: []string{"-i", "99"},
			initialDB: io.Database{
				Size:  1,
				Tasks: []io.Task{{ID: 0, Title: "A Task"}},
			},
			expectedStderr:   "deleting task: deleting task: invalid index 99: out of bounds",
			expectedExitCode: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dbFile := setupTestDB(t, tc.initialDB)

			stdout, stderr, exitCode := runTestCommand(t, "RunDelete", tc.args, dbFile)

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

			if tc.taskShouldBeDeleted {
				db := readTestDB(t, dbFile)
				if len(db.Tasks) > 0 && !db.Tasks[0].IsDeleted {
					t.Errorf("task was not marked as deleted")
				}
			}
		})
	}
}