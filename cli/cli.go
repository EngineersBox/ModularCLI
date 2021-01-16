package cli

import (
	"flag"
	"fmt"
	"os"
)

type Command struct {
	FlagSet *flag.FlagSet
	Flags   map[string]TypedArgument
}

type CLI struct {
	Commands map[string]*Command
}

func (c *CLI) AddCommand(name string, handler flag.ErrorHandling) {
	c.Commands[name] = &Command{
		FlagSet: flag.NewFlagSet(name, handler),
		Flags:   make(map[string]TypedArgument),
	}
}

func (c *CLI) AddCommandArgs(cmdName string, arg *Argument) error {
	var currentCmd = c.Commands[cmdName]
	switch arg.Type {
	case TypeString:
		currentCmd.Flags[arg.Name] = StringArgument{
			Value:    currentCmd.FlagSet.String(arg.Name, arg.DefaultValue.(string), arg.HelpMsg),
			Argument: arg,
		}
		break
	case TypeBool:
		currentCmd.Flags[arg.Name] = BoolArgument{
			Value:    currentCmd.FlagSet.Bool(arg.Name, arg.DefaultValue.(bool), arg.HelpMsg),
			Argument: arg,
		}
		break
	case TypeInt:
		currentCmd.Flags[arg.Name] = IntArgument{
			Value:    currentCmd.FlagSet.Int(arg.Name, arg.DefaultValue.(int), arg.HelpMsg),
			Argument: arg,
		}
		break
	case TypeInvalid:
		return fmt.Errorf("invalid arugment type: %d", arg.Type)
	}
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

func (c *CLI) Parse() error {
	if len(os.Args) < 2 {
		return fmt.Errorf("insufficient arguments")
	}
	for cmdName, cmd := range c.Commands {
		if cmdName == os.Args[1] {
			err := cmd.FlagSet.Parse(os.Args[2:])
			if err != nil {
				return err
			}
			err = validateArguments(cmd.Flags)
			if err != nil {
				return err
			}
			return nil
		}
	}
	return fmt.Errorf("unknown command: %s", os.Args[1])
}

func CreateCLI(commands map[string]SubCommand) (*CLI, error) {
	var newCLI = CLI{
		Commands: make(map[string]*Command),
	}
	for name, cmd := range commands {
		newCLI.AddCommand(name, cmd.ErrorHandler)
		for _, arg := range cmd.Arguments {
			err := newCLI.AddCommandArgs(name, arg)
			if err != nil {
				return nil, err
			}
		}
	}
	return &newCLI, nil
}
