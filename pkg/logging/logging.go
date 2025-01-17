package logging

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"strings"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger        *zap.Logger
	logLevel      = zap.NewAtomicLevel()
	loggingLogger *zap.Logger
	loggerConfig  *LoggingConfig
)

type LoggingConfig struct {
	Level      string `mapstructure:"level"`
	FilePath   string `mapstructure:"file_path"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
	DebugMode  bool   `mapstructure:"debug_mode"`
	AddTime    bool   `mapstructure:"add_time"`
}

// NewLoggingConfig creates a new logging configuration.
func NewLoggingConfig(level, filePath string, maxSize, maxBackups, maxAge int, compress, debugMode, addTime bool) *LoggingConfig {
	return &LoggingConfig{
		Level:      level,
		FilePath:   filePath,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		Compress:   compress,
		DebugMode:  debugMode,
		AddTime:    addTime,
	}
}

func GetLogger(name string) *zap.Logger {
	if logger == nil {
		log.Fatal("logger is not initialized")
	}

	return logger.Named(name)
}

// NewLogger initializes a zap.Logger instance if it has not been initialized
// already and returns the same instance for subsequent calls.
func NewLogger(cfg *LoggingConfig) *zap.Logger {
	stdout := zapcore.AddSync(os.Stdout)

	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   cfg.FilePath,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	})

	level := zap.InfoLevel

	if cfg.DebugMode {
		level = zap.DebugLevel
	}

	cfgLevel := cfg.Level
	if cfgLevel != "" && !cfg.DebugMode {
		levelFromEnv, err := zapcore.ParseLevel(cfgLevel)
		if err != nil {
			log.Println(
				fmt.Errorf("invalid level, defaulting to INFO: %w", err),
			)
		}

		level = levelFromEnv
	}

	logLevel.SetLevel(level)

	productionCfg := zap.NewProductionEncoderConfig()

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	if cfg.AddTime {
		productionCfg.TimeKey = "timestamp"
		productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		productionCfg.TimeKey = ""  // Remove timestamp
		developmentCfg.TimeKey = "" // Remove timestamp
	}

	fileEncoder := zapcore.NewJSONEncoder(productionCfg)
	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)

	var gitRevision string

	buildInfo, ok := debug.ReadBuildInfo()
	if ok {
		for _, v := range buildInfo.Settings {
			if v.Key == "vcs.revision" {
				gitRevision = v.Value
				break
			}
		}
	}

	// log to multiple destinations (console and file)
	// extra fields are added to the JSON output alone
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, logLevel),
		zapcore.NewCore(fileEncoder, file, logLevel).
			With(
				[]zapcore.Field{
					zap.String("git_revision", gitRevision),
					zap.String("go_version", buildInfo.GoVersion),
				},
			),
	)

	logger = zap.New(core)

	loggerConfig = cfg
	loggingLogger = logger.Named("main")

	return logger
}

func IsValidLogLevel(level string) bool {
	switch strings.ToLower(level) {
	case "debug", "info", "warn", "error", "dpanic", "panic", "fatal":
		return true
	}

	return false
}

// SetLogLevel dynamically updates the log level at runtime
func SetLogLevel(level string) error {
	if loggerConfig.DebugMode {
		loggingLogger.Warn("Debug mode is enabled. Log level cannot be changed at runtime.")
		return nil
	}

	parsedLevel, err := zapcore.ParseLevel(level)
	if err != nil {
		return fmt.Errorf("invalid log level: %s", level)
	}

	// oldLevel := logLevel.Level()
	// loggingLogger.Debug("Changing log level", zap.String("old_level", oldLevel.String()), zap.String("new_level", parsedLevel.String()))
	logLevel.SetLevel(parsedLevel)
	return nil
}
