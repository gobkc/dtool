package lib

import (
	"errors"
	"github.com/vishvananda/netlink"
	"net"
	"os/exec"
)

func AddMacVLan(devName string, parentDevName string, macAddr string) error {
	var err error

	/*根据设备名找到设备*/
	link, err := netlink.LinkByName(parentDevName)
	if err != nil {
		return err
	}

	mac, err := net.ParseMAC(macAddr)
	if err != nil {
		return errors.New("在添加macvlan时无法解析mac地址：" + err.Error())
	}

	mv := netlink.NewLinkAttrs()
	mv.Name = devName
	mv.ParentIndex = link.Attrs().Index
	mv.HardwareAddr = mac

	macVLan := netlink.Macvlan{
		LinkAttrs: mv,
		Mode:      netlink.MACVLAN_MODE_BRIDGE,
	}

	/*先删除同名的macVLan，防止出错*/
	netlink.LinkDel(&macVLan)

	/*添加macVLan*/
	if err = netlink.LinkAdd(&macVLan); err != nil {
		return errors.New("在添加macvlan时报错:" + err.Error())
	}

	err = netlink.LinkSetDown(&macVLan)
	if err != nil {
		return err
	}

	/*启用macVLan*/
	if err = netlink.LinkSetUp(&macVLan); err != nil {
		return errors.New("在启用macvlan时报错:" + devName + " " + err.Error())
	}

	return nil
}

func BatchSetMacVLanAndPromisc(cmdString string) error {
	cmd := exec.Command("bash", "-c", "modprobe macvlan && "+cmdString)
	if _, err := cmd.Output(); err != nil {
		return errors.New("启用macvlan和promisc失败,详细原因:" + err.Error())
	}
	return nil
}

func OpenMacVLanMode(eth string) error {
	cmdString := "ifconfig " + eth + " promisc"
	if err := BatchSetMacVLanAndPromisc(cmdString); err != nil {
		return err
	}
	return nil
}

func KillByName(name string) error {
	cli := exec.Command("bash", "-c", "ps -ef|grep "+name+"|grep -v grep|cut -c 9-15|xargs kill -9")
	if _, err := cli.Output(); err != nil {
		return errors.New(name + "已被关闭")
	}
	return nil
}
