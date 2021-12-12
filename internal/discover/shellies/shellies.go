package shellies

import (
	"context"
	"encoding/json"
	"errors"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
	"hkserver/configs"
	"hkserver/internal/registry"
)

type service struct {
	mqtt   mqtt.Client
	logger *zap.Logger
	init   bool
}

func New(client mqtt.Client, logger *zap.Logger) *service {
	return &service{
		mqtt:   client,
		logger: logger.With(zap.String("service", "shellies")),
		init:   false,
	}
}

func (s *service) Init(ctx context.Context) error {
	s.logger.Debug("Start subscription to shellies/announce topic")
	s.mqtt.Subscribe("shellies/announce", 1, func(client mqtt.Client, message mqtt.Message) {
		var payload Payload
		if err := json.Unmarshal(message.Payload(), &payload); err != nil {
			s.logger.Warn("Found shelly but cannot convert payload", zap.String("payload", string(message.Payload())))
		} else {
			s.logger.Debug("Found shelly", zap.String("ID", payload.ID))
			registry.Add(payload.toDevice())
		}
	})
	s.init = true
	return nil
}

func (s *service) Scan(ctx context.Context) error {
	if !s.init {
		return errors.New("service is not initialized")
	}

	s.mqtt.Publish("shellies/command", 1, false, "announce")
	return nil
}

type Payload struct {
	ID              string `json:"id"`
	Model           string `json:"model"`
	Mac             string `json:"mac"`
	Ip              string `json:"ip"`
	NewFirmware     bool   `json:"new_fw"`
	FirmwareVersion string `json:"fw_ver"`
}

func (p Payload) toDevice() *configs.Device {
	return &configs.Device{
		Type:             "shelly1",
		Model:            p.ID,
		SerialNumber:     p.Mac,
		Ip:               p.Ip,
		FirmwareRevision: p.FirmwareVersion,
	}
}
