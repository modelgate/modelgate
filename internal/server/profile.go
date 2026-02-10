package server

import (
	"github.com/samber/do/v2"
)

type Profile struct {
	Name      string
	Addr      string
	Container do.Injector
}
