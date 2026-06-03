package ga

import (
	"math/rand"
	"sort"
)

const numberOfSelectors = 3

func Run[T any](cfg Config[T]) (Result[T], error) {
	runner := newRunner(cfg)

	if err := cfg.Validate(); err != nil {
		return Result[T]{}, err
	}
	return runner.Run(), nil
}

type runner[T any] struct {
	cfg      Config[T]
	rng      *rand.Rand
	selector Selector[T]

	initializer PopulationInitializer[T]
	scorer      PopulationScorer[T]
	sorter      PopulationSorter[T]
	stats       PopulationStatsCalculator[T]
	next        NextGenerationBuilder[T]

	best           Scored[T]
	bestGeneration int
	stagnation     int
}

func newRunner[T any](cfg Config[T]) *runner[T] {
	selector := cfg.Select
	if selector == nil {
		selector = TournamentSelector[T](numberOfSelectors)
	}

	ops := defaultEvolutionOps[T]{}

	return &runner[T]{
		cfg:         cfg,
		rng:         rand.New(rand.NewSource(cfg.Seed)),
		selector:    selector,
		initializer: ops,
		scorer:      ops,
		sorter:      ops,
		stats:       ops,
		next:        ops,
	}
}

func (r *runner[T]) Run() Result[T] {
	population := initialPopulation(r.rng, r.cfg.PopulationSize, r.cfg.Generate)

	for generation := 0; generation < r.cfg.Generations; generation++ {
		scored := r.scoreAndSort(population)

		stats, improved := r.updateBest(generation, scored)

		r.emitGeneration(stats, improved)

		if r.targetReached() {
			return r.result()
		}

		population = r.nextGeneration(scored)
	}

	return r.result()
}

func (r *runner[T]) scoreAndSort(population []T) []Scored[T] {
	scored := scorePopulation(population, r.cfg.Fitness)

	sort.Slice(scored, func(i, j int) bool {
		return better(r.cfg.Direction, scored[i].Score, scored[j].Score)
	})

	return scored
}

func (r *runner[T]) updateBest(
	generation int,
	scored []Scored[T],
) (GenerationStats[T], bool) {
	averageScore, worstScore := populationStats(scored)

	// improved := generation == 0 || scored[0].Score > r.best.Score
	improved := generation == 0 || better(r.cfg.Direction, scored[0].Score, r.best.Score)

	if improved {
		r.best = scored[0]
		r.bestGeneration = generation
		r.stagnation = 0
	} else {
		r.stagnation++
	}

	return GenerationStats[T]{
		Generation:    generation,
		BestCandidate: r.best.Candidate,
		BestScore:     r.best.Score,
		AverageScore:  averageScore,
		WorstScore:    worstScore,
		Stagnation:    r.stagnation,
	}, improved
}

func (r *runner[T]) emitGeneration(stats GenerationStats[T], improved bool) {
	if !improved {
		return
	}

	if r.cfg.OnGeneration != nil {
		r.cfg.OnGeneration(stats)
	}

	if r.cfg.OnImprovement != nil {
		r.cfg.OnImprovement(stats)
	}
}

func (r *runner[T]) targetReached() bool {
	return r.cfg.TargetScore > 0 && r.best.Score >= r.cfg.TargetScore
}

func (r *runner[T]) nextGeneration(scored []Scored[T]) []T {
	return nextGeneration(r.rng, r.cfg, scored, r.selector)
}

func (r *runner[T]) result() Result[T] {
	return Result[T]{
		Best:       r.best.Candidate,
		BestScore:  r.best.Score,
		Generation: r.bestGeneration,
	}
}
