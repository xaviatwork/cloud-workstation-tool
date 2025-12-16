package main

import (
	"cloudworkstation/internal/cloudworkstation"
	"context"
	"flag"
	"fmt"

	workstations "cloud.google.com/go/workstations/apiv1"
	workstationspb "cloud.google.com/go/workstations/apiv1/workstationspb"
)

func stop(args []string) error {
	fs := flag.NewFlagSet("stop", flag.ExitOnError)
	quiet := fs.Bool("quiet", false, "Do not output Cloud Workstation details")
	if err := fs.Parse(args); err != nil {
		return err
	}

	// Load configuration from ~/cloud-workstation/cw.cfg
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
	case workstationspb.Workstation_STATE_RUNNING:
		fmt.Println("Stopping Cloud Workstation ...")
		if !*quiet {
			cw.ShowConfig()
		}
		cw.Stop(c)
	case workstationspb.Workstation_STATE_STOPPED:
		fmt.Println("Workstation is already stopped")
	default:
		fmt.Printf("Workstation is %s\nPlease Wait.\n", w.State.String())
	}
	return nil
}
