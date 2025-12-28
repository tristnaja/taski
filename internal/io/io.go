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
	var jsonData []byte
	var err error
	var existingDB Database
	var file *os.File

	file, err = os.OpenFile(fileName, os.O_RDWR, 0644)

	if err != nil {
		file.Close()
		os.WriteFile(fileName, nil, 0644)
		file, err = os.OpenFile(fileName, os.O_RDWR, 0644)

		if err != nil {
			return fmt.Errorf("opening file: %w", err)
		}
	}

	defer file.Close()

	data, err := io.ReadAll(file)

	if err != nil {
		return fmt.Errorf("turning file to bytes: ")
	}

	json.Unmarshal(data, &existingDB)

	existingDB.Tasks = append(existingDB.Tasks, *task)
	existingDB.Capacity++

	jsonData, err = json.MarshalIndent(existingDB, "", "\t")

	if err != nil {
		return fmt.Errorf("json indenting: %w", err)
	}

	err = os.WriteFile(fileName, jsonData, 0644)

	if err != nil {
		return fmt.Errorf("adding task: %w", err)
	}

	return nil
}

func ReadTask(fileName string) (*Database, error) {
	var result Database

	message, err := os.ReadFile(fileName)

	if err != nil {
		return &Database{}, fmt.Errorf("reading task: %w", err)
	}

	json.Unmarshal(message, &result)

	return &result, nil
}

func ChangeTask(fileName string, taskIndex int, newTitle string, newDescription string) error {
	if taskIndex < 0 {
		return errors.New("Index cannot be < 0")
	}

	if newTitle == "" && newDescription == "" {
		return errors.New("No value is changed")
	}

	var db Database
	var err error
	var file *os.File

	file, err = os.OpenFile(fileName, os.O_RDWR, 0644)

	if err != nil {
		return fmt.Errorf("opening file: %w", err)
	}

	defer file.Close()

	data, err := io.ReadAll(file)

	if err != nil {
		return fmt.Errorf("turning file to bytes: ")
	}

	json.Unmarshal(data, &db)

	if newTitle == "" {
		newTitle = db.Tasks[taskIndex].Title
	}

	if newDescription == "" {
		newDescription = db.Tasks[taskIndex].Description
	}

	db.Tasks[taskIndex].Title = newTitle
	db.Tasks[taskIndex].Description = newDescription
	db.Tasks[taskIndex].Date = time.Now()

	jsonData, err := json.MarshalIndent(db, "", "\t")

	if err != nil {
		return fmt.Errorf("reading db: %w", err)
	}

	err = os.WriteFile(fileName, jsonData, 0644)

	return nil
}

func RemoveTask(fileName string, taskIndex int) error {
	if taskIndex < 0 {
		return errors.New("Index cannot be < 0")
	}

	var db Database
	var err error
	var file *os.File

	file, err = os.OpenFile(fileName, os.O_RDWR, 0644)

	if err != nil {
		return fmt.Errorf("opening file: %w", err)
	}

	defer file.Close()

	data, err := io.ReadAll(file)

	if err != nil {
		return fmt.Errorf("turning file to bytes: ")
	}

	json.Unmarshal(data, &db)

	db.Tasks = append(db.Tasks[:taskIndex], db.Tasks[taskIndex+1:]...)
	db.Capacity--

	jsonData, err := json.MarshalIndent(db, "", "\t")

	if err != nil {
		return fmt.Errorf("reading db: %w", err)
	}

	err = os.WriteFile(fileName, jsonData, 0644)

	return nil
}
