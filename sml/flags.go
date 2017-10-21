package sml

import (
	"flag"
)

func newFlags() *flag.FlagSet {
	return flag.NewFlagSet("sml", flag.ExitOnError)
}
