package znet

import "lean-zinx/ziface"

type Request struct {
	// 连接
	conn ziface.IConnection
	// 数据
	msg ziface.IMessage
}

// GetConnection 得到当前连接
func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

// GetData 得到消息数据
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

// GetMsgID 获取消息ID
func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMessageId()
}
