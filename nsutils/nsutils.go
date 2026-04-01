package nsutils

import (
	"fmt"
	"os"
)

/*
Ensures we are inside a NEW mount namespace.

If not, we abort to avoid modifying the host system.
*/
func MustBeInNewMountNS() {
	self, _ := os.Readlink("/proc/self/ns/mnt")
	init, _ := os.Readlink("/proc/1/ns/mnt")

	fmt.Println("Self MNT NS:", self)
	fmt.Println("Init MNT NS:", init)

	if self == init {
		fmt.Println("❌ Refusing to mount /proc: still in host namespace")
		os.Exit(1)
	}

	fmt.Println("✅ Safe: running in isolated mount namespace")
}
