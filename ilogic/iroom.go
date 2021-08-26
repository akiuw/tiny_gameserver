package ilogic

type IRoom interface {
	GetRoomID() (rid int)
	AddPlayer(pid int) (err error)
	NotifyAllPlayers()
	PlayerLeaveRoom(pid int)
}
