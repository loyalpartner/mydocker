package container

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"
)

func RunContainerInitProcess() error {
	args := readArgs()
	if args == nil || len(args) == 0 {
		return fmt.Errorf("Run container get user command error, command is nil")
	}

	setupMount()

	path, err := exec.LookPath(args[0])
	if err != nil {
		log.Errorf("Exec loop path err %v", err)
	}

	log.Infof("Find path %s", path)

	if err := syscall.Exec(path, args, os.Environ()); err != nil {
		log.Errorf(err.Error())
	}

	return nil
}

func readArgs() []string {
	pipe := os.NewFile(uintptr(3), "pipe")
	msg, err := io.ReadAll(pipe)
	if err != nil {
		log.Errorf("init read pipe error %v", err)
		return nil
	}

	return strings.Split(string(msg), " ")
}

func setupMount() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Errorf("Get current location error %v", err)
	}

	log.Infof("Current location is %s", pwd)

	// TODO: pivoRoot

	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")

	syscall.Mount("tmpfs", "/dev", "tmpfs", syscall.MS_NOSUID|syscall.MS_STRICTATIME, "mode=755")

}
