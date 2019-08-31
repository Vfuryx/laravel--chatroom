package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

// 实例 websocket 服务端
var ws = func() *websocket.Upgrader {
	return &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		ReadBufferSize:  2048,
		WriteBufferSize: 2048,
	}
}

// 存储客户端的集合 格式是 ["用户名"] -> 客户端结构体
var ClientMap map[string]*Client = make(map[string]*Client, 0)

// 客户端结构体 存储 客户端的资料
type Client struct {
	Conn  *websocket.Conn //链接
	PWD   string          //密码
	Queue chan []byte     //队列
	Close bool
}

// 定义消息和类型常量 要与前端配合
const (
	// 消息源
	MSG_SOURCE_PRIVATE     = 1
	MSG_SOURCE_GROUP       = 2
	MSG_SOURCE_PING_PONG   = 3
	MSG_SOURCE_USER_LIST   = 4
	MSG_SOURCE_ADD_USER    = 5
	MSG_SOURCE_REMOVE_USER = 6
	// 消息类型
	MSG_TYPE_DEFAULT  = 1
	MSG_TYPE_EMORICON = 2
	MSG_TYPE_SHARE    = 3
	MSG_TYPE_REDPACK  = 3
)

// 消息的结构提 定义消息的格式
type MSG struct {
	Source  int       // 源  1: 私信 2: 群聊 3: ping pong 4:用户列表 5 添加用户 6 移除用户
	From    string    // 发送方 用户名
	To      string    // 接收方 源为 1 则是用户名， 为 2 则是群id
	Type    int       // 消息类型 处理方式的不同 1:普通消息; 2:表情; 3:分享;  4:红包;
	Content string    // 消息内容
	SendAt  time.Time // 发送时间
}

/**
 * 发送消息给客户端
 */
func (c *Client) sendproc() {
	for {
		if c.Close {
			return
		}

		select {
		case msg := <-c.Queue:

			err := c.Conn.WriteMessage(websocket.TextMessage, msg);
			if err != nil {
				log.Println(err)
				// 错误立即关闭websocket与客户端的链接

				c.Close = true
				c.Conn.Close()
				break
			}
		}
	}
}

/**
 * 接收 客户端发送的消息
 */
func (c *Client) recvproc() {
	for {
		_, content, err := c.Conn.ReadMessage()

		if err != nil {
			log.Println(err)
			// 错误立即关闭websocket与客户端的链接
			c.Close = true
			c.Conn.Close()

			break
		}
		// 消息处理
		MSGHandle(content)
		//Broadcast(content)
		//c.Queue <- content
	}
}

// 发送消息
// name 用户名
// msg  消息
func SendMsg(name string, msg []byte) {

	for userName, client := range ClientMap {
		if userName == name && client.Close == false {
			client.Queue <- msg
		}
	}
}

// 消息处理函数
// 根绝接收个消息，判断消息源，匹配要对应处理的函数
func MSGHandle(content []byte) {
	var m MSG
	// 解析消息 将 json 转为 MSG结构
	json.Unmarshal(content, &m)

	//判断源
	switch m.Source {
	case MSG_SOURCE_PRIVATE: 		// 私信
		sendToUser(m.To, content)
	case MSG_SOURCE_GROUP:  		// 群聊

	case MSG_SOURCE_PING_PONG:		// 心跳

	case MSG_SOURCE_USER_LIST:		// 获取用户列表
		getUserList(m.From)

	//case MSG_SOURCE_ADD_USER:		// 添加用户
	//	addUser()

	//case MSG_SOURCE_REMOVE_USER:	// 移除用户

	}
}

// 私信
//一对一发送
// to 是指要接收者名称
// content 接收的内容
func sendToUser(to string, content []byte) {
	if ClientMap[to].Close == false {
		ClientMap[to].Queue <- content
	}
}

