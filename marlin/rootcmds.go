package marlin

import (
	"encoding/json"
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
	{Text: "connect", Description: "Connect to a different host"},
	{Text: "exit", Description: "Exit this program"},
}

func rootCompleter(args []string) []prompt.Suggest {
	if len(args) <= 1 {
		return prompt.FilterHasPrefix(rootCommands, args[0], true)
	}

	if len(args) == 2 {
		switch args[0] {
		case "application":
			return prompt.FilterHasPrefix(getApplicationPrompts(), args[1], true)
		}
	}
	return []prompt.Suggest{}
}

func getApplicationPrompts() []prompt.Suggest {
	resp, success := MarlinApi.getApplications()
	if success {
		var result []map[string]interface{}
		if err := json.Unmarshal([]byte(resp), &result); err == nil {
			var res []prompt.Suggest
			for _, d := range result {
				res = append(res, prompt.Suggest{Text: d["name"].(string)})
			}
			return res
		}
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
	if success {
		h := []string{"name", "appId", "apiKey", "numIndexes"}
		displayTableJson(resp, h)
	}
}

func createApp(args []string) {
	if len(args) != 2 {
		fmt.Println("Error: create-application <App Name>\n")
		return
	}
	resp, success := MarlinApi.createApplication(args[1])
	displayResult(resp, success)
}

func chooseApplication(args []string) {
	if len(args) != 2 {
		fmt.Println("Error: application <App Name>\n")
		return
	}
	resp, success := MarlinApi.getApplication(args[1])
	if success {
		var result map[string]interface{}
		if err := json.Unmarshal([]byte(resp), &result); err == nil {
			MarlinApi.AppId = result["appId"].(string)
			MarlinApi.ApiKey = result["apiKey"].(string)
			CliState.ActiveApp = args[1]
			CliState.CurrentContext = APP
		}
	} else {
		fmt.Println("Error:  Could not retrieve application details")
	}
}

func connectToHost(args []string) {
	if len(args) != 2 {
		fmt.Println("Error: connect <http(s)://hostname:port>\n")
		return
	}
	MarlinApi.Host = args[1]
	CliState.Host = args[1]

	connected, version := MarlinApi.Connect()
	CliState.Connected = connected
	if connected {
		fmt.Println("Connected to", CliState.Host, "Marlin version", version)
		CliState.CliPrefix = getCliPrefix()
	} else {
		CliState.CliPrefix = ">>> "
		fmt.Println("Failed to connect to ", CliState.Host, ". Please use connect command to connect")
	}
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
	case "create-application":
		createApp(args)
	case "info":
		displayInfo()
	case "application":
		chooseApplication(args)
	case "connect":
		connectToHost(args)
	default:
		displayInvalidCommand(args[0])
	}

}
