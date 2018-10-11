package marlin

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/c-bata/go-prompt/completer"
	"io/ioutil"
	"os"
	"os/user"
	"strings"
	"time"
)

var indexCommands = []prompt.Suggest{
	{Text: "load", Description: "Load the index with data"},
	{Text: "clear", Description: "Clear the index"},
	{Text: "reindex", Description: "Reindex the index data"},
	{Text: "mapping", Description: "Display index mapping information"},
	{Text: "help", Description: "Displays list of available commands"},
	{Text: "exit", Description: "Exit this program"},
	{Text: "..", Description: "Exit index context back to application context"},
}

var jsonFileCompleter = completer.FilePathCompleter{
	IgnoreCase: true,
	Filter: func(fi os.FileInfo) bool {
		if fi.IsDir() {
			return true
		}
		if strings.HasSuffix(fi.Name(), ".json") || strings.HasSuffix(fi.Name(), ".js") {
			return true
		}
		return false
	},
}

func indexCompleter(args []string, d prompt.Document) []prompt.Suggest {
	if len(args) <= 1 {
		return prompt.FilterHasPrefix(indexCommands, args[0], true)
	}
	switch args[0] {
	case "load":
		return jsonFileCompleter.Complete(d)
	}
	return []prompt.Suggest{}
}

func exitIndexContext() {
	CliState.CurrentContext = APP
}

// Check if every line is a json object by loading the first line alone
func isNewLineJson(path string) bool {
	inFile, _ := os.Open(path)
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	if scanner.Scan() {
		var j interface{}
		if json.Unmarshal(scanner.Bytes(), &j) == nil {
			return true
		}
	}
	return false
}

func loadNewLineJsonFile(path string) {
	start := time.Now()
	inFile, _ := os.Open(path)
	defer inFile.Close()

	scanner := bufio.NewScanner(inFile)
	const maxCapacity = 512 * 1024 * 100
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	start2 := time.Now()
	numperpush := 10000

	// Push in chunks of 10000 objects
	pushjs := make([]interface{}, 0)
	count := 0
	trials := 0

	for scanner.Scan() {
		var j interface{}
		if json.Unmarshal(scanner.Bytes(), &j) != nil {
			continue
		}
		pushjs = append(pushjs, j)
		count++
		if count >= numperpush {
			for {
				time.Sleep(100 * time.Millisecond)
				if MarlinApi.getNumIndexJobs() == 0 {
					break
				}
			}
			count = 0
			trials++
			jjs, _ := json.Marshal(pushjs)
			body, success := MarlinApi.addObjectsToIndex(string(jjs[:]))
			if !success {
				fmt.Println("Failed to post data, stopping upload")
				return
			}
			elapsed2 := time.Since(start2)
			fmt.Println("Data Add Response", body, trials*numperpush, elapsed2)
			pushjs = make([]interface{}, 0)
			start2 = time.Now()
		}
	}
	// THe last chunk
	if len(pushjs) > 0 {
		jjs, _ := json.Marshal(pushjs)
		body, success := MarlinApi.addObjectsToIndex(string(jjs[:]))
		if !success {
			fmt.Println("Failed to post data, stopping upload")
			return
		}
		fmt.Println("Done Adding", string(body), trials*numperpush+count)
	}
	elapsed := time.Since(start)
	// Lookup and sort took
	fmt.Println("Adding took", elapsed)
}

// TODO; Cleanup.. and add support to interrupt upload
func loadFile(in string) {
	path := in[5:]
	path = strings.TrimSpace(path)
	if strings.HasPrefix(path, "~") {
		usr, _ := user.Current()
		path = strings.TrimPrefix(path, "~")
		path = usr.HomeDir + path
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("Invalid path", path)
		fmt.Println("")
		return
	}

	if isNewLineJson(path) {
		loadNewLineJsonFile(path)
		return
	}

	// It is a json array?
	if file, err := ioutil.ReadFile(path); err == nil {
		var j interface{}
		if json.Unmarshal(file, &j) != nil {
			fmt.Println("Failed to load file")
			return
		}
		js := j.([]interface{})
		limit := 10000
		for len(js) > 0 {
			var pushjs []interface{}
			if len(js) > limit {
				pushjs = js[:10000]
				js = js[10000:]
			} else {
				pushjs = js[:len(js)]
				js = js[:0]
			}
			jjs, _ := json.Marshal(pushjs)
			body, success := MarlinApi.addObjectsToIndex(string(jjs[:]))
			if !success {
				fmt.Println("Failed to load file")
				return
			}
			fmt.Println("Data Add Response", string(body))
		}
		return
	}
	fmt.Println("Failed to load file", path)

}

func performIndexCommand(in string) {
	args := strings.Split(in, " ")
	switch args[0] {
	case "help":
		displayHelp(indexCommands)
	case "load":
		loadFile(in)
	case "..":
		exitIndexContext()
	default:
		displayInvalidCommand(args[0])
	}

}
