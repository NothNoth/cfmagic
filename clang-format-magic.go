package main

import (
	"fmt"
	"os"
)

// JSON taken from https://zed0.co.uk/clang-format-configurator/doc/

const configs_file = "configs.json"

type ConfigEntry struct {
	Type    string   `json:"type"`
	Doc     string   `json:"doc"`
	Options []string `json:"options"`
}

func main() {
	var clangPath string
	var perfectSource string
	if len(os.Args) != 3 {
		fmt.Println("Usage: " + os.Args[0] + " <clang format binary> <perfect source code>")
		return
	}
	clangPath = os.Args[1]
	perfectSource = os.Args[2]

	//fmt.Println(entries)
	entries, err := loadConfig(clangPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	magic(entries, perfectSource)
}
