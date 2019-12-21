package znet

import (
	"fmt"
	"zinx/utils"
	"zinx/ziface"
)

type MsgHandler struct {
	APIs map[uint32] ziface.IRouter   //msgId对应的处理方法
	TaskQueue []chan ziface.IRequest  //消息队列(一个worker对应处理一个消息队列)
	WorkerPoolSize uint32             //worker数量
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		APIs:make(map[uint32] ziface.IRouter),
		WorkerPoolSize:utils.GlobalObject.WorkerPoolSize,
		TaskQueue:make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}


//调度,执行对应Router消息处理方法
func (mh *MsgHandler) DoMsgHandler(request ziface.IRequest) {
	handler, ok := mh.APIs[request.GetMsgID()];
	if !ok {
		fmt.Println("api msgId =", request.GetMsgID(), " is not found")
		return
	}

	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

//添加消息的具体处理逻辑
func (mh *MsgHandler) AddRouter(msgID uint32, router ziface.IRouter) {
	if _, ok := mh.APIs[msgID]; ok {
		fmt.Println("repeated api, msgId =", msgID)
		return
	}

	mh.APIs[msgID] = router
	fmt.Println("add api msgId =", msgID, " success")
}


//启动worker工作池
func (mh *MsgHandler) StartWorkerPool() {
	//根据workerPoolSize来分别开启worker
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		//给当前的worker对应的消息队列开辟空间
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		//启动当前worker
		go mh.StartOneWorker(i, mh.TaskQueue[i])

	}

}


//启动一个worker工作流程
func (mh *MsgHandler) StartOneWorker(workID int, taskQueue chan ziface.IRequest) {
	fmt.Println(" WorkerID=", workID, " starting....")
	//阻塞等待消息到来处理
	for {
		select {
		case request := <- taskQueue:
			mh.DoMsgHandler(request)
		}
	}

}

//将请求交给TaskQueue
func (mh *MsgHandler) SendMsgToTaskQueue(request ziface.IRequest) {
	//1.消息平均分配给不同的worker
	//根据连接ID(connID)
	workID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("ConnID =", request.GetConnection().GetConnID(), " to workerID = ", workID)

	//2.消息交给这个worker对应的taskQueue
	mh.TaskQueue[workID] <- request
}