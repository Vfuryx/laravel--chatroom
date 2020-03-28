package main

import (
	"chatroom/service"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
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

// 登录授权检查
func WsAuth(name, pwd string) error {
	c, ok := service.S.CCM.Get(name)
	if ok {
		if c.PWD != pwd {
			return errors.New("用户已经存,密码不正确")

		}
		// 别的地方登录，强迫下线
		if c.Close == false {
			if err := service.Offline(name); err != nil {
				return err
			}
		}
	}
	return nil
}

func upgradeWebSocket(ctx *gin.Context) {

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
	c := &service.Client{
		Conn:  conn,
		Name:  name,
		PWD:   pwd,
		Queue: make(chan []byte, 50),
		Close: false,
	}
	service.S.CCM.Set(name, c)

	// 设置客户端触发关闭事件
	conn.SetCloseHandler(service.CloseHandle(name))

	//c, _ = service.S.GetClient(name)
	// 开启接收消息协程
	go c.Recv(service.MSGHandle)
	// 开启发送消息协程
	go c.Send()

	fmt.Println(name + " 上线了")

	// 广播所有登录消息
	service.AddUser(name)
}

func main() {

	r := gin.Default()

	r.GET("/ws", upgradeWebSocket)

	r.LoadHTMLFiles("./public/index.tmpl")

	r.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.tmpl", nil)
	})

	r.Static("/public", "./public")

	if err := r.Run(":9091"); err != nil {
		panic(err)
	}

}
