package main

import (
	"bufio"
	"cloudworkstation/internal/cloudworkstation"
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

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
		fmt.Printf("Workstation is %s\nPlease Wait.\n", w.State.String())
	}
	return nil
}

func stop(args []string) error {
	fs := flag.NewFlagSet("stop", flag.ExitOnError)
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
	case workstationspb.Workstation_STATE_RUNNING:
		cw.ShowConfig()
		cw.Stop(c)
	case workstationspb.Workstation_STATE_STOPPED:
		fmt.Println("Worstation is already stopped")
	default:
		fmt.Printf("Workstation is %s\nPlease Wait.\n", w.State.String())
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
	case "stop":
		if err := stop(cmdArgs); err != nil {
			log.Fatalf("command stop failed: %s\n", err)
		}
	case "tunnel":
		if err := tunnel(cmdArgs); err != nil {
			log.Fatalf("command tunnel failed: %s\n", err)
		}
	default:
		log.Printf("unknown command: %s\n", args[1])
		return ErrUnknownCommand
	}
	return nil
}

func tunnel(args []string) error {
	fs := flag.NewFlagSet("stop", flag.ExitOnError)
	if err := fs.Parse(args); err != nil {
		return err
	}

	// Load configuration from ~/cloud-worstation/cw.cfg
	cw, err := cloudworkstation.Load(nil)
	if err != nil {
		return err
	}

	// Construct the gcloud command
	cmd := exec.CommandContext(context.Background(), "gcloud",
		"workstations", "start-tcp-tunnel",
		"--project", cw.ProjectId,
		"--region", cw.Region,
		"--cluster", cw.Cluster,
		"--config", cw.ClusterConfig,
		"--local-host-port=localhost:"+cw.Local,
		cw.Name,
		cw.Remote,
	)

	// Create pipes to capture output
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to start tunnel: %v", err)
	}

	fmt.Printf("Starting tunnel to %s:%s on localhost:%s...\n", cw.Name, cw.Remote, cw.Local)

	// Parse stderr to detect when the tunnel is actually ready
	// gcloud outputs "Listening on port [...]" to stderr when ready
	scanner := bufio.NewScanner(stderr)
	ready := make(chan bool)

	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			// fmt.Println("[gcloud]", line) // optional: print gcloud logs
			if strings.Contains(line, "Listening on port") {
				ready <- true
			}
			if strings.Contains(line, "ERROR") {
				fmt.Println(line)
				ready <- false
			}
		}
	}()

	// Wait for the "Listening" signal
	// <-ready
	if <-ready {
		fmt.Println("\n✅ Tunnel is ready! You can now connect to localhost:" + cw.Local)
	} else {
		fmt.Println("\n💥 Tunnel failed\n")
	}

	// Keep the program running to keep the tunnel open
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}

func main() {
	args := os.Args
	if err := run(args); err != nil {
		log.Fatalln(err)
	}
}
