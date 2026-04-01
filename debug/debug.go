package debug

import (
	"fmt"
	"os"
	"os/exec"
)

/*
Prints detailed information about the current process.

Used to understand:
- where we are running (host vs container)
- namespace IDs
- process hierarchy
*/
func Info(stage string) {
	fmt.Println("===================================")
	fmt.Println("STAGE:", stage)

	fmt.Printf("PID: %d\n", os.Getpid())
	fmt.Printf("PPID: %d\n", os.Getppid())

	exe, _ := os.Readlink("/proc/self/exe")
	fmt.Println("Executable:", exe)

	pidNs, _ := os.Readlink("/proc/self/ns/pid")
	utsNs, _ := os.Readlink("/proc/self/ns/uts")
	mntNs, _ := os.Readlink("/proc/self/ns/mnt")

	fmt.Println("Namespaces:")
	fmt.Println("  PID:", pidNs)
	fmt.Println("  UTS:", utsNs)
	fmt.Println("  MNT:", mntNs)

	fmt.Println("===================================")

	fmt.Println("Process tree:")
	cmd := exec.Command("ps", "-o", "pid,ppid,comm")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	fmt.Println("===================================")
}
