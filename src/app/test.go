package app

import "app/command"

//
import "ants/gnet"
import "ants/conf"
import "ants/gutil"
import "fmt"
import "strings"

//世界删除
func Test_remove_player(uid int) {
	if tx, ok := gnet.Socket(conf.GetRouter(conf.PORT_WORLD).Addr); ok {
		var code int16 = 1
		tx.Send(gnet.NewPackArgs(command.SERVER_WORLD_KICK_PLAYER, code, int32(uid)))
		//>>>
		tx.Join(func(b []byte) {
			pack := gnet.NewPackBytes(b)
			if pack.ReadShort() == 0 {
				fmt.Println("踢出用户成功:", pack.ReadInt())
			} else {
				fmt.Println("踢出用户失败:", pack.ReadInt())
			}
		})
	}
}

//通知世界派送消息
func Test_send_all() {
	if tx, ok := gnet.Socket(conf.GetRouter(conf.PORT_WORLD).Addr); ok {
		var uid int32 = 100000
		cmd := command.CLIENT_NOTICE_CHANNEL
		var cid int32 = 10086
		var fromid int32 = uid
		var mtype int16 = 0
		message := gutil.Ltoa(gutil.GetNano())
		message += "|"
		for i := 0; i < 10; i++ {
			message += "abcde"
		}
		tx.Send(gnet.NewPackArgs(command.SERVER_WORLD_NOTICE_PLAYERS, uid, cmd, cid, fromid, mtype, message))
		tx.Join(func(b []byte) {
			println("on read")
		})
	}
}

//获取在线用户
func Test_get_online() {
	if tx, ok := gnet.Socket(conf.GetRouter(conf.PORT_WORLD).Addr); ok {
		tx.Send(gnet.NewPackArgs(command.SERVER_WORLD_GET_ONLINE_NUM))
		tx.Join(func(b []byte) {
			pack := gnet.NewPackBytes(b)
			fmt.Println("当前在线人数 socket:", pack.ReadInt())
		})
	}
}

func Test_max_login(idx int) {
	i := idx
	for Test_login_send(i, "") {
		gutil.Sleep(5)
		i++
		if i > 5000 {
			return
		}
	}
}

func Test_login_send(idx int, pwd string) bool {
	uid := int32(idx)
	if tx, ok := gnet.Socket(conf.GetRouter(conf.PORT_GATE).Addr); ok {
		go test_socket(uid, pwd, tx)
		return true
	}
	return false
}

func test_socket(uid int32, pwd string, tx gnet.Context) {
	tx.Send(gnet.NewPackArgs(command.CLIENT_LOGON, uid, pwd))
	//
	t := gutil.GetNano()
	tx.Join(func(bits []byte) {
		packet := gnet.NewPackBytes(bits)
		switch packet.Cmd() {
		case gnet.EVENT_HEARTBEAT_PINT:
			{
				tx.Send(gnet.NewPackArgs(gnet.EVENT_HEARTBEAT_PINT))
			}
		case command.CLIENT_LOGON:
			{
				//packet.Print()
				code := packet.ReadShort()
				//body := packet.ReadBytes(0)
				//info
				//				if code == 0 {
				//					name := packet.ReadString()
				//					exp := packet.ReadInt()
				//					money := packet.ReadInt64()
				//					vipexp := packet.ReadInt()
				//					viptype := packet.ReadInt()
				//					pion := packet.ReadInt()
				//					fmt.Println("登录成功, 用户数据:", name, exp, money, vipexp, viptype, pion)
				//				}
				fmt.Println("客户端登录: err=", code, ",UID=", uid, ",runtime=", gutil.NanoStr(gutil.GetNano()-t))
				//改名
				//psend := gnet.NewPackTopic(command.CLIENT_CHANGE_NAME, conf.TOPIC_HALL, "我不是谁，谁不是我")
				//tx.Send(psend)
				//
				//psend1 := gnet.NewPackTopic(command.CLIENT_JOIN_CHANNEL, conf.TOPIC_CHAT, int32(10086), "test1")
				//tx.Send(psend1)
				//str := gutil.Int64ToString(gutil.GetTimer())
				//psend2 := gnet.NewPackTopic(command.CLIENT_NOTICE_CHANNEL, conf.TOPIC_CHAT, int32(10086), int16(1), str)
				//tx.Send(psend2)
				psend3 := gnet.NewPackTopic(command.CLIENT_ENTER_ROOM, conf.TOPIC_GAME, 100)
				tx.Send(psend3)
				psend5 := gnet.NewPackTopic(command.CLIENT_TEXAS_SITDOWN, conf.TOPIC_GAME, 100, uid, int8(1), int32(1024))
				tx.Send(psend5)
				//psend4 := gnet.NewPackTopic(command.CLIENT_QUIT_CHANNEL, conf.TOPIC_CHAT, int32(10086))
				//tx.Send(psend4)
			}
		case command.CLIENT_NOTICE_CHANNEL:
			{
				cid, fromid, _, message := packet.ReadInt(), packet.ReadInt(), packet.ReadShort(), packet.ReadString()
				strs := strings.Split(message, "|")
				chat_time := gutil.Atol(strs[0])
				fmt.Println(uid, ">收到[", fromid, "]在频道[", cid, "] 消息延迟:", gutil.NanoStr(gutil.GetNano()-chat_time), ", size=", len(strs[1]))
			}
		case command.CLIENT_JOIN_CHANNEL:
			{
				cid, fuid := packet.ReadInt(), packet.ReadInt()
				fmt.Println(fuid, "进入房间:", cid)
			}
		default:
			fmt.Println("客户端未处理:", packet.Cmd())
		}
	})
	fmt.Println("关闭了:", uid)
}

func Test(b int) {
	fmt.Println("200勇士1秒入侵服务器:", b)
	gutil.Sleep(1000)
	for i := b * 200; i < b*200+200; i++ {
		go Test_login_send(i, "abc123")
	}
}
