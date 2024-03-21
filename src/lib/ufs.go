package lib

import (
	"fmt"
	"os/exec"
	"path"
)

type UfsOp int

const (
	UfsOpCreate UfsOp = 1
	UfsOpClean  UfsOp = 2
)

func UfsHandler(ns string, op UfsOp) error {
	scriptPath, err := GetScriptPath()
	if err != nil {
		return err
	}
	rootfsPath, err := GetRootFSPath()
	if err != nil {
		return err
	}

	var (
		rootfs       = rootfsPath
		rootfsUpper  = path.Join(rootfsPath, "../", fmt.Sprintf("%s_%s", ROOT_FS_UPPER, ns))
		rootfsWork   = path.Join(rootfsPath, "../", fmt.Sprintf("%s_%s", ROOT_FS_WORK, ns))
		rootfsMerged = path.Join(rootfsPath, "../", fmt.Sprintf("%s_%s", ROOT_FS_MERGED, ns))
		cmd          *exec.Cmd
	)

	switch op {
	case UfsOpCreate:
		cmd = CreateCMD(BASH, scriptPath, BASH_FUNC_CREATE_CONTAINER_UFS,
			rootfs, rootfsUpper, rootfsWork, rootfsMerged)
	case UfsOpClean:
		cmd = CreateCMD(BASH, scriptPath, BASH_FUNC_REMOVE_CONTAINER_UFS,
			rootfsUpper, rootfsWork, rootfsMerged)
	default:
		return fmt.Errorf("invalid ufs op %d", op)
	}

	if err := cmd.Run(); err != nil {
		return err
	}
	fmt.Printf("cmd: %v executed successfully\n", cmd)

	return nil
}
