package imports

type RPC struct {
}

type Args_FetchPkgIndex struct {
}
type Reply_FetchPkgIndex struct {
	PkgIndexMap map[string][]Pkg
}

func (r *RPC) FetchPkgIndex(args *Args_FetchPkgIndex, reply *Reply_FetchPkgIndex) error {
	pkgIndex.Lock()
	reply.PkgIndexMap = pkgIndex.m
	pkgIndex.Unlock()
	return nil
}
func (c *rpcClient) fetchPkgIndex() {
	var args Args_FetchPkgIndex
	var reply Reply_FetchPkgIndex
	err := c.Call("RPC.FetchPkgIndex", &args, &reply)
	if err != nil {
		panic(err)
	}
	pkgIndex.Lock()
	pkgIndex.m = reply.PkgIndexMap
	pkgIndex.Unlock()
	return
}
