package ziface

type IServer interface {
	// Start 开启服务器
	Start()
	// Stop 停止服务器
	Stop()
	// Serve 运行服务器
	Serve()
	// AddRouter 给服务器注册一个路由方法
	AddRouter(msgID uint32, router IRouter)
	// GetConnMgr 获取ConnManager
	GetConnMgr() IConnManager
	// SetOnConnStart 注册Hook
	SetOnConnStart(func(connection IConnection))
	// SetOnConnStop 注册销毁前Hook
	SetOnConnStop(func(connection IConnection))
	// CallOnConnStart 调用Hook
	CallOnConnStart(connection IConnection)
	// CallOnConnStop 调用销毁前Hook
	CallOnConnStop(connection IConnection)
}
