package ga

import (
	"math/rand"
	"sort"
)

func initialPopulation[T any](
	rng *rand.Rand,
	size int,
	generate Generator[T],
) []T {
	population := make([]T, size)

	for i := range population {
		population[i] = generate(rng)
	}

	return population
}

func populationStats[T any](scored []Scored[T]) (average float64, worst float64) {
	if len(scored) == 0 {
		return 0, 0
	}

	total := 0.0

	for _, item := range scored {
		total += item.Score
	}

	return total / float64(len(scored)), scored[len(scored)-1].Score
}

func scorePopulation[T any](
	population []T,
	fitness FitnessFunc[T],
) []Scored[T] {
	scored := make([]Scored[T], len(population))

	for i, candidate := range population {
		scored[i] = Scored[T]{
			Candidate: candidate,
			Score:     fitness(candidate),
		}
	}

	return scored
}

func nextGeneration[T any](
	rng *rand.Rand,
	cfg Config[T],
	scored []Scored[T],
	selector Selector[T],
) []T {
	next := make([]T, 0, cfg.PopulationSize)

	for i := 0; i < cfg.EliteCount; i++ {
		next = append(next, scored[i].Candidate)
	}

	for len(next) < cfg.PopulationSize {
		a := selector(rng, scored, cfg.Direction)
		b := selector(rng, scored, cfg.Direction)

		child := a

		if rng.Float64() < cfg.CrossoverRate {
			child = cfg.Crossover(rng, a, b)
		}

		if rng.Float64() < cfg.MutationRate {
			child = cfg.Mutate(rng, child)
		}

		next = append(next, child)
	}

	return next
}

func better(direction OptimizationDirection, a, b float64) bool {
	if direction == Minimize {
		return a < b
	}

	return a > b
}

type defaultEvolutionOps[T any] struct{}

func (defaultEvolutionOps[T]) InitialPopulation(
	rng *rand.Rand,
	size int,
	generate Generator[T],
) []T {
	return initialPopulation(rng, size, generate)
}

func (defaultEvolutionOps[T]) Score(
	population []T,
	fitness FitnessFunc[T],
) []Scored[T] {
	return scorePopulation(population, fitness)
}

func (defaultEvolutionOps[T]) Sort(scored []Scored[T]) {
	sort.Slice(scored, func(i, j int) bool {
		return scored[i].Score > scored[j].Score
	})
}

func (defaultEvolutionOps[T]) Stats(scored []Scored[T]) (float64, float64) {
	return populationStats(scored)
}

func (defaultEvolutionOps[T]) NextGeneration(
	rng *rand.Rand,
	cfg Config[T],
	scored []Scored[T],
	selector Selector[T],
) []T {
	return nextGeneration(rng, cfg, scored, selector)
}
