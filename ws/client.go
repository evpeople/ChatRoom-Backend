package ws

import (
	"bytes"
	"evpeople/ChatRoom/db"
	"evpeople/ChatRoom/middleware"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
	// welcomeMessage = "The User %s is comming\n"
	welcomeMessage = `{"from":"%s","type":"coming"}`
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
	totalId = 0
	idMap   = map[int]int{}
)

type Client struct {
	hub *Hub

	usr *middleware.User
	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}
type Message struct {
	messageId uint64
	userId    uint64
	message   string
	time      time.Time
}

// serveWs handles websocket requests from the peer.
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	first := false
	err = middleware.TokenValid(r)
	if err != nil {
		logrus.Debug(err)
	}
	user, err := middleware.ExtractTokenMetadata(r)
	if err != nil {
		logrus.Debug(err)
	}
	if _, ok := idMap[int(user.ID)]; !ok {
		first = true
		totalId++
		idMap[int(user.ID)] = 1
	} else {
		idMap[int(user.ID)]++
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256), usr: user}
	client.hub.register <- client
	if first {
		welcome := fmt.Sprintf(welcomeMessage, user.Username)
		fmt.Println(welcome)
		hub.broadcast <- []byte(welcome)
	}
	if err != nil {
		return
	}

	go client.writePump()
	go client.readPump()
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		//构造发送的信息
		tmpMessage := Message{userId: c.usr.ID, message: string(message)}
		db.DB.Exec("insert into message(userID,message)"+" VALUES(?,?)", tmpMessage.userId, tmpMessage.message)
		// message = bytes.TrimSpace(bytes.Replace([]byte(c.usr.Username+" 说"+string(message)), newline, space, -1))
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.hub.broadcast <- append(message[:len(message)-1], []byte(`,"from":"`+c.usr.Username+`"}`)...)
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {

		case message, ok := <-c.send:
			// fmt.Println(string(message) + "fsadfasdg")
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			// fmt.Println(message)
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
