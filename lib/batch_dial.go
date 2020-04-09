package lib

import (
	"errors"
	"fmt"
	"log"
	"time"
)

func BatchDialing(accounts []DialAccount, eth string, interval int, flow int, fix bool, debug bool) error {
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
	fmt.Println("批量拨号结束，开始跑流量")

	netMap := GetAllDst()
	for i, row := range accounts {
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
					if debug{
						log.Println(row.Ppp, "在获取IP和目标地址时发生错误，重试成功,ip:", netRow.Ip, " dst:", netRow.Dst)
					}
					break
				}
				if debug{
					log.Println(row.Ppp, "的账户：", row.Account, "第", t, "次重试失败，2秒后开始继续尝试")
				}
				time.Sleep(time.Second * 2)
			}

			if netRow.Ip == "" || netRow.Dst == "" {
				if debug{
					log.Println(row.Ppp, "的账户：", row.Account, "重试失败")
				}
			}
		}

		if ok {
			accounts[i].IsActive = true
			var rtId int
			var err error

			//同步RT TABLE
			if rtId = SetRtTable(row.MacVLan); rtId == 0 {
				if debug{
					log.Println(row.Ppp, "同步RT TABLE失败")
				}
				accounts[i].IsActive = false
				continue
			}

			DelRuleByTableId(rtId)

			if err = SetPPPOE(row.Ppp, rtId); err != nil {
				if debug{
					log.Println(row.Ppp, "设置PPPOE失败")
				}
				accounts[i].IsActive = false
				continue
			}

			if err := AddDefaultRoute(row.Ppp, netRow.Dst, rtId); err != nil {
				if debug{
					log.Println(row.Ppp, "设置默认路由失败")
				}
				accounts[i].IsActive = false
				continue
			}

			if err := AddDns("114.114.114.114/24", rtId, 5); err != nil {
				if debug{
					log.Println(row.Ppp, "设置DNS失败")
				}
				accounts[i].IsActive = false
				continue
			}

			fmt.Println("ppp:", row.Ppp, " ip:", netRow.Ip, " dst:", netRow.Dst, " 设定成功")
		} else {
			accounts[i].IsActive = false
		}
	}
	if debug{
		log.Println("所有帐号设定结束")
	}
	fmt.Println(fmt.Sprintf("\n%c[1;33;40m可使用账户列表：%c[0m", 0x1B, 0x1B))
	for _, row := range accounts {
		if row.IsActive{
			fmt.Println(row.Account,"      可以使用")
		}
	}
	fmt.Println(fmt.Sprintf("\n%c[1;33;40m不可使用账户列表：%c[0m", 0x1B, 0x1B))
	for _, row := range accounts {
		if !row.IsActive{
			fmt.Println(row.Account,"      不可使用")
		}
	}

	log.Println("5秒后开始跑流量")
	time.Sleep(5 * time.Second)

	tick := time.NewTicker(time.Second)
	var t = 0
	for {
		select {
		//此处在等待channel中的信号，因此执行此段代码时会阻塞120秒
		case <-tick.C:
			if t == 0 {
				fmt.Printf("\n")
				var requestUrl = "http://mirrors.163.com/mysql/Downloads/MySQL-4.1/MySQL-4.1.21-0.glibc23.src.rpm"
				for _, row := range accounts {
					account := row.Account
					if row.IsActive {
						SetDefaultRouter(row.Ppp)
						log.Println("当前账户：", account, " 使用PPP:", row.Ppp)
						// todo 这里等待处理，只是把问题绕过去了
						//ip := netMap[row.Ppp]
						reTry := 3
						for i := 1; i <= reTry; i++ {
							time1 := time.Now().Unix()
							Download2(requestUrl, func(length, downLen int64) {
								fmt.Print("\r", account, "开始跑流量:", B2M(downLen), "mb                ")
							})
							fmt.Print("\n")
							time2 := time.Now().Unix()
							if time2-time1 > 10 {
								break
							}
							fmt.Println("跑流量时间不够，重试第", i, "次")
						}
						fmt.Println("结束跑流量,等待中...")
						time.Sleep(time.Second * 3)
					}
				}
				t = interval * 60 * 60
			} else {
				fmt.Printf("\r还有%v小时%v分%v秒后执行下一队列，间隔时间：%v hour       ",
					t/60/60, 60-(interval*60-t/60)%60, 60-(interval*60*60-t)%60, interval)
			}
			t--
		}
	}
	return nil
}
