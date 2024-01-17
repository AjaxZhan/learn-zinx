package ziface

/*
	IRequest接口：把Client请求的连接信息和数据进行封装
*/

type IRequest interface {
	// GetConnection 得到当前连接
	GetConnection() IConnection
	// GetData 得到消息数据
	GetData() []byte
	// GetMsgID 获取消息ID
	GetMsgID() uint32
}
