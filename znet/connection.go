package znet

import (
	"fmt"
	"io"
	"my/zinx/ziface"
	"net"
)

type Connection struct {
	//链接的套接字
	Conn *net.TCPConn

	//链接的id
	ConnId uint32

	//链接是否关闭
	IsClosed bool

	//与当前链接所绑定的处理方法
	//HandleApi ziface.HandleFun

	//绑定的router处理算法
	//Router ziface.IRouter

	//消息管理 msgid 和对应的处理业务的api
	Msghandler ziface.IMsgHandle

	//等待链接退出的channel
	ExitChan chan bool

	//用于读端和写端通信的管道
	msgChan chan []byte
}

func NewConnection(conn *net.TCPConn, conid uint32, msghandle ziface.IMsgHandle) *Connection {
	c := &Connection{
		Conn:       conn,
		ConnId:     conid,
		IsClosed:   false,
		Msghandler: msghandle,
		msgChan:    make(chan []byte),
		ExitChan:   make(chan bool, 1),
	}
	return c
}

func (c *Connection) StartReader() {
	fmt.Println("[Connection Reader Start]")

	//fmt.Println("Connid is ", c.ConnId)
	defer fmt.Println("[Reader exit]Connid is ", c.ConnId, ", Remote addr is ", c.Conn.RemoteAddr().String())
	defer c.Stop()
	//创建datapack实例用于拆包

	for {
		dp := NewDataPack()
		//第一次读取包头部，给msg的msgid和msglen
		headbyte := make([]byte, dp.GetHeadLen())
		cnt, err := io.ReadFull(c.GetTcpConn(), headbyte)
		if err != nil {
			fmt.Println("read pack head error ", err)
			if cnt == 0 {
				return
			}
			continue
		}

		msg, err := dp.UnPack(headbyte)
		if err != nil {
			fmt.Println("unpack errno ", err)
		}

		//fmt.Println("message read first msglen is ", msg.GetMsgLen())
		//第二次读取包内容，给msg的data
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			_, err := io.ReadFull(c.GetTcpConn(), data)
			if err != nil {
				fmt.Println("read pack data error ", err)
				return
			}
		}
		msg.SetMsgData(data)

		//将消息封装成request包
		req := Request{
			Conn: c,
			msg:  msg,
		}

		//调用设置好的router方法
		go c.Msghandler.DoRouter(&req)

	}
}

// 写端
func (c *Connection) StartWriter() {
	fmt.Println("[Connection Writer Start]")
	defer fmt.Println("[Connection Writer exit]")
	for {
		select {
		case data := <-c.msgChan:
			_, err := c.Conn.Write(data)
			if err != nil {
				fmt.Println("[Connection Writer]  write error ", err)
				return
			}
		case <-c.ExitChan:
			return
		}
	}
}

// 开始链接
func (c *Connection) Start() {
	fmt.Println("[Conn start] connection id is ", c.ConnId)

	//TODO 不断读取数据
	go c.StartReader()

	//读端开始运行
	go c.StartWriter()
}

// 结束链接
func (c *Connection) Stop() {
	fmt.Println("this connection stop...")

	//判断isclosed的值
	if c.IsClosed {
		return
	}

	//连接关闭
	c.Conn.Close()
	c.ExitChan <- true

	close(c.ExitChan)
	close(c.msgChan)
}

// 返回链接的套接字
func (c *Connection) GetTcpConn() *net.TCPConn {
	return c.Conn
}

// 返回链接的ID
func (c *Connection) GetConnId() uint32 {
	return c.ConnId
}

// 返回对方的addr信息
func (c *Connection) GetRemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// 发送数据
func (c *Connection) PackSend(msgid uint32, data []byte) error {

	//创建封包的datapack
	dp := NewDataPack()
	msg := NewMessage(msgid, data)

	fmt.Println("[send] msg len is ", msg.GetMsgLen())
	binaryMsg, err := dp.Pack(msg)
	if err != nil {
		fmt.Println("datapack pack error ", err)
		return err
	}

	//向管道写入数据，发送给写端
	c.msgChan <- binaryMsg
	return nil
}
