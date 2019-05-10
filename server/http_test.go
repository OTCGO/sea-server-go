package server

import (
	"github.com/OTCGO/sea-server-go/pkg/httptest"
)

func init()  {
	Init()
	Run()

	httptest.GinEngine = g
}