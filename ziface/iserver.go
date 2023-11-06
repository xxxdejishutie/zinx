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
}
