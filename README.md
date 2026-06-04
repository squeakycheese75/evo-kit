# evo-kit

A lightweight, type-safe genetic algorithm toolkit for Go.

## Features

- Generic candidate types using Go generics
- Tournament and roulette-wheel selection
- Maximization and minimization support
- Elitism
- Target score stopping
- Stagnation-based stopping
- Generation statistics and history tracking
- Extensible architecture for custom evolutionary operators

## Installation

```bash
go get github.com/squeakycheese75/evo-kit
```

## Quick Start

```go
result, err := ga.Run(cfg)
if err != nil {
    panic(err)
}

fmt.Printf("best score: %.2f\n", result.BestScore)
```

## Examples

### String Matching

```bash
go run ./examples/string_match
```

### Knapsack Optimization

```bash
go run ./examples/knapsack
```

### Query Plan Optimization

```bash
go run ./examples/query_plan
```
