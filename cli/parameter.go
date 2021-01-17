package cli

import (
	"fmt"
	"os"
)

type Parameter struct {
	Type         ArgType
	Name         string
	Position     int
	Value        interface{}
	ValidateFunc ParamValidateFunc
}

type ParamValidateFunc func(Parameter) error

func (p Parameter) IsPresent() bool {
	return len(os.Args) > (p.Position + 1)
}

func (p Parameter) GetInt() *int {
	val, ok := (p.Value).(int)
	if !ok {
		return nil
	}
	return &val
}
func (p Parameter) GetBool() *bool {
	val, ok := (p.Value).(bool)
	if !ok {
		return nil
	}
	return &val
}
func (p Parameter) GetString() *string {
	val := fmt.Sprint(p.Value)
	return &val
}
