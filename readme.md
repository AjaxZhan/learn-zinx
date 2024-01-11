
## V0.1 基础server模块

> V0.1版本实现了基本的Server服务，实现了消息返送的过程。

框架分为抽象层和实现层，分别定义在ziface和znet里面。

demo文件夹存放基于框架写的server端和client端用于测试。

基于Golang的net包进行开发：
- net.ResolveTCPAddr：创建TCP的addr
- net.ListenTCP：监听TCPAddr，得到listener
- listener.AcceptTCP：接收客户端的连接

客户端：
- net.Dial：连接服务器，得到conn

其他包：fmt、time

