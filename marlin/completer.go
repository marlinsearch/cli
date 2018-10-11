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
		return appCompleter(args)
	case INDEX:
		return indexCompleter(args, d)
	}
	return rootCompleter(args)
}
