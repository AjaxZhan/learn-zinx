package utils

import (
	"encoding/json"
	"lean-zinx/ziface"
	"os"
)

/*
	存储框架的全局参数
	部分参数由用户通过json配置
*/

type GlobalObj struct {
	/*
		Server配置
	*/
	TcpServer ziface.IServer // 全局server对象
	Host      string         // 服务器监听IP
	TcpPort   int            // 服务器监听的端口
	Name      string         // 服务器名
	/*
		Zinx配置
	*/
	Version        string // zinx版本号
	MaxConn        int    // 当前服务器允许的最大连接数
	MaxPackageSize uint32 // 数据包最大值

}

// GlobalObject 定义对外GlobalObj
var GlobalObject *GlobalObj

// Reload 从zinx.json加载配置文件
func (g *GlobalObj) Reload() {
	data, err := os.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}
	// 解析JSON
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

// 导包时，初始化GlobalObj
func init() {
	// 默认配置
	GlobalObject = &GlobalObj{
		Name:           "ZinxServerApp",
		Version:        "V0.6",
		Host:           "0.0.0.0",
		TcpPort:        8999,
		MaxPackageSize: 512,
		MaxConn:        1000,
	}
	// 尝试从config/zinx.json加载用户自定义参数
	GlobalObject.Reload()
}
