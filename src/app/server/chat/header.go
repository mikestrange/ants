package chat

import "ants/gnet"

//会话头
type GameHeader struct {
	UserID    int
	GateID    int
	SessionID uint64
}

func NewHeader(pack gnet.IByteArray) *GameHeader {
	this := new(GameHeader)
	this.UnPack(pack)
	return this
}

func (this *GameHeader) UnPack(pack gnet.IByteArray) {
	pack.ReadValue(&this.UserID, &this.GateID, &this.SessionID)
}

func (this *GameHeader) SerID() int {
	return int(this.GateID)
}

func (this *GameHeader) UID() int {
	return int(this.UserID)
}
