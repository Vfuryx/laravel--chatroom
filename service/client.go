package service

import (
	"chatroom/utils"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

// 客户端结构体 存储 客户端的资料
type Client struct {
	Conn  *websocket.Conn // 链接
	Name  string          // 用户名
	PWD   string          // 密码
	Queue chan []byte     // 队列
	Close bool
}

type Clients struct {
	cs map[string]*Client
	sync.RWMutex
}

/**
 * 发送消息给客户端
 */
func (c *Client) Send() {
	for {
		if c.Close {
			return
		}

		select {
		case msg := <-c.Queue:

			err := c.Conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Println(err)
				// 错误立即关闭websocket与客户端的链接

				c.Close = true
				_ = c.Conn.Close()
				break
			}
		}
	}
}

/**
 * 接收 客户端发送的消息
 */
func (c *Client) Recv(handel func(name string, content []byte)) {
	for {
		_, content, err := c.Conn.ReadMessage()

		if err != nil {
			log.Println(err)
			// 错误立即关闭websocket与客户端的链接
			c.Close = true
			_ = c.Conn.Close()

			break
		}
		// 消息处理
		handel(c.Name, content)
		//Broadcast(content)
		//c.Queue <- content
	}
}

func (c Clients) GetClient(name string) (client *Client, ok bool) {
	c.RLock()
	client, ok = c.cs[name]
	c.RUnlock()
	return
}

func (c Clients) SetClient(name string, client *Client) {
	c.Lock()
	c.cs[name] = client
	c.Unlock()
}

func (c Clients) RangeClients(f func(key string, client *Client)) {
	c.Lock()
	for key, value := range c.cs {
		f(key, value)
	}
	c.Unlock()
}

//func (c Clients) SendToClient(name string, msg []byte) {
//	client, ok := c.GetClient(name)
//	if !ok || client.Close != false {
//		return
//	}
//	client.Queue <- msg
//}

// 并发 map
type ConcurrentClientMap struct {
	shards []*Clients
	len    int
}

func NewConcurrentClientMap(num int) *ConcurrentClientMap {
	m := make([]*Clients, num)
	for i := 0; i < num; i++ {
		m[i] = &Clients{cs: make(map[string]*Client)}
	}
	return &ConcurrentClientMap{
		shards: m,
		len:    num,
	}
}

func (ccm ConcurrentClientMap) GetShard(key string) *Clients {
	return ccm.shards[uint(utils.Fnv32(key))%uint(ccm.len)]
}

func (ccm ConcurrentClientMap) Get(key string) (c *Client, ok bool) {
	shard := ccm.GetShard(key)
	c, ok = shard.GetClient(key)
	return
}

func (ccm ConcurrentClientMap) Set(key string, client *Client) {
	shard := ccm.GetShard(key)
	shard.SetClient(key, client)
	return
}

func (ccm ConcurrentClientMap) Range(f func(key string, client *Client)) {
	var wg sync.WaitGroup
	wg.Add(ccm.len)
	for _, shard := range ccm.shards {
		go func(s *Clients) {
			defer wg.Done()
			s.RangeClients(f)
		}(shard)
	}

	wg.Wait()
}

func (ccm ConcurrentClientMap) SendToClient(key string, msg []byte) {
	shard, ok := ccm.Get(key)
	if !ok || shard.Close != false {
		return
	}
	shard.Queue <- msg
}
