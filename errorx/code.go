package errorx

/**
 * 错误码采用AABBCC形式
 * 其中：
 * A：服务或业务线编码，30开始，每个服务或业务线唯一
 *  10 - 29：预留，用作未来系统全局服务模块，如全局库包、注册中心、认证中心，服务网关、服务监控等
 *  30：辅导中心
 * B：错误类型编码，分为程序类、通用业务类和实体业务类，每种类型可能拥有多个子类型
 *  00 - 29：程序类错误类型编码
 *  00: 全局暂未分类异常，如请求、响应、编码等
 *  01：关系数据库，如mysql、tidb等
 *  02：非关系数据库，如mongdb、hbase
 *  03：缓存类数据库，如redis、memcache
 *  04：消息队列类，如kafka、rabbitmq
 *  05：搜索引擎类，如elasticsearch、solr
 *  06：网关请求框架或中间件，如curl、http request、rpcx、grpc、jwt、oath
 *  30 - 39：通用业务类错误类型编码
 *  30：路由分发，如请求方法不存在、路由不存在
 *  31：签名认证，如登陆、认证、签名校验等
 *  40 - 99：实体业务类错误类型编码
 *  40：老师
 *  41：学生
 *  42：班级
 * C：错误顺序号，可自行定义，或利用系统自动增加
 */

//服务编码或业务线编码
const ServiceCode = 10

//模块编码 - 程序类
const (
	GlobalUnknownException int = iota
	RelationDb
	NoRelationDb
	CacheDb
	MessageQueue
	SearchEngine
	GatewayRegister
)

/*************************详细错误码************************************************/
//全局暂未分类异常
const (
	GinxResponseTypeErr int = iota + ServiceCode*10000 + GlobalUnknownException*100
	methodNotFoundErr
	routerNotFoundErr
	JsonParseCodeErr
	Base64ParseCodeErr
	PemParseCodeErr
)

//Gateway类错误码列表
const (
	RpcEndpointErr int = iota + ServiceCode*10000 + GatewayRegister*100
	RpcSelectClientErr
	RpcReturnErr
	RpcCodeErr
	HttpRequestReturnErr
	JwtCreateSignStringErr
	JwtTokenMalformedErr
	JwtTokenExpiredErr
	JwtTokenNotValidYetErr
	JwtTokenInvalidErr
	JwtSignMethodErr
	AppleTokenInvalidErr
)

//MQ类错误码列表
const (
	KafkaInitProducerErr int = iota + ServiceCode*10000 + MessageQueue*100
	KafkaCloseProducerErr
	KafkaProducerConfigErr
	KafkaProducerWriterErr
	KafkaProducerWriterChanErr
	KafkaCloseConsumerErr
	KafkaInitConsumerReaderErr
	KafkaConsumerConfigErr
	KafkaConsumerMessageErr
	KafkaConsumerClusterReaderErr
)

//缓存类错误码列表
const (
	RedisInitConfigErr int = iota + ServiceCode*10000 + CacheDb*100
	RedisConnectErr
)
