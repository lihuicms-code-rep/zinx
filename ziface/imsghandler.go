package ziface

//消息管理抽象层
type IMsgHandler interface {
	//调度,执行对应Router消息处理方法
	DoMsgHandler(request IRequest)

	//添加消息的具体处理逻辑
	AddRouter(msgID uint32, router IRouter)

	//启动Worker工作池
	StartWorkerPool()

	//将消息发送给消息队列
	SendMsgToTaskQueue(IRequest)
}
