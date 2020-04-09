package lib

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"syscall"
)

type DialRuntime struct {
	Password string
	Status   bool
	WaitFlow int
	Ppp      string
	PppUnit  int
	MacVLan  string
	Mac      string
}

func DialOne(eth string, ppp string, macVLan string, mac string, account string, password string, runNum *chan int, debug bool, useMacVLan ...interface{}) {
	if err := AddMacVLan(macVLan, eth, mac); err != nil {
		if debug {
			log.Printf("在给账户：%s拨号时，添加macvlan失败，详细原因：%s", account, err.Error())
		}
		return
	}
	pppUnit := strings.Replace(ppp, "ppp", "", 1)
	var cmdString string
	if len(useMacVLan) > 0 {
		cmdString = fmt.Sprintf("/usr/sbin/pppd pty '/usr/sbin/pppoe -p /var/run/%s-pppoe.pid.pppoe -I %s -T 80 -U  -m 1412' unit %s ipparam %s linkname %s noipdefault noauth default-asyncmap hide-password nodetach mtu 1492 mru 1492 noaccomp nodeflate usepeerdns nopcomp novj novjccomp user %s password '%s' lcp-echo-interval 20 lcp-echo-failure 3", macVLan, eth, pppUnit, ppp, ppp, account, password)
	} else {
		cmdString = fmt.Sprintf("/usr/sbin/pppd pty '/usr/sbin/pppoe -p /var/run/%s-pppoe.pid.pppoe -I %s -T 80 -U  -m 1412' unit %s ipparam %s linkname %s noipdefault noauth default-asyncmap hide-password nodetach mtu 1492 mru 1492 noaccomp nodeflate usepeerdns nopcomp novj novjccomp user %s password '%s' lcp-echo-interval 20 lcp-echo-failure 3", macVLan, macVLan, pppUnit, ppp, ppp, account, password)
	}
	if debug {
		fmt.Println("执行的CMD:\n     ", cmdString)
	}
	ctx, _ := context.WithCancel(context.Background())

	cmd := exec.CommandContext(ctx, "bash", "-c", cmdString)
	cmd.SysProcAttr = &syscall.SysProcAttr{}

	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	cmd.Start()
	*runNum <- 1
	cmd.Wait()
}
