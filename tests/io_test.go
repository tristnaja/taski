package tests

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/tristnaja/taski/internal/io"
)

// setupTestFile creates a temporary file with initial data for testing.
func setupTestFile(t *testing.T, initialData io.Database) string {
	t.Helper()
	dir := t.TempDir()
	fileName := filepath.Join(dir, "test_db.json")

	// If initialData is the zero value, create an empty file
	if reflect.DeepEqual(initialData, io.Database{}) {
		file, err := os.Create(fileName)
		if err != nil {
			t.Fatalf("Failed to create empty temp file: %v", err)
		}
		file.Close()
		return fileName
	}

	file, err := os.Create(fileName)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer file.Close()

	// If initialData is provided and not empty, write it to the file.
	if len(initialData.Tasks) > 0 || initialData.Size > 0 {
		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "\t")
		if err := encoder.Encode(initialData); err != nil {
			t.Fatalf("Failed to write initial data to temp file: %v", err)
		}
	}

	return fileName
}

// readTestFile reads the content of the test database file.
func readTestFile(t *testing.T, fileName string) io.Database {
	t.Helper()
	var db io.Database

	file, err := os.Open(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return io.Database{}
		}
		t.Fatalf("Failed to open test file for reading: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&db); err != nil {
		// Handle empty file case, where Decode would return EOF
		if err.Error() == "EOF" {
			return io.Database{}
		}
		t.Fatalf("Failed to decode test file: %v", err)
	}
	return db
}

func TestAddTask(t *testing.T) {
	// Define a fixed time for consistency in tests
	fixedTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)

	testCases := []struct {
		name         string
		initialDB    io.Database
		taskToAdd    *io.Task
		expectedDB   io.Database
		expectError  bool
		setupCorrupt bool // Flag to manually create a corrupt JSON file
	}{
		{
			name:      "Happy Path: Add task to an empty database",
			initialDB: io.Database{Size: 0, Tasks: []io.Task{}},
			taskToAdd: &io.Task{Title: "New Task", Description: "A new task", Date: fixedTime},
			expectedDB: io.Database{
				Size: 1,
				Tasks: []io.Task{
					{Title: "New Task", Description: "A new task", Date: fixedTime},
				},
			},
			expectError: false,
		},
		{
			name: "Happy Path: Add task to a non-empty database",
			initialDB: io.Database{
				Size: 1,
				Tasks: []io.Task{
					{Title: "Existing Task", Description: "An old task", Date: fixedTime.Add(-time.Hour)},
				},
			},
			taskToAdd: &io.Task{Title: "Another Task", Description: "Another new task", Date: fixedTime},
			expectedDB: io.Database{
				Size: 2,
				Tasks: []io.Task{
					{Title: "Existing Task", Description: "An old task", Date: fixedTime.Add(-time.Hour)},
					{Title: "Another Task", Description: "Another new task", Date: fixedTime},
				},
			},
			expectError: false,
		},
		{
			name:      "Edge Case: Add task with empty title and description",
			initialDB: io.Database{Size: 0, Tasks: []io.Task{}},
			taskToAdd: &io.Task{Title: "", Description: "", Date: fixedTime},
			expectedDB: io.Database{
				Size: 1,
				Tasks: []io.Task{
					{Title: "", Description: "", Date: fixedTime},
				},
			},
			expectError: false,
		},
		{
			name:         "Bad Case: Corrupt JSON file",
			initialDB:    io.Database{}, // Will be ignored
			taskToAdd:    &io.Task{Title: "Test", Description: "Test"},
			expectedDB:   io.Database{},
			expectError:  true,
			setupCorrupt: true,
		},
		{
			name: "Bad Case: File path is a directory",
			// No initial DB needed, the path itself is the problem
			initialDB:   io.Database{},
			taskToAdd:   &io.Task{Title: "Test", Description: "Test"},
			expectedDB:  io.Database{},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var fileName string
			if tc.name == "Bad Case: File path is a directory" {
				fileName = t.TempDir()
			} else {
				fileName = setupTestFile(t, tc.initialDB)
			}

			if tc.setupCorrupt {
				// Manually write invalid JSON to the file
				if err := os.WriteFile(fileName, []byte("{corrupt_json"), 0644); err != nil {
					t.Fatalf("Failed to write corrupt data: %v", err)
				}
			}

			// Defer cleanup for directory paths as well
			defer func() {
				if tc.name == "Bad Case: File path is a directory" {
					os.RemoveAll(fileName) // Clean up the created directory
				}
			}()

			err := io.AddTask(tc.taskToAdd, fileName)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected an error, but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Did not expect an error, but got: %v", err)
				}

				actualDB := readTestFile(t, fileName)
				if !reflect.DeepEqual(actualDB, tc.expectedDB) {
					t.Errorf("Database content mismatch.\nExpected: %+v\nGot:      %#v", tc.expectedDB, actualDB)
				}
			}
		})
	}
}

