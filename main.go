package main

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"github.com/Happy2018new/nemc-tan-lobby-solver/protocol/encoding"
	"github.com/Happy2018new/nemc-tan-lobby-solver/protocol/packet"
	"github.com/OmineDev/flowers-for-machines/core/minecraft/raknet"
	"github.com/df-mc/go-nethernet"
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
	str := `fee301000084467b526bdff284bc3fc4ac361337884210b163c23fe75117ee2f7f2f19f014a771e7b50c4861707079323031386e6577`
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

	userToken := "nPs8yqZC9Ldr0b17"
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

	conn, err := raknet.Dial("117.147.202.234:10000")
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
	bs, err = hex.DecodeString(`fee3018c78fc0baa0f7b4a0cf8fc74c7ddfaacbcc6e8cbb67908fe868325609d42efb400565511dde44f54bbc97377360f3bd74fc3f184f39eec`)
	conn.Write(bs)

	bs, err = dec.Decode()
	if err != nil {
		panic(err)
	}
	fmt.Println(bs)

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

	// nethernet.Dialer{}.DialContext()
	// ll := nethernet.Listener{}
	nethernet.Dialer{}.DialContext()
	kfc, err := raknet.Dialer{}.DialContext(context.Background(), "10.99.20.169:19146")
	// kfc, err := raknet.Dial("10.99.20.169:19146")
	if err != nil {
		panic(err)
	}
	fmt.Println(kfc)

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
