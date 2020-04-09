package lib

import (
	"errors"
	"fmt"
	"log"
	"time"
)

/*批量检查帐号*/
func CheckAccounts(accounts []DialAccount, eth string, fix bool, debug bool) error {
	if len(accounts) > 244 {
		return errors.New("最大的帐号数不能超过244个")
	}

	var runNum = make(chan int, 252)
	for i, v := range accounts {
		pppNumber := i + 10
		ppp := fmt.Sprintf("ppp%v", pppNumber)
		macVLan := fmt.Sprintf("macvlan%v", pppNumber)
		password := v.Password
		account := v.Account
		accounts[i].IsActive = false
		mac := GetRandMac()
		accounts[i].Ppp = ppp
		accounts[i].MacVLan = macVLan
		accounts[i].Mac = mac
		if fix {
			go DialOne(eth, ppp, macVLan, mac, account, password, &runNum, debug, 1)
		} else {
			go DialOne(eth, ppp, macVLan, mac, account, password, &runNum, debug)
		}
	}

	i := 0
	countAccount := len(accounts)
Loop:
	for {
		select {
		case <-runNum:
			i++
			if i == countAccount {
				break Loop
			}
		case <-time.After(5 * time.Second):
			break Loop
		}
	}
	fmt.Println("批量拨号结束，开始验证帐号可用性")

	netMap := GetAllDst()
	for _, row := range accounts {
		netRow, ok := netMap[row.Ppp]
		if netRow.Ip == "" || netRow.Dst == "" {
			log.Println(row.Ppp, "的账户：", row.Account, "在获取IP和目标地址时发生错误，开始重试")
			tryTimes := 5
			for t := 1; t <= tryTimes; t++ {
				log.Println(row.Ppp, "的账户：", row.Account, "第", t, "次重试")
				if t == 3 {
					log.Println("启用修复拨号")
					var runNum = make(chan int, 252)
					KillByName(row.Ppp)
					go DialOne(eth, row.Ppp, row.MacVLan, row.Mac, row.Account, row.Password, &runNum, debug, 1)
					<-runNum
					log.Println(row.Ppp + "口的账户：" + row.Account + "尝试修复完毕\n")
				}
				netRow.Dst, _ = GetDst(row.Ppp)
				netRow.Ip, _ = GetIP(row.Ppp)
				if netRow.Dst != "" && netRow.Ip != "" {
					ok = true
					log.Println(row.Ppp, "在获取IP和目标地址时发生错误，重试成功,ip:", netRow.Ip, " dst:", netRow.Dst)
					break
				}
				log.Println(row.Ppp, "的账户：", row.Account, "第", t, "次重试失败，2秒后开始继续尝试")
				time.Sleep(time.Second * 2)
			}

			if netRow.Ip == "" || netRow.Dst == "" {
				log.Println(row.Ppp, "的账户：", row.Account, "重试失败")
			}
		}

		if ok {
			fmt.Printf("%s  可以使用\n", row.Account)
		} else {
			fmt.Printf("%s  不能使用\n", row.Account)
		}
	}
	return nil
}
