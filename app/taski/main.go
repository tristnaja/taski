package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/tristnaja/taski/app/cmd"
	"github.com/tristnaja/taski/internal/io"
)

func main() {
	log.SetPrefix("taski: ")
	log.SetFlags(0)
	fileName := "../../bin/data.json"
	trashDue := 30 * 24 * time.Hour

	io.CleanUp(fileName, trashDue)

	if len(os.Args) < 2 {
		fmt.Println("usage: taski <cmd> <args>")
		fmt.Println("usage: taski add --title <title> --desc <desc>")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add":
		cmd.RunAdd(os.Args[2:], fileName)
	case "change":
		cmd.RunChange(os.Args[2:], fileName)
	case "delete":
		cmd.RunDelete(os.Args[2:], fileName)
	case "restore":
		cmd.RunRestore(os.Args[2:], fileName)
	case "view":
		cmd.RunView(os.Args[2:], fileName)
	default:
		fmt.Println("unknown command, usable: add, change, restore, delete, view")
	}
}
