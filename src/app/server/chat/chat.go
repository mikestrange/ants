package chat

import "fat/gutil"
import "fat/gnet/nsc"
import "fat/gnet"
import "app/config"
import "app/command"
import "fat/gsys"

var router nsc.IRemoteScheduler
var mode gutil.IModeAccessor

func init() {
	//基础设施(单线程处理)
	router, mode = command.SetRouter(config.CHAT_PORT, events, gsys.NewChannel())
}

//服务器的启动
func ServerLaunch(port int) {
	gnet.ListenAndRunServer(port, func(conn interface{}) {
		gnet.LoopConnWithPing(gnet.NewConn(conn), func(tx gnet.INetContext, data interface{}) {
			mode.Done(data.(gnet.ISocketPacket).Cmd(), tx, data)
		})
	})
}
