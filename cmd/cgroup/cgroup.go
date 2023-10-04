package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	sc "syscall"
)

const cgroupMemoryHirarchyMount = "/sys/fs/cgroup/memory"

func main() {

	if os.Args[0] == "/proc/self/exe" {

		pid := sc.Getpid()

		fmt.Printf("current pid %d", pid)
		fmt.Println()

		command := `stress --vm-bytes 200m --vm-keep -m 1`
		cmd := exec.Command("sh", "-c", command)

		redirectFd(cmd)

		if err := cmd.Run(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	cmd := exec.Command("/proc/self/exe")
	cmd.SysProcAttr = &sc.SysProcAttr{
		Cloneflags: sc.CLONE_NEWUTS |
			sc.CLONE_NEWPID |
			sc.CLONE_NEWNS,
	}

	redirectFd(cmd)

	if err := cmd.Start(); err != nil {
		fmt.Println("ERROR", err)
		os.Exit(1)
	}

	cpid := cmd.Process.Pid // child process id
	fmt.Printf("pid is %v\n", sc.Getpid())
	fmt.Printf("child pid is %v", cpid)
	fmt.Println()

	d := path.Join(cgroupMemoryHirarchyMount, "testmemorylimit2")
	os.Mkdir(d, 0755)
	os.WriteFile(path.Join(d, "tasks"), []byte(strconv.Itoa(cpid)), 0644)
	os.WriteFile(path.Join(d, "memory.limit_in_bytes"), []byte("100m"), 0644)
	cmd.Process.Wait()
}

func redirectFd(cmd *exec.Cmd) {
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
}
