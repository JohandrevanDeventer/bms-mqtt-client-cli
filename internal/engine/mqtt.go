package engine

import (
	"fmt"
	"strings"
	"time"

	mqttclient "github.com/JohandrevanDeventer/bms-mqtt-client-cli/pkg/mqtt"
	"go.uber.org/zap"
)

var (
	mqttStartTime time.Time
	mqttEndTime   time.Time
)

func (e *Engine) initMQTTClient() {
	e.logger.Info("Initializing MQTT client")

	e.mqttStatePersistStop()
}

func (e *Engine) connectMQTTClient() error {
	config := mqttclient.MQTTConfig{
		Broker:                e.cfg.App.Mqtt.Broker,
		Port:                  e.cfg.App.Mqtt.Port,
		ClientID:              e.cfg.App.Mqtt.ClientId,
		Topic:                 e.cfg.App.Mqtt.Topic,
		Qos:                   e.cfg.App.Mqtt.Qos,
		CleanSession:          e.cfg.App.Mqtt.CleanSession,
		KeepAlive:             e.cfg.App.Mqtt.KeepAlive,
		ReconnectOnDisconnect: e.cfg.App.Mqtt.ReconnectOnFailure,
		Username:              e.cfg.App.Mqtt.Username,
		Password:              e.cfg.App.Mqtt.Password,
	}

	e.client = mqttclient.NewMQTTClient(config)
	if err := e.client.Connect(); err != nil {
		return e.handleMqttConnectionError(err, config.Username, config.Password)
	}

	return nil
}

func (e *Engine) tryMQTTConnection(retryInterval int) {
	e.logger.Info("Attempting to connect to MQTT broker")

	if retryInterval > 60 {
		e.logger.Warn("Exceeded maximum retry interval of 60 seconds. Resetting to 60 seconds")
		retryInterval = 60
	}

	for {
		if e.client != nil && e.client.Client.IsConnected() {
			e.client.Disconnect()
		}

		if e.statePersister.Get("mqtt.status") == "connected" {
			e.mqttStatePersistStop()
		}

		if err := e.connectMQTTClient(); err != nil {
			e.logger.Error("Error connecting to MQTT broker", zap.Error(err))
			e.logger.Info("Retrying in 5 seconds")
			time.Sleep(5 * time.Second)
			continue
		}

		e.client.Subscribe()
		e.mqttStatePersistStart()
		break
	}
}

// Handle MQTT connection error
func (e *Engine) handleMqttConnectionError(err error, username, password string) error {
	if strings.Contains(err.Error(), "bad user name or password") {
		if username == "" {
			return fmt.Errorf("bad MQTT credentials detected: No username provided")
		}
		if password == "" {
			return fmt.Errorf("bad MQTT credentials detected: No password provided")
		}
		return fmt.Errorf("bad MQTT credentials detected: %w", err)
	}
	return fmt.Errorf("error connecting to MQTT broker: %w", err)
}

// mqttStatePersistStart persists the state of the MQTT connection
func (e *Engine) mqttStatePersistStart() {
	mqttStartTime = time.Now()

	e.WriteToLogFile("./connections/connections.log", fmt.Sprintf("%s: MQTT connection started\n", mqttStartTime.Format(time.RFC3339)))

	e.statePersister.Set("mqtt", map[string]interface{}{})
	e.statePersister.Set("mqtt.status", "connected")
	e.statePersister.Set("mqtt.start_time", startTime.Format(time.RFC3339))
	e.statePersister.Set("mqtt.topic", e.cfg.App.Mqtt.Topic)
	e.statePersister.Set("mqtt.client_id", e.client.Config.ClientID)
}

// mqttStatePersistStop persists the state of the MQTT connection
func (e *Engine) mqttStatePersistStop() {
	e.statePersister.Set("mqtt.status", "disconnected")

	if !mqttStartTime.IsZero() {
		mqttEndTime = time.Now()

		duration := mqttEndTime.Sub(startTime)

		e.WriteToLogFile("./connections/connections.log", fmt.Sprintf("%s: MQTT connection stopped\n", mqttEndTime.Format(time.RFC3339)))

		e.statePersister.Set("mqtt.end_time", mqttEndTime.Format(time.RFC3339))
		e.statePersister.Set("mqtt.duration", duration.String())
	}
}
