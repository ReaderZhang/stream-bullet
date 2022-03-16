package comet

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

type Connection struct {
	Ws   *websocket.Conn
	Sc   chan []byte
	Data *Data
}

type Data struct {
	Ip       string   `json:"ip"`
	Id       int64    `json:"id"`
	User     string   `json:"user"`
	Content  string   `json:"content"`
	UserList []string `json:"user_list"`
}

var wu = &websocket.Upgrader{
	ReadBufferSize:  512,
	WriteBufferSize: 512,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Handler(w http.ResponseWriter, r *http.Request, header http.Header) {
	ws, err := wu.Upgrade(w, r, header)
	if err != nil {
		return
	}

	connection := &Connection{
		Sc:   make(chan []byte, 256),
		Ws:   ws,
		Data: &Data{},
	}

	MessageHub.addChannel <- connection

	go connection.writer()
	connection.reader()

	defer func() {
		fmt.Println("server is closing ...")
	}()
}

func (c *Connection) writer() {
	//取出消息并写入
	for message := range c.Sc {
		fmt.Println(message)
		c.Ws.WriteMessage(websocket.TextMessage, message)
	}
	c.Ws.Close()
}

func (c *Connection) reader() {
	for {
		//接收消息
		_, message, err := c.Ws.ReadMessage()
		if err != nil {
			MessageHub.addChannel <- c
			break
		}
		json.Unmarshal(message, &c.Data)
		//TODO 写入kafka存储
		MessageHub.data <- message
	}
}
