package gnet

//private 基础包
type ISocketPacket interface {
	IByteArray
	//get
	Cmd() int
	Topic() int
	VerID() int
	Format() int
	//set
	SetCmd(cmd int)
	SetTopic(int)
	SetVersion(int)
	SetFormat(int)
	//write
	WriteBegin(cmd int) interface{}
	WriteBeginWithTopic(cmd int, topic int) interface{}
	WriteEnd() interface{}
	//body
	SetBodyBegin()
	BodySize() int
	BodyPos() int
	//read
	ReadBegin() interface{}
	ReadEnd() interface{}
	GetBody() []byte
	//
	Print()
}

type SocketPacket struct {
	ByteArray
	m_cmd    int32 //命令符
	m_topic  uint8 //一个游戏不需要那么多服务器
	m_verid  uint8 //不需要那么多版本
	m_format uint8 //解析方式
	//m_time   int64
	m_pbody int
}

func NewPacket() ISocketPacket {
	this := &SocketPacket{}
	this.InitByteArray()
	return this
}

//通知网关
func NewPackArgs(cmd int, args ...interface{}) ISocketPacket {
	this := NewPacket()
	this.WriteBegin(cmd)
	this.WriteValue(args...)
	this.WriteEnd()
	return this
}

//通知其他
func NewPackTopic(cmd int, topic int, args ...interface{}) ISocketPacket {
	this := NewPacket()
	this.WriteBeginWithTopic(cmd, topic)
	this.WriteValue(args...)
	this.WriteEnd()
	return this
}

//直接引用字段
func NewPackBytes(bits []byte) ISocketPacket {
	this := &SocketPacket{}
	this.InitByteArrayWithBits(bits)
	this.ReadBegin()
	return this
}

func (this *SocketPacket) Cmd() int {
	return int(this.m_cmd)
}

func (this *SocketPacket) Topic() int {
	return int(this.m_topic)
}

func (this *SocketPacket) VerID() int {
	return int(this.m_verid)
}

func (this *SocketPacket) Format() int {
	return int(this.m_format)
}

//set
func (this *SocketPacket) SetCmd(cmd int) {
	this.m_cmd = int32(cmd)
}

func (this *SocketPacket) SetVersion(val int) {
	this.m_verid = uint8(val)
}

func (this *SocketPacket) SetTopic(val int) {
	this.m_topic = uint8(val)
}

func (this *SocketPacket) SetFormat(val int) {
	this.m_format = uint8(val)
}

//body
func (this *SocketPacket) BodySize() int {
	return this.Length() - this.m_pbody
}

func (this *SocketPacket) BodyPos() int {
	return this.m_pbody
}

func (this *SocketPacket) SetBodyBegin() {
	this.SetPos(this.m_pbody)
}

//write
func (this *SocketPacket) WriteBegin(cmd int) interface{} {
	return this.WriteBeginWithTopic(cmd, 0)
}

func (this *SocketPacket) WriteBeginWithTopic(cmd int, topic int) interface{} {
	this.Reset() //每次都会冲掉
	this.SetCmd(cmd)
	this.SetTopic(topic)
	//this.m_time = gutil.GetNano()
	this.WriteValue(this.m_cmd, this.m_topic, this.m_verid, this.m_format)
	this.m_pbody = this.Pos()
	return this
}

func (this *SocketPacket) WriteEnd() interface{} {
	this.SetBegin()
	return this
}

//reads
func (this *SocketPacket) ReadBegin() interface{} {
	this.SetBegin()
	this.ReadValue(&this.m_cmd, &this.m_topic, &this.m_verid, &this.m_format)
	this.m_pbody = this.Pos()
	return this
}

func (this *SocketPacket) ReadEnd() interface{} {
	this.SetEnd()
	return this
}

func (this *SocketPacket) GetBody() []byte {
	this.SetPos(this.BodyPos())
	return this.ReadBytes(0)
}

func (this *SocketPacket) Print() {
	//fmt.Println(this.Cmd(), "网络消耗:", gutil.NanoStr(gutil.GetNano()-this.m_time))
}
