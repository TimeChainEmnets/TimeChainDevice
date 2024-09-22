package sensor

import (
	"fmt"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/host/v3"
)

type GPIOSensor struct {
	pin gpio.PinIO
}

func NewGPIOSensor(pinName string) (*GPIOSensor, error) {
	// 初始化 periph 主机
	if _, err := host.Init(); err != nil {
		return nil, fmt.Errorf("failed to initialize periph: %v", err)
	}

	// 获取指定的 GPIO 引脚
	pin := gpioreg.ByName(pinName)
	if pin == nil {
		return nil, fmt.Errorf("failed to find pin %s", pinName)
	}

	// 配置引脚为输入模式
	if err := pin.In(gpio.PullUp, gpio.NoEdge); err != nil {
		return nil, fmt.Errorf("failed to configure pin as input: %v", err)
	}

	return &GPIOSensor{pin: pin}, nil
}

func (s *GPIOSensor) Read() (float64, error) {
	// 读取 GPIO 状态
	// 这里假设我们使用的是数字传感器，返回 0 或 1
	// 如果是模拟传感器，可能需要使用 ADC（模数转换器）
	value := 0.0
	if s.pin.Read() == gpio.High {
		value = 1.0
	}

	return value, nil
}

func (s *GPIOSensor) Close() error {
	// GPIO 引脚通常不需要特别的关闭操作
	return nil
}
