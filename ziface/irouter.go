package ziface

/*
	router中的抽象
	所有的数据都是router
*/

type IRouter interface {
	//处理数据之前需要使用的方法
	PreHandle(request IRequest)

	//处理数据中需要使用的方法
	Handle(request IRequest)

	//处理数据后要使用的方法
	PostHandle(request IRequest)
}
