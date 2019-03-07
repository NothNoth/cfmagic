package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
)

// JSON taken from https://zed0.co.uk/clang-format-configurator/doc/

const configs_file = "configs.json"
const minStdDevForMutationBoost = 5.0
const scoreUnitialized = math.MaxUint32

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

func evalPopulation(clangPath string, perfectSourceData []byte, pop *Population, fromIdx int, toIdx int, doneCH chan int) error {
	for i := fromIdx; i < toIdx; i++ {
		if pop.population[i].score != scoreUnitialized {
			continue
		}
		err := pop.population[i].UpdateScore(clangPath, perfectSourceData)
		if err != nil {
			fmt.Println(err)
			pop.population[i].score = scoreUnitialized
		}
	}
	doneCH <- 1
	return nil
}

func cfmagic(clangPath string, configEntries map[string]*ConfigEntry, perfectSource string, populationSize uint32, mutationRate uint32) {
	if (populationSize < 2) || (mutationRate > 100) {
		fmt.Println("Invalid settings")
		return
	}

	perfectSourceData, err := ioutil.ReadFile(perfectSource)
	if err != nil {
		fmt.Println(err)
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
	threads := runtime.NumCPU() - 1
	if threads == 0 {
		threads = 1
	}

	threads = 1

	fmt.Printf("Will use %d threads\n", threads)
	chunks := len(pop.population) / threads
	for {
		pop.generation++

		for z := 0; z < len(pop.population); z++ {
			pop.population[z].score = scoreUnitialized
		}
		//Eval everyone
		doneCh := make(chan int, threads)
		for t := 0; t < threads; t++ {
			var from int
			var to int
			from = t * chunks
			if t == threads-1 {
				to = len(pop.population)
			} else {
				to = (t + 1) * chunks
			}

			go evalPopulation(clangPath, perfectSourceData, &pop, from, to, doneCh)
		}
		// Wait for all goroutines to exit
		for i := 0; i < threads; i++ {
			_ = <-doneCh
		}

		//Sort
		pop.sort()
		stdDev := pop.getStdDev(len(pop.population) / 2)
		fmt.Printf("Best score for generation %d: %d (stdDev: %f)\n", pop.generation, pop.population[0].score, stdDev)
		for z := 0; z < len(pop.population); z++ {
			fmt.Printf("#%d score : %d\n", z, pop.population[z].score)
		}
		if pop.population[0].score == 0 {
			fmt.Println("Found perfect configuration, stopping here.")
			break
		}

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
