package cli

import (
	"os"
	"strings"
)

type Argument struct {
	Type         ArgType
	Name         string
	DefaultValue interface{}
	Required     bool
	HelpMsg      string
	ValidateFunc ArgValidateFunc
	Value        interface{}
	TypedArgument
}

type ArgValidateFunc func(TypedArgument) error

func (a Argument) IsPresent() bool {
	var present = false
	for _, arg := range os.Args {
		present = present || strings.Contains(arg, "--"+a.Name)
	}
	return present
}

type TypedArgument interface {
	GetString() *string
	GetBool() *bool
	GetInt() *int
	GetArg() *Argument
}

type StringArgument struct {
	Value *string
	*Argument
}

func (s StringArgument) GetString() *string { return s.Value }
func (StringArgument) GetBool() *bool       { return nil }
func (StringArgument) GetInt() *int         { return nil }
func (s StringArgument) GetArg() *Argument  { return s.Argument }

type IntArgument struct {
	Value *int
	*Argument
}

func (i IntArgument) GetInt() *int      { return i.Value }
func (IntArgument) GetBool() *bool      { return nil }
func (IntArgument) GetString() *string  { return nil }
func (i IntArgument) GetArg() *Argument { return i.Argument }

type BoolArgument struct {
	Value *bool
	*Argument
}

func (b BoolArgument) GetBool() *bool    { return b.Value }
func (BoolArgument) GetString() *string  { return nil }
func (BoolArgument) GetInt() *int        { return nil }
func (b BoolArgument) GetArg() *Argument { return b.Argument }
