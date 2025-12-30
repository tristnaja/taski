package cmd

import (
	"flag"
	"fmt"
	"time"

	"github.com/tristnaja/taski/internal/io"
)

func RunAdd(args []string, fileName string) error {
	cmd := flag.NewFlagSet("add", flag.ContinueOnError)
	var title string
	var description string

	cmd.StringVar(&title, "title", "", "Task Title")
	cmd.StringVar(&title, "t", "", "Task Title (shorthand)")
	cmd.StringVar(&description, "desc", "", "Task Description")
	cmd.StringVar(&description, "d", "", "Task Description (shorthand)")

	err := cmd.Parse(args)

	if err != nil {
		return fmt.Errorf("parsing arguments: %w", err)
	}

	if title == "" || description == "" {
		cmd.Usage()
		return fmt.Errorf("unfilled arguments")
	}

	task := io.Task{
		Title:       title,
		Description: description,
		Date:        time.Now(),
		IsDeleted:   false,
	}

	err = io.AddTask(task, fileName)

	if err != nil {
		return fmt.Errorf("adding task: %v\n", err)
	}

	fmt.Println("Added New Task:")
	fmt.Printf("Title: %v\n", title)
	fmt.Printf("Description: %v\n", description)
	fmt.Println("\nTo view, type: taski view")

	return nil
}
