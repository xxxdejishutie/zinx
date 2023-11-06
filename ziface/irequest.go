package ziface

type IRequest interface {
	//获取链接
	GetConnection() IConnection

	GetData() []byte

	//获取消息的magid
	GetMsgId() uint32
}
