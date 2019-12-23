package znet

import "github.com/lihuicms-code-rep/zinx/ziface"

//实现router时,先嵌入BaseRouter基类, 然后对这个基类的方法进行实现
//BaseRouter已经把IRouter已经实现了
//只要继承BaseRouter重写自己想实现的方法即可
//这里设计思路参考于:beego
type BaseRouter struct {}

//处理业务之前的钩子
func (br *BaseRouter) PreHandle(request ziface.IRequest) {

}

//处理业务的方法
func (br *BaseRouter) Handle(request ziface.IRequest) {

}

//处理业务之后的钩子
func (br *BaseRouter) PostHandle(request ziface.IRequest) {

}


