package znet

import (
	"fmt"
	"lean-zinx/utils"
	"lean-zinx/ziface"
)

type MsgHandler struct {
	// 存放每个MsgID对应的处理方法
	Apis map[uint32]ziface.IRouter
	// 负责Worker取任务的消息对垒
	TaskQueue []chan ziface.IRequest
	// 业务工作Worker池的worker数量
	WorkPoolSize uint32
}

// NewMsgHandle 初始化后一个Handler
func NewMsgHandle() *MsgHandler {
	return &MsgHandler{
		Apis:         make(map[uint32]ziface.IRouter),
		WorkPoolSize: utils.GlobalObject.WorkPoolSize,
		TaskQueue:    make([]chan ziface.IRequest, utils.GlobalObject.WorkPoolSize),
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

// StartWorkerPool 启动一个Worker工作池，只需要开启一次
func (m *MsgHandler) StartWorkerPool() {
	for i := 0; i < int(m.WorkPoolSize); i++ {
		// 给当前Worker对应的Channel开辟空间
		m.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		// 尝试启动go程
		go m.StartOneWorker(i, m.TaskQueue[i])
	}
}

// StartOneWorker 启动一个Worker进程
func (m *MsgHandler) StartOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("Worker:", workerID, " is started ...")
	// 阻塞消息，等待channel传递的消息
	for {
		select {
		case request := <-taskQueue:
			//出列一个客户端的request，执行的handler方法
			m.DoMsgHandle(request)
		}
	}
}

// SendMsgToTaskQueue 将消息交给消息队列，由Worker进行处理
func (m *MsgHandler) SendMsgToTaskQueue(request ziface.IRequest) {
	// 将消息平均分配给不同的Worker，根据客户端建立的ConnID来进行分配(按理说要加个requestID，根据这个分配)
	workID := request.GetConnection().GetConnID() % m.WorkPoolSize
	fmt.Println("Add ConnID=", request.GetConnection().GetConnID(), " requestMsgID=", request.GetMsgID(),
		" to workerID:", workID)
	// 将消息发送给对应Worker的taskQueue
	m.TaskQueue[workID] <- request
}
