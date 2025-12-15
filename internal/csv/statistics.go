package csv

import (
	"fmt"
	"sort"
)

// calculateMedian computes the median value of a slice of integers
func calculateMedian(values []int) float64 {
	if len(values) == 0 {
		return 0
	}
	// Create a copy to avoid modifying the original slice
	sorted := make([]int, len(values))
	copy(sorted, values)
	sort.Ints(sorted)

	// Find median
	middle := len(sorted) / 2
	if len(sorted)%2 == 0 {
		// Even number of values, average the two middle values
		return float64(sorted[middle-1]+sorted[middle]) / 2.0
	}
	// Odd number of values, return the middle value
	return float64(sorted[middle])
}

// calculateMean computes the mean value of a slice of integers
func calculateMean(values []int) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0
	for _, v := range values {
		sum += v
	}
	return float64(sum) / float64(len(values))
}

func generateNetworkStatistics(networks map[string]*networkStats) []networkMedian {
	networkMedians := []networkMedian{}

	for key, stats := range networks {
		median := calculateMedian(stats.ConnectionTimes)
		mean := calculateMean(stats.ConnectionTimes)

		networkMedians = append(networkMedians, networkMedian{
			Address:     key, // Use the key as the address
			Location:    fmt.Sprintf("%s, %s", stats.City, stats.Country),
			Network:     stats.Network,
			MedianTime:  median,
			MinTime:     stats.MinTime,
			MaxTime:     stats.MaxTime,
			MeanTime:    mean,
			SampleCount: len(stats.ConnectionTimes),
		})
	}

	// Sort networks by median time (ascending)
	sort.Slice(networkMedians, func(i, j int) bool {
		return networkMedians[i].MedianTime < networkMedians[j].MedianTime
	})
	return networkMedians
}
