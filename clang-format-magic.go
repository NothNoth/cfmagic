package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

// JSON taken from https://zed0.co.uk/clang-format-configurator/doc/

const configs_file = "configs.json"

type ConfigEntry struct {
	Type    string   `json:"type"`
	Doc     string   `json:"doc"`
	Options []string `json:"options"`
}

func getClangFormatVersion(clangCmd string) string {

	out, err := exec.Command(clangCmd, "-version").Output()
	if err != nil {
		log.Fatal(err)
	}

	pos := strings.Index(string(out), "version ")
	if pos == -1 {
		return ""
	}

	version := string(out)[pos+8:]

	pos = strings.Index(version, " ")
	if pos == -1 {
		return ""
	}

	version = version[:pos]
	return version
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: " + os.Args[0] + " <clang format binary>")
		return
	}
	//Parse clang version
	version := getClangFormatVersion(os.Args[1])
	if len(version) == 0 {
		fmt.Println("Failed to fetch clang version")
		return
	}
	fmt.Println("Using clang version: " + version)
	configsData, err := ioutil.ReadFile(configs_file)
	if err != nil {
		fmt.Println("Failed to load " + configs_file)
		return
	}

	//Read clang configs : load as raw json (just extract keys and store internal value for later parsing)
	var data map[string]json.RawMessage
	err = json.Unmarshal(configsData, &data)
	if err != nil {
		fmt.Println("Failed to parse " + configs_file)
	}

	//Parse Values for version key
	var availableVersions []string
	err = json.Unmarshal(data["versions"], &availableVersions)
	if err != nil {
		fmt.Println("Failed to parse " + configs_file)
	}

	//Match with current local version
	configVersionEntry := ""
	for _, v := range availableVersions {
		if strings.Index(version, v) == 0 {
			configVersionEntry = v
			break
		}
	}
	if len(configVersionEntry) == 0 {
		fmt.Println("Failed to match version " + version + " with known versions, will use HEAD instead")
		configVersionEntry = "HEAD"
	}
	fmt.Println("Will be using configuration settings for version " + configVersionEntry)

	//Now parse keys for given version
	configVersionitems := data[configVersionEntry]

	var entries map[string]*ConfigEntry
	err = json.Unmarshal(configVersionitems, &entries)
	if err != nil {
		fmt.Println("Failed to parse " + configs_file)
	}

	for k, v := range entries {
		fmt.Println(k + ": [" + v.Doc[:16] + "...] " + v.Type)
		if len(v.Options) != 0 {
			fmt.Println(v.Options)
		}
	}
	//fmt.Println(entries)
}
