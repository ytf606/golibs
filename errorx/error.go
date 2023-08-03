package errorx

import (
	"github.com/pkg/errors"
)

var (
	New          = errors.New
	Errorf       = errors.Errorf
	Wrap         = errors.Wrap
	Wrapf        = errors.Wrapf
	WithStack    = errors.WithStack
	WithMessage  = errors.WithMessage
	WithMessagef = errors.WithMessagef
)

var (
	ErrMethodNotAllow = New500Response(methodNotFoundErr, "method not allow")
	ErrNotFound       = New500Response(routerNotFoundErr, "router not found")

	ErrJwtTokenMalformed   = New401Response(JwtTokenMalformedErr, "That's not even a token of jwt")
	ErrJwtTokenExpired     = New401Response(JwtTokenExpiredErr, "Token of jwt is expired")
	ErrJwtTokenNotValidYet = New401Response(JwtTokenNotValidYetErr, "Token of jwt not active yet")
	ErrJwtTokenInvalid     = New401Response(JwtTokenInvalidErr, "Token of jwt invalid")
	ErrJwtSignMethod       = New401Response(JwtSignMethodErr, "Sign method of jwt invalid")

	ErrKafkaProducerConfig      = New500Response(KafkaProducerConfigErr, "Async error for kafka writer mode")
	ErrKafkaConsumerGroupReader = New500Response(KafkaConsumerMessageErr, "kafka consumer message read chan error")
	ErrKafkaConsumerConfig      = New500Response(KafkaConsumerConfigErr, "CommitMessage cannot run in consumer group")

	ErrRedisInitConfig = New500Response(RedisInitConfigErr, "")
	ErrRedisConnect    = New500Response(RedisConnectErr, "")
)
