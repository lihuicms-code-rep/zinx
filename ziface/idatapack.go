package ziface

//消息封装抽象层
type IDataPack interface {
	GetHeadLen() uint32              //获取头部长度
	Pack(IMessage) ([]byte, error)   //封包
	Unpack([]byte) (IMessage, error) //拆包
}
