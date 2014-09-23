package imports

import (
	"fmt"
	"syscall"
	"unsafe"
)

var kernel32 = syscall.NewLazyDLL("kernel32.dll")
var proc_get_module_file_name = kernel32.NewProc("GetModuleFileNameW")

// Full path of the current executable
func getExecutableFilename() string {
	b := make([]uint16, syscall.MAX_PATH)

	ret, _, err := syscall.Syscall(proc_get_module_file_name.Addr(), 3,
		0, uintptr(unsafe.Pointer(&b[0])), uintptr(len(b)))
	if int(ret) == 0 {
		panic(fmt.Sprintf("GetModuleFileNameW : err %d", int(err)))
	}
	return syscall.UTF16ToString(b)
}

func getSocketAddr() string {
	return "localhost:41414"
}

func getNetwork() string {
	return "tcp"
}
