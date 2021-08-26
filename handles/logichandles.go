package handles

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"gameserver/pb"
	"gameserver/pb/actions"
	"gameserver/pb/msgerr"
	"unsafe"

	"github.com/panjf2000/gnet"
)

func S2CSendErrorMessage(error_code msgerr.ERROR_CODE, c gnet.Conn) {
	data := &msgerr.Error{
		Err: error_code,
	}
	mh := &pb.MessageHeader{
		ID:    pb.S2C_PLAYER_LOGIN,
		Lenth: int32(unsafe.Sizeof(*data)),
	}

	buf := &bytes.Buffer{}

	d := pb.Encode(data)
	err := binary.Write(buf, binary.LittleEndian, mh)
	if err != nil {
		fmt.Println("player login faild!", err)
		return
	}
	msg := append(buf.Bytes(), d...)
	c.AsyncWrite(msg)
	return

}

func SendPlayerLoginAck(a *actions.C2SPlayerAction, c gnet.Conn) error {

	data := &actions.S2CPlayerActionAck{
		Pid: a.Pid,
		Ack: actions.ACTION_ACK_LOGIN_SUCCESS,
	}
	fmt.Println("----->", data)
	return pb.SendToClinet(pb.S2C_PLAYER_LOGIN, data, int32(unsafe.Sizeof(*data)), c)
}

func SendPlayerJoinRoomAck(pid int, rid int, c gnet.Conn) error {

	data := &actions.S2CJoinRoomAck{
		Rid: int32(rid),
		Ack: actions.JOINROOM_ACK_JOINROOM_SUCCESS,
	}
	return pb.SendToClinet(pb.S2C_PLAYER_JOIN_ROOM, data, int32(unsafe.Sizeof(*data)), c)
}
