package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	MQTTClientConfig MQTTClientConfig `json:"mqtt_client_config"`
	MQTTBrokerConfig MQTTBrokerConfig `json:"mqtt_broker_config"`
}

type MQTTBrokerConfig struct {
	MQTTConfig struct {
		Address         string `json:"address"`
		DeviceInfoTopic string `json:"device_info_topic"`
		Port            int    `json:"port"`
	} `json:"mqtt_config"`
	DeviceConfig struct {
		ScanInterval int `json:"scan_interval"`
	} `json:"device_config"`
}

type MQTTClientConfig struct {
	MQTTConfig struct {
		Broker          string `json:"broker"`
		ClientID        string `json:"client_id"`
		DeviceInfoTopic string `json:"device_info_topic"`
		Username        string `json:"username"`
		Password        string `json:"password"`
	} `json:"mqtt_config"`
	SensorConfig SensorConfig `json:"sensor_config"`
}

type SensorConfig struct {
	DeviceID        string  `json:"device_id"`
	SensorID        string  `json:"sensor_id"`
	SensorType      string  `json:"sensor_type"`
	SampleRate      int     `json:"sample_rate"`
	PublishInterval int     `json:"publish_interval"`
	BatchSize       int     `json:"batch_size"`
	MinValue        float64 `json:"min_value"` // 传感器最小值
	MaxValue        float64 `json:"max_value"` // 传感器最大值
	Unit            string  `json:"unit"`      // 测量单位
}

func LoadConfig(cliCfgFileName string, brokerCfgFileName string) *Config {
	cliCfg, err := LoadClientConfig(cliCfgFileName)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	brokerCfg, err := LoadBrokerConfig(brokerCfgFileName)
	if err != nil {
		log.Fatalf("Failed to create MQTT client: %v", err)
	}
	cfg := Config{*cliCfg, *brokerCfg}
	return &cfg
}

func LoadClientConfig(fileName string) (*MQTTClientConfig, error) {
	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	// 构造配置文件的相对路径
	configPath := filepath.Join(currentDir, fileName)
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config MQTTClientConfig
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// LoadServerConfig 函数用于加载 JSON 配置文件
func LoadBrokerConfig(fileName string) (*MQTTBrokerConfig, error) {
	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	// 构造配置文件的相对路径
	configPath := filepath.Join(currentDir, fileName)

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config MQTTBrokerConfig
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
