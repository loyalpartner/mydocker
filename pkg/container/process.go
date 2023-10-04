package container

import (
	"os"
	"os/exec"
	sc "syscall"

	"github.com/sirupsen/logrus"
)

func NewParentProcess(tty bool) (*exec.Cmd, *os.File) {
	reader, writer, err := os.Pipe()
	if err != nil {
		logrus.Errorf("New pipe error %v", err)
		return nil, nil
	}

	cmd := exec.Command("/proc/self/exe", "init")
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
	cmd.ExtraFiles = []*os.File{reader}
	return cmd, writer
}
