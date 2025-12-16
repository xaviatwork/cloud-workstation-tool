package main

import (
	"cloudworkstation/internal/cloudworkstation"
	"context"
	"flag"
	"fmt"

	workstations "cloud.google.com/go/workstations/apiv1"
	workstationspb "cloud.google.com/go/workstations/apiv1/workstationspb"
)

func start(args []string) error {
	fs := flag.NewFlagSet("start", flag.ExitOnError)
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
	case workstationspb.Workstation_STATE_STOPPED:
		fmt.Println("Starting Cloud Workstation ...")
		if !*quiet {
			cw.ShowConfig()
		}
		cw.Start(c)
	case workstationspb.Workstation_STATE_RUNNING:
		fmt.Println("Workstation is already running")
	default:
		// opReq := &longrunningpb.GetOperationRequest{
		// 	Name: cw.FullResourceName(),
		// }
		// op, err := c.GetOperation(ctx, opReq)
		// if err != nil {
		// 	return err
		// }
		// for {
		// 	fmt.Printf("operation: %s\n", op.Name)
		// 	if op.Done {
		// 		fmt.Println("operation completed!")
		// 		break
		// 	}
		// 	time.Sleep(5 * time.Second)
		// }
		fmt.Printf("Workstation is %s\nWait a little bit and try again.\n", w.State.String())
	}
	return nil
}
