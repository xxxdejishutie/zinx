package znet

import (
	"fmt"
	"my/zinx/ziface"
)

type MsgHandle struct {
	Apis map[uint32]ziface.IRouter
}

func NewMsgHandle() ziface.IMsgHandle {
	return &MsgHandle{}
}

// 调用id对应的router方法
func (m *MsgHandle) DoRouter(req ziface.IRequest) {
	hand, ok := m.Apis[req.GetMsgId()]
	if !ok {
		panic(fmt.Sprintf("msgid[%d] Router is not rigster", req.GetMsgId()))
	}

	hand.PreHandle(req)
	hand.Handle(req)
	hand.PostHandle(req)
}

// 添加路由的方法
func (m *MsgHandle) AddRouter(msgid uint32, router ziface.IRouter) {
	_, ok := m.Apis[msgid]
	if ok {
		fmt.Println("this msgid is exist")
		return
	}
	m.Apis[msgid] = router
	fmt.Println("Msgid is", msgid, " Router Register Succ")
}
