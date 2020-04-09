package lib

import (
	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
	"log"
	"net"
)

/*设置PPPOE*/
func SetPPPOE(eth string, ruleTableID int) error {
	//未考虑重复规则情况
	link, err := netlink.LinkByName(eth)
	if err != nil {
		return err
	}
	addrList, err := netlink.AddrList(link, unix.AF_INET)
	for _, addr := range addrList {
		/*addr.Peer获取IP+网关*/
		log.Println(addr.Peer)
		/*addr.IP 只获取IP*/
		log.Println(addr.IP)
		route := netlink.Route{
			LinkIndex: link.Attrs().Index,
			Dst:       addr.Peer,
			Src:       addr.IP,
			Table:     ruleTableID,
			Scope:     unix.RT_SCOPE_LINK,
			Protocol:  unix.RTPROT_KERNEL,
		}
		// 添加路由
		if err := netlink.RouteAdd(&route); err != nil {
			log.Println(err)
			return err
		}
		// 添加策略
		rule := netlink.NewRule()
		_, srcnet, _ := net.ParseCIDR(addr.IP.String())
		rule.Src = srcnet
		rule.Dst = addr.Peer
		rule.Table = ruleTableID
		netlink.RuleAdd(rule)
	}
	return nil
}