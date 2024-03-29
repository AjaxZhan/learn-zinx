package znet

type Message struct {
	Id      uint32
	DataLen uint32
	Data    []byte
}

func (m *Message) GetMessageId() uint32 {
	return m.Id
}
func (m *Message) GetDataLen() uint32 {
	return m.DataLen
}
func (m *Message) GetData() []byte {
	return m.Data
}

func (m *Message) SetMessageId(id uint32) {
	m.Id = id
}
func (m *Message) SetDataLen(len uint32) {
	m.DataLen = len
}
func (m *Message) SetData(data []byte) {
	m.Data = data
}

// NewMessagePack 创建Message
func NewMessagePack(id uint32, data []byte) *Message {
	return &Message{
		Id:      id,
		DataLen: uint32(len(data)),
		Data:    data,
	}
}
