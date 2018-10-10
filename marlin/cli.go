package marlin

import (
	"flag"
	"fmt"
	"github.com/c-bata/go-prompt"
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
Marlin Search
`

func changeLivePrefix() (string, bool) {
	return CliState.CliPrefix, true
}

func Init(args []string) {
	f := flag.NewFlagSet("main", flag.ExitOnError)
	hostPtr := f.String("m", "http://localhost:9002", "Host to connect to Eg., http://localhost:9002")
	appIdPtr := f.String("a", "abcdefgh", "The master application id")
	apiKeyPtr := f.String("k", "12345678901234567890123456789012", "The master api key")
	f.Parse(args)

	fmt.Println(splash)
	fmt.Println("Please use `exit` or `Ctrl-D` to exit this program.\n")

	defer fmt.Println("Bye!")

	CliState.Host = *hostPtr
	CliState.MasterAppId = *appIdPtr
	CliState.MasterApiKey = *apiKeyPtr
	MarlinApi = &Api{CliState.Host, CliState.MasterAppId, CliState.MasterApiKey}

	connected, version := MarlinApi.Connect()
	CliState.Connected = connected
	if connected {
		fmt.Println("Connected to", CliState.Host, "Marlin version", version)
		CliState.CliPrefix = CliState.Host + "> "
	} else {
		CliState.CliPrefix = ">>> "
		fmt.Println("Failed to connect to ", CliState.Host, ". Please use connect command to connect")
	}

	p := prompt.New(
		Executor,
		Completer,
		prompt.OptionPrefix(">>> "),
		prompt.OptionLivePrefix(changeLivePrefix),
		prompt.OptionInputTextColor(prompt.Yellow),
		//prompt.OptionSuggestionTextColor(prompt.White),
		//prompt.OptionSuggestionBGColor(prompt.DarkGray),
		prompt.OptionTitle("Marlin CLI "),
	)
	p.Run()
}
