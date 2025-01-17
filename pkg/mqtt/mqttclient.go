package mqttclient

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/JohandrevanDeventer/bms-mqtt-client-cli/pkg/logging"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var logger *zap.Logger

// MQTTConfig is the configuration for the MQTT client
type MQTTConfig struct {
	Broker                string
	Port                  int
	ClientID              string
	Topic                 string
	Qos                   byte
	CleanSession          bool
	KeepAlive             int
	ReconnectOnDisconnect bool
	Username              string
	Password              string
}

// MQTTClient is the interface for the MQTT client
type MQTTClient struct {
	mu     sync.Mutex
	Client mqtt.Client
	Config MQTTConfig
	ctx    context.Context
	cancel context.CancelFunc
}

func NewMQTTClient(config MQTTConfig) *MQTTClient {
	logger = logging.GetLogger("mqtt")
	ctx, cancel := context.WithCancel(context.Background())

	// Generate a new ClientID
	newClientId := generateClientID(config.ClientID)
	config.ClientID = newClientId

	return &MQTTClient{
		Config: config,
		ctx:    ctx,
		cancel: cancel,
	}
}

// generateClientID creates a random ClientID.
func generateClientID(baseID string) string {
	uuidPart := strings.Split(uuid.New().String(), "-")[0] // Extract the first part of the UUID
	return fmt.Sprintf("%s-%s", baseID, uuidPart)
}

func (m *MQTTClient) Connect() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	logger.Info("Connecting to MQTT broker", zap.String("broker", m.Config.Broker), zap.Int("port", m.Config.Port))
	logger.Debug("MQTT client configuration", zap.String("client_id", m.Config.ClientID), zap.String("topic", m.Config.Topic), zap.Uint8("qos", m.Config.Qos), zap.Bool("clean_session", m.Config.CleanSession), zap.Int("keep_alive", m.Config.KeepAlive))

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", m.Config.Broker, m.Config.Port))
	opts.SetClientID(m.Config.ClientID)
	opts.SetCleanSession(m.Config.CleanSession)
	opts.SetKeepAlive(time.Duration(m.Config.KeepAlive) * time.Second)
	opts.SetUsername(m.Config.Username)
	opts.SetPassword(m.Config.Password)

	opts.OnConnect = m.onConnect
	opts.OnConnectionLost = m.onConnectionLost

	m.Client = mqtt.NewClient(opts)
	token := m.Client.Connect()
	if token.Wait() && token.Error() != nil {
		logger.Error("Error connecting to MQTT broker", zap.Error(token.Error()))
		return fmt.Errorf("error connecting to MQTT broker: %v", token.Error())
	}

	return nil
}

func (m *MQTTClient) Disconnect() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.Client != nil && m.Client.IsConnected() {
		m.Client.Disconnect(250)
		logger.Info("Disconnected from MQTT broker")
	}

	m.cancel()
}

func (m *MQTTClient) Subscribe() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.Client == nil || !m.Client.IsConnected() {
		return fmt.Errorf("client is not connected")
	}

	token := m.Client.Subscribe(m.Config.Topic, byte(m.Config.Qos), m.onMessage)
	token.Wait()
	logger.Info("Subscribed to topic", zap.String("topic", m.Config.Topic))
	return token.Error()
}

func (m *MQTTClient) onConnect(client mqtt.Client) {
	logger.Info("Connected to MQTT broker")
}

func (m *MQTTClient) onConnectionLost(client mqtt.Client, err error) {
	logger.Error("Connection lost. Attempting to reconnect...", zap.Error(err))

	for {
		select {
		case <-m.ctx.Done():
			logger.Warn("Context canceled, stopping reconnection attempts")
			return
		default:
			if err := m.Connect(); err != nil {
				logger.Warn("Reconnection failed. Retrying...", zap.Error(err))
			} else {
				logger.Info("Reconnected to MQTT broker")
				return
			}
			time.Sleep(5 * time.Second) // Wait before retrying
		}
	}
}

func (m *MQTTClient) onMessage(client mqtt.Client, msg mqtt.Message) {
	topic := msg.Topic()

	logger.Info("Received message", zap.Uint16("message_id", msg.MessageID()), zap.String("topic", topic))
	logger.Debug("Message payload", zap.Uint16("message_id", msg.MessageID()), zap.String("payload", string(msg.Payload())))
}
