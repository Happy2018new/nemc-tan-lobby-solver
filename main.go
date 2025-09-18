package main

import (
	"context"
	"fmt"

	"github.com/Happy2018new/nemc-tan-lobby-solver/bunker/auth"
	"github.com/Happy2018new/nemc-tan-lobby-solver/minecraft"
	"github.com/Happy2018new/nemc-tan-lobby-solver/protocol/service"
)

func main() {
	if false {
		client, err := auth.CreateClient(&auth.ClientOptions{
			AuthServer: "http://127.0.0.1:80",
		})
		if err != nil {
			panic(err)
		}
		wrapper := auth.NewServerAccessWrapper(client, "YOUR FB TOKEN")

		listenConfig, listener, roomID, err := service.Listen(wrapper, "献给机械の花束")
		if err != nil {
			panic(err)
		}
		defer listenConfig.CloseRoom()

		fmt.Printf("[SUCCESS] Create room: %d", roomID)
		for {
			fmt.Println(listener.Accept())
		}
	}

	if false {
		client, err := auth.CreateClient(&auth.ClientOptions{
			AuthServer: "http://127.0.0.1:80",
		})
		if err != nil {
			panic(err)
		}

		wrapper := auth.NewClientAccessWrapper(client, "124674", "", "YOUR FB TOKEN")
		netConn, tanLobbyLoginResp, err := service.Dial(wrapper)
		if err != nil {
			panic(err)
		}
		fmt.Printf("[INFO] Tan lobby login response: %#v\n", tanLobbyLoginResp)

		serverConn, err := minecraft.DialContext(context.Background(), netConn)
		if err != nil {
			panic(err)
		}
		defer serverConn.Close()

		for {
			fmt.Println(serverConn.ReadPacket())
		}
	}
}
