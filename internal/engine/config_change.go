package engine

import (
	"strings"
	"time"

	"github.com/JohandrevanDeventer/bms-mqtt-client-cli/internal/config"
	"github.com/JohandrevanDeventer/bms-mqtt-client-cli/pkg/logging"
	"go.uber.org/zap"
)

func (e *Engine) appConfigChangeCallback() {
	oldCfg, err := config.CloneConfig(e.cfg)
	if err != nil {
		e.logger.Error("failed to clone config", zap.Error(err))
		return
	}

	newCfg := config.GetConfig()

	e.handleAppConfigChanged(oldCfg, newCfg)
}

// General helper for comparing simple fields
func (e *Engine) hasConfigFieldChanged(oldVal, newVal interface{}) bool {
	return oldVal != newVal
}

func (e *Engine) handleAppConfigChanged(oldCfg, newCfg *config.Config) {

	// Handle changes for the Logging config
	if e.hasLoggingConfigChanged(oldCfg.App.Logging, newCfg.App.Logging) {
		e.logger.Warn("Logging configuration changed, but the logger only supports changing the log level at runtime. Please restart the application to apply the changes.")
	}

	// Handle changes to the log level
	if e.hasConfigFieldChanged(oldCfg.App.Logging.Level, newCfg.App.Logging.Level) {
		e.handleLogLevelChange(oldCfg.App.Logging.Level, newCfg.App.Logging.Level)
	}

	if e.hasMQTTConfigChanged(oldCfg.App.Mqtt, newCfg.App.Mqtt) {
		e.handleMQTTConfigChanged(oldCfg, newCfg)
	}
}

// ========================================= Logging =============================================================

// General function for comparing complex struct fields (e.g., Logging struct)
func (e *Engine) hasLoggingConfigChanged(oldLogging, newLogging config.LoggingConfig) bool {
	return oldLogging.FilePath != newLogging.FilePath ||
		oldLogging.MaxSize != newLogging.MaxSize ||
		oldLogging.MaxBackups != newLogging.MaxBackups ||
		oldLogging.MaxAge != newLogging.MaxAge ||
		oldLogging.Compress != newLogging.Compress ||
		oldLogging.AddTime != newLogging.AddTime
}

// Handle log level changes
func (e *Engine) handleLogLevelChange(oldLevel, newLevel string) {
	// Check for invalid log level
	if newLevel == "" {
		e.logger.Warn("Logging level cannot be empty. Please provide a valid logging level.")
		return
	}

	// Apply the new log level if valid
	if e.cfg.Flags.DebugMode {
		e.logger.Warn("Debug mode is enabled. The logging level cannot be changed at runtime.")
		return
	}

	// Check if the new log level is valid
	if !logging.IsValidLogLevel(newLevel) {
		e.logger.Warn("Logging level changed in the config, but invalid log level provided. Valid log levels: 'debug', 'info', 'warn', 'error', 'dpanic', 'panic', 'fatal'", zap.String("log_level", newLevel))
		return
	}

	if strings.EqualFold(newLevel, "debug") || strings.EqualFold(newLevel, "info") {
		err := logging.SetLogLevel(newLevel)
		if err != nil {
			e.logger.Error("failed to set log level", zap.Error(err))
			return
		}

		e.logger.Info("Logging level changed", zap.String("old_level", oldLevel), zap.String("new_level", newLevel))
	} else {
		e.logger.Info("Logging level changed", zap.String("old_level", oldLevel), zap.String("new_level", newLevel))
		err := logging.SetLogLevel(newLevel)
		if err != nil {
			e.logger.Error("failed to set log level", zap.Error(err))
			return
		}
	}
}

// ========================================= MQTT =============================================================

