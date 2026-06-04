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

	if cfg.Islands != nil {
		return newIslandRunner(cfg).Run(), nil
	}

	return newRunner(cfg).Run(), nil
}

type islandRunner[T any] struct {
	cfg     Config[T]
	islands []*runner[T]
}

func (r *islandRunner[T]) step() {
	for _, island := range r.islands {
		if island.stopReason != "" {
			continue
		}

		if island.generation >= island.cfg.Generations {
			island.stopReason = StopReasonGenerationLimit
			continue
		}

		island.step()
	}
}

func (r *islandRunner[T]) Run() Result[T] {
	for generation := 0; generation < r.cfg.Generations; generation++ {
		r.step()

		if r.shouldMigrate(generation) {
			r.migrate()
		}

		if r.allStopped() {
			break
		}
	}

	return r.result()
}

func (r *islandRunner[T]) shouldMigrate(generation int) bool {
	if r.cfg.Islands == nil {
		return false
	}

	if r.cfg.Islands.MigrationInterval <= 0 {
		return false
	}

	if generation == 0 {
		return false
	}

	return generation%r.cfg.Islands.MigrationInterval == 0
}

func (r *islandRunner[T]) result() Result[T] {
	bestIsland := r.islands[0]

	history := make([]GenerationStats[T], 0)

	for _, island := range r.islands {
		history = append(history, island.history...)

		betterScore := isBetter(
			r.cfg.Direction,
			island.best.Score,
			bestIsland.best.Score,
		)

		sameScoreEarlier := island.best.Score == bestIsland.best.Score &&
			island.bestGeneration < bestIsland.bestGeneration

		if betterScore || sameScoreEarlier {
			bestIsland = island
		}
	}

	result := bestIsland.result()
	result.History = history

	return result
}

func (r *islandRunner[T]) allStopped() bool {
	for _, island := range r.islands {
		if island.stopReason == "" {
			return false
		}
	}

	return true
}

func (r *islandRunner[T]) migrate() {
	if r.cfg.Islands == nil {
		return
	}

	if r.cfg.Islands.MigrationCount <= 0 {
		return
	}

	if len(r.islands) < 2 {
		return
	}

	migrants := make([][]T, len(r.islands))

	for i, island := range r.islands {
		scored := island.scoreAndSort(island.population)

		count := r.cfg.Islands.MigrationCount
		if count > len(scored) {
			count = len(scored)
		}

		migrants[i] = make([]T, count)

		for j := 0; j < count; j++ {
			migrants[i][j] = scored[j].Candidate
		}
	}

	for i := range r.islands {
		target := (i + 1) % len(r.islands)

		r.injectMigrants(r.islands[target], migrants[i])
	}
}

func (r *islandRunner[T]) injectMigrants(
	island *runner[T],
	migrants []T,
) {
	if len(migrants) == 0 {
		return
	}

	scored := island.scoreAndSort(island.population)

	for i, migrant := range migrants {
		targetIndex := len(scored) - 1 - i
		if targetIndex < 0 {
			return
		}

		scored[targetIndex].Candidate = migrant
		scored[targetIndex].Score = island.cfg.Fitness(migrant)
	}

	nextPopulation := make([]T, len(scored))

	for i, candidate := range scored {
		nextPopulation[i] = candidate.Candidate
	}

	island.population = nextPopulation
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

func newIslandRunner[T any](cfg Config[T]) *islandRunner[T] {
	islands := make([]*runner[T], cfg.Islands.Count)

	for i := range islands {
		islandCfg := cfg
		islandCfg.Seed += int64(i)
		islandCfg.Islands = nil

		r := newRunner(islandCfg)
		r.island = i
		r.init()

		islands[i] = r
	}

	return &islandRunner[T]{
		cfg:     cfg,
		islands: islands,
	}
}

func (r *runner[T]) Run() Result[T] {
	r.init()

	for r.generation < r.cfg.Generations {
		if !r.step() {
			return r.result()
		}
	}

	r.stopReason = StopReasonGenerationLimit

	return r.result()
}
func (r *runner[T]) init() {
	r.population = r.ops.InitialPopulation(
		r.rng,
		r.cfg.PopulationSize,
		r.cfg.Generate,
	)
}

func (r *runner[T]) step() bool {
	scored := r.scoreAndSort(r.population)

	stats, improved := r.updateBest(r.generation, scored)

	r.history = append(r.history, stats)
	r.emitGeneration(stats, improved)

	if r.targetReached() {
		r.stopReason = StopReasonTargetReached
		return false
	}

	if r.stagnated() {
		r.stopReason = StopReasonStagnated
		return false
	}

	r.population = r.nextGeneration(scored)
	r.generation++

	return true
}

func (r *runner[T]) scoreAndSort(population []T) []Scored[T] {
	scored := r.ops.Score(population, r.cfg.Fitness, r.cfg.Workers)
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
		Island:        r.island,
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
