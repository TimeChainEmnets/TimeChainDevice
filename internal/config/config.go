package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	MQTTConfig   MQTTConfig   `json:"mqtt_config"`
	SensorConfig SensorConfig `json:"sensor_config"`
}

type MQTTConfig struct {
	Broker   string `json:"broker"`
	ClientID string `json:"client_id"`
	Topic    string `json:"topic"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type SensorConfig struct {
	DeviceID     string  `json:"device_id"`
	DeviceType   string  `json:"device_type"`
	ReadInterval int     `json:"read_interval"` // 读取间隔（秒）
	MinValue     float64 `json:"min_value"`     // 传感器最小值
	MaxValue     float64 `json:"max_value"`     // 传感器最大值
	Unit         string  `json:"unit"`          // 测量单位
}

func Load() (*Config, error) {
	file, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
