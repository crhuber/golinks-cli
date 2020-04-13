package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
)

// GoLinksConfig does stuff
type GoLinksConfig struct {
	Hostname string `json:"hostname"`
	Port     string `json:"port"`
	Protocol string `json:"protocol"`
}

// GolinksError blah
type GolinksError struct {
	Message string
}

func (e *GolinksError) Error() string {
	return fmt.Sprintf("message %s", e.Message)
}

var home, err = os.UserHomeDir()
var golinksDir = "/.golinks"
var glConf = filepath.Join(home, golinksDir, "golinks.json")

// GoLinksSearchResult blah
type GoLinksSearchResult struct {
	Results []struct {
		DestinationURL string `json:"destination_url"`
		DateAdded      string `json:"date_added"`
		Keyword        string `json:"keyword"`
		Description    string `json:"description"`
		Views          int    `json:"views"`
	} `json:"results"`
}

func dirExists(dir string) bool {
	info, err := os.Stat(dir)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func initialize(hostname, port, protocol string) {
	gd := filepath.Join(home, golinksDir)
	gconf := filepath.Join(home, golinksDir, "kelp.json")
	if dirExists(gd) == false {
		fmt.Println("Creating Golinks dir...")
		err := os.Mkdir(gd, 0777)
		if err != nil {
			fmt.Println(err)
		}
	}

	// create empty config
	if fileExists(gconf) == false {
		glc := GoLinksConfig{
			Hostname: hostname,
			Port:     port,
			Protocol: protocol,
		}
		bs, err := json.MarshalIndent(glc, "", " ")
		if err != nil {
			fmt.Println(bs)
		}
		ioutil.WriteFile(glConf, bs, 0644)
	}

	fmt.Println("🌱 Golinks Initialized!")
}

func loadGolinksConfig() GoLinksConfig {
	bs, _ := ioutil.ReadFile(glConf)
	var glc GoLinksConfig
	err := json.Unmarshal(bs, &glc)
	if err != nil {
		fmt.Println(err)
	}
	return glc
}

func queryAPI(keyword, hostname, port, protocol string) (string, error) {
	fmt.Println("\nSearching ...")
	url := fmt.Sprintf("%s://%s:%s/api/search?q=%s", protocol, hostname, port, keyword)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	glsr := GoLinksSearchResult{}

	if err := json.Unmarshal(body, &glsr); err != nil {
		panic(err)
	}

	var options []string
	for _, result := range glsr.Results {
		options = append(options, result.Keyword)
	}
	// the answers will be written to this struct
	answers := struct {
		Result string `survey:"result"`
	}{}
	if len(options) > 0 {
		// the questions to ask
		var qs = []*survey.Question{
			{
				Name: "result",
				Prompt: &survey.Select{
					Message: "Open golink:",
					Options: options,
					Default: options[0],
				},
			},
		}

		// perform the questions
		err = survey.Ask(qs, &answers)
		if err == terminal.InterruptErr {
			fmt.Println("interrupted")
			os.Exit(0)
		} else if err != nil {
			panic(err)
		}

		fmt.Printf("Opening %s.", answers.Result)

	} else {
		fmt.Println("No results.")
		return "", &GolinksError{
			Message: "No Results",
		}
	}
	err = nil

	return answers.Result, err
}

func openBrowser(url string) {
	fmt.Printf("Opening %s", url)

	switch runtime.GOOS {
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}

func queryBrowse(keyword string) {
	glc := loadGolinksConfig()
	keyword, err := queryAPI(keyword, glc.Hostname, glc.Port, glc.Protocol)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	url := fmt.Sprintf("%s://%s:%s/%s", glc.Protocol, glc.Hostname, glc.Port, keyword)
	openBrowser(url)
}

func main() {
	Cli()
}
