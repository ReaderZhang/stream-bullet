package comet

import "encoding/json"

var MessageHub = hub{
	connections: make(map[*Connection]bool),
	delChannel:  make(chan *Connection),
	data:        make(chan []byte),
	addChannel:  make(chan *Connection),
}

var UserList = []string{}

type hub struct {
	//当前在线连接的信息
	connections map[*Connection]bool
	//删除连接
	delChannel chan *Connection
	//加入connection
	addChannel chan *Connection
	//数据传递
	data chan []byte
}

func (h *hub) Run() {
	for {
		select {
		//用户连接
		case c := <-MessageHub.addChannel:
			MessageHub.connections[c] = true
			c.Data.Ip = c.Ws.RemoteAddr().String()
			c.Data.UserList = UserList
			data, _ := json.Marshal(c.Data)
			c.Sc <- data
		case c := <-MessageHub.delChannel:
			if _, ok := MessageHub.connections[c]; ok {
				//TODO 删除用户
				close(c.Sc)
			}
		case data := <-MessageHub.data:
			for c := range MessageHub.connections {
				select {
				//发送消息
				case c.Sc <- data:
				default:
					//TODO 删除connection的信息
					close(c.Sc)
				}
			}
		}

	}
}
