package marlin

import (
	"flag"
	"fmt"
	"github.com/c-bata/go-prompt"
	"strings"
)

type Context int

const (
	ROOT Context = iota
	APP
	INDEX
)

var CliState struct {
	Host           string
	MasterAppId    string
	MasterApiKey   string
	Connected      bool
	CliPrefix      string
	ActiveApp      string
	ActiveIndex    string
	CurrentContext Context
}

var MarlinApi *Api
var splash = `
              __            
|V| _ __|o__ (_  _  _ ___|_ 
| |(_|| ||| |__)(/_(_||(_| |
                                                                          
`

func changeLivePrefix() (string, bool) {
	CliState.CliPrefix = getCliPrefix()
	return CliState.CliPrefix, true
}

func getCliPrefix() string {
	if CliState.Connected == false {
		return ">>> "
	}
	htok := strings.Split(CliState.Host, "/")
	prefix := htok[len(htok)-1]
	prefix += "> "
	switch CliState.CurrentContext {
	case ROOT:
		return prefix
	case APP:
		return CliState.ActiveApp + "@" + prefix
	case INDEX:
		return CliState.ActiveApp + "/" + CliState.ActiveIndex + "@" + prefix
	}
	return prefix
}

func Init(args []string) {
	f := flag.NewFlagSet("main", flag.ExitOnError)
	hostPtr := f.String("m", "http://localhost:9002", "Host to connect to Eg., http://localhost:9002")
	appIdPtr := f.String("a", "abcdefgh", "The master application id")
	apiKeyPtr := f.String("k", "12345678901234567890123456789012", "The master api key")
	f.Parse(args)

	fmt.Println(splash)

	defer fmt.Println("Bye!")

	CliState.Host = *hostPtr
	CliState.MasterAppId = *appIdPtr
	CliState.MasterApiKey = *apiKeyPtr
	MarlinApi = &Api{CliState.Host, CliState.MasterAppId, CliState.MasterApiKey}

	connected, version := MarlinApi.Connect()
	CliState.Connected = connected
	if connected {
		fmt.Println("Connected to", CliState.Host, "Marlin version", version)
		CliState.CliPrefix = getCliPrefix()
	} else {
		CliState.CliPrefix = ">>> "
		fmt.Println("Failed to connect to ", CliState.Host, ". Please use connect command to connect")
	}
	fmt.Println("Please use `exit` or `Ctrl-D` to exit this program.\n")

	p := prompt.New(
		Executor,
		Completer,
		prompt.OptionPrefix(">>> "),
		prompt.OptionCompletionWordSeparator(" /"),
		prompt.OptionLivePrefix(changeLivePrefix),
		prompt.OptionInputTextColor(prompt.Yellow),
		//prompt.OptionSuggestionTextColor(prompt.White),
		//prompt.OptionSuggestionBGColor(prompt.DarkGray),
		prompt.OptionTitle("Marlin CLI "),
	)
	p.Run()
}
