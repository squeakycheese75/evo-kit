# evo-kit

A lightweight, type-safe genetic algorithm toolkit for Go.

Evolve solutions to optimization problems using Go generics.

## Why?

Many optimization problems can be expressed as:

- Generate a candidate solution
- Score how good it is
- Mutate and combine candidates
- Repeat until a good solution is found

evo-kit handles the evolutionary loop so you can focus on the problem itself.

## Features

- Go generics support
- Maximization and minimization
- Tournament and roulette selection
- Elitism
- Stagnation detection
- Target score stopping
- Parallel fitness evaluation with configurable worker count
- Generation history and statistics

## Installation

```bash
go get github.com/squeakycheese75/evo-kit
```

## Quick Example

```go
result, err := ga.Run(ga.Config[string]{
    PopulationSize: 100,
    Generations:    500,
    MutationRate:   0.1,
    CrossoverRate:  0.7,
    EliteCount:     2,

    Generate: func(rng *rand.Rand) string {
        return randomString(rng, len(target))
    },

    Fitness: func(candidate string) float64 {
        return score(candidate, target)
    },

    Mutate: mutate,
    Crossover: crossover,
})

if err != nil {
    log.Fatal(err)
}

fmt.Printf("best=%q score=%.0f\n",
    result.Best,
    result.BestScore,
)
```

Example output:

```text
best="hello world" score=11
```

## Included Examples

| Example | Description |
|----------|-------------|
| string_match | Evolve a string until it matches a target |
| knapsack | Maximize value within a weight constraint |
| query_plan | Optimize execution order and parallelism |

Run an example:

```bash
go run ./examples/query_plan
```

Example output:

```text
gen=0 cost=9522.53
gen=31 cost=7645.53

best cost: 7645.53
improvement: 19.71%
stop reason: stagnated
```

## Configuration

```go
cfg := ga.Config[Candidate]{
    PopulationSize: 200,
    Generations:    500,

    MutationRate:   0.1,
    CrossoverRate:  0.7,

    EliteCount:     2,

    Direction:      ga.Minimize,

    MaxStagnation:  50,

    Workers:        runtime.NumCPU(),
}
```

## Parallel Fitness Evaluation

Fitness evaluation can be parallelised with `Workers`.

```go
cfg := ga.Config[Candidate]{
    PopulationSize: 500,
    Generations:    200,
    Workers:        runtime.NumCPU(),

    Generate:  generate,
    Fitness:   expensiveFitness,
    Mutate:    mutate,
    Crossover: crossover,
}
```

A value of `0` or `1` scores candidates sequentially.

When `Workers > 1`, the fitness function must be safe to call concurrently.

Typical scaling for expensive fitness functions:

| Workers | Time |
|----------|------|
| 1 | 1.18s |
| 2 | 590ms |
| 4 | 296ms |
| 8 | 154ms |

## When should I use a genetic algorithm?

Genetic algorithms work well when:

- Exhaustive search is too expensive
- The search space is large
- An approximate solution is acceptable
- The fitness function is easy to evaluate

Examples include scheduling, routing, portfolio optimization, query planning, game AI, and parameter tuning.
