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

func (cw *CloudWorkstation) FullResourceName() string {
	return fmt.Sprintf("projects/%s/locations/%s/workstationClusters/%s/workstationConfigs/%s/workstations/%s",
		cw.ProjectId,
		cw.Region,
		cw.Cluster,
		cw.ClusterConfig,
		cw.Name)
}

func (cw *CloudWorkstation) Stop() error {
	stopParams := []string{"workstations", "stop", cw.Name,
		"--project", cw.ProjectId,
		"--region", cw.Region,
		"--cluster", cw.Cluster,
		"--config", cw.ClusterConfig,
	}
	gcloudStopCmd := exec.CommandContext(context.Background(), "gcloud", stopParams...)

	if err := gcloudStopCmd.Run(); err != nil {
		return fmt.Errorf("Failed to stop Cloud Workstation %s: %v", cw.Name, err)
	}
	return nil
}

func (cw *CloudWorkstation) Start() error {
	startParams := []string{"workstations", "start", cw.Name,
		"--project", cw.ProjectId,
		"--region", cw.Region,
		"--cluster", cw.Cluster,
		"--config", cw.ClusterConfig,
	}
	gcloudStartCmd := exec.CommandContext(context.Background(), "gcloud", startParams...)

	if err := gcloudStartCmd.Run(); err != nil {
		return fmt.Errorf("Failed to start Cloud Workstation %s: %v", cw.Name, err)
	}
	return nil
}

type TunnelPorts struct {
	Local  string `json:"local"`
	Remote string `json:"remote"`
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
		fmt.Println("\n✅ Tunnel is ready! You can now connect to localhost:" + cw.Local)
	} else {
		fmt.Println("\n💥 Tunnel failed")
		os.Exit(1)
	}
	// Keep the program running to keep the tunnel open
	if err := gcloudTunnelCmd.Wait(); err != nil {
		return err
	}
	return nil
}
