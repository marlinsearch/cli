package marlin

import (
	"encoding/json"
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/olekukonko/tablewriter"
	"os"
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

func displayTableJson(jstr string, header []string) {
	var result []map[string]interface{}
	if err := json.Unmarshal([]byte(jstr), &result); err == nil {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader(header)
		for _, d := range result {
			var row []string
			for _, h := range header {
				row = append(row, d[h].(string))
			}
			table.Append(row)
		}
		table.Render()
		fmt.Println("")
	} else {
		fmt.Println(err)
	}
}

func listApps() {
	resp, success := MarlinApi.getApplications()
	if success {
		h := []string{"name", "appId", "apiKey"}
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
	}

}
