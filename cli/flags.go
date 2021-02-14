package cli

import (
	"os"
	"strings"
)

type Flag struct {
	Type         FlagType
	Name         string
	DefaultValue interface{}
	Required     bool
	HelpMsg      string
	ValidateFunc FlagValidateFunc
	Value        interface{}
	TypedArgument
}

type FlagValidateFunc func(TypedArgument) error

func (a Flag) IsPresent() bool {
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
	GetArg() *Flag
}

type StringArgument struct {
	Value *string
	*Flag
}

func (s StringArgument) GetString() *string { return s.Value }
func (StringArgument) GetBool() *bool       { return nil }
func (StringArgument) GetInt() *int         { return nil }
func (s StringArgument) GetArg() *Flag      { return s.Flag }

type IntArgument struct {
	Value *int
	*Flag
}

func (i IntArgument) GetInt() *int     { return i.Value }
func (IntArgument) GetBool() *bool     { return nil }
func (IntArgument) GetString() *string { return nil }
func (i IntArgument) GetArg() *Flag    { return i.Flag }

type BoolArgument struct {
	Value *bool
	*Flag
}

func (b BoolArgument) GetBool() *bool   { return b.Value }
func (BoolArgument) GetString() *string { return nil }
func (BoolArgument) GetInt() *int       { return nil }
func (b BoolArgument) GetArg() *Flag    { return b.Flag }
