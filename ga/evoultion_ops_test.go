package ga

import (
	"math/rand"
	"reflect"
	"testing"
)

func TestScorePopulationParallelMatchesSequential(t *testing.T) {
	population := []int{1, 2, 3, 4, 5}

	fitness := func(candidate int) float64 {
		return float64(candidate * candidate)
	}

	sequential := scorePopulation(population, fitness)
	parallel := scorePopulationParallel(population, fitness, 2)

	if !reflect.DeepEqual(sequential, parallel) {
		t.Fatalf(
			"expected parallel scoring to match sequential\nsequential=%+v\nparallel=%+v",
			sequential,
			parallel,
		)
	}
}

func TestConfigValidateRejectsNegativeWorkers(t *testing.T) {
	cfg := Config[int]{
		PopulationSize: 100,
		Generations:    100,
		MutationRate:   0.1,
		CrossoverRate:  0.7,
		EliteCount:     2,
		Workers:        -1,

		Generate: func(rng *rand.Rand) int {
			return 0
		},

		Fitness: func(candidate int) float64 {
			return float64(candidate)
		},

		Mutate: func(rng *rand.Rand, candidate int) int {
			return candidate
		},

		Crossover: func(rng *rand.Rand, a, b int) int {
			return a
		},
	}

	err := cfg.Validate()

	if err == nil {
		t.Fatal("expected validation error")
	}

	expected := "workers must be greater than or equal to zero"

	if err.Error() != expected {
		t.Fatalf("expected %q, got %q", expected, err.Error())
	}
}

func TestScorePopulationParallelWithMoreWorkersThanPopulation(t *testing.T) {
	population := []int{1, 2, 3}

	fitness := func(candidate int) float64 {
		return float64(candidate * candidate)
	}

	sequential := scorePopulation(population, fitness)
	parallel := scorePopulationParallel(population, fitness, 100)

	if !reflect.DeepEqual(sequential, parallel) {
		t.Fatalf(
			"expected parallel scoring to match sequential\nsequential=%+v\nparallel=%+v",
			sequential,
			parallel,
		)
	}
}
