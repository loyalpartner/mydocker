package main

import (
	"log"
	"os"
	"os/exec"
	sc "syscall"
)

func main() {

	cmd := exec.Command("sh")

	cmd.SysProcAttr = &sc.SysProcAttr{
		Cloneflags: sc.CLONE_NEWUTS |
			sc.CLONE_NEWNS |
			sc.CLONE_NEWIPC |
			sc.CLONE_NEWUSER |
			sc.CLONE_NEWNET |
			sc.CLONE_NEWPID,
	}

	// cmd.SysProcAttr.Credential = &sc.Credential{
	// 	Uid: uint32(1),
	// 	Gid: uint32(1),
	// }
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	os.Exit(-1)

}
