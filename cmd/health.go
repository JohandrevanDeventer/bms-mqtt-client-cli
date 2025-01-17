/*
Copyright © 2025 Johandré van Deventer <johandre.vandeventer@rubiconsa.com>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// healthCmd represents the health command
var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("health called")

		// Read the content of the JSON file
		filePath := "./persist/persist.json"
		data, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatalf("Error reading file: %v", err)
		}

		// Declare a variable to store the raw JSON
		var raw json.RawMessage
		if err := json.Unmarshal(data, &raw); err != nil {
			log.Fatalf("Error unmarshalling JSON: %v", err)
		}

		// Print the raw JSON (as bytes)
		fmt.Println("Raw JSON:", string(raw))
	},
}

func init() {
	rootCmd.AddCommand(healthCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// healthCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// healthCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
