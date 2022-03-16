package ws

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
			for client := range h.clients {
				select {
				//这个send的意思是发送给当前存在的连接上。
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
