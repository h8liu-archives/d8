package jobman

type worker struct {
	addr string
}

func newWorker(addr string) *worker {
	ret := new(worker)
	ret.addr = addr
	return ret
}
