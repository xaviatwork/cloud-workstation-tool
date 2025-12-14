package cloudworkstation

import (
	"context"
	"fmt"
	"time"

	workstations "cloud.google.com/go/workstations/apiv1"
	workstationspb "cloud.google.com/go/workstations/apiv1/workstationspb"
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

func (cw *CloudWorkstation) Start(c *workstations.Client) error {
	ctx := context.Background()
	req := &workstationspb.StartWorkstationRequest{
		Name: cw.FullResourceName(),
	}
	op, err := c.StartWorkstation(ctx, req)
	if err != nil {
		return err
	}

	spin := NewSpinner(fmt.Sprintf("Workstation '%s'", cw.Name))
	cancel := spin.Start(context.Background())
	defer cancel()

	for {
		_, err := op.Poll(ctx)
		if err != nil {
			return err
		}
		if op.Done() {
			spin.Stop("Started!")
			break
		}
		spin.UpdateMessage("Starting ...")
		time.Sleep(5 * time.Second)
	}
	return nil
}

func (cw *CloudWorkstation) Stop(c *workstations.Client) error {
	ctx := context.Background()
	req := &workstationspb.StopWorkstationRequest{
		Name: cw.FullResourceName(),
	}
	op, err := c.StopWorkstation(ctx, req)
	if err != nil {
		return err
	}

	spin := NewSpinner(fmt.Sprintf("Workstation '%s'", cw.Name))
	cancel := spin.Start(context.Background())
	defer cancel()

	for {
		_, err := op.Poll(ctx)
		if err != nil {
			return err
		}
		if op.Done() {
			spin.Stop("Stopped!")
			break
		}
		spin.UpdateMessage("Stopping ...")
		time.Sleep(5 * time.Second)
	}
	return nil
}

type TunnelPorts struct {
	Local  string `json:"local"`
	Remote string `json:"remote"`
}
