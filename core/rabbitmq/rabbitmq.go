package rabbitmq

import (
	"fmt"
	"github.com/dlintw/goconf"
	"sync"
)

var ProducerRabbitMq *BaseMq
var producerRabbitMqonce sync.Once
var ConsumerRabbitMq *BaseMq
var consumerRabbitMqonce sync.Once

func getProducerRabbitMq(uri string) *BaseMq {
	producerRabbitMqonce.Do(func() {
		ProducerRabbitMq = &BaseMq{
			MqConnection: &MqConnection{MqUri: uri},
		}
		ProducerRabbitMq.Init()
	})
	return ProducerRabbitMq
}

func getConsumerRabbitMq(uri string) *BaseMq {
	consumerRabbitMqonce.Do(func() {
		ConsumerRabbitMq = &BaseMq{
			MqConnection: &MqConnection{MqUri: uri},
		}
		ConsumerRabbitMq.Init()
	})
	return ConsumerRabbitMq
}

func InitProducerMqConnection(conf *goconf.ConfigFile) *BaseMq {
	uri, _ := conf.GetString("rabbitmq", "rabbitmq_uri")
	if uri == "" {
		fmt.Println("未启用rabbimq")
		return nil
	}
	return getProducerRabbitMq(uri)
}

func InitConsumerMqConnection(conf *goconf.ConfigFile) *BaseMq {
	uri, _ := conf.GetString("rabbitmq", "rabbitmq_uri")
	if uri == "" {
		fmt.Println("未启用rabbimq")
		return nil
	}
	return getConsumerRabbitMq(uri)
}
