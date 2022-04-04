package ws

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
)

type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}
type messageDetail struct {
	To     []string `json:"to"`
	Dialog string   `json:"dialog"`
	Type   string   `json:"type"`
	Detail string   `json:"detail"`
	From   string   `json:"from"`
}
type messageSend struct {
	From   string `json:"from"`
	Dialog string `json:"dialog"`
	Type   string `json:"type"`
	Detail string `json:"detail"`
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				idMap[int(client.usr.ID)]--
				if idMap[int(client.usr.ID)] == 0 {
					delete(idMap, int(client.usr.ID))
				}
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			//从 broadcast获取了消息后，在每一个当前存在的链接上发送消息。
			var messageDetail messageDetail
			err := json.Unmarshal(message, &messageDetail)
			if err != nil {
				logrus.Debug(messageDetail)
				logrus.Debug("wrong message is ", string(message))
				logrus.Debug("json err", err)
			}
			toMap := make(map[string]bool)
			for _, v := range messageDetail.To {
				// logrus.Println("to", k, v)
				toMap[v] = true
			}
			logrus.Println("dialog", messageDetail.Dialog)
			logrus.Println("detail", string(messageDetail.Detail))
			var sendMessage messageSend
			if messageDetail.Type == "users" {
				usrs := ""
				for k, v := range h.clients {
					if v {
						usrs += k.usr.Username + "\n"
					}
				}
				sendMessage = messageSend{Dialog: messageDetail.Dialog, Detail: usrs, From: messageDetail.From, Type: messageDetail.Type}
			} else {
				sendMessage = messageSend{Dialog: messageDetail.Dialog, Detail: messageDetail.Detail, From: messageDetail.From, Type: messageDetail.Type}
			}
			//TODO issue #1 ,此处应该序列化json 数据
			/*{
				"to":{"evpeople","verso"},
				"dialog":"聊天大厅"
				"type":"text",
				"text":"xxxxxxxxxxxxx"
			}*/
			//然后从json中提取出 to 的信息。
			//client 有usr的信息，所以可以与 to 的内容做匹配，从而实现发送到指定的人。
			//如果比较闲，可以给报文增加一个"timeStamp":时间戳的项。
			//另规定：当头部中"to":{}时发送到所有连接到服务器的用户。
			//
			tmp, err := json.Marshal(sendMessage)
			if err != nil {
				logrus.Debug(err)
			}
			for client := range h.clients {
				if len(toMap) != 0 {
					if _, ok := toMap[client.usr.Username]; ok {
						select {
						//这个send的意思是发送给当前存在的连接上。

						case client.send <- tmp:
						default:
							close(client.send)
							delete(h.clients, client)
						}
					}
				} else {
					//TO 为空时
					select {
					//这个send的意思是发送给当前存在的连接上。
					case client.send <- tmp:
					default:
						close(client.send)
						delete(h.clients, client)
					}
				}
			}
		}
	}
}
