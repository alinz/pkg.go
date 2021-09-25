package generator

import "github.com/rs/xid"

type UUID xid.ID

func (u UUID) String() string {
	return xid.ID(u).String()
}

func NewUUID() UUID {
	return UUID(xid.New())
}
