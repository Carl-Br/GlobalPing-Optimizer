package csv

import (
	"fmt"
)

func groupByCityAndNetwork(measurements []MeasurementData) map[string]*networkStats {
	networks := make(map[string]*networkStats)

	unfinishedResultsCount := 0
	for _, measurement := range measurements {
		for _, result := range measurement.Results {
			// Skip results that don't have timing data
			if result.Result.Status != "finished" {
				unfinishedResultsCount++
				continue
			}

			// Create a key based on city and network
			key := fmt.Sprintf("%s-%s", result.Probe.City, result.Probe.Network)

			// Create new Network entry if it doesn't exists
			if _, exists := networks[key]; !exists {
				networks[key] = &networkStats{
					Address:         key, // Use the key as the address
					City:            result.Probe.City,
					Country:         result.Probe.Country,
					Network:         result.Probe.Network,
					ConnectionTimes: []int{},
					MinTime:         999999, // Initialize with a high value
					MaxTime:         0,
				}
			}

			// Add connection time to Network stats
			totalTime := result.Result.Timings.Total
			networks[key].ConnectionTimes = append(
				networks[key].ConnectionTimes,
				totalTime,
			)

			// Update min/max
			if totalTime < networks[key].MinTime {
				networks[key].MinTime = totalTime
			}
			if totalTime > networks[key].MaxTime {
				networks[key].MaxTime = totalTime
			}
		}
	}
	fmt.Println("Number of unfinished results:", unfinishedResultsCount)
	return networks
}
