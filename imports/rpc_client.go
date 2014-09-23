package imports

import (
	"fmt"
	"net/rpc"
	"os"
	"time"
)

var gClient *rpcClient

type rpcClient struct {
	*rpc.Client
}

func newRpcClient() *rpcClient {
	addr := getSocketAddr()
	network := getNetwork()

	// client
	client, err := rpc.Dial(network, addr)
	if err != nil {
		if network == "unix" && fileExists(addr) {
			os.Remove(addr)
		}
		err = tryRunServer()
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			return nil
		}
		client, err = tryConnect(network, addr)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			return nil
		}
	}

	return &rpcClient{client}

}

func tryRunServer() error {
	path := getExecutableFilename()
	args := []string{os.Args[0], "-s"}
	procattr := os.ProcAttr{Dir: "", Env: os.Environ(), Files: []*os.File{}}
	p, err := os.StartProcess(path, args, &procattr)
	if err != nil {
		return err
	}
	return p.Release()
}

func tryConnect(network, address string) (client *rpc.Client, err error) {
	t := 0
	for {
		client, err = rpc.Dial(network, address)
		if err != nil && t < 1000 {
			time.Sleep(10 * time.Millisecond)
			t += 10
			continue
		}
		break
	}

	return
}
