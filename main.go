package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "devcert [space separated domain names]",
		Short: "Self-signed trusted certificates for local development.",
		Long:  `Generate self-signed, trusted certificates for local development.`,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := devcertExec(args)
			return err
		},
	}

	var infoCmd = &cobra.Command{
		Use:   "info [.crt file]",
		Short: "Print the certificate information.",
		Long:  "Print out the certificate information that the .crt file contains.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := devcertInfo(args)
			return err
		},
	}

	// Add the info command to the CLI commands
	rootCmd.AddCommand(infoCmd)

	var uninstallCmd = &cobra.Command{
		Use:   "uninstall",
		Short: "Uninstall devcert.",
		Long:  "Uninstall devcert, cleanup the CA and certificate files. This command uninstalls the devcert command,\nif you want to uninstall a crt file, remove it from the ~/.devcert folder.",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := uninstallDevcert()
			return err
		},
	}

	// Add the uninstall command to the CLI commands
	rootCmd.AddCommand(uninstallCmd)

	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}
