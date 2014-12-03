package jobman

type job struct {
	name    string
	state   int
	total   int
	crawled int
	sample  string
	err     string
	birth   string
	death   string

	doms []string
}

func newJob(name string) *job {
	ret := new(job)
	ret.name = name
	return ret
}
