package mqtt

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"timechain-device/internal/config"
	"timechain-device/pkg/models"
)

type Client struct {
	client      mqtt.Client
	topic       string
	config      config.MQTTConfig
	isConnected bool
}

func NewClient(config config.MQTTConfig) (*Client, error) {
	opts := mqtt.NewClientOptions().
		AddBroker(config.Broker).
		SetClientID(config.ClientID).
		SetUsername(config.Username).
		SetPassword(config.Password).
		SetAutoReconnect(true).
		SetOnConnectHandler(onConnect).
		SetConnectionLostHandler(onConnectionLost)

	client := mqtt.NewClient(opts)

	mqttClient := &Client{
		client: client,
		topic:  config.Topic,
		config: config,
	}

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("failed to connect to MQTT broker: %v", token.Error())
	}

	mqttClient.isConnected = true
	return mqttClient, nil
}

func onConnect(client mqtt.Client) {
	log.Println("Connected to MQTT broker")
}

func onConnectionLost(client mqtt.Client, err error) {
	log.Printf("Connection lost to MQTT broker: %v", err)
}

func (c *Client) PublishData(data models.SensorData) error {
	if !c.isConnected {
		return fmt.Errorf("not connected to MQTT broker")
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal sensor data: %v", err)
	}

	token := c.client.Publish(c.topic, 0, false, payload)
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to publish message: %v", token.Error())
	}

	return nil
}

func (c *Client) Subscribe(handler func(mqtt.Client, mqtt.Message)) error {
	token := c.client.Subscribe(c.topic, 0, handler)
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to subscribe to topic %s: %v", c.topic, token.Error())
	}
	log.Printf("Subscribed to topic: %s", c.topic)
	return nil
}

func (c *Client) Disconnect() {
	if c.isConnected {
		c.client.Disconnect(250)
		c.isConnected = false
		log.Println("Disconnected from MQTT broker")
	}
}

func (c *Client) ReconnectIfNeeded() error {
	if !c.isConnected {
		log.Println("Attempting to reconnect to MQTT broker...")
		if token := c.client.Connect(); token.Wait() && token.Error() != nil {
			return fmt.Errorf("failed to reconnect to MQTT broker: %v", token.Error())
		}
		c.isConnected = true
		log.Println("Reconnected to MQTT broker")
	}
	return nil
}

func (c *Client) PublishDeviceStatus(status models.DeviceStatus) error {
	payload, err := json.Marshal(status)
	if err != nil {
		return fmt.Errorf("failed to marshal device status: %v", err)
	}

	statusTopic := fmt.Sprintf("%s/status", c.topic)
	token := c.client.Publish(statusTopic, 1, false, payload)
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to publish device status: %v", token.Error())
	}

	return nil
}
