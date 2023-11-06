package znet

import (
	"my/zinx/ziface"
)

type Request struct {

	//链接
	Conn ziface.IConnection

	//数据
	msg ziface.IMessage
}

func (r *Request) GetConnection() ziface.IConnection {
	return r.Conn
}

func (r *Request) GetData() []byte {
	return r.msg.GetMsgData()
}

func (r *Request) GetMsgId() uint32 {
	return r.msg.GetMsgId()
}
