package container

import (
	"os"
	"os/exec"
	sc "syscall"
)

func NewParentProcess(tty bool, command string) *exec.Cmd {
	args := []string{"init", command}
	cmd := exec.Command("/proc/self/exe", args...)

	cmd.SysProcAttr = &sc.SysProcAttr{
		Cloneflags: sc.CLONE_NEWUTS |
			sc.CLONE_NEWNS |
			sc.CLONE_NEWIPC |
			sc.CLONE_NEWNET |
			sc.CLONE_NEWPID,
	}

	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	return cmd
}
