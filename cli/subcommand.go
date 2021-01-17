package cli

import "flag"

type SubCommand struct {
	Arguments    []*Argument
	Parameters   []*Parameter
	ErrorHandler flag.ErrorHandling
}
