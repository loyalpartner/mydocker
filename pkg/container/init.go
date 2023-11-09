package container

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"
	sc "syscall"

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

	if err := sc.Exec(path, args, os.Environ()); err != nil {
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
	if err := pivotRoot(pwd); err != nil {
		log.Errorf("pivotRoot %v", err)
	}

	defaultMountFlags := sc.MS_NOEXEC | sc.MS_NOSUID | sc.MS_NODEV
	sc.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")

	sc.Mount("tmpfs", "/dev", "tmpfs", sc.MS_NOSUID|sc.MS_STRICTATIME, "mode=755")

}

func pivotRoot(root string) error {
	err := sc.Mount("", "/", "bind", syscall.MS_PRIVATE|syscall.MS_REC, "")
	if err != nil {
		return fmt.Errorf("Mount rootfs to itself error: %v", err)
	}

	err = sc.Mount(root, root, "bind", sc.MS_BIND|sc.MS_REC, "")
	if err != nil {
		return fmt.Errorf("Mount rootfs to itself error: %v", err)
	}

	pivotDir := path.Join(root, ".pivot_root")
	if err := os.Mkdir(pivotDir, 0777); err != nil {
		return err
	}

	if err := sc.PivotRoot(root, pivotDir); err != nil {
		return fmt.Errorf("pivot_root %v", err)
	}

	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("chdir / %v", err)
	}

	pivotDir = path.Join("/", ".pivot_root")
	if err := sc.Unmount(pivotDir, sc.MNT_DETACH); err != nil {
		return fmt.Errorf("unmount pivot_root dir %v", err)
	}

	return os.RemoveAll(pivotDir)
}
