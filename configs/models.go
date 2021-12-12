package configs

import "github.com/brutella/hc/accessory"

type Config struct {
	LogLevel  string    `json:"logLevel" mapstruct:"logLevel"`
	Http      Http      `json:"http"`
	Homekit   Homekit   `json:"homekit"`
	Mqtt      Mqtt      `json:"mqtt"`
	Devices   Devices   `json:"devices"`
	Discovery Discovery `json:"discovery"`
}

type Http struct {
	Port int `json:"port,omitempty"`
}

type Mqtt struct {
	Username string   `json:"username"`
	Password string   `json:"password"`
	Brokers  []string `json:"brokers"`
}

type Homekit struct {
	Enable      bool   `json:"enable"`
	Name        string `json:"name"`
	Pin         string `json:"pin"`
	StoragePath string `json:"storagePath"`
	Port        string `json:"port"`
}

type Device struct {
	// Name is the name display in apple home application
	Name             string                 `json:"name"`
	SerialNumber     string                 `json:"serialNumber"`
	Manufacturer     string                 `json:"manufacturer"`
	Model            string                 `json:"model"`
	FirmwareRevision string                 `json:"firmwareRevision"`
	ID               uint64                 `json:"id"`
	Type             string                 `json:"type"`
	Ip               string                 `json:"ip"`
	Options          map[string]interface{} `json:"options"`
}

func (d *Device) Info() accessory.Info {
	return accessory.Info{
		Name:             d.Name,
		SerialNumber:     d.SerialNumber,
		Manufacturer:     d.Manufacturer,
		Model:            d.Model,
		FirmwareRevision: d.FirmwareRevision,
		ID:               d.ID,
	}
}

type Devices = []*Device

type Discovery struct {
	Enable bool `json:"enable"`
}
