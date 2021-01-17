package main

import (
	"flag"
	"fmt"
	"github.com/EngineersBox/ModularCLI/cli"
	"log"
	"strings"
)

var commands = map[string]cli.SubCommand{
	"test": {
		ErrorHandler: flag.ExitOnError,
		Arguments: []*cli.Argument{
			{
				Type:         cli.TypeString,
				Name:         "file",
				DefaultValue: "",
				HelpMsg:      "File to import (*.ext)",
				Required:     true,
				ValidateFunc: func(arg cli.TypedArgument) error {
					if !strings.Contains(*arg.GetString(), ".ext") {
						return fmt.Errorf("filetype must be .ext")
					}
					return nil
				},
			},
			{
				Type:         cli.TypeBool,
				Name:         "recursive",
				DefaultValue: false,
				HelpMsg:      "Whether to recursively import files [default: false]",
				Required:     false,
			},
			{
				Type:         cli.TypeInt,
				Name:         "count",
				DefaultValue: 4,
				HelpMsg:      "How many files to count [default: 4]",
				Required:     false,
			},
		},
		Parameters: []*cli.Parameter{
			{
				Type:     cli.TypeString,
				Name:     "instance_type",
				Position: 0,
				ValidateFunc: func(param cli.Parameter) error {
					if !strings.Contains(*param.GetString(), "s3") {
						return fmt.Errorf("invalid instance type, must be 's3")
					}
					return nil
				},
			},
		},
	},
}

func main() {
	schematicCli, err := cli.CreateCLI(commands)
	if err != nil {
		log.Fatal(err)
	}
	err = schematicCli.Parse()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(*schematicCli.Commands["test"].Params["instance_type"].GetString())
	fmt.Println(*schematicCli.Commands["test"].Flags["file"].GetString())
	fmt.Println(*schematicCli.Commands["test"].Flags["recursive"].GetBool())
	fmt.Println(*schematicCli.Commands["test"].Flags["count"].GetInt())
}
