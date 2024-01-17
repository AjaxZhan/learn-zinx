package znet

import (
	"errors"
	"fmt"
	"io"
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
	// 消息的管理msgID和对应的处理业务API关系
	MsgHandler ziface.IMsgHandler
}

// NewConnection 初始化连接的方法
func NewConnection(conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandler) *Connection {
	c := &Connection{
		Conn:       conn,
		ConnID:     connID,
		IsClosed:   false,
		ExitChan:   make(chan bool, 1),
		MsgHandler: msgHandler,
	}
	return c
}

// StartReader 连接的读业务
func (c *Connection) StartReader() {
	fmt.Println("Reader GoRoutine is running...")
	defer fmt.Println("connID=", c.ConnID, " Reader is exit, the remote addr is ", c.RemoteAddr().String())
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

		// 根据绑定好的msgID找到处理API业务执行
		go c.MsgHandler.DoMsgHandle(&req)

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
	_, err = c.Conn.Write(binaryMsg)
	if err != nil {
		fmt.Println("msgId:", msgId, "error:", err)
		return errors.New("conn write error")
	}

	return nil
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
