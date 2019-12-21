package ziface

//连接管理模块抽象层
type IConnManager interface {
	//添加连接
	Add(conn IConnection)
	//删除连接
	Remove(conn IConnection)
	//根据ConnId获取具体连接
	Get(connId uint32) (IConnection, error)
	//得到当前连接数
	Len() int
	//清除连接
	ClearConn()
}
