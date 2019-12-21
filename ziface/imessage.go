package ziface


//Message抽象层,Getter&Setter
type IMessage interface {
	GetMsgId() uint32
	GetMsgLen() uint32
	GetData() []byte

	SetMsgId(uint32)
	SetDataLen(uint32)
	SetData([]byte)
}