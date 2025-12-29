package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/tristnaja/taski/internal/io"
)

func RunChange(args []string, fileName string) {
	cmd := flag.NewFlagSet("change", flag.ExitOnError)
	var index int
	var title string
	var description string

	cmd.IntVar(&index, "index", -1, "Index of The Targeted Task")
	cmd.IntVar(&index, "i", -1, "Index of The Targeted Task (shorthand)")
	cmd.StringVar(&title, "title", "", "New Task Title")
	cmd.StringVar(&title, "t", "", "New Task Title (shorthand)")
	cmd.StringVar(&description, "desc", "", "New Task Description")
	cmd.StringVar(&description, "d", "", "New Task Description (shorthand)")

	err := cmd.Parse(args)

	if index == -1 && title == "" || index == -1 && description == "" {
		cmd.Usage()
		os.Exit(1)
	}

	err = io.ChangeTask(fileName, index, title, description)

	if err != nil {
		fmt.Printf("changing task: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Changed Task:")
	fmt.Printf("Title: %v\n", title)
	fmt.Printf("Index: %d\n", index)
	fmt.Printf("Description: %v\n", description)
	fmt.Println("\nTo view, type: taski view")
}
