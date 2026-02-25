package eserver

import (
	"errors"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/nats-io/nats.go"
	"go-additional-tools/enats-web/enats"
	"strings"
	"time"
)

// matchRoute 判断 path 是否匹配 pattern，并返回参数
// pattern 格式：
//   - 静态段：直接写，如 "user"
//   - 参数段：以 : 开头，如 ":name"  或 以 { 开头并以 } 结尾，如 "{field}.html"（支持后缀）
//   - 通配符段：以 * 开头，如 "*any"（只能出现在最后，匹配剩余所有路径）
func matchRoute(pattern, path string) (bool, map[string]string) {
	params := make(map[string]string)

	// 分割 pattern 和 path
	patSegs := strings.Split(pattern, "/")
	pathSegs := strings.Split(path, "/")

	// 如果 pattern 以 "/" 开头，那么分割后第一个元素为空，需要去掉
	if len(patSegs) > 0 && patSegs[0] == "" {
		patSegs = patSegs[1:]
	}
	if len(pathSegs) > 0 && pathSegs[0] == "" {
		pathSegs = pathSegs[1:]
	}

	i := 0
	j := 0
	for i < len(patSegs) && j < len(pathSegs) {
		pat := patSegs[i]
		seg := pathSegs[j]

		// 通配符段（*）必须匹配剩余所有路径
		if strings.HasPrefix(pat, "*") {
			// 提取参数名
			paramName := pat[1:] // 去掉 *
			// 将剩余所有 pathSegs 从 j 开始合并成一个字符串（保持斜杠）
			rest := strings.Join(pathSegs[j:], "/")
			params[paramName] = rest
			return true, params
		}

		// 处理参数段（: 或 {}）
		if strings.HasPrefix(pat, ":") {
			// 简单参数，如 :name
			paramName := pat[1:]
			params[paramName] = seg
			i++
			j++
			continue
		}
		if strings.HasPrefix(pat, "{") && strings.Contains(pat, "}") {
			// 处理 {field}.html 这种带后缀的参数
			closeIdx := strings.Index(pat, "}")
			paramName := pat[1:closeIdx]
			suffix := pat[closeIdx+1:] // } 后面的部分，可能是 ".html"
			if suffix != "" {
				// 要求当前段必须以 suffix 结尾，且去除后缀后作为参数值
				if !strings.HasSuffix(seg, suffix) {
					return false, nil
				}
				params[paramName] = strings.TrimSuffix(seg, suffix)
			} else {
				// 没有后缀，如 {subj}
				params[paramName] = seg
			}
			i++
			j++
			continue
		}

		// 静态段必须精确匹配
		if pat != seg {
			return false, nil
		}
		i++
		j++
	}

	// 如果 pattern 和 path 都消耗完毕，则匹配成功
	if i == len(patSegs) && j == len(pathSegs) {
		return true, params
	}

	// 如果 pattern 还剩一个通配符段（*），它可以匹配零个路径段
	if i == len(patSegs)-1 && strings.HasPrefix(patSegs[i], "*") {
		paramName := patSegs[i][1:]
		params[paramName] = "" // 空字符串
		return true, params
	}

	return false, nil
}

func GfRouter(group *ghttp.RouterGroup) {
	cli := enats.Nats("http")

	group.ALL("/{subj}/*act", func(r *ghttp.Request) {
		request, err := cli.WebRequest(r, 60*time.Second)
		if err != nil {
			if errors.Is(err, nats.ErrNoResponders) {
				r.Response.Status = 404
				r.Response.Write("无相应的服务")
			} else if errors.Is(err, nats.ErrTimeout) {
				r.Response.Write("请求超时")
			} else {
				r.Response.Write(err.Error())
			}
			return
		}

		header := enats.GetMsgHeader(request)

		for s, i := range header {
			for i2, s2 := range i {
				if i2 == 0 {
					r.Response.Header().Set(s, s2)
				} else {
					r.Response.Header().Add(s, s2)
				}
			}
		}
		r.Response.Status = gconv.Int(enats.GetStatusMsg(request))
		msg := enats.GetErrorMsg(request)
		if msg == "" {
			r.Response.Write(request.Data)
		} else {
			r.Response.Write(msg)
		}

	})
}
