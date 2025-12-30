package tests

// NOTE: AI generated because im lazy

import (
	"reflect"
	"testing"
	"time"

	"github.com/tristnaja/taski/internal/io"
)

func TestAddTask(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name        string
		initialDB   io.Database
		taskToAdd   io.Task
		expectedDB  io.Database
		expectError bool
	}{
		{
			name:      "add to empty db",
			initialDB: io.Database{Tasks: []io.Task{}},
			taskToAdd: io.Task{Title: "New Task", Description: "A description"},
			expectedDB: io.Database{
				Size: 1,
				Tasks: []io.Task{
					{ID: 0, Title: "New Task", Description: "A description"},
				},
			},
			expectError: false,
		},
		{
			name: "add to existing db",
			initialDB: io.Database{
				Size:  1,
				Tasks: []io.Task{{ID: 0, Title: "First Task", Description: "Desc 1"}},
			},
			taskToAdd: io.Task{Title: "Second Task", Description: "Desc 2"},
			expectedDB: io.Database{
				Size: 2,
				Tasks: []io.Task{
					{ID: 0, Title: "First Task", Description: "Desc 1"},
					{ID: 1, Title: "Second Task", Description: "Desc 2"},
				},
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dbFile := setupTestDB(t, tc.initialDB)

			err := io.AddTask(tc.taskToAdd, dbFile)

			if (err != nil) != tc.expectError {
				t.Errorf("AddTask() error = %v, expectError %v", err, tc.expectError)
				return
			}

			if !tc.expectError {
				resultDB := readTestDB(t, dbFile)
				// We don't care about the date, so we ignore it in comparison
				for i := range resultDB.Tasks {
					resultDB.Tasks[i].Date = time.Time{}
				}
				if !reflect.DeepEqual(resultDB.Size, tc.expectedDB.Size) {
					t.Errorf("AddTask() got size = %v, want %v", resultDB.Size, tc.expectedDB.Size)
				}

				if !reflect.DeepEqual(resultDB.Tasks, tc.expectedDB.Tasks) {
					t.Errorf("AddTask() got tasks = %+v, want %+v", resultDB.Tasks, tc.expectedDB.Tasks)
				}
			}
		})
	}
}

func TestReadTask(t *testing.T) {
	t.Parallel()
	deletedAt := time.Now()

	testCases := []struct {
		name        string
		initialDB   io.Database
		expectedDB  io.Database
		expectError bool
	}{
		{
			name: "read only active tasks",
			initialDB: io.Database{
				Size: 2,
				Tasks: []io.Task{
					{ID: 0, Title: "Active Task", IsDeleted: false},
					{ID: 1, Title: "Deleted Task", IsDeleted: true, DeletedAt: &deletedAt},
				},
			},
			expectedDB: io.Database{
				Size:  2, // ReadTask doesn't filter size
				Tasks: []io.Task{{ID: 0, Title: "Active Task", IsDeleted: false}},
			},
			expectError: false,
		},
		{
			name:      "read from empty db",
			initialDB: io.Database{Tasks: []io.Task{}},
			expectedDB: io.Database{
				Tasks: []io.Task{},
			},
			expectError: false,
		},
		{
			name: "read from db with only deleted tasks",
			initialDB: io.Database{
				Size: 1,
				Tasks: []io.Task{
					{ID: 0, Title: "Deleted Task", IsDeleted: true, DeletedAt: &deletedAt},
				},
			},
			expectedDB: io.Database{
				Size:  1,
				Tasks: []io.Task{},
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dbFile := setupTestDB(t, tc.initialDB)

			resultDB, err := io.ReadTask(dbFile)

			if (err != nil) != tc.expectError {
				t.Errorf("ReadTask() error = %v, expectError %v", err, tc.expectError)
				return
			}

			if !tc.expectError {
				// Normalize slices for comparison (nil vs empty slice)
				if len(resultDB.Tasks) == 0 {
					resultDB.Tasks = []io.Task{}
				}
				if !reflect.DeepEqual(resultDB, tc.expectedDB) {
					t.Errorf("ReadTask() got = %+v, want %+v", resultDB, tc.expectedDB)
				}
			}
		})
	}
}

