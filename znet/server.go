package znet

import (
	"fmt"
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
}

// Start 启动服务器
func (s *Server) Start() {
	fmt.Printf("[Start] Server Listener at IP: %s,Port %d, is starting\n", s.IP, s.Port)

	// 异步
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
		fmt.Println("start Zinx server success ", s.Name, " succ listening...")
		// 阻塞等待客户端连接，处理客户端业务
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept error: ", err)
				continue
			}
			// 已经和客户端建立连接，执行业务
			// 这里做一个最基本的512字节长度的回显业务，用go来承载
			go func() {
				for {
					buf := make([]byte, 512)
					cnt, err := conn.Read(buf)
					if err != nil {
						fmt.Println("Recv buf error:", err)
						continue
					}
					if _, err := conn.Write(buf[0:cnt]); err != nil {
						fmt.Println("Write back buf error:", err)
						continue
					}

				}
			}()
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

// NewServer 初始化Server
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8899,
	}
	return s

}
