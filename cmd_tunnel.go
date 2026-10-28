package main

import (
	"bufio"
	"cloudworkstation/internal/cloudworkstation"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func tunnel(args []string) error {
	fs := flag.NewFlagSet("stop", flag.ExitOnError)
	up := fs.Bool("up", false, "Start the Cloud Workstation if it is stopped, then create the tunnel")
	quiet := fs.Bool("quiet", false, "do not output Cloud Workstation details")
	if err := fs.Parse(args); err != nil {
		return err
	}

	// Load configuration from ~/cloud-workstation/cw.cfg
	cw, err := cloudworkstation.Load(nil)
	if err != nil {
		return err
	}

	if *up {
		fmt.Printf("Starting Cloud Workstation %s ...\n", cw.Name)
		startArgs := []string{}
		if *quiet {
			startArgs = append(startArgs, "-quiet")
		}
		if err := start(startArgs); err != nil {
			return err
		}
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
			// if strings.Contains(line, "ERROR") {
			// 	fmt.Println(line)
			// 	ready <- false
			// }
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
		}
	}()

	// Wait for the "Listening" signal
	// <-ready
	if <-ready {
		fmt.Println("\n✅ Tunnel is ready! You can now connect to localhost:" + cw.Local)
	} else {
		fmt.Println("\n💥 Tunnel failed")
	}

	// Keep the program running to keep the tunnel open
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}
