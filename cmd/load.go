/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	handler "gcat/internal/handler"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// loadCmd represents the load command
var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "Load new catalog from disk",
	Long:  `Use this command to load a new catalog from a disk drive`,
	Run: func(cmd *cobra.Command, args []string) {
		path := "."
		user_path := cmd.Flag("path").Value.String()
		if user_path != "" {
			path = user_path
		}
		if path[0] != '/' {
			currentDir, err := os.Getwd()
			if err != nil {
				panic(err)
			}
			path = filepath.Join(currentDir, path)
		}
		fmt.Println("path:", path)
		handler.Load(path, "main")
	},
}

func init() {
	rootCmd.AddCommand(loadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	loadCmd.Flags().String("path", "p", "path to load")
}
