package ziface

type IMsgHandle interface {
	//调用id对应的router方法
	DoRouter(IRequest)

	//添加路由的方法
	AddRouter(uint32, IRouter)
}
