package main

import (
	"cloudworkstation/internal/cloudworkstation"
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	workstations "cloud.google.com/go/workstations/apiv1"
	workstationspb "cloud.google.com/go/workstations/apiv1/workstationspb"
)

var ErrMissingCommand = errors.New("missing command")
var ErrUnknownCommand = errors.New("unknown command")

func start(args []string) error {
	fs := flag.NewFlagSet("start", flag.ExitOnError)
	if err := fs.Parse(args); err != nil {
		return err
	}

	// Load configuration from ~/cloud-worstation/cw.cfg
	cw, err := cloudworkstation.Load(nil)
	if err != nil {
		return err
	}

	// Prepare Workstation Client
	ctx := context.Background()
	c, err := workstations.NewClient(ctx)
	if err != nil {
		return err
	}
	// Prepare request
	getReq := &workstationspb.GetWorkstationRequest{
		Name: cw.FullResourceName(),
	}
	// Execute request
	w, err := c.GetWorkstation(ctx, getReq)
	if err != nil {
		return err
	}

	// Do something with the response
	switch w.State {
	case workstationspb.Workstation_STATE_STOPPED:
		cw.ShowConfig()
		cw.Start(c)
	case workstationspb.Workstation_STATE_RUNNING:
		fmt.Println("Worstation is already running")
	default:
		fmt.Println("not stopped or running")
	}
	return nil
}

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
