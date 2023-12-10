/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"demo-tui/internal/lfscache"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// lfscacheCmd represents the lfscache command
var lfscacheCmd = &cobra.Command{
	Use:   "lfscache",
	Short: "Caculate file' pointer from source dir  and add to git repo cache.",
	Long: `Caculate file' pointer from source dir  and add to git repo cache.
It will move file to .git/lfs/objects/aa/bb/file-sha256. And, generate pointer file to the repo root.
`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		viper.AutomaticEnv()
		viper.SetEnvPrefix("dtui")
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		if err := lfscache.LfsCache(viper.GetString(sourceDir), viper.GetString(repoRootDir), viper.GetInt(parallels)); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(lfscacheCmd)
	lfscacheCmd.PersistentFlags().String(sourceDir, "", "location of source files for add to lfs")
	lfscacheCmd.PersistentFlags().String(repoRootDir, "", "root dir  of lfs cached repo")
	lfscacheCmd.PersistentFlags().Int(parallels, 4, "concurrence of hash computing")
}
