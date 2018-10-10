package marlin

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"os"
)

func Executor(in string) {
	if in == "q" || in == "quit" || in == "exit" {
		os.Exit(0)
	}
	switch CliState.CurrentContext {
	case ROOT:
		performRootCommand(in)
	case APP:
	case INDEX:
	}
}

func displayHelp(commands []prompt.Suggest) {
	fmt.Println("\nAvailable Commands:\n")
	for _, command := range commands {
		fmt.Printf("  %-20s %s\n", command.Text, command.Description)
	}
	fmt.Println("")
}
