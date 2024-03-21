package lib

import (
	"fmt"
	"net"
)

func GetIPv4Addr() (addr string, err error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", fmt.Errorf("ip address not found")
}

func CreateBridge(bridgeName, bridgeIPv4Addr string) error {
	scriptPath, err := GetScriptPath()
	if err != nil {
		return err
	}
	cmd := CreateCMD(BASH, scriptPath, BASH_FUNC_CREATE_BRIDGE, bridgeName, bridgeIPv4Addr)
	if err := cmd.Run(); err != nil {
		return err
	}
	fmt.Printf("cmd: %v executed successfully\n", cmd)

	return nil
}

func RemoveBridge(bridgeName string) error {
	cmd := CreateCMD(IP, "link", "delete", bridgeName)
	if err := cmd.Run(); err != nil {
		return err
	}
	fmt.Printf("cmd: %v executed successfully\n", cmd)

	return nil
}

func CreateContainerNetwork(nc *NetworkConfig) error {
	scriptPath, err := GetScriptPath()
	if err != nil {
		return err
	}
	vethPeerName := fmt.Sprintf("%s-br", nc.VethName)
	cmd := CreateCMD(BASH, scriptPath, BASH_FUNC_CREATE_CONTAINER_NETWORK,
		nc.VethName, nc.VethAddr, vethPeerName, nc.NS, nc.BridgeName, nc.DefaultRoute)
	if err := cmd.Run(); err != nil {
		return err
	}
	fmt.Printf("cmd: %v executed successfully\n", cmd)

	return nil
}

func RemoveNetworkNamespace(ns string) error {
	cmd := CreateCMD("ip", "netns", "del", ns)
	if err := cmd.Run(); err != nil {
		return err
	}
	fmt.Printf("cmd: %v executed successfully\n", cmd)
	return nil
}
