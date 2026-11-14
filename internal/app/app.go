package app

import (
	"cloudworkstation/internal/cloudworkstation"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func app() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cw",
		Short: "cw is CLI tool to start/stop or create a tunnel to a Cloud Workstation",
	}
	// add global flags
	cmd.PersistentFlags().StringP("conf", "c", "", "Configuration file")
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

func runTunnel(cmd *cobra.Command, args []string) error {
	configFile, err := cmd.Flags().GetString("conf")
	if err != nil {
		return err
	}
	if configFile == "" {
		configFile = cloudworkstation.DefaultConfigPath()
	}
	cw, err := cloudworkstation.LoadConfig(configFile)
	if err != nil {
		return err
	}

	onlyTunnel, err := cmd.Flags().GetBool("only")
	if err != nil {
		return err
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
	configFile, err := cmd.Flags().GetString("conf")
	if err != nil {
		return err
	}
	if configFile == "" {
		configFile = cloudworkstation.DefaultConfigPath()
	}
	cw, err := cloudworkstation.LoadConfig(configFile)
	if err != nil {
		return err
	}

	if err := cw.Start(); err != nil {
		return err
	}

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
	configFile, err := cmd.Flags().GetString("conf")
	if err != nil {
		return err
	}
	if configFile == "" {
		configFile = cloudworkstation.DefaultConfigPath()
	}
	cw, err := cloudworkstation.LoadConfig(configFile)
	if err != nil {
		return err
	}

	if err := cw.Stop(); err != nil {
		return err
	}

	return nil
}
