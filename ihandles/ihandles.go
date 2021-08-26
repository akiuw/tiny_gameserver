package ihandles

import "github.com/panjf2000/gnet"

type IHandles interface {
	GetHandle(id int) func(data interface{}, c gnet.Conn)
}
