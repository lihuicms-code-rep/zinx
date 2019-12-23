package znet

import (
	"errors"
	"fmt"
	"sync"
	"github.com/lihuicms-code-rep/zinx/ziface"
)

type ConnManager struct {
	connections map[uint32]ziface.IConnection    //所有连接信息
	connLock    sync.RWMutex                     //读写锁
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections:make(map[uint32]ziface.IConnection),
	}
}

//添加连接
func (connMgr *ConnManager) Add(conn ziface.IConnection) {
	//加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//将conn加入
	connMgr.connections[conn.GetConnID()] = conn

	fmt.Println("connection add to ConnManager successfully....", conn.GetConnID())
}

//删除连接
func (connMgr *ConnManager) Remove(conn ziface.IConnection) {
	//加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//将conn删除
	delete(connMgr.connections, conn.GetConnID())

	fmt.Println("connection remove from ConnManager successfully....", conn.GetConnID())


}

//根据ConnId获取具体连接
func (connMgr *ConnManager) Get(connId uint32) (ziface.IConnection, error) {
	//加读锁
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.connections[connId]; ok {
		return conn, nil
	}

	return nil, errors.New("connection not found")
}

//得到当前连接数
func (connMgr *ConnManager) Len() int {
	return len(connMgr.connections)
}

//清除连接
func (connMgr *ConnManager) ClearConn() {
	//加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//conn停止工作并删除
	for connID, conn := range connMgr.connections {
		conn.Stop()
		delete(connMgr.connections, connID)
	}

}
