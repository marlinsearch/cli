package marlin

import (
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
		performAppCommand(in)
	case INDEX:
		performIndexCommand(in)
	}
}
