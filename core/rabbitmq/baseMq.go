package rabbitmq

import (
	"GoH/core/log"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/streadway/amqp"
	"strconv"
	"sync"
	"time"
)

type MqConnection struct {
	Lock       sync.RWMutex
	Connection *amqp.Connection
	MqUri      string
}

type ChannelContext struct {
	Exchange     string
	ExchangeType string
	RoutingKey   string
	Reliable     bool
	Durable      bool
	ChannelId    string
	Channel      *amqp.Channel
}
type BaseMq struct {
	MqConnection *MqConnection
	//rabbitMq通道缓存
	ChannelContexts map[string]*ChannelContext
}

func (bmq *BaseMq) Init() {
	bmq.ChannelContexts = make(map[string]*ChannelContext)
}

// One would typically keep a channel of publishings, a sequence number, and a
// set of unacknowledged sequence numbers and loop until the publishing channel
// is closed.
func (bmq *BaseMq) confirmOne(confirms <-chan amqp.Confirmation) {
	log.Info("waiting for confirmation of one publishing")
	if confirmed := <-confirms; confirmed.Ack {
		log.Info("confirmed delivery with delivery tag: %d", confirmed.DeliveryTag)
	} else {
		log.Error("failed delivery of delivery tag: %d", confirmed.DeliveryTag)
	}
}

/*
 * get md5 from channel context
 */
func (bmq *BaseMq) generateChannelId(channelContext *ChannelContext) string {
	stringTag := channelContext.Exchange + ":" + channelContext.ExchangeType + ":" + channelContext.RoutingKey + ":" +
		strconv.FormatBool(channelContext.Durable) + ":" + strconv.FormatBool(channelContext.Reliable)
	hasher := md5.New()
	hasher.Write([]byte(stringTag))
	return hex.EncodeToString(hasher.Sum(nil))
}

/*
1. use old connection to generate channel
2. update connection then channel
*/
func (bmq *BaseMq) refreshConnectionAndChannel(channelContext *ChannelContext) error {
	bmq.MqConnection.Lock.Lock()
	defer bmq.MqConnection.Lock.Unlock()
	var err error

	if bmq.MqConnection.Connection != nil {
		channelContext.Channel, err = bmq.MqConnection.Connection.Channel()
	} else {
		log.Error("connection not init,dial first time..")
		err = errors.New("connection nil")
	}

	// reconnect connection
	if err != nil {
		for {
			bmq.MqConnection.Connection, err = amqp.Dial(bmq.MqConnection.MqUri)
			if err != nil {
				log.Error("connect mq get connection error,retry..." + bmq.MqConnection.MqUri)
				time.Sleep(10 * time.Second)
			} else {
				log.Info("connection。。。。。..")
				channelContext.Channel, _ = bmq.MqConnection.Connection.Channel()
				break

			}
		}
	}

	if err = channelContext.Channel.ExchangeDeclare(
		channelContext.Exchange,     // name
		channelContext.ExchangeType, // type
		channelContext.Durable,      // durable
		false,                       // auto-deleted
		false,                       // internal
		false,                       // noWait
		nil,                         // arguments
	); err != nil {
		log.Error("channel exchange deflare failed refreshConnectionAndChannel again", err)
		return err
	}

	//add channel to channel cache
	bmq.ChannelContexts[channelContext.ChannelId] = channelContext
	return nil
}

/*
publish message
*/
func (bmq *BaseMq) Publish(channelContext *ChannelContext, body string) error {
	channelContext.ChannelId = bmq.generateChannelId(channelContext)
	if bmq.ChannelContexts[channelContext.ChannelId] == nil {
		bmq.refreshConnectionAndChannel(channelContext)
	} else {
		channelContext = bmq.ChannelContexts[channelContext.ChannelId]
	}
	for {
		if err := channelContext.Channel.Publish(
			channelContext.Exchange,   // publish to an exchange
			channelContext.RoutingKey, // routing to 0 or more queues
			false,                     // mandatory
			false,                     // immediate
			amqp.Publishing{
				Headers:         amqp.Table{},
				ContentType:     "application/json",
				ContentEncoding: "",
				Body:            []byte(body),
				DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
				Priority:        0,              // 0-9
				// a bunch of application/implementation-specific fields
			},
		); err != nil {
			log.Error("send message failed refresh connection")
			time.Sleep(10 * time.Second)
			bmq.refreshConnectionAndChannel(channelContext)
		}
	}
	return nil
}

func (bmq *BaseMq) Consumer(channelContext *ChannelContext) <-chan amqp.Delivery {
	channelContext.ChannelId = bmq.generateChannelId(channelContext)
	if bmq.ChannelContexts[channelContext.ChannelId] == nil {
		bmq.refreshConnectionAndChannel(channelContext)
	} else {
		channelContext = bmq.ChannelContexts[channelContext.ChannelId]
	}
	for {
		msgs, err := channelContext.Channel.Consume(
			channelContext.RoutingKey, // queue
			"",                        // consumer
			true,                      // auto-ack
			false,                     // exclusive
			false,                     // no-local
			false,                     // no-wait
			nil,                       // args
		)
		if err != nil {
			log.Error(err)
			log.Error("Failed to register a consumer")
			bmq.refreshConnectionAndChannel(channelContext)
		}
		return msgs
	}
}
