package producer

import (
	"GoH/core/rabbitmq"
	"fmt"
)

var exchange string = "go-exchange"
var exchangeType string = "direct"
var routingKey string = "go"

func SendMQMessage(message string) {
	if rabbitmq.ProducerRabbitMq == nil {
		return
	}
	channleContxt := rabbitmq.ChannelContext{Exchange: exchange, ExchangeType: exchangeType, RoutingKey: routingKey, Reliable: true, Durable: false}
	for {
		fmt.Println("sending message")
		rabbitmq.ProducerRabbitMq.Publish(&channleContxt, message)
	}
}
