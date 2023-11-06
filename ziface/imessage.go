package ziface

type IMessage interface {
	//设置消息id
	SetMsgId(uint32)

	//设置消息长度
	SetMsgLen(uint32)

	//设置消息数据
	SetMsgData([]byte)

	//获取消息id
	GetMsgId() uint32

	//获取消息长度
	GetMsgLen() uint32

	//获取消息数据
	GetMsgData() []byte
}
