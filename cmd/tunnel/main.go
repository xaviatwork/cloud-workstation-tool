package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

type CloudWorkstationConfig struct {
	ProjectId     string `json:"project_id"`
	Region        string `json:"region"`
	Cluster       string `json:"cluster"`
	ClusterConfig string `json:"cluster_config"`
	Name          string `json:"name"`
	LocalPort     string `json:"local_port"`
}

func readConfig(path2config string, cwcfg *CloudWorkstationConfig) {
	data, err := os.ReadFile(path2config)
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
	err = json.Unmarshal(data, &cwcfg)
	if err != nil {
		log.Fatalf("Invalid configuration data: %v", err)
	}
	fmt.Printf("Configuration:\n")
	fmt.Printf("- Workstation: %s\n", cwcfg.Name)
	fmt.Printf("- Project: %s\n", cwcfg.ProjectId)
	fmt.Printf("- LocalPort: %s\n", cwcfg.LocalPort)
	fmt.Printf("- Cluster: %s\n", cwcfg.Cluster)
	fmt.Printf("- Config: %s\n", cwcfg.ClusterConfig)
	fmt.Printf("- Region: %s\n", cwcfg.Region)
}

func main() {
	// Read configuration from ~/.config/cloud-workstation-manager/cw.conf
	var cwcfg CloudWorkstationConfig
	var cfgfile string
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Unable to determine user's home directory: %s\n", err)
		os.Exit(1)
	}
	cfgfile = path.Join(home, ".config", "cloud-workstation-config", "cw.cfg")
	readConfig(cfgfile, &cwcfg)
	// Configuration
	// project := "vwgitc-gcpf-workstation-p-0001"
	// region := "europe-west3"
	// cluster := "gcpf-cw-cluster"
	// config := "default-config"
	// workstation := "workstation-xavi"
	remotePort := "22" // Port running on the workstation (SSH)
	// localPort := "8910" // Port to expose locally

	// Construct the gcloud command
	cmd := exec.CommandContext(context.Background(), "gcloud",
		"workstations", "start-tcp-tunnel",
		"--project", cwcfg.ProjectId,
		"--region", cwcfg.Region,
		"--cluster", cwcfg.Cluster,
		"--config", cwcfg.ClusterConfig,
		"--local-host-port=localhost:"+cwcfg.LocalPort,
		cwcfg.Name,
		remotePort,
	)

	// Create pipes to capture output
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to start tunnel: %v", err)
	}

	fmt.Printf("Starting tunnel to %s:%s on localhost:%s...\n", cwcfg.Name, remotePort, cwcfg.LocalPort)

	// Parse stderr to detect when the tunnel is actually ready
	// gcloud outputs "Listening on port [...]" to stderr when ready
	scanner := bufio.NewScanner(stderr)
	ready := make(chan bool)

	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println("[gcloud]", line) // optional: print gcloud logs
			if strings.Contains(line, "Listening on port") {
				ready <- true
			}
		}
	}()

	// Wait for the "Listening" signal
	<-ready
	fmt.Println("✅ Tunnel is ready! You can now connect to localhost:" + cwcfg.LocalPort)

	// Keep the program running to keep the tunnel open
	cmd.Wait()
}
