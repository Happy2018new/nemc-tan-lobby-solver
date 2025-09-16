package main

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/Happy2018new/nemc-tan-lobby-solver/bunker/auth"
	"github.com/Happy2018new/nemc-tan-lobby-solver/minecraft"
	"github.com/Happy2018new/nemc-tan-lobby-solver/protocol/login"
)

func main() {
	client, err := auth.CreateClient(&auth.ClientOptions{
		AuthServer: "https://yorha.eulogist-api.icu",
	})
	if err != nil {
		panic(err)
	}

	wrapper := auth.NewAccessWrapper(client, "ROOM ID", "", "YOUR FB TOKEN")
	netConn, tanLobbyLoginResp, err := login.Dial(wrapper)
	if err != nil {
		panic(err)
	}
	fmt.Printf("[INFO] Tan lobby login response: %#v\n", tanLobbyLoginResp)

	serverConn, err := minecraft.DialContext(context.Background(), netConn)
	if err != nil {
		panic(err)
	}
	defer serverConn.Close()

	listenConfig := minecraft.ListenConfig{
		AuthenticationDisabled: true,
	}
	listener, err := listenConfig.Listen("raknet", ":19132")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	fmt.Println("[INFO] Server starting on 127.0.0.1:19132")
	c, err := listener.Accept()
	if err != nil {
		panic(err)
	}
	handleConn(c.(*minecraft.Conn), serverConn, listener)
}

// handleConn handles a new incoming minecraft.Conn from the minecraft.Listener passed.
func handleConn(clientConn *minecraft.Conn, serverConn *minecraft.Conn, listener *minecraft.Listener) {
	var spawnWaitGroup sync.WaitGroup
	var normalWaitGroup sync.WaitGroup

	spawnWaitGroup.Add(2)
	go func() {
		if err := clientConn.StartGame(serverConn.GameData()); err != nil {
			panic(err)
		}
		spawnWaitGroup.Done()
	}()
	go func() {
		if err := serverConn.DoSpawn(); err != nil {
			panic(err)
		}
		spawnWaitGroup.Done()
	}()
	spawnWaitGroup.Wait()

	normalWaitGroup.Add(2)
	go func() {
		defer normalWaitGroup.Done()
		defer listener.Disconnect(clientConn, "connection lost")
		defer serverConn.Close()
		for {
			pk, err := clientConn.ReadPacket()
			if err != nil {
				return
			}
			if err := serverConn.WritePacket(pk); err != nil {
				var disc minecraft.DisconnectError
				if ok := errors.As(err, &disc); ok {
					_ = listener.Disconnect(clientConn, disc.Error())
				}
				return
			}
		}
	}()
	go func() {
		defer normalWaitGroup.Done()
		defer serverConn.Close()
		defer listener.Disconnect(clientConn, "connection lost")
		for {
			pk, err := serverConn.ReadPacket()
			if err != nil {
				var disc minecraft.DisconnectError
				if ok := errors.As(err, &disc); ok {
					_ = listener.Disconnect(clientConn, disc.Error())
				}
				return
			}
			if err := clientConn.WritePacket(pk); err != nil {
				return
			}
		}
	}()
	normalWaitGroup.Wait()
}
