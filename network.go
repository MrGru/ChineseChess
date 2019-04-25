package main

import (
	"context"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gopherjs/gopherjs/js"
)

func Connect() {
	client := &http.Client{Transport: &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			host, port, err := net.SplitHostPort(addr)
			ws := js.Global.Get("WebSocket").New("ws://" + host + ":5000/dial" + port)
			conn := newWSConn(ws)
			return conn, nil
		},
	}}
	resp, err := client.Get("http://127.0.0.1:5001")

	bs, _ := ioutil.ReadAll(resp.Body)
	js.Global.Get("document").Call("write", string(bs))
}

type wsconn struct {
	ws  *js.Object
	rdr *ChannelReader
}

type ChannelReader struct {
	Message string
	Data    []byte
}

var _ net.Conn = (*wsconn)(nil)

func newWSConn(ws *js.Object) *wsconn {
	ws.Set("binaryType", "arraybuffer")
	out := make(chan []byte, 1)
	ws.Call("addEventListener", "message", func(evt *js.Object) {
		out <- toBytes(evt.Get("data"))
	})
	rdr := NewChannelReader(out)
	return &wsconn{
		ws:  ws,
		rdr: rdr,
	}
}

func (c *wsconn) Read(b []byte) (n int, err error) {
	n, err = c.rdr.Read(b)
	return n, err
}

func (c *wsconn) Write(b []byte) (n int, err error) {
	buf := js.NewArrayBuffer(b)
	c.ws.Call("send", buf)
	return len(b), nil
}

func (c *wsconn) Close() error {
	c.ws.Call("close")
	return nil
}

func (c *wsconn) LocalAddr() net.Addr {
	return websocketAddress{c.ws.Get("url").String()}
}

func (c *wsconn) RemoteAddr() net.Addr {
	return websocketAddress{c.ws.Get("url").String()}
}

func (c *wsconn) SetDeadline(t time.Time) error {
	c.SetReadDeadline(t)
	c.SetWriteDeadline(t)
	return nil
}

func (c *wsconn) SetReadDeadline(t time.Time) error {
	c.rdr.SetDeadline(t)
	return nil
}

func (c *wsconn) SetWriteDeadline(t time.Time) error {
	log.Println("SetWriteDeadline not implemented")
	return nil
}

func toBytes(obj *js.Object) []byte {
	return js.Global.Get("Uint8Array").New(obj).Interface().([]byte)
}

type websocketAddress struct {
	url string
}

func (wsa websocketAddress) Network() string {
	return "ws"
}

func (wsa websocketAddress) String() string {
	return wsa.url
}
