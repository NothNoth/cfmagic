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

//const configs_file = "test.json"

type ConfigEntry struct {
	Type string `json:"type"`
	Doc  string `json:"doc"`
}
type VersionConfigsStruct struct {
	ConfigEntries map[string]*ConfigEntry `json:"-"`
}

type Configs struct {
	Versions       []string                        `json:"versions"`
	VersionConfigs map[string]VersionConfigsStruct `json:"-"`
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
	var configs Configs
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

	//Read clang configs
	err = json.Unmarshal(configsData, &configs)
	if err != nil {
		fmt.Println("Failed to parse " + configs_file)
	}

	//Match with current local version
	configVersionEntry := ""
	for _, v := range configs.Versions {
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

	fmt.Println(configs)
}
