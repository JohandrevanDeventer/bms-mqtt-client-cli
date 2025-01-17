package config

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/JohandrevanDeventer/bms-mqtt-client-cli/pkg/text_style"
	"github.com/JohandrevanDeventer/bms-mqtt-client-cli/utils"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var (
	debounceTimer *time.Timer
	mu            sync.Mutex
	lastHash      uint64
	debounceMu    sync.Mutex
)

// Config represents the application configuration
func InitConfig() (newFiles []string, existingFiles []string, err error) {
	// Initialize the system configuration file and save it to a file
	systemCfgExists := false
	systemCfgExists, err = InitSystemConfig()
	if systemCfgExists {
		existingFiles = append(existingFiles, systemConfigFilePath)
	} else {
		newFiles = append(newFiles, systemConfigFilePath)
	}

	if err != nil {
		return newFiles, existingFiles, err
	}

	// Initialize the application configuration file and save it to a file
	appCfgExists := false
	appCfgExists, err = InitAppConfig()
	if appCfgExists {
		existingFiles = append(existingFiles, appConfigFilePath)
	} else {
		newFiles = append(newFiles, appConfigFilePath)
	}

	if err != nil {
		return newFiles, existingFiles, err
	}

	return newFiles, existingFiles, nil
}

// GetConfig returns the application configuration
func GetConfig() *Config {
	return &Config{
		PersistFilePath: persistFilePath,
		Flags:           GetFlagsConfig(),
		System:          GetSystemConfig(),
		App:             GetAppConfig(),
	}
}

// SaveConfig saves the configuration
func SaveConfig() error {
	err := SaveAppConfig(false)
	if err != nil {
		return err
	}

	return nil
}

// loadConfig loads the configuration from a file
func loadConfig(path string, target interface{}) error {
	// Create a new viper instance
	v := viper.New()

	// Check for file extension
	ext := filepath.Ext(path)
	if ext == "" {
		return fmt.Errorf("config file must have an extension")
	}

	// Set the config file
	v.SetConfigFile(path)

	// Read the config file
	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}

	// Unmarshal the configuration file into the struct
	if err := v.Unmarshal(&target); err != nil {
		return fmt.Errorf("error unmarshalling config file: %w", err)
	}

	return nil
}

// saveConfig saves the configuration to a file
func saveConfig(path string, config interface{}, createFile bool) error {
	// Check if the file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if !createFile {
			return fmt.Errorf("file does not exist: %s", path)
		}
	}

	// Open the file for writing
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Encode the configuration to the file
	encoder := yaml.NewEncoder(file)
	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}

	return nil
}

// PrintInfo prints the application information
func PrintInfo(versionOnly bool) {
	flagsCfg := GetConfig().Flags
	systemCfg := GetConfig().System

	goVersion := strings.Replace(runtime.Version(), "go", "", 1)

	fullVersion := fmt.Sprintf("v%s-%d", systemCfg.AppVersion, systemCfg.BuildNumber)

	fmt.Printf("Welcome to %s!\n", text_style.ColorText(text_style.Green, (text_style.BoldText(systemCfg.AppName))))
	fmt.Printf("Built with Go %s\n", text_style.ColorText(text_style.Yellow, (text_style.BoldText(goVersion))))
	fmt.Printf("Running version %s\n", text_style.ColorText(text_style.Magenta, (text_style.BoldText(fullVersion))))

	contributors := strings.Join(systemCfg.Contributors, ", ")
	fmt.Printf("Developed by %s\n", text_style.ColorText(text_style.Cyan, (text_style.BoldText(contributors))))
	fmt.Printf("Release date: %s\n", text_style.ColorText(text_style.Blue, (text_style.BoldText(systemCfg.ReleaseDate))))

	fmt.Println("")

	if !versionOnly {
		time.Sleep(time.Duration(utils.GetRandomNumber(500, 2000)) * time.Millisecond)

		switch strings.ToLower(flagsCfg.Environment) {
		case "development":
			fmt.Println(text_style.ColorText(text_style.Red, (text_style.BoldText("Running in Development mode"))))
		case "testing":
			fmt.Println(text_style.ColorText(text_style.Yellow, (text_style.BoldText("Running in Testing mode"))))
		case "production":
			fmt.Println(text_style.ColorText(text_style.Green, (text_style.BoldText("Running in Production mode"))))
		default:
			fmt.Println(text_style.ColorText(text_style.Blue, (text_style.BoldText("Running in Default mode"))))
		}

		fmt.Println("")

		if flagsCfg.DebugMode {
			fmt.Println(text_style.ColorText(text_style.Red, (text_style.BoldText("Debug mode enabled"))))

			fmt.Println("")
		}

	}
}

func WatchConfig(path string, callback func(), debounce int) error {
	v := viper.New()

	// Check for file extension
	ext := filepath.Ext(path)
	if ext == "" {
		return fmt.Errorf("config file must have an extension")
	}

	// Set the config file
	v.SetConfigFile(path)

	debounceDuration := time.Duration(debounce) * time.Millisecond

	v.OnConfigChange(func(e fsnotify.Event) {
		if debounceTimer != nil {
			debounceTimer.Stop()
		}

		debounceTimer = time.AfterFunc(debounceDuration, func() {
			callback()
		})
	})

	v.WatchConfig()

	return nil
}

func watchConfigFileWithPolling(path string, callback func(), interval, debounceDuration time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		func() {
			mu.Lock()
			defer mu.Unlock()

			// Open and hash the file
			file, err := os.Open(path)
			if err != nil {
				// Log or handle file open error as needed
				return
			}
			defer file.Close()

			hash := fnv.New64a()
			if _, err := io.Copy(hash, file); err != nil {
				// Log or handle file read error as needed
				return
			}

			currentHash := hash.Sum64()
			if lastHash != currentHash {
				lastHash = currentHash

				// Handle debounce logic
				debounceMu.Lock()
				defer debounceMu.Unlock()

				if debounceTimer != nil {
					if !debounceTimer.Stop() {
						// Drain the channel only if the timer has already fired
						select {
						case <-debounceTimer.C:
						default:
						}
					}
				}

				// Set a new debounce timer
				debounceTimer = time.AfterFunc(debounceDuration, func() {
					// Ensure callback execution is thread-safe
					go callback()
				})
			}
		}()
	}
}

// Clone function to create a deep copy of the config in TOML format
func CloneConfig(cfg *Config) (*Config, error) {
	// Serialize the current config to TOML format
	var buf bytes.Buffer
	err := yaml.NewEncoder(&buf).Encode(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to encode config to TOML: %w", err)
	}

	// Now, deserialize the TOML back into a new config struct
	var clonedConfig Config
	err = yaml.Unmarshal(buf.Bytes(), &clonedConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal cloned config: %w", err)
	}

	return &clonedConfig, nil
}
