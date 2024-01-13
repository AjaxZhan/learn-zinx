package ziface

/*
	路由抽象接口，路由内的数据为IRequest
*/

type IRouter interface {
	// PreHandle 处理业务之前hook
	PreHandle(request IRequest)
	// Handle 处理业务主方法
	Handle(request IRequest)
	// PostHandle 处理业务之后的hook
	PostHandle(request IRequest)
}
