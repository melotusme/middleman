package middleman

import (
	"net"
	"log"
	"bytes"
	"fmt"
	"net/url"
	"strings"
	"io"
)

func httpHandler(c net.Conn) {
	if c == nil {
		return
	}
	defer c.Close()

	var b [1024]byte
	n, err := c.Read(b[:])
	if err != nil {
		log.Println(err)
		return
	}
	var method, host, address string
	requestLIine := string(b[:bytes.IndexByte(b[:], '\n')])
	log.Println(requestLIine)
	fmt.Sscanf(requestLIine, "%s%s", &method, &host)
	hostPortURL, err := url.Parse(host)
	if err != nil {
		log.Println(err)
		return
	}
	if hostPortURL.Port() == "443" {
		address = hostPortURL.Scheme + ":443"
	} else {
		if strings.Index(hostPortURL.Host, ":") == -1 {
			address = hostPortURL.Host + ":80"
		} else {
			address = hostPortURL.Host
		}
	}

	server, err := net.Dial("tcp", address)
	if err != nil {
		log.Println(err)
		return
	}
	if method == "CONNECT" {
		fmt.Fprint(c, "HTTP/1.1 200 Connection established\r\n\r\n")
	} else {
		server.Write(b[:n])
	}
	go io.Copy(server, c)
	io.Copy(c, server)
}
