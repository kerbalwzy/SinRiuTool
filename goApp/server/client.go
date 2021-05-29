package server

import (
	"bytes"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/kerbalwzy/SinRiuTool/goApp/utils"
	"log"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
	logger  = utils.GetLogger()
)

type Client struct {
	conn     *websocket.Conn
	send     chan []byte
	recv     chan []byte
	shut     chan struct{}
	pingFail int
}

func (obj *Client) readPump() {
	defer func() {
		obj.conn.Close()
		logger.Debug(fmt.Sprintf("%p return readPump", obj))
	}()
	obj.conn.SetReadLimit(maxMessageSize)
	obj.conn.SetReadDeadline(time.Now().Add(pongWait))
	obj.conn.SetPongHandler(func(appData string) error {
		obj.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		select {
		case _, ok := <-obj.shut:
			if !ok {
				return
			}
		default:
			_, message, err := obj.conn.ReadMessage()
			if nil != err {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("Error:%v", err)
				}
				close(obj.shut)
			}
			message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
			logger.Debug(fmt.Sprintf("%p recevie from %s : %s", obj, obj.conn.RemoteAddr(), message))
			obj.recv <- message
		}

	}
}

func (obj *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		obj.conn.Close()
		logger.Debug(fmt.Sprintf("%p return writePump", obj))
	}()

	for {
		select {
		case _, ok := <-obj.shut:
			if !ok {
				return
			}
		case message, ok := <-obj.send:
			obj.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				obj.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := obj.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(obj.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-obj.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			obj.conn.SetWriteDeadline(time.Now().Add(writeWait))
			logger.Debug(fmt.Sprintf("%p ping client %s", obj, obj.conn.RemoteAddr()))
			if err := obj.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				close(obj.shut)
			}
		}
	}
}

// 未来如果需要支持多客户端时可以用得上
func NewClient(conn *websocket.Conn) *Client {
	cli := &Client{
		conn: conn,
		send: make(chan []byte, 256),
		recv: make(chan []byte, 10),
		shut: make(chan struct{}),
	}
	go cli.writePump()
	go cli.readPump()
	return cli
}

var singleClient *Client

// 唯一客户端
func SingleNewClient(conn *websocket.Conn) *Client {
	if singleClient == nil {
		singleClient = NewClient(conn)
	} else {
		oldClient := singleClient
		logger.Debug(fmt.Sprintf("SingleNewClient old %p", oldClient))
		singleClient = NewClient(conn)
		time.AfterFunc(time.Millisecond*500, func() {
			_ = oldClient.conn.Close()
			oldClient = nil
		})
	}
	logger.Debug(fmt.Sprintf("SingleNewClient new %p", singleClient))
	return singleClient
}
