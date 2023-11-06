package znet

import "my/zinx/ziface"

type BaseRouter struct{}

//加入一个基类，可以让router类在实现的过程中可以不全部的三个方法都实现

// 处理数据之前需要使用的方法
func (r *BaseRouter) PreHandle(ziface.IRequest) {}

// 处理数据中需要使用的方法
func (r *BaseRouter) Handle(ziface.IRequest) {}

// 处理数据后要使用的方法
func (r *BaseRouter) PostHandle(ziface.IRequest) {}
