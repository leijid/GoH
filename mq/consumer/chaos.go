package consumer

import (
	"GoH/core/rabbitmq"
	"fmt"
)

var exchange string = "go-exchange"
var exchangeType string = "direct"
var routingKey string = "go"

func ReceiveMQMessage() {
	//初始化rabbitmq
	if rabbitmq.ConsumerRabbitMq == nil {
		return
	}
	channleContxt := rabbitmq.ChannelContext{Exchange: exchange, ExchangeType: exchangeType, RoutingKey: routingKey, Reliable: true, Durable: false}
	msgs := rabbitmq.ConsumerRabbitMq.Consumer(&channleContxt)
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			fmt.Println(string(d.Body))
		}
	}()
	<-forever
}
