package ziface

//将客户端请求的连接信息和连接数据封装抽象层

type IRequest interface {
	//得到当前连接
	GetConnection() IConnection

	//得到请求的数据
	GetData() []byte

	//得到请求消息ID
	GetMsgID() uint32
}
