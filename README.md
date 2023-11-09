# zinx
## 一.使用流程
### 1.先创建server对象
```go
    ser := znet.NewServer(YourServerName)
```
    
### 2.为业务设置处理方法
    定义一个自己的处理方法类，继承znet.base，按照需要重写prehandle ,handle ,posthandle 管理函数，参数为ziface.irequest

```go
    type Myrouter struct {
	znet.BaseRouter
    }

    func (this *Myrouter) Handle(req ziface.IRequest) {
	    fmt.Println("call Handle..")

	    fmt.Println("msgid is ", req.GetMsgId(), "msgdata is ", string(req.GetData()))

	    err := req.GetConnection().PackSend(1, []byte("ping...ping...ping"))
	    if err != nil {
		    fmt.Println(err)
	    }
    }
```
### 3.将自定义处理方法，与对应的任务id注册到server对象中
    添加路由，多个不同处理方法需要保证业务号不同

```go
	ser.AddRouter(0, &Myrouter{})
    ser.AddRouter(1, &HelloRouter{})
```

### 4.启动server

```go
    ser.start()
```
    
### 5.总体代码

```go
    import (
	"fmt"
	"my/zinx/ziface"
	"my/zinx/znet"
    )   

    type Myrouter struct {
	znet.BaseRouter
    }

    func (this *Myrouter) Handle(req ziface.IRequest) {
	    fmt.Println("call Handle..")

	    fmt.Println("msgid is ", req.GetMsgId(), "msgdata is ", string(req.GetData()))

	    err := req.GetConnection().PackSend(1, []byte("ping...ping...ping"))
	    if err != nil {
		    fmt.Println(err)
	    }
    }
    func main() {
	//创建zinx服务器,使用zinx的api
	ser := znet.NewServer("[zinx 0.2]")

	//添加路由
	ser.AddRouter(0, &Myrouter{})
	ser.AddRouter(1, &HelloRouter{})

	//运行zinx服务器
	ser.Serve()
    }
```

## 二.客户端通信格式
### 1.包的定义
    HLV 格式，包的头部4个字节为包的大小，之后的四个字节为消息对应ID，之后为消息内容
### 2.使用方法
	直接发送对应包格式的包即可
### 3.注意事项
	degers