package io

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

type Task struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Date        time.Time  `json:"date"`
	IsDeleted   bool       `json:"is_deleted"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type Database struct {
	Size  int    `json:"size"`
	Tasks []Task `json:"tasks"`
}

func AddTask(task *Task, fileName string) error {
	db, err := readJSON(fileName)

	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	db.Tasks = append(db.Tasks, *task)
	db.Size++

	err = writeJSON(fileName, db)

	if err != nil {
		return fmt.Errorf("writing into file: %w", err)
	}

	return nil
}

func ReadTask(fileName string) (Database, error) {
	var filteredDB Database

	db, err := readJSON(fileName)

	if err != nil {
		return Database{}, fmt.Errorf("reading file: %w", err)
	}

	for _, task := range db.Tasks {
		if task.IsDeleted == false {
			filteredDB.Tasks = append(filteredDB.Tasks, task)
		}
	}

	filteredDB.Size = db.Size

	return filteredDB, nil
}

func ChangeTask(fileName string, taskIndex int, newTitle string, newDescription string) error {
	if taskIndex < 0 {
		return errors.New("Index cannot be < 0")
	}

	if newTitle == "" && newDescription == "" {
		return errors.New("No value is changed")
	}

	db, err := readJSON(fileName)

	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	if newTitle == "" {
		newTitle = db.Tasks[taskIndex].Title
	}

	if newDescription == "" {
		newDescription = db.Tasks[taskIndex].Description
	}

	db.Tasks[taskIndex].Title = newTitle
	db.Tasks[taskIndex].Description = newDescription
	db.Tasks[taskIndex].Date = time.Now()

	err = writeJSON(fileName, db)

	if err != nil {
		return fmt.Errorf("writing into file: %w", err)
	}

	return nil
}

func RemoveTask(fileName string, taskIndex int) error {
	err := softDelete(fileName, taskIndex)

	if err != nil {
		return fmt.Errorf("deleting task: %w", err)
	}

	return nil
}

func RestoreTask(fileName string, taskIndex int) error {
	err := restoreTask(fileName, taskIndex)

	if err != nil {
		return fmt.Errorf("restoring task: %w", err)
	}

	return nil
}

func CleanUp(fileName string, retention time.Duration) error {
	now := time.Now()
	var keptTasks []Task

	db, err := readJSON(fileName)

	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	for _, task := range db.Tasks {
		if task.IsDeleted == true {
			keptTasks = append(keptTasks, task)
		}

		if task.DeletedAt == nil {
			keptTasks = append(keptTasks, task)
		}

		if now.Sub(*task.DeletedAt) < retention {
			keptTasks = append(keptTasks, task)
		}

		// NOTE: deleted task will not be kept into the db
	}

	db.Tasks = keptTasks

	err = writeJSON(fileName, db)

	if err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	return nil
}

func RestoreAll(fileName string, retention time.Duration) error {
	db, err := readJSON(fileName)

	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	for index := range db.Tasks {
		db.Tasks[index].IsDeleted = false
		db.Tasks[index].DeletedAt = nil
	}

	err = writeJSON(fileName, db)

	if err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	return nil
}

func writeJSON(fileName string, db Database) error {
	file, err := os.OpenFile(fileName, os.O_TRUNC|os.O_RDWR, 0644)

	if err != nil {
		return fmt.Errorf("opening file: %w", err)
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "\t")

	encoder.Encode(db)

	return nil
}

func readJSON(fileName string) (Database, error) {
	var result Database

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDONLY, 0644)

	if err != nil {
		return Database{}, fmt.Errorf("opening file: %w", err)
	}

	defer file.Close()

	decoder := json.NewDecoder(file)

	err = decoder.Decode(&result)

	if err != nil {
		if errors.Is(err, io.EOF) {
			result = Database{}
		} else {
			return Database{}, fmt.Errorf("decoding file: %w", err)
		}
	}

	return result, nil
}

func softDelete(fileName string, taskIndex int) error {
	now := time.Now()

	db, err := readJSON(fileName)

	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	db.Tasks[taskIndex].IsDeleted = true
	db.Tasks[taskIndex].DeletedAt = &now
	db.Size--

	err = writeJSON(fileName, db)

	if err != nil {
		return fmt.Errorf("writing into file: %w", err)
	}

	return nil
}

func restoreTask(fileName string, taskIndex int) error {
	db, err := readJSON(fileName)

	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	db.Tasks[taskIndex].IsDeleted = false
	db.Tasks[taskIndex].DeletedAt = nil
	db.Size++

	err = writeJSON(fileName, db)

	if err != nil {
		return fmt.Errorf("writing into file: %w", err)
	}

	return nil
}
