package main

import (
	"os"
	"strings"

	"github.com/loyalpartner/mydocker/pkg/cgroups"
	"github.com/loyalpartner/mydocker/pkg/cgroups/subsystems"
	"github.com/loyalpartner/mydocker/pkg/container"
	log "github.com/sirupsen/logrus"
)

func Run(tty bool, args []string, res *subsystems.ResourceConfig) {
	parent, writer := container.NewParentProcess(tty)
	if parent == nil {
		log.Errorf("New parent process error")
	}

	if err := parent.Start(); err != nil {
		log.Error(err)
	}

	cm := cgroups.NewCgroupManager("mydocker-cgroup")
	defer cm.Destroy()
	cm.Set(res)
	cm.Apply(parent.Process.Pid)

	sendInitCommand(args, writer)
	if err := parent.Wait(); err != nil {
		log.Error(err)
	}

}

func sendInitCommand(args []string, writer *os.File) {
	cmd := strings.Join(args, " ")
	log.Infof("command all is %s", cmd)
	writer.WriteString(cmd)
	writer.Close()
}
