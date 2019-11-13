package simplelb

import (
	"net"

	"github.com/valyala/fasthttp"
)

// NewReverseProxy ...
func NewReverseProxy(isTLS bool, host string) *ReverseProxy {
	return &ReverseProxy{client: &fasthttp.HostClient{Addr: host, IsTLS: isTLS}}
}

// ReverseProxy reverse handler using fasthttp.HostClient
type ReverseProxy struct {
	client *fasthttp.HostClient
}

// ServeHTTP ReverseProxy to serve
// ref to: https://golang.org/src/net/http/httputil/reverseproxy.go#L169
func (p *ReverseProxy) ServeHTTP(ctx *fasthttp.RequestCtx) {
	// prepare request(replace headers and some URL host)
	SetXffHeader(ctx)

	// to save all response header
	res := &ctx.Response
	headerSaver := saveHeaders(res)
	headerHop := MakeHeaderHop()

	req := &ctx.Request
	headerHop.Del(&req.Header)

	// ctx.Logger().Printf("recv a requests to proxy to: %s", p.client.Addr)
	if err := p.client.Do(req, res); err != nil {
		ctx.Logger().Printf("could not proxy: %v\n", err)
		return
	}

	// response to client
	headerHop.Del(&res.Header)
	headerSaver.set(&res.Header)
}

// SetXffHeader sets X-Forwarded-For header for the request.
func SetXffHeader(ctx *fasthttp.RequestCtx) {
	if clientIP, _, err := net.SplitHostPort(ctx.RemoteAddr().String()); err == nil {
		ctx.Request.Header.Add("X-Forwarded-For", clientIP)
	}
}
