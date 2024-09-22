package models

import (
	"time"
)

// SensorData 代表从传感器读取的单个数据点
type SensorData struct {
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
}

// BatchData 代表一批传感器数据，可用于批量发送
type BatchData struct {
	DeviceID   string       `json:"device_id"`
	Unit       string       `json:"unit"`
	DataPoints []SensorData `json:"data_points"`
}

// DeviceStatus 表示设备的当前状态
type DeviceStatus struct {
	DeviceID     string    `json:"device_id"`
	Status       string    `json:"status"` // 例如: "online", "offline", "error"
	LastSeen     time.Time `json:"last_seen"`
	BatteryLevel float64   `json:"battery_level,omitempty"` // 如果适用
}
