package ziface

//连接管理模块抽象层

type IConnManager interface {
	//添加连接
	Add(conn IConnection)
	//删除连接
	Remote(conn IConnection)
	//根据connID获取连接
	Get(connID uint32) (IConnection, error)
	//得到当前连接的总条数
	Len() int
	//清除并终止当前所有连接
	ClearConn()
}
