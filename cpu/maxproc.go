package cpu

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"runtime"
)

func SetMaxProc() int {
	var err error
	var p, l int

	p, err = cpu.Counts(false)
	if err != nil {
		p = -1
	}

	l, err = cpu.Counts(true)
	if err != nil {
		l = -1
	}

	if p > l {
		runtime.GOMAXPROCS(p)
		return p
	}

	if p < l {
		runtime.GOMAXPROCS(l)
		return l
	}

	runtime.GOMAXPROCS(1)
	return 1
}
