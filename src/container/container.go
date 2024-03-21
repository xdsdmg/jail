/* Container module */

package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path"
	"strconv"
	"strings"
	"syscall"

	"github.com/xdsdmg/jail/lib"
	"golang.org/x/sys/unix"
)

const (
	ChildFlag = "child"
	Flag      = syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS /* Ref: https://lwn.net/Articles/531114/ */
)

var networkConfig *lib.NetworkConfig

func printArg(args []string) {
	for i, arg := range args {
		fmt.Printf("[INFO] [printArg] index: %d, arg: %s\n", i, arg)
	}
}

// getNetworkConfig gets network configuration from server through UNIX domain socket.
func getNetworkConfig() (networkConfig *lib.NetworkConfig, err error) {
	c, err := net.Dial("unix", lib.SOCK_FILE)
	if err != nil {
		return nil, err
	}

	_, err = c.Write([]byte("hi\n"))
	if err != nil {
		return nil, err
	}

	buf := make([]byte, 512)
	n, err := c.Read(buf)
	if err != nil {
		return nil, err
	}

	data := buf[0:n]
	fmt.Printf("[INFO] [getNetworkConfig] received: %s, len: %dByte\n", string(data), n)

	err = json.Unmarshal(data, &networkConfig)
	if err != nil {
		return nil, err
	}

	return networkConfig, nil
}

// getWorkDirPath gets the path of work directory for container.
func getWorkDirPath() (rootfsPath string, err error) {
	ip, err := lib.GetIPv4Addr()
	if err != nil {
		return "", err
	}
	segs := strings.Split(ip, ".")
	if len(segs) != 4 {
		return "", fmt.Errorf("invalid ip %s", ip)
	}
	num, err := strconv.ParseInt(segs[3], 10, 32)
	if err != nil {
		return "", err
	}
	p, err := lib.GetRootFSPath()
	if err != nil {
		return "", err
	}

	return path.Join(p, "../", fmt.Sprintf("rootfs_merged_ns%d", num-2)), nil
}

// The logic of parent process.
func parent() {
	nc, err := getNetworkConfig()
	if err != nil {
		panic(fmt.Errorf("get net config failed, error: %+v", err))
	}
	networkConfig = nc

	/* Create union file system */
	err = lib.UfsHandler(nc.NS, lib.UfsOpCreate)
	if err != nil {
		panic(fmt.Errorf("create ufs failed, error: %+v", err))
	}

	/* Config container network */
	err = lib.CreateContainerNetwork(nc)
	if err != nil {
		panic(fmt.Errorf("create veth failed, error: %+v", err))
	}

	/* Set network namespace */
	netnsPath := fmt.Sprintf("/var/run/netns/%s", networkConfig.NS)
	fd, err := os.OpenFile(netnsPath, os.O_RDONLY, 0755)
	if err != nil {
		panic(fmt.Errorf("open file failed, error: %+v", err))
	}
	unix.Setns(int(fd.Fd()), 0)

	/* Create child process */
	var (
		binaryName = os.Args[0]
		childTask  = lib.BASH // The command will be executed by the child process
	)
	cmd := lib.CreateCMD(binaryName, ChildFlag, childTask)
	cmd.SysProcAttr = &syscall.SysProcAttr{Cloneflags: Flag}
	if err := cmd.Run(); err != nil {
		panic(fmt.Errorf("cmd %v run failed, error: %+v", cmd, err))
	}
}

// The logic of child process.
//
// os.Args[0]: the binary name
// os.Args[1]: "child"
// os.Args[2]: the task will be executed by the child process
func child() {
	rootfsPath, err := getWorkDirPath()
	if err != nil {
		panic(fmt.Sprintf("get rootfs path failed, error: %+v\n", err))
	}

	syscall.Sethostname([]byte("container"))

	/* chroot */
	err = syscall.Chroot(rootfsPath)
	if err != nil {
		panic(fmt.Sprintf("chroot failed, error: %v\n", err))
	}

	/* Change current dir to "/" */
	err = syscall.Chdir("/")
	if err != nil {
		panic(fmt.Sprintf("chdir failed, error: %v\n", err))
	}

	/* Mount proc fs, tools like ps, top, etc will read /proc */
	err = syscall.Mount("proc", "/proc", "proc",
		uintptr(syscall.MS_NOEXEC|syscall.MS_NOSUID|syscall.MS_NODEV), "")
	if err != nil {
		panic(fmt.Sprintf("mount proc failed, error: %v\n", err))
	}

	/* Execute task */
	if len(os.Args) < 3 {
		panic("no task specified by the parent process")
	}
	task := os.Args[2]
	err = syscall.Exec(task, nil, os.Environ())
	if err != nil {
		panic(fmt.Sprintf("child cmd %s run failed, error: %v\n", lib.BASH, err))
	}
}

func main() {
	fmt.Printf("[INFO] begin pid: %d\n", os.Getpid())

	if !(len(os.Args) > 2 && os.Args[1] == ChildFlag) {
		parent()
	} else {
		child()
	}

	if networkConfig == nil {
		panic("network config not found")
	}
	/* Remove Network config */
	if err := lib.RemoveNetworkNamespace(networkConfig.NS); err != nil {
		panic(fmt.Errorf("clean network config failed, error: %+v", err))
	}
	/* Remove ufs */
	if err := lib.UfsHandler(networkConfig.NS, lib.UfsOpClean); err != nil {
		panic(fmt.Errorf("clean ufs failed, error: %+v", err))
	}

	fmt.Printf("[INFO] end pid: %d\n", os.Getpid())
}
