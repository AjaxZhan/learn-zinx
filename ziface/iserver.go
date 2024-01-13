package ziface

type IServer interface {
	// Start 开启服务器
	Start()
	// Stop 停止服务器
	Stop()
	// Serve 运行服务器
	Serve()
	// AddRouter 给服务器注册一个路由方法
	AddRouter(router IRouter)
}
