package ga

import "math/rand"

type FitnessFunc[T any] func(candidate T) float64

type Generator[T any] func(rng *rand.Rand) T

type Mutator[T any] func(rng *rand.Rand, candidate T) T

type Crossover[T any] func(rng *rand.Rand, a, b T) T

type Config[T any] struct {
	PopulationSize int
	Generations    int
	MutationRate   float64
	CrossoverRate  float64
	EliteCount     int
	Seed           int64

	Generate      Generator[T]
	Fitness       FitnessFunc[T]
	Mutate        Mutator[T]
	Crossover     Crossover[T]
	TargetScore   float64
	OnGeneration  func(stats GenerationStats[T])
	OnImprovement func(stats GenerationStats[T])
}

type Scored[T any] struct {
	Candidate T
	Score     float64
}

type Result[T any] struct {
	Best       T
	BestScore  float64
	Generation int
}

type GenerationStats[T any] struct {
	Generation    int
	BestCandidate T
	BestScore     float64
	AverageScore  float64
	WorstScore    float64
	Stagnation    int
}
