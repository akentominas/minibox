package container

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func Run(cmdArgs []string) {
	if len(cmdArgs) == 0 {
		fmt.Println("No command specified")
		os.Exit(1)
	}

	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, cmdArgs...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | // new hostname namespace
			    syscall.CLONE_NEWPID | // new PID namespace
			    syscall.CLONE_NEWNS, // new mount namespace - filesystem isolation
	}

	fmt.Println(">>> RUN parent")
	
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

}

func Child(cmdArgs []string) {

	fmt.Println("Inside container")

	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	
	}

	fmt.Println(">>> RUN child")
}
