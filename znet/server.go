package znet

import (
	"fmt"
	"my/zinx/utils"
	"my/zinx/ziface"
	"net"
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

			dealconn := NewConnection(cnn, cid, s.MsgHandler)
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

}

func (s *Server) AddRouter(msgid uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgid, router)
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
	}
	return s
}
