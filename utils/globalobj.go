package utils

import (
	"encoding/json"
	"fmt"
	"my/zinx/ziface"
	"os"
)

/*
	存储关于zinx框架的全局参数，供其他模块使用
	一些参数可以由用户通过zinx.json配置
*/

type GlobalObj struct {
	Server  ziface.ISever
	Name    string //服务名称
	TcpPort int    //绑定的端口号
	Host    string //主机的ip地址

	Version        string //zinx的版本
	Maxconn        int    //允许链接的最大数量
	MaxPackageSize int    //一次接受的文件包最大大小
}

func (g *GlobalObj) Reload() {
	//读取用户的配置文件
	data, err := os.ReadFile("conf/zinx.json")
	if err != nil {
		fmt.Println("errno is ", err)
		panic(err)
	}

	//解析成json文件,必须使用全局变量的地址
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}

}

// 定义一个全局变量，供其他包访问
var GlobalObject *GlobalObj

func init() {

	//如果配置文件没有加载，默认的值
	GlobalObject := &GlobalObj{
		Name:           "ZinxServerApp",
		TcpPort:        8999,
		Host:           "0.0.0.0",
		Version:        "zinx v0.4",
		Maxconn:        1000,
		MaxPackageSize: 512,
	}
	//可以从zinx.json中读取配置信息，加载用户定义的参数
	GlobalObject.Reload()

}
