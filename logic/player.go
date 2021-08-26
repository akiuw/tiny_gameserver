package logic

import (
	"gameserver/pb/frames"
	"sync"
)

var g_playerMutex *sync.RWMutex
var PlayerMap map[int]*Player
var g_maxPid = 0

type Player struct {
	Pid    int
	Rid    int
	Frames map[int32]*frames.Frame
}

func init() {
	PlayerMap = make(map[int]*Player)
	g_playerMutex = &sync.RWMutex{}
	g_roomMutex = &sync.Mutex{}
}

func FindPlayer(id int) *Player {
	g_playerMutex.RLock()
	defer g_playerMutex.RLocker().Unlock()
	p, ok := PlayerMap[id]
	if ok {
		return p
	}
	return nil
}

func NewPlayer() (p *Player) {
	g_maxPid++
	p = &Player{
		Pid:    g_maxPid,
		Frames: make(map[int32]*frames.Frame),
	}
	g_playerMutex.Lock()
	defer g_playerMutex.Unlock()
	PlayerMap[p.Pid] = p
	return
}

func (p *Player) GetPlayerID() (pid int) {
	return p.Pid

}

func (p *Player) GetRommID() (rid int) {
	return p.Rid

}
