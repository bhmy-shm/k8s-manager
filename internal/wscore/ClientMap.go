package wscore

import (
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"time"
)

var ClientMap *ClientMapStruct

func init() {
	ClientMap = &ClientMapStruct{}
}

type ClientMapStruct struct {
	data sync.Map //  key 是客户端IP  value 就是 WsClient连接对象
}

func (this *ClientMapStruct) Store(conn *websocket.Conn) {
	wsClient := NewWsClient(conn)
	this.data.Store(conn.RemoteAddr().String(), wsClient)
	go wsClient.Ping(time.Second * 1)
	go wsClient.ReadLoop() //处理读 循环
}

//向所有客户端 发送消息--发送deployment列表

func (this *ClientMapStruct) SendAll(v interface{}) {
	this.data.Range(func(key, value interface{}) bool {
		c := value.(*WsClient).conn
		err := c.WriteJSON(v)
		if err != nil {
			this.Remove(c)
			log.Println(err)
		}
		return true
	})
}
func (this *ClientMapStruct) Remove(conn *websocket.Conn) {
	this.data.Delete(conn.RemoteAddr().String())
}
