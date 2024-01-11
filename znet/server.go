package znet

// Server 定义IServer接口的实现
type Server struct {
	// 服务器名字
	Name string
	// IP版本
	IPVersion string
	// IP地址
	IP string
	// 端口
	Port int
}

func (s *Server) Start() {

}

func (s *Server) Stop() {

}

func (s *Server) Serve() {

}
