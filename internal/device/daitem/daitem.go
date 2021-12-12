package daitem

import (
	"encoding/json"
	"fmt"
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/characteristic"
	"github.com/brutella/hc/service"
	"github.com/stefanopulze/daitem"
	"go.uber.org/zap"
	"hkserver/configs"
)

type daitemFactory struct {
	logger *zap.Logger
}

func NewDaitemFactory(logger *zap.Logger) *daitemFactory {
	return &daitemFactory{
		logger: logger,
	}
}

func (f daitemFactory) Id() string {
	return "Gate"
}

func (f daitemFactory) HandleType() string {
	return "daitem"
}

func (f daitemFactory) CreateAccessory(d *configs.Device) (*accessory.Accessory, error) {
	b, _ := json.Marshal(d.Options)
	var opts Options
	json.Unmarshal(b, &opts)

	options, _ := daitem.DefaultOptions(
		opts.Email,
		opts.Password,
		opts.MasterCode,
	)
	dc := daitem.NewClient(options)

	acc := accessory.New(d.Info(), accessory.TypeSecuritySystem)
	sec := service.NewSecuritySystem()
	acc.AddService(sec.Service)

	go func() {
		state := characteristic.SecuritySystemCurrentStateDisarmed

		if active, err := dc.Status(); err != nil {
			f.logger.Warn("Cannot get current daitem status")
		} else if active {
			state = characteristic.SecuritySystemTargetStateStayArm
		}

		sec.SecuritySystemCurrentState.SetValue(state)
		sec.SecuritySystemTargetState.SetValue(state)
	}()

	sec.SecuritySystemTargetState.OnValueRemoteUpdate(func(i int) {
		f.logger.Debug(fmt.Sprintf("SecuritySystemTargetState: %d", i))

		alarm := i != characteristic.SecuritySystemTargetStateDisarm

		dc.TurnAlarm(alarm)
		sec.SecuritySystemCurrentState.SetValue(i)
	})

	return acc, nil
}
