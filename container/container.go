package container

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"minibox/debug"
	"minibox/nsutils"
)

/*
Run = PARENT PROCESS (runs on the host)

This function:
1. Validates input
2. Re-executes the SAME binary (/proc/self/exe)
3. Creates NEW namespaces for the child process

IMPORTANT:
This process itself is NOT inside the namespaces.
It asks the kernel to create a NEW process that is.
*/
func Run(cmdArgs []string) {
	if len(cmdArgs) == 0 {
		fmt.Println("No command specified")
		os.Exit(1)
	}

	debug.Info("RUN (parent)")

	/*
		Re-exec current binary:

		/proc/self/exe = path to currently running binary

		We pass:
		"child" → tells program to execute Child() next
		cmdArgs → actual command (e.g. bash)

		This creates:
		minibox (parent) → minibox (child in new namespaces)
	*/
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, cmdArgs...)...)

	// Connect terminal (so bash is interactive)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	/*
		Namespace creation happens HERE:

		These flags tell the kernel:
		"Create the new process in isolated environments"
	*/
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags:
		// UTS namespace → separate hostname
		syscall.CLONE_NEWUTS |

			// PID namespace → separate process tree
			syscall.CLONE_NEWPID |

			// Mount namespace → separate filesystem mounts
			syscall.CLONE_NEWNS,
	}

	// Start the child process
	if err := cmd.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

/*
Child = runs INSIDE the container

This process:
- is already inside new namespaces
- becomes PID 1 (init process inside container)
- sets up the environment
- finally replaces itself with the target command (bash)
*/
func Child(cmdArgs []string) {
	debug.Info("CHILD (before setup)")

	/*
		Make mount namespace PRIVATE

		Why:
		Without this, mount changes could propagate to the host

		MS_PRIVATE → no sharing with host
		MS_REC     → apply recursively
	*/
	if err := syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, ""); err != nil {
		fmt.Println("Error making mounts private:", err)
		os.Exit(1)
	}

	/*
		SAFETY CHECK (optional but recommended):

		Ensures we are NOT in the host mount namespace.
		If we are, we abort to avoid breaking the system.
	*/
	nsutils.MustBeInNewMountNS()

	/*
		Mount /proc filesystem

		Why:
		Tools like ps, top read from /proc

		Without this:
		ps would show host processes

		With this:
		ps shows ONLY container processes
	*/
	if err := syscall.Mount("proc", "/proc", "proc", 0, ""); err != nil {
		fmt.Println("Error mounting /proc:", err)
		os.Exit(1)
	}

	/*
		Set hostname (UTS namespace)

		This proves isolation:
		host hostname ≠ container hostname
	*/
	if err := syscall.Sethostname([]byte("minibox")); err != nil {
		fmt.Println("Error setting hostname:", err)
	}

	debug.Info("CHILD (before exec)")

	/*
		Find the binary to execute (e.g. /bin/bash)
	*/
	binary, err := exec.LookPath(cmdArgs[0])
	if err != nil {
		fmt.Println("Error finding binary:", err)
		os.Exit(1)
	}

	fmt.Println("Executing:", binary)

	/*
		CRITICAL STEP: Replace process using exec

		This does NOT create a new process.
		It REPLACES the current one.

		Result:
		minibox (PID 1) → becomes bash (PID 1)

		This is why bash becomes PID 1.
	*/
	if err := syscall.Exec(binary, cmdArgs, os.Environ()); err != nil {
		fmt.Println("Exec error:", err)
		os.Exit(1)
	}
}
