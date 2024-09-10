package sensor

import (
	"math/rand"
	"time"
	"timechain-device/internal/config"
	"timechain-device/pkg/models"
)

type Sensor struct {
	config     config.SensorConfig
	lastValue  float64
	trend      float64
	noiseLevel float64
}

func NewSensor(config config.SensorConfig) *Sensor {
	midValue := (config.MinValue + config.MaxValue) / 2
	return &Sensor{
		config:     config,
		lastValue:  midValue,
		trend:      0,
		noiseLevel: (config.MaxValue - config.MinValue) * 0.01, // 1% of range as noise
	}
}

func (s *Sensor) readSensorValue() float64 {
	// 更新趋势
	s.trend += rand.Float64()*0.2 - 0.1 // 趋势在 -0.1 到 0.1 之间随机变化

	// 添加一些随机噪声
	noise := (rand.Float64()*2 - 1) * s.noiseLevel

	// 计算新值
	newValue := s.lastValue + s.trend + noise

	// 确保值在有效范围内
	if newValue < s.config.MinValue {
		newValue = s.config.MinValue
		s.trend = 0 // 重置趋势
	} else if newValue > s.config.MaxValue {
		newValue = s.config.MaxValue
		s.trend = 0 // 重置趋势
	}

	s.lastValue = newValue
	return newValue
}

// 在传感器读取函数中
func (s *Sensor) ReadData() models.SensorData {
	return models.SensorData{
		DeviceID:   s.config.DeviceID,
		DeviceType: s.config.DeviceType,
		Timestamp:  time.Now(),          //  time.Now().Unix(),
		Value:      s.readSensorValue(), // 实际读取传感器值的函数
		Unit:       "celsius",           // 或其他适当的单位
	}
}
