package io

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Task struct {
	Task      string    `json:"task"`
	Date      time.Time `json:"date"`
	IsDeleted bool      `json:"is_deleted"`
}

type Database struct {
	Capacity int    `json:"capacity"`
	Tasks    []Task `json:"tasks"`
}

func AddTask(db *Database, task *Task, fileName string) error {
	var jsonData []byte
	var err error

	// NOTE: The value of isDBExist is just error value of ReadTask() func
	existingDB, isDBExist := ReadTask(fileName)

	// NOTE: if not nill (there IS an error), then there is no existing DB
	if isDBExist != nil {
		db.Tasks = append(db.Tasks, *task)
		db.Capacity++

		fmt.Println(db.Tasks)

		jsonData, err = json.MarshalIndent(db, "", "\t")

		if err != nil {
			return fmt.Errorf("json indenting: %w", err)
		}
		// NOTE: if nill (there is NO error), then there is an existing DB
	} else {
		existingDB.Tasks = append(existingDB.Tasks, *task)
		existingDB.Capacity++

		jsonData, err = json.MarshalIndent(existingDB, "", "\t")

		if err != nil {
			return fmt.Errorf("json indenting: %w", err)
		}
	}

	err = os.WriteFile(fileName, jsonData, 0644)

	if err != nil {
		return fmt.Errorf("adding task: %w", err)
	}

	return nil
}

func ReadTask(fileName string) (Database, error) {
	var result Database

	message, err := os.ReadFile(fileName)

	if err != nil {
		return Database{}, fmt.Errorf("reading task: %w", err)
	}

	json.Unmarshal(message, &result)

	return result, nil
}
