package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/tristnaja/taski/internal/io"
)

func RunRestore(args []string, fileName string) {
	cmd := flag.NewFlagSet("restore", flag.ExitOnError)
	var all bool
	var index int

	cmd.BoolVar(&all, "all", false, "Restore All? (t or f)")
	cmd.BoolVar(&all, "a", false, "Restore All? <t or f> (shorthand)")
	cmd.IntVar(&index, "i", -1, "Index of The Targeted Task (shorthand)")

	err := cmd.Parse(args)

	if all == false && index == -1 {
		cmd.Usage()
		os.Exit(1)
	}

	if all == true && index != -1 {
		fmt.Println("When restoring all, you do not need to input an index")
	}

	if all == false {
		err = io.RestoreTask(fileName, index)
		if err != nil {
			fmt.Printf("restoring task: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Task Restored")
		fmt.Printf("Index: %d\n", index)
		fmt.Println("\nTo view, type: taski view")
		fmt.Println("To restore, type: taski restore")
	} else {
		err = io.RestoreAll(fileName)

		if err != nil {
			fmt.Printf("deleting task: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("All Tasks in Trash is Restored")
		fmt.Println("\nTo view, type: taski view")
		fmt.Println("To restore, type: taski restore")
	}
}
