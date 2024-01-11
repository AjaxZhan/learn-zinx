package main

import (
	"fmt"
	"net"
	"time"
)

// 模拟客户端
func main() {
	fmt.Println("Client Start...")
	time.Sleep(1 * time.Second)
	// 连接服务器，得到conn
	conn, err := net.Dial("tcp", "127.0.0.1:8899")
	if err != nil {
		fmt.Println("Connect Server error : ", err)
		return
	}
	// 使用write方法写数据
	for {
		_, err := conn.Write([]byte("Hello, zinx V0.1 ..."))
		if err != nil {
			fmt.Println("Write error:", err)
			return
		}
		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Read err : ", err)
			return
		}
		fmt.Printf("Server Return Message%s, cnt= %d: \n", buf, cnt)
		// CPU阻塞，CPU干其它事情
		time.Sleep(1 * time.Second)
	}

}
