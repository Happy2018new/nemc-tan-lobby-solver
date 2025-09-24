package main

import (
	"context"
	"fmt"

	"github.com/Happy2018new/nemc-tan-lobby-solver/bunker/auth"
	"github.com/Happy2018new/nemc-tan-lobby-solver/core/nethernet"
	"github.com/Happy2018new/nemc-tan-lobby-solver/core/raknet"
	"github.com/Happy2018new/nemc-tan-lobby-solver/minecraft"
	"github.com/Happy2018new/nemc-tan-lobby-solver/protocol/service"
)

func main() {
	if false {
		client, err := auth.CreateClient(&auth.ClientOptions{
			AuthServer: "AUTH SERVER ADDRESS",
		})
		if err != nil {
			panic(err)
		}
		wrapper := auth.NewServerAccessWrapper(client, "YOUR FB TOKEN")

		listenConfig, listener, roomID, err := service.Listen(wrapper, "ROOM NAME")
		if err != nil {
			panic(err)
		}
		defer listenConfig.CloseRoom()

		fmt.Printf("[SUCCESS] Create room: %d", roomID)
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
		client, err := auth.CreateClient(&auth.ClientOptions{
			AuthServer: "http://127.0.0.1:80",
		})
		if err != nil {
			panic(err)
		}

		wrapper := auth.NewClientAccessWrapper(client, "ROOM ID", "", "YOUR FB TOKEN")
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
