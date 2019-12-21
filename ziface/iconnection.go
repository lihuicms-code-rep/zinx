package ziface

import "net"

//定义连接模块的抽象层
type IConnection interface {
	//启动连接
	Start()

	//停止连接
	Stop()

	//获取当前连接所绑定的socket
	GetTCPConnection() *net.TCPConn

	//获取当前连接的ID
	GetConnID() uint32

	//获取当前连接对端客户端的地址
	RemoteAddr() net.Addr

	//发送数据给客户端
	SendMsg(uint32, []byte) error

	//设置连接属性
	SetProperty(string, interface{})

	//获取连接属性
	GetProperty(string) (interface{}, error)

	//移除连接属性
	RemoveProperty(string)
}

//处理业务连接的方法
type HandleFunc func(*net.TCPConn, []byte, int) error
