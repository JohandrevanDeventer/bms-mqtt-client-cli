package config

import (
	"fmt"
	"os"

	"github.com/JohandrevanDeventer/bms-mqtt-client-cli/utils"
)

var systemConfigFilePath = fmt.Sprintf("%s/%s", configRoot, systemConfigFile)

var systemConfig *SystemConfig

var defaultSystemConfig = SystemConfig{
	AppName:      "Rubicon BMS MQTT Client CLI",
	AppVersion:   "1.0.0",
	BuildNumber:  1,
	ReleaseDate:  "2025-01-01",
	Contributors: []string{"Johandré van Deventer"},
}

// InitSystemConfig initializes the system configuration
func InitSystemConfig() (fileExists bool, err error) {
	// Check if the config directory exists, if not create it
	if utils.FileExists(systemConfigFilePath) {
		return true, nil
	}

	// Create the configuration directory
	os.Mkdir(configRoot, 0o770)

	// This is just to set the default values
	GetSystemConfig()

	// Save the system configuration
	err = SaveSystemConfig(true)
	if err != nil {
		return false, err
	}

	return false, nil
}

// GetSystemConfig returns the system configuration
func GetSystemConfig() *SystemConfig {
	err := loadConfig(systemConfigFilePath, &systemConfig)
	if err != nil {
		systemConfig = &defaultSystemConfig
	}
	return systemConfig
}

// SaveSystemConfig saves the system configuration
func SaveSystemConfig(createFile bool) error {
	err := saveConfig(systemConfigFilePath, systemConfig, createFile)
	if err != nil {
		return err
	}

	return nil
}

// IncrementBuildNumber increments the build number
func IncrementBuildNumber() {
	systemConfig.BuildNumber++
	_ = SaveSystemConfig(false)
}
