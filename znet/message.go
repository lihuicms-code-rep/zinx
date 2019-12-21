package znet

//处理数据的基本单位
//size   4     4     6       字节
//理解:| 201 |  6  | "iamagg" |
type Message struct {
	DataLen uint32 //消息长度(是指具体内容的长度)
	Id      uint32 //消息ID
	Data    []byte //消息内容
}

func NewMessage(id uint32, data []byte) *Message {
	return &Message{
		Id:id,
		DataLen:uint32(len(data)),
		Data:data,
	}
}

func (m *Message) GetMsgId() uint32 {
	return m.Id
}

func (m *Message) GetMsgLen() uint32 {
	return m.DataLen
}

func (m *Message) GetData() []byte {
	return m.Data
}

func (m *Message) SetMsgId(id uint32) {
	m.Id = id
}

func (m *Message) SetDataLen(len uint32) {
	m.DataLen = len
}

func (m *Message) SetData(data []byte) {
	m.Data = data
}
