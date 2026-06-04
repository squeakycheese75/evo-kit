package ga

import (
	"math/rand"
	"testing"
	"time"
)

const benchmarkTarget = "hello world"

var benchmarkAlphabet = []byte("abcdefghijklmnopqrstuvwxyz ")

func BenchmarkStringMatchTournament2(b *testing.B) {
	benchmarkStringMatch(b, TournamentSelector[string](2))
}

func BenchmarkStringMatchTournament3(b *testing.B) {
	benchmarkStringMatch(b, TournamentSelector[string](3))
}

func BenchmarkStringMatchTournament5(b *testing.B) {
	benchmarkStringMatch(b, TournamentSelector[string](5))
}

func BenchmarkStringMatchRoulette(b *testing.B) {
	benchmarkStringMatch(b, RouletteSelector[string]())
}

func BenchmarkStringMatchPopulation50(b *testing.B) {
	benchmarkStringMatchPopulation(
		b,
		50,
		TournamentSelector[string](3),
	)
}

func BenchmarkStringMatchPopulation100(b *testing.B) {
	benchmarkStringMatchPopulation(
		b,
		100,
		TournamentSelector[string](3),
	)
}

func BenchmarkStringMatchPopulation500(b *testing.B) {
	benchmarkStringMatchPopulation(
		b,
		500,
		TournamentSelector[string](3),
	)
}

func BenchmarkStringMatchPopulation1000(b *testing.B) {
	benchmarkStringMatchPopulation(
		b,
		1000,
		TournamentSelector[string](3),
	)
}

func benchmarkStringMatchPopulation(
	b *testing.B,
	populationSize int,
	selector Selector[string],
) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cfg := Config[string]{
			PopulationSize: populationSize,
			Generations:    500,
			MutationRate:   0.1,
			CrossoverRate:  0.7,
			EliteCount:     2,
			Seed:           int64(i),
			TargetScore:    float64(len(benchmarkTarget)),
			Select:         selector,

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

			Crossover: func(rng *rand.Rand, a, c string) string {
				point := rng.Intn(len(a))

				return a[:point] + c[point:]
			},
		}

		if _, err := Run(cfg); err != nil {
			b.Fatal(err)
		}
	}
}

func benchmarkStringMatch(b *testing.B, selector Selector[string]) {
	for i := 0; i < b.N; i++ {
		cfg := Config[string]{
			PopulationSize: 100,
			Generations:    500,
			MutationRate:   0.1,
			CrossoverRate:  0.7,
			EliteCount:     2,
			Seed:           int64(i),
			TargetScore:    float64(len(benchmarkTarget)),
			Select:         selector,

			Generate: func(rng *rand.Rand) string {
				return benchmarkRandomString(rng, len(benchmarkTarget))
			},

			Fitness: func(candidate string) float64 {
				return benchmarkScore(candidate, benchmarkTarget)
			},

			Mutate: func(rng *rand.Rand, candidate string) string {
				chars := []byte(candidate)
				pos := rng.Intn(len(chars))
				chars[pos] = benchmarkRandomChar(rng)

				return string(chars)
			},

			Crossover: func(rng *rand.Rand, a, c string) string {
				point := rng.Intn(len(a))

				return a[:point] + c[point:]
			},
		}

		if _, err := Run(cfg); err != nil {
			b.Fatal(err)
		}
	}
}

func benchmarkRandomString(rng *rand.Rand, length int) string {
	chars := make([]byte, length)

	for i := range chars {
		chars[i] = benchmarkRandomChar(rng)
	}

	return string(chars)
}

func benchmarkRandomChar(rng *rand.Rand) byte {
	return benchmarkAlphabet[rng.Intn(len(benchmarkAlphabet))]
}

func benchmarkScore(candidate string, target string) float64 {
	score := 0.0

	for i := range target {
		if candidate[i] == target[i] {
			score++
		}
	}

	return score
}

func BenchmarkWorkers1(b *testing.B) {
	benchmarkWorkers(b, 1)
}

func BenchmarkWorkers2(b *testing.B) {
	benchmarkWorkers(b, 2)
}

func BenchmarkWorkers4(b *testing.B) {
	benchmarkWorkers(b, 4)
}

func BenchmarkWorkers8(b *testing.B) {
	benchmarkWorkers(b, 8)
}

func benchmarkWorkers(b *testing.B, workers int) {
	target := "hello world"

	cfg := Config[string]{
		PopulationSize: 100,
		Generations:    100,
		MutationRate:   0.1,
		CrossoverRate:  0.7,
		EliteCount:     2,
		Workers:        workers,
		Direction:      Maximize,
		Seed:           42,

		Generate: func(rng *rand.Rand) string {
			buf := make([]byte, len(target))

			for i := range buf {
				buf[i] = byte('a' + rng.Intn(26))
			}

			return string(buf)
		},

		Fitness: func(candidate string) float64 {
			time.Sleep(100 * time.Microsecond)

			score := 0

			for i := range target {
				if candidate[i] == target[i] {
					score++
				}
			}

			return float64(score)
		},

		Mutate: func(rng *rand.Rand, candidate string) string {
			buf := []byte(candidate)

			idx := rng.Intn(len(buf))
			buf[idx] = byte('a' + rng.Intn(26))

			return string(buf)
		},

		Crossover: func(rng *rand.Rand, a, b string) string {
			point := rng.Intn(len(a))

			return a[:point] + b[point:]
		},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := Run(cfg)

		if err != nil {
			b.Fatal(err)
		}
	}
}
