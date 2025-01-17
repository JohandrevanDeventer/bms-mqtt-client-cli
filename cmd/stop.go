/*
Copyright © 2025 Johandré van Deventer <johandre.vandeventer@rubiconsa.com>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the Rubicon BMS MQTT Client",
	Long:  `This command stops the Rubicon BMS MQTT Client.`,
	Run: func(cmd *cobra.Command, args []string) {
		stopFilePath := "./tmp/stop_signal"
		if _, err := os.Create(stopFilePath); err != nil {
			fmt.Println("Failed to create stop file:", err)
			return
		}
		fmt.Println("Stop file created at", stopFilePath)
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// stopCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// stopCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
