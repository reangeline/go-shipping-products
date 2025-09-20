package order

import (
	"errors"
	"math"
	"sort"
)

// To calculate the best package combination
type PackCalculator interface {
	Calculate(quantity int, packs []Pack) (Combination, error)
}

type packCalculator struct{}

func NewPackCalculator() PackCalculator { return &packCalculator{} }

func (pc *packCalculator) Calculate(quantity int, packs []Pack) (Combination, error) {
	// Basic validations
	if quantity <= 0 {
		return Combination{}, errors.New("quantity must be > 0")
	}
	if len(packs) == 0 {
		return Combination{}, errors.New("at least one pack size is required")
	}

	// Normalize and remove duplicated packs
	sizeSet := make(map[int]struct{})
	var sizes []int
	maxPack := 0
	for _, p := range packs {
		if p.Size <= 0 {
			return Combination{}, errors.New("pack size must be > 0")
		}
		if _, ok := sizeSet[p.Size]; !ok {
			sizeSet[p.Size] = struct{}{}
			sizes = append(sizes, p.Size)
			if p.Size > maxPack {
				maxPack = p.Size
			}
		}
	}

	if len(sizes) == 0 {
		return Combination{}, errors.New("no valid pack sizes")
	}
	sort.Ints(sizes)

	// Otiumization: scale for GCD (reduce the DP size)
	g := gcdAll(sizes)
	qScaled := (quantity + g - 1) / g      // ceil(quantity/g)
	sizesScaled := make([]int, len(sizes)) // packs in “unit of g”
	for i, s := range sizes {
		sizesScaled[i] = s / g
	}
	maxPackScaled := sizesScaled[len(sizesScaled)-1]

	// DP scaled size:
	// dp[t] = minimum packs to sum exactly t
	// prev[t] = last pack used to get to t
	upper := qScaled + maxPackScaled - 1
	const inf = math.MaxInt32
	dp := make([]int, upper+1)
	prev := make([]int, upper+1)
	for i := range dp {
		dp[i], prev[i] = inf, -1
	}
	dp[0] = 0

	for t := 0; t <= upper; t++ {
		if dp[t] == inf {
			continue
		}
		for _, s := range sizesScaled {
			if nt := t + s; nt <= upper && dp[t]+1 < dp[nt] {
				dp[nt] = dp[t] + 1
				prev[nt] = s
			}
		}
	}

	// Get the less total total >= qScaled; in a tie, fewer packages
	bestTotalScaled, bestPacks := -1, inf
	for t := qScaled; t <= upper; t++ {
		if dp[t] == inf {
			continue
		}
		if bestTotalScaled == -1 || t < bestTotalScaled || (t == bestTotalScaled && dp[t] < bestPacks) {
			bestTotalScaled, bestPacks = t, dp[t]
		}
	}
	if bestTotalScaled == -1 {
		return Combination{}, errors.New("no feasible combination found")
	}

	// Reconstructs counts (to scale)
	countsScaled := make(map[int]int)
	for t := bestTotalScaled; t > 0; {
		s := prev[t]
		if s <= 0 {
			return Combination{}, errors.New("internal reconstruction error")
		}
		countsScaled[s]++
		t -= s
	}

	// “De-scale” to real values
	counts := make(map[int]int)
	for sScaled, c := range countsScaled {
		counts[sScaled*g] = c
	}
	totalItems := bestTotalScaled * g
	leftover := totalItems - quantity

	return Combination{
		ItemsByPack: counts,
		TotalItems:  totalItems,
		TotalPacks:  bestPacks,
		Leftover:    leftover,
	}, nil
}

// Helpers GCD (Euclides) greatest common divisor
func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}
func gcdAll(nums []int) int {
	g := nums[0]
	for _, n := range nums[1:] {
		g = gcd(g, n)
	}
	return g
}
