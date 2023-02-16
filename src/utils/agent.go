package utils

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"os"
)

func GetMACAddress() (string, error) {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		panic(err.Error())
	}
	mac, macerr := "", errors.New("无法获取到正确的MAC地址")
	for i := 0; i < len(netInterfaces); i++ {
		//fmt.Println(netInterfaces[i])
		if (netInterfaces[i].Flags&net.FlagUp) != 0 && (netInterfaces[i].Flags&net.FlagLoopback) == 0 {
			addrs, _ := netInterfaces[i].Addrs()
			for _, address := range addrs {
				ipnet, ok := address.(*net.IPNet)
				//fmt.Println(ipnet.IP)
				if ok && ipnet.IP.IsGlobalUnicast() {
					// 如果IP是全局单拨地址，则返回MAC地址
					mac = netInterfaces[i].HardwareAddr.String()
					return mac, nil
				}
			}
		}
	}
	return mac, macerr
}

func GetPacketHash() string {
	path := "/opt/pokeball/attachment/pokeball.tar.gz"
	pFile, err := os.Open(path)
	if err != nil {
		log.Errorf("open err ，path=%v, err=%v", path, err)
		return "dev"
	}
	defer pFile.Close()
	md5h := md5.New()
	_, err = io.Copy(md5h, pFile)
	if err != nil {
		log.Errorf("md5 sum err  err=%v", err)
		return "error"
	}

	return hex.EncodeToString(md5h.Sum(nil))
}
