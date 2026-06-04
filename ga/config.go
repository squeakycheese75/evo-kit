package ga

import "fmt"

// OptimizationDirection determines whether scores are maximized or minimized.
type OptimizationDirection int

const (
	// Maximize prefers larger scores.
	Maximize OptimizationDirection = iota

	// Minimize prefers smaller scores.
	Minimize
)

// Config controls the behaviour of a genetic algorithm run.
type Config[T any] struct {
	// PopulationSize is the number of candidates maintained in each generation.
	PopulationSize int

	// Generations is the maximum number of generations to run.
	Generations int

	// MutationRate is the probability that a child candidate will be mutated.
	MutationRate float64

	// CrossoverRate is the probability that a child candidate will be created
	// by crossover instead of cloning a selected parent.
	CrossoverRate float64

	// EliteCount is the number of top-scoring candidates preserved unchanged
	// into the next generation.
	EliteCount int

	// Seed controls random number generation.
	Seed int64

	// MaxStagnation stops the run after N generations without improvement.
	// A value of zero disables stagnation-based stopping.
	MaxStagnation int

	Generate  Generator[T]
	Fitness   FitnessFunc[T]
	Mutate    Mutator[T]
	Crossover Crossover[T]
	Select    Selector[T]

	// Direction determines whether scores should be maximized or minimized.
	Direction OptimizationDirection

	// TargetScore causes the run to stop once reached.
	TargetScore float64

	// OnGeneration is called after each generation.
	OnGeneration func(stats GenerationStats[T])

	// OnImprovement is called whenever a new best candidate is found.
	OnImprovement func(stats GenerationStats[T])

	// Workers controls how many goroutines are used to score candidates.
	// A value of zero or one scores candidates sequentially.
	// Fitness functions must be safe for concurrent use when Workers is greater than one.
	Workers int
}

// Validate checks whether the config contains the required values.
func (cfg Config[T]) Validate() error {
	if cfg.PopulationSize <= 0 {
		return fmt.Errorf("population size must be greater than zero")
	}

	if cfg.Generations <= 0 {
		return fmt.Errorf("generations must be greater than zero")
	}

	if cfg.EliteCount < 0 || cfg.EliteCount > cfg.PopulationSize {
		return fmt.Errorf("elite count must be between 0 and population size")
	}

	if cfg.Generate == nil {
		return fmt.Errorf("generate function is required")
	}

	if cfg.Fitness == nil {
		return fmt.Errorf("fitness function is required")
	}

	if cfg.Mutate == nil {
		return fmt.Errorf("mutate function is required")
	}

	if cfg.Crossover == nil {
		return fmt.Errorf("crossover function is required")
	}

	if cfg.Workers < 0 {
		return fmt.Errorf("workers must be greater than or equal to zero")
	}

	return nil
}
