package ga

import "math/rand"

func tournament[T any](rng *rand.Rand, population []Scored[T], size int) T {
	best := population[rng.Intn(len(population))]

	for i := 1; i < size; i++ {
		candidate := population[rng.Intn(len(population))]
		if candidate.Score > best.Score {
			best = candidate
		}
	}

	return best.Candidate
}
