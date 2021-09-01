package network

import (
	"errors"
	"fmt"
	"gameserver/ihandles"
	"gameserver/ilogic"
	"gameserver/pb"
	"gameserver/pb/actions"
	"gameserver/pb/frames"
	"time"
	"unsafe"

	"github.com/panjf2000/gnet"
	"github.com/panjf2000/gnet/pool/goroutine"
	"google.golang.org/protobuf/proto"
)

var recvbuffer []byte

var g_instance *GameServer

type GameServer struct {
	*gnet.EventServer
	tick time.Duration
	pool *goroutine.Pool
	l    []ilogic.IGameLoop
	h    ihandles.IHandles
}

type PoolFunction struct {
	data interface{}
	conn gnet.Conn
	f    func(d interface{}, c gnet.Conn)
}

func init() {

	recvbuffer = make([]byte, 0)
}

func GetInstance() *GameServer {
	if g_instance == nil {
		g_instance = &GameServer{}
	}
	return g_instance
}

func (pf *PoolFunction) function() {
	pf.f(pf.data, pf.conn)
}

func (gs *GameServer) RegisterProtocol(h ihandles.IHandles) {

	pb.RegistProtocol(pb.ACTION, h.GetHandle(pb.ACTION))
	pb.RegistProtocol(pb.FRAMES, h.GetHandle(pb.FRAMES))
}

// AddLogic 添加游戏object(实现update的对象)每帧都会调用
func (gs *GameServer) AddLogic(g ilogic.IGameLoop) {
	gs.l = append(gs.l, g)
}

// SetTimer 设置服务器的帧间隔时间(整数)
func (gs *GameServer) SetTimer(t time.Duration) {
	gs.tick = t
}

func (gs *GameServer) PoolInit(p *goroutine.Pool) {
	gs.pool = p
}

func (gs *GameServer) Tick() (delay time.Duration, action gnet.Action) {
	//循环通知所有的gameobject
	for _, v := range gs.l {
		v.RunFrame()
	}

	delay = gs.tick
	return
}

// React 拆包发放带各个模块
func (gs *GameServer) React(data []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	fmt.Println(data)
	lenth := 0
	if len(recvbuffer) > 0 {
		data = append(recvbuffer, data...)
		recvbuffer = recvbuffer[0:0]
	}
	for lenth < len(data) {
		mh, err := GetMessageHeader(data[lenth:])
		if err != nil {
			fmt.Println(err)
			return
		}
		start := int(unsafe.Sizeof(*mh)) + lenth
		end := int(unsafe.Sizeof(*mh)) + lenth + int(mh.Lenth)

		if len(data[start:]) < int(mh.Lenth) {
			//不够的先缓存
			recvbuffer = append(recvbuffer, data[start:]...)
			return
		}
		d := MessageBodyUnmarshal(mh, data[start:end])
		lenth += int(unsafe.Sizeof(*mh)) + int(mh.Lenth)
		pf := &PoolFunction{
			data: d,
			conn: c,
			f:    pb.ProtocolMap[mh.ID],
		}
		gs.pool.Submit(pf.function)
	}
	return
}

// GetMessageHeader 得到协议头
func GetMessageHeader(data []byte) (messageheader *pb.MessageHeader, err error) {

	messageheader = (*pb.MessageHeader)(unsafe.Pointer(&data[0]))
	if messageheader.ID <= pb.MESSAGE_BEGIN || messageheader.ID >= pb.MESSAGE_END {
		return nil, errors.New(fmt.Sprintf("error: protocol message id! ID : %d", messageheader.ID))
	}
	return messageheader, nil
}

func MessageBodyUnmarshal(messageheader *pb.MessageHeader, data []byte) interface{} {

	switch messageheader.ID {
	case pb.ACTION:
		{
			msbody := &actions.C2SPlayerAction{}
			proto.Unmarshal(data, msbody)
			fmt.Println(msbody)
			return msbody
		}
	case pb.FRAMES:
		{
			msbody := &frames.Frame{}
			proto.Unmarshal(data, msbody)
			return msbody

		}
	}
	return nil
}
