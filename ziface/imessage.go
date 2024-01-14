package ziface

/*
	将请求的消息封装到Message中，定义一个抽象层模块
*/

type IMessage interface {
	GetMessageId() uint32
	GetDataLen() uint32
	GetData() []byte

	SetMessageId(uint32)
	SetDataLen(uint32)
	SetData([]byte)
}
