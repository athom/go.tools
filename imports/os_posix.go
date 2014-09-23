package imports

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Full path of the current executable
func getExecutableFilename() string {
	// try readlink first
	path, err := os.Readlink("/proc/self/exe")
	if err == nil {
		return path
	}
	// use argv[0]
	path = os.Args[0]
	if !filepath.IsAbs(path) {
		cwd, _ := os.Getwd()
		path = filepath.Join(cwd, path)
	}
	if fileExists(path) {
		return path
	}
	// Fallback : use "gocode" and assume we are in the PATH...
	path, err = exec.LookPath("goimports")
	if err == nil {
		return path
	}
	return ""
}

func getSocketAddr() string {
	user := os.Getenv("USER")
	if user == "" {
		user = "all"
	}
	return filepath.Join(os.TempDir(), fmt.Sprintf("goimports-daemon.%s", user))
}

func getNetwork() string {
	return "unix"
}
