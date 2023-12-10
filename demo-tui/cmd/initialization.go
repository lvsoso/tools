/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"demo-tui/internal/tui"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// initializationCmd represents the initialization command
var initializationCmd = &cobra.Command{
	Use:     "initialization",
	Short:   "Initialization demo tui config.",
	Long:    "Initialization demo tui config.",
	Example: "demo-tui init",
	Aliases: []string{"i", "init"},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		viper.AutomaticEnv()
		viper.SetEnvPrefix("dtui")
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		if err := tea.NewProgram(tui.NewInitPrompt(viper.GetString(cfgPath), homeDir)).Start(); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initializationCmd)
}
