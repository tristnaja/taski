package tests

import (
	"strings"
	"testing"
	"time"

	"github.com/tristnaja/taski/internal/io"
)

func TestRunRestore(t *testing.T) {
	deletedAt := time.Now()
	testCases := []struct {
		name                 string
		args                 []string
		initialDB            io.Database
		expectedStdout       string
		expectedStderr       string
		allShouldBeRestored  bool
		taskShouldBeRestored bool
		targetIndex          int
		expectedExitCode     int
	}{
		{
			name: "restore single task",
			args: []string{"-i", "0"},
			initialDB: io.Database{
				Size: 0,
				Tasks: []io.Task{
					{ID: 0, IsDeleted: true, DeletedAt: &deletedAt},
					{ID: 1, IsDeleted: true, DeletedAt: &deletedAt},
				},
			},
			expectedStdout:       "Task Restored",
			taskShouldBeRestored: true,
			targetIndex:          0,
			expectedExitCode:     0,
		},
		{
			name: "restore all tasks",
			args: []string{"--all"},
			initialDB: io.Database{
				Size: 0,
				Tasks: []io.Task{
					{ID: 0, IsDeleted: true, DeletedAt: &deletedAt},
					{ID: 1, IsDeleted: true, DeletedAt: &deletedAt},
				},
			},
			expectedStdout:      "All Tasks in Trash is Restored",
			allShouldBeRestored: true,
			expectedExitCode:    0,
		},
		{
			name:             "no flags",
			args:             []string{},
			initialDB:        io.Database{},
			expectedStderr:   "Usage of restore:|unfilled arguments",
			expectedExitCode: 1,
		},
		{
			name:             "all and index flag",
			args:             []string{"-a", "-i", "0"},
			initialDB:        io.Database{},
			expectedStderr:   "When restoring all, you do not need to input an index",
			expectedExitCode: 1,
		},
		{
			name:             "invalid index",
			args:             []string{"-i", "99"},
			initialDB:        io.Database{Tasks: []io.Task{{ID: 0, IsDeleted: true, DeletedAt: &deletedAt}}},
			expectedStderr:   "restoring task: restoring task: invalid index 99: out of bounds",
			expectedExitCode: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dbFile := setupTestDB(t, tc.initialDB)

			stdout, stderr, exitCode := runTestCommand(t, "RunRestore", tc.args, dbFile)

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

			if tc.allShouldBeRestored {
				db := readTestDB(t, dbFile)
				for _, task := range db.Tasks {
					if task.IsDeleted {
						t.Errorf("task with ID %d was not restored", task.ID)
					}
				}
			}

			if tc.taskShouldBeRestored {
				db := readTestDB(t, dbFile)
				if db.Tasks[tc.targetIndex].IsDeleted {
					t.Errorf("task with index %d was not restored", tc.targetIndex)
				}
			}
		})
	}
}