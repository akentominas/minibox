package container

import (
	"fmt"
	"os"
	"os/exec"
)

func Run(cmdArgs []string) {
	if len(cmdArgs) == 0 {
		fmt.Println("No command specified")
		os.Exit(1)
	}

	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Running command:", cmdArgs)

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

}
