package command

import (
	"bytes"
	"fmt"
	"github.com/hpehl/whatunga/model"
	"strings"
)

var helpUsage = "help [command]"

var help = Command{
	"help",
	"Displays this help message or prints detailed help on requested commands.",
	helpUsage,
	"Displays this help message or prints detailed help on requested commands.",
	// tab completer
	func(_ *model.Project, query, _ string) []string {
		var results []string
		for key, _ := range Registry {
			if strings.HasPrefix(key, query) {
				results = append(results, key)
			}
		}
		return results
	},
	// action
	func(_ *model.Project, args []string) error {
		if len(args) == 0 {
			// general help
			fmt.Print("Commands:\n\n")
			var buffer bytes.Buffer
			Registry.forEach(func(cmd Command) {
				buffer.WriteString(fmt.Sprintf("    %s%s\n", pad(cmd.Name, 20), cmd.Description))
			})
			fmt.Print(buffer.String())
			fmt.Print("\nMore command help available using \"help <command>\"\n")
		} else if len(args) == 1 {
			cmd, valid := Registry[args[0]]
			if !valid {
				fmt.Printf("Unknown command: \"%s\"\n", args[0])
			} else {
				fmt.Printf("Usage: %s\n\n%s\n", cmd.Usage, cmd.UsageDescription)
			}
		} else {
			return fmt.Errorf("Too many arguments. Usage: %s", helpUsage)
		}
		return nil
	},
}

func pad(str string, n int) string {
	stringLength := len(str)
	padLength := n - stringLength
	return fmt.Sprintf("%s%s", str, strings.Repeat(" ", padLength))
}
