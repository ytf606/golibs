package kafka

import (
	"context"
	"sync"
	"time"

	"github.com/ytf606/golibs/errorx"
	"github.com/ytf606/golibs/logx"
	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
)

type consumer struct {
	clusterReader  *cluster.Consumer
	consumerConfig *ConsumerConfig
	kafkaConfig    *sarama.Config
	errors         chan *Error
	messages       chan *Message
	onceChan       sync.Once
}

// 暂不支持TLS
func NewConsumer(ctx context.Context, config *ConsumerConfig) Consumer {
	mqConfig := sarama.NewConfig()
	mqConfig.Net.SASL.Enable = false
	mqConfig.Net.SASL.Handshake = true
	mqConfig.Net.TLS.Enable = false
	mqConfig.Metadata.Retry.Max = DefaultRetry
	mqConfig.Version = config.Version

	mqConfig.Consumer.Offsets.Retry.Max = DefaultRetry
	mqConfig.Consumer.Return.Errors = DefaultConsumerError
	mqConfig.Consumer.Offsets.Initial = config.Offset
	// 提交间隔
	mqConfig.Consumer.Offsets.CommitInterval = DefaultCommitInterval
	mqConfig.ChannelBufferSize = DefaultChannelSize
	consumer := &consumer{
		consumerConfig: config,
		kafkaConfig:    mqConfig,
		errors:         make(chan *Error, DefaultChannelSize),
		messages:       make(chan *Message, DefaultChannelSize),
	}
	go consumer.CleanErrorChan(ctx)
	return consumer
}

// 关闭连接
func (this *consumer) Close(ctx context.Context) (err error) {
	tag := "[kafka_consumer_Close]"
	if this.clusterReader != nil {
		err = this.clusterReader.Close()
	}
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logx.Ex(ctx, tag, "close kafka consumer panic err:%+v", err)
			}
		}()
		this.onceChan.Do(func() {
			time.Sleep(10 * time.Millisecond)
			close(this.errors)
			close(this.messages)
		})
	}()
	return errorx.Wrap500Response(err, errorx.KafkaCloseConsumerErr, "")
}

// 获取错误信息
// 默认情况下，来自订阅主题和分区的消息和错误都是多路复用的，并通过使用者的messages()和errors()通道提供。
// sarama.ConsumerModeMultiplex:多路复用
// sarama.ConsumerModePartitions:需要低级访问的用户可以启用ConsumerModePartitions()通道上公开各个分区的consumermodepartition
func (this *consumer) Errors(ctx context.Context) <-chan *Error {

	return this.errors
}

// 启用消费组模式，需要设置groupID，此时设置的partition是无效的
func (this *consumer) GroupReader(ctx context.Context) (<-chan *Message, error) {
	tag := "[kafka_consumer_GroupReader]"
	if err := this.initReader(ctx); err != nil {
		logx.Ex(ctx, tag, "initReader failed err:%+v", err)
		return nil, err
	}
	// consume errors
	go func() {
		for {
			select {
			case err := <-this.clusterReader.Errors():
				if DefaultConsumerError && err != nil {
					this.errors <- &Error{
						Err:  errorx.Wrap500Response(err, errorx.KafkaConsumerClusterReaderErr, ""),
						Time: time.Now(),
					}
				}
			case <-ctx.Done():
				logx.Ix(ctx, tag, "done ctx Errors done------")
				return
			}
		}
	}()
	go func() {
		for {
			select {
			case ntf := <-this.clusterReader.Notifications():
				logx.Ix(ctx, tag, "Rebalanced:%+v", ntf)
			case <-ctx.Done():
				logx.Ix(ctx, tag, "done ctx Notifications done------")
				return
			}
		}
	}()
	go func() {
		// consume messages, watch signals
		for {
			select {
			case msg, ok := <-this.clusterReader.Messages():
				this.tranMessage(ok, msg, errorx.ErrKafkaConsumerGroupReader)
			case <-ctx.Done():
				logx.Ix(ctx, tag, "kafka consumer GroupReader Done...")
				return
				//default:
			}
		}
	}()
	return this.messages, nil
}

// 提交消息,仅仅支持group模式，不支持patition模式
func (this *consumer) CommitMessage(ctx context.Context, mess ...*Message) error {

	if err := this.initReader(ctx); err != nil {
		return err
	}
	if len(mess) > 0 {
		for _, msg := range mess {
			this.clusterReader.MarkOffset(&sarama.ConsumerMessage{
				Key:       msg.Key,
				Value:     msg.Value,
				Topic:     msg.Topic,
				Partition: int32(msg.Partition),
				Offset:    msg.Offset,
				Timestamp: msg.Time,
			}, "")
		}
	}
	return nil
}

// 清除错误信息
func (this *consumer) CleanErrorChan(ctx context.Context) {
	tag := "[kafka_consumer_CleanErrorChan]"
	for {
		if len(this.Errors(ctx)) > 0 && len(this.Errors(ctx)) > int(DefaultChannelSize/5*4) {
			this.errors = make(chan *Error, DefaultChannelSize)
			logx.Ix(ctx, tag, "clean ErrorChan...")
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func (this *consumer) initReader(ctx context.Context) (err error) {
	if this.clusterReader == nil {
		if len(this.consumerConfig.GroupID) == 0 ||
			len(this.consumerConfig.Topic) == 0 ||
			len(this.consumerConfig.Brokers) == 0 {
			return errorx.ErrKafkaConsumerConfig
		}
		tps, tag := []string{this.consumerConfig.Topic}, "[kakfa_consumer_initReader]"
		config := cluster.NewConfig()
		config.Config = *this.kafkaConfig
		config.Consumer.Return.Errors = DefaultConsumerError
		config.Group.Return.Notifications = true
		this.clusterReader, err = cluster.NewConsumer(this.consumerConfig.Brokers, this.consumerConfig.GroupID, tps, config)
		if err != nil {
			logx.Ex(ctx, tag, "kafka consumer group init failed err:%+v", err)
			return errorx.Wrap500Response(err, errorx.KafkaInitConsumerReaderErr, "")
		}
	}
	return nil
}

func (this *consumer) tranMessage(ok bool, msg *sarama.ConsumerMessage, err error) {
	if msg == nil {
		return
	}
	if ok {
		this.messages <- &Message{
			Topic:     msg.Topic,
			Partition: int(msg.Partition),
			Offset:    msg.Offset,
			Key:       msg.Key,
			Value:     msg.Value,
			Time:      msg.BlockTimestamp,
		}
	} else {
		this.errors <- &Error{
			Topic:     msg.Topic,
			Partition: int(msg.Partition),
			Offset:    msg.Offset,
			Key:       msg.Key,
			Value:     msg.Value,
			Time:      msg.BlockTimestamp,
			Err:       err,
		}
	}
}
