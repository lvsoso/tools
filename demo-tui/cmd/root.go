/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"errors"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "demo-tui",
	Short: "A tool with tui.",
	Long:  `A tool with tui.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		viper.AutomaticEnv()
		viper.SetEnvPrefix("dtui")

		if _, err := os.Stat(viper.GetString(cfgPath)); errors.Is(err, os.ErrNotExist) {
			return errors.New(err.Error() + ": please run init to configure demo-tui\n")
		}
		return nil
	},
}

func Execute() error {
	dir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	rootCmd.PersistentFlags().String(cfgPath, dir+cfgDir+cfgFile, "location of the dtui config file")

	return rootCmd.ExecuteContext(context.Background())
}
