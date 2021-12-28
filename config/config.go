package config

import "github.com/jothflee/bluetooth"

type Device struct {
	MAC    string            `yaml:"mac"`
	Name   string            `yaml:"name"`
	Type   DeviceType        `yaml:"type"`
	Device *bluetooth.Device `yaml:"-"`
}

type DeviceType string

const (
	DeviceTypeIN  DeviceType = "in"
	DeviceTypeOUT DeviceType = "out"
)

type Config struct {
	Devices map[string]*Device `yaml:"devices"`
}

func (c *Config) AddDevice(d *Device) {
	c.Devices[d.MAC] = d
}
func (d *Device) Disconnect() {
	if d.Device != nil {
		d.Device.Disconnect()
		d.Device = nil
	}
}

func NewConfig() *Config {
	return &Config{
		Devices: map[string]*Device{},
	}
}
func NewDevice(MAC string, deviceType DeviceType) Device {
	return Device{
		MAC:  MAC,
		Type: deviceType,
	}
}
