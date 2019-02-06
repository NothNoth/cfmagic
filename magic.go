package main

import (
	"fmt"
	"sort"
)

const populationSize = 100

type ConfigValue struct {
	entry string
	value string
}

type Individual struct {
	configValues []ConfigValue
	score        uint32
}

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

func genPopulation(configEntries map[string]*ConfigEntry) Population {
	var pop Population
	for i := 0; i < populationSize; i++ {
		pop.population = append(pop.population, genIndividual(configEntries))
	}
	pop.generation = 0
	return pop
}

func magic(clangPath string, configEntries map[string]*ConfigEntry, perfectSource string) {

	pop := genPopulation(configEntries)
	fmt.Println("Population ready")

	//Eval everyone
	for i := 0; i < len(pop.population); i++ {
		err := pop.population[i].UpdateScore(clangPath, perfectSource)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	//Sort
	sort.Sort(AllIndividuals(pop.population))
	fmt.Println(pop)
	//fmt.Println("All done")
}
