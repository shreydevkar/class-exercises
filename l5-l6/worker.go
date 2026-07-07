package main

import (
	"strconv"
	"strings"
)

// Map takes a filepath (key) and its contents (value) and emits a
// (year, temperature) KVPair for every reading found in the file.
// Each line looks like: "Moscow,2020-10-06,14.2"
func Map(key string, value string) []KVPair {

	output := make([]KVPair, 0)

	// Iterate over each line of the "value" string.
	for _, line := range strings.Split(strings.TrimSuffix(value, "\n"), "\n") {
		// Trim any trailing carriage return (Windows CRLF line endings).
		line = strings.TrimSpace(line)

		// Split "City,YYYY-MM-DD,temp" into its three fields.
		fields := strings.Split(line, ",")
		if len(fields) != 3 {
			continue
		}

		date := fields[1]
		temperature := fields[2]

		// The year is the part of the date before the first "-".
		year := strings.Split(date, "-")[0]

		output = append(output, KVPair{key: year, value: temperature})
	}

	return output
}

// Reduce takes a year (key) and all temperatures recorded that year (value)
// and returns the maximum temperature.
func Reduce(key string, value []string) float64 {
	max := 0.0
	for i, temp := range value {
		val, err := strconv.ParseFloat(temp, 64)
		if err != nil {
			continue
		}
		if i == 0 || val > max {
			max = val
		}
	}
	return max
}
