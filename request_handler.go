package middleman

import (
	"net"
	"bytes"
	"fmt"
	"net/url"
	"strings"
	"io"
	"log"
	"strconv"
)

type RequestHandler func(c net.Conn) ()

var RequestHandleManager = make(map[string]RequestHandler)

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

func socks5Handler(client net.Conn) {
	if client == nil {
		return
	}
	defer client.Close()
	var b [1024]byte
	n, err := client.Read(b[:])
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("%d bytes from client", n)

	if b[0] == 0x05 {
		client.Write([]byte{0x05, 0x00})
		n, err := client.Read(b[:])
		var host, port string
		fmt.Print(string(b[:]))
		switch b[3] {
		case 0x01:
			host = net.IPv4(b[4], b[5], b[6], b[7]).String()
		case 0x03:
			host = string(b[5 : n-2])
		case 0x04:
			host = net.IP{b[4], b[5], b[6], b[7], b[8], b[9], b[10], b[11], b[12], b[13], b[14], b[15], b[16], b[17], b[18], b[19]}.String()
		}
		port = strconv.Itoa(int(b[n-2])<<8 | int(b[n-1]))
		server, err := net.Dial("tcp", net.JoinHostPort(host, port))
		if err != nil {
			log.Println(err)
			return
		}
		defer server.Close()
		client.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

		go io.Copy(server, client)
		io.Copy(client, server)
	} else {
		client.Write([]byte("this is socks5 proxy"))
	}
}

func init() {
	RequestHandleManager["http"] = httpHandler
	RequestHandleManager["socks5"] = socks5Handler
}
