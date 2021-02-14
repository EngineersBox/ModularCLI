package cli

type FlagType int

const (
	TypeInvalid FlagType = iota
	TypeInt
	TypeBool
	TypeString
)
