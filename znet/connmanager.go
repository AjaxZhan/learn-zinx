package znet

import (
	"errors"
	"fmt"
	"lean-zinx/ziface"
	"sync"
)

type ConnManager struct {
	// 维护连接集合的Map
	connections map[uint32]ziface.IConnection
	// 读写锁
	connLock sync.RWMutex
}

// NewConnManager 创建当前连接的方法
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

// AddConn 添加连接
func (connMgr *ConnManager) AddConn(conn ziface.IConnection) {
	// 上写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()
	// 将conn加入到map中
	connMgr.connections[conn.GetConnID()] = conn
	fmt.Println("Add Conn to manager successfully, connID=", conn.GetConnID(), " present num:", len(connMgr.connections))
}

// RemoveConn 删除连接
func (connMgr *ConnManager) RemoveConn(conn ziface.IConnection) {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	delete(connMgr.connections, conn.GetConnID())
	fmt.Println("connID:", conn.GetConnID(), " delete successfully!")

}

// GetConn 根据ID获取连接
func (connMgr *ConnManager) GetConn(connID uint32) (ziface.IConnection, error) {
	// 加读锁
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()
	if conn, ok := connMgr.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not found")
	}
}

// GetLen 得到连接总数
func (connMgr *ConnManager) GetLen() int {
	return len(connMgr.connections)
}

// ClearConn 清除并终止所有连接
func (connMgr *ConnManager) ClearConn() {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	// 删除conn并停止conn的工作
	for connID, conn := range connMgr.connections {
		conn.Stop()
		delete(connMgr.connections, connID)
	}
	fmt.Println("Clear all connections successfully!")
}
