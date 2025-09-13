package main

import (
	"bytes"
	"context"
	"fmt"

	"github.com/Happy2018new/nemc-tan-lobby-solver/bunker/auth"
	"github.com/Happy2018new/nemc-tan-lobby-solver/minecraft"
	"github.com/Happy2018new/nemc-tan-lobby-solver/protocol/login"
)

type SingleReader struct {
	*bytes.Buffer
}

func NewSingleReader(bs []byte) *SingleReader {
	return &SingleReader{Buffer: bytes.NewBuffer(bs)}
}

func (s *SingleReader) ReadPacket() ([]byte, error) {
	return s.Bytes(), nil
}

func main() {
	client, err := auth.CreateClient(&auth.ClientOptions{
		AuthServer: "http://127.0.0.1:8080",
	})
	if err != nil {
		panic(err)
	}

	wrapper := auth.NewAccessWrapper(client, "484434", "", "YOUR FB TOKEN")
	netConn, err := login.Dial(wrapper)
	if err != nil {
		panic(err)
	}

	conn, err := minecraft.DialContext(context.Background(), netConn)
	if err != nil {
		panic(err)
	}
	for {
		pk, err := conn.ReadPacket()
		if err != nil {
			panic(err)
		}
		fmt.Printf("%#v\n", pk)
	}
}
