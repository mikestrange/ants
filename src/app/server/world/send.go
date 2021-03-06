package world

import "app/command"
import "ants/gcode"
import "app/conf"

//踢人
func packet_kick_player(code int, player *GamePlayer) interface{} {
	return gcode.NewPackArgs(command.SERVER_KICK_PLAYER, player.UserID, player.SessionID, int16(code))
}

//登录返回
func packet_logon_result(code int, uid int, session uint64, body []byte) interface{} {
	return gcode.NewPackArgs(command.SERVER_LOGON_RESULT, int16(code), uid, session, body)
}

//直接发送给客户端
func packet_send_client(cmd int, uid int, session uint64, body []byte) interface{} {
	return gcode.NewPackTopic(cmd, conf.TOPIC_CLIENT, uid, session, body)
}
