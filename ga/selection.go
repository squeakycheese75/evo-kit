package ga

import "math/rand"

func TournamentSelector[T any](size int) Selector[T] {
	return func(
		rng *rand.Rand,
		population []Scored[T],
		direction OptimizationDirection,
	) T {
		best := population[rng.Intn(len(population))]

		for i := 1; i < size; i++ {
			candidate := population[rng.Intn(len(population))]
			if better(direction, candidate.Score, best.Score) {
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
	return func(
		rng *rand.Rand,
		population []Scored[T],
		direction OptimizationDirection,
	) T {
		weights := make([]float64, len(population))
		totalWeight := 0.0

		minScore := population[0].Score
		maxScore := population[0].Score

		for _, candidate := range population {
			if candidate.Score < minScore {
				minScore = candidate.Score
			}

			if candidate.Score > maxScore {
				maxScore = candidate.Score
			}
		}

		for i, candidate := range population {
			var weight float64

			if direction == Minimize {
				weight = (maxScore - candidate.Score) + 1
			} else {
				weight = (candidate.Score - minScore) + 1
			}

			weights[i] = weight
			totalWeight += weight
		}

		if totalWeight <= 0 {
			return population[rng.Intn(len(population))].Candidate
		}

		selectionPoint := rng.Float64() * totalWeight

		runningTotal := 0.0

		for i, candidate := range population {
			runningTotal += weights[i]

			if runningTotal >= selectionPoint {
				return candidate.Candidate
			}
		}

		return population[len(population)-1].Candidate
	}
}
