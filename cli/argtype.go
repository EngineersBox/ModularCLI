package cli

type ArgType int

const (
	TypeInvalid ArgType = iota
	TypeInt
	TypeBool
	TypeString
)
