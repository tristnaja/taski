package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/tristnaja/taski/internal/io"
)

func RunView(args []string, fileName string) {
	cmd := flag.NewFlagSet("view", flag.ExitOnError)
	err := cmd.Parse(args)

	db, err := io.ReadTask(fileName)

	if err != nil {
		fmt.Printf("viewing task: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Here is your Tasks:")
	for index, task := range db.Tasks {
		fmt.Printf("%d. %v\n", (index + 1), task.Title)
		fmt.Printf("index to target: %d\n", task.ID)
		fmt.Printf("Date: %v\n", task.Date.Format("02 Jan 2006, 15:04"))
		fmt.Printf("%v\n\n", task.Description)
	}
	fmt.Println("\nYou can Interact with your Tasks with:")
	fmt.Println("1. Adding new Task: \ntaski add --title <title> -desc <description>")
	fmt.Println("\n2. Changing Task: \ntaski change --index <index> --title <title> -desc <description>")
	fmt.Println("\n3. Deleting Task: \ntaski delete --index <index>")
	fmt.Println("\n4. Restoring Tasks: \ntaski delete --mode <mode> --index <index>")
	fmt.Println("\n5. Viewing Tasks: \ntaski view")
}
