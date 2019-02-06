package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
)

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

func loadConfig(clangPath string) (map[string]*ConfigEntry, error) {

	var entries map[string]*ConfigEntry

	//Parse clang version
	version := getClangFormatVersion(clangPath)
	if len(version) == 0 {
		return nil, errors.New("Failed to fetch clang version")
	}
	fmt.Println("Using clang version: " + version)
	configsData, err := ioutil.ReadFile(configs_file)
	if err != nil {
		return nil, errors.New("Failed to load " + configs_file)
	}

	//Read clang configs : load as raw json (just extract keys and store internal value for later parsing)
	var data map[string]json.RawMessage
	err = json.Unmarshal(configsData, &data)
	if err != nil {
		return nil, errors.New("Failed to parse " + configs_file)
	}

	//Parse Values for version key
	var availableVersions []string
	err = json.Unmarshal(data["versions"], &availableVersions)
	if err != nil {
		return nil, errors.New("Failed to parse " + configs_file)
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

	err = json.Unmarshal(configVersionitems, &entries)
	if err != nil {
		return nil, errors.New("Failed to parse " + configs_file)
	}

	for k, v := range entries {
		fmt.Println(k + ": [" + v.Doc[:16] + "...] " + v.Type)
		if len(v.Options) != 0 {
			fmt.Println(v.Options)
		}
	}

	return entries, nil
}
