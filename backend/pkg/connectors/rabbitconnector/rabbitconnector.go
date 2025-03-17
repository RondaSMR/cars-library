package rabbitconnector

import (
	"fmt"
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

package rabbitconnector

import (
"context"
"fmt"

amqp "github.com/rabbitmq/amqp091-go"
"go.uber.org/zap"
)

const (
	textContentType = "text/plain"
)

type RabbitConfig struct {
	Host     string `form:"host"`
	Port     string `form:"port"`
	Username string `form:"username"`
	Password string `form:"password"`
	Path     string `form:"path"`
}

type Connector interface {
	GetConnection() *amqp.Connection
	GetChannel() *amqp.Channel
	CloseConnection() error
	Consume(
		queryName string,
		consumerFunc func(amqp.Delivery),
		opts ...QueueOption,
	) error
	Publish(
		ctx context.Context,
		queryName string,
		body []byte,
		queueOpts ...QueueOption,
	) error
}

type RabbitConnector struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

func NewConnector(config *RabbitConfig) (*RabbitConnector, error) {
	if config == nil {
		return nil, fmt.Errorf("config is empty")
	}

	conn, err := amqp.Dial(fmt.Sprintf(`amqp://%v:%v@%v:%v/%v`,
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Path),
	)
	if err != nil {
		return nil, fmt.Errorf("error with connection to rabbitmq: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("error with opening channel: %w", err)
	}

	return &RabbitConnector{
		connection: conn,
		channel:    ch,
	}, nil
}

func (c RabbitConnector) GetConnection() *amqp.Connection {
	return c.connection
}

func (c RabbitConnector) GetChannel() *amqp.Channel {
	return c.channel
}

func (c RabbitConnector) CloseConnection() error {
	return c.connection.Close()
}

type queueConfig struct {
	queueName     string
	consumer      string
	autoAck       bool
	exclusive     bool
	noLocal       bool
	noWait        bool
	durableQueue  bool
	synchronous   bool
	args          amqp.Table
	qos           *qosConfig
	publishConfig *publishConfig
}

// publishConfig describes how to publish a message
type publishConfig struct {
	exchange   string
	routingKey string
	mandatory  bool
	immediate  bool
	publishing amqp.Publishing
}

// qosConfig describes how to receive message
type qosConfig struct {
	prefetchCount int
	prefetchSize  int
	global        bool
}

func defaultConfig(queryName string) queueConfig {
	return queueConfig{
		queueName:    queryName,
		consumer:     "",
		autoAck:      true,
		exclusive:    false,
		noLocal:      false,
		durableQueue: false,
		synchronous:  false,
		noWait:       false,
		args:         nil,
		qos:          nil,
	}
}

// defaultPublishConfig returns a default publish configuration
//
// default content type is text/plain
func defaultPublishConfig(queryName string, body []byte) queueConfig {
	conf := defaultConfig(queryName)
	conf.publishConfig = &publishConfig{
		exchange:   "",
		routingKey: queryName,
		mandatory:  false,
		immediate:  false,
		publishing: amqp.Publishing{
			ContentType: textContentType,
			Body:        body,
		},
	}
	return conf
}

type QueueOption interface {
	apply(queueConfig) queueConfig
}

type queueOptionFunc func(queueConfig) queueConfig

func (fn queueOptionFunc) apply(cfg queueConfig) queueConfig {
	return fn(cfg)
}
func WithConsumer(consumer string) QueueOption {
	return queueOptionFunc(func(cfg queueConfig) queueConfig {
		cfg.consumer = consumer
		return cfg
	})
}

func WithAutoAck(autoAck bool) QueueOption {
	return queueOptionFunc(func(cfg queueConfig) queueConfig {
		cfg.autoAck = autoAck
		return cfg
	})
}

func WithExclusive(exclusive bool) QueueOption {
	return queueOptionFunc(func(cfg queueConfig) queueConfig {
		cfg.exclusive = exclusive
		return cfg
	})
}

func WithNoLocal(noLocal bool) QueueOption {
	return queueOptionFunc(func(cfg queueConfig) queueConfig {
		cfg.noLocal = noLocal
		return cfg
	})
}

func WithNoWait(noWait bool) QueueOption {
	return queueOptionFunc(func(cfg queueConfig) queueConfig {
		cfg.noWait = noWait
		return cfg
	})
}

func WithArgs(args amqp.Table) QueueOption {
	return queueOptionFunc(func(cfg queueConfig) queueConfig {
		cfg.args = args
		return cfg
	})
}

func WithDurableQueue(durableQueue bool) QueueOption {
	return queueOptionFunc(func(cfg queueConfig) queueConfig {
		cfg.durableQueue = durableQueue
		return cfg
	})
}

func WithSynchronous() QueueOption {
	return queueOptionFunc(func(cfg queueConfig) queueConfig {
		cfg.synchronous = true
		return cfg
	})
}

func WithDeliveryMode(deliveryMode uint8) QueueOption {
	return queueOptionFunc(func(cfg queueConfig) queueConfig {
		if cfg.publishConfig == nil {
			panic("isn't configure to publish")
		}
		cfg.publishConfig.publishing.DeliveryMode = deliveryMode
		return cfg
	})
}

func WithContentType(contentType string) QueueOption {
	return queueOptionFunc(func(cfg queueConfig) queueConfig {
		if cfg.publishConfig == nil {
			panic("isn't configure to publish")
		}
		cfg.publishConfig.publishing.ContentType = contentType
		return cfg
	})
}

func WithReplyTo(replyTo string) QueueOption {
	return queueOptionFunc(func(cfg queueConfig) queueConfig {
		if cfg.publishConfig == nil {
			panic("isn't configure to publish")
		}
		cfg.publishConfig.publishing.ReplyTo = replyTo
		return cfg
	})
}

func WithCorrelationId(correlationID string) QueueOption {
	return queueOptionFunc(func(cfg queueConfig) queueConfig {
		if cfg.publishConfig == nil {
			panic("isn't configure to publish")
		}
		cfg.publishConfig.publishing.CorrelationId = correlationID
		return cfg
	})
}

func WithExchange() {
	//	TODO implement
}

func WithQueueBinding() {
	//	TODO implement
}

func WithQos(
	prefetchCount int,
	prefetchSize int,
	global bool,
) QueueOption {
	return queueOptionFunc(func(cfg queueConfig) queueConfig {
		if prefetchCount <= 0 {
			panic(fmt.Sprintf("prefetchCount = %v, must be positive", prefetchCount))
		}

		if prefetchSize < 0 {
			panic(fmt.Sprintf("prefetchSize = %v, must be non-negative", prefetchCount))
		}

		cfg.qos = &qosConfig{
			prefetchCount: prefetchCount,
			prefetchSize:  prefetchSize,
			global:        global,
		}
		return cfg
	})
}

func (c RabbitConnector) Consume(
	queryName string,
	consumerFunc func(amqp.Delivery),
	opts ...QueueOption,
) error {
	config := defaultConfig(queryName)

	for _, opt := range opts {
		config = opt.apply(config)
	}

	q, err := c.GetChannel().QueueDeclare(
		config.queueName,
		config.durableQueue,
		config.exclusive,
		config.autoAck,
		config.noWait,
		config.args,
	)
	if err != nil {
		return fmt.Errorf("error with queue declare: %w", err)
	}

	if config.qos != nil {
		if err = c.channel.Qos(
			config.qos.prefetchCount,
			config.qos.prefetchSize,
			config.qos.global,
		); err != nil {
			return fmt.Errorf("error with qos: %w", err)
		}
	}

	delivery, err := c.GetChannel().Consume(
		q.Name,
		config.consumer,
		config.autoAck,
		config.exclusive,
		config.noLocal,
		config.noWait,
		config.args,
	)
	if err != nil {
		return fmt.Errorf("error with consume: %w", err)
	}

	go func() {
		for d := range delivery {
			if !config.synchronous {
				go consumerFunc(d)
			} else {
				consumerFunc(d)
			}
		}
	}()

	return nil
}

func (c RabbitConnector) Publish(
	ctx context.Context,
	queryName string,
	body []byte,
	queueOpts ...QueueOption,
) error {
	zap.L().Debug("publish")
	config := defaultPublishConfig(queryName, body)

	for _, opt := range queueOpts {
		config = opt.apply(config)
	}

	q, err := c.GetChannel().QueueDeclare(
		config.queueName,
		config.durableQueue,
		false,
		false,
		config.noWait,
		config.args,
	)
	if err != nil {
		return fmt.Errorf("error with queue declare: %w", err)
	}

	zap.L().Info(fmt.Sprintf("publishing to queue %v", q.Name))
	zap.L().Info(fmt.Sprintf("publishing to queue %v", config.queueName))

	return c.channel.PublishWithContext(
		ctx,
		"",
		q.Name,
		config.publishConfig.mandatory,
		config.publishConfig.immediate,
		config.publishConfig.publishing,
	)
}

