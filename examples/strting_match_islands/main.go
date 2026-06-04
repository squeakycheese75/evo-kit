package main

import (
	"fmt"
	"log"
	"math/rand"
	"strings"

	"github.com/squeakycheese75/evo-kit/ga"
)

const target = "genetic islands"

var alphabet = []byte("abcdefghijklmnopqrstuvwxyz ")

func main() {
	cfg := ga.Config[string]{
		PopulationSize: 40,
		Generations:    300,
		MutationRate:   0.15,
		CrossoverRate:  0.7,
		EliteCount:     2,
		Seed:           42,
		TargetScore:    float64(len(target)),
		Direction:      ga.Maximize,
		MaxStagnation:  80,

		Islands: &ga.IslandConfig{
			Count:             4,
			MigrationInterval: 25,
			MigrationCount:    2,
		},

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

		OnImprovement: func(stats ga.GenerationStats[string]) {
			fmt.Printf(
				"island=%d gen=%d score=%.0f candidate=%q\n",
				stats.Island,
				stats.Generation,
				stats.BestScore,
				stats.BestCandidate,
			)
		},
	}

	result, err := ga.Run(cfg)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println()
	fmt.Printf("best: %q\n", result.Best)
	fmt.Printf("score: %.0f/%d\n", result.BestScore, len(target))
	fmt.Printf("best generation: %d\n", result.Generation)
	fmt.Printf("total island generations: %d\n", len(result.History))
	fmt.Printf("stop reason: %s\n", result.StopReason)
	fmt.Printf("migration every=%d count=%d\n",
		cfg.Islands.MigrationInterval,
		cfg.Islands.MigrationCount,
	)

	printIslandSummary(result.History)
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

func printIslandSummary(history []ga.GenerationStats[string]) {
	bestByIsland := make(map[int]ga.GenerationStats[string])

	for _, stats := range history {
		current, ok := bestByIsland[stats.Island]
		if !ok || stats.BestScore > current.BestScore {
			bestByIsland[stats.Island] = stats
		}
	}

	fmt.Println()
	fmt.Println("island summary:")

	for island := 0; island < len(bestByIsland); island++ {
		stats := bestByIsland[island]

		fmt.Printf(
			"island=%d best=%.0f candidate=%q\n",
			island,
			stats.BestScore,
			stats.BestCandidate,
		)
	}

	fmt.Println()
	fmt.Println(strings.Repeat("-", 40))
}
