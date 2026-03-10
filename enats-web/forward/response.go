package forward

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/nats-io/nats.go"
	"go-additional-tools/enats-web/enats"
	"net"
	"strings"
)

var (
	defNatsFuncMap = make(map[string]func(subj, queue string, msg *nats.Msg))
	scheme         = "http:"
	ip             string
	addr           string
	host           string
)

func RegisterResponseFunc(name string, f func(subj, queue string, msg *nats.Msg)) {
	defNatsFuncMap[name] = f
}

func defNatsRequest(subj, queue string, msg *nats.Msg) {

	if host == "" {
		ip = g.Cfg().MustGet(gctx.GetInitCtx(), "server.ip", "").String()
		addr = g.Cfg().MustGet(gctx.GetInitCtx(), "server.address", "").String()
		scheme = g.Cfg().MustGet(gctx.GetInitCtx(), "server.scheme", "").String()
		if scheme == "" {
			scheme = "http"
		}
		if ip == "" {
			ip = getLocalIp()
		}
		if addr == "" {
			addr = ":8080"
		}
		addr = strings.ToLower(addr)
		if strings.HasPrefix(addr, "http") {
			host = addr
		} else if strings.HasPrefix(addr, ":") {
			host = scheme + "://" + ip + addr
		} else if strings.Contains(addr, ":") {
			host = scheme + "://" + addr
		} else {
			host = scheme + "://" + ip + ":" + addr
		}

	}

	msg.Respond([]byte(host + "/" + string(msg.Data)))
}

func getLocalIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("获取地址失败:", err)
		return "127.0.0.1"
	}

	for _, addr := range addrs {
		// 检查是否为 IP 地址（不是接口名称）
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			// 只取 IPv4 地址
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}
	return "127.0.0.1"
}

func NatsResponse(name string) {

	natsx := enats.Nats("http")

	if f, ok := defNatsFuncMap[name]; ok {
		natsx.Subscribe(webNatsSubPre+name, webNatsQueue, f)
	} else {
		natsx.Subscribe(webNatsSubPre+name, webNatsQueue, defNatsRequest)
	}
}
