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
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	IsDeleted   bool      `json:"is_deleted"`
}

type Database struct {
	Capacity int    `json:"capacity"`
	Tasks    []Task `json:"tasks"`
}

func AddTask(task *Task, fileName string) error {
	db, err := ReadTask(fileName)

	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	db.Tasks = append(db.Tasks, *task)
	db.Capacity++

	file, err := os.OpenFile(fileName, os.O_TRUNC|os.O_RDWR, 0644)

	if err != nil {
		return fmt.Errorf("opening file: %w", err)
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "\t")

	err = encoder.Encode(db)

	if err != nil {
		return fmt.Errorf("encoding file: %w", err)
	}

	return nil
}

func ReadTask(fileName string) (Database, error) {
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

func ChangeTask(fileName string, taskIndex int, newTitle string, newDescription string) error {
	if taskIndex < 0 {
		return errors.New("Index cannot be < 0")
	}

	if newTitle == "" && newDescription == "" {
		return errors.New("No value is changed")
	}

	db, err := ReadTask(fileName)

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

func RemoveTask(fileName string, taskIndex int) error {
	if taskIndex < 0 {
		return errors.New("Index cannot be < 0")
	}

	var db Database
	var err error
	var file *os.File

	db, err = ReadTask(fileName)

	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	db.Tasks = append(db.Tasks[:taskIndex], db.Tasks[taskIndex+1:]...)
	db.Capacity--

	file, err = os.OpenFile(fileName, os.O_TRUNC|os.O_RDWR, 0644)

	if err != nil {
		return fmt.Errorf("opening file: %w", err)
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "\t")

	encoder.Encode(db)

	return nil
}
