package kafka

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/ytf606/golibs/errorx"
	"github.com/ytf606/golibs/logx"
	"github.com/Shopify/sarama"
)

type producer struct {
	asyncProducer   sarama.AsyncProducer
	syncProducer    sarama.SyncProducer
	producterConfig *ProducerConfig
	kafkaConfig     *sarama.Config
	errors          chan *Error

	onceErr sync.Once
}

func NewProducer(ctx context.Context, config *ProducerConfig) (Producer, error) {
	tag := "[kafka_produce_NewProducer]"
	mqConfig := sarama.NewConfig()
	mqConfig.Net.SASL.Enable = false
	mqConfig.Net.SASL.Handshake = true
	mqConfig.Net.TLS.Enable = false
	mqConfig.Metadata.Retry.Max = DefaultRetry
	mqConfig.Producer.Return.Errors = DefaultProducerError
	mqConfig.Producer.RequiredAcks = DefaultAcks
	mqConfig.ChannelBufferSize = DefaultChannelSize

	Producter := &producer{
		errors:          make(chan *Error, DefaultChannelSize),
		producterConfig: config,
	}
	var err error
	if config.Async {
		mqConfig.Producer.Return.Successes = false
		Producter.asyncProducer, err = sarama.NewAsyncProducer(config.Brokers, mqConfig)
	} else {
		mqConfig.Producer.Return.Successes = true
		mqConfig.Producer.Return.Errors = true
		Producter.syncProducer, err = sarama.NewSyncProducer(config.Brokers, mqConfig)
	}
	if err != nil {
		logx.Ex(ctx, tag, "NewProducer failed err:%+v, config:%+v", err, config)
		return nil, errorx.Wrap500Response(err, errorx.KafkaInitProducerErr, "")
	}
	Producter.kafkaConfig = mqConfig
	go Producter.CleanErrorChan(ctx)
	return Producter, nil
}

func (this *producer) Close(ctx context.Context) (err error) {
	tag := "[kafka_produce_Close]"
	if this.producterConfig.Async {
		err = this.asyncProducer.Close()
	} else {
		err = this.syncProducer.Close()
	}
	go func() {
		defer func() {
			err := recover()
			if err != nil {
				logx.Ex(ctx, tag, "close kafka product panic err:%+v", err)
			}
		}()
		time.Sleep(10 * time.Millisecond)
		close(this.errors)
	}()
	return errorx.Wrap500Response(err, errorx.KafkaCloseProducerErr, "")
}

func (this *producer) SyncWriter(ctx context.Context, mess ...*ProducerMessage) (rerr error) {
	tag := "[kafka_produce_SyncWriter]"
	if this.producterConfig.Async {
		return errorx.ErrKafkaProducerConfig
	}
	msgs := this.formatMessage(mess...)
	if len(msgs) == 0 {
		return
	}
	for _, msg := range msgs {
		if msg == nil {
			continue
		}
		partition, offset, err := this.syncProducer.SendMessage(msg)
		if err != nil {
			logx.Ex(ctx, tag, "SyncWriter failed err:%+v", err)
			rerr = errorx.Wrap500Response(err, errorx.KafkaProducerWriterErr, "")
			key, _ := msg.Key.Encode()
			value, _ := msg.Value.Encode()
			this.errors <- &Error{
				Key:       key,
				Value:     value,
				Topic:     msg.Topic,
				Partition: int(partition),
				Offset:    offset,
				Err:       rerr,
				Time:      time.Now(),
			}
			// 连接失败情况
			if strings.Contains(err.Error(), "getsockopt") && strings.Contains(err.Error(), "connection") {
				syncProducer, err := sarama.NewSyncProducer(this.producterConfig.Brokers, this.kafkaConfig)
				logx.Ix(ctx, tag, "retry conn...brokers:%+v, config:%+v", this.producterConfig.Brokers, this.kafkaConfig)
				if err != nil {
					logx.Ex(ctx, tag, "NewKafkaProducter reconn failed err:%+v", err)
				} else {
					this.syncProducer = syncProducer
				}
			}
		}

	}
	return rerr
}

//@todo msg指针
func (this *producer) AsyncWriter(ctx context.Context, mess ...*ProducerMessage) error {
	if !this.producterConfig.Async {
		return errorx.ErrKafkaProducerConfig
	}
	msgs := this.formatMessage(mess...)
	if len(msgs) == 0 {
		return nil
	}
	pchan := this.asyncProducer.Input()
	for _, msg := range msgs {
		if msg == nil {
			continue
		}
		pchan <- msg
	}
	return nil
}

func (this *producer) Errors(ctx context.Context) <-chan *Error {
	tag := "[kafka_produce_Errors]"
	if !this.producterConfig.Async {
		return this.errors
	}
	this.onceErr.Do(func() {
		//after := time.After(10 * time.Millisecond)
		go func() {
			defer func() {
				if err := recover(); err != nil {
					logx.Ix(ctx, tag, "closed errChan panic err:%+v", err)
					return
				}
			}()
			for {
				//perrChan := this.asyncProducer.Errors()
				select {
				case pe := <-this.asyncProducer.Errors():
					if pe == nil {
						continue
					}
					mqerr := &Error{
						Topic:     pe.Msg.Topic,
						Partition: int(pe.Msg.Partition),
						Offset:    pe.Msg.Offset,
						Err:       errorx.Wrap500Response(pe.Err, errorx.KafkaProducerWriterChanErr, ""),
					}
					if pe.Msg != nil && pe.Msg.Key != nil {
						key, _ := pe.Msg.Key.Encode()
						mqerr.Key = key
					}
					if pe.Msg != nil && pe.Msg.Value != nil {
						value, _ := pe.Msg.Value.Encode()
						mqerr.Value = value
					}

					this.errors <- mqerr
					//case <-after:

				}
			}
		}()
	})
	return this.errors
}

func (this *producer) CleanErrorChan(ctx context.Context) {
	tag := "[kafka_produce_CleanErrorChan]"
	for {
		if len(this.Errors(ctx)) > 0 && len(this.Errors(ctx)) > int(DefaultChannelSize/5*4) {
			logx.Ix(ctx, tag, "clean ErrorChan...length:%d", len(this.Errors(ctx)))
			this.errors = make(chan *Error, DefaultChannelSize)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func (this *producer) formatMessage(mess ...*ProducerMessage) []*sarama.ProducerMessage {
	megs := make([]*sarama.ProducerMessage, 0, len(mess))
	for _, m := range mess {
		msg := &sarama.ProducerMessage{
			Topic: this.producterConfig.Topic,
			Value: sarama.ByteEncoder(m.Value),
		}
		if m.Key == nil || len(m.Key) < 1 {
			msg.Key = nil
		} else {
			msg.Key = sarama.ByteEncoder(m.Key)
		}
		megs = append(megs, msg)
	}
	return megs
}
