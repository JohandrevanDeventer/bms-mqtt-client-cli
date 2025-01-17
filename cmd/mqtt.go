/*
Copyright © 2025 Johandré van Deventer <johandre.vandeventer@rubiconsa.com>
*/
package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"time"

	"github.com/JohandrevanDeventer/bms-mqtt-client-cli/internal/config"
	"github.com/JohandrevanDeventer/bms-mqtt-client-cli/pkg/text_style"
	"github.com/JohandrevanDeventer/bms-mqtt-client-cli/utils"
	"github.com/spf13/cobra"
)

var (
	mqttBroker             string
	mqttClientID           string
	mqttPort               int
	mqttTopic              string
	mqttQos                byte
	mqttCleanSession       bool
	mqttKeepAlive          int
	mqttReconnectOnFailure bool
	mqttUsername           string
	mqttPassword           string
)

// mqttCmd represents the mqtt command
var mqttCmd = &cobra.Command{
	Use:   "mqtt",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		bindMQTTFlags()
	},
}

func init() {
	rootCmd.AddCommand(mqttCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mqttCmd.PersistentFlags().String("foo", "", "A help for foo")
	mqttCmd.PersistentFlags().StringVar(&mqttBroker, "broker", "", "MQTT Broker URL")
	mqttCmd.PersistentFlags().StringVar(&mqttClientID, "client-id", "", "MQTT Client ID")
	mqttCmd.PersistentFlags().IntVar(&mqttPort, "port", 0, "MQTT Port")
	mqttCmd.PersistentFlags().StringVar(&mqttTopic, "topic", "", "MQTT Topic")
	mqttCmd.PersistentFlags().Uint8Var(&mqttQos, "qos", 0, "MQTT QoS")
	mqttCmd.PersistentFlags().BoolVar(&mqttCleanSession, "clean-session", false, "MQTT Clean Session")
	mqttCmd.PersistentFlags().IntVar(&mqttKeepAlive, "keep-alive", 0, "MQTT Keep Alive")
	mqttCmd.PersistentFlags().BoolVar(&mqttReconnectOnFailure, "reconnect-on-failure", false, "MQTT Reconnect on Failure")
	mqttCmd.PersistentFlags().StringVar(&mqttUsername, "username", "", "MQTT Username")
	mqttCmd.PersistentFlags().StringVar(&mqttPassword, "password", "", "MQTT Password")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mqttCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func bindMQTTFlags() {
	newFlag := false

	if mqttBroker != "" && mqttBroker != cfg.App.Mqtt.Broker {
		cfg.App.Mqtt.Broker = mqttBroker
		newFlag = true
	}

	if mqttClientID != "" && mqttClientID != cfg.App.Mqtt.ClientId {
		cfg.App.Mqtt.ClientId = mqttClientID
		newFlag = true
	}

	if mqttPort != 0 && mqttPort != cfg.App.Mqtt.Port {
		cfg.App.Mqtt.Port = mqttPort
		newFlag = true
	}

	if mqttTopic != "" && mqttTopic != cfg.App.Mqtt.Topic {
		cfg.App.Mqtt.Topic = mqttTopic
		newFlag = true
	}

	if mqttQos != 0 && mqttQos != cfg.App.Mqtt.Qos {
		cfg.App.Mqtt.Qos = mqttQos
		newFlag = true
	}

	if mqttCleanSession && mqttCleanSession != cfg.App.Mqtt.CleanSession {
		cfg.App.Mqtt.CleanSession = mqttCleanSession
		newFlag = true
	}

	if mqttKeepAlive != 0 && mqttKeepAlive != cfg.App.Mqtt.KeepAlive {
		cfg.App.Mqtt.KeepAlive = mqttKeepAlive
		newFlag = true
	}

	if mqttReconnectOnFailure && mqttReconnectOnFailure != cfg.App.Mqtt.ReconnectOnFailure {
		cfg.App.Mqtt.ReconnectOnFailure = mqttReconnectOnFailure
		newFlag = true
	}

	if mqttUsername != "" && mqttUsername != cfg.App.Mqtt.Username {
		cfg.App.Mqtt.Username = mqttUsername
		newFlag = true
	}

	if mqttPassword != "" && mqttPassword != cfg.App.Mqtt.Password {
		cfg.App.Mqtt.Password = mqttPassword
		newFlag = true
	}

	if newFlag {
		// Save the configuration

		fmt.Print("Updating MQTT configuration -> ")

		err := config.SaveConfig()
		if err != nil {
			var pathErr *fs.PathError
			// Check if the error is of type *fs.PathError
			if !errors.As(err, &pathErr) {
				fmt.Println(text_style.ColorText(text_style.Red, fmt.Sprintf("Cannot update MQTT configuration: %s. Restart the program with the --init flag to initialize the config files to enable changes at runtime.", err.Error())))
				os.Exit(1)
			}

		}

		time.Sleep(time.Duration(utils.GetRandomNumber(100, 500)) * time.Millisecond)

		fmt.Println(text_style.ColorText(text_style.Green, "MQTT configuration updated successfully"))

		time.Sleep(time.Duration(utils.GetRandomNumber(100, 500)) * time.Millisecond)
	}
}
