package lib

import "fmt"

const (
	BRIDGE            = "br0"
	ADDR_PREFIX       = "172.16.0"
	ADDR_SUFFIX_BEGIN = 2
)

const (
	IP       = "/usr/bin/ip"
	BASH     = "/bin/bash"
	IPTABLES = "/usr/sbin/iptables"
)

var (
	BRIDGE_ADDR   = fmt.Sprintf("%s.1/24", ADDR_PREFIX)
	DEFAULT_ROUTE = fmt.Sprintf("%s.1", ADDR_PREFIX)
)

const (
	SOCK_FILE      = "/opt/tmp/my_container.sock"
	ROOT_FS        = "rootfs"
	ROOT_FS_UPPER  = "rootfs_upper"
	ROOT_FS_WORK   = "rootfs_work"
	ROOT_FS_MERGED = "rootfs_merged"
)

// Bash script function
const (
	BASH_FUNC_CREATE_BRIDGE            = "create_bridge"
	BASH_FUNC_CREATE_CONTAINER_NETWORK = "create_container_network"
	BASH_FUNC_CREATE_CONTAINER_UFS     = "create_container_ufs"
	BASH_FUNC_REMOVE_CONTAINER_UFS     = "remove_container_ufs"
)
