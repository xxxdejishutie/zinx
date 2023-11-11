package ziface

import "net"

type IConnection interface {
	//开始链接
	Start()
	//结束链接
	Stop()

	//返回链接的套接字
	GetTcpConn() *net.TCPConn

	//返回链接的ID
	GetConnId() uint32

	//返回对方的addr信息
	GetRemoteAddr() net.Addr

	//发送数据
	PackSend(uint32, []byte) error

	//设置链接的属性
	SetProperty(string, interface{})
	//获取链接的属性
	GetProperty(string) (interface{}, error)
	//移除链接的属性
	RemoteProperty(string)
}

type HandleFun func(*net.TCPConn, []byte, int) error
