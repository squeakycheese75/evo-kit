package ga

import (
	"math/rand"
	"testing"
)

func TestPopulationSizes(t *testing.T) {
	sizes := []int{
		50,
		100,
		500,
		1000,
	}

	for _, size := range sizes {
		var totalGenerations int

		const runs = 100

		for run := 0; run < runs; run++ {
			result, err := Run(Config[string]{
				PopulationSize: size,
				Generations:    500,
				MutationRate:   0.1,
				CrossoverRate:  0.7,
				EliteCount:     2,
				Seed:           int64(run),
				TargetScore:    float64(len(benchmarkTarget)),
				Select:         TournamentSelector[string](3),

				Generate: func(rng *rand.Rand) string {
					return benchmarkRandomString(
						rng,
						len(benchmarkTarget),
					)
				},

				Fitness: func(candidate string) float64 {
					return benchmarkScore(
						candidate,
						benchmarkTarget,
					)
				},

				Mutate: func(rng *rand.Rand, candidate string) string {
					chars := []byte(candidate)

					pos := rng.Intn(len(chars))
					chars[pos] = benchmarkRandomChar(rng)

					return string(chars)
				},

				Crossover: func(rng *rand.Rand, a, b string) string {
					point := rng.Intn(len(a))

					return a[:point] + b[point:]
				},
			})
			if err != nil {
				t.Fatal(err)
			}

			totalGenerations += result.Generation
		}

		average := float64(totalGenerations) / runs

		t.Logf("population=%d avg_generation=%.2f", size, average)
	}
}
