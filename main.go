// main
package main

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
	"os"
	"runtime"
	"sync"
)

const (
	LFATAL = iota
	LERROR
	LWARNING
	LINFO
	LDEBUG
)

const (
	Version = "2.1-dev"
)

var (
	listenAddr, forwardAddr strSlice
	pv                      bool // print version

	listenArgs  []Args
	forwardArgs []Args
)

func init() {
	flag.Var(&listenAddr, "L", "listen address, can listen on multiple ports")
	flag.Var(&forwardAddr, "F", "forward address, can make a forward chain")
	flag.BoolVar(&pv, "V", false, "print version")
	flag.Parse()
}

func main() {
	defer glog.Flush()

	if flag.NFlag() == 0 {
		flag.PrintDefaults()
		return
	}
	if pv {
		fmt.Fprintf(os.Stderr, "gost %s (%s)\n", Version, runtime.Version())
		return
	}

	listenArgs = parseArgs(listenAddr)
	forwardArgs = parseArgs(forwardAddr)

	if len(listenArgs) == 0 {
		glog.Exitln("no listen addr")
	}

	var wg sync.WaitGroup
	for _, args := range listenArgs {
		wg.Add(1)
		go func(arg Args) {
			defer wg.Done()
			glog.V(LERROR).Infoln(listenAndServe(arg))
		}(args)
	}
	wg.Wait()
}
