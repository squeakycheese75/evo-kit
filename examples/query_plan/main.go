package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/squeakycheese75/evo-kit/ga"
)

type Dataset struct {
	Name       string
	Rows       int
	FetchMs    int
	Filterable bool
}

type QueryPlan struct {
	Order       []Dataset
	Parallelism int
}

var datasets = []Dataset{
	{Name: "google_sheets_leads", Rows: 10_000, FetchMs: 500, Filterable: true},
	{Name: "stripe_sales", Rows: 2_000, FetchMs: 300, Filterable: true},
	{Name: "geoip_api", Rows: 50_000, FetchMs: 1200, Filterable: false},
	{Name: "crm_accounts", Rows: 5_000, FetchMs: 700, Filterable: true},
	{Name: "support_tickets", Rows: 25_000, FetchMs: 900, Filterable: true},
	{Name: "billing_events", Rows: 80_000, FetchMs: 1500, Filterable: true},
	{Name: "email_events", Rows: 120_000, FetchMs: 1800, Filterable: false},
	{Name: "product_catalog", Rows: 1_000, FetchMs: 200, Filterable: true},
}

func main() {
	cfg := ga.Config[QueryPlan]{
		PopulationSize: 20,
		Generations:    200,
		MutationRate:   0.2,
		CrossoverRate:  0.7,
		EliteCount:     2,
		Seed:           time.Now().UnixNano(),
		Select:         ga.TournamentSelector[QueryPlan](3),

		Generate: func(rng *rand.Rand) QueryPlan {
			return randomPlan(rng)
		},

		Fitness: func(plan QueryPlan) float64 {
			return estimatedCost(plan)
		},

		Mutate: func(rng *rand.Rand, plan QueryPlan) QueryPlan {
			next := clonePlan(plan)

			switch rng.Intn(2) {
			case 0:
				swapDatasets(rng, next.Order)
			case 1:
				next.Parallelism = 1 + rng.Intn(4)
			}

			return next
		},
		Direction: ga.Minimize,

		Crossover: func(rng *rand.Rand, a, b QueryPlan) QueryPlan {
			return crossoverPlan(rng, a, b)
		},

		OnImprovement: func(stats ga.GenerationStats[QueryPlan]) {
			cost := estimatedCost(stats.BestCandidate)

			fmt.Printf(
				"gen=%d cost=%.2f parallelism=%d order=%v\n",
				stats.Generation,
				cost,
				stats.BestCandidate.Parallelism,
				datasetNames(stats.BestCandidate.Order),
			)
		},
	}

	result, err := ga.Run(cfg)
	if err != nil {
		panic(err)
	}

	fmt.Println()
	fmt.Printf("best cost: %.2f\n", estimatedCost(result.Best))
	fmt.Printf("parallelism: %d\n", result.Best.Parallelism)
	fmt.Printf("order: %v\n", datasetNames(result.Best.Order))
	fmt.Printf("generation: %d\n", result.Generation)

	fmt.Println("explanation:")

	for _, line := range explainPlan(result.Best) {
		fmt.Printf("- %s\n", line)
	}
}

func randomPlan(rng *rand.Rand) QueryPlan {
	order := make([]Dataset, len(datasets))
	copy(order, datasets)

	rng.Shuffle(len(order), func(i, j int) {
		order[i], order[j] = order[j], order[i]
	})

	return QueryPlan{
		Order:       order,
		Parallelism: 1 + rng.Intn(4),
	}
}

func estimatedCost(plan QueryPlan) float64 {
	cost := 0.0
	remainingRows := 100_000.0

	for i, dataset := range plan.Order {
		fetchCost := float64(dataset.FetchMs)
		rowCost := float64(dataset.Rows) * 0.01

		if dataset.Filterable && i > 0 {
			rowCost *= 0.5
		}

		joinCost := remainingRows * float64(dataset.Rows) * 0.000001

		cost += fetchCost + rowCost + joinCost

		if dataset.Filterable {
			remainingRows *= 0.6
		} else {
			remainingRows *= 0.9
		}
	}

	parallelismBenefit := 1 + float64(plan.Parallelism)*0.15
	parallelismPenalty := float64(plan.Parallelism-1) * 80

	return (cost / parallelismBenefit) + parallelismPenalty
}

func swapDatasets(rng *rand.Rand, order []Dataset) {
	i := rng.Intn(len(order))
	j := rng.Intn(len(order))

	order[i], order[j] = order[j], order[i]
}

func crossoverPlan(rng *rand.Rand, a, b QueryPlan) QueryPlan {
	size := len(a.Order)
	start := rng.Intn(size)
	end := start + rng.Intn(size-start)

	childOrder := make([]Dataset, 0, size)
	used := make(map[string]bool, size)

	for i := start; i <= end; i++ {
		childOrder = append(childOrder, a.Order[i])
		used[a.Order[i].Name] = true
	}

	for _, dataset := range b.Order {
		if used[dataset.Name] {
			continue
		}

		childOrder = append(childOrder, dataset)
	}

	return QueryPlan{
		Order:       childOrder,
		Parallelism: a.Parallelism,
	}
}

func clonePlan(plan QueryPlan) QueryPlan {
	order := make([]Dataset, len(plan.Order))
	copy(order, plan.Order)

	return QueryPlan{
		Order:       order,
		Parallelism: plan.Parallelism,
	}
}

func datasetNames(datasets []Dataset) []string {
	names := make([]string, 0, len(datasets))

	for _, dataset := range datasets {
		names = append(names, dataset.Name)
	}

	return names
}

func explainPlan(plan QueryPlan) []string {
	explanations := make([]string, 0)

	for i, dataset := range plan.Order {
		if dataset.Filterable && i < 4 {
			explanations = append(
				explanations,
				fmt.Sprintf("%s is filterable and appears early", dataset.Name),
			)
		}

		if !dataset.Filterable && i >= len(plan.Order)-2 {
			explanations = append(
				explanations,
				fmt.Sprintf("%s is not filterable and appears late", dataset.Name),
			)
		}
	}

	return explanations
}
