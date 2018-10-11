package marlin

import (
	"encoding/json"
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cast"
	"os"
)

func displayHelp(commands []prompt.Suggest) {
	fmt.Println("\nAvailable Commands:\n")
	for _, command := range commands {
		fmt.Printf("  %-20s %s\n", command.Text, command.Description)
	}
	fmt.Println("")
}

func displayTableJson(jstr string, header []string) {
	var result []map[string]interface{}
	if err := json.Unmarshal([]byte(jstr), &result); err == nil {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader(header)
		for _, d := range result {
			var row []string
			for _, h := range header {
				row = append(row, cast.ToString(d[h]))
			}
			table.Append(row)
		}
		table.Render()
		fmt.Println("")
	} else {
		fmt.Println(err)
	}
}

func displayInvalidCommand(cmd string) {
	fmt.Println("Error: Command", "'"+cmd+"'", "not found.  Please try 'help' to see possible commands")
}