// 获取用户列表
// client 的 close 属性为 false 即还没关闭，是在线用户
// 整合在线用户信息 发给客户端
func getUserList(from string) {
	content := []gin.H{}

	for name, client := range ClientMap {
		if client.Close == false {

			content = append(content, gin.H{"name": name})
		}
	}

	data, err := json.Marshal(content)
	if err != nil {
		data = []byte{}
	}

	msg := MSG{
		Source:  MSG_SOURCE_USER_LIST,
		From:    from,
		To:      from,
		Type:    MSG_TYPE_DEFAULT,
		Content: string(data),
		SendAt:  time.Now(),
	}
	res, _ := json.Marshal(msg)

	if ClientMap[from].Close == false {
		ClientMap[from].Queue <- res
	}
}

// 有新用户
// 广播给所有在线用户，有新用户登录
// name 新用户名称
func addUser(name string) {

	msg := MSG{
		Source:  MSG_SOURCE_ADD_USER,
		From:    name,
		To:      "",
		Type:    MSG_TYPE_DEFAULT,
		Content: name,
		SendAt:  time.Now(),
	}
	res, _ := json.Marshal(msg)

	for n, client := range ClientMap {
		if n != name && client.Close == false {
			client.Queue <- res
		}

	}
}

// 广播
func Broadcast(msg []byte) {
	for _, client := range ClientMap {
		if client.Close == false {
			client.Queue <- msg
		}
	}
}

// 登录授权检查
func WsAuth(name, pwd string) error {

	c, ok := ClientMap[name];
	if ok {
		if c.PWD != pwd {
			return errors.New("用户已经存,密码不正确")

		}
		// 别的地方登录，强迫下线
		if c.Close == false {
			if err := Offline(name); err != nil {
				return err
			}
		}
	}
	return nil
}

// 离线处理
func Offline(name string) error {

	c, ok := ClientMap[name];
	if ok == false {
		return errors.New(name + " 无法离线")
	}

	c.Conn.WriteControl(
		websocket.PingMessage,
		websocket.FormatCloseMessage(
			websocket.CloseNormalClosure,
			"可以关闭",
		),
		time.Now().Add(time.Second),
	)

	c.Close = true

	c.Conn.Close()

	msg := MSG{
		Source:  MSG_SOURCE_REMOVE_USER,
		From:    name,
		To:      "",
		Type:    MSG_TYPE_DEFAULT,
		Content: name,
		SendAt:  time.Now(),
	}
	res, err := json.Marshal(msg)

	if err != nil {

	}

	Broadcast(res)

	return nil
}

func main() {

	r := gin.Default()

	r.GET("/ws", func(ctx *gin.Context) {

		name := ctx.Query("name")
		pwd := ctx.Query("password")

		// 名称不能为空
		if name == "" || pwd == "" {
			ctx.JSON(http.StatusMisdirectedRequest, gin.H{"code": 10000, "msg": "名称和密码不能为空"})
			return
		}

		// 用户验证
		if err := WsAuth(name, pwd); err != nil {
			ctx.JSON(http.StatusMisdirectedRequest, gin.H{"code": 10000, "msg": err.Error()})
			return
		}

		// 升级为 websocket
		conn, err := ws().Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			panic(err)
			return
		}

		// 存入map
		ClientMap[name] = &Client{
			conn,
			pwd,
			make(chan []byte, 50),
			false,
		}

		// 设置客户端触发关闭事件
		conn.SetCloseHandler(func(code int, text string) error {
			log.Println("我关闭了88")

			if err := Offline(name); err != nil {
				return err
			}

			return nil
		})

		// 开启接收消息协程
		go ClientMap[name].recvproc()
		// 开启发送消息协程
		go ClientMap[name].sendproc()

		fmt.Println(name + " 上线了")

		// 广播所有登录消息
		addUser(name)

		fmt.Println(ClientMap)
	})

	r.LoadHTMLFiles("./public/index.tmpl")

	r.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.tmpl", nil)
	})

	r.Static("/public", "./public")

	r.Run(":9090")
}
