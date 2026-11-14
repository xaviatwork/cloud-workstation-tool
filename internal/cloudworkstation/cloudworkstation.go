package cloudworkstation

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type TunnelPorts struct {
	Local  string `json:"local"`
	Remote string `json:"remote"`
}

type CloudWorkstation struct {
	ProjectId     string `json:"project_id"`
	Region        string `json:"region"`
	Cluster       string `json:"cluster"`
	ClusterConfig string `json:"cluster_config"`
	Name          string `json:"name"`
	TunnelPorts
}

func (cw CloudWorkstation) ShowConfig() {
	fmt.Printf("%-15s: %s\n", "Name", cw.Name)
	fmt.Printf("%-15s: %s\n", "ProjectID", cw.ProjectId)
	fmt.Printf("%-15s: %s\n", "Region", cw.Region)
	fmt.Printf("%-15s: %s\n", "Cluster", cw.Cluster)
	fmt.Printf("%-15s: %s\n", "Cluster Config", cw.ClusterConfig)
}

func (cw *CloudWorkstation) start_stop(action string) error {
	cwParams := []string{"workstations", action,
		cw.Name,
		"--project", cw.ProjectId,
		"--region", cw.Region,
		"--cluster", cw.Cluster,
		"--config", cw.ClusterConfig,
	}
	gcloudCmd := exec.CommandContext(context.Background(), "gcloud", cwParams...)

	if err := gcloudCmd.Run(); err != nil {
		return fmt.Errorf("Failed to %s Cloud Workstation %s: %v", action, cw.Name, err)
	}
	return nil
}

func (cw *CloudWorkstation) Stop() error {
	return cw.start_stop("stop")
}

func (cw *CloudWorkstation) Start() error {
	return cw.start_stop("start")
}

func (cw *CloudWorkstation) Tunnel(onlyTunnel bool) error {
	tunnelParams := []string{"workstations", "start-tcp-tunnel",
		"--project", cw.ProjectId,
		"--region", cw.Region,
		"--cluster", cw.Cluster,
		"--config", cw.ClusterConfig,
		"--local-host-port=localhost:" + cw.Local,
	}
	if !onlyTunnel {
		tunnelParams = append(tunnelParams, "--start-workstation")
	}
	tunnelParams = append(tunnelParams, cw.Name, cw.Remote)

	gcloudTunnelCmd := exec.CommandContext(context.Background(), "gcloud", tunnelParams...)
	// Create pipes to capture output
	stderr, _ := gcloudTunnelCmd.StderrPipe()
	if err := gcloudTunnelCmd.Start(); err != nil {
		log.Fatalf("Failed to start tunnel: %v", err)
	}
	fmt.Printf("Starting tunnel to %s:%s on localhost:%s...\n", cw.Name, cw.Remote, cw.Local)
	spin := NewSpinner(fmt.Sprintf("Workstation '%s'", cw.Name))
	cancel := spin.Start(context.Background())
	defer cancel()

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
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
		}
	}()

	// Wait for the "Listening" signal
	// <-ready
	if <-ready {
		spin.Stop()
		fmt.Println("\n✅ Tunnel is ready! You can now connect to localhost:" + cw.Local)
	} else {
		spin.Stop()
		fmt.Println("\n💥 Tunnel failed")
		os.Exit(1)
	}
	// Keep the program running to keep the tunnel open
	if err := gcloudTunnelCmd.Wait(); err != nil {
		return err
	}
	return nil
}
