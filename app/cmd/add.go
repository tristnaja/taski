package cmd

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/tristnaja/taski/internal/io"
)

func RunAdd(args []string, fileName string) {
	cmd := flag.NewFlagSet("add", flag.ExitOnError)
	var title string
	var description string

	cmd.StringVar(&title, "title", "", "Task Title")
	cmd.StringVar(&title, "t", "", "Task Title (shorthand)")
	cmd.StringVar(&description, "desc", "", "Task Description")
	cmd.StringVar(&description, "d", "", "Task Description (shorthand)")

	err := cmd.Parse(args)

	if title == "" || description == "" {
		fmt.Println("usage: taski add --title <title> --desc <description>")
		cmd.Usage()
		os.Exit(1)
	}

	task := io.Task{
		Title:       title,
		Description: description,
		Date:        time.Now(),
		IsDeleted:   false,
	}

	err = io.AddTask(task, fileName)

	if err != nil {
		fmt.Printf("adding task: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Added New Task:")
	fmt.Printf("Title: %v\n", title)
	fmt.Printf("Description: %v\n", description)
	fmt.Println("\nTo view, type: taski view")
}
