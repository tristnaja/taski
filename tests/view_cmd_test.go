package tests

// NOTE: AI generated because im lazy

import (
	"strings"
	"testing"
	"time"

	"github.com/tristnaja/taski/internal/io"
)

func TestRunView(t *testing.T) {
	testCases := []struct {
		name           string
		args           []string
		initialDB      io.Database
		expectedStdout []string
		expectedExitCode int
	}{
		{
			name: "view with tasks",
			args: []string{},
			initialDB: io.Database{
				Size: 1,
				Tasks: []io.Task{{
					ID:          0,
					Title:       "My First Task",
					Description: "This is a task.",
					Date:        time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
				}},
			},
			expectedStdout: []string{
				"Here is your Tasks:",
				"1. My First Task",
				"index to target: 0",
				"Date: 01 Jan 2024, 12:00",
				"This is a task.",
			},
			expectedExitCode: 0,
		},
		{
			name: "view with no tasks",
			args: []string{},
			initialDB: io.Database{
				Tasks: []io.Task{},
			},
			expectedStdout: []string{
				"Here is your Tasks:",
				"You can Interact with your Tasks with:",
			},
			expectedExitCode: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dbFile := setupTestDB(t, tc.initialDB)

			stdout, _, exitCode := runTestCommand(t, "RunView", tc.args, dbFile)

			for _, expected := range tc.expectedStdout {
				if !strings.Contains(stdout, expected) {
					t.Errorf("expected stdout to contain %q, got %q", expected, stdout)
				}
			}

			if exitCode != tc.expectedExitCode {
				t.Errorf("expected exit code %d, got %d", tc.expectedExitCode, exitCode)
			}
		})
	}
}