func TestChangeTask(t *testing.T) {
	t.Parallel()
	initialTasks := []io.Task{{ID: 0, Title: "Original Title", Description: "Original Desc"}}

	testCases := []struct {
		name           string
		taskIndex      int
		newTitle       string
		newDescription string
		expectedTitle  string
		expectedDesc   string
		expectError    bool
	}{
		{
			name:           "change both title and description",
			taskIndex:      0,
			newTitle:       "New Title",
			newDescription: "New Desc",
			expectedTitle:  "New Title",
			expectedDesc:   "New Desc",
			expectError:    false,
		},
		{
			name:           "change only title",
			taskIndex:      0,
			newTitle:       "Only Title Changed",
			newDescription: "",
			expectedTitle:  "Only Title Changed",
			expectedDesc:   "Original Desc",
			expectError:    false,
		},
		{
			name:           "change only description",
			taskIndex:      0,
			newTitle:       "",
			newDescription: "Only Desc Changed",
			expectedTitle:  "Original Title",
			expectedDesc:   "Only Desc Changed",
			expectError:    false,
		},
		{
			name:           "change nothing",
			taskIndex:      0,
			newTitle:       "",
			newDescription: "",
			expectError:    true,
		},
		{
			name:        "invalid negative index",
			taskIndex:   -1,
			expectError: true,
		},
		{
			name:        "invalid out of bounds index",
			taskIndex:   len(initialTasks),
			newTitle:    "trigger error",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			initialDB := io.Database{Size: len(initialTasks), Tasks: initialTasks}
			dbFile := setupTestDB(t, initialDB)

			err := io.ChangeTask(dbFile, tc.taskIndex, tc.newTitle, tc.newDescription)

			if (err != nil) != tc.expectError {
				t.Errorf("ChangeTask() error = %v, expectError %v", err, tc.expectError)
				return
			}

			if !tc.expectError {
				db := readTestDB(t, dbFile)
				if db.Tasks[tc.taskIndex].Title != tc.expectedTitle {
					t.Errorf("ChangeTask() title got = %q, want %q", db.Tasks[tc.taskIndex].Title, tc.expectedTitle)
				}
				if db.Tasks[tc.taskIndex].Description != tc.expectedDesc {
					t.Errorf("ChangeTask() description got = %q, want %q", db.Tasks[tc.taskIndex].Description, tc.expectedDesc)
				}
				if db.Tasks[tc.taskIndex].Date.IsZero() {
					t.Errorf("ChangeTask() date was not updated")
				}
			}
		})
	}
}

func TestRemoveTask(t *testing.T) {
	t.Parallel()
	initialDB := io.Database{
		Size:  1,
		Tasks: []io.Task{{ID: 0, Title: "Task to delete", IsDeleted: false}},
	}

	testCases := []struct {
		name        string
		taskIndex   int
		expectError bool
	}{
		{
			name:        "remove valid task",
			taskIndex:   0,
			expectError: false,
		},
		{
			name:        "invalid out of bounds index",
			taskIndex:   1,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dbFile := setupTestDB(t, initialDB)

			err := io.RemoveTask(dbFile, tc.taskIndex)

			if (err != nil) != tc.expectError {
				t.Errorf("RemoveTask() error = %v, expectError %v", err, tc.expectError)
				return
			}

			if !tc.expectError {
				db := readTestDB(t, dbFile)
				if !db.Tasks[tc.taskIndex].IsDeleted {
					t.Errorf("RemoveTask() task was not marked as deleted")
				}
				if db.Tasks[tc.taskIndex].DeletedAt == nil {
					t.Errorf("RemoveTask() DeletedAt was not set")
				}
				if db.Size != initialDB.Size-1 {
					t.Errorf("RemoveTask() DB size was not decremented, got %d, want %d", db.Size, initialDB.Size-1)
				}
			}
		})
	}
}

