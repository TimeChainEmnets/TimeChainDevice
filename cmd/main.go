package main

import (
	"context"
	"log"
	"timechain-device/internal/config"
	"timechain-device/internal/mqtt"
	"timechain-device/internal/sensor"
)

func main() {
	// 加载服务器端配置信息
	cfg := config.LoadConfig("clientConfig.json", "brokerConfig.json")

	// 更新 MQTT 配置以匹配服务端配置
	cfg.MQTTClientConfig.MQTTConfig.DeviceInfoTopic = cfg.MQTTBrokerConfig.MQTTConfig.DeviceInfoTopic

	mqttClient, err := mqtt.NewClient(cfg.MQTTClientConfig)
	if err != nil {
		log.Fatalf("Failed to create MQTT client: %v", err)
	}
	defer mqttClient.Disconnect()

	// 初始化 GPIO 传感器
	// 假设我们使用的是 GPIO4 引脚
	gpioSensor, err := sensor.NewGPIOSensor("GPIO4")
	//if err != nil {
	//	log.Fatalf("Failed to initialize GPIO sensor: %v", err)
	//}
	//defer func(gpioSensor *sensor.GPIOSensor) {
	//	err := gpioSensor.Close()
	//	if err != nil {
	//		log.Fatalf("close gpio sensor failed %s", err)
	//	}
	//}(gpioSensor)

	// s := sensor.NewSensor(cfg.MQTTClientConfig.SensorConfig, gpioSensor)
	s := sensor.NewSensor(cfg.MQTTClientConfig.SensorConfig, gpioSensor)

	// 创建一个可取消的上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 启动周期性发布
	mqttClient.StartPeriodicPublish(ctx, s)

	// 保持主程序运行
	select {}
}
