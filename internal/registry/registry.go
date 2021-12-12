package registry

import (
	"context"
	"github.com/brutella/hc/accessory"
	"go.uber.org/zap"
	"hkserver/configs"
	"hkserver/internal/device"
)

var devices = make(configs.Devices, 0)
var log *zap.Logger

func Initialize(ctx context.Context, logger *zap.Logger, devs configs.Devices) {
	log = logger
	log.Debug("Initialize device registry")

	for _, d := range devs {
		Add(d)
	}
}

func Add(d *configs.Device) {
	var found *configs.Device

	for _, device := range devices {
		if device.Ip == d.Ip {
			found = device
			break
		}
	}

	if found != nil {
		found.Model = d.Model
		found.FirmwareRevision = d.FirmwareRevision
		found.Manufacturer = d.Manufacturer
		found.SerialNumber = d.SerialNumber
		log.Info("Device already present, merge configuration")

		// TODO notify change
	} else {
		log.Debug("Adding device", zap.String("name", d.Name))
		devices = append(devices, d)
	}
}

func GetAllDevices() configs.Devices {
	return devices
}

func GetAllAccessories() []*accessory.Accessory {
	factory := device.NewFactory(log)

	accs := make([]*accessory.Accessory, 0)

	for _, dev := range devices {
		if acc, err := factory.Build(dev); err != nil {
			//if acc, err := device.BuildAccessory(dev); err != nil {
			log.Warn(err.Error())
		} else {
			accs = append(accs, acc)
		}
	}

	return accs
}

func Size() int {
	return len(devices)
}