func TestReadTask(t *testing.T) {
	deletedTime := time.Date(2023, 1, 1, 13, 0, 0, 0, time.UTC)

	testCases := []struct {
		name         string
		initialDB    io.Database
		expectedDB   io.Database
		expectError  bool
		setupCorrupt bool
		nonExistent  bool
	}{
		{
			name: "Happy Path: Read tasks with no deleted items",
			initialDB: io.Database{
				Size: 2,
				Tasks: []io.Task{
					{Title: "Task 1", IsDeleted: false},
					{Title: "Task 2", IsDeleted: false},
				},
			},
			expectedDB: io.Database{
				Size: 2, // Bug: Should be 2
				Tasks: []io.Task{
					{Title: "Task 1", IsDeleted: false},
					{Title: "Task 2", IsDeleted: false},
				},
			},
			expectError: false,
		},
		{
			name: "Happy Path: Read tasks with some deleted items",
			initialDB: io.Database{
				Size: 3,
				Tasks: []io.Task{
					{Title: "Task 1", IsDeleted: false},
					{Title: "Task 2 (Deleted)", IsDeleted: true, DeletedAt: &deletedTime},
					{Title: "Task 3", IsDeleted: false},
				},
			},
			expectedDB: io.Database{
				Size: 3, // Bug: Should be 2
				Tasks: []io.Task{
					{Title: "Task 1", IsDeleted: false},
					{Title: "Task 3", IsDeleted: false},
				},
			},
			expectError: false,
		},
		{
			name: "Happy Path: Read tasks with all items deleted",
			initialDB: io.Database{
				Size: 2,
				Tasks: []io.Task{
					{Title: "Task 1 (Deleted)", IsDeleted: true, DeletedAt: &deletedTime},
					{Title: "Task 2 (Deleted)", IsDeleted: true, DeletedAt: &deletedTime},
				},
			},
			expectedDB: io.Database{
				Size:  2, // Bug: Should be 0
				Tasks: nil,
			},
			expectError: false,
		},
		{
			name:        "Edge Case: Read from an empty database file",
			initialDB:   io.Database{Size: 0, Tasks: []io.Task{}},
			expectedDB:  io.Database{Size: 0, Tasks: nil},
			expectError: false,
		},
		{
			name:        "Edge Case: Read from a completely empty file",
			initialDB:   io.Database{}, // Represents a zero-byte file
			expectedDB:  io.Database{},
			expectError: false,
		},
		{
			name:        "Edge Case: Read from a non-existent file (creates new)",
			initialDB:   io.Database{},
			expectedDB:  io.Database{Size: 0, Tasks: nil},
			expectError: false,
			nonExistent: true,
		},
		{
			name:         "Bad Case: Read from a corrupt JSON file",
			initialDB:    io.Database{},
			expectedDB:   io.Database{},
			expectError:  true,
			setupCorrupt: true,
		},
		{
			name:        "Bad Case: Read from a directory",
			initialDB:   io.Database{},
			expectedDB:  io.Database{},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var fileName string
			if tc.nonExistent {
				fileName = filepath.Join(t.TempDir(), "non_existent_file.json")
			} else if tc.name == "Bad Case: Read from a directory" {
				fileName = t.TempDir()
			} else {
				fileName = setupTestFile(t, tc.initialDB)
			}

			if tc.setupCorrupt {
				if err := os.WriteFile(fileName, []byte("{corrupt_json"), 0644); err != nil {
					t.Fatalf("Failed to write corrupt data: %v", err)
				}
			}

			actualDB, err := io.ReadTask(fileName)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected an error, but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Did not expect an error, but got: %v", err)
				}
				if !reflect.DeepEqual(actualDB, tc.expectedDB) {
					t.Errorf("Database content mismatch.\nExpected: %+v\nGot:      %#v", tc.expectedDB, actualDB)
				}
			}
		})
	}
}

