package lib

import (
	"errors"
	"fmt"
	"github.com/vishvananda/netlink"
	"net"
	"os/exec"
)

/*批量删除策略 根据策略表ID*/
func DelRuleByTableId(tableID int) {
	rule := netlink.NewRule()
	rule.Table = tableID
	for {
		if err := netlink.RuleDel(rule); err != nil {
			break
		}
	}
}

/*添加默认路由*/
func AddDefaultRoute(ppp string, gatewayIp string, tableId int) error {
	_, inet, _ := net.ParseCIDR(fmt.Sprintf("%s/32", gatewayIp))
	_, cidr, _ := net.ParseCIDR("0.0.0.0/0")

	peerLink, err := netlink.LinkByName(ppp)
	if err != nil {
		return errors.New("在获取网卡设备时发生错误:" + err.Error())
	}
	defaultRoute := &netlink.Route{
		LinkIndex: peerLink.Attrs().Index,
		Gw:        inet.IP,
		Dst:       cidr,
		Table:     tableId,
	}
	if err := netlink.RouteAdd(defaultRoute); err != nil {
		return errors.New("在添加默认路由时发生错误:" + err.Error())
	}
	return nil
}

/*添加策略*/
func AddRule(src string, tableID int, fref int) error {
	var err error

	/*分割字符串IP地址.使之符合要求*/
	_, srcIpNet, err := net.ParseCIDR(fmt.Sprintf("%s/%s", src, "32"))
	if err != nil {
		return errors.New("添加策略时，因分割字符串转ipNet时发生错误")
	}

	rule := netlink.NewRule()
	rule.Table = tableID
	rule.Src = srcIpNet
	rule.Priority = fref
	rule.Invert = false

	/*添加rule*/
	if err = netlink.RuleAdd(rule); err != nil {
		return errors.New("在添加策略时报错:" + err.Error())
	}

	return nil
}

/*获取默认路由名称*/
func GetDefaultRouterName() string {
	cli := exec.Command("bash", "-c", "ip route list|grep 'default dev'|awk '{print($3)}'")
	var name string
	if s, err := cli.Output(); err != nil {
		cli = exec.Command("bash", "-c", "ip route list|grep 'default via '|awk '{print($5)}'")
		s, err = cli.Output()
		name = string(s)
	} else {
		name = string(s)
	}
	return name
}

/*删除默认路由*/
func DelDefaultRouter(name string) error {
	cli := exec.Command("bash", "-c", "ip route del default dev "+name)
	if _, err := cli.Output(); err != nil {
		return errors.New(name + "已被关闭")
	}
	return nil
}

/*设置默认路由*/
func SetDefaultRouter(ppp string) error {
	routeName := GetDefaultRouterName()
	if routeName != "" {
		DelDefaultRouter(routeName)
	}
	cli := exec.Command("bash", "-c", "ip route add default dev "+ppp)
	if _, err := cli.Output(); err != nil {
		return err
	}
	return nil
}