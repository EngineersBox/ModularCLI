package cli

import "flag"

type SubCommand struct {
	Flags        []*Flag
	Parameters   []*Parameter
	ErrorHandler flag.ErrorHandling
}
