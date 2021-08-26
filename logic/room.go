package logic

import (
	"errors"
	"fmt"
	"gameserver/network"
	"gameserver/pb"
	"gameserver/pb/frames"
	"math/rand"
	"sync"
	"time"
	"unsafe"
)

var RoomMap *sync.Map
var g_roomMutex *sync.Mutex

type Room struct {
	*GameLoop
	Rid           int
	Players       []int
	StartTime     int64
	EndTime       int64
	RoomMaxPlayer int
	FrameID       int32
}

func NewRoom() *Room {
	r := &Room{
		Rid:           rand.Intn(100),
		StartTime:     time.Now().Unix(),
		EndTime:       0,
		RoomMaxPlayer: 2,
		Players:       make([]int, 0),
		FrameID:       0,
	}
	RoomMap.Store(r.Rid, r)
	gs := network.GetInstance()
	gs.AddLogic(r)
	go r.RunFrame()
	return r
}

func (r *Room) AddPlayer(pid int) (err error) {
	for v := range r.Players {
		if v == pid {
			return errors.New(fmt.Sprintf("player allready in this room : id = %d", r.Rid))
		}
	}
	r.Players = append(r.Players, pid)
	return nil
}

func (r *Room) NotifyAllPlayers(f *frames.Frame) {
	for v := range r.Players {
		p := FindPlayer(v)
		if p == nil {

			fmt.Printf("no player id = %d\n", v)
		}
		p.Frames[f.Fid] = f
	}

}

//tick函数会每帧唤醒每个gameobject的update函数
func (r *Room) RunFrame() {
	r.ch <- 0
	for {
		<-r.ch
		r.Update()
	}
}

func FindRoom(rid int) *Room {

	r, ok := RoomMap.Load(rid)
	if ok {
		return r.(*Room)
	}
	return nil

}

//每帧更新的逻辑
func (r *Room) Update() {
	//通知房间里面所有的玩家更新
	for _, v := range r.Players {
		p := FindPlayer(v)
		ConnLocker.RLocker().Lock()
		defer ConnLocker.RLocker().Unlock()
		f := p.Frames[r.FrameID]
		c, ok := Connections[p.Pid]
		if ok {
			pb.SendToClinet(pb.FRAMES, f, int32(unsafe.Sizeof(*f)), c)
		} else {
			//如果玩家操作里面不存在该帧号,广播空操作给整个房间
			f.Dir = frames.DIRCTION_DOWN
			f.Opt = frames.OPERATOR_NONE
			f.Pid = int32(p.Pid)
			pb.SendToClinet(pb.FRAMES, f, int32(unsafe.Sizeof(*f)), c)
		}
	}
	r.FrameID++
}
