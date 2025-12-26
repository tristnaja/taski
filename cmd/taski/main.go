package main

import (
	"fmt"
	"log"
	"time"

	"github.com/tristnaja/taski/internal/io"
)

func main() {
	log.SetPrefix("taski: ")
	log.SetFlags(0)
	fileName := "data.json"

	task := io.Task{
		Task:      "Test",
		Date:      time.Now(),
		IsDeleted: false,
	}

	task2 := io.Task{
		Task:      "Test2",
		Date:      time.Now(),
		IsDeleted: false,
	}

	db := io.Database{}

	err := io.AddTask(&db, &task, fileName)

	if err != nil {
		log.Fatal(err)
	}

	err = io.AddTask(&db, &task2, fileName)

	if err != nil {
		log.Fatal(err)
	}

	result, err := io.ReadTask(fileName)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
}
