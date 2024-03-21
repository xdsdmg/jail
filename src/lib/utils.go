package lib

import (
	"os"
	"os/exec"
	"path"
)

func CreateCMD(name string, args ...string) (cmd *exec.Cmd) {
	cmd = exec.Command(name, args...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}

func GetScriptPath() (scriptPath string, err error) {
	binPath, err := os.Executable()
	if err != nil {
		return "", err
	}
	return path.Join(binPath, "../", "scripts/start.sh"), nil
}

func GetRootFSPath() (rootfsPath string, err error) {
	binPath, err := os.Executable()
	if err != nil {
		return "", err
	}

	return path.Join(binPath, "../", ROOT_FS), nil
}
