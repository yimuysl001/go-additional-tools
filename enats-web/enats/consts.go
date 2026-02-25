package enats

import (
	"sync"

	"github.com/gogf/gf/v2/container/gmap"
	nats "github.com/nats-io/nats.go"
)

var (
	// DefaultPort is the default port for the connect server.
	localNatsClient = gmap.NewKVMap[string, *NatsCli](true)
	natsConfigs     = gmap.NewKVMap[string, *NatsConfig](true)
)

const (
	natsDefKey         = "default"
	natsPath           = "nats"
	webPre             = "web."
	consumerKeyPath    = "project.name"
	RespondErrorHeader = "NATS-ERROR"
	RespondStatus      = "NATS-STATUS"
	RespondHeader      = "NATS-HEADER"
	RespondConsumerKey = "NATS-Consumer"
	RespondTraceId     = "NATS-TRACE-ID"
	RespondTraceHeader = "NATS-TRACE-ID"
	RequestUrlPath     = "NATS-URL-PATH"
	RequestMethod      = "NATS-METHOD"
	RequestParam       = "NATS-PARAM"
	RequestMsg         = "NATS-MSG"
)

type HandlerFunc func(ctx *MsgContext) error

// NatsCli NATS客户端结构体
type NatsCli struct {
	ns     *nats.Conn                    // NATS连接对象
	jsName string                        // JetStream名称
	js     nats.JetStreamContext         // JetStream上下文
	err    error                         // 错误信息
	subMap map[string]*nats.Subscription // 订阅映射表
	mu     sync.RWMutex
}
