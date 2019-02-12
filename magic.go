package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"os/signal"
	"sort"
	"syscall"
)

var populationSize uint32
var baseMutationRate uint32
var mutationRate uint32

const minStdDevForMutationBoost = 5.0

type Population struct {
	population []Individual
	generation uint32
}

type AllIndividuals []Individual

func (pop AllIndividuals) Len() int {
	return len(pop)
}

func (a AllIndividuals) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a AllIndividuals) Less(i, j int) bool { return a[i].score < a[j].score }

func (pop Population) String() string {
	var s string

	for idx, p := range pop.population {
		s += fmt.Sprintf("#%d: %s %d\n", idx, p.String(), p.score)
	}

	return s
}

func (pop Population) mix(configEntries map[string]*ConfigEntry) {
	var i uint32
	for i = 0; i < populationSize/2; i++ {
		var motherID int
		fatherID := rand.Intn(int(populationSize / 2))
		for {
			motherID := rand.Intn(int(populationSize / 2))
			if fatherID != motherID {
				break
			}
		}
		baby := pop.population[motherID].mix(&pop.population[fatherID], configEntries)
		pop.population[i+populationSize/2] = baby
	}
}

func genPopulation(configEntries map[string]*ConfigEntry) Population {
	var pop Population
	var i uint32
	for i = 0; i < populationSize; i++ {
		pop.population = append(pop.population, genIndividual(configEntries))
	}
	pop.generation = 0

	fmt.Println(pop)
	return pop
}

func (pop *Population) getStdDev(topN int) float64 {
	var avg float64
	var stdDev float64
	var count float64

	for idx, p := range pop.population {
		if idx > topN {
			break
		}
		//Ignore erroneous scores
		if p.score != math.MaxUint32 {
			avg += float64(p.score)
			count += 1.0
		}
	}
	avg /= count

	stdDev = 0.0
	for idx, p := range pop.population {
		if idx > topN {
			break
		}
		//Ignore erroneous scores
		if p.score != math.MaxUint32 {
			dev := math.Abs(float64(p.score) - avg)
			stdDev += dev * dev
		}
	}

	return math.Sqrt(stdDev / count)
}

func magic(clangPath string, configEntries map[string]*ConfigEntry, perfectSource string, popSize uint32, mutRate uint32) {

	populationSize = popSize
	baseMutationRate = mutRate
	mutationRate = baseMutationRate

	pop := genPopulation(configEntries)
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
		sort.Sort(AllIndividuals(pop.population))
		stdDev := pop.getStdDev(len(pop.population) / 2)
		fmt.Printf("Best score for generation %d: %d (stdDev: %f)\n", pop.generation, pop.population[0].score, stdDev)

		//Boost mutations if top population is too homogenous
		if stdDev < minStdDevForMutationBoost {
			mutationRate = 100
		} else {
			mutationRate = baseMutationRate
		}

		//Mix
		pop.mix(configEntries)

		if done == true {
			break
		}
	}

	fmt.Printf("Best configuration has score: %d\n", pop.population[0].score)
	ioutil.WriteFile(".clang-format", []byte(pop.population[0].toClangFormatConfigFile()), os.ModePerm)
	fmt.Println("Written to .clang-format")
}
