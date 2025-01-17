package config

import (
	"fmt"
	"os"
	"time"

	"github.com/JohandrevanDeventer/bms-mqtt-client-cli/utils"
)

var appConfigFilePath = fmt.Sprintf("%s/%s", configRoot, appConfigFile)

var appConfig *AppConfig

var defaultAppConfig = AppConfig{
	Logging: defaultLoggingConfig,
	Mqtt:    defaultMQTTConfig,
}

var defaultLoggingConfig = LoggingConfig{
	Level:      "info",
	FilePath:   "./logs/app.log",
	MaxSize:    100,
	MaxBackups: 3,
	MaxAge:     28,
	Compress:   true,
	AddTime:    true,
}

var defaultMQTTConfig = MqttConfig{
	Broker:             "broker.emqx.io",
	ClientId:           "bms-mqtt-client-cli",
	Port:               1883,
	Topic:              "bms",
	Qos:                0,
	CleanSession:       true,
	KeepAlive:          60,
	ReconnectOnFailure: true,
	Username:           "",
	Password:           "",
}

// InitAppConfig initializes the application configuration
func InitAppConfig() (fileExists bool, err error) {
	// Check if the configuration file exists
	if utils.FileExists(appConfigFilePath) {
		return true, nil
	}

	// Create the configuration directory
	os.Mkdir(configRoot, 0o770)

	// This is just to set the default values
	GetAppConfig()

	// Save the application configuration
	err = SaveAppConfig(true)
	if err != nil {
		return false, err
	}

	return false, nil
}

// GetAppConfig returns the application configuration
func GetAppConfig() *AppConfig {
	err := loadConfig(appConfigFilePath, &appConfig)
	if err != nil {
		appConfig = &defaultAppConfig
	}
	return appConfig
}

// SaveAppConfig saves the application configuration
func SaveAppConfig(createFile bool) error {
	err := saveConfig(appConfigFilePath, appConfig, createFile)
	if err != nil {
		return err
	}

	return nil
}

// WatchAppConfigFileWithPolling watches the application configuration file for changes using polling
func WatchAppConfigFileWithPolling(callback func(), interval, debounceDuration time.Duration) {
	watchConfigFileWithPolling(appConfigFilePath, callback, interval, debounceDuration)
}
