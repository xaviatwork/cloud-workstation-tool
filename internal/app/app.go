package app

import (
	"cloudworkstation/internal/cloudworkstation"
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version string = "v2.0"

func app() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "cw",
		Version: version,
		Short:   "cw is CLI tool to start/stop or create a tunnel to a Cloud Workstation (using gcloud)",
		Long: `cw is a wrapper around gcloud to start/stop Cloud Workstations and to create a tunnel to them.
v2.0 relies on gcloud completely, to avoid having to set ADC for starting/stopping the Cloud Workstation.`,
	}
	// add global flags
	cmd.PersistentFlags().StringP("conf", "c", "", "Configuration file")
	cmd.PersistentFlags().BoolP("quiet", "q", false, "Don't display Cloud Workstation information")
	// add subcommands
	cmd.AddCommand(tunnel(), start(), stop())
	return cmd
}

func Execute() {
	cmd := app()

	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func tunnel() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tunnel",
		Short: "starts the Cloud Workstation and creates a tunnel to it",
		RunE:  runTunnel,
	}
	cmd.PersistentFlags().Bool("only", false, "Create only the tunnel (do not start the Cloud Workstation)")
	return cmd
}

func loadConfigHelper(cmd *cobra.Command) (*cloudworkstation.CloudWorkstation, error) {
	configFile, err := cmd.Flags().GetString("conf")
	if err != nil {
		return nil, err
	}
	if configFile == "" {
		configFile = cloudworkstation.DefaultConfigPath()
	}
	return cloudworkstation.LoadConfig(configFile)
}

func runTunnel(cmd *cobra.Command, args []string) error {
	cw, err := loadConfigHelper(cmd)
	if err != nil {
		return err
	}

	onlyTunnel, err := cmd.Flags().GetBool("only")
	if err != nil {
		return err
	}

	quiet, _ := cmd.Flags().GetBool("quiet")
	if !quiet {
		cw.ShowConfig()
	}

	if err := cw.Tunnel(onlyTunnel); err != nil {
		return err
	}

	return nil
}

func start() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "starts the Cloud Workstation if it's stopped",
		RunE:  runStart,
	}

	return cmd
}

func runStart(cmd *cobra.Command, args []string) error {
	cw, err := loadConfigHelper(cmd)
	if err != nil {
		return err
	}
	quiet, _ := cmd.Flags().GetBool("quiet")
	if !quiet {
		cw.ShowConfig()
	}

	spin := cloudworkstation.NewSpinner(fmt.Sprintf("Workstation '%s'", cw.Name))
	cancel := spin.Start(context.Background())
	defer cancel()

	if err := cw.Start(); err != nil {
		return err
	}

	spin.Stop()
	fmt.Println("Cloud Workstation is running 🚀")
	return nil
}

func stop() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "stops the Cloud Workstation if it's running",
		RunE:  runStop,
	}

	return cmd
}

func runStop(cmd *cobra.Command, args []string) error {
	cw, err := loadConfigHelper(cmd)
	if err != nil {
		return err
	}

	quiet, _ := cmd.Flags().GetBool("quiet")
	if !quiet {
		cw.ShowConfig()
	}
	spin := cloudworkstation.NewSpinner(fmt.Sprintf("Workstation '%s'", cw.Name))
	cancel := spin.Start(context.Background())
	defer cancel()

	if err := cw.Stop(); err != nil {
		return err
	}

	spin.Stop()
	fmt.Println("Cloud Workstation is stopped ⏸️")
	return nil
}
