package logic

import (
	"gameserver/pb/frames"
	"sync"

	"github.com/panjf2000/gnet"
)

var frameCount = 0
var zsetkey = "serverframe"
var FramesMap map[int]map[int]*frames.Frame

//索引为pid [pid]conn
var Connections map[int]gnet.Conn
var ConnLocker *sync.RWMutex

type GameLoop struct {
	ch chan int
}

type Connection struct {
	Pid int
	c   gnet.Conn
}

func init() {
	Connections = make(map[int]gnet.Conn)
	FramesMap = make(map[int]map[int]*frames.Frame)
	RoomMap = new(sync.Map)
	ConnLocker = new(sync.RWMutex)
}

func (gl *GameLoop) RunFrame() {}
