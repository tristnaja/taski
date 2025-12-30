package cmd

import (
	"flag"
	"fmt"

	"github.com/tristnaja/taski/internal/io"
)

func RunDelete(args []string, fileName string) error {
	cmd := flag.NewFlagSet("delete", flag.ContinueOnError)
	var index int

	cmd.IntVar(&index, "index", -1, "Index of The Targeted Task")
	cmd.IntVar(&index, "i", -1, "Index of The Targeted Task (shorthand)")

	err := cmd.Parse(args)

	if err != nil {
		return fmt.Errorf("parsing arguments: %w", err)
	}

	if index == -1 {
		cmd.Usage()
		return fmt.Errorf("unfilled arguments")
	}

	err = io.RemoveTask(fileName, index)

	if err != nil {
		return fmt.Errorf("deleting task: %v\n", err)
	}

	fmt.Println("Deleted Task:")
	fmt.Printf("Index: %d\n", index)
	fmt.Println("\nTo view, type: taski view")
	fmt.Println("To restore, type: taski restore")

	return nil
}
