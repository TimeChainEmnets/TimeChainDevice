package models

import (
	"time"
)

// SensorData 代表从传感器读取的单个数据点
type SensorData struct {
	DeviceID   string    `json:"device_id"`
	DeviceType string    `json:"device_type"`
	Timestamp  time.Time `json:"timestamp"`
	Value      float64   `json:"value"`
	Unit       string    `json:"unit"`
}

// BatchData 代表一批传感器数据，可用于批量发送
type BatchData struct {
	DeviceID   string       `json:"device_id"`
	DataPoints []SensorData `json:"data_points"`
}

// DeviceStatus 表示设备的当前状态
type DeviceStatus struct {
	DeviceID     string    `json:"device_id"`
	Status       string    `json:"status"` // 例如: "online", "offline", "error"
	LastSeen     time.Time `json:"last_seen"`
	BatteryLevel float64   `json:"battery_level,omitempty"` // 如果适用
}

// DeviceConfig 表示设备的配置信息
type DeviceConfig struct {
	DeviceID       string  `json:"device_id"`
	SampleRate     int     `json:"sample_rate"`         // 采样率（每分钟采样次数）
	ReportInterval int     `json:"report_interval"`     // 报告间隔（秒）
	Threshold      float64 `json:"threshold,omitempty"` // 数据报告阈值，如果适用
}
