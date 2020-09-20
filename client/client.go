package client

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/jmattheis/shargo/encrypt"
	"github.com/jmattheis/shargo/proto"
	"image"
	"image/png"
	"io"
	"net"
	"time"
)

func Client(addr string, password []byte, images chan<- image.Image) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	enc, err := encrypt.Encrypt(password, &proto.Packet{
		Control: proto.ControlHello,
		Payload: []byte(time.Now().Format(time.RFC3339)),
	})
	if err != nil {
		return fmt.Errorf("could not encrypt hello: %s", err)
	}

	_, err = conn.Write(enc)
	if err != nil {
		return fmt.Errorf("could not write hello: %s", err)
	}
	buf := bytes.Buffer{}

	for {
		buf.Reset()
		size := new(int64)
		err := binary.Read(conn, binary.LittleEndian, size)
		if err != nil {
			return fmt.Errorf("read failed: %s", err)
		}
		fmt.Println("read from 2", *size)
		n, err := buf.ReadFrom(io.LimitReader(conn, *size))
		if err != nil {
			return fmt.Errorf("read failed: %s", err)
		}
		fmt.Println("read from 3", *size, n)

		if err != nil {
			return fmt.Errorf("read failed: %s", err)
		}

		p, err := encrypt.Decrypt(password, buf.Bytes()[0:n])
		if err != nil {
			return fmt.Errorf("could not decrypt: %s %d", err, n)
		}
		if p.Control != proto.ControlImage {
			return fmt.Errorf("unknown control %d", p.Control)
		}

		img, err := png.Decode(bytes.NewBuffer(p.Payload))
		if err != nil {
			return fmt.Errorf("could not decode png: %s", err)
		}

		images <- img
	}
}
