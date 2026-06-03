package ga

import (
	"math/rand"
	"testing"
)

func validConfig() Config[int] {
	return Config[int]{
		PopulationSize: 10,
		Generations:    10,
		MutationRate:   0.1,
		CrossoverRate:  0.7,
		EliteCount:     1,
		Seed:           42,

		Generate: func(rng *rand.Rand) int {
			return rng.Intn(10)
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
}

func TestConfigValidateAcceptsValidConfig(t *testing.T) {
	cfg := validConfig()

	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected valid config, got error: %v", err)
	}
}

func TestConfigValidateRequiresPopulationSize(t *testing.T) {
	cfg := validConfig()
	cfg.PopulationSize = 0

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error")
	}
}

func TestConfigValidateRequiresGenerations(t *testing.T) {
	cfg := validConfig()
	cfg.Generations = 0

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error")
	}
}

func TestConfigValidateRejectsInvalidEliteCount(t *testing.T) {
	cfg := validConfig()
	cfg.EliteCount = cfg.PopulationSize + 1

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error")
	}
}

func TestConfigValidateRequiresGenerate(t *testing.T) {
	cfg := validConfig()
	cfg.Generate = nil

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error")
	}
}

func TestConfigValidateRequiresFitness(t *testing.T) {
	cfg := validConfig()
	cfg.Fitness = nil

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error")
	}
}

func TestConfigValidateRequiresMutate(t *testing.T) {
	cfg := validConfig()
	cfg.Mutate = nil

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error")
	}
}

func TestConfigValidateRequiresCrossover(t *testing.T) {
	cfg := validConfig()
	cfg.Crossover = nil

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error")
	}
}
