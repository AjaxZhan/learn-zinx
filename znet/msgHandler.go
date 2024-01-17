package znet

import (
	"fmt"
	"lean-zinx/ziface"
)

type MsgHandler struct {
	// 存放每个MsgID对应的处理方法
	Apis map[uint32]ziface.IRouter
}

// NewMsgHandle 初始化后一个Handler
func NewMsgHandle() *MsgHandler {
	return &MsgHandler{
		Apis: make(map[uint32]ziface.IRouter),
	}
}

func (m *MsgHandler) DoMsgHandle(request ziface.IRequest) {
	msgID := request.GetMsgID()
	handler, ok := m.Apis[msgID]
	if !ok {
		fmt.Println("MsgID:", msgID, " 's handler not found!")
		return
	}
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// AddRouter 为消息添加具体的处理路由
func (m *MsgHandler) AddRouter(msgID uint32, router ziface.IRouter) {
	// 判断当前ID是否已经被注册
	if _, ok := m.Apis[msgID]; ok {
		// ID已经注册了
		panic(fmt.Sprintf("Repeat API! MsgID=%d\n", msgID))
	}
	// 添加绑定关系
	m.Apis[msgID] = router
	fmt.Println("Add API successfully! msgID:", msgID)
}
