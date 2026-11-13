package app

import (
	"fmt"
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
