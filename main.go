package main

import (
	"fmt"

	"github.com/OmineDev/flowers-for-machines/core/minecraft/raknet"
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
