package app

import (
	"cloudworkstation/internal/cloudworkstation"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// var rootCmd = &cobra.Command{
// 	Use:   "cw",
// 	Short: "cw is CLI tool to start/stop or create a tunnel to a Cloud Workstation",
// }

func app() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cw",
		Short: "cw is CLI tool to start/stop or create a tunnel to a Cloud Workstation",
	}
	// add global flags
	cmd.PersistentFlags().StringP("conf", "c", "", "Configuration file")
	// add subcommands
	cmd.AddCommand(tunnel(), start(), stop())
	// TODO
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
		Short: "creates a tunnel to the remote Cloud Workstation",
		RunE:  runTunnel,
	}

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
	log.Println("using configFile", configFile)
	log.Println("cw.Name", cw.Name)

	fmt.Println("running the tunnel command")
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
	fmt.Println("start the Cloud Workstation")
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
	fmt.Println("stop the Cloud Workstation")
	return nil
}
