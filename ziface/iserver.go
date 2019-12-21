package ziface

//定义服务接口
type IServer interface {
	Start()                                 //启动服务方法
	Stop()                                  //停止服务方法
	Serve()                                 //服务运行方法
	AddRouter(msgId uint32, router IRouter) //服务注册路由方法,供客户端连接使用
	GetConnMgr() IConnManager               //获取服务管理
	SetOnConnStart(func(IConnection))       //设置OnConnStart
	CallOnConnStart(IConnection)            //调用OnConnStart
	SetOnConnStop(func(IConnection))       //设置OnConnStart
	CallOnConnStop(IConnection)            //调用OnConnStart
}
