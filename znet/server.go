package znet

import (
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
	// 消息管理模块，绑定MsgID和对应的API关系
	MsgHandler ziface.IMsgHandler
	// 该Server的连接管理器
	ConnManager ziface.IConnManager
	// Conn创建之前的Hook
	OnConnStart func(conn ziface.IConnection)
	// Conn创建之后的Hook
	OnConnStop func(conn ziface.IConnection)
}

// NewServer 初始化Server
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:        utils.GlobalObject.Name,
		IPVersion:   "tcp4",
		IP:          utils.GlobalObject.Host,
		Port:        utils.GlobalObject.TcpPort,
		MsgHandler:  NewMsgHandle(),
		ConnManager: NewConnManager(),
	}
	return s

}

// CallBackToClient 暂时写死的业务方法，后面理应由框架使用者传入
//func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
//	fmt.Println("[Conn Handle] CallBackToClient ...")
//	if _, err := conn.Write(data[:cnt]); err != nil {
//		fmt.Println("Write back buf error ", err)
//		return errors.New("CallBackToClient error")
//	}
//
//	return nil
//}

// Start 启动服务器
func (s *Server) Start() {
	// TODO 制作一个日志模块
	fmt.Printf("[Zinx]ServerName:%s, listener at IP:%s,Port:%d,is starting...\n",
		utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)
	fmt.Printf("[Zinx] Version %s, MaxConn:%d, MaxPackageSize:%d,\n",
		utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPackageSize)
	// server的协程
	go func() {
		// 开启消息队列和worker工作池
		s.MsgHandler.StartWorkerPool()
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

			// 判断是否超过最大连接个数，如果超过则关闭
			if s.ConnManager.GetLen() >= utils.GlobalObject.MaxConn {
				// TODO 给客户端响应一个超出连接的错误包
				_ = conn.Close()
				fmt.Println("Too many connection, len of conn exceed maxConn")
				continue
			}

			// 将处理新连接的业务方法和conn绑定，得到我们的连接
			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++

			// 启动业务
			go dealConn.Start()
		}
	}()

}

func (s *Server) Stop() {
	// 将服务器资源、状态、开辟的连接信息进行停止或者回收
	fmt.Println("[STOP] zinx server name:", s.Name)
	s.ConnManager.ClearConn()
}

func (s *Server) Serve() {
	// 启动服务
	s.Start()

	// TODO 为什么不在Start内阻塞？这样可以在启动后做一些额外的业务，这样Server才有意义。Start只进行监听和处理业务功能。
	// Serve阻塞
	select {}
}

func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("Add Router successfully!")
}

// GetConnMgr 获取ConnManager
func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnManager
}

// SetOnConnStart 注册Hook
func (s *Server) SetOnConnStart(hook func(connection ziface.IConnection)) {
	s.OnConnStart = hook
}

// SetOnConnStop 注册销毁前Hook
func (s *Server) SetOnConnStop(hook func(connection ziface.IConnection)) {
	s.OnConnStop = hook
}

// CallOnConnStart 调用Hook
func (s *Server) CallOnConnStart(connection ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("Call On Connection Start.")
		s.OnConnStart(connection)
	}
}

// CallOnConnStop 调用销毁前Hook
func (s *Server) CallOnConnStop(connection ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("Call On Connection Stop.")
		s.OnConnStop(connection)
	}
}
