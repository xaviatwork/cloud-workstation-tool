package main

import (
	"errors"
	"log"
	"os"
)

var ErrMissingCommand = errors.New("missing command")
var ErrUnknownCommand = errors.New("unknown command")

func start(args []string) error {
	log.Printf("args: %v\n", args)
	return nil
}

func run(args []string) error {
	if len(args) <= 1 {
		return ErrMissingCommand
	}
	cmdArgs := args[2:]
	switch args[1] {
	case "start":
		log.Printf("command: %s, args: %v\n", args[1], cmdArgs)
		if err := start(args); err != nil {
			log.Fatalf("command start failed: %s\n", err)
		}
	default:
		log.Printf("unknown command: %s\n", args[1])
		return ErrUnknownCommand
	}
	return nil
}

func main() {
	args := os.Args
	if err := run(args); err != nil {
		log.Fatalln(err)
	}
}
