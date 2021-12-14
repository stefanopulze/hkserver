package shelly

import (
	"fmt"
	"github.com/brutella/hc/accessory"
	pahoMqtt "github.com/eclipse/paho.mqtt.golang"
	"hkserver/configs"
	"hkserver/internal/characteristic"
	"hkserver/internal/transport"
	"strconv"
)

type shelly1PmFactory struct {
}

func NewShelly1PMFactory() *shelly1PmFactory {
	return &shelly1PmFactory{}
}

func (f shelly1PmFactory) Id() string {
	return "Shelly1"
}

func (f shelly1PmFactory) HandleType() string {
	return "shelly1pm"
}

func (f shelly1PmFactory) CreateAccessory(d *configs.Device) (*accessory.Accessory, error) {
	info := d.Info()
	info.Manufacturer = "Shelly"
	acc := accessory.NewLightbulb(info)

	pow := characteristic.NewPower(0)
	pow.Type = characteristic.TypePower
	pow.Unit = "W"
	pow.Description = "Instantaneous power in Watts"

	energy := characteristic.NewPower(0)
	energy.Type = characteristic.TypeTotalPower
	energy.Unit = "Wm"
	energy.Description = "Watt-minute"

	acc.Lightbulb.AddCharacteristic(pow.Characteristic)
	acc.Lightbulb.AddCharacteristic(energy.Characteristic)

	acc.Lightbulb.On.OnValueRemoteUpdate(func(b bool) {
		command := "off"
		if b {
			command = "on"
		}

		transport.Mqtt.Publish(fmt.Sprintf("shellies/%s/relay/0/command", d.Model), 0, false, command)
	})

	// register
	topic := fmt.Sprintf("shellies/%s/relay/0", d.Model)
	transport.Mqtt.Subscribe(topic, 0, func(client pahoMqtt.Client, message pahoMqtt.Message) {
		status := string(message.Payload()) == "on"
		acc.Lightbulb.On.SetValue(status)
	})

	transport.Mqtt.Subscribe(fmt.Sprintf("%s/power", topic), 0, func(client pahoMqtt.Client, message pahoMqtt.Message) {
		p, err := strconv.Atoi(string(message.Payload()))
		if err != nil {
			p = 0
		}

		pow.SetValue(p)
	})

	transport.Mqtt.Subscribe(fmt.Sprintf("%s/energy", topic), 0, func(client pahoMqtt.Client, message pahoMqtt.Message) {
		p, err := strconv.Atoi(string(message.Payload()))
		if err != nil {
			p = 0
		}

		energy.SetValue(p)
	})

	return acc.Accessory, nil
}
