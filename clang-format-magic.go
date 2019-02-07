package main

import (
	"fmt"
	"os"
	"strconv"
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
	if len(os.Args) != 5 {
		fmt.Println("Usage: " + os.Args[0] + " <clang format binary> <perfect source code> <population size> <mutation rate>")
		return
	}
	clangPath = os.Args[1]
	perfectSource = os.Args[2]
	pop, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Println("Error : invalid population size (default: 100)")
		return
	}
	mut, err := strconv.Atoi(os.Args[4])
	if err != nil {
		fmt.Println("Error : invalid mutation rate (default: 4)")
		return
	}

	//fmt.Println(entries)
	entries, err := loadConfig(clangPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	magic(clangPath, entries, perfectSource, uint32(pop), uint32(mut))
}
