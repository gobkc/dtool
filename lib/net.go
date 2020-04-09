package lib

import (
	"log"
	"net"
	"os/exec"
	"strings"
)

type NetSpace struct {
	Ip  string //ip地址
	Dst string //目标地址
}

/*获取目标地址*/
func GetDst(eth string) (ip string, err error) {
	cli := exec.Command("bash", "-c", "ifconfig "+eth+"|grep -E 'inet'| head -n 1|awk '{print($6)}'")
	if ipAddr, err := cli.Output(); err != nil {
		log.Println(err)
		return "", err
	} else {
		ip := strings.TrimSpace(string(ipAddr))
		return ip, err
	}
}

/*获取IP*/
func GetIP(eth string) (ip string, err error) {
	cli := exec.Command("bash", "-c", "ifconfig "+eth+"|grep -E 'inet'| head -n 1|awk '{print($2)}'")
	if ipAddr, err := cli.Output(); err != nil {
		log.Println(err)
		return "", err
	} else {
		ip := strings.TrimSpace(string(ipAddr))
		return ip, err
	}
}

/*获取所有网卡的IP和目标地址*/
func GetAllDst() map[string]NetSpace {
	iFaces, _ := net.Interfaces()
	var netMap = make(map[string]NetSpace)
	for _, iFace := range iFaces {
		iFaceAddr, _ := iFace.Addrs()
		for _, ip := range iFaceAddr {
			iFaceReal := ip.(*net.IPNet)
			address := iFaceReal.IP.String()
			if strings.Count(address, ":") < 2 {
				ipV4 := address
				dst, _ := GetDst(iFace.Name)
				/*重试机制*/
				netMap[iFace.Name] = NetSpace{Ip: ipV4, Dst: dst}
			}
		}
	}
	return netMap
}