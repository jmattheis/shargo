package host

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/jmattheis/shargo/encrypt"
	"github.com/jmattheis/shargo/proto"
	"image"
	"image/png"
	"net"
)

func Server(addr string, images <-chan image.Image, password []byte) error {
	srv, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("could not listen on %s", addr)
	}
	fmt.Println("listening on ", addr)

	newConn := make(chan net.Conn)
	closeConn := make(chan net.Conn)
	go func() {
		for {
			conn, err := srv.Accept()
			if err != nil {
				continue
			}

			go func() {
				defer func() {
					closeConn <- conn
				}()

				authed := false
				for {
					buf := make([]byte, 1024)
					n, err := conn.Read(buf)
					if err != nil {
						return
					}
					p, err := encrypt.Decrypt(password, buf[0:n])
					if err != nil {
						fmt.Println("bad decrypt", err)
						return
					}
					if p.Control != proto.ControlHello {
						fmt.Println("bad control", p.Control)
						return
					}
					if !authed {
						fmt.Println("received valid hello", conn.RemoteAddr().String())
						authed = true
						newConn <- conn
					}
				}
			}()

		}
	}()
	connected := map[string]net.Conn{}
	for {
		buf := &bytes.Buffer{}
		select {
		case screen := <-images:
			buf.Reset()
			_ = png.Encode(buf, screen)

			encrypted, err := encrypt.Encrypt(password, &proto.Packet{
				Control: proto.ControlImage,
				Payload: buf.Bytes(),
			})
			buf.Reset()
			_ = binary.Write(buf, binary.LittleEndian, int64(len(encrypted)))
			buf.Write(encrypted)
			if err != nil {
				fmt.Println("could not encrypt", err)
				return nil
			}

			for _, conn := range connected {
				go func(conn net.Conn) {
					if _, err := conn.Write(buf.Bytes()); err != nil {
						closeConn <- conn
					}
				}(conn)
			}
		case conn := <-newConn:
			connected[conn.RemoteAddr().String()] = conn
		case conn := <-closeConn:
			_ = conn.Close()
			delete(connected, conn.RemoteAddr().String())
		}

	}
}
