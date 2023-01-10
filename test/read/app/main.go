package main

import (
	"golang.org/x/sys/windows"
	"golang.zx2c4.com/wintun"
	"golang.zx2c4.com/wireguard/tun"
	"log"
	"net/netip"
	"simpletun/winipcfg"
	"sync/atomic"
	"time"
)

func main() {
	id := &windows.GUID{
		0x0000000,
		0xFFFF,
		0xFFFF,
		[8]byte{0xFF, 0xe9, 0x76, 0xe5, 0x8c, 0x74, 0x06, 0x3e},
	}
	_ = wintun.Uninstall()
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
	log.Println(">> running...")
	ticker := time.NewTicker(time.Second)
	var cnt uint32 = 0
	var total uint32 = 0
	go func() {
		for {
			<-ticker.C
			log.Printf(">> %6d pkt/s %6d MB/s\n", atomic.LoadUint32(&cnt), total/1024/1024)
			atomic.StoreUint32(&cnt, 0)
			atomic.StoreUint32(&total, 0)
		}
	}()
	for {
		n = 2048
		n, err = dev.Read(buf, 0)
		if err != nil {
			panic(err)
		}
		atomic.AddUint32(&cnt, 1)
		atomic.AddUint32(&total, uint32(n))
	}

}
