package main

import "lean-zinx/znet"

/**
基于Zinx 框架来开发的服务端应用程序
*/

func main() {
	// 创建Server句柄
	s := znet.NewServer("[zine V0.1]")
	// 启动服务
	s.Serve()
}
