package forward

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
	"go-additional-tools/enats-web/enats"
	"net/url"
	"strings"
	"time"
)

var (
	forwardDefaultClient, _ = client.NewClient()
)

type FeignParams struct {
	Method string
	Path   string
	Name   string
}

func defNats(cloudName string, method string, path string, body any) func(ctx context.Context, r *protocol.Request) error {
	return func(ctx context.Context, r *protocol.Request) error {
		request, err := enats.Nats("http").Request(ctx, webNatsSubPre+cloudName, []byte(path), 5*time.Second)
		if err != nil {
			return err
		}
		var requestUrl = string(request.Data)

		r.SetMethod(method)
		if strings.EqualFold(method, "POST") {
			r.SetBody(gjson.New(body).MustToJson())
		} else {
			for s, a := range gjson.New(body).Map() {
				if strings.Contains(requestUrl, "?") {
					requestUrl += "&" + s + "=" + url.QueryEscape(gconv.String(a))
				} else {
					requestUrl += "?" + s + "=" + url.QueryEscape(gconv.String(a))
				}
			}
		}
		r.SetRequestURI(requestUrl)
		return nil
	}
}

func CloudNatsBytes(ctx context.Context, httpReq *ghttp.Request, cloudName string, method string, path string, body any) ([]byte, error) {
	return CloudForwardBytes(ctx, httpReq, defNats(cloudName, method, path, body))
}

func CloudNatsBody[T any](ctx context.Context, httpReq *ghttp.Request, cloudName string, method string, path string, body any) (T, error) {
	var t T
	err := CloudForwardBody(ctx, httpReq, &t, defNats(cloudName, method, path, body))
	return t, err
}

func CloudForwardBytes(ctx context.Context, httpReq *ghttp.Request, f func(ctx context.Context, r *protocol.Request) error) ([]byte, error) {
	var hertzReq = protocol.AcquireRequest()
	defer protocol.ReleaseRequest(hertzReq)
	// 设置方法
	hertzReq.SetMethod(httpReq.Method)

	// 设置 Host
	hertzReq.SetHost(httpReq.Host)
	// 设置 URI
	if httpReq.RequestURI != "" {
		hertzReq.SetRequestURI(httpReq.RequestURI)
	} else if httpReq.URL != nil {
		hertzReq.SetRequestURI(httpReq.URL.String())
	}
	// 设置请求头（保留多值）
	for key, strings := range httpReq.Header {
		for i, s := range strings {
			if i == 0 {
				hertzReq.Header.Add(key, s)
			} else {
				hertzReq.Header.Set(key, s)
			}

		}
	}

	for key, strings := range httpReq.Trailer {
		for i, s := range strings {
			if i == 0 {
				hertzReq.Header.Add(key, s)
			} else {
				hertzReq.Header.Set(key, s)
			}

		}
	}

	// 处理请求体
	if httpReq.Body != nil {
		hertzReq.SetBodyStream(httpReq.Body, int(httpReq.ContentLength))

		// 可选：手动设置 Content-Length（SetBody 通常会自动处理）
		// hertzReq.Header.SetContentLength(len(bodyBytes))
	} else {
		// 无 Body 时设置空 Body（hertz 默认已处理）
	}
	hertzReq.SetFormDataFromValues(httpReq.Form)
	if httpReq.PostForm != nil {
		hertzReq.SetFormDataFromValues(httpReq.PostForm)
	}
	//hertzReq.SetQueryString(httpReq.URL.RawQuery)

	err := f(ctx, hertzReq)
	if err != nil {
		return nil, err
	}

	response := protocol.AcquireResponse()
	defer protocol.ReleaseResponse(response)

	err = forwardDefaultClient.Do(ctx, hertzReq, response)
	if err != nil {
		return nil, err
	}
	return response.Body(), nil
}

func CloudForwardBody[T any](ctx context.Context, httpReq *ghttp.Request, t *T, f func(ctx context.Context, r *protocol.Request) error) error {
	bytes, err := CloudForwardBytes(ctx, httpReq, f)
	if err != nil {
		return err
	}
	return gjson.New(bytes).Scan(t)
}

func CloudForward(ctx context.Context, httpReq *ghttp.Request, f func(ctx context.Context, r *protocol.Request) error) (*protocol.Response, error) {
	var hertzReq = protocol.AcquireRequest()
	defer protocol.ReleaseRequest(hertzReq)
	// 设置方法
	hertzReq.SetMethod(httpReq.Method)

	// 设置 Host
	hertzReq.SetHost(httpReq.Host)
	// 设置 URI
	if httpReq.RequestURI != "" {
		hertzReq.SetRequestURI(httpReq.RequestURI)
	} else if httpReq.URL != nil {
		hertzReq.SetRequestURI(httpReq.URL.String())
	}
	// 设置请求头（保留多值）
	for key, strings := range httpReq.Header {
		for i, s := range strings {
			if i == 0 {
				hertzReq.Header.Add(key, s)
			} else {
				hertzReq.Header.Set(key, s)
			}

		}
	}

	for key, strings := range httpReq.Trailer {
		for i, s := range strings {
			if i == 0 {
				hertzReq.Header.Add(key, s)
			} else {
				hertzReq.Header.Set(key, s)
			}

		}
	}

	// 处理请求体
	if httpReq.Body != nil {
		hertzReq.SetBodyStream(httpReq.Body, int(httpReq.ContentLength))

		// 可选：手动设置 Content-Length（SetBody 通常会自动处理）
		// hertzReq.Header.SetContentLength(len(bodyBytes))
	} else {
		// 无 Body 时设置空 Body（hertz 默认已处理）
	}
	hertzReq.SetFormDataFromValues(httpReq.Form)
	if httpReq.PostForm != nil {
		hertzReq.SetFormDataFromValues(httpReq.PostForm)
	}
	//hertzReq.SetQueryString(httpReq.URL.RawQuery)

	err := f(ctx, hertzReq)
	if err != nil {
		return nil, err
	}

	response := &protocol.Response{}

	err = forwardDefaultClient.Do(ctx, hertzReq, response)
	return response, err

}
