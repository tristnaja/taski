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
	trashDue := 30 * 24 * time.Hour

	io.CleanUp(fileName, trashDue)

	task := io.Task{
		Title:       "Test",
		Description: "Halo",
		Date:        time.Now(),
		IsDeleted:   false,
	}

	task2 := io.Task{
		Title:       "Test2",
		Description: "HaloHalo",
		Date:        time.Now(),
		IsDeleted:   false,
	}

	err := io.AddTask(&task, fileName)

	if err != nil {
		log.Fatal(err)
	}

	err = io.AddTask(&task2, fileName)

	if err != nil {
		log.Fatal(err)
	}

	result, err := io.ReadTask(fileName)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)

	for index, value := range result.Tasks {
		fmt.Printf("\n%d. %v\n", (index + 1), value.Title)
		fmt.Printf("%v\n", value.Description)
		fmt.Printf("Date: %v\n", value.Date)
	}

	err = io.ChangeTask(fileName, 0, "Title3", "")

	if err != nil {
		log.Fatal(err)
	}

	err = io.ChangeTask(fileName, 1, "Title4", "HaloHaloHalo")

	if err != nil {
		log.Fatal(err)
	}

	result, err = io.ReadTask(fileName)

	if err != nil {
		log.Fatal(err)
	}

	for index, value := range result.Tasks {
		fmt.Printf("\n%d. %v\n", (index + 1), value.Title)
		fmt.Printf("%v\n", value.Description)
		fmt.Printf("Date: %v\n", value.Date)
	}
	//
	err = io.RemoveTask(fileName, 0)

	if err != nil {
		log.Fatal(err)
	}

	err = io.RemoveTask(fileName, 1)

	if err != nil {
		log.Fatal(err)
	}

	result, err = io.ReadTask(fileName)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println()

	for index, value := range result.Tasks {
		fmt.Printf("\n%d. %v\n", (index + 1), value.Title)
		fmt.Printf("%v\n", value.Description)
		fmt.Printf("Date: %v\n", value.Date)
	}
}
