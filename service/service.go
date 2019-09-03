package service

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	json "github.com/json-iterator/go"
	"log"
	"time"
)

// 存储客户端的集合 格式是 ["用户名"] -> 客户端结构体
var ClientMap map[string]*Client = make(map[string]*Client)

// 客户端结构体 存储 客户端的资料
type Client struct {
	Conn  *websocket.Conn // 链接
	Name  string		  // 用户名
	PWD   string          // 密码
	Queue chan []byte     // 队列
	Close bool
}

var GroupMap map[string]*Group = make(map[string]*Group)

// 群组结构体
type Group struct {
	Name string				// 名称
	Type int				// 公共或者私有 0 / 1
	Admin string			// 群主
	List map[string]*Client	// 私有群的用户列表
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
	MSG_SOURCE_CREATE_PUBLIC_GROUP = 7
	MSG_SOURCE_CREATE_PRIVATE_GROUP = 8
	MSG_SOURCE_GROUP_LIST = 9

	// 消息类型
	MSG_TYPE_DEFAULT  = 1
	MSG_TYPE_EMORICON = 2
	MSG_TYPE_SHARE    = 3
	MSG_TYPE_REDPACK  = 4

	// 群组类型
	GROUP_TYPE_PUBLIC = 0
	GROUP_TYPE_PRIVATE = 1
)

// 消息的结构提 定义消息的格式
type MSG struct {
	Source  int       // 源  1: 私信 2: 群聊 3: ping pong 4:用户列表 5:添加用户 6:移除用户 7:创建群组 8:私有群组 9: 群组列表
	From    string    // 发送方 用户名
	To      string    // 接收方 源为 1 则是用户名， 为 2 则是群id
	Type    int       // 消息类型 处理方式的不同 1:普通消息; 2:表情; 3:分享;  4:红包;
	Content string    // 消息内容
	SendAt  time.Time // 发送时间
}

/**
 * 发送消息给客户端
 */
func (c *Client) Sendproc() {
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
func (c *Client) Recvproc() {
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
		MSGHandle(c.Name, content)
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
// fromName 信息接收源
// 根绝接收个消息，判断消息源，匹配要对应处理的函数
func MSGHandle(fromName string, content []byte) {
	var m MSG
	// 解析消息 将 json 转为 MSG结构
	json.Unmarshal(content, &m)

	//判断源
	switch m.Source {
	case MSG_SOURCE_PRIVATE: 		// 私信
		sendToUser(m.To, content)
	case MSG_SOURCE_GROUP:  		// 群聊
		sendGroup(fromName, m)

	case MSG_SOURCE_PING_PONG:		// 心跳

	case MSG_SOURCE_USER_LIST:		// 获取用户列表
		getUserList(fromName)

	case MSG_SOURCE_CREATE_PUBLIC_GROUP:
		createPublicGroup(fromName, m.Content)

	case MSG_SOURCE_GROUP_LIST:
		getGroupList(fromName)

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

// 群聊
func sendGroup(from string, m MSG)  {
	g, ok := GroupMap[m.To]
	if !ok {
		return
	}

	msg := MSG{
		Source: MSG_SOURCE_GROUP,
		From: from,
		To: g.Name,
		Type: m.Type,
		Content: m.Content,
		SendAt: m.SendAt,
	}

	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}


	switch g.Type {
	case GROUP_TYPE_PUBLIC:
		Broadcast(b)
	case GROUP_TYPE_PRIVATE:

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

// 获取群组
func getGroupList(from string)  {
	content := []gin.H{}
	for name, group := range GroupMap {
		if group.Type == GROUP_TYPE_PUBLIC {
			content = append(content, gin.H{"name": name})
		}
	}

	data, err := json.Marshal(content)
	if err != nil {
		data = []byte{}
	}

	msg := MSG{
		Source:  MSG_SOURCE_GROUP_LIST,
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

// 创建公共群组
func createPublicGroup(from string, content string)  {
	fmt.Println("group")
	// 判断是否存在公共群组(存在则忽略)
	if _, ok := GroupMap[content]; ok {
		return
	}

	// 创建公共群组
	GroupMap[content] = &Group{
		Name: content,
		Type: GROUP_TYPE_PUBLIC,
		Admin: from,
		List: map[string]*Client{},
	}

	msg := MSG{
		Source:	MSG_SOURCE_CREATE_PUBLIC_GROUP,
		From: from,
		To: "",
		Type: MSG_TYPE_DEFAULT,
		Content: content,
		SendAt: time.Now(),
	}

	m, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	Broadcast(m)
}

// 有新用户
// 广播给所有在线用户，有新用户登录
// name 新用户名称
func AddUser(name string) {

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

// 离线处理
func Offline(name string) error {

	c, ok := ClientMap[name];
	if ok == false {
		return errors.New(name + " 无法离线")
	}

	err := c.Conn.WriteControl(
		websocket.PingMessage,
		websocket.FormatCloseMessage(
			websocket.CloseNormalClosure,
			"可以关闭",
		),
		time.Now().Add(time.Second),
	)

	if err != nil {

	}

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