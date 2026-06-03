package ga

import "math/rand"

func TournamentSelector[T any](size int) Selector[T] {
	return func(rng *rand.Rand, population []Scored[T]) T {
		best := population[rng.Intn(len(population))]

		for i := 1; i < size; i++ {
			candidate := population[rng.Intn(len(population))]
			if candidate.Score > best.Score {
				best = candidate
			}
		}

		return best.Candidate
	}
}

// RouletteSelector assumes all fitness scores are non-negative.
// For fitness functions that may return negative values,
// use TournamentSelector or normalize scores first.
func RouletteSelector[T any]() Selector[T] {
	return func(rng *rand.Rand, population []Scored[T]) T {
		totalFitness := 0.0

		for _, candidate := range population {
			totalFitness += candidate.Score
		}

		if totalFitness <= 0 {
			return population[rng.Intn(len(population))].Candidate
		}

		selectionPoint := rng.Float64() * totalFitness

		runningTotal := 0.0

		for _, candidate := range population {
			runningTotal += candidate.Score

			if runningTotal >= selectionPoint {
				return candidate.Candidate
			}
		}

		// Handle floating point edge cases
		return population[len(population)-1].Candidate
	}
}
