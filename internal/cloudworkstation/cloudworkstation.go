package cloudworkstation

import (
	"context"
	"fmt"
	"time"

	workstations "cloud.google.com/go/workstations/apiv1"
	workstationspb "cloud.google.com/go/workstations/apiv1/workstationspb"
	"github.com/yarlson/pin"
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

	s := pin.New("Loading ...",
		pin.WithSpinnerColor(pin.ColorCyan),
		pin.WithTextColor(pin.ColorYellow),
		pin.WithDoneSymbol('✔'),
		pin.WithDoneSymbolColor(pin.ColorGreen),
		pin.WithPrefix(fmt.Sprintf("Workstation '%s'", cw.Name)),
		pin.WithSeparatorColor(pin.ColorGray),
		pin.WithPosition(pin.PositionRight),
	)
	cancel := s.Start(context.Background())
	defer cancel()

	for {
		_, err := op.Poll(ctx)
		if err != nil {
			return err
		}
		if op.Done() {
			s.Stop("Running!")
			// fmt.Printf("Workstation state: %s\n", w.State.String())
			break
		}
		s.UpdateMessage("Starting ...")
		time.Sleep(5 * time.Second)
	}
	return nil
}

type TunnelPorts struct {
	Local  string `json:"local"`
	Remote string `json:"remote"`
}
