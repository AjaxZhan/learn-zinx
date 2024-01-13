package main

import (
	"fmt"
	"lean-zinx/ziface"
	"lean-zinx/znet"
)

/**
基于Zinx 框架来开发的服务端应用程序
*/

// PingRouter Ping test 自定义路由
type PingRouter struct {
	// 继承BaseRouter
	znet.BaseRouter
}

func (pr *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("Call PreHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("Before ping ... \n"))
	if err != nil {
		fmt.Println("Call back BeforePing error.")
	}
}
func (pr *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Handle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("Ping ... \n"))
	if err != nil {
		fmt.Println("Call back Ping error.")
	}
}
func (pr *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("Call PostHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("After ping ...\n"))
	if err != nil {
		fmt.Println("Call back AfterPing error.")
	}
}

func main() {
	// 创建Server句柄
	s := znet.NewServer("[zine V0.3]")

	// 添加自定义的Router
	s.AddRouter(&PingRouter{})

	// 启动服务
	s.Serve()
}
