package znet

import (
	"errors"
	"fmt"
	"io"
	"my/zinx/utils"
	"my/zinx/ziface"
	"net"
	"sync"
)

type Connection struct {
	//关联的TCPServer
	TCPserver ziface.ISever

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

	//链接配置
	property map[string]interface{}
	//保护链接配置的读写锁
	propertyLock sync.RWMutex
}

func NewConnection(tcpserver ziface.ISever, conn *net.TCPConn, conid uint32, msghandle ziface.IMsgHandle) *Connection {
	c := &Connection{
		TCPserver:  tcpserver,
		Conn:       conn,
		ConnId:     conid,
		IsClosed:   false,
		Msghandler: msghandle,
		msgChan:    make(chan []byte),
		ExitChan:   make(chan bool, 1),
		property:   make(map[string]interface{}),
	}

	//将连接加入到连接管理
	c.TCPserver.GetConnMgr().Add(c)

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

		//工作池开启了就用工作池，否则直接处理
		if utils.GlobalObject.WorkPoolSize > 0 {
			c.Msghandler.SendMsgToTaskQueue(&req)
		} else {
			//调用设置好的router方法
			go c.Msghandler.DoRouter(&req)
		}

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

	//调用OnConnStart函数
	c.TCPserver.CallOnConnStart(c)
}

// 结束链接
func (c *Connection) Stop() {
	fmt.Println("this connection stop...")

	//判断isclosed的值
	if c.IsClosed {
		return
	}

	c.TCPserver.CallOnConnStop(c)

	//在连接管理中删除连接
	c.TCPserver.GetConnMgr().Remote(c)

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

// 设置链接的属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = value
}

// 获取链接的属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("GetProperty error")
	}
}

// 移除链接的属性
func (c *Connection) RemoteProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}
