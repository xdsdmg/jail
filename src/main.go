// The logic of main.go:
//
//	After executing `./jail` command，jail will create a process
//	with clone flags to execute jail again.

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"syscall"
)

const (
	ChildFlag  = "child"
	CMD        = "/bin/bash"
	Flag       = syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET
	RootfsName = "rootfs_merged"
)

func createCMD(args []string) (cmd *exec.Cmd) {
	if len(args) == 0 {
		panic("len of args is zero")
	} else if len(args) < 1 {
		cmd = exec.Command(args[0])
	} else {
		cmd = exec.Command(args[0], args[1:]...)
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}

func runCMD(tip, script string) error {
	if len(tip) <= 0 {
		fmt.Println(tip)
	}

	cmd := exec.Command("/bin/bash", "-cx", script)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func getRootfsPath() (rootfsPath string, err error) {
	workDirPath, err := os.Getwd()
	if err != nil {
		return "", err
	}
	rootfsPath = path.Join(workDirPath, os.Args[0], fmt.Sprintf("../../%s", RootfsName))

	return rootfsPath, nil
}

func printArg(args []string) {
	for i, arg := range args {
		fmt.Printf("[INFO] index: %d, arg: %s\n", i, arg)
	}
}

func main() {
	printArg(os.Args)

	fmt.Printf("[INFO] begin pid: %d\n", os.Getpid())

	if len(os.Args) > 1 && os.Args[1] == ChildFlag {
		syscall.Sethostname([]byte("container"))

		/* chroot */
		rootfsPath, err := getRootfsPath()
		if err != nil {
			fmt.Printf("[ERROR] get rootfs path failed, error: %v\n", err)
			return
		}
		err = syscall.Chroot(rootfsPath)
		if err != nil {
			fmt.Printf("[ERROR] chroot failed, error: %v\n", err)
			return
		}

		/* Change current dir to "/" */
		err = syscall.Chdir("/")
		if err != nil {
			fmt.Printf("[ERROR] chdir failed, error: %v\n", err)
			return
		}

		/* Mount proc fs, tools like ps, top, etc will read /proc */
		err = syscall.Mount("proc", "/proc", "proc",
			uintptr(syscall.MS_NOEXEC|syscall.MS_NOSUID|syscall.MS_NODEV), "")
		if err != nil {
			fmt.Printf("[ERROR] mount proc failed, error: %v\n", err)
			return
		}

		/* Exec bash */
		err = syscall.Exec(CMD, nil, os.Environ())
		if err != nil {
			fmt.Printf("[ERROR] child cmd %s run failed, error: %v\n", CMD, err)
			return
		}
	} else {
		cmd := createCMD([]string{os.Args[0], ChildFlag, CMD})
		cmd.SysProcAttr = &syscall.SysProcAttr{Cloneflags: Flag}
		if err := cmd.Run(); err != nil {
			fmt.Printf("[ERROR] parent cmd %s run failed, error: %v\n", CMD, err)
			return
		}
	}

	fmt.Printf("[INFO] end pid: %d\n", os.Getpid())
}
