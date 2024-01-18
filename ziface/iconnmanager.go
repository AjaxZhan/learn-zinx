package ziface

/*
	连接管理模块抽象层
*/

type IConnManager interface {
	// AddConn 添加连接
	AddConn(conn IConnection)
	// RemoveConn 删除连接
	RemoveConn(conn IConnection)
	// GetConn 根据ID获取连接
	GetConn(connID uint32) (IConnection, error)
	// GetLen 得到连接总数
	GetLen() int
	// ClearConn 清除并终止所有连接
	ClearConn()
}
