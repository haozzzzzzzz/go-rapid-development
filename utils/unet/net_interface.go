package unet

import (
	"github.com/sirupsen/logrus"
	"net"
	"strings"
)

func GetInterfaceIP(name string) (ip string, err error) {
	inter, err := net.InterfaceByName(name)
	if err != nil {
		logrus.Errorf("get interface failed. name: %s, error: %s", name, err)
		return
	}

	addrs, err := inter.Addrs()
	if err != nil {
		logrus.Errorf("get addresses failed. error: %s", err)
		return
	}

	if len(addrs) <= 0 {
		return
	}

	ipaddress := addrs[0].String()
	parts := strings.Split(ipaddress, "/")
	if len(parts) == 0 {
		return
	}

	ip = strings.TrimSpace(parts[0])

	return
}
