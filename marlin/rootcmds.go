package marlin

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"strings"
)

var rootCommands = []prompt.Suggest{
	{Text: "info", Description: "Display Marlin Information"},
	{Text: "list-applications", Description: "List all available applications"},
	{Text: "create-application", Description: "Create an application"},
	{Text: "application", Description: "Choose an application, enter into app context"},
	{Text: "help", Description: "Displays list of available commands"},
	{Text: "exit", Description: "Exit this program"},
}

func rootCompleter(args []string) []prompt.Suggest {
	if len(args) <= 1 {
		return prompt.FilterHasPrefix(rootCommands, args[0], true)
	}
	return []prompt.Suggest{}
}

func displayResult(resp string, success bool) {
	if success {
		fmt.Println(resp)
	} else {
		fmt.Println("Failed to execute command")
	}
	fmt.Println("")
}

func displayInfo() {
	resp, success := MarlinApi.getInfo()
	displayResult(resp, success)
}

func listApps() {
	resp, success := MarlinApi.getApplications()
	displayResult(resp, success)
}

func performRootCommand(in string) {
	args := strings.Split(in, " ")
	switch args[0] {
	case "help":
		displayHelp(rootCommands)
	case "ls":
		listApps()
	case "list-applications":
		listApps()
	case "info":
		displayInfo()
	}

}
