package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/squeakycheese75/evo-kit/ga"
)

type Item struct {
	Name   string
	Value  float64
	Weight float64
}

type Solution []bool

var items = []Item{
	{Name: "Laptop", Value: 500, Weight: 5},
	{Name: "Camera", Value: 300, Weight: 3},
	{Name: "Headphones", Value: 150, Weight: 1},
	{Name: "Coffee", Value: 60, Weight: 2},
	{Name: "Jacket", Value: 200, Weight: 4},
	{Name: "Book", Value: 80, Weight: 2},
	{Name: "Shoes", Value: 180, Weight: 3},
	{Name: "Powerbank", Value: 120, Weight: 1},
}

const maxWeight = 10.0
const overweightPenaltyPerKg = 100.0

func main() {
	cfg := ga.Config[Solution]{
		PopulationSize: 100,
		Generations:    300,
		MutationRate:   0.1,
		CrossoverRate:  0.7,
		EliteCount:     2,
		Seed:           time.Now().UnixNano(),

		Generate: func(rng *rand.Rand) Solution {
			return randomSolution(rng)
		},

		Fitness: func(candidate Solution) float64 {
			value, weight := totals(candidate)

			if weight > maxWeight {
				return value - ((weight - maxWeight) * overweightPenaltyPerKg)
			}

			return value
		},

		Mutate: func(rng *rand.Rand, candidate Solution) Solution {
			next := clone(candidate)
			pos := rng.Intn(len(next))
			next[pos] = !next[pos]

			return next
		},

		Crossover: func(rng *rand.Rand, a, b Solution) Solution {
			point := rng.Intn(len(a))

			child := make(Solution, len(a))
			copy(child[:point], a[:point])
			copy(child[point:], b[point:])

			return child
		},
		Select: ga.TournamentSelector[Solution](3),

		OnImprovement: func(stats ga.GenerationStats[Solution]) {
			value, weight := totals(stats.BestCandidate)

			fmt.Printf(
				"gen=%d fitness=%.2f value=%.2f weight=%.2f selected=%v\n",
				stats.Generation,
				stats.BestScore,
				value,
				weight,
				selectedNames(stats.BestCandidate),
			)
		},
	}

	result, err := ga.Run(cfg)
	if err != nil {
		panic("unexpected error")
	}

	value, weight := totals(result.Best)

	fmt.Println()
	fmt.Printf("best fitness: %.2f\n", result.BestScore)
	fmt.Printf("value: %.2f\n", value)
	fmt.Printf("weight: %.2f / %.2f\n", weight, maxWeight)
	fmt.Printf("items: %v\n", selectedNames(result.Best))
	fmt.Printf("generation: %d\n", result.Generation)
}

func randomSolution(rng *rand.Rand) Solution {
	solution := make(Solution, len(items))

	for i := range solution {
		solution[i] = rng.Intn(2) == 1
	}

	return solution
}

func totals(solution Solution) (value float64, weight float64) {
	for i, selected := range solution {
		if !selected {
			continue
		}

		value += items[i].Value
		weight += items[i].Weight
	}

	return value, weight
}

func selectedNames(solution Solution) []string {
	names := make([]string, 0)

	for i, selected := range solution {
		if !selected {
			continue
		}

		names = append(names, items[i].Name)
	}

	return names
}

func clone(solution Solution) Solution {
	next := make(Solution, len(solution))
	copy(next, solution)

	return next
}
