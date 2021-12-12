package mqtt

import (
	"context"
	"errors"
	pahoMqtt "github.com/eclipse/paho.mqtt.golang"
	"hkserver/configs"
	"sync"
)

func New(ctx context.Context, config configs.Mqtt, wg *sync.WaitGroup) (pahoMqtt.Client, error) {
	wg.Add(1)
	brokersCnt := len(config.Brokers)
	if brokersCnt == 0 {
		return nil, errors.New("empty mqtt broker list")
	}

	mqttOptions := pahoMqtt.NewClientOptions()
	mqttOptions.AddBroker(config.Brokers[0])

	if brokersCnt > 1 {
		for i := 1; i < brokersCnt; i++ {
			mqttOptions.AddBroker(config.Brokers[i])
		}
	}

	if len(config.Username) > 0 && len(config.Password) > 0 {
		mqttOptions.SetUsername(config.Username)
		mqttOptions.SetPassword(config.Password)
	}

	client := pahoMqtt.NewClient(mqttOptions)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	go func() {
		<-ctx.Done()
		client.Disconnect(1)
		wg.Done()
	}()

	return client, nil
}
