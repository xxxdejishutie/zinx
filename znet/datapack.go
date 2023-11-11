package znet

import (
	"bytes"
	"encoding/binary"
	"errors"

	"zinx.mod/ziface"
)

type DataPack struct{}

func NewDataPack() ziface.IDataPack {
	dp := &DataPack{}
	return dp
}

// 获取包头长度
func (d *DataPack) GetHeadLen() uint32 {

	return 8
}

// 封包
func (d *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	//创捷个byte类型的缓冲
	databyte := bytes.NewBuffer([]byte{})

	//按照大端存储，将msglen写入缓冲
	err := binary.Write(databyte, binary.LittleEndian, msg.GetMsgLen())
	if err != nil {
		return nil, err
	}

	//按照大端存储，将msgid写入缓冲
	err = binary.Write(databyte, binary.LittleEndian, msg.GetMsgId())
	if err != nil {
		return nil, err
	}

	//按照大端存储，将msgdata写入缓冲
	err = binary.Write(databyte, binary.LittleEndian, msg.GetMsgData())
	if err != nil {
		return nil, err
	}
	return databyte.Bytes(), nil
}

// 拆包 将包的head读出来，然后根据head的data长度，再进行一次读
func (d *DataPack) UnPack(buffer []byte) (ziface.IMessage, error) {
	databyte := bytes.NewBuffer(buffer)

	msg := &Message{}

	if err := binary.Read(databyte, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	if err := binary.Read(databyte, binary.LittleEndian, &msg.MsgId); err != nil {
		return nil, err
	}

	if msg.DataLen > 1000 { //uint32(utils.GlobalObject.MaxPackageSize) {
		return nil, errors.New("this package is large")
	}
	return msg, nil
}
