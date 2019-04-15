package db

import "github.com/hzxiao/goutil"

type Status struct {
	ID           int
	Name         string
	UpdateHeight int
}

type Block struct {
	Height      int
	SysFee      int
	TotalSysFee int
	Raw         goutil.Map
}
