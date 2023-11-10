package ziface

type ISever interface {
	//开启服务器
	Start()

	//运行服务器
	Stop()

	//关闭服务器
	Serve()

	//增加一个路由
	AddRouter(uint32, IRouter)

	//获取connmanger
	GetConnMgr() IConnManager

	// 设置OnConnStart的钩子函数，在连接结束的时候调用
	SetOnConnStart(func(IConnection))

	// 设置OnConnStop的钩子函数，在连接结束的时候调用
	SetOnConnStop(func(IConnection))
	// 调用OnConnStart函数
	CallOnConnStart(conn IConnection)

	// 调用OnConnStop函数
	CallOnConnStop(IConnection)
}
