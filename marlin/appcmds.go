package marlin

import (
	"encoding/json"
	"fmt"
	"github.com/c-bata/go-prompt"
	"strconv"
	"strings"
)

var appCommands = []prompt.Suggest{
	{Text: "list-indexes", Description: "List all available indexes"},
	{Text: "create-index", Description: "Create an index"},
	{Text: "index", Description: "Choose an index, enter into index context"},
	{Text: "help", Description: "Displays list of available commands"},
	{Text: "exit", Description: "Exit this program"},
	{Text: "..", Description: "Exit application context back to root context"},
}

func appCompleter(args []string) []prompt.Suggest {
	if len(args) <= 1 {
		return prompt.FilterHasPrefix(appCommands, args[0], true)
	}
	switch args[0] {
	case "create-index":
		if len(args) == 3 {
			return []prompt.Suggest{
				{Text: "num-shards"},
			}
		}
	case "index":
		return prompt.FilterHasPrefix(getIndexPrompts(), args[1], true)
	}
	return []prompt.Suggest{}
}

func getIndexPrompts() []prompt.Suggest {
	resp, success := MarlinApi.getIndexes()
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

func createIndex(args []string) {
	if len(args) < 2 || len(args[1]) == 0 {
		fmt.Println("Error: create-index <Index Name> num-shards <No of shards>\n")
		return
	}
	numShards := 1
	if len(args) == 4 {
		if i, err := strconv.Atoi(args[3]); err != nil {
			numShards = i
		}
	}

	resp, success := MarlinApi.createIndex(args[1], numShards)
	displayResult(resp, success)
}

func listIndexes() {
	resp, success := MarlinApi.getIndexes()
	if success {
		h := []string{"name", "numShards"}
		displayTableJson(resp, h)
	}
}

func isValidIndex(name string) bool {
	resp, success := MarlinApi.getIndexes()
	if success {
		var result []map[string]interface{}
		if err := json.Unmarshal([]byte(resp), &result); err == nil {
			for _, d := range result {
				if d["name"].(string) == name {
					return true
				}
			}
		}
	}
	return false
}

func chooseIndex(args []string) {
	if len(args) != 2 {
		fmt.Println("Error: application <App Name>\n")
		return
	}
	if isValidIndex(args[1]) {
		CliState.ActiveIndex = args[1]
		CliState.CurrentContext = INDEX
	} else {
		fmt.Println("Error:  Could not retrieve index details")
	}
}

func exitAppContext() {
	MarlinApi.AppId = CliState.MasterAppId
	MarlinApi.ApiKey = CliState.MasterApiKey
	CliState.CurrentContext = ROOT
}

func performAppCommand(in string) {
	args := strings.Split(in, " ")
	switch args[0] {
	case "help":
		displayHelp(appCommands)
	case "create-index":
		createIndex(args)
	case "list-indexes":
		listIndexes()
	case "ls":
		listIndexes()
	case "index":
		chooseIndex(args)
	case "..":
		exitAppContext()
	default:
		displayInvalidCommand(args[0])
	}
}