func (e *Engine) hasMQTTConfigChanged(oldMQTT, newMQTT config.MqttConfig) bool {
	return oldMQTT.Broker != newMQTT.Broker ||
		oldMQTT.Port != newMQTT.Port ||
		oldMQTT.ClientId != newMQTT.ClientId ||
		oldMQTT.Topic != newMQTT.Topic ||
		oldMQTT.Qos != newMQTT.Qos ||
		oldMQTT.CleanSession != newMQTT.CleanSession ||
		oldMQTT.KeepAlive != newMQTT.KeepAlive ||
		oldMQTT.ReconnectOnFailure != newMQTT.ReconnectOnFailure ||
		oldMQTT.Username != newMQTT.Username ||
		oldMQTT.Password != newMQTT.Password
}

func (e *Engine) handleMQTTConfigChanged(oldCfg, newCfg *config.Config) {
	if oldCfg.App.Mqtt.Broker != newCfg.App.Mqtt.Broker {
		e.logger.Debug("MQTT broker changed", zap.String("old_broker", oldCfg.App.Mqtt.Broker), zap.String("new_broker", newCfg.App.Mqtt.Broker))
	}

	if oldCfg.App.Mqtt.Port != newCfg.App.Mqtt.Port {
		e.logger.Debug("MQTT port changed", zap.Int("old_port", oldCfg.App.Mqtt.Port), zap.Int("new_port", newCfg.App.Mqtt.Port))
	}

	if oldCfg.App.Mqtt.ClientId != newCfg.App.Mqtt.ClientId {
		e.logger.Debug("MQTT client ID changed", zap.String("old_client_id", oldCfg.App.Mqtt.ClientId), zap.String("new_client_id", newCfg.App.Mqtt.ClientId))
	}

	if oldCfg.App.Mqtt.Topic != newCfg.App.Mqtt.Topic {
		e.logger.Debug("MQTT topic changed", zap.String("old_topic", oldCfg.App.Mqtt.Topic), zap.String("new_topic", newCfg.App.Mqtt.Topic))
	}

	if oldCfg.App.Mqtt.Qos != newCfg.App.Mqtt.Qos {
		e.logger.Debug("MQTT QoS changed", zap.Uint8("old_qos", oldCfg.App.Mqtt.Qos), zap.Uint8("new_qos", newCfg.App.Mqtt.Qos))
	}

	if oldCfg.App.Mqtt.CleanSession != newCfg.App.Mqtt.CleanSession {
		e.logger.Debug("MQTT clean session changed", zap.Bool("old_clean_session", oldCfg.App.Mqtt.CleanSession), zap.Bool("new_clean_session", newCfg.App.Mqtt.CleanSession))
	}

	if oldCfg.App.Mqtt.KeepAlive != newCfg.App.Mqtt.KeepAlive {
		e.logger.Debug("MQTT keep alive changed", zap.Int("old_keep_alive", oldCfg.App.Mqtt.KeepAlive), zap.Int("new_keep_alive", newCfg.App.Mqtt.KeepAlive))
	}

	if oldCfg.App.Mqtt.ReconnectOnFailure != newCfg.App.Mqtt.ReconnectOnFailure {
		e.logger.Debug("MQTT reconnect on failure changed", zap.Bool("old_reconnect_on_failure", oldCfg.App.Mqtt.ReconnectOnFailure), zap.Bool("new_reconnect_on_failure", newCfg.App.Mqtt.ReconnectOnFailure))
	}

	if oldCfg.App.Mqtt.Username != newCfg.App.Mqtt.Username {
		e.logger.Debug("MQTT username changed", zap.String("old_username", oldCfg.App.Mqtt.Username), zap.String("new_username", newCfg.App.Mqtt.Username))
	}

	if oldCfg.App.Mqtt.Password != newCfg.App.Mqtt.Password {
		e.logger.Debug("MQTT password changed", zap.String("old_password", oldCfg.App.Mqtt.Password), zap.String("new_password", newCfg.App.Mqtt.Password))
	}

	e.logger.Debug("MQTT configuration changed. Restarting MQTT connection")
	e.client.Disconnect()
	time.Sleep(1000 * time.Millisecond)
	e.tryMQTTConnection(5)
}
