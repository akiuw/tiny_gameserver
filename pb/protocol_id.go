package pb

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"runtime"

	"github.com/panjf2000/gnet"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var ProtocolMap map[int32]func(obj interface{}, c gnet.Conn)

type MessageHeader struct {
	ID    int32
	Lenth int32
}

func init() {
	ProtocolMap = make(map[int32]func(obj interface{}, c gnet.Conn))
}

const (
	MESSAGE_BEGIN = iota
	ACTION
	FRAMES
	MESSAGE_END
	ACK_BEGIN
	S2C_PLAYER_LOGIN
	S2C_PLAYER_JOIN_ROOM
)

func RegistProtocol(protoid int32, f func(obj interface{}, c gnet.Conn)) {
	ProtocolMap[protoid] = f
}

func Encode(msg interface{}) []byte {
	data, err := proto.Marshal(msg.(protoreflect.ProtoMessage))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return data
}

func NewError(str string) error {
	funcName, file, line, ok := runtime.Caller(0)
	if ok {
		return errors.New(fmt.Sprintf("file: %s, funcname: %s, line: %d", file, runtime.FuncForPC(funcName).Name(), line) + str)
	}
	return errors.New("now pc pointer : 0")
}

func SendToClinet(messageid int, data interface{}, lenth int32, c gnet.Conn) error {
	mh := &MessageHeader{
		ID:    int32(messageid),
		Lenth: lenth,
	}

	buf := &bytes.Buffer{}

	d := Encode(data)
	err := binary.Write(buf, binary.LittleEndian, mh)
	if err != nil {
		return NewError("sent to client error : binary.Write")
	}
	msg := append(buf.Bytes(), d...)
	c.AsyncWrite(msg)
	return nil
}
