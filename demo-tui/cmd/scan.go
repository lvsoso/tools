/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"demo-tui/internal/scan"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan a drectory or file.",
	Long:  `Scan a drectory or file.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		viper.AutomaticEnv()
		viper.SetEnvPrefix("dtui")
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := scan.RunScan(context.Background(), viper.GetString(targetDir)); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
}
