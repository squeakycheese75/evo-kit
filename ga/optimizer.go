package ga

import (
	"math/rand"
	"sort"
)

func Run[T any](cfg Config[T]) Result[T] {
	rng := rand.New(rand.NewSource(cfg.Seed))

	population := initialPopulation(rng, cfg.PopulationSize, cfg.Generate)

	var best Scored[T]
	bestGeneration := 0
	stagnation := 0

	for generation := 0; generation < cfg.Generations; generation++ {
		scored := scorePopulation(population, cfg.Fitness)

		sort.Slice(scored, func(i, j int) bool {
			return scored[i].Score > scored[j].Score
		})

		averageScore, worstScore := populationStats(scored)

		improved := generation == 0 || scored[0].Score > best.Score

		if improved {
			best = scored[0]
			bestGeneration = generation
			stagnation = 0

			stats := GenerationStats[T]{
				Generation:    generation,
				BestCandidate: best.Candidate,
				BestScore:     best.Score,
				AverageScore:  averageScore,
				WorstScore:    worstScore,
				Stagnation:    stagnation,
			}

			if cfg.OnGeneration != nil {
				cfg.OnGeneration(stats)
			}

			if cfg.OnImprovement != nil {
				cfg.OnImprovement(stats)
			}

			if cfg.TargetScore > 0 && best.Score >= cfg.TargetScore {
				return Result[T]{
					Best:       best.Candidate,
					BestScore:  best.Score,
					Generation: bestGeneration,
				}
			}
		} else {
			stagnation++
		}

		population = nextGeneration(rng, cfg, scored)
	}

	return Result[T]{
		Best:       best.Candidate,
		BestScore:  best.Score,
		Generation: bestGeneration,
	}
}

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
) []T {
	next := make([]T, 0, cfg.PopulationSize)

	for i := 0; i < cfg.EliteCount; i++ {
		next = append(next, scored[i].Candidate)
	}

	for len(next) < cfg.PopulationSize {
		a := tournament(rng, scored, 3)
		b := tournament(rng, scored, 3)

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
