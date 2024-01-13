
## V0.1 基础server模块

> V0.1版本实现了基本的Server服务，实现了消息回显的过程。

文件结构：
- demo： 存放基于框架写的server端和client端用于测试。
- ziface：抽象层
- znet：实现层

基于Golang的net包进行开发：
- net.ResolveTCPAddr：创建TCP的addr
- net.ListenTCP：监听TCPAddr，得到listener
- listener.AcceptTCP：接收客户端的连接，得到conn
- net.Dial：客户端连接服务器，得到conn

server的定义：
- 属性：
  - 服务器名
  - IP地址
  - IP版本
  - 端口
- 方法
  - 启动
  - 暂停
  - 运行

## V0.2 连接封装和业务绑定

连接封装：`connection.go`
- 属性：
  - conn套接字 *net.Conn
  - 连接ID uint32
  - 当前连接状态 bool
  - 与连接绑定的业务方法 handleFunc
  - 异步过程，有个channel捕获退出信号 chan bool
- 方法：
  - 启动连接：启动两个go程
  - 停止连接：关闭套接字和管道
  - 获取当前连接的conn对象
  - 得到连接ID
  - 得到客户端连接地址和端口
  - 发送数据

系统架构设计上分为两个Go程去执行任务，一个负责读业务，一个负责写业务。
V0.2实现了简单的读业务，即回显功能。

server.go在start方法中启动一个go程创建套接字，通过NewConnection实例化connection对象，并通过go程开启这个connection对象的start方法。

connection.go的start方法负责开启读业务和写业务这两个Go程。
目前只实现了读业务的Go程，且业务处理方法的回调在server.go写死，实现写回功能。





