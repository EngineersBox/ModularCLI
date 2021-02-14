package cli

import (
	"flag"
	"fmt"
	"os"
)

type Command struct {
	FlagSet *flag.FlagSet
	Flags   map[string]TypedArgument
	Params  map[string]*Parameter
}

type CLI struct {
	Commands map[string]*Command
}

func (c *CLI) AddCommand(name string, handler flag.ErrorHandling) {
	c.Commands[name] = &Command{
		FlagSet: flag.NewFlagSet(name, handler),
		Flags:   make(map[string]TypedArgument),
		Params:  make(map[string]*Parameter),
	}
}

func (c *CLI) AddCommandArgs(cmdName string, flag *Flag) error {
	var currentCmd = c.Commands[cmdName]
	switch flag.Type {
	case TypeString:
		currentCmd.Flags[flag.Name] = StringArgument{
			Value: currentCmd.FlagSet.String(flag.Name, flag.DefaultValue.(string), flag.HelpMsg),
			Flag:  flag,
		}
	case TypeBool:
		currentCmd.Flags[flag.Name] = BoolArgument{
			Value: currentCmd.FlagSet.Bool(flag.Name, flag.DefaultValue.(bool), flag.HelpMsg),
			Flag:  flag,
		}
	case TypeInt:
		currentCmd.Flags[flag.Name] = IntArgument{
			Value: currentCmd.FlagSet.Int(flag.Name, flag.DefaultValue.(int), flag.HelpMsg),
			Flag:  flag,
		}
	case TypeInvalid:
		return fmt.Errorf("invalid arugment type: %d", flag.Type)
	}
	return nil
}

func (c *CLI) AddCommandParameters(cmdName string, param *Parameter) error {
	var currentCmd = c.Commands[cmdName]
	if param.Type == TypeInvalid {
		return fmt.Errorf("invalid parameter type: %d", param.Type)
	}
	currentCmd.Params[param.Name] = param
	return nil
}

func validateArguments(arguments map[string]TypedArgument) error {
	for _, arg := range arguments {
		argument := arg.GetArg()
		if argument.Required && !argument.IsPresent() {
			return fmt.Errorf("required argument not present: %s", argument.Name)
		}
		if argument.ValidateFunc == nil {
			continue
		}
		err := argument.ValidateFunc(arg)
		if err != nil {
			return err
		}
	}
	return nil
}

func validateParameters(parameters map[string]*Parameter) error {
	for _, param := range parameters {
		if param.ValidateFunc == nil {
			continue
		}
		err := param.ValidateFunc(*param)
		if err != nil {
			return err
		}
	}
	return nil
}

func removeIndex(s []string, index int) []string {
	ret := make([]string, 0)
	ret = append(ret, s[:index]...)
	return append(ret, s[index+1:]...)
}

func parseParameters(osArgs []string, params map[string]*Parameter) ([]string, error) {
	if len(params) == 0 {
		return osArgs, nil
	}
	var tempArgs = make([]string, len(osArgs))
	copy(tempArgs, osArgs)
	for _, param := range params {
		if param.Position >= len(tempArgs) {
			return os.Args, fmt.Errorf("parameter position [%v] exceeds the provided number of values: %v", param.Position, len(tempArgs))
		}
		param.Value = tempArgs[param.Position]
		tempArgs = removeIndex(tempArgs, param.Position)
	}
	return tempArgs, nil
}

func (c *CLI) Parse() error {
	if len(os.Args) < 2 {
		return fmt.Errorf("insufficient arguments")
	}
	for cmdName, cmd := range c.Commands {
		if cmdName != os.Args[1] {
			continue
		}
		newArgs, err := parseParameters(os.Args[2:], cmd.Params)
		if err != nil {
			return err
		}
		err = validateParameters(cmd.Params)
		if err != nil {
			return err
		}
		err = cmd.FlagSet.Parse(newArgs)
		if err != nil {
			return err
		}
		err = validateArguments(cmd.Flags)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("unknown command: %s", os.Args[1])
}

func CreateCLI(commands map[string]SubCommand) (*CLI, error) {
	var newCLI = CLI{
		Commands: make(map[string]*Command),
	}
	for name, cmd := range commands {
		newCLI.AddCommand(name, cmd.ErrorHandler)
		for _, arg := range cmd.Flags {
			err := newCLI.AddCommandArgs(name, arg)
			if err != nil {
				return nil, err
			}
		}
		for _, param := range cmd.Parameters {
			err := newCLI.AddCommandParameters(name, param)
			if err != nil {
				return nil, err
			}
		}
	}
	return &newCLI, nil
}
