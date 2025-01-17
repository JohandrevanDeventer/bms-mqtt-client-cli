/*
Copyright © 2025 Johandré van Deventer <johandre.vandeventer@rubiconsa.com>
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/JohandrevanDeventer/bms-mqtt-client-cli/internal/config"
	"github.com/JohandrevanDeventer/bms-mqtt-client-cli/internal/engine"
	"github.com/JohandrevanDeventer/bms-mqtt-client-cli/pkg/logging"
	"github.com/JohandrevanDeventer/bms-mqtt-client-cli/pkg/persist"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var logger *zap.Logger

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the Rubicon BMS MQTT Client",
	Long:  `This command starts the Rubicon BMS MQTT Client.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("start called")

		cfg = config.GetConfig()

		initLogger(cfg)

		statePersister, err := initPersist(cfg)
		if err != nil {
			logger.Fatal("failed to initialize the state persister", zap.Error(err))
		}

		// Graceful shutdown handling
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		svc := engine.NewEngine(cfg, logger, statePersister)

		// Goroutine to handle stop signals or stop file detection
		go func() {
			select {
			case <-ctx.Done(): // Handle system interrupt (e.g., Ctrl+C)
				logger.Warn("Received signal to stop the application")
			case <-svc.StopFileDetected(): // Stop file detected by Engine
				logger.Warn("Stop file detected, shutting down application")
			}

			// Ensure application cleanup and shutdown
			// svc.Cleanup() // Cleanup resources
			svc.Stop() // Stop the engine
			stop()     // Cancel the context
		}()

		// go func() {
		// 	<-ctx.Done()
		// 	svc.Stop()
		// }()

		defer func() {
			if r := recover(); r != nil {
				logger.Error("recovered from panic", zap.Any("panic", r))
			}
		}()

		svc.Run(ctx)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initLogger configures the logger based on the app config
func initLogger(cfg *config.Config) {
	// Create a new logging config with values from the app config
	loggingConfig := logging.NewLoggingConfig(
		cfg.App.Logging.Level,
		cfg.App.Logging.FilePath,
		cfg.App.Logging.MaxSize,
		cfg.App.Logging.MaxBackups,
		cfg.App.Logging.MaxAge,
		cfg.App.Logging.Compress,
		cfg.Flags.DebugMode,
		cfg.App.Logging.AddTime,
	)

	// Get a new logger based on the config
	logger = logging.NewLogger(loggingConfig)
	logger = logging.GetLogger("main")
}

// initPersist initializes the file persister.
func initPersist(cfg *config.Config) (*persist.FilePersister, error) {
	var err error

	statePersister, err := persist.NewFilePersister(cfg.PersistFilePath)
	if err != nil {
		return nil, err
	}

	return statePersister, nil
}
