package device

import (
	"fmt"
	"github.com/brutella/hc/accessory"
	"go.uber.org/zap"
	"hkserver/configs"
	"hkserver/internal/device/daitem"
	"hkserver/internal/device/gate"
	"hkserver/internal/device/shelly"
)

// AccessoryFactory is the base interface to build a hc accessory
type AccessoryFactory interface {
	Id() string
	HandleType() string
	CreateAccessory(*configs.Device) (*accessory.Accessory, error)
}

type factory struct {
	factories map[string]AccessoryFactory
	logger    *zap.Logger
}

func NewFactory(logger *zap.Logger) *factory {
	f := factory{
		factories: make(map[string]AccessoryFactory),
		logger:    logger.With(zap.String("component", "factory")),
	}

	f.Register(shelly.NewShelly1Factory())
	f.Register(shelly.NewShelly1PMFactory())
	f.Register(gate.NewGateFactory(logger))
	f.Register(daitem.NewDaitemFactory(logger))

	return &f
}

func (f *factory) Register(af AccessoryFactory) error {
	t := af.HandleType()

	if _, ok := f.factories[t]; ok {
		e := fmt.Sprintf("Accessory factory is already registered for type: %s", t)
		f.logger.Warn(e)
		return fmt.Errorf(e)
	}

	f.factories[t] = af
	f.logger.Debug(fmt.Sprintf("Accessory factory register for type: %s", t))
	return nil
}

func (f *factory) Build(d *configs.Device) (*accessory.Accessory, error) {
	fact, ok := f.factories[d.Type]
	if !ok {
		return nil, fmt.Errorf("unknow accessory type %s, cannot build", d.Type)
	}

	acc, err := fact.CreateAccessory(d)
	if err != nil {
		return nil, fmt.Errorf("error creating accessory %s: %s", d.Name, err)
	}

	return acc, nil
}
