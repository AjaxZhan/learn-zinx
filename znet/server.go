package znet

import (
	"errors"
	"fmt"
	"lean-zinx/utils"
	"lean-zinx/ziface"
	"net"
)

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
	// 添加Router对象
	Router ziface.IRouter
}

// CallBackToClient 暂时写死的业务方法，后面理应由框架使用者传入
func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
	fmt.Println("[Conn Handle] CallBackToClient ...")
	if _, err := conn.Write(data[:cnt]); err != nil {
		fmt.Println("Write back buf error ", err)
		return errors.New("CallBackToClient error")
	}

	return nil
}

// Start 启动服务器
func (s *Server) Start() {
	// TODO 制作一个日志模块
	fmt.Printf("[Zinx]ServerName:%s, listener at IP:%s,Port:%d,is starting...\n",
		utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)
	fmt.Printf("[Zinx] Version %s, MaxConn:%d, MaxPackageSize:%d,\n",
		utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPackageSize)
	// server的协程
	go func() {
		// 获取TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp address error : ", err)
			return
		}
		// 监听服务器的地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("Listen ", s.IPVersion, " error:", err)
			return
		}
		// 比较邋遢的ConnID
		var cid uint32 = 0
		fmt.Println("start Zinx server success ", s.Name, " succ listening...")
		// 阻塞等待客户端连接，处理客户端业务
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept error: ", err)
				continue
			}
			// 将处理新连接的业务方法和conn绑定，得到我们的连接
			dealConn := NewConnection(conn, cid, s.Router)
			cid++

			// 启动业务
			go dealConn.Start()
		}
	}()

}

func (s *Server) Stop() {
	// TODO 将服务器资源、状态、开辟的连接信息进行停止或者回收

}

func (s *Server) Serve() {
	// 启动服务
	s.Start()

	// TODO 为什么不在Start内阻塞？这样可以在启动后做一些额外的业务，这样Server才有意义。Start只进行监听和处理业务功能。
	// Serve阻塞
	select {}
}

func (s *Server) AddRouter(router ziface.IRouter) {
	s.Router = router
	fmt.Println("Add Router successfully!")
}

// NewServer 初始化Server
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      utils.GlobalObject.Name,
		IPVersion: "tcp4",
		IP:        utils.GlobalObject.Host,
		Port:      utils.GlobalObject.TcpPort,
		Router:    nil,
	}
	return s

}
