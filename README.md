# ModularCLI
A GoLang modular CLI

## Example

`main.go`:
```Go
package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/EngineersBox/ModularCLI/cli"
)

var commands = map[string]cli.SubCommand{
    "create": {
        ErrorHandler: flag.ExitOnError,
        Arguments: []*cli.Argument{
            {
                Type:         cli.TypeString,
                Name:         "input",
                DefaultValue: "",
                HelpMsg:      "File to read from",
                Required:     true,
            },
            {
                Type:         cli.TypeBool,
                Name:         "count",
                DefaultValue: false,
                HelpMsg:      "How many instances to create",
                Required:     false,
            },
        },
    },
    "dataset": {
        ErrorHandler: flag.ExitOnError,
        Arguments: []*cli.Argument{
            {
                Type:         cli.TypeString,
                Name:         "from",
                DefaultValue: "",
                HelpMsg:      "URL to retrieve data from",
                Required:     true,
				ValidateFunc: func(arg cli.TypedArgument) error {
					if !strings.Contains(*arg.GetString(), "http") {
						return fmt.Errorf("url must use HTTP protocol")
					}
					return nil
				},
            },
            {
                Type:         cli.TypeBool,
                Name:         "recursive",
                DefaultValue: false,
                HelpMsg:      "Whether to import nested directories [default: false]",
                Required:     false,
            },
            {
                Type:         cli.TypeBool,
                Name:         "count",
                DefaultValue: false,
                HelpMsg:      "How many files to read",
                Required:     false,
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
	
    fmt.Println(*schematicCli.Commands["dataset"].Flags["from"].GetString())
    fmt.Println(*schematicCli.Commands["dataset"].Flags["count"].GetInt())
    fmt.Println(*schematicCli.Commands["dataset"].Flags["recursive"].GetBool())
}
```

Usage:
```bash
go build -o out/cli_test main.go
./out/cli_test dataset --from=localhost:8080 --count=28 --recursive
```

Result:
```text
localhost:8080
28
true
```