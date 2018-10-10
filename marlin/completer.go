package marlin

import (
	"github.com/c-bata/go-prompt"
	"strings"
)

func Completer(d prompt.Document) []prompt.Suggest {
	if d.TextBeforeCursor() == "" {
		return []prompt.Suggest{}
	}
	args := strings.Split(d.TextBeforeCursor(), " ")
	switch CliState.CurrentContext {
	case ROOT:
		return rootCompleter(args)
	case APP:
	case INDEX:
	}
	return rootCompleter(args)
}
