package lib

import (
	"errors"
	"github.com/vishvananda/netlink"
	"net"
)

/*添加DNS*/
func AddDns(src string, tableID int, fref int) error {
	var err error

	/*分割字符串IP地址.使之符合要求*/
	_, srcIpNet, err := net.ParseCIDR(src)
	if err != nil {
		return errors.New("添加DNS时，因分割字符串转ipNet时发生错误")
	}

	rule := netlink.NewRule()
	rule.Table = tableID
	rule.Dst = srcIpNet
	rule.Priority = fref
	rule.Invert = false

	/*添加dns*/
	if err = netlink.RuleAdd(rule); err != nil {
		return errors.New("在添加DNS时报错:" + err.Error())
	}

	return nil
}