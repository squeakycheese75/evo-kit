package ga

import (
	"fmt"
	"math/rand"
)

const (
	Maximize OptimizationDirection = iota
	Minimize
)

type StopReason string

const (
	StopReasonGenerationLimit StopReason = "generation_limit"
	StopReasonTargetReached   StopReason = "target_reached"
	StopReasonStagnated       StopReason = "stagnated"
)

type PopulationInitializer[T any] interface {
	InitialPopulation(rng *rand.Rand, size int, generate Generator[T]) []T
}

type PopulationScorer[T any] interface {
	Score(population []T, fitness FitnessFunc[T]) []Scored[T]
}

type PopulationSorter[T any] interface {
	Sort(scored []Scored[T])
}

type PopulationStatsCalculator[T any] interface {
	Stats(scored []Scored[T]) (average float64, worst float64)
}

type NextGenerationBuilder[T any] interface {
	NextGeneration(
		rng *rand.Rand,
		cfg Config[T],
		scored []Scored[T],
		selector Selector[T],
	) []T
}

type FitnessFunc[T any] func(candidate T) float64

type Generator[T any] func(rng *rand.Rand) T

type Mutator[T any] func(rng *rand.Rand, candidate T) T

type Crossover[T any] func(rng *rand.Rand, a, b T) T

type Selector[T any] func(
	rng *rand.Rand,
	population []Scored[T],
	direction OptimizationDirection,
) T

type OptimizationDirection int

type Config[T any] struct {
	PopulationSize int
	Generations    int
	MutationRate   float64
	CrossoverRate  float64
	EliteCount     int
	Seed           int64
	MaxStagnation  int

	Generate  Generator[T]
	Fitness   FitnessFunc[T]
	Mutate    Mutator[T]
	Crossover Crossover[T]
	Select    Selector[T]
	Direction OptimizationDirection

	TargetScore   float64
	OnGeneration  func(stats GenerationStats[T])
	OnImprovement func(stats GenerationStats[T])
}

func (cfg Config[T]) Validate() error {
	if cfg.PopulationSize <= 0 {
		return fmt.Errorf("population size must be greater than zero")
	}

	if cfg.Generations <= 0 {
		return fmt.Errorf("generations must be greater than zero")
	}

	if cfg.EliteCount < 0 || cfg.EliteCount > cfg.PopulationSize {
		return fmt.Errorf("elite count must be between 0 and population size")
	}

	if cfg.Generate == nil {
		return fmt.Errorf("generate function is required")
	}

	if cfg.Fitness == nil {
		return fmt.Errorf("fitness function is required")
	}

	if cfg.Mutate == nil {
		return fmt.Errorf("mutate function is required")
	}

	if cfg.Crossover == nil {
		return fmt.Errorf("crossover function is required")
	}

	return nil
}

type Scored[T any] struct {
	Candidate T
	Score     float64
}

type Result[T any] struct {
	Best       T
	BestScore  float64
	Generation int
	History    []GenerationStats[T]
	StopReason StopReason
}

type GenerationStats[T any] struct {
	Generation    int
	BestCandidate T
	BestScore     float64
	AverageScore  float64
	WorstScore    float64
	Stagnation    int
}
