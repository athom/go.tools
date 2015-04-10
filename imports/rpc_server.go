package imports

import (
	"log"
	"net"
	"net/rpc"
	"os"
	"runtime"
)

const (
	DAEMON_CLOSE = iota
)

var g_daemon *daemon

type daemon struct {
	listener net.Listener
	cmd_in   chan int
}

func newDaemon(network string, address string) *daemon {
	var err error

	d := new(daemon)
	d.listener, err = net.Listen(network, address)
	if err != nil {
		panic(err)
	}

	d.cmd_in = make(chan int, 1)
	pkgIndexOnce.Do(loadPkgIndex)
	return d
}

func RunServer() int {
	addr := getSocketAddr()
	network := getNetwork()
	if network == "unix" {
		if fileExists(addr) {
			log.Printf("unix socket: '%s' already exists\n", addr)
			return 1
		}
	}

	g_daemon = newDaemon(network, addr)
	// cleanup unix socket file
	if network == "unix" {
		defer os.Remove(addr)
	}

	rpc.Register(new(RPC))

	// scan GOPATH for changes
	//fm := newFilesMonitor()
	//fm.Monitor()

	// serv rpc
	g_daemon.loop()

	return 0
}

func (this *daemon) loop() {
	conn_in := make(chan net.Conn)
	go func() {
		for {
			c, err := this.listener.Accept()
			if err != nil {
				panic(err)
			}
			conn_in <- c
		}
	}()
	for {
		// handle connections or server CMDs (currently one CMD)
		select {
		case c := <-conn_in:
			rpc.ServeConn(c)
			runtime.GC()
		case cmd := <-this.cmd_in:
			switch cmd {
			case DAEMON_CLOSE:
				return
			}
		}
	}
}

func (this *daemon) close() {
	this.cmd_in <- DAEMON_CLOSE
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return true
}
