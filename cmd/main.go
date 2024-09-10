package main

import (
	"log"
	"time"
	"timechain-device/internal/config"
	"timechain-device/internal/mqtt"
	"timechain-device/internal/sensor"
	"timechain-device/pkg/models"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	mqttClient, err := mqtt.NewClient(cfg.MQTTConfig)
	if err != nil {
		log.Fatalf("Failed to create MQTT client: %v", err)
	}
	defer mqttClient.Disconnect()

	// 将传感器配置载入，定义出传感器实例
	s := sensor.NewSensor(cfg.SensorConfig)

	for {
		if err := mqttClient.ReconnectIfNeeded(); err != nil {
			log.Printf("Failed to reconnect: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		data := s.ReadData()
		if err := mqttClient.PublishData(data); err != nil {
			log.Printf("Failed to publish data: %v", err)
		}

		status := models.DeviceStatus{
			DeviceID: cfg.SensorConfig.DeviceID,
			Status:   "online",
			LastSeen: time.Now(),
		}
		if err := mqttClient.PublishDeviceStatus(status); err != nil {
			log.Printf("Failed to publish status: %v", err)
		}

		time.Sleep(time.Duration(cfg.SensorConfig.ReadInterval) * time.Second)
	}
}
