/*
Copyright © 2025 Johandré van Deventer <johandre.vandeventer@rubiconsa.com>
*/
package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"time"

	"github.com/JohandrevanDeventer/bms-mqtt-client-cli/internal/config"
	"github.com/JohandrevanDeventer/bms-mqtt-client-cli/pkg/text_style"
	"github.com/JohandrevanDeventer/bms-mqtt-client-cli/utils"
	"github.com/spf13/cobra"
)

var (
	loggingLogLevel   string
	loggingFilePath   string
	loggingMaxSize    int
	loggingMaxBackups int
	loggingMaxAge     int
	loggingCompress   bool
	loggingAddTime    bool
)

// loggingCmd represents the logging command
var loggingCmd = &cobra.Command{
	Use:   "logging",
	Short: "Change the logging configuration",
	Long: `Change the logging configuration of the application.
All the configurations that can be changed are optional and can be seen under the flags section.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("logging called")
		bindLoggingFlags()
	},
}

func init() {
	rootCmd.AddCommand(loggingCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loggingCmd.PersistentFlags().String("foo", "", "A help for foo")
	loggingCmd.PersistentFlags().StringVar(&loggingLogLevel, "level", "", "Logging level")
	loggingCmd.PersistentFlags().StringVar(&loggingFilePath, "file-path", "", "Logging file path")
	loggingCmd.PersistentFlags().IntVar(&loggingMaxSize, "max-size", 0, "Logging max size")
	loggingCmd.PersistentFlags().IntVar(&loggingMaxBackups, "max-backups", 0, "Logging max backups")
	loggingCmd.PersistentFlags().IntVar(&loggingMaxAge, "max-age", 0, "Logging max age")
	loggingCmd.PersistentFlags().BoolVar(&loggingCompress, "compress", false, "Logging compress")
	loggingCmd.PersistentFlags().BoolVar(&loggingAddTime, "add-time", false, "Logging add time")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loggingCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func bindLoggingFlags() {
	newFlag := false

	if loggingLogLevel != "" && loggingLogLevel != cfg.App.Logging.Level {
		cfg.App.Logging.Level = loggingLogLevel
		newFlag = true
	}

	if loggingFilePath != "" && loggingFilePath != cfg.App.Logging.FilePath {
		cfg.App.Logging.FilePath = loggingFilePath
		newFlag = true
	}

	if loggingMaxSize != 0 && loggingMaxSize != cfg.App.Logging.MaxSize {
		cfg.App.Logging.MaxSize = loggingMaxSize
		newFlag = true
	}

	if loggingMaxBackups != 0 && loggingMaxBackups != cfg.App.Logging.MaxBackups {
		cfg.App.Logging.MaxBackups = loggingMaxBackups
		newFlag = true
	}

	if loggingMaxAge != 0 && loggingMaxAge != cfg.App.Logging.MaxAge {
		cfg.App.Logging.MaxAge = loggingMaxAge
		newFlag = true
	}

	if loggingCompress && loggingCompress != cfg.App.Logging.Compress {
		cfg.App.Logging.Compress = loggingCompress
		newFlag = true
	}

	if loggingAddTime && loggingAddTime != cfg.App.Logging.AddTime {
		cfg.App.Logging.AddTime = loggingAddTime
		newFlag = true
	}

	if newFlag {
		// Save the configuration
		// config.SaveConfig()

		fmt.Print("Updating logging configuration -> ")

		err := config.SaveConfig()
		if err != nil {
			var pathErr *fs.PathError
			// Check if the error is of type *fs.PathError
			if !errors.As(err, &pathErr) {
				fmt.Println(text_style.ColorText(text_style.Red, fmt.Sprintf("Failed to save configuration: %s", err)))
				os.Exit(1)
			}

		}

		time.Sleep(time.Duration(utils.GetRandomNumber(100, 500)) * time.Millisecond)

		fmt.Println(text_style.ColorText(text_style.Green, "Logging configuration updated successfully"))

		time.Sleep(time.Duration(utils.GetRandomNumber(100, 500)) * time.Millisecond)
	}
}
