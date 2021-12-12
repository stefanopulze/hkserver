package shelly

import (
	"fmt"
	"github.com/brutella/hc/accessory"
	pahoMqtt "github.com/eclipse/paho.mqtt.golang"
	"hkserver/configs"
	"hkserver/internal/transport"
)

type shelly1Factory struct {
}

func NewShelly1Factory() *shelly1Factory {
	return &shelly1Factory{}
}

func (f shelly1Factory) Id() string {
	return "Shelly1"
}

func (f shelly1Factory) HandleType() string {
	return "shelly1"
}

func (f shelly1Factory) CreateAccessory(d *configs.Device) (*accessory.Accessory, error) {
	acc := accessory.NewLightbulb(d.Info())

	acc.Lightbulb.On.OnValueRemoteUpdate(func(b bool) {
		command := "off"
		if b {
			command = "on"
		}

		transport.Mqtt.Publish(fmt.Sprintf("shellies/%s/relay/0/command", d.Model), 0, false, command)
	})

	// register
	transport.Mqtt.Subscribe(fmt.Sprintf("shellies/%s/relay/0", d.Model), 0, func(client pahoMqtt.Client, message pahoMqtt.Message) {
		//log.Info("Incoming message from mqtt", zap.String("payload", string(message.Payload())))
		status := string(message.Payload()) == "on"
		acc.Lightbulb.On.SetValue(status)
	})

	return acc.Accessory, nil
}
