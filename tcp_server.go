package wstunnel

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/url"

	log "github.com/fangdingjun/go-log/v5"
)

type tcpServer struct {
	addr    string
	remote  string
	payload string
}

var l net.Listener

func StopWSTunnel() {
	l.Close()
}

func (srv *tcpServer) run() {
	var err error
	l, err = net.Listen("tcp", srv.addr)
	if err != nil {
		log.Errorln(err)
		return
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Error(err)
			return
		}
		go srv.serve(conn, srv.payload)
	}
}

func (srv *tcpServer) serve(c net.Conn, payload string) {
	defer c.Close()

	u, _ := url.Parse(srv.remote)

	log.Debugf("connected from %s, forward to %s", c.RemoteAddr(), srv.remote)

	defer func() {
		log.Debugf("from %s, finished", c.RemoteAddr())
	}()

	if u.Scheme == "ws" || u.Scheme == "wss" {
		dialer.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		header := make(http.Header)
		header.Add("Host", payload)
		header.Add("X-Online-Host", payload)
		header.Add("X-Forward-Host", payload)
		header.Add("Referer", "https://"+payload)
		header.Add("Origin", "https://"+payload)
		header.Add("User-Agent", "Mozilla/5.0 (Linux; Android 8.0.0; SM-G960F Build/R16NW) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.84 Mobile Safari/537.36")
		conn1, resp, err := dialer.Dial(srv.remote, header)
		if err != nil {
			log.Errorln(err)
			return
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusSwitchingProtocols {
			log.Errorf("dial remote ws %d", resp.StatusCode)
			return
		}
		defer conn1.Close()

		forwardWS2TCP(conn1, c)
		return
	}

	if u.Scheme == "tcp" {
		conn1, err := net.Dial("tcp", u.Host)
		if err != nil {
			log.Errorln(err)
			return
		}
		defer conn1.Close()

		forwardTCP2TCP(c, conn1)
		return
	}

	log.Errorf("unsupported scheme %s", u.Scheme)
}
