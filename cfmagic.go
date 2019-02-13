package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

// JSON taken from https://zed0.co.uk/clang-format-configurator/doc/

const configs_file = "configs.json"
const minStdDevForMutationBoost = 5.0

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
		fmt.Println("Example: " + os.Args[0] + "/usr/bin/clang-format-6.0 perfect.c 20 4")
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

	cfmagic(clangPath, entries, perfectSource, uint32(pop), uint32(mut))
}

func cfmagic(clangPath string, configEntries map[string]*ConfigEntry, perfectSource string, populationSize uint32, mutationRate uint32) {
	if (populationSize < 2) || (mutationRate > 100) {
		fmt.Println("Invalid settings")
		return
	}

	pop := genPopulation(populationSize, configEntries)
	fmt.Printf("Population size : %d | Mutation rate %d %%\n", populationSize, mutationRate)

	done := false
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		s := <-sigc
		fmt.Println("Received signal " + s.String())
		done = true
	}()

	fmt.Println("Hit CTRL+C to end iterations and save the current best result")

	for {
		pop.generation++
		//Eval everyone
		for i := 0; i < len(pop.population); i++ {
			err := pop.population[i].UpdateScore(clangPath, perfectSource)
			if err != nil {
				fmt.Println(err)
				//return
				pop.population[i].score = uint32(math.MaxUint32)
			}
		}

		//Sort
		pop.sort()
		stdDev := pop.getStdDev(len(pop.population) / 2)
		fmt.Printf("Best score for generation %d: %d (stdDev: %f)\n", pop.generation, pop.population[0].score, stdDev)

		//Mix
		//Boost mutations if top population is too homogenous
		if stdDev < minStdDevForMutationBoost {
			pop.mix(2*mutationRate, configEntries)
		} else {
			pop.mix(mutationRate, configEntries)
		}

		if done == true {
			break
		}
	}

	fmt.Printf("Best configuration has score: %d\n", pop.population[0].score)
	ioutil.WriteFile(".clang-format", []byte(pop.population[0].toClangFormatConfigFile()), os.ModePerm)
	fmt.Println("Written to .clang-format")
}
