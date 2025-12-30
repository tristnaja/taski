package cmd

import (
	"flag"
	"fmt"

	"github.com/tristnaja/taski/internal/io"
)

func RunRestore(args []string, fileName string) error {
	cmd := flag.NewFlagSet("restore", flag.ContinueOnError)
	var all bool
	var index int

	cmd.BoolVar(&all, "all", false, "Restore All? (t or f)")
	cmd.BoolVar(&all, "a", false, "Restore All? <t or f> (shorthand)")
	cmd.IntVar(&index, "index", -1, "Index of The Targeted Task")
	cmd.IntVar(&index, "i", -1, "Index of The Targeted Task (shorthand)")

	err := cmd.Parse(args)

	if err != nil {
		return fmt.Errorf("parsing arguments: %w", err)
	}

	if all == false && index == -1 {
		cmd.Usage()
		return fmt.Errorf("unfilled arguments")
	}

	if all == true && index != -1 {
		return fmt.Errorf("When restoring all, you do not need to input an index")
	}

	if all == false {
		err = io.RestoreTask(fileName, index)
		if err != nil {
			return fmt.Errorf("restoring task: %v\n", err)
		}

		fmt.Println("Task Restored")
		fmt.Printf("Index: %d\n", index)
		fmt.Println("\nTo view, type: taski view")
		fmt.Println("To restore, type: taski restore")
	} else {
		err = io.RestoreAll(fileName)

		if err != nil {
			return fmt.Errorf("restoring task: %v\n", err)
		}

		fmt.Println("All Tasks in Trash is Restored")
		fmt.Println("\nTo view, type: taski view")
		fmt.Println("To restore, type: taski restore")
	}

	return nil
}
