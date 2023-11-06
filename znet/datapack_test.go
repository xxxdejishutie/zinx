package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
	"time"
)

// 只是测试datapack拆包
func TestDataPack(t *testing.T) {
	//模拟服务器段

	//创建监听套接字
	listener, err := net.Listen("tcp", "0.0.0.0:8999")
	if err != nil {
		fmt.Println("create listen errno ", err)
	}
	go func() {
		for {
			//获得链接
			cnn, err := listener.Accept()
			if err != nil {
				fmt.Println("listen accept errno")
				//break
			}
			go func(conn net.Conn) {
				//创建datapack实例
				dp := NewDataPack()
				for {
					//从链接中读取
					//第一次读取，把head内容读出来
					headbyte := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headbyte)
					if err != nil {
						fmt.Println("read error", err)
						break
					}

					msghead, err := dp.UnPack(headbyte)
					if err != nil {
						fmt.Println("datapack unpack error", err)
						return
					}
					//第二次读取，按照包中的datalen把数据读出来
					if msghead.GetMsgLen() > 0 {
						msg := msghead.(*Message)
						msg.Data = make([]byte, msg.DataLen)

						//根据msgdatalen读取数据
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("read data error ", err)
						}

						//完整的一个消息读完了
						fmt.Printf("--->Recv id is %d,Recv datalen is %d,Recv data is %s\n", msg.MsgId, msg.DataLen, msghead.GetMsgData())
					}
				}
			}(cnn)
		}
	}()
	//模拟客户端的代码
	clientcnn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client dial error", err)
	}
	for {
		//创建datapack实例用于封包
		dp := NewDataPack()

		//创建message1 包，并通过模块打包
		msg1 := Message{
			MsgId:   1,
			DataLen: 4,
			Data:    []byte{'z', 'i', 'n', 'x'},
		}

		sendData1, err := dp.Pack(&msg1)
		if err != nil {
			fmt.Println("datapack is error")
			return
		}

		//创建message2 包，并通过模块打包
		msg2 := Message{
			MsgId:   2,
			DataLen: 10,
			Data:    []byte{'h', 'e', 'l', 'l', 'o', '!', 's', 's', 'd', 's'},
		}

		sendData2, err := dp.Pack(&msg2)

		if err != nil {
			fmt.Println("datapack is error")
			return
		}

		//把两个包黏在一起，一起发送

		sendData1 = append(sendData1, sendData2...)

		clientcnn.Write(sendData1)
		time.Sleep(1 * time.Second)
	}
	select {}
}
