package ziface

import "net"

// 定义连接模块的抽象层

type IConnection interface {
	// Start 启动连接，让当前连接进入工作状态
	Start()
	// Stop 停止连接，结束当前连接的工作
	Stop()
	// GetTCPConnection 获取当前连接绑定的conn
	GetTCPConnection() *net.TCPConn
	// GetConnID 获取当前连接的连接ID
	GetConnID() uint32
	// RemoteAddr 获取远程的TCP状态，包括IP和Port
	RemoteAddr() net.Addr
	// SendMsg 发送数据给远程客户端
	SendMsg(msgId uint32, data []byte) error
	// SetProperty 设置连接属性
	SetProperty(key string, value interface{})
	// GetProperty 获取连接属性
	GetProperty(key string) (interface{}, error)
	// RemoveProperty 移除连接属性
	RemoveProperty(key string)
}

// HandleFunc 定义一个处理连接业务的抽象方法，和连接绑定
// 参数：连接、数据、数据长度
//type HandleFunc func(*net.TCPConn, []byte, int) error
