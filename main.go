package main

import (
	"github.com/hzxiao/goutil/version"
	"github.com/spf13/pflag"
)

var (
	ver = pflag.BoolP("version", "v", false, "show version info.")
)

func main() {
	pflag.Parse()
	if *ver {
		err := version.Print()
		if err != nil {
			panic(err)
		}
		return
	}
}
