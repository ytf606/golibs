package kafka

import (
	"time"

	"github.com/Shopify/sarama"
)

const (
	MessageTypeEarliest   = -1 //从第一个开始
	MessageTypeLatest     = -2 //从最新的开始
	DefaultRetry          = 3
	DefaultAcks           = sarama.WaitForAll
	DefaultProducerError  = false
	DefaultConsumerError  = true
	DefaultChannelSize    = 256
	DefaultCommitInterval = 1 * time.Second
)

type Error struct {
	Topic     string
	Partition int
	Key       []byte
	Value     []byte
	Err       error
	Offset    int64

	Time time.Time
}

// 返回信息
type Message struct {
	Topic     string
	Partition int
	Offset    int64
	Key       []byte
	Value     []byte

	// If not set at the creation, Time will be automatically set when
	// writing the message.
	Time time.Time
}

type ProducerMessage struct {
	Key   []byte
	Value []byte
}

type ProducerConfig struct {
	// kafka代理地址
	Brokers []string
	// 主题
	Topic string

	// 是否异步
	Async bool
}

type ConsumerConfig struct {
	// kafka代理地址
	Brokers []string
	// 主题
	Topic string
	// consumer group id
	GroupID string
	// kafka 版本
	Version sarama.KafkaVersion

	// -1:从第一个开始，-2:从最新一个开始
	Offset int64
}
