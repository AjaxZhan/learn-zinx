package znet

import "lean-zinx/ziface"

/*
	实现router时先嵌入BaseRouter，框架使用者根据需求重写这个方法。
*/

type BaseRouter struct {
}

// PreHandle 处理业务之前hook
func (br *BaseRouter) PreHandle(request ziface.IRequest) {

}

// Handle 处理业务主方法
func (br *BaseRouter) Handle(request ziface.IRequest) {

}

// PostHandle 处理业务之后的hook
func (br *BaseRouter) PostHandle(request ziface.IRequest) {

}
