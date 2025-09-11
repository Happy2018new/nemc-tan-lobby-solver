package main

import (
	"fmt"

	"github.com/Happy2018new/nemc-tan-lobby-solver/raknet"
)

func main() {
	conn, err := raknet.Dial("117.147.202.226:10000")
	if err != nil {
		panic(err)
	}
	for {
		bs, err := conn.ReadPacket()
		if err != nil {
			panic(err)
		}
		fmt.Println(bs)
	}
}
