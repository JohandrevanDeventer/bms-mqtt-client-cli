package engine

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/JohandrevanDeventer/bms-mqtt-client-cli/internal/config"
	mqttclient "github.com/JohandrevanDeventer/bms-mqtt-client-cli/pkg/mqtt"
	"github.com/JohandrevanDeventer/bms-mqtt-client-cli/pkg/persist"
	"go.uber.org/zap"
)

var (
	startTime time.Time
	endTime   time.Time
)

type Engine struct {
	cfg            *config.Config
	logger         *zap.Logger
	statePersister *persist.FilePersister
	client         *mqttclient.MQTTClient
	stopFileChan   chan struct{}
}

func NewEngine(cfg *config.Config, logger *zap.Logger, statePersister *persist.FilePersister) *Engine {
	return &Engine{
		cfg:            cfg,
		logger:         logger,
		statePersister: statePersister,
		stopFileChan:   make(chan struct{}), // Initialize stop file channel
	}
}

func (e *Engine) Run(ctx context.Context) {
	defer e.Cleanup()

	e.logger.Info("Starting application")

	// Create tmp directory
	tmpDirPath := "./tmp"
	if err := os.MkdirAll(tmpDirPath, os.ModePerm); err != nil {
		e.logger.Fatal("Failed to create tmp directory", zap.String("directory", tmpDirPath), zap.Error(err))
	} else {
		e.logger.Info("Tmp directory created", zap.String("directory", tmpDirPath))
	}

	startTime = time.Now()

	e.statePersister.Set("app", map[string]interface{}{})
	e.statePersister.Set("app.status", "running")
	e.statePersister.Set("app.name", e.cfg.System.AppName)
	e.statePersister.Set("app.version", fmt.Sprintf("%s-%d", e.cfg.System.AppVersion, e.cfg.System.BuildNumber))
	e.statePersister.Set("app.environment", e.cfg.Flags.Environment)
	e.statePersister.Set("app.start_time", startTime.Format(time.RFC3339))

	e.start()

	// Main Engine logic
	<-ctx.Done()
}

func (e *Engine) start() {
	go config.WatchAppConfigFileWithPolling(e.appConfigChangeCallback, 50*time.Millisecond, 50*time.Millisecond)

	stopFilePath := "./tmp/stop_signal"
	e.WatchStopFile(stopFilePath)

	e.initMQTTClient()

	go func() { e.tryMQTTConnection(5) }()
}

func (e *Engine) Cleanup() {
	// Perform Cleanup
	e.logger.Debug("Cleaning up")
	defer e.logger.Debug("Cleanup complete")

	// Disconnect MQTT client and set status to disconnected
	e.client.Disconnect()
	e.mqttStatePersistStop()

	// Delete the `tmp` directory if it exists
	tmpDir := "./tmp"
	if _, err := os.Stat(tmpDir); err == nil {
		err := os.RemoveAll(tmpDir)
		if err != nil {
			e.logger.Error("Failed to delete tmp directory", zap.String("directory", tmpDir), zap.Error(err))
		} else {
			e.logger.Info("Tmp directory deleted successfully", zap.String("directory", tmpDir))
		}
	} else if os.IsNotExist(err) {
		e.logger.Info("Tmp directory does not exist, skipping deletion", zap.String("directory", tmpDir))
	} else {
		e.logger.Error("Error checking tmp directory", zap.String("directory", tmpDir), zap.Error(err))
	}
}

func (e *Engine) Stop() {
	endTime = time.Now()

	duration := endTime.Sub(startTime)

	e.logger.Info("Stopping application")
	e.statePersister.Set("app.status", "stopped")
	e.statePersister.Set("app.end_time", endTime.Format(time.RFC3339))
	e.statePersister.Set("app.duration", duration.String())
}

func (e *Engine) WatchStopFile(stopFilePath string) {
	go func() {
		ticker := time.NewTicker(1 * time.Second) // Polling interval
		defer ticker.Stop()

		for {
			select {
			case <-e.stopFileChan: // Stop watching if channel is closed
				return
			default:
				if _, err := os.Stat(stopFilePath); err == nil {
					close(e.stopFileChan) // Signal stop file detection
					return
				}
				time.Sleep(1 * time.Second)
			}
		}
	}()
}

func (e *Engine) StopFileDetected() <-chan struct{} {
	return e.stopFileChan
}
