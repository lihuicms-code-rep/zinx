package znet

import (
	"fmt"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

//iServer接口类的实现
type Server struct {
	Name        string                   //服务名称
	IPVersion   string                   //tcp4或其他
	IP          string                   //绑定IP地址
	Port        int                      //绑定端口
	MsgHandler  ziface.IMsgHandler       //服务消息管理模块
	ConnMgr     ziface.IConnManager      //服务连接管理模块
	OnConnStart func(ziface.IConnection) //连接建立时要处理的hook业务,这里设计成两个函数类型的属性
	OnConnStop  func(ziface.IConnection) //连接断开之前要处理的hook业务
}

//初始化Sever模块
func NewServer(name string) ziface.IServer {
	return &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  utils.GlobalObject.IPVersion,
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.Port,
		MsgHandler: NewMsgHandler(),
		ConnMgr:    NewConnManager(),
	}
}

func (s *Server) Start() {
	fmt.Printf("[START] Server listening at IP:%s, Port:%d\n", s.IP, s.Port)
	go s.serverLogic()

}

func (s *Server) Stop() {
	//将服务器的资源,连接等回收
	fmt.Println("[STOP] zinx Server, name ", s.Name)
	s.ConnMgr.ClearConn()
}

func (s *Server) Serve() {
	s.Start()

	//阻塞住
	select {}
}

//获取服务连接管理
func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}


//对外暴露
//注册OnConnStart 函数
func (s *Server) SetOnConnStart(hookFunc func(connection ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

//调用
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("CallOnConnStart......")
		s.OnConnStart(conn)
	}
}

//注册OnConnStop
func (s *Server) SetOnConnStop(hookFunc func(connection ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

//调用
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("CallOnConnStop.......")
		s.OnConnStop(conn)
	}
}

//服务器逻辑,3步
func (s *Server) serverLogic() {
	//1.获取TCP addr(socket句柄)
	addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		fmt.Println("resolve tcp addr err ", err)
		return
	}

	//2.监听服务器地址
	listener, err := net.ListenTCP(s.IPVersion, addr)
	if err != nil {
		fmt.Println("listen", s.IP, " err", err)
		return
	}

	fmt.Println("[Success] start zinx server  ", s.Name, " success, now listening version ", utils.GlobalObject.Version)
	var cid uint32

	//3.阻塞的等待客户端进行连接,处理连接业务(读写)
	for {
		//3.0 开启消息队列和worker工作池
		s.MsgHandler.StartWorkerPool()

		//3.1 如果客户端连接事件过来,阻塞会返回
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("listener accept err", err)
			continue
		}

		//3.2 判断是否超过最大连接数
		if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
			fmt.Println("too many connection, max conn =", utils.GlobalObject.MaxConn)
			conn.Close()
			continue
		}

		//3.3 将处理新连接的业务方法和conn绑定
		dealConn := NewConnection(s, conn, cid, s.MsgHandler)
		cid++
		dealConn.Start()
	}
}

//连接建立后的具体业务
func (s *Server) connHandler(conn *net.TCPConn) {
	if conn == nil {
		return
	}

	//阻塞等待客户端数据
	for {
		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("receive  buf err ", err)
			continue
		}

		fmt.Printf("receive client data:%s, cnt=%d\n", buf, cnt)

		//读到就进行简单的回写
		if _, err := conn.Write(buf[:cnt]); err != nil {
			fmt.Println("write back buf err ", err)
			continue
		}
	}
}

func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgId, router)
	fmt.Println("server add router success ...")
}
