package ga

import (
	"math/rand"
)

// StopReason describes why a run terminated.
type StopReason string

const (
	StopReasonGenerationLimit StopReason = "generation_limit"
	StopReasonTargetReached   StopReason = "target_reached"
	StopReasonStagnated       StopReason = "stagnated"
)

type FitnessFunc[T any] func(candidate T) float64

type Generator[T any] func(rng *rand.Rand) T

type Mutator[T any] func(rng *rand.Rand, candidate T) T

type Crossover[T any] func(rng *rand.Rand, a, b T) T

type Selector[T any] func(
	rng *rand.Rand,
	population []Scored[T],
	direction OptimizationDirection,
) T

type Scored[T any] struct {
	Candidate T
	Score     float64
}

type runner[T any] struct {
	cfg      Config[T]
	rng      *rand.Rand
	selector Selector[T]

	ops evolutionOps[T]

	best           Scored[T]
	bestGeneration int
	stagnation     int

	history    []GenerationStats[T]
	stopReason StopReason

	population []T
	generation int
	island     int
}

// Result contains the outcome of a genetic algorithm run.
type Result[T any] struct {
	Best      T
	BestScore float64

	// Generation is the generation in which the best solution was found.
	Generation int

	// History contains statistics for each generation.
	History []GenerationStats[T]

	// StopReason indicates why the run terminated.
	StopReason StopReason
}

// GenerationStats contains statistics for a single generation.
type GenerationStats[T any] struct {
	Island        int
	Generation    int
	BestCandidate T
	BestScore     float64
	AverageScore  float64
	WorstScore    float64

	// Stagnation is the number of generations since the last improvement.
	Stagnation int
}
