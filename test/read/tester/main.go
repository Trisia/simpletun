package main

import (
	"crypto/rand"
	"net"
)

const DataSize = 32

//const DataSize = 128
//const DataSize = 516
//const DataSize = 1420

func main() {
	dst, err := net.ResolveUDPAddr("udp", "192.168.0.220:7878")
	if err != nil {
		panic(err)
	}
	conn, err := net.DialUDP("udp", nil, dst)
	if err != nil {
		panic(err)
	}
	buf := make([]byte, DataSize)
	_, _ = rand.Reader.Read(buf)
	for {
		_, err := conn.Write(buf)
		if err != nil {
			panic(err)
		}
		//fmt.Println(">> Send data")
	}
}
