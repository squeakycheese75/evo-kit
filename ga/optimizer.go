package ga

import (
	"math/rand"
)

const defaultNoOfSelectors = 3

// Run executes a genetic algorithm using cfg and returns the best result found.
func Run[T any](cfg Config[T]) (Result[T], error) {
	if err := cfg.Validate(); err != nil {
		return Result[T]{}, err
	}

	return newRunner(cfg).Run(), nil
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
}

func newRunner[T any](cfg Config[T]) *runner[T] {
	selector := cfg.Select
	if selector == nil {
		selector = TournamentSelector[T](defaultNoOfSelectors)
	}

	ops := defaultEvolutionOps[T]{}

	return &runner[T]{
		cfg:      cfg,
		rng:      rand.New(rand.NewSource(cfg.Seed)),
		selector: selector,
		ops:      ops,
		history:  make([]GenerationStats[T], 0, cfg.Generations),
	}
}

func (r *runner[T]) Run() Result[T] {
	population := r.ops.InitialPopulation(
		r.rng,
		r.cfg.PopulationSize,
		r.cfg.Generate,
	)

	for generation := 0; generation < r.cfg.Generations; generation++ {
		scored := r.scoreAndSort(population)

		stats, improved := r.updateBest(generation, scored)

		r.history = append(r.history, stats)
		r.emitGeneration(stats, improved)

		if r.targetReached() {
			r.stopReason = StopReasonTargetReached
			return r.result()
		}

		if r.stagnated() {
			r.stopReason = StopReasonStagnated
			return r.result()
		}

		population = r.nextGeneration(scored)
	}

	r.stopReason = StopReasonGenerationLimit

	return r.result()
}

func (r *runner[T]) scoreAndSort(population []T) []Scored[T] {
	scored := r.ops.Score(population, r.cfg.Fitness)
	r.ops.Sort(r.cfg.Direction, scored)

	return scored
}

func (r *runner[T]) updateBest(
	generation int,
	scored []Scored[T],
) (GenerationStats[T], bool) {
	averageScore, worstScore := r.ops.Stats(scored)

	improved := generation == 0 || isBetter(r.cfg.Direction, scored[0].Score, r.best.Score)

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
	if r.cfg.OnGeneration != nil {
		r.cfg.OnGeneration(stats)
	}

	if improved && r.cfg.OnImprovement != nil {
		r.cfg.OnImprovement(stats)
	}
}

func (r *runner[T]) targetReached() bool {
	return r.cfg.TargetScore > 0 &&
		isBetterOrEqual(r.cfg.Direction, r.best.Score, r.cfg.TargetScore)
}

func (r *runner[T]) nextGeneration(scored []Scored[T]) []T {
	return r.ops.NextGeneration(
		r.rng,
		r.cfg,
		scored,
		r.selector,
	)
}

func (r *runner[T]) result() Result[T] {
	return Result[T]{
		Best:       r.best.Candidate,
		BestScore:  r.best.Score,
		Generation: r.bestGeneration,
		History:    r.history,
		StopReason: r.stopReason,
	}
}

func (r *runner[T]) stagnated() bool {
	return r.cfg.MaxStagnation > 0 &&
		r.stagnation >= r.cfg.MaxStagnation
}
