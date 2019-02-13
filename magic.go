package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"os/signal"
	"syscall"
)

const minStdDevForMutationBoost = 5.0

func magic(clangPath string, configEntries map[string]*ConfigEntry, perfectSource string, populationSize uint32, mutationRate uint32) {

	pop := genPopulation(populationSize, configEntries)
	fmt.Println("Population ready")
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

		//Boost mutations if top population is too homogenous
		if stdDev < minStdDevForMutationBoost {
			//Mix
			pop.mix(2*mutationRate, configEntries)

		} else {
			//Mix
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
