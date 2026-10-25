package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var version = "dev"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Shows current version of awsaudit",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Fprintln(cmd.OutOrStdout(), version)
		fmt.Fprintln(cmd.OutOrStdout(), "region:", viper.GetString(keyRegion))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
