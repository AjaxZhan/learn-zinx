package ziface

/*
	消息管理模块抽象层
*/

type IMsgHandler interface {
	// DoMsgHandle 调度/执行对应的Router消息处理方法
	DoMsgHandle(request IRequest)
	// AddRouter 为消息添加具体的处理路由
	AddRouter(msgID uint32, router IRouter)
	// StartWorkerPool 启动Worker工作池
	StartWorkerPool()
	// SendMsgToTaskQueue 将消息交给消息队列，由Worker进行处理
	SendMsgToTaskQueue(request IRequest)
}
