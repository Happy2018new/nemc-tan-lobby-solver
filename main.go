package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"github.com/Happy2018new/nemc-tan-lobby-solver/bunker/auth"
	"github.com/Happy2018new/nemc-tan-lobby-solver/protocol/encoding"
	"github.com/Happy2018new/nemc-tan-lobby-solver/protocol/login"
	"github.com/Happy2018new/nemc-tan-lobby-solver/protocol/packet"
	"github.com/OmineDev/flowers-for-machines/core/minecraft/raknet"
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
	cli, err := auth.CreateClient(&auth.ClientOptions{
		AuthServer: "http://127.0.0.1:8080",
	})
	if err != nil {
		panic(err)
	}

	wrapper := auth.NewAccessWrapper(cli, "484434", "", "...")
	// resp, err := wrapper.RaknetServerList()
	// fmt.Printf("%#v, %v\n", resp, err)
	sss, err := login.Dial(wrapper)
	fmt.Println(sss, err)
	return

	str := `fee301000084467b5293e06c151956a4decdd8b3fc5f51e9511be51cceada7a6d355877b47dc7d92100c4861707079323031386e6577`
	bs, err := hex.DecodeString(str)
	if err != nil {
		panic(err)
	}
	dec, err := packet.NewDecoder(NewSingleReader(bs))
	if err != nil {
		panic(err)
	}

	pks, err := dec.Decode()
	if err != nil {
		panic(err)
	}

	buf := bytes.NewBuffer(pks)
	header := packet.Header{}
	header.Read(buf)

	r := encoding.NewReader(buf)
	loginRequest := packet.TanLoginRequest{}
	loginRequest.Marshal(r)
	fmt.Printf("%#v\n", loginRequest)

	userToken := "gnEI3C0IIcGnuInu"
	encryptedUserToken := MD5Sum([]byte(userToken))
	encryptKeyBytes := []byte(string(encryptedUserToken) + string(loginRequest.Rand))
	decryptKeyBytes := []byte(string(loginRequest.Rand) + string(encryptedUserToken))

	// str = `fee3011c8b9d3334db7b2f056fdec7f9fe59189ac3c5bee1f6b999c4c0daecb8a163b9a4797e66d9ce6d68915c97ef6965aa2bf4`
	// bs, err = hex.DecodeString(str)
	// if err != nil {
	// 	panic(err)
	// }
	// dec, err = packet.NewDecoder(NewSingleReader(bs))
	// if err != nil {
	// 	panic(err)
	// }
	// dec.EnableEncryption(encryptKeyBytes, decryptKeyBytes)
	// dec.EnableEncryption(encryptKeyBytes, encryptKeyBytes)

	// pks, err = dec.Decode()
	// if err != nil {
	// 	panic(err)
	// }

	// buf = bytes.NewBuffer(pks)
	// header = packet.Header{}
	// header.Read(buf)

	// r = encoding.NewReader(buf)
	// pk := packet.TanEnterRoomRequest{}
	// pk.Marshal(r)
	// return

	conn, err := raknet.Dial("117.147.202.212:10000")
	if err != nil {
		panic(err)
	}
	dec, err = packet.NewDecoder(conn)
	if err != nil {
		panic(err)
	}

	conn.Write(bs)
	pkData, err := dec.Decode()
	if err != nil {
		panic(err)
	}
	fmt.Println(pkData)

	dec.EnableEncryption(encryptKeyBytes, decryptKeyBytes)
	// dec.EnableEncryption(decryptKeyBytes, decryptKeyBytes)
	bs, err = hex.DecodeString(`fee30173f6c8738e06732ca6a001c204f5293707d5aa2f78b340b6b36ae717119d0f53c5a12327b60e9db09d3852a684eab48737d32e90130e87`)
	conn.Write(bs)
	{
		ddddd, err := packet.NewDecoder(NewSingleReader(bs))
		if err != nil {
			panic(err)
		}
		ddddd.EnableEncryption(encryptKeyBytes, encryptKeyBytes)
		pkData, err = ddddd.Decode()
		if err != nil {
			panic(err)
		}
		buf = bytes.NewBuffer(pkData)
		header = packet.Header{}
		header.Read(buf)
		r = encoding.NewReader(buf)
		pk := packet.TanEnterRoomRequest{}
		pk.Marshal(r)
		fmt.Printf("%#v\n", pk)
	}

	bs, err = dec.Decode()
	if err != nil {
		panic(err)
	}
	{
		buf = bytes.NewBuffer(bs)
		header = packet.Header{}
		header.Read(buf)
		r = encoding.NewReader(buf)
		pk := packet.TanEnterRoomResponse{}
		pk.Marshal(r)
		fmt.Printf("%#v\n", pk)
	}

	bs, err = dec.Decode()
	if err != nil {
		panic(err)
	}
	buf = bytes.NewBuffer(bs)
	header = packet.Header{}
	header.Read(buf)
	r = encoding.NewReader(buf)
	pk := packet.TanNotifyServerReady{}
	pk.Marshal(r)
	fmt.Printf("%#v\n", pk) //10.99.20.169|19146

	// d := nethernet.Dialer{}
	// d.DialContext(context.Background(), 0, 17508800308048208044, nil)

	// 2219211602 (user id)
	// 16134042324431409537
	// 11417861250525914572

	// nethernet.Dialer{}.DialContext()
	// ll := nethernet.Listener{}
	// nethernet.Dialer{}.DialContext()
	// kfc, err := raknet.Dialer{}.DialContext(context.Background(), "10.99.20.169:19146")
	// kfc, err := raknet.Dial("10.99.20.169:19146")
	if err != nil {
		panic(err)
	}
	// fmt.Println(kfc)

	for {
		bs, err := dec.Decode()
		if err != nil {
			panic(err)
		}
		fmt.Println(bs)
	}
}

func MD5Sum(data []byte) []byte {
	result := md5.Sum(data)
	return result[:]
}
