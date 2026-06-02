package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/squeakycheese75/evo-kit/ga"
)

const target = "hello world"

var alphabet = []byte("abcdefghijklmnopqrstuvwxyz ")

func main() {
	cfg := ga.Config[string]{
		PopulationSize: 100,
		Generations:    500,
		MutationRate:   0.1,
		CrossoverRate:  0.7,
		EliteCount:     2,
		// Seed:           42,
		Seed: time.Now().UnixNano(),

		Generate: func(rng *rand.Rand) string {
			return randomString(rng, len(target))
		},

		Fitness: func(candidate string) float64 {
			return score(candidate, target)
		},

		Mutate: func(rng *rand.Rand, candidate string) string {
			chars := []byte(candidate)
			pos := rng.Intn(len(chars))
			chars[pos] = randomChar(rng)

			return string(chars)
		},

		Crossover: func(rng *rand.Rand, a, b string) string {
			point := rng.Intn(len(a))

			return a[:point] + b[point:]
		},
		TargetScore: float64(len(target)),
		OnImprovement: func(stats ga.GenerationStats[string]) {
			fmt.Printf(
				"gen=%d best=%.0f avg=%.2f worst=%.0f candidate=%q\n",
				stats.Generation,
				stats.BestScore,
				stats.AverageScore,
				stats.WorstScore,
				stats.BestCandidate,
			)
		},
	}

	result := ga.Run(cfg)

	fmt.Printf("best: %q\n", result.Best)
	fmt.Printf("score: %.0f/%d\n", result.BestScore, len(target))
	fmt.Printf("generation: %d\n", result.Generation)
}

func randomString(rng *rand.Rand, length int) string {
	chars := make([]byte, length)

	for i := range chars {
		chars[i] = randomChar(rng)
	}

	return string(chars)
}

func randomChar(rng *rand.Rand) byte {
	return alphabet[rng.Intn(len(alphabet))]
}

func score(candidate string, target string) float64 {
	score := 0.0

	for i := range target {
		if candidate[i] == target[i] {
			score++
		}
	}

	return score
}
