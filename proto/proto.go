package proto

import (
	"gopkg.in/mgo.v2/bson"
	"io"
)

type Control byte

func (c Control) Write(w io.Writer, payload []byte) error {
	b, err := bson.Marshal(&Packet{Control: c, Payload: payload})
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}

const (
	ControlHello Control = iota
	ControlImage
)

type Packet struct {
	Control Control
	Payload []byte
}
