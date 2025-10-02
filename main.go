package main

import (
	"context"
	"fmt"

	"github.com/Happy2018new/nemc-tan-lobby-solver/bunker"
	"github.com/Happy2018new/nemc-tan-lobby-solver/core/nethernet"
	"github.com/Happy2018new/nemc-tan-lobby-solver/core/raknet"
	"github.com/Happy2018new/nemc-tan-lobby-solver/minecraft"
	"github.com/Happy2018new/nemc-tan-lobby-solver/protocol/service"
)

func main() {
	if false {
		listenConfig, listener, roomID, err := service.Listen(
			service.DefaultRoomConfig("来和我一起玩吧！", "", 10, service.PlayerPermissionMember),
			bunker.NewAccessWrapper("AUTH SERVER ADDRESS", "YOUR FB TOKEN", "PE AUTH (CAN BE EMPTY)", "SA AUTH (CAN BE EMPTY)"),
		)
		if err != nil {
			panic(err)
		}
		defer listenConfig.CloseRoom()

		fmt.Printf("[SUCCESS] Create room: %d\n", roomID)
		for {
			clientConn, err := listener.Accept()
			if err != nil {
				panic(err)
			}

			serverConn, err := raknet.Dial("BDS SERVER ADDRESS")
			if err != nil {
				panic(err)
			}

			go func() {
				defer clientConn.Close()
				defer serverConn.Close()
				for {
					pkData, err := clientConn.(*nethernet.Conn).ReadPacket()
					if err != nil {
						return
					}
					serverConn.Write(append([]byte{0xfe}, pkData...))
				}
			}()

			go func() {
				defer clientConn.Close()
				defer serverConn.Close()
				for {
					pkData, err := serverConn.ReadPacket()
					if err != nil {
						return
					}
					clientConn.Write(pkData[1:])
				}
			}()
		}
	}

	if false {
		netConn, tanLobbyLoginResp, err := service.Dial(
			"ROOM ID",
			"ROOM PASSCODE",
			bunker.NewAccessWrapper("AUTH SERVER ADDRESS", "YOUR FB TOKEN", "PE AUTH (CAN BE EMPTY)", "SA AUTH (CAN BE EMPTY)"),
		)
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
