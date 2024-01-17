package main

import (
	"fmt"
	"io"
	"lean-zinx/znet"
	"net"
	"time"
)

// 模拟客户端
func main() {
	fmt.Println("Client01 Start...")
	time.Sleep(1 * time.Second)
	// 连接服务器，得到conn
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("Connect Server error : ", err)
		return
	}
	for {
		dp := znet.NewDataPack()
		// 封包
		binaryMessage, err := dp.Pack(znet.NewMessagePack(1, []byte("Zinx V0.6 client Message...")))
		if err != nil {
			fmt.Println("Client pack error:", err)
			return
		}
		// 写数据
		_, err = conn.Write(binaryMessage)
		if err != nil {
			fmt.Println("Client write message error:", err)
			return
		}
		// 服务器回复Message：ID=1
		// 先读取头部
		binaryHead := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, binaryHead)
		if err != nil {
			fmt.Println("Client read head error:", err)
		}
		// 将二进制Head拆包到Message
		msg, err := dp.UnPack(binaryHead)
		if err != nil {
			fmt.Println("Client Unpack error:", err)
		}
		// 再进行第二次读取，读取Data
		if msg.GetDataLen() > 0 {
			msg := msg.(*znet.Message)
			msg.Data = make([]byte, msg.GetDataLen())
			_, err = io.ReadFull(conn, msg.Data)
			if err != nil {
				fmt.Println("Client read data error:", err)
			}
			fmt.Println("Recv Message ID:", msg.Id, " Len=", msg.GetMessageId(), " data=", string(msg.GetData()))
		}
		// CPU阻塞，CPU干其它事情
		time.Sleep(1 * time.Second)
	}

}
