package znet

import (
	"errors"
	"fmt"
	"io"
	"lean-zinx/utils"
	"lean-zinx/ziface"
	"net"
	"sync"
)

// Connection 连接模块
type Connection struct {
	// conn隶属于的server
	TcpServer ziface.IServer
	// conn对象
	Conn *net.TCPConn
	// 连接ID
	ConnID uint32
	// 连接状态
	IsClosed bool
	// 告知连接是否退出的 channel（Reader告知Writer）
	ExitChan chan bool
	// 无缓冲管道，用于读、写协程之间的通信
	msgChan chan []byte
	// 该连接处理方法的Router
	Router ziface.IRouter
	// 消息的管理msgID和对应的处理业务API关系
	MsgHandler ziface.IMsgHandler
	// 连接属性集合
	property map[string]interface{}
	// 保护连接属性的互斥锁
	propertyLock sync.RWMutex
}

// NewConnection 初始化连接的方法
func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandler) *Connection {
	c := &Connection{
		TcpServer:  server,
		Conn:       conn,
		ConnID:     connID,
		IsClosed:   false,
		ExitChan:   make(chan bool, 1),
		MsgHandler: msgHandler,
		msgChan:    make(chan []byte),
		property:   make(map[string]interface{}),
	}

	//将conn加入到connManager中
	c.TcpServer.GetConnMgr().AddConn(c)
	return c
}

// StartWriter 写消息的的协程，专门用于将消息发送给客户端
func (c *Connection) StartWriter() {
	fmt.Println("[Writer GoRoutine is running...]")
	defer fmt.Println(c.RemoteAddr().String(), " [conn Writer exit]")
	// 阻塞等待Channel的消息
	for {
		select {
		case data := <-c.msgChan:
			// 写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error:", err)
				return
			}
		case <-c.ExitChan:
			// Reader已经退出, Writer也要退出
			return
		}

	}
}

// StartReader 读消息的协程
func (c *Connection) StartReader() {
	fmt.Println("[Reader GoRoutine is running...]")
	defer fmt.Println("connID=", c.ConnID, " [Reader is exit], the remote addr is ", c.RemoteAddr().String())
	defer c.Stop()

	for {
		// 创建拆包解包对象
		dp := NewDataPack()
		// 读取二进制流头部
		dataHead := make([]byte, dp.GetHeadLen())
		_, err := io.ReadFull(c.GetTCPConnection(), dataHead)
		if err != nil {
			fmt.Println("Read Message Head error:", err)
			break
		}
		// 将Msg得到msgID和msgDataLen，封装为Message
		msg, err := dp.UnPack(dataHead)
		if err != nil {
			fmt.Println("unpack error:", err)
			break
		}
		// 再次读取Message的Data
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			_, err := io.ReadFull(c.GetTCPConnection(), data)
			if err != nil {
				fmt.Println("read msg data error:", err)
				break
			}
		}
		msg.SetData(data)

		// 封装request
		req := Request{
			conn: c,
			msg:  msg,
		}

		// 判断是否开启工作池
		if utils.GlobalObject.WorkPoolSize > 0 {
			// 已经开启工作池机制，将消息发送给工作池
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			// 没有开启工作池，直接开启一个go程即可
			go c.MsgHandler.DoMsgHandle(&req)
		}

	}

}

// SendMsg 提供发送数据方法，将要发送给客户端的数据进行封包再发送
func (c *Connection) SendMsg(msgId uint32, data []byte) error {

	if c.IsClosed == true {
		return errors.New("connection closed when send msg")
	}

	// 封包
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMessagePack(msgId, data))
	if err != nil {
		fmt.Println("Send Message error when Pack:", err, " MessageId:", msgId)
		return errors.New("pack error msg")
	}

	// 将数据发送给客户端
	c.msgChan <- binaryMsg

	//_, err = c.Conn.Write(binaryMsg)
	//if err != nil {
	//	fmt.Println("msgId:", msgId, "error:", err)
	//	return errors.New("conn write error")
	//}

	return nil
}

// Start 启动连接，让当前连接进入工作状态
func (c *Connection) Start() {
	fmt.Println("Connection Start ... ,ConnID=", c.ConnID)
	// 启动从当前连接读数据的业务的携程
	go c.StartReader()
	// 启动写数据业务
	go c.StartWriter()
	// 调用创建连接之前的Hook
	c.TcpServer.CallOnConnStart(c)
}

// Stop 停止连接，结束当前连接的工作
func (c *Connection) Stop() {
	// 调用开发者注册的销毁前Hook
	c.TcpServer.CallOnConnStop(c)
	fmt.Println("Connection Stop...,ConnId=", c.ConnID)
	if c.IsClosed == true {
		return
	}
	// 关闭socket连接
	c.IsClosed = true
	_ = c.Conn.Close()

	// 告知Writer已经关闭
	c.ExitChan <- true

	// 将当前连接从管理器中删除掉
	c.TcpServer.GetConnMgr().RemoveConn(c)
	// 关闭管道
	close(c.ExitChan)
	close(c.msgChan)
}

// GetTCPConnection 获取当前连接绑定的conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// GetConnID 获取当前连接的连接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// RemoteAddr 获取远程的TCP状态，包括IP和Port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// SetProperty 设置连接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	// 添加一个连接属性
	c.property[key] = value

}

// GetProperty 获取连接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	if value, ok := c.property[key]; ok {
		return value, nil
	}
	return nil, errors.New("no property found")
}

// RemoveProperty 移除连接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}
