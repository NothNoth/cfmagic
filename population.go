package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
)

type Population struct {
	population     []Individual
	generation     uint32
	populationSize uint32
}

func genPopulation(populationSize uint32, configEntries map[string]*ConfigEntry) Population {
	var pop Population
	var i uint32
	for i = 0; i < populationSize; i++ {
		pop.population = append(pop.population, genIndividual(configEntries))
	}
	pop.generation = 0
	pop.populationSize = populationSize
	fmt.Println(pop)
	return pop
}

func (pop Population) String() string {
	var s string

	for idx, p := range pop.population {
		s += fmt.Sprintf("#%d: %s %d\n", idx, p.String(), p.score)
	}

	return s
}

func (pop Population) mix(mutationRate uint32, configEntries map[string]*ConfigEntry) {
	var i uint32
	for i = 0; i < pop.populationSize/2; i++ {
		var motherID int
		fatherID := rand.Intn(int(pop.populationSize / 2))
		for {
			motherID := rand.Intn(int(pop.populationSize / 2))
			if fatherID != motherID {
				break
			}
		}
		baby := pop.population[motherID].mix(&pop.population[fatherID], mutationRate, configEntries)
		pop.population[i+pop.populationSize/2] = baby
	}
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

func (pop *Population) sort() {
	sort.Sort(pop)
}

func (pop Population) Swap(i, j int) {
	pop.population[i], pop.population[j] = pop.population[j], pop.population[i]
}

func (pop Population) Less(i, j int) bool {
	return pop.population[i].score < pop.population[j].score
}

func (pop Population) Len() int {
	return int(pop.populationSize)
}
