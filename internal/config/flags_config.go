package config

import (
	"fmt"
	"os"

	"github.com/JohandrevanDeventer/bms-mqtt-client-cli/utils"
)

var flagsConfigFilePath = fmt.Sprintf("%s/%s", configRoot, flagsConfigFile)

var flagsConfig *FlagsConfig

var defaultFlagsConfig = FlagsConfig{
	Environment: "production",
	DebugMode:   false,
}

// InitFlagsConfig initializes the flags configuration
func InitFlagsConfig() (fileExists bool) {
	// Check if the config directory exists. If not, create it
	if utils.FileExists(flagsConfigFilePath) {
		return true
	}

	// Create the config directory
	os.Mkdir(configRoot, 0o770)

	// This is just to set the default values
	GetFlagsConfig()

	// We don't need to save the flags configuration

	return false
}

// GetFlagsConfig returns the flags configuration
func GetFlagsConfig() *FlagsConfig {
	err := loadConfig(flagsConfigFilePath, &flagsConfig)
	if err != nil {
		flagsConfig = &defaultFlagsConfig
	}
	return flagsConfig
}

// SaveFlagsConfig saves the flags configuration
func SaveFlagsConfig(createFile bool) error {
	err := saveConfig(flagsConfigFilePath, flagsConfig, createFile)
	if err != nil {
		return err
	}

	return nil
}
