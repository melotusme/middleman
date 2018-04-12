package middleman

import (
	"net"
)

type RequestHandler func(c net.Conn) ()

var RequestHandleManager = make(map[string]RequestHandler)

func init() {
	RequestHandleManager["http"] = httpHandler
	RequestHandleManager["socks5"] = socks5Handler
}
