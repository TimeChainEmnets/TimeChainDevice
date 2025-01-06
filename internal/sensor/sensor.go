package sensor

import (
	"context"
	"log"
	"math/rand"
	"time"
	"github.com/TimeChainEmnets/TimeChainDevice/internal/config"
	"github.com/TimeChainEmnets/TimeChainDevice/pkg/models"
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

	// value, err := s.reader.Read()
	// log.Print(value)
	return newValue
}

func (s *Sensor) ReadData() (models.SensorData, error) {
	value := s.readSensorValue()

	return models.SensorData{
		SensorID: s.config.SensorID,
		DeviceID: s.config.DeviceID,
		Timestamp: time.Now().Unix(),
		Type:	 s.config.SensorType,
		Value:     value,
		Unit:      s.config.Unit,
	}, nil
}

func (s *Sensor) CollectBatchData(ctx context.Context) <-chan models.SensorData {
	out := make(chan models.SensorData)
	go func() {
		defer close(out)
		ticker := time.NewTicker(time.Duration(s.config.SampleRate) * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				data, err := s.ReadData()
				if err != nil {
					log.Fatalf("Read data from sensor failed! %s", err)
				}
				out <- data
			}
		}
	}()
	return out
}
