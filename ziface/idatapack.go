package ziface

/*
	定义封包拆包的模块，直接面向TCP，针对TLV格式的封装
	拆包：先读取固定字节的消息内容长度，再根据消息内容长度再次读写。
*/

type IDataPack interface {
	// GetHeadLen 获取包的头的长度的方法
	GetHeadLen() uint32
	// Pack 封包
	Pack(msg IMessage) ([]byte, error)
	// UnPack 拆包
	UnPack([]byte) (IMessage, error)
}
