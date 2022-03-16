package server

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

//用户在线名单列表
var User_list = []string{}

type Data struct {
	Ip       string   `json:"ip"`
	Id       int64    `json:"id"`
	User     string   `json:"user"`
	Content  string   `json:"content"`
	Start    string   `json:"start"`
	Type     string   `json:"type"`
	UserList []string `json:"user_list"`
}

var wu = &websocket.Upgrader{
	ReadBufferSize:  512,
	WriteBufferSize: 512,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//websocket服务
func Myws(w http.ResponseWriter, r *http.Request) {
	//协议升级
	ws, err := wu.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	//创建连接
	c := &Connection{
		Sc:   make(chan []byte, 256),
		Ws:   ws,
		Data: &Data{},
	}

	H.r <- c
	go c.writer()
	c.reader()
	//退出登陆
	defer logout(c)
}

//数据写入器
func (c *Connection) writer() {
	//取出消息并写入
	for message := range c.Sc {
		fmt.Println(message, '\n')
		c.Ws.WriteMessage(websocket.TextMessage, message)
	}
	c.Ws.Close()
}

//数据读取器
func (c *Connection) reader() {
	for {
		//接受ws信息
		_, message, err := c.Ws.ReadMessage()
		if err != nil {
			H.r <- c
			break
		}
		json.Unmarshal(message, &c.Data)
		//解析信息类型
		switch c.Data.Type {
		case "login":
			c.Data.User = c.Data.Content
			c.Data.Ip = c.Data.User
			//在线人数增加
			User_list = append(User_list, c.Data.User)
			c.Data.UserList = User_list
			data_b, _ := json.Marshal(c.Data)
			//发送消息
			H.b <- data_b
		case "user":
			c.Data.Type = "user"
			data_b, _ := json.Marshal(c.Data)
			H.b <- data_b
		case "logout":
			c.Data.Type = "logout"
			User_list = del(User_list, c.Data.User)
			data_b, _ := json.Marshal(c.Data)
			//删除连接
			H.b <- data_b
			H.r <- c
		default:
			fmt.Print("=============default============")
		}
	}
}

//删除登出的用户，维护在线用户名单
func del(slice []string, user string) []string {
	count := len(slice)
	if count == 0 {
		return slice
	}
	if count == 1 && slice[0] == user {
		return []string{}
	}
	var n_slice = []string{}
	for i := range slice {
		if slice[i] == user && i == count {
			return slice[:count]
		} else if slice[i] == user {
			n_slice = append(slice[:i], slice[i+1:]...)
			break
		}
	}
	return n_slice
}

//退出
func logout(c *Connection) {
	c.Data.Type = "logout"
	User_list = del(User_list, c.Data.User)
	c.Data.UserList = User_list
	c.Data.Content = c.Data.User
	data_b, _ := json.Marshal(c.Data)
	H.b <- data_b
	H.r <- c
}

var H = hub{
	c: make(map[*Connection]bool),
	u: make(chan *Connection),
	b: make(chan []byte),
	r: make(chan *Connection),
}

type hub struct {
	//当前在线connection的信息
	c map[*Connection]bool
	//删除connection
	u chan *Connection
	//传递数据
	b chan []byte
	//加入connection
	r chan *Connection
}

func (h *hub) Run() {
	for {
		select {
		//用户连接，添加connection信息
		case c := <-h.r:
			h.c[c] = true
			c.Data.Ip = c.Ws.RemoteAddr().String()
			c.Data.Type = "handshake"
			c.Data.UserList = User_list
			data_b, _ := json.Marshal(c.Data)
			//发送给写入器
			c.Sc <- data_b
		//删除指定
		case c := <-h.u:
			if _, ok := h.c[c]; ok {
				delete(h.c, c)
				close(c.Sc)
			}
		//发送消息
		case data := <-h.b:
			for c := range h.c {
				select {
				//发送数据
				case c.Sc <- data:
				default:
					delete(h.c, c)
					close(c.Sc)
				}
			}
		}
	}
}
