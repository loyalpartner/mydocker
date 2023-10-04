package main

import (
	"os"

	"github.com/loyalpartner/mydocker/pkg/container"
	log "github.com/sirupsen/logrus"
)

func Run(tty bool, command string) {
	parent := container.NewParentProcess(tty, command)
	if err := parent.Start(); err != nil {
		log.Error(err)
	}

	if err := parent.Wait(); err != nil {
		log.Error(err)
	}

	os.Exit(-1)

}
