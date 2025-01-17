/*
Copyright © 2025 Johandré van Deventer <johandre.vandeventer@rubiconsa.com>
*/
package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/JohandrevanDeventer/bms-mqtt-client-cli/internal/config"
	"github.com/JohandrevanDeventer/bms-mqtt-client-cli/pkg/text_style"
	"github.com/JohandrevanDeventer/bms-mqtt-client-cli/utils"
	"github.com/spf13/cobra"
)

var cfg *config.Config

var (
	rootInitConfig  bool
	rootEnvironment string
	rootDebugMode   bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bms-mqtt-client-cli",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cfg = config.GetConfig()

		time.Sleep(time.Duration(utils.GetRandomNumber(500, 2000)) * time.Millisecond)

		// Initialize the configuration file
		if rootInitConfig {
			initConfig()
			bindRootFlags()
			os.Exit(0)
		}

		bindRootFlags()

		if cmd.Use == "bms-mqtt-client-cli" || cmd.Use == "start" {
			config.PrintInfo(false)
		} else {
			config.PrintInfo(true)
		}

		time.Sleep(time.Duration(utils.GetRandomNumber(500, 2000)) * time.Millisecond)
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cmd.Help())
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// Define persistent flags for the root command
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.bms-mqtt-client-cli.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&rootInitConfig, "init", "i", false, "Initialize the configuration file")
	rootCmd.PersistentFlags().StringVarP(&rootEnvironment, "environment", "e", "", "Environment to run the application in")
	rootCmd.PersistentFlags().BoolVarP(&rootDebugMode, "debug", "x", false, "Enable debug mode")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// bindRootFlags binds the root flags to the configuration
func bindRootFlags() {
	if rootEnvironment != "" && rootEnvironment != cfg.Flags.Environment {
		cfg.Flags.Environment = rootEnvironment
	}

	if rootDebugMode && rootDebugMode != cfg.Flags.DebugMode {
		cfg.Flags.DebugMode = rootDebugMode
	}

}

// initConfig initializes the configuration file
func initConfig() {
	fmt.Println(text_style.BoldText("Initializing configuration file(s)..."))

	newFiles, existingFiles, err := config.InitConfig()
	if err != nil {
		fmt.Println(text_style.ColorText(text_style.Red, fmt.Sprintf("Failed to initialize configuration file: %s", err)))
		os.Exit(1)
	}

	if len(newFiles) > 0 {
		for _, file := range newFiles {
			fmt.Println(text_style.ColorText(text_style.Green, fmt.Sprintf("-> Configuration file created: %s", file)))
		}
	}

	if len(existingFiles) > 0 {
		for _, file := range existingFiles {
			fmt.Println(text_style.ColorText(text_style.Yellow, fmt.Sprintf("-> Configuration file already exists: %s", file)))
		}
	}

	fmt.Println("")
}
