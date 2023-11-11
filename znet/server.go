package znet

import (
	"fmt"
	"net"

	"zinx.mod/utils"
	"zinx.mod/ziface"
)

type Server struct {
	//服务器的名字
	Name string
	//ip版本
	IPVersion string
	//ip号
	IP string
	//端口号
	Port int

	//增加的路由算法
	//Router ziface.IRouter

	//增加处理业务api
	MsgHandler ziface.IMsgHandle

	//连接管理器
	ConnMgr ziface.IConnManager

	//OnConnStart的钩子函数，在创建连接的时候调用
	OnConnStart func(ziface.IConnection)
	//OnConnStop的钩子函数，在连接结束的时候调用
	OnConnStop func(ziface.IConnection)
}

// 开启服务器
func (s *Server) Start() {
	//获取一个addr

	fmt.Printf("[zinx]ServerName is %s , TcpHost is %s ,TcpPort is %d , Version is %s\n", utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort, utils.GlobalObject.Version)
	fmt.Printf("[zinx]ServerMaxserver is %d,MaxPackageSize is %d\n", utils.GlobalObject.Maxconn, utils.GlobalObject.MaxPackageSize)
	fmt.Printf("[zinx]Workpoolsize is %d\n", utils.GlobalObject.WorkPoolSize)
	fmt.Println("[start] server has start ,  listen at  ", s.IP, " : ", s.Port)

	go func() {

		//开启工作池
		s.MsgHandler.StartWorkerPool()

		addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("net addr create errno")
		}
		//监听的listen

		listener, err := net.ListenTCP(s.IPVersion, addr)

		if err != nil {
			fmt.Printf("ipversion 【%s】 is errno ", s.IPVersion)
		}

		fmt.Println("server start succ , start listening ...")

		var cid uint32 = 0
		//阻塞监听是否有连接，并处理任务
		for {
			cnn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("accept errno", err)
				continue
			}
			fmt.Println("connect succ")

			//判断现在连接是不是超过了最大连接
			if s.ConnMgr.Len() >= utils.GlobalObject.Maxconn {
				//TODO 去给客户端回应一个错误
				fmt.Println("Too Many Connection")
				cnn.Close()
				continue
			}

			//创建connection连接
			dealconn := NewConnection(s, cnn, cid, s.MsgHandler)
			dealconn.Start()
			cid++
		}
	}()
}

// 运行服务器
func (s *Server) Serve() {

	s.Start()

	///做一些额外的功能

	select {}
}

// 关闭服务器
func (s *Server) Stop() {
	//释放资源

	fmt.Println("[STOP] Server ClearConn")
	s.ConnMgr.ClearConn()
}

func (s *Server) AddRouter(msgid uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgid, router)
}

func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

// 初始化server模块
func NewServer(name string) ziface.ISever {
	s := &Server{
		Name:      utils.GlobalObject.Name,
		IPVersion: "tcp",
		IP:        utils.GlobalObject.Host,
		Port:      utils.GlobalObject.TcpPort,
		//Router:    nil,
		MsgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(),
	}
	return s
}

// 设置OnConnStart的钩子函数，在连接结束的时候调用
func (s *Server) SetOnConnStart(UserOnConnStart func(ziface.IConnection)) {
	s.OnConnStart = UserOnConnStart
}

// 设置OnConnStop的钩子函数，在连接结束的时候调用
func (s *Server) SetOnConnStop(UserOnConnStop func(ziface.IConnection)) {
	s.OnConnStop = UserOnConnStop
}

// 调用OnConnStart函数
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	fmt.Println("Call -----> OnConnStart")
	s.OnConnStart(conn)
}

// 调用OnConnStop函数
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	fmt.Println("Call -----> OnConnStop")
	s.OnConnStop(conn)
}
