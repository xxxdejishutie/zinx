package znet

import (
	"fmt"
	"my/zinx/utils"
	"my/zinx/ziface"
)

type MsgHandle struct {
	//msgip对应的业务处理api
	Apis map[uint32]ziface.IRouter

	//工作池缓冲区
	WorkPoolQueue []chan ziface.IRequest

	//工作池gotine数量
	WorkPoolSize uint32
}

func NewMsgHandle() ziface.IMsgHandle {
	return &MsgHandle{
		Apis:          make(map[uint32]ziface.IRouter),
		WorkPoolQueue: make([]chan ziface.IRequest, utils.GlobalObject.WorkPoolSize),
		WorkPoolSize:  utils.GlobalObject.WorkPoolSize,
	}
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

// 创建一个WorkerPool
func (mh *MsgHandle) StartWorkerPool() {
	//循环创建配置的工作池数量的Worker
	for i := 0; i < int(mh.WorkPoolSize); i++ {
		//为每个Worker对应的chan初始化
		mh.WorkPoolQueue[i] = make(chan ziface.IRequest, int(utils.GlobalObject.MaxWorkerTaskLen))

		//开启一个worker的gotinue
		go mh.startoneWorker(i, mh.WorkPoolQueue[i])

	}
}

func (mh *MsgHandle) startoneWorker(WorkerID int, TaskChan chan ziface.IRequest) {
	fmt.Println("[Worker] worker id ", WorkerID, " start ...")

	for {
		select {
		//如果又消息到来，出列的是一个客户端的request，执行当前request对应的router
		case req := <-TaskChan:
			mh.DoRouter(req)
		}
	}
}

// 添加任务到WorkerPool中
func (mh *MsgHandle) SendMsgToTaskQueue(req ziface.IRequest) {
	//将req平均分配给worker taskqueue
	//根据客户端建立的connid分配
	workerid := req.GetConnection().GetConnId() % utils.GlobalObject.WorkPoolSize
	fmt.Printf("[Req Connection Id] %d Send To [WorkerId] %d\n", req.GetConnection().GetConnId(), workerid)

	mh.WorkPoolQueue[workerid] <- req
}
