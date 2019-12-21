package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"zinx/ziface"
)

//存储一切有关zinx的全局参数
type GlobalObj struct {
	TCPServer ziface.IServer //zinx全局Server
	IPVersion string         //网络版本
	Host      string         //主机监听IP
	Port      int            //主机监听端口
	Name      string         //服务器名称

	Version        string   //服务器版本
	MaxConn          int    //最大连接数
	MaxPackageSize   uint32 //允许最大数据包大小
	WorkerPoolSize   uint32 //工作worker数量
	MaxWorkerTaskLen uint32 //一个队列最大处理量
}

//对外访问对象
var GlobalObject *GlobalObj

//初始化GlobalObject
func init() {
	GlobalObject = &GlobalObj{
		Host:           "0.0.0.0",
		IPVersion:      "tcp4",
		Port:           7777,
		Name:           "ZinxServerApp",
		Version:        "V0.4",
		MaxConn:        1000,
		MaxPackageSize: 1024,
		WorkerPoolSize:10,
		MaxWorkerTaskLen:100,
	}

	GlobalObject.Reload()
}

//加载配置文件数据
func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		fmt.Println("read zinx.json error ", err)
		return
	}

	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}
