package cluster

import (
	"ants/actor"
	"ants/gnet"
	"ants/gsys"
)

//分布式基础
type INodeRouter interface {
	actor.IBoxRef
	//路由数据
	Data() IDataRoute
}

//远程调度器
type NodeRouter struct {
	actor.BaseBox
	gsys.Locked
	data IDataRoute
	conn gnet.Context
}

func NewRouter(data IDataRoute) INodeRouter {
	this := new(NodeRouter)
	this.SetActor(this) //自身行动
	this.data = data
	return this
}

func NewRouterPort(port int) INodeRouter {
	return NewRouter(NewData(port))
}

//interface INodeRouter(处理方式就在这里)
func (this *NodeRouter) OnMessage(args ...interface{}) {
	conn, ok := this.Context()
	if ok {
		conn.Send(args...)
	}
}

func (this *NodeRouter) UnContext() {
	this.Lock()
	this.conn = nil
	this.Unlock()
}

func (this *NodeRouter) Context() (gnet.Context, bool) {
	this.Lock()
	defer this.Unlock()
	if this.conn == nil {
		if conn, ok := gnet.Socket(this.Data().Addr()); ok {
			this.conn = conn
			gnet.RunWithAgent(this)
		} else {
			return nil, false
		}
	}
	return this.conn, true
}

func (this *NodeRouter) Data() IDataRoute {
	return this.data
}

func (this *NodeRouter) OnHandle(b []byte) {
	//处理回执(基本不用处理)
	pack := gnet.NewPackBytes(b)
	switch pack.Cmd() {
	case gnet.EVENT_HEARTBEAT_PINT:
		this.conn.Send(gnet.NewPackArgs(gnet.EVENT_HEARTBEAT_PINT))
	default:
		println("node not handle:", pack.Cmd())
	}
}

func (this *NodeRouter) SetHandle(f func([]byte)) {

}

//====
func (this *NodeRouter) Run() {
	this.conn.Join(this.OnHandle)
	this.UnContext()
}

func (this *NodeRouter) OnClose() {
	//关闭
}
