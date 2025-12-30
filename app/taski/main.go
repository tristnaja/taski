package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/tristnaja/taski/app/cmd"
	"github.com/tristnaja/taski/internal/io"
)

func main() {
	log.SetPrefix("taski: ")
	log.SetFlags(0)
	exe, err := os.Executable()

	if err != nil {
		errLog := fmt.Errorf("locating executables: %w", err)
		log.Fatal(errLog)
	}

	fileDir := filepath.Dir(exe)
	fileName := filepath.Join(fileDir, "data.json")
	trashDue := 30 * 24 * time.Hour

	err = io.CleanUp(fileName, trashDue)

	if err != nil {
		log.Printf("cleanup failed: %v", err)
	}

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: taski <cmd> <args>")
		errLog := fmt.Errorf("parsing args: arguments not enough")
		log.Fatal(errLog)
	}

	switch os.Args[1] {
	case "add":
		err = cmd.RunAdd(os.Args[2:], fileName)
	case "change":
		err = cmd.RunChange(os.Args[2:], fileName)
	case "delete":
		err = cmd.RunDelete(os.Args[2:], fileName)
	case "restore":
		err = cmd.RunRestore(os.Args[2:], fileName)
	case "view":
		err = cmd.RunView(os.Args[2:], fileName)
	default:
		log.Fatal("unknown command, usable: add, change, restore, delete, view")
	}

	if err != nil {
		log.Fatal(err)
	}
}
