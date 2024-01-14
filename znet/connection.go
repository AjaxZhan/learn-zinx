package znet

import (
	"fmt"
	"lean-zinx/utils"
	"lean-zinx/ziface"
	"net"
)

// Connection 连接模块
type Connection struct {
	// conn对象
	Conn *net.TCPConn
	// 连接ID
	ConnID uint32
	// 连接状态
	IsClosed bool
	// 告知连接是否退出的 channel
	ExitChan chan bool
	// 该连接处理方法的Router
	Router ziface.IRouter
}

// NewConnection 初始化连接的方法
func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connID,
		IsClosed: false,
		ExitChan: make(chan bool, 1),
		Router:   router,
	}
	return c
}

// StartReader 连接的读业务
func (c *Connection) StartReader() {
	fmt.Println("Reader GoRoutine is running...")
	defer fmt.Println("connID=", c.ConnID, " Reader is exit, the remote addr is ", c.RemoteAddr().String())
	defer c.Stop()

	for {
		// 读客户端数据到buf中，最大512字节
		buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("Recv buf err", err)
			// 这个包读失败
			continue
		}

		// 封装request
		req := Request{
			conn: c,
			data: buf,
		}

		// 开GO程，从路由中找到注册的conn对应的router调用
		// 模板方法设计模式
		go func(request ziface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)

		// 调用当前连接所绑定的HandleAPI
		//if err := c.handleAPI(c.Conn, buf, cnt); err != nil {
		//	fmt.Println("ConnID", c.ConnID, " handle is error:", err)
		//	break
		//}
	}

}

// Start 启动连接，让当前连接进入工作状态
func (c *Connection) Start() {
	fmt.Println("Connection Start ... ,ConnID=", c.ConnID)
	// 启动从当前连接读数据的业务的携程
	go c.StartReader()
	// TODO 启动写数据业务
}

// Stop 停止连接，结束当前连接的工作
func (c *Connection) Stop() {
	fmt.Println("Connection Stop...,ConnId=", c.ConnID)
	if c.IsClosed == true {
		return
	}
	// 关闭socket连接

	c.IsClosed = true
	_ = c.Conn.Close()

	// 关闭管道
	close(c.ExitChan)
}

// GetTCPConnection 获取当前连接绑定的conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// GetConnID 获取当前连接的连接ID
func (c *Connection) GetConnID() uint32 {
	return c.GetConnID()
}

// RemoteAddr 获取远程的TCP状态，包括IP和Port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// Send 发送数据给远程客户端
func (c *Connection) Send(data []byte) error {
	return nil
}
