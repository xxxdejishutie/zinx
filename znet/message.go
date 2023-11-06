package znet

import "my/zinx/ziface"

type Message struct {

	//消息ID
	MsgId uint32
	//消息长度
	DataLen uint32
	//消息内容
	Data []byte
}

func NewMessage(msgid uint32, data []byte) ziface.IMessage {
	return &Message{
		MsgId:   msgid,
		DataLen: uint32(len(data)),
		Data:    data,
	}
}

// 设置消息id
func (m *Message) SetMsgId(msgid uint32) {
	m.MsgId = msgid
}

// 设置消息长度
func (m *Message) SetMsgLen(msglen uint32) {
	m.DataLen = msglen
}

// 设置消息数据
func (m *Message) SetMsgData(data []byte) {
	m.Data = data
}

// 获取消息id
func (m *Message) GetMsgId() uint32 {
	return m.MsgId
}

// 获取消息长度
func (m *Message) GetMsgLen() uint32 {
	return m.DataLen
}

// 获取消息数据
func (m *Message) GetMsgData() []byte {
	return m.Data
}
