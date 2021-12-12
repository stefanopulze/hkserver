package gate

import (
	"fmt"
	"github.com/brutella/hc/accessory"
	"go.uber.org/zap"
	"hkserver/configs"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type gateFactory struct {
	logger *zap.Logger
}

func NewGateFactory(logger *zap.Logger) *gateFactory {
	return &gateFactory{
		logger: logger,
	}
}

func (f gateFactory) Id() string {
	return "Gate"
}

func (f gateFactory) HandleType() string {
	return "gate"
}

func (f gateFactory) CreateAccessory(d *configs.Device) (*accessory.Accessory, error) {
	acc := accessory.NewSwitch(d.Info())

	acc.Switch.On.OnValueRemoteUpdate(func(b bool) {
		command := "off"
		if b {
			command = "on"
		}

		// Send http request
		uri := fmt.Sprintf("http://%s/relay/%d", d.Ip, 0)
		postData := url.Values{}
		postData.Set("turn", command)

		resp, err := http.Post(uri, "application/x-www-form-urlencoded", strings.NewReader(postData.Encode()))
		if err != nil {
			f.logger.Error(fmt.Sprintf("Cannot update accessory status: %s", err), zap.String("accessory", d.Name))
			return
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			f.logger.Error(fmt.Sprintf("Cannot parse accessory response: %s", err), zap.String("accessory", d.Name))
			return
		} else {
			f.logger.Debug(fmt.Sprintf("Accessory response: %s", body), zap.String("accessory", d.Name))
			acc.Switch.On.SetValue(true)
		}

		go func() {
			time.Sleep(3 * time.Millisecond)
			acc.Switch.On.SetValue(false)
		}()
	})

	return acc.Accessory, nil
}
