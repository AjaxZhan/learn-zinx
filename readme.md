# Learn-zinx

This repository is used to record my process of learning the zinx framework.

zinx framework repository address：[https://github.com/aceld/zinx](https://github.com/aceld/zinx)

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

## V0.3 基础router模块(单一Router)

### Request请求封装

目的：将conn和数据绑定在一起，以request作为请求的原子操作。
属性：
- 连接 iconnection
- 请求数据
方法：
- 获取当前连接
- 获取当前数据

### Router模块的定义

IRouter：抽象层
- 方法：
  - 处理业务之前的方法
  - 处理业务的方法
  - 处理业务之后的方法

BaseRouter：实现层
实现router时先嵌入BaseRouter，框架使用者根据需求重写这个方法。

### zinx集成Router

1. IServer增添添加路由方法
2. Server添加Router成员，去掉之前的回调。
3. Connection类绑定Router成员，去掉handleAPI
4. Connection中调用Router处理业务——模板方法设计模式

### 测试——使用ZinxV0.3开发

1. 创建server
2. 创建自定义Router，继承BaseRouter，重写三个方法
3. 添加Router
4. 启动Server

## V0.4 全局配置模块

这里以JSON格式为配置文件，用户编写zinx.json。
放到/config/zinx.json中。

内容包括：
- Name：服务器名
- Host：监听的IP地址
- TcpPort：端口
- MaxConn：允许最大客户端数量

### 创建全局配置模块

`utils/globalobj.go`
- 提供全局globalobj对象
- init：读取json配置，序列化到globalobj中

## V0.5 消息封装

### 定义消息结构

定义消息结构：Message
- 属性：
  - 消息ID
  - 消息长度
  - 消息内容
- 方法：
  - set
  - get

### 解决TCP粘包问题

消息的TLV序列化：
一个完整的消息 = Head(DataLen + ID) + Body(Data)

