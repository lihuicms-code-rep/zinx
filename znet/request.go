package znet

import "zinx/ziface"

type Request struct {
	//与客户端连接
	conn ziface.IConnection

	//客户端请求数据
	msg ziface.IMessage
}


func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.msg.GetData()
}


func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgId()
}