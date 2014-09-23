package imports

import (
	"go/build"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type scanCallback func(path string)

// filesmMonitor aims to update packages info to the daemon server
// currently only new created pkg will be updated to the map of pkgIndex
type filesMonitor struct {
	sync.Mutex
	pathTimeMap map[string]time.Time
}

func newFilesMonitor() *filesMonitor {
	r := &filesMonitor{}
	r.pathTimeMap = make(map[string]time.Time)
	return r
}

var newCreatedCallback = func(path string) {
	dir := filepath.Dir(path)
	buildPkg, err := build.ImportDir(dir, 0)
	//skip invalid go package
	if err != nil {
		return
	}

	shortName := buildPkg.Name
	dir = buildPkg.Dir
	importpath := buildPkg.ImportPath

	pkgIndex.Lock()
	defer pkgIndex.Unlock()
	pkgs, ok := pkgIndex.m[shortName]
	if ok {
		for _, pkg := range pkgs {
			if pkg.Importpath == buildPkg.ImportPath && pkg.Dir == dir {
				return
			}
		}
	}
	pkgIndex.m[shortName] = append(
		pkgIndex.m[shortName],
		Pkg{
			Importpath: importpath,
			Dir:        dir,
		},
	)
}

func (this *filesMonitor) Monitor() {
	go func() {
		srcDir := filepath.Join(build.Default.GOPATH, "src")
		this.scanChanges(srcDir, newCreatedCallback)
	}()
	return
}

// TODO more callback events like modifed and deleted
func (this *filesMonitor) walker(loop bool, createdCallback scanCallback) func(path string, info os.FileInfo, err error) error {
	return func(path string, info os.FileInfo, err error) error {
		if info == nil || info.IsDir() {
			return nil
		}

		// ignore non go files
		if filepath.Ext(path) != ".go" {
			return nil
		}

		this.Lock()
		_, ok := this.pathTimeMap[path]
		if loop && !ok {
			createdCallback(path)
		}
		this.pathTimeMap[path] = info.ModTime()
		this.Unlock()
		return nil
	}
}

func (this *filesMonitor) scanChanges(watchPath string, createdCallback scanCallback) {
	// need to walk first round so that the walker in the loop can detect new created files
	filepath.Walk(watchPath, this.walker(false, createdCallback))
	for {
		filepath.Walk(watchPath, this.walker(true, createdCallback))
		time.Sleep(500 * time.Millisecond)
	}
}
