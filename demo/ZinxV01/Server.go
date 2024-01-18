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

type HelloZinxRouter struct {
	znet.BaseRouter
}

//	func (pr *PingRouter) PreHandle(request ziface.IRequest) {
//		fmt.Println("Call PreHandle...")
//		_, err := request.GetConnection().GetTCPConnection().Write([]byte("Before ping ... \n"))
//		if err != nil {
//			fmt.Println("Call back BeforePing error.")
//		}
//	}
func (pr *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call PingRouter Handle...")
	// 先读取客户端的数据
	fmt.Println("recv from client , msgId:", request.GetMsgID(), " data=", string(request.GetData()))
	// 回写消息
	err := request.GetConnection().SendMsg(200, []byte("ping..."))
	if err != nil {
		fmt.Println(err)
	}
}

func (hzr *HelloZinxRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call HelloZinxRouter Handle...")
	// 先读取客户端的数据
	fmt.Println("recv from client , msgId:", request.GetMsgID(), " data=", string(request.GetData()))
	// 回写消息
	err := request.GetConnection().SendMsg(201, []byte("Hello Zinx..."))
	if err != nil {
		fmt.Println(err)
	}
}

//func (pr *PingRouter) PostHandle(request ziface.IRequest) {
//	fmt.Println("Call PostHandle...")
//	_, err := request.GetConnection().GetTCPConnection().Write([]byte("After ping ...\n"))
//	if err != nil {
//		fmt.Println("Call back AfterPing error.")
//	}
//}

// DoConnectionBegin 创建连接后执行
func DoConnectionBegin(conn ziface.IConnection) {
	fmt.Println("==> DoConnectionBegin is called.")
	err := conn.SendMsg(202, []byte("服务器创建连接后方法被执行！"))
	if err != nil {
		fmt.Println("DoConnectionBegin error:", err)
	}
}

func DoConnectionLost(conn ziface.IConnection) {
	fmt.Println("==> DoConnectionLost is called. connID=", conn.GetConnID())
}

// 创建连接断开前执行

func main() {
	// 创建Server句柄
	s := znet.NewServer("[zine V0.9]")

	// 注册连接创建的钩子和销毁的钩子
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	// 添加自定义的Router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloZinxRouter{})
	// 启动服务
	s.Serve()
}
