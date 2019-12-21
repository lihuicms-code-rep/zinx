package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"zinx/utils"
	"zinx/ziface"
)

//连接模块抽象层的实现
type Connection struct {
	//conn隶属于哪个server(可以明晰连接与服务的关系)
	TCPServer ziface.IServer

	//连接的socket
	Conn *net.TCPConn

	//连接ID
	ConnID uint32

	//连接状态
	isClosed bool

	//异步捕捉退出,由Reader告知Writer
	ExitChan chan bool

	//该连接处理对应Router
	MsgHandler ziface.IMsgHandler

	//用于读写goroutine之间消息通信
	msgChan chan []byte

	//连接属性集合
	property map[string]interface{}

	//保护连接属性修改锁
	propertyLock sync.RWMutex
}

//初始化
func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandler) *Connection {
	c := &Connection{
		TCPServer:  server,
		Conn:       conn,
		ConnID:     connID,
		isClosed:   false,
		ExitChan:   make(chan bool, 1),
		MsgHandler: msgHandler,
		msgChan:    make(chan []byte),
		property:   make(map[string]interface{}),
	}

	//将conn加入到ConnManager中
	c.TCPServer.GetConnMgr().Add(c)

	return c
}

//连接读业务
func (c *Connection) StartReader() {
	fmt.Println(" [Reader] goroutine is running")
	defer fmt.Println("ConnID ", c.ConnID, " [Read Writer] is exit")
	defer c.Stop()

	for {
		//读取客户端数据
		dp := NewDataPack()

		//读取客户端Head部分,8个字节的二进制流
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read Head error ", err)
			break
		}

		//将得到的Head部分进行拆包,得到msgID和msgDataLen
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack headData error ", err)
			break
		}

		//根据dataLen再次读取data,将以上数据组织成msg结构
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msgData error ", err)
				break
			}
		}
		msg.SetData(data)

		//当前conn的Request对象
		req := Request{
			conn: c,
			msg:  msg,
		}

		//已经开启工作池,直接交给工作worker
		if utils.GlobalObject.WorkerPoolSize > 0 {
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			go c.MsgHandler.DoMsgHandler(&req)
		}
	}
}

//连接写业务
func (c *Connection) StartWriter() {
	fmt.Println("[Writer] goroutine is running")
	defer fmt.Println(c.RemoteAddr().String(), " [conn Writer] exit")

	//阻塞等待channel消息
	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("send data err ", err)
				return
			}
		case <-c.ExitChan:
			return
		}
	}
}

func (c *Connection) Start() {
	fmt.Println("[Start] Conn start... ConnID ", c.ConnID)
	//启动当前连接的读写业务
	go c.StartReader()
	go c.StartWriter()

	//创建连接之后要添加的hook业务
	c.TCPServer.CallOnConnStart(c)
}

func (c *Connection) Stop() {
	fmt.Println("[Stop] Conn ... ConnID ", c.ConnID)

	if c.isClosed == true {
		return
	}

	c.isClosed = true

	//销毁连接之前添加hook业务
	c.TCPServer.CallOnConnStop(c)

	//关闭socket连接
	c.Conn.Close()

	//告知Writer退出
	c.ExitChan <- true

	//将当前连接删除
	c.TCPServer.GetConnMgr().Remove(c)

	//回收资源
	close(c.ExitChan)
	close(c.msgChan)

}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()

}

func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("connection already closed")
	}

	//将data进行封包处理
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMessage(msgId, data))
	if err != nil {
		fmt.Println("pack msg error ", err)
		return errors.New("pack msg error")
	}

	//将数据发送到管道
	c.msgChan <- binaryMsg

	return nil
}

//设置连接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = value
}

//获取连接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	}

	return nil, errors.New("no property not found")
}

//移除连接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}
