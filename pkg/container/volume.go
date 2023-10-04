package container

import (
	"os"
	"os/exec"
	"path"

	log "github.com/sirupsen/logrus"
)

func NewWorkspace(root string, mnt string) {
	CreateReadonlyLayer(root)
	CreateWritableLayer(root)
	CreateMountPoint(root, mnt)
}

func CreateReadonlyLayer(root string) {
	busybox := path.Join(root, "busybox")
	tarfile := path.Join(root, "busybox.tar")

	exist, err := isExist(busybox)
	if err != nil {
		log.Infof("fail to judge whether dir %s exists. %v", busybox, err)
	}

	if !exist {
		if err := os.Mkdir(busybox, 0777); err != nil {
			log.Errorf("Mkdir dir %s error. %v", busybox, err)
		}

		cmd := exec.Command("tar", "-xvf", tarfile, "-C", busybox)
		if _, err := cmd.CombinedOutput(); err != nil {
			log.Errorf("Uncompress dir %s error %v", busybox, err)
		}
	}
}

func CreateWritableLayer(root string) {
	path := path.Join(root, "writeLayer/")
	if err := os.Mkdir(path, 0777); err != nil {
		log.Errorf("Mkdir dir %s error. %v", path, err)
	}
}

func CreateMountPoint(root string, mnt string) {
	if err := os.Mkdir(mnt, 0777); err != nil {
		log.Errorf("Mkdir dir %s error. %v", mnt, err)
	}

	// dirs := "dirs=" + root + "writeLayer:" + root + "busybox"
	// cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", mnt)
	// TODO: use btrfs filesystem
	cmd := exec.Command("mount", "-t", "tmpfs", "swap", mnt)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("%v", err)
	}
}

func isExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func DeleteWorkspace(root string, mnt string) {
	DeleteMountPoint(root, mnt)
	DeleteWritebaleLayer(root)
}

func DeleteMountPoint(root string, mnt string) {
	cmd := exec.Command("umount", mnt)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("%v", err)
	}
	if err := os.RemoveAll(mnt); err != nil {
		log.Errorf("Remove dir %s error %v", mnt, err)
	}
}

func DeleteWritebaleLayer(root string) {
	path := path.Join(root, "writeLayer")
	if err := os.RemoveAll(path); err != nil {
		log.Errorf("Remove dir %s error %v", path, err)
	}
}
