package main

import (
	"flag"
	"fmt"
	"log"
	"new-dtool/lib"
	"os"
)

func main() {
	eth := flag.String("e", "eth0", "指定网卡，默认是eth0,如果不是此网卡，请手动指定")
	config := flag.String("c", "", "配置文件 如:1.txt")
	flow := flag.Int("f", 100, "每次最多跑多少流量，单位(Mb),默认100Mb")
	interval := flag.Int("i", 1, "间隔多少小时跑一次流量")
	fix := flag.Bool("fix", false, "修复模式，不使用macvlan，仅使用eth0来实现多拨。兼容性最高")
	check := flag.Bool("check", false, "不跑流量，只验证账户是否可用")
	debug := flag.Bool("debug", false, "调试模式，开发测试阶段使用")
	flag.Usage = func() {
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, fmt.Sprintf("%c[1;32;40m\n注意事项：%c[0m\n", 0x1B, 0x1B))
		fmt.Fprintf(os.Stderr, fmt.Sprintf("%c[1;33;40m (1) 本程序必须在多网卡环境下使用%c[0m\n", 0x1B, 0x1B))
		fmt.Fprintf(os.Stderr, fmt.Sprintf("%c[1;33;40m (2) 网卡必须支持promisc模式%c[0m\n", 0x1B, 0x1B))
		fmt.Fprintf(os.Stderr, fmt.Sprintf("%c[1;33;40m (3) 母机必须开启允许arp欺骗%c[0m\n", 0x1B, 0x1B))
		fmt.Fprintf(os.Stderr, fmt.Sprintf("%c[1;33;40m (4) 云服务器内核必须支持MACVLAN%c[0m\n", 0x1B, 0x1B))
		fmt.Fprintf(os.Stderr, fmt.Sprintf("%c[1;33;40m (5) 如果使用-e=eth0指定错误了网卡，请重启服务器，否则已经生成的MACVLAN无法撤销%c[0m\n", 0x1B, 0x1B))

		fmt.Fprintf(os.Stderr, fmt.Sprintf("%c[1;32;40m\n使用示例：%c[0m\n", 0x1B, 0x1B))
		fmt.Fprintf(os.Stderr, fmt.Sprintf("%c[1;33;40m (1) dtool -c=1.txt%c[0m\n", 0x1B, 0x1B))
		fmt.Fprintf(os.Stderr, fmt.Sprintf("%c[1;33;40m (2) dtool -c=1.txt -e=eth1%c[0m\n", 0x1B, 0x1B))
		fmt.Fprintf(os.Stderr, fmt.Sprintf("%c[1;33;40m (3) dtool -c=1.txt -fix%c[0m\n", 0x1B, 0x1B))
		fmt.Fprintf(os.Stderr, fmt.Sprintf("%c[1;33;40m (4) dtool -c=1.txt --check%c[0m\n", 0x1B, 0x1B))
	}
	flag.Parse()

	//读取帐号配置文件
	conf, err := lib.ReadConfFile(*config);
	if err != nil {
		log.Fatalf("配置文件读取失败,详细原因:%s", err.Error())
	}

	/*开启macVLan*/
	if err := lib.OpenMacVLanMode(*eth); err != nil {
		log.Fatalf("启用macvlan失败")
	}

	/*如果是检查帐号模式，仅作检查*/
	if *check {
		if err := lib.CheckAccounts(conf, *eth, *fix,*debug); err != nil {
			log.Fatalln(err.Error())
		}
		return
	}else{
		if err := lib.BatchDialing(conf, *eth,*interval,*flow, *fix,*debug); err != nil {
			log.Fatalln(err.Error())
		}
	}
}
