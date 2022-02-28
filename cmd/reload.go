/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"gcat/internal/entity"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// reloadCmd represents the reload command
var reloadCmd = &cobra.Command{
	Use:   "reload",
	Short: "reload catalog from disk",
	Long:  `Use this command to reload your catalog from disk.`,
	Run: func(cmd *cobra.Command, args []string) {
		var config entity.Config
		viper.Unmarshal(&config)
		fmt.Println("reload catalog", config.Paths[0].Name)
	},
}

func init() {
	rootCmd.AddCommand(reloadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// reloadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// reloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
