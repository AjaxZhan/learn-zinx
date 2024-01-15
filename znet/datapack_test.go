package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

/*
封包拆包的单元测试
*/
func TestDataPack(t *testing.T) {
	// 模拟的服务器
	listener, err := net.Listen("tcp", "127.0.0.1:7788")
	if err != nil {
		fmt.Println("Server listen error:", err)
		return
	}
	go func() {
		for {
			conn, err2 := listener.Accept()
			if err2 != nil {
				fmt.Println("server accept error:", err2)
			}
			go func(conn net.Conn) {
				// 处理客户端请求
				/* 拆包 */
				dp := NewDataPack()
				for {
					headData := make([]byte, dp.GetHeadLen())
					_, err3 := io.ReadFull(conn, headData)
					if err3 != nil {
						fmt.Println("Read head error.")
						break
					}
					msg, err3 := dp.UnPack(headData)
					if err3 != nil {
						fmt.Println("Server unpack error:", err3)
						return
					}
					if msg.GetDataLen() > 0 {
						// msg 内有数据，需要进行第二次读
						msg := msg.(*Message)
						// 根据dataLen开辟空间，再次从IO读取
						msg.Data = make([]byte, msg.GetDataLen())
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("Server unpack data error:", err)
							return
						}
						// 完整消息读取完毕，打印
						fmt.Println("===>Recv msgId:", msg.Id, " DataLen:", msg.GetDataLen(), " data:", msg.GetData())
					}
				}
				// 1. 根据
			}(conn)
		}
	}()
	// 模拟客户端
	conn, err := net.Dial("tcp", "127.0.0.1:7788")
	if err != nil {
		fmt.Println("client dial error:", err)
		return
	}
	// 创建封包对象
	dp := NewDataPack()

	// 模拟粘包过程，封装两个msg包
	msg1 := &Message{
		Id:      1,
		DataLen: 4,
		Data:    []byte{'z', 'i', 'n', 'x'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("Client pack message1 error:", err)
		return
	}
	msg2 := &Message{
		Id:      2,
		DataLen: 9,
		Data:    []byte{'c', 'a', 'g', 'u', 'r', 'z', 'h', 'a', 'n'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("Client pack message2 error:", err)
		return
	}
	sendData1 = append(sendData1, sendData2...)
	_, err = conn.Write(sendData1)
	if err != nil {
		fmt.Println("client write error:", err)
	}

	// 客户端阻塞
	select {}
}