func TestRestoreTask(t *testing.T) {
	t.Parallel()
	deletedAt := time.Now()
	initialDB := io.Database{
		Size: 1, // Size is active tasks
		Tasks: []io.Task{
			{ID: 0, Title: "Deleted Task", IsDeleted: true, DeletedAt: &deletedAt},
			{ID: 1, Title: "Active Task", IsDeleted: false},
		},
	}

	testCases := []struct {
		name              string
		taskIndex         int
		initialSize       int
		expectedSize      int
		shouldBeRestored  bool
		expectError       bool
		finalIsDeletedState bool
	}{
		{
			name:              "restore deleted task",
			taskIndex:         0,
			initialSize:       1,
			expectedSize:      2,
			shouldBeRestored:  true,
			expectError:       false,
			finalIsDeletedState: false,
		},
		{
			name:              "restore active task (no change)",
			taskIndex:         1,
			initialSize:       1,
			expectedSize:      1, // Should not change
			shouldBeRestored:  false,
			expectError:       false,
			finalIsDeletedState: false,
		},
		{
			name:        "invalid out of bounds index",
			taskIndex:   2,
			initialSize: 1,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dbFile := setupTestDB(t, initialDB)

			err := io.RestoreTask(dbFile, tc.taskIndex)

			if (err != nil) != tc.expectError {
				t.Errorf("RestoreTask() error = %v, expectError %v", err, tc.expectError)
				return
			}

			if !tc.expectError {
				db := readTestDB(t, dbFile)
				if db.Tasks[tc.taskIndex].IsDeleted != tc.finalIsDeletedState {
					t.Errorf("RestoreTask() final deleted state got = %v, want %v", db.Tasks[tc.taskIndex].IsDeleted, tc.finalIsDeletedState)
				}
				if db.Size != tc.expectedSize {
					t.Errorf("RestoreTask() DB size got = %d, want %d", db.Size, tc.expectedSize)
				}
				if tc.shouldBeRestored {
					if db.Tasks[tc.taskIndex].DeletedAt != nil {
						t.Errorf("RestoreTask() DeletedAt was not cleared")
					}
				}
			}
		})
	}
}

func TestCleanUp(t *testing.T) {
	t.Parallel()
	now := time.Now()
	deletedOld := now.Add(-2 * time.Hour)
	deletedRecent := now.Add(-30 * time.Minute)

	initialDB := io.Database{
		Size: 3,
		Tasks: []io.Task{
			{ID: 0, Title: "Active Task"},
			{ID: 1, Title: "Recently Deleted", IsDeleted: true, DeletedAt: &deletedRecent},
			{ID: 2, Title: "Old Deleted", IsDeleted: true, DeletedAt: &deletedOld},
		},
	}

	testCases := []struct {
		name            string
		retention       time.Duration
		expectedTaskIDs []int
		expectError     bool
	}{
		{
			name:            "cleanup old tasks",
			retention:       time.Hour,
			expectedTaskIDs: []int{0, 1}, // Active and recently deleted are kept
			expectError:     false,
		},
		{
			name:            "cleanup with zero retention",
			retention:       0,
			expectedTaskIDs: []int{0}, // Only active tasks kept
			expectError:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dbFile := setupTestDB(t, initialDB)

			err := io.CleanUp(dbFile, tc.retention)

			if (err != nil) != tc.expectError {
				t.Errorf("CleanUp() error = %v, expectError %v", err, tc.expectError)
				return
			}

			if !tc.expectError {
				db := readTestDB(t, dbFile)
				var resultingIDs []int
				for _, task := range db.Tasks {
					resultingIDs = append(resultingIDs, task.ID)
				}
				if !reflect.DeepEqual(resultingIDs, tc.expectedTaskIDs) {
					t.Errorf("CleanUp() resulting task IDs got = %v, want %v", resultingIDs, tc.expectedTaskIDs)
				}
			}
		})
	}
}

func TestRestoreAll(t *testing.T) {
	t.Parallel()
	deletedAt := time.Now()
	initialDB := io.Database{
		Size: 1, // 1 active task
		Tasks: []io.Task{
			{ID: 0, Title: "Deleted 1", IsDeleted: true, DeletedAt: &deletedAt},
			{ID: 1, Title: "Active"},
			{ID: 2, Title: "Deleted 2", IsDeleted: true, DeletedAt: &deletedAt},
		},
	}

	dbFile := setupTestDB(t, initialDB)
	err := io.RestoreAll(dbFile)

	if err != nil {
		t.Fatalf("RestoreAll() returned an unexpected error: %v", err)
	}

	db := readTestDB(t, dbFile)
	if db.Size != len(db.Tasks) {
		t.Errorf("RestoreAll() size got = %d, want %d", db.Size, len(db.Tasks))
	}

	for _, task := range db.Tasks {
		if task.IsDeleted {
			t.Errorf("RestoreAll() found a task that was not restored: ID %d", task.ID)
		}
		if task.DeletedAt != nil {
			t.Errorf("RestoreAll() found a task where DeletedAt was not cleared: ID %d", task.ID)
		}
	}
}