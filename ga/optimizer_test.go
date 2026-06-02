package ga

import (
	"math/rand"
	"testing"
)

func TestRunFindsBestCandidate(t *testing.T) {
	cfg := Config[int]{
		PopulationSize: 20,
		Generations:    50,
		MutationRate:   0.2,
		CrossoverRate:  0.7,
		EliteCount:     1,
		Seed:           42,
		TargetScore:    10,

		Generate: func(rng *rand.Rand) int {
			return rng.Intn(10)
		},

		Fitness: func(candidate int) float64 {
			return float64(candidate)
		},

		Mutate: func(rng *rand.Rand, candidate int) int {
			if candidate >= 10 {
				return candidate
			}

			return candidate + 1
		},

		Crossover: func(rng *rand.Rand, a, b int) int {
			if a > b {
				return a
			}

			return b
		},
	}

	result := Run(cfg)

	if result.BestScore < 10 {
		t.Fatalf("expected best score >= 10, got %.2f", result.BestScore)
	}
}
