package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"lean-zinx/utils"
	"lean-zinx/ziface"
)

/*
	封包拆包具体模块
	前面四个字节数据长，再来各四个字节数据ID
*/

type DataPack struct {
}

func NewDataPack() *DataPack {
	return &DataPack{}
}

// GetHeadLen 获取包的头的长度的方法
func (dp *DataPack) GetHeadLen() uint32 {
	// DataLen + ID 一共8个字节，即两个uint32
	// 后面可以用常量代替
	return 8
}

// Pack 封包
func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	// 创建一个存放bytes字节的缓冲
	dataBuf := bytes.NewBuffer([]byte{})

	// 将DataLen和ID写入
	// 小端传输
	err := binary.Write(dataBuf, binary.LittleEndian, msg.GetDataLen())
	if err != nil {
		return nil, err
	}
	err = binary.Write(dataBuf, binary.LittleEndian, msg.GetMessageId())
	if err != nil {
		return nil, err
	}
	// 写入Data
	err = binary.Write(dataBuf, binary.LittleEndian, msg.GetData())
	if err != nil {
		return nil, err
	}

	return dataBuf.Bytes(), nil
}

// UnPack 拆包：将包的Head读取，根据Head信息的数据长度读后续
func (dp *DataPack) UnPack(binaryData []byte) (ziface.IMessage, error) {

	// 创建从输入二进制数据的IO Reader
	reader := bytes.NewReader(binaryData)
	// 只解压Head信息
	msg := &Message{}

	// 读Len
	err := binary.Read(reader, binary.LittleEndian, &msg.DataLen)
	if err != nil {
		return nil, err
	}

	// 判断是否超过最大包长
	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("too large package size")
	}

	// 读ID
	err = binary.Read(reader, binary.LittleEndian, &msg.DataLen)
	if err != nil {
		return nil, err
	}

	return msg, nil

}
