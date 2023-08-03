package kafka

import "context"

type Producer interface {
	// 关闭连接
	Close(ctx context.Context) error
	// 同步写
	SyncWriter(ctx context.Context, mes ...*ProducerMessage) error
	// 异步写
	AsyncWriter(ctx context.Context, mes ...*ProducerMessage) error
	// 获取错误信息
	Errors(ctx context.Context) <-chan *Error
}

type Consumer interface {
	// 关闭连接
	Close(ctx context.Context) error
	// 获取错误信息
	Errors(ctx context.Context) <-chan *Error
	// 消费组
	GroupReader(ctx context.Context) (<-chan *Message, error)
	// 提交消息
	CommitMessage(ctx context.Context, mess ...*Message) error
}
