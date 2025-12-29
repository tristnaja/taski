package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/tristnaja/taski/internal/io"
)

func RunDelete(args []string, fileName string) {
	cmd := flag.NewFlagSet("delete", flag.ExitOnError)
	var index int

	cmd.IntVar(&index, "index", -1, "Index of The Targeted Task")
	cmd.IntVar(&index, "i", -1, "Index of The Targeted Task (shorthand)")

	err := cmd.Parse(args)

	if index == -1 {
		cmd.Usage()
		os.Exit(1)
	}

	err = io.RemoveTask(fileName, index)

	if err != nil {
		fmt.Printf("deleting task: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Deleted Task:")
	fmt.Printf("Index: %d\n", index)
	fmt.Println("\nTo view, type: taski view")
	fmt.Println("To restore, type: taski restore")
}
