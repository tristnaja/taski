package io

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

type Task struct {
	Title       string    `json:"title"`
	Description string    `josn:"description"`
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

	// NOTE: The value of isDBExist is just error value of ReadTask() func
	existingDB, isDBExist := ReadTask(fileName)

	// NOTE: if not nill (there IS an error), then there is no existing DB
	if isDBExist != nil {
		db := Database{}

		db.Tasks = append(db.Tasks, *task)
		db.Capacity++

		fmt.Println(db.Tasks)

		jsonData, err = json.MarshalIndent(db, "", "\t")

		// NOTE: if nill (there is NO error), then there is an existing DB
	} else {
		existingDB.Tasks = append(existingDB.Tasks, *task)
		existingDB.Capacity++

		jsonData, err = json.MarshalIndent(existingDB, "", "\t")
	}

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

func ChangeTask(filename string, taskIndex int, newTitle string, newDescription string) error {
	if newTitle == "" && newDescription == "" {
		return errors.New("No value is changed")
	}

	db, err := ReadTask(filename)

	if err != nil {
		return fmt.Errorf("reading db: %w", err)
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

	jsonData, err := json.MarshalIndent(db, "", "\t")

	if err != nil {
		return fmt.Errorf("reading db: %w", err)
	}

	err = os.WriteFile(filename, jsonData, 0644)

	return nil
}

func RemoveTask(filename string, taskIndex int) error {
	db, err := ReadTask(filename)

	if err != nil {
		return fmt.Errorf("reading db: %w", err)
	}

	db.Tasks = append(db.Tasks[:taskIndex], db.Tasks[taskIndex+1:]...)

	jsonData, err := json.MarshalIndent(db, "", "\t")

	if err != nil {
		return fmt.Errorf("reading db: %w", err)
	}

	err = os.WriteFile(filename, jsonData, 0644)

	return nil
}
