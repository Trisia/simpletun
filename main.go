package main

import (
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/sys/windows"
	"golang.zx2c4.com/wireguard/tun"
	"log"
	"net/netip"
	"simpletun/winipcfg"
)

func main() {
	id := &windows.GUID{
		0x0000000,
		0xFFFF,
		0xFFFF,
		[8]byte{0xFF, 0xe9, 0x76, 0xe5, 0x8c, 0x74, 0x06, 0x3e},
	}
	ifname := "MyNIC"
	dev, err := tun.CreateTUNWithRequestedGUID(ifname, id, 0)
	if err != nil {
		panic(err)
	}
	defer dev.Close()
	// 保存原始设备句柄
	nativeTunDevice := dev.(*tun.NativeTun)

	// 获取LUID用于配置网络
	link := winipcfg.LUID(nativeTunDevice.LUID())

	ip, err := netip.ParsePrefix("10.0.0.77/24")
	if err != nil {
		panic(err)
	}
	err = link.SetIPAddresses([]netip.Prefix{ip})
	if err != nil {
		panic(err)
	}

	n := 2048
	buf := make([]byte, n)

	for {
		n = 2048
		n, err = dev.Read(buf, 0)
		if err != nil {
			panic(err)
		}
		const ProtocolICMP = 1
		header, err := ipv4.ParseHeader(buf[:n])
		if err != nil {
			continue
		}
		if header.Protocol == ProtocolICMP {
			log.Println("Src:", header.Src, " dst:", header.Dst)
			msg, _ := icmp.ParseMessage(ProtocolICMP, buf[header.Len:])
			log.Println(">> ICMP:", msg.Type)
		}
	}
}
