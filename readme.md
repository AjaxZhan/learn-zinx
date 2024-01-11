
## V0.1 基础server模块

> V0.1版本实现了基本的Server服务，实现了消息回显的过程。

框架分为抽象层和实现层，分别定义在ziface和znet里面。

demo文件夹存放基于框架写的server端和client端用于测试。

基于Golang的net包进行开发：
- net.ResolveTCPAddr：创建TCP的addr
- net.ListenTCP：监听TCPAddr，得到listener
- listener.AcceptTCP：接收客户端的连接，得到conn

客户端：
- net.Dial：连接服务器，得到conn

其他包：fmt、time

## V0.2 连接封装和业务绑定

连接封装：
- 方法：
  - 启动连接
  - 停止连接
  - 获取当前连接的conn对象
  - 得到连接ID
  - 得到客户端连接地址和端口
  - 发送数据
- 属性：
  - conn套接字 *net.Conn
  - 连接ID uint32
  - 当前连接状态 bool
  - 与连接绑定的业务方法 handleFunc
  - 异步过程，有个channel捕获退出信号 chan bool



