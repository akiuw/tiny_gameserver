package handles

import (
	"fmt"
	"gameserver/logic"
	"gameserver/pb"
	"gameserver/pb/actions"
	"gameserver/pb/frames"
	"gameserver/pb/msgerr"

	"github.com/panjf2000/gnet"
)

type Handle struct{}

var g_handles *Handle
var count = 0

func GetInstance() *Handle {
	if g_handles == nil {
		g_handles = &Handle{}
	}
	return g_handles
}

func C2SPlayerAction(data interface{}, c gnet.Conn) {
	s := data.(*actions.C2SPlayerAction)
	switch s.Action {
	case actions.ACTION_LOGIN:
		{
			if s.Pid != 0 {
				S2CSendErrorMessage(msgerr.ERROR_CODE_PLAYER_LOGIN_ERROR, c)
				break
			}
			p := logic.NewPlayer()
			s.Pid = int32(p.Pid)

			err := SendPlayerLoginAck(s, c)
			if err != nil {
				S2CSendErrorMessage(msgerr.ERROR_CODE_PLAYER_LOGIN_ERROR, c)
				break
			}
			//将连接存入Connections
			logic.ConnLocker.Lock()
			defer logic.ConnLocker.Unlock()
			logic.Connections[p.Pid] = c
			fmt.Println("player :", p.Pid, "login success")

			break
		}
	case actions.ACTION_JOINROOM:
		{
			logic.RoomMap.Range(func(k, v interface{}) bool {
				d := v.(*logic.Room)
				if len(d.Players) < 2 {
					err := SendPlayerJoinRoomAck(int(s.Pid), d.Rid, c)
					if err != nil {
						fmt.Println(err)
						return true
					}
					d.Players = append(d.Players, int(s.Pid))
				}
				return true
			})

			r := logic.NewRoom()
			if r == nil {
				S2CSendErrorMessage(msgerr.ERROR_CODE_PLATER_JOIN_ROOM_ERROR, c)
				break
			}
			logic.RoomMap.Store(r.Rid, r)

			err := SendPlayerJoinRoomAck(int(s.Pid), r.Rid, c)
			if err != nil {
				fmt.Println(err)
				S2CSendErrorMessage(msgerr.ERROR_CODE_PLATER_JOIN_ROOM_ERROR, c)
				break

			}
			r.Players = append(r.Players, int(s.Pid))
			break
		}
	case actions.ACTION_LOGOUT:
		{
			break
		}
	}

}
func FrameToBuffer(data interface{}, c gnet.Conn) {
	s := data.(*frames.Frame)
	fmt.Println("in FrameToBuffer", s)
	p := logic.FindPlayer(int(s.Pid))
	if p == nil {
		fmt.Printf("no player : %v\n", s)
		return
	}
	rid := p.GetRommID()
	r := logic.FindRoom(rid)
	if r == nil {
		fmt.Printf("player not in room : %d\n", rid)
		return
	}

	r.NotifyAllPlayers(s)

}

func (h *Handle) GetHandle(id int) func(data interface{}, c gnet.Conn) {
	switch id {
	case pb.ACTION:
		return C2SPlayerAction
	case pb.FRAMES:
		return FrameToBuffer
	}
	return nil
}
