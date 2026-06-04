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

	result, err := Run(cfg)
	if err != nil {
		t.Fatalf("unexpected error")
	}

	if result.BestScore < 10 {
		t.Fatalf("expected best score >= 10, got %.2f", result.BestScore)
	}
}

func TestRunStopsWhenTargetScoreReached(t *testing.T) {
	cfg := Config[int]{
		PopulationSize: 10,
		Generations:    100,
		MutationRate:   1,
		CrossoverRate:  0,
		EliteCount:     1,
		Seed:           42,
		TargetScore:    5,

		Generate: func(rng *rand.Rand) int {
			return 0
		},

		Fitness: func(candidate int) float64 {
			return float64(candidate)
		},

		Mutate: func(rng *rand.Rand, candidate int) int {
			return candidate + 1
		},

		Crossover: func(rng *rand.Rand, a, b int) int {
			return a
		},
	}

	result, err := Run(cfg)
	if err != nil {
		t.Fatalf("unexpected error")
	}

	if result.BestScore < 5 {
		t.Fatalf("expected score >= 5, got %.2f", result.BestScore)
	}

	if result.Generation >= cfg.Generations {
		t.Fatalf("expected early stop, got generation %d", result.Generation)
	}
}

func TestRunCallsOnImprovement(t *testing.T) {
	improvements := 0

	cfg := Config[int]{
		PopulationSize: 10,
		Generations:    20,
		MutationRate:   1,
		CrossoverRate:  0,
		EliteCount:     1,
		Seed:           42,
		TargetScore:    3,

		Generate: func(rng *rand.Rand) int {
			return 0
		},

		Fitness: func(candidate int) float64 {
			return float64(candidate)
		},

		Mutate: func(rng *rand.Rand, candidate int) int {
			return candidate + 1
		},

		Crossover: func(rng *rand.Rand, a, b int) int {
			return a
		},

		OnImprovement: func(stats GenerationStats[int]) {
			improvements++
		},
	}

	Run(cfg)

	if improvements == 0 {
		t.Fatal("expected OnImprovement to be called")
	}
}

func TestRunPreservesEliteCandidate(t *testing.T) {
	cfg := Config[int]{
		PopulationSize: 10,
		Generations:    5,
		MutationRate:   1,
		CrossoverRate:  0,
		EliteCount:     1,
		Seed:           42,

		Generate: func(rng *rand.Rand) int {
			return rng.Intn(10)
		},

		Fitness: func(candidate int) float64 {
			return float64(candidate)
		},

		Mutate: func(rng *rand.Rand, candidate int) int {
			return 0
		},

		Crossover: func(rng *rand.Rand, a, b int) int {
			return a
		},
	}

	result, err := Run(cfg)
	if err != nil {
		t.Fatalf("unexpected error")
	}

	if result.BestScore <= 0 {
		t.Fatalf("expected elite candidate to survive, got %.2f", result.BestScore)
	}
}

func TestRunStopsWhenStagnated(t *testing.T) {
	cfg := validConfig()
	cfg.Generations = 100
	cfg.MaxStagnation = 5
	cfg.TargetScore = 0

	cfg.Generate = func(rng *rand.Rand) int {
		return 1
	}

	cfg.Fitness = func(candidate int) float64 {
		return float64(candidate)
	}

	cfg.Mutate = func(rng *rand.Rand, candidate int) int {
		return candidate
	}

	result, err := Run(cfg)
	if err != nil {
		t.Fatal(err)
	}

	if result.StopReason != StopReasonStagnated {
		t.Fatalf("expected stop reason %q, got %q", StopReasonStagnated, result.StopReason)
	}

	if len(result.History) != cfg.MaxStagnation+1 {
		t.Fatalf("expected %d generations run, got %d", cfg.MaxStagnation+1, len(result.History))
	}
}

func TestRunStopReasonTargetReached(t *testing.T) {
	cfg := validConfig()
	cfg.TargetScore = 5

	cfg.Generate = func(rng *rand.Rand) int {
		return 5
	}

	result, err := Run(cfg)
	if err != nil {
		t.Fatal(err)
	}

	if result.StopReason != StopReasonTargetReached {
		t.Fatalf("expected stop reason %q, got %q", StopReasonTargetReached, result.StopReason)
	}
}
