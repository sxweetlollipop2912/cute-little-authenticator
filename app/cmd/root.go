package cmd

import (
	"fmt"
	"github.com/99designs/keyring"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:     "",
	Version: "0.0.1",
	Short:   "cute little authenticator iz here :>",
}

func Execute() {
	rootCmd.PersistentFlags().BoolVarP(&keyring.Debug, "verbose", "v", false, "verbose output")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
