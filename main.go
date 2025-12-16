package main

import (
	"errors"
	"log"
	"os"
)

func main() {
	args := os.Args
	if err := run(args); err != nil {
		log.Fatalln(err)
	}
}

var ErrMissingCommand = errors.New("missing command")
var ErrUnknownCommand = errors.New("unknown command")

func run(args []string) error {
	if len(args) <= 1 {
		return ErrMissingCommand
	}
	cmdArgs := args[2:]
	switch args[1] {
	case "start":
		if err := start(cmdArgs); err != nil {
			log.Fatalf("command start failed: %s\n", err)
		}
	case "stop":
		if err := stop(cmdArgs); err != nil {
			log.Fatalf("command stop failed: %s\n", err)
		}
	case "tunnel":
		if err := tunnel(cmdArgs); err != nil {
			log.Fatalf("command tunnel failed: %s\n", err)
		}
	default:
		log.Printf("Unknown command: %s\n", args[1])
		return ErrUnknownCommand
	}
	return nil
}
