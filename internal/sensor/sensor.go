package sensor

import (
	"context"
	"log"
	"math/rand"
	"time"
	"timechain-device/internal/config"
	"timechain-device/pkg/models"
)

type SensorReader interface {
	Read() (float64, error)
}

type Sensor struct {
	config     config.SensorConfig
	lastValue  float64
	trend      float64
	noiseLevel float64
	reader     SensorReader
}

func NewSensor(config config.SensorConfig) *Sensor {
	// , reader SensorReader
	midValue := (config.MinValue + config.MaxValue) / 2
	return &Sensor{
		config:     config,
		lastValue:  midValue,
		trend:      0,
		noiseLevel: (config.MaxValue - config.MinValue) * 0.01, // 1% of range as noise
		// reader:     reader,
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
func (s *Sensor) ReadData() (models.SensorData, error) {
	//value, err := s.reader.Read()
	//if err != nil {
	//	return models.SensorData{}, err
	//}
	//
	//// 可以在这里添加一些数据验证或处理
	//if value < s.config.MinValue {
	//	value = s.config.MinValue
	//} else if value > s.config.MaxValue {
	//	value = s.config.MaxValue
	//}
	//
	//s.lastValue = value
	//
	//return models.SensorData{
	//	Timestamp: time.Now(),
	//	Value:     value,
	//}, nil
	return models.SensorData{
		Timestamp: time.Now(),          //  time.Now().Unix(),
		Value:     s.readSensorValue(), // 实际读取传感器值的函数
	}, nil
}

func (s *Sensor) CollectBatchData(ctx context.Context) <-chan models.BatchData {
	out := make(chan models.BatchData)
	go func() {
		defer close(out)
		ticker := time.NewTicker(time.Duration(s.config.SampleRate) * time.Second)
		defer ticker.Stop()

		var dataPoints []models.SensorData
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				data, err := s.ReadData()
				if err != nil {
					log.Fatalf("Read data from sensor failed! %s", err)
				}
				dataPoints = append(dataPoints, data)

				if len(dataPoints) >= s.config.BatchSize {
					batchData := models.BatchData{
						DeviceID:   s.config.DeviceID,
						Unit:       s.config.Unit,
						DataPoints: dataPoints,
					}
					out <- batchData
					dataPoints = nil // Reset the slice
				}
			}
		}
	}()
	return out
}
